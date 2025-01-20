// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"sync"

	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/dustin/go-humanize"
	"github.com/go-resty/resty/v2"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

// MaxUploadParts is the maximum allowed number of parts in a multipart upload
// on Amazon S3.
const MaxUploadParts int32 = 10000

// MinUploadPartSize is the minimum allowed part size when uploading a part to
// Amazon S3.
const MinUploadPartSize int64 = 1024 * 1024 * 5

// MaxUploadPartSize is the maximum allowed part size when uploading a part to
// Amazon S3.
const MaxUploadPartSize int64 = 1024 * 1024 * 1024 * 5

// DefaultUploadPartSize is the default part size to buffer chunks of a
// payload into.
const DefaultUploadPartSize = 1024 * 1024 * 100

// DefaultUploadConcurrency is the default number of goroutines to spin up when
// using Upload().
const DefaultUploadConcurrency = 8

// MaxConcurrency is the maximum number of goroutines to upload parts in parallel.
const MaxConcurrency = 64

const processCompletePartWeights = 1

type UploadFailure interface {
	error

	UploadID() string
}

type uploadError struct {
	err error

	// ID for multipart upload which failed.
	uploadID string
}

// batchItemError returns the string representation of the error.
func (m *uploadError) Error() string {
	var extra string
	if m.err != nil {
		extra = fmt.Sprintf(", cause: %v", m.err)
	}
	return fmt.Sprintf("upload failed, upload id: %s%s", m.uploadID, extra)
}

// Unwrap returns the underlying error that cause the upload failure
func (m *uploadError) Unwrap() error {
	return m.err
}

// UploadID returns the id of the S3 upload which failed.
func (m *uploadError) UploadID() string {
	return m.uploadID
}

type PutObjectInput struct {
	FileName      *string
	DatabaseName  *string
	TableName     *string
	ContentLength *int64
	ClusterID     string
	Body          ReaderAtSeeker

	OnProgress func(ratio float64)
}

type Uploader interface {
	Upload(ctx context.Context, input *PutObjectInput) (string, error)
	SetConcurrency(concurrency int) error
}

// The UploaderImpl structure that calls Upload(). It is safe to call Upload()
// on this structure for multiple objects and across concurrent goroutines.
// Mutating the UploaderImpl's properties is not safe to be done concurrently.
type UploaderImpl struct {
	// The buffer size (in bytes) to use when buffering data into chunks and
	// sending them as parts to S3. The minimum allowed part size is 5MB, and
	// if this value is set to zero, the DefaultUploadPartSize value will be used.
	PartSize int64

	// The number of goroutines to spin up in parallel per call to Upload when
	// sending parts. If this is set to zero, the DefaultUploadConcurrency value
	// will be used.
	//
	// The concurrency pool is not shared between calls to Upload.
	Concurrency int

	// Setting this value to true will cause the SDK to avoid calling
	// CancelUpload on a failure, leaving all successfully uploaded
	// parts on S3 for manual recovery.
	//
	// Note that storing parts of an incomplete multipart upload counts towards
	// space usage on S3 and will add additional costs if not cleaned up.
	LeavePartsOnError bool

	// MaxUploadParts is the max number of parts which will be uploaded to S3.
	// Will be used to calculate the part size of the object to be uploaded.
	// E.g: 5GB file, with MaxUploadParts set to 100, will upload the file
	// as 100, 50MB parts. With a limited of s3.MaxUploadParts (10,000 parts).
	//
	// MaxUploadParts must not be used to limit the total number of bytes uploaded.
	// Use a type like to io.LimitReader (https://golang.org/pkg/io/#LimitedReader)
	// instead. An io.LimitReader is helpful when uploading an unbounded reader
	// to S3, and you know its maximum size. Otherwise, the reader's io.EOF returned
	// error must be used to signal end of stream.
	//
	// Defaults to package const's MaxUploadParts value.
	MaxUploadParts int32

	// Defines the buffer strategy used when uploading a part
	BufferProvider manager.ReadSeekerWriteToProvider

	// partPool allows for the re-usage of streaming payload part buffers between upload calls
	partPool byteSlicePool

	client     cloud.TiDBCloudClient
	httpClient *resty.Client
}

// NewUploader creates a new UploaderImpl instance to upload objects to S3. Pass In
// additional functional options to customize the uploader's behavior. Requires a
// cloud.TiDBCloudClient.
func NewUploader(client cloud.TiDBCloudClient) Uploader {
	httpClient := resty.New()
	debug := os.Getenv(config.DebugEnv) != ""
	httpClient.SetDebug(debug)
	u := &UploaderImpl{
		PartSize:          DefaultUploadPartSize,
		Concurrency:       DefaultUploadConcurrency,
		LeavePartsOnError: false,
		MaxUploadParts:    MaxUploadParts,
		BufferProvider:    defaultUploadBufferProvider(),
		client:            client,
		httpClient:        httpClient,
	}

	u.partPool = newByteSlicePool(u.PartSize)

	return u
}

// Upload uploads an object to S3, intelligently buffering large
// files into smaller chunks and sending them in parallel across multiple
// goroutines. You can configure the buffer size and concurrency through the
// UploaderImpl parameters.
// It is safe to call this method concurrently across goroutines.
func (u *UploaderImpl) Upload(ctx context.Context, input *PutObjectInput) (string, error) {
	i := uploader{cfg: u, ctx: ctx, in: input}

	return i.upload()
}

func (u *UploaderImpl) SetConcurrency(concurrency int) error {
	if concurrency > MaxConcurrency {
		return errors.New("concurrency must be less than 64")
	}
	if concurrency <= 0 {
		return errors.New("concurrency must be greater than 0")
	}
	u.Concurrency = concurrency
	return nil
}

// internal structure to manage an upload to S3.
type uploader struct {
	ctx context.Context
	cfg *UploaderImpl

	in *PutObjectInput

	readerPos int64 // current reader position
	totalSize int64 // set to -1 if the size is not known
}

// internal logic for deciding whether to upload a single part or use a
// multipart upload.
func (u *uploader) upload() (string, error) {
	if err := u.init(); err != nil {
		return "", fmt.Errorf("unable to initialize upload: %w", err)
	}
	defer u.cfg.partPool.Close()

	if u.cfg.PartSize < MinUploadPartSize {
		return "", fmt.Errorf("part size must be at least %s bytes", humanize.IBytes(uint64(MinUploadPartSize)))
	}
	if u.cfg.PartSize > MaxUploadPartSize {
		return "", fmt.Errorf("part size must be at most %s bytes", humanize.IBytes(uint64(MaxUploadPartSize)))
	}

	// Do one read to determine if we have more than one part
	reader, _, cleanup, err := u.nextReader()
	if err == io.EOF { // single part
		sg := singerUploader{uploader: u}
		return sg.upload(reader, cleanup)
	} else if err != nil {
		cleanup()
		return "", fmt.Errorf("read upload data failed: %w", err)
	}

	mu := multiUploader{uploader: u}
	return mu.upload(reader, cleanup)
}

// init will initialize all default options.
func (u *uploader) init() error {
	if u.cfg.Concurrency == 0 {
		u.cfg.Concurrency = DefaultUploadConcurrency
	}
	if u.cfg.PartSize == 0 {
		u.cfg.PartSize = DefaultUploadPartSize
	}
	if u.cfg.MaxUploadParts == 0 {
		u.cfg.MaxUploadParts = MaxUploadParts
	}

	// Try to get the total size for some optimizations
	u.initSize()

	// If PartSize was changed or partPool was never setup then we need to allocate a new pool
	// so that we return []byte slices of the correct size
	poolCap := u.cfg.Concurrency + 1
	if u.cfg.partPool == nil || u.cfg.partPool.SliceSize() != u.cfg.PartSize {
		u.cfg.partPool = newByteSlicePool(u.cfg.PartSize)
		u.cfg.partPool.ModifyCapacity(poolCap)
	} else {
		u.cfg.partPool = &returnCapacityPoolCloser{byteSlicePool: u.cfg.partPool}
		u.cfg.partPool.ModifyCapacity(poolCap)
	}

	return nil
}

// initSize tries to detect the total stream size, setting u.totalSize. If
// the size is not known, totalSize is set to -1.
func (u *uploader) initSize() {
	u.totalSize = *u.in.ContentLength
	// Try to adjust partSize if it is too small and account for
	// integer division truncation.
	if u.totalSize/u.cfg.PartSize >= int64(u.cfg.MaxUploadParts) {
		// Add one to the part size to account for remainders
		// during the size calculation. e.g. odd number of bytes.
		u.cfg.PartSize = (u.totalSize / int64(u.cfg.MaxUploadParts)) + 1
	}
}

// nextReader returns a seekable reader representing the next packet of data.
// This operation increases the shared u.readerPos counter, but note that it
// does not need to be wrapped in a mutex because nextReader is only called
// from the main thread.
func (u *uploader) nextReader() (io.ReadSeeker, int, func(), error) {
	var err error
	n := u.cfg.PartSize
	if u.totalSize >= 0 {
		bytesLeft := u.totalSize - u.readerPos

		if bytesLeft <= u.cfg.PartSize {
			err = io.EOF
			n = bytesLeft
		}
	}

	var (
		reader  io.ReadSeeker
		cleanup func()
	)

	reader = io.NewSectionReader(u.in.Body, u.readerPos, n)
	if u.cfg.BufferProvider != nil {
		reader, cleanup = u.cfg.BufferProvider.GetWriteTo(reader)
	} else {
		cleanup = func() {}
	}

	u.readerPos += n

	return reader, int(n), cleanup, err
}

// internal structure to manage a specific multipart upload to S3.
type singerUploader struct {
	*uploader
	uploadID string
}

// upload contains upload logic for uploading a single chunk via
// a regular PutObject request. Multipart requests require at least two
// parts, or at least 5MB of data.
func (u *singerUploader) upload(r io.ReadSeeker, cleanup func()) (string, error) {
	defer cleanup()
	res, err := u.cfg.client.StartUpload(u.ctx, u.in.ClusterID, u.in.FileName,
		u.in.DatabaseName, u.in.TableName, aws.Int32(1))
	if err != nil {
		return "", err
	}
	u.uploadID = *res.UploadId

	resp, err := u.cfg.httpClient.R().
		SetContext(u.ctx).
		SetContentLength(true).
		SetBody(r).Put(res.UploadUrl[0])

	if err != nil {
		u.fail()
		return "", &uploadError{
			err:      err,
			uploadID: *res.UploadId,
		}
	}

	if !resp.IsSuccess() {
		u.fail()
		return "", &uploadError{
			err:      fmt.Errorf("upload failed, code: %s, reason: %s", resp.Status(), string(resp.Body())),
			uploadID: *res.UploadId,
		}
	}

	err = u.complete()
	if err != nil {
		return "", err
	}

	return *res.UploadId, nil
}

func (u *singerUploader) complete() error {
	var parts []imp.CompletePart
	err := u.cfg.client.CompleteUpload(u.ctx, u.in.ClusterID, &u.uploadID, &parts)
	if err != nil {
		u.fail()
		return err
	}
	return nil
}

func (u *singerUploader) fail() {
	ctx := context.WithoutCancel(u.ctx)
	err := u.cfg.client.CancelUpload(ctx, u.in.ClusterID, &u.uploadID)
	if err != nil {
		log.Warn("failed to abort singlePart upload", zap.Error(err))
		return
	}
}

// internal structure to manage a specific multipart upload to S3.
type multiUploader struct {
	*uploader
	wg       sync.WaitGroup
	m        sync.Mutex
	err      error
	uploadID string
	//parts    completedParts
	parts completedParts
	urls  []string
}

// keeps track of a single chunk of data being sent to S3.
type chunk struct {
	buf     io.ReadSeeker
	num     int32
	cleanup func()
}

// completedParts is a wrapper to make parts sortable by their part number,
// since S3 required this list to be sent in sorted order.
type completedParts []imp.CompletePart

func (a completedParts) Len() int      { return len(a) }
func (a completedParts) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a completedParts) Less(i, j int) bool {
	return a[i].PartNumber < a[j].PartNumber
}

// upload will perform a multipart upload using the firstBuf buffer containing
// the first chunk of data.
func (u *multiUploader) upload(firstBuf io.ReadSeeker, cleanup func()) (string, error) {
	partNumber := math.Ceil(float64(u.totalSize) / float64(u.cfg.PartSize))
	url, err := u.cfg.client.StartUpload(u.ctx, u.in.ClusterID, u.in.FileName, u.in.DatabaseName, u.in.TableName, aws.Int32(int32(partNumber)))
	if err != nil {
		cleanup()
		return "", err
	}
	u.uploadID = *url.UploadId
	u.urls = url.UploadUrl

	// Create the workers
	ch := make(chan chunk, u.cfg.Concurrency)
	for i := 0; i < u.cfg.Concurrency; i++ {
		u.wg.Add(1)
		go u.readChunk(ch)
	}

	// Send part 1 to the workers
	var num int32 = 1
	ch <- chunk{buf: firstBuf, num: num, cleanup: cleanup}

	// Read and queue the rest of the parts
	for u.getErr() == nil && err == nil {
		var (
			reader       io.ReadSeeker
			nextChunkLen int
			ok           bool
		)

		reader, nextChunkLen, cleanup, err = u.nextReader()
		ok, err = u.shouldContinue(num, nextChunkLen, err)
		if !ok {
			cleanup()
			if err != nil {
				u.setErr(err)
			}
			break
		}

		num++

		ch <- chunk{buf: reader, num: num, cleanup: cleanup}
	}

	// Close the channel, wait for workers, and complete upload
	close(ch)
	u.wg.Wait()
	u.complete()

	if err := u.getErr(); err != nil {
		return "", &uploadError{
			err:      err,
			uploadID: u.uploadID,
		}
	}
	return u.uploadID, nil
}

func (u *multiUploader) shouldContinue(part int32, nextChunkLen int, err error) (bool, error) {
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("read multipart upload data failed, %w", err)
	}

	if nextChunkLen == 0 {
		// No need to upload empty part, if file was empty to start
		// with empty single part would have been created and never
		// started multipart upload.
		return false, nil
	}

	part++
	// This upload exceeded maximum number of supported parts, error now.
	if part > u.cfg.MaxUploadParts || part > MaxUploadParts {
		var msg string
		if part > u.cfg.MaxUploadParts {
			msg = fmt.Sprintf("exceeded total allowed configured MaxUploadParts (%d). Adjust PartSize to fit in this limit",
				u.cfg.MaxUploadParts)
		} else {
			msg = fmt.Sprintf("exceeded total allowed S3 limit MaxUploadParts (%d). Adjust PartSize to fit in this limit",
				MaxUploadParts)
		}
		return false, fmt.Errorf(msg)
	}

	return true, err
}

// readChunk runs in worker goroutines to pull chunks off of the ch channel
// and send() them as UploadPart requests.
func (u *multiUploader) readChunk(ch chan chunk) {
	defer u.wg.Done()
	for {
		data, ok := <-ch

		if !ok {
			break
		}

		if u.getErr() == nil {
			if err := u.send(data); err != nil {
				u.setErr(err)
			}
		}

		data.cleanup()
	}
}

// send performs an UploadPart request and keeps track of the completed
// part information.
func (u *multiUploader) send(c chunk) error {
	resp, err := u.cfg.httpClient.R().
		SetContext(u.ctx).
		SetContentLength(true).
		SetBody(c.buf).Put(u.urls[c.num-1])

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return util.InterruptError
		}
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("upload failed, code: %s, reason: %s", resp.Status(), string(resp.Body()))
	}

	var completed imp.CompletePart
	completed.PartNumber = c.num
	etag := resp.Header().Get("ETag")
	completed.Etag = etag

	u.m.Lock()
	u.parts = append(u.parts, completed)
	u.in.OnProgress(float64(len(u.parts)) / float64(len(u.urls)+processCompletePartWeights))
	u.m.Unlock()

	return nil
}

// getErr is a thread-safe getter for the error object
func (u *multiUploader) getErr() error {
	u.m.Lock()
	defer u.m.Unlock()

	return u.err
}

// setErr is a thread-safe setter for the error object
func (u *multiUploader) setErr(e error) {
	u.m.Lock()
	defer u.m.Unlock()

	u.err = e
}

// fail will abort the multipart unless LeavePartsOnError is set to true.
func (u *multiUploader) fail() {
	if u.cfg.LeavePartsOnError {
		return
	}

	ctx := context.WithoutCancel(u.ctx)
	err := u.cfg.client.CancelUpload(ctx, u.in.ClusterID, &u.uploadID)
	if err != nil {
		log.Warn("failed to abort multipart upload", zap.Error(err))
		return
	}
}

// complete successfully completes a multipart upload and returns the response.
func (u *multiUploader) complete() {
	if u.getErr() != nil {
		u.fail()
		return
	}

	sort.Sort(u.parts)
	cp := []imp.CompletePart(u.parts)
	err := u.cfg.client.CompleteUpload(u.ctx, u.in.ClusterID, &u.uploadID, &cp)

	if err != nil {
		u.setErr(err)
		u.fail()
	}
}

type ReaderAtSeeker interface {
	io.ReaderAt
	io.ReadSeeker
}
