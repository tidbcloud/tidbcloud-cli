// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fs provides the filesystem client for TiDB Cloud FS operations.
package fs

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// PartSize is the default S3 multipart part size (8 MiB).
	PartSize = 8 * 1024 * 1024

	// DefaultSmallFileThreshold matches the server's threshold for direct PUT vs multipart.
	DefaultSmallFileThreshold = 50_000 // 50,000 bytes

	// Upload concurrency limits
	uploadMaxConcurrency = 16
	uploadMaxBufferBytes = 256 * 1024 * 1024 // 256 MB
)

// Part represents a single part in a multipart upload.
type Part struct {
	Number int64
	Size   int64
}

// CalcParts calculates the parts needed for a file of given size.
func CalcParts(totalSize, partSize int64) []Part {
	if totalSize <= 0 {
		return nil
	}
	if partSize <= 0 {
		partSize = PartSize
	}
	parts := make([]Part, 0, (totalSize+partSize-1)/partSize)
	var offset int64
	for offset < totalSize {
		sz := partSize
		if offset+sz > totalSize {
			sz = totalSize - offset
		}
		parts = append(parts, Part{
			Number: int64(len(parts) + 1),
			Size:   sz,
		})
		offset += sz
	}
	return parts
}

// Client is the FS HTTP client.
type Client struct {
	baseURL            string
	httpClient         *http.Client
	clusterID          string
	zeroInstanceID     string
	smallFileThreshold int64 // 0 means use DefaultSmallFileThreshold
}

// GetBaseURL returns the base URL of the client.
func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// GetClusterID returns the associated TiDB Cloud cluster ID.
func (c *Client) GetClusterID() string {
	return c.clusterID
}

// GetZeroInstanceID returns the associated TiDB Zero instance ID.
func (c *Client) GetZeroInstanceID() string {
	return c.zeroInstanceID
}

// ProvisionResult is the response from the FS provision endpoint.
type ProvisionResult struct {
	TenantID string `json:"tenant_id"`
	APIKey   string `json:"api_key"`
	Status   string `json:"status"`
}

// Provision initializes the FS tenant for the associated database.
func (c *Client) Provision(ctx context.Context, user, password string) (*ProvisionResult, error) {
	body, err := json.Marshal(map[string]string{
		"user":     user,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/provision", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusAccepted {
		return nil, readError(resp)
	}
	var result ProvisionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode provision response: %w", err)
	}
	return &result, nil
}

// NewClient creates a new FS client.
// The httpClient should already be configured with authentication (e.g., Bearer or Digest transport).
// clusterID and zeroInstanceID are sent as X-TIDBCLOUD-CLUSTER-ID and X-TIDBCLOUD-ZERO-INSTANCE-ID headers.
func NewClient(baseURL string, httpClient *http.Client, clusterID, zeroInstanceID string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{
		baseURL:        strings.TrimRight(baseURL, "/"),
		httpClient:     httpClient,
		clusterID:      clusterID,
		zeroInstanceID: zeroInstanceID,
	}
}

// FileInfo represents a file entry from a directory listing.
type FileInfo struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isDir"`
}

// StatResult represents file metadata from HEAD.
type StatResult struct {
	Size     int64
	IsDir    bool
	Revision int64
}

func (c *Client) url(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return c.baseURL + "/v1/fs" + path
}

// RawPost sends a raw POST request to the specified endpoint.
func (c *Client) RawPost(endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, c.baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.doWithClient(c.httpClient, req)
}

func (c *Client) doWithClient(client *http.Client, req *http.Request) (*http.Response, error) {
	// Inject TiDB Cloud native provider headers. Cluster ID takes precedence over Zero instance ID.
	if c.clusterID != "" {
		req.Header.Set("X-TIDBCLOUD-CLUSTER-ID", c.clusterID)
	} else if c.zeroInstanceID != "" {
		req.Header.Set("X-TIDBCLOUD-ZERO-INSTANCE-ID", c.zeroInstanceID)
	}
	return client.Do(req)
}

// Write uploads data to a remote path.
func (c *Client) Write(path string, data []byte) error {
	req, err := http.NewRequest(http.MethodPut, c.url(path), bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return readError(resp)
	}
	return nil
}

// Read downloads a file's content.
func (c *Client) Read(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.url(path), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	return io.ReadAll(resp.Body)
}

// List returns the entries in a directory.
func (c *Client) List(path string) ([]FileInfo, error) {
	// Use an explicit value to avoid intermediaries dropping bare "?list".
	req, err := http.NewRequest(http.MethodGet, c.url(path)+"?list=1", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var result struct {
		Entries []FileInfo `json:"entries"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return result.Entries, nil
}

// Stat returns metadata for a path.
func (c *Client) Stat(path string) (*StatResult, error) {
	req, err := http.NewRequest(http.MethodHead, c.url(path), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("not found: %s", path)
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	s := &StatResult{
		IsDir: resp.Header.Get("X-FS-IsDir") == "true",
	}
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		s.Size, _ = strconv.ParseInt(cl, 10, 64)
	}
	if rev := resp.Header.Get("X-FS-Revision"); rev != "" {
		s.Revision, _ = strconv.ParseInt(rev, 10, 64)
	}
	return s, nil
}

// Delete removes a file or directory.
func (c *Client) Delete(path string) error {
	req, err := http.NewRequest(http.MethodDelete, c.url(path), nil)
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return readError(resp)
	}
	return nil
}

// Copy performs a server-side zero-copy (same file_id, new path).
func (c *Client) Copy(srcPath, dstPath string) error {
	req, err := http.NewRequest(http.MethodPost, c.url(dstPath)+"?copy", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Dat9-Copy-Source", srcPath)
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return readError(resp)
	}
	return nil
}

// Rename moves/renames a file or directory (metadata-only).
func (c *Client) Rename(oldPath, newPath string) error {
	req, err := http.NewRequest(http.MethodPost, c.url(newPath)+"?rename", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Dat9-Rename-Source", oldPath)
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return readError(resp)
	}
	return nil
}

// Mkdir creates a directory.
func (c *Client) Mkdir(path string) error {
	req, err := http.NewRequest(http.MethodPost, c.url(path)+"?mkdir", nil)
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return readError(resp)
	}
	return nil
}

func readError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	var errResp struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
		return fmt.Errorf("%s", errResp.Error)
	}
	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
}

// SQL executes a SQL query.
func (c *Client) SQL(query string) ([]map[string]interface{}, error) {
	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/v1/sql", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var rows []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return rows, nil
}

// SearchResult represents a search result.
type SearchResult struct {
	Path      string   `json:"path"`
	Name      string   `json:"name"`
	SizeBytes int64    `json:"size_bytes"`
	Score     *float64 `json:"score,omitempty"`
}

// Grep searches for a pattern in files.
func (c *Client) Grep(query, pathPrefix string, limit int) ([]SearchResult, error) {
	u := c.url(pathPrefix) + "?grep=" + url.QueryEscape(query)
	if limit > 0 {
		u += "&limit=" + strconv.Itoa(limit)
	}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var results []SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// Find searches for files matching criteria.
func (c *Client) Find(pathPrefix string, params url.Values) ([]SearchResult, error) {
	params.Set("find", "")
	u := c.url(pathPrefix) + "?" + params.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var results []SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// ProgressFunc is called after each part upload completes.
type ProgressFunc func(partNumber, totalParts int, bytesUploaded int64)

// UploadPlan is the server's 202 response for large file uploads.
type UploadPlan struct {
	UploadID string    `json:"upload_id"`
	PartSize int64     `json:"part_size"`
	Parts    []PartURL `json:"parts"`
}

// PartURL is a presigned URL for uploading one part.
type PartURL struct {
	Number         int               `json:"number"`
	URL            string            `json:"url"`
	Size           int64             `json:"size"`
	ChecksumSHA256 string            `json:"checksum_sha256,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	ExpiresAt      string            `json:"expires_at"`
}

// WriteStream uploads data from a reader. For small files (size < threshold),
// it does a direct PUT with body. For large files, it sends a Content-Length-only
// PUT to get a 202 with presigned URLs, then uploads parts concurrently.
func (c *Client) WriteStream(ctx context.Context, path string, r io.Reader, size int64, progress ProgressFunc) error {
	threshold := int64(DefaultSmallFileThreshold)
	if c.smallFileThreshold > 0 {
		threshold = c.smallFileThreshold
	}
	if size < threshold {
		// Small file: direct PUT with body
		data, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("read data: %w", err)
		}
		return c.Write(path, data)
	}
	ra, ok := r.(io.ReaderAt)
	if !ok {
		return fmt.Errorf("large uploads require an io.ReaderAt (seekable source) to compute per-part checksums")
	}
	checksums, err := computePartChecksumsFromReaderAt(ra, size, PartSize)
	if err != nil {
		return fmt.Errorf("compute part checksums: %w", err)
	}
	plan, err := c.initiateUpload(ctx, path, size, checksums)
	if err != nil {
		return err
	}
	return c.uploadParts(ctx, plan, ra, progress)
}

type uploadInitiateRequest struct {
	Path          string   `json:"path"`
	TotalSize     int64    `json:"total_size"`
	PartChecksums []string `json:"part_checksums"`
}

type uploadResumeRequest struct {
	PartChecksums []string `json:"part_checksums"`
}

func (c *Client) initiateUpload(ctx context.Context, path string, size int64, checksums []string) (UploadPlan, error) {
	plan, resp, err := c.initiateUploadByBody(ctx, path, size, checksums)
	if err == nil {
		return plan, nil
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
			return c.initiateUploadLegacy(ctx, path, size, checksums)
		}
		if resp.StatusCode == http.StatusBadRequest && strings.Contains(strings.ToLower(err.Error()), "unknown upload action") {
			return c.initiateUploadLegacy(ctx, path, size, checksums)
		}
		return UploadPlan{}, err
	}
	return UploadPlan{}, err
}

func (c *Client) initiateUploadByBody(ctx context.Context, path string, size int64, checksums []string) (UploadPlan, *http.Response, error) {
	body, err := json.Marshal(uploadInitiateRequest{Path: path, TotalSize: size, PartChecksums: checksums})
	if err != nil {
		return UploadPlan{}, nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/uploads/initiate", bytes.NewReader(body))
	if err != nil {
		return UploadPlan{}, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	if err != nil {
		return UploadPlan{}, nil, err
	}
	if resp.StatusCode != http.StatusAccepted {
		return UploadPlan{}, resp, readError(resp)
	}
	var plan UploadPlan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		_ = resp.Body.Close()
		return UploadPlan{}, nil, fmt.Errorf("decode upload plan: %w", err)
	}
	_ = resp.Body.Close()
	return plan, nil, nil
}

func (c *Client) initiateUploadLegacy(ctx context.Context, path string, size int64, checksums []string) (UploadPlan, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.url(path), http.NoBody)
	if err != nil {
		return UploadPlan{}, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-FS-Content-Length", fmt.Sprintf("%d", size))
	if len(checksums) > 0 {
		req.Header.Set("X-FS-Part-Checksums", strings.Join(checksums, ","))
	}

	resp, err := c.do(req)
	if err != nil {
		return UploadPlan{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		return UploadPlan{}, readError(resp)
	}

	var plan UploadPlan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		return UploadPlan{}, fmt.Errorf("decode upload plan: %w", err)
	}
	return plan, nil
}

func uploadParallelism(partSize int64) int {
	if partSize <= 0 {
		partSize = PartSize
	}
	byMemory := int(uploadMaxBufferBytes / partSize)
	if byMemory < 1 {
		byMemory = 1
	}
	return min(byMemory, uploadMaxConcurrency)
}

func checksumParallelism(partSize int64, partCount int) int {
	if partSize <= 0 {
		partSize = PartSize
	}
	byMemory := int(uploadMaxBufferBytes / partSize)
	if byMemory < 1 {
		byMemory = 1
	}
	return min(runtime.GOMAXPROCS(0), partCount, byMemory)
}

func (c *Client) uploadParts(ctx context.Context, plan UploadPlan, ra io.ReaderAt, progress ProgressFunc) error {
	stdPartSize := plan.PartSize
	if stdPartSize <= 0 && len(plan.Parts) > 0 {
		stdPartSize = plan.Parts[0].Size
	}
	if stdPartSize <= 0 {
		stdPartSize = PartSize
	}
	maxConcurrency := uploadParallelism(stdPartSize)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for _, part := range plan.Parts {
		select {
		case err := <-errCh:
			cancel()
			wg.Wait()
			return err
		default:
		}

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Wait()
			return ctx.Err()
		}

		wg.Add(1)
		go func(p PartURL) {
			defer wg.Done()
			defer func() { <-sem }()

			data := make([]byte, p.Size)
			offset := int64(p.Number-1) * stdPartSize
			n, err := ra.ReadAt(data, offset)
			if err != nil && err != io.EOF {
				select {
				case errCh <- fmt.Errorf("read part %d: %w", p.Number, err):
				default:
				}
				cancel()
				return
			}
			if int64(n) != p.Size {
				select {
				case errCh <- fmt.Errorf("short read for part %d: got %d want %d", p.Number, n, p.Size):
				default:
				}
				cancel()
				return
			}

			_, err = c.uploadOnePart(ctx, p, data)
			if err != nil {
				select {
				case errCh <- fmt.Errorf("part %d: %w", p.Number, err):
				default:
				}
				cancel()
				return
			}

			if progress != nil {
				progress(p.Number, len(plan.Parts), int64(len(data)))
			}
		}(part)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
	}

	return c.completeUpload(ctx, plan.UploadID)
}

func (c *Client) uploadOnePart(ctx context.Context, part PartURL, data []byte) (string, error) {
	checksum := part.ChecksumSHA256
	if checksum == "" {
		h := sha256.Sum256(data)
		checksum = base64.StdEncoding.EncodeToString(h[:])
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, part.URL, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	for k, v := range part.Headers {
		if strings.EqualFold(k, "host") {
			continue
		}
		req.Header.Set(k, v)
	}
	req.ContentLength = int64(len(data))
	req.Header.Set("x-amz-checksum-sha256", checksum)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return resp.Header.Get("ETag"), nil
}

func (c *Client) completeUpload(ctx context.Context, uploadID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/v1/uploads/"+uploadID+"/complete", nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		return readError(resp)
	}
	return nil
}

// ReadStream reads a file, following 302 redirects for large files.
func (c *Client) ReadStream(ctx context.Context, path string) (io.ReadCloser, error) {
	noRedirectClient := *c.httpClient
	noRedirectClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(path), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.doWithClient(&noRedirectClient, req)
	if err != nil {
		return nil, err
	}

	switch {
	case resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusTemporaryRedirect:
		_ = resp.Body.Close()
		location := resp.Header.Get("Location")
		if location == "" {
			return nil, fmt.Errorf("302 without Location header")
		}
		req2, err := http.NewRequestWithContext(ctx, http.MethodGet, location, nil)
		if err != nil {
			return nil, err
		}
		resp2, err := c.httpClient.Do(req2)
		if err != nil {
			return nil, err
		}
		if resp2.StatusCode >= 300 {
			defer func() { _ = resp2.Body.Close() }()
			return nil, readError(resp2)
		}
		return resp2.Body, nil

	case resp.StatusCode >= 300:
		defer func() { _ = resp.Body.Close() }()
		return nil, readError(resp)

	default:
		return resp.Body, nil
	}
}

// ReadStreamRange reads a byte range from a remote file.
func (c *Client) ReadStreamRange(ctx context.Context, path string, offset, length int64) (io.ReadCloser, error) {
	if length <= 0 {
		return io.NopCloser(bytes.NewReader(nil)), nil
	}

	noRedirectClient := *c.httpClient
	noRedirectClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(path), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.doWithClient(&noRedirectClient, req)
	if err != nil {
		return nil, err
	}

	switch {
	case resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusTemporaryRedirect:
		_ = resp.Body.Close()
		location := resp.Header.Get("Location")
		if location == "" {
			return nil, fmt.Errorf("302 without Location header")
		}
		req2, err := http.NewRequestWithContext(ctx, http.MethodGet, location, nil)
		if err != nil {
			return nil, err
		}
		req2.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, offset+length-1))
		resp2, err := c.httpClient.Do(req2)
		if err != nil {
			return nil, err
		}

		switch resp2.StatusCode {
		case http.StatusPartialContent:
			return resp2.Body, nil
		case http.StatusRequestedRangeNotSatisfiable:
			defer func() { _ = resp2.Body.Close() }()
			return io.NopCloser(bytes.NewReader(nil)), nil
		default:
			if resp2.StatusCode >= 300 {
				defer func() { _ = resp2.Body.Close() }()
				return nil, readError(resp2)
			}
			return c.sliceBody(resp2.Body, offset, length)
		}

	case resp.StatusCode >= 300:
		defer func() { _ = resp.Body.Close() }()
		return nil, readError(resp)

	default:
		return c.sliceBody(resp.Body, offset, length)
	}
}

func (c *Client) sliceBody(rc io.ReadCloser, offset, length int64) (io.ReadCloser, error) {
	if offset > 0 {
		if _, err := io.CopyN(io.Discard, rc, offset); err != nil {
			_ = rc.Close()
			if err == io.EOF {
				return io.NopCloser(strings.NewReader("")), nil
			}
			return nil, fmt.Errorf("skip to offset: %w", err)
		}
	}
	return &limitedReadCloser{r: io.LimitReader(rc, length), c: rc}, nil
}

type limitedReadCloser struct {
	r io.Reader
	c io.Closer
}

func (l *limitedReadCloser) Read(p []byte) (int, error) { return l.r.Read(p) }
func (l *limitedReadCloser) Close() error               { return l.c.Close() }

// UploadMeta is the server's response for querying active uploads.
type UploadMeta struct {
	UploadID   string `json:"upload_id"`
	PartsTotal int    `json:"parts_total"`
	Status     string `json:"status"`
	ExpiresAt  string `json:"expires_at"`
}

// ResumeUpload queries for an incomplete upload and resumes it.
func (c *Client) ResumeUpload(ctx context.Context, path string, r io.ReaderAt, totalSize int64, progress ProgressFunc) error {
	meta, err := c.queryUpload(ctx, path)
	if err != nil {
		return err
	}

	checksums, err := computePartChecksumsFromReaderAt(r, totalSize, PartSize)
	if err != nil {
		return fmt.Errorf("compute part checksums: %w", err)
	}
	plan, err := c.requestResume(ctx, meta.UploadID, checksums)
	if err != nil {
		return err
	}

	if len(plan.Parts) == 0 {
		return c.completeUpload(ctx, plan.UploadID)
	}

	if err := c.uploadMissingParts(ctx, *plan, r, meta.PartsTotal, progress); err != nil {
		return err
	}

	return c.completeUpload(ctx, plan.UploadID)
}

func (c *Client) queryUpload(ctx context.Context, path string) (*UploadMeta, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/v1/uploads?path="+path+"&status=UPLOADING", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}

	var envelope struct {
		Uploads []UploadMeta `json:"uploads"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode upload meta: %w", err)
	}
	if len(envelope.Uploads) == 0 {
		return nil, fmt.Errorf("no active upload for %s", path)
	}
	return &envelope.Uploads[0], nil
}

func (c *Client) requestResume(ctx context.Context, uploadID string, checksums []string) (*UploadPlan, error) {
	plan, resp, err := c.requestResumeByBody(ctx, uploadID, checksums)
	if err == nil {
		return plan, nil
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode == http.StatusBadRequest && strings.Contains(strings.ToLower(err.Error()), "missing x-fs-part-checksums header") {
			return c.requestResumeLegacy(ctx, uploadID, checksums)
		}
		return nil, err
	}
	return nil, err
}

func (c *Client) requestResumeByBody(ctx context.Context, uploadID string, checksums []string) (*UploadPlan, *http.Response, error) {
	body, err := json.Marshal(uploadResumeRequest{PartChecksums: checksums})
	if err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/v1/uploads/"+uploadID+"/resume", bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode == http.StatusGone {
		_ = resp.Body.Close()
		return nil, nil, fmt.Errorf("upload %s has expired", uploadID)
	}
	if resp.StatusCode >= 300 {
		return nil, resp, readError(resp)
	}

	var plan UploadPlan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		_ = resp.Body.Close()
		return nil, nil, fmt.Errorf("decode resume plan: %w", err)
	}
	_ = resp.Body.Close()
	return &plan, nil, nil
}

func (c *Client) requestResumeLegacy(ctx context.Context, uploadID string, checksums []string) (*UploadPlan, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/v1/uploads/"+uploadID+"/resume", nil)
	if err != nil {
		return nil, err
	}
	if len(checksums) > 0 {
		req.Header.Set("X-FS-Part-Checksums", strings.Join(checksums, ","))
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusGone {
		return nil, fmt.Errorf("upload %s has expired", uploadID)
	}
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}

	var plan UploadPlan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		return nil, fmt.Errorf("decode resume plan: %w", err)
	}
	return &plan, nil
}

func (c *Client) uploadMissingParts(ctx context.Context, plan UploadPlan, r io.ReaderAt, totalParts int, progress ProgressFunc) error {
	stdPartSize := plan.PartSize
	if stdPartSize <= 0 {
		stdPartSize = PartSize
	}
	maxConcurrency := uploadParallelism(stdPartSize)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sem := make(chan struct{}, maxConcurrency)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	for _, part := range plan.Parts {
		select {
		case err := <-errCh:
			cancel()
			wg.Wait()
			return err
		default:
		}

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Wait()
			return ctx.Err()
		}

		wg.Add(1)
		go func(p PartURL) {
			defer wg.Done()
			defer func() { <-sem }()

			data := make([]byte, p.Size)
			offset := int64(p.Number-1) * stdPartSize
			n, err := r.ReadAt(data, offset)
			if err != nil && err != io.EOF {
				select {
				case errCh <- fmt.Errorf("read part %d at offset %d: %w", p.Number, offset, err):
				default:
				}
				cancel()
				return
			}
			if int64(n) != p.Size {
				select {
				case errCh <- fmt.Errorf("short read for part %d at offset %d: got %d want %d", p.Number, offset, n, p.Size):
				default:
				}
				cancel()
				return
			}

			_, err = c.uploadOnePart(ctx, p, data)
			if err != nil {
				select {
				case errCh <- fmt.Errorf("part %d: %w", p.Number, err):
				default:
				}
				cancel()
				return
			}
			if progress != nil {
				progress(p.Number, totalParts, int64(len(data)))
			}
		}(part)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
	}

	return nil
}

func computePartChecksumsFromReaderAt(r io.ReaderAt, totalSize int64, partSize int64) ([]string, error) {
	if totalSize <= 0 {
		return nil, nil
	}
	parts := CalcParts(totalSize, partSize)
	checksums := make([]string, len(parts))
	workers := checksumParallelism(partSize, len(parts))

	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once
	partCh := make(chan int, len(parts))

	for i := range parts {
		partCh <- i
	}
	close(partCh)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, partSize)
			for i := range partCh {
				p := parts[i]
				data := buf[:p.Size]
				offset := int64(p.Number-1) * partSize
				n, err := r.ReadAt(data, offset)
				if err != nil && err != io.EOF {
					errOnce.Do(func() { firstErr = fmt.Errorf("read part %d: %w", p.Number, err) })
					return
				}
				if int64(n) != p.Size {
					errOnce.Do(func() {
						firstErr = fmt.Errorf("short read for part %d: got %d want %d", p.Number, n, p.Size)
					})
					return
				}
				h := sha256.Sum256(data)
				checksums[i] = base64.StdEncoding.EncodeToString(h[:])
			}
		}()
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	return checksums, nil
}

// PatchPlan mirrors the server's response for a PATCH request.
type PatchPlan struct {
	UploadID    string          `json:"upload_id"`
	PartSize    int64           `json:"part_size"`
	UploadParts []*PatchPartURL `json:"upload_parts"`
	CopiedParts []int           `json:"copied_parts"`
}

// PatchPartURL describes one dirty part the client must upload.
type PatchPartURL struct {
	Number      int               `json:"number"`
	URL         string            `json:"url"`
	Size        int64             `json:"size"`
	Headers     map[string]string `json:"headers,omitempty"`
	ExpiresAt   string            `json:"expires_at"`
	ReadURL     string            `json:"read_url,omitempty"`
	ReadHeaders map[string]string `json:"read_headers,omitempty"`
}

// PatchFile performs a partial update of a large file.
func (c *Client) PatchFile(ctx context.Context, path string, newSize int64, dirtyParts []int, readPart func(partNumber int, partSize int64, origData []byte) ([]byte, error), progress ProgressFunc) error {
	reqBody, err := json.Marshal(map[string]any{
		"new_size":    newSize,
		"dirty_parts": dirtyParts,
	})
	if err != nil {
		return fmt.Errorf("marshal patch request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.url(path), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusAccepted {
		return readError(resp)
	}

	var plan PatchPlan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		return fmt.Errorf("decode patch plan: %w", err)
	}

	const maxConcurrency = 4
	sem := make(chan struct{}, maxConcurrency)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	totalParts := len(plan.UploadParts) + len(plan.CopiedParts)

	for _, part := range plan.UploadParts {
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Wait()
			return ctx.Err()
		}

		select {
		case err := <-errCh:
			wg.Wait()
			return err
		default:
		}

		wg.Add(1)
		go func(p *PatchPartURL) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := c.uploadPatchPart(ctx, p, readPart); err != nil {
				select {
				case errCh <- fmt.Errorf("part %d: %w", p.Number, err):
				default:
				}
				return
			}

			if progress != nil {
				progress(p.Number, totalParts, p.Size)
			}
		}(part)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
	}

	return c.completeUpload(ctx, plan.UploadID)
}

func (c *Client) uploadPatchPart(ctx context.Context, part *PatchPartURL, readPart func(int, int64, []byte) ([]byte, error)) error {
	var origData []byte
	if part.ReadURL != "" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, part.ReadURL, nil)
		if err != nil {
			return fmt.Errorf("create read request: %w", err)
		}
		for k, v := range part.ReadHeaders {
			if !strings.EqualFold(k, "host") {
				req.Header.Set(k, v)
			}
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("download original part: %w", err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("download original part: HTTP %d", resp.StatusCode)
		}
		origData, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read original part body: %w", err)
		}
	}

	data, err := readPart(part.Number, part.Size, origData)
	if err != nil {
		return fmt.Errorf("readPart callback: %w", err)
	}

	h := sha256.Sum256(data)
	checksum := base64.StdEncoding.EncodeToString(h[:])

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, part.URL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	for k, v := range part.Headers {
		if strings.EqualFold(k, "host") {
			continue
		}
		req.Header.Set(k, v)
	}
	req.ContentLength = int64(len(data))
	req.Header.Set("x-amz-checksum-sha256", checksum)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload part: HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// VaultSecret holds secret metadata returned by the management API.
type VaultSecret struct {
	Name       string    `json:"name"`
	SecretType string    `json:"secret_type"`
	Revision   int64     `json:"revision"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// VaultTokenIssueResponse is returned when issuing a scoped capability token.
type VaultTokenIssueResponse struct {
	Token     string    `json:"token"`
	TokenID   string    `json:"token_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// VaultAuditEvent is an audit event returned by the vault audit API.
type VaultAuditEvent struct {
	EventID    string         `json:"event_id"`
	EventType  string         `json:"event_type"`
	TokenID    string         `json:"token_id,omitempty"`
	AgentID    string         `json:"agent_id,omitempty"`
	TaskID     string         `json:"task_id,omitempty"`
	SecretName string         `json:"secret_name,omitempty"`
	FieldName  string         `json:"field_name,omitempty"`
	Adapter    string         `json:"adapter,omitempty"`
	Detail     map[string]any `json:"detail,omitempty"`
	Timestamp  time.Time      `json:"timestamp"`
}

func (c *Client) vaultURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return c.baseURL + "/v1/vault" + path
}

// CreateVaultSecret creates a new secret via the management API.
func (c *Client) CreateVaultSecret(ctx context.Context, name string, fields map[string]string) (*VaultSecret, error) {
	body, err := json.Marshal(map[string]any{
		"name":       name,
		"fields":     fields,
		"created_by": "ticloud-cli",
	})
	if err != nil {
		return nil, fmt.Errorf("marshal secret create request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.vaultURL("/secrets"), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var sec VaultSecret
	if err := json.NewDecoder(resp.Body).Decode(&sec); err != nil {
		return nil, fmt.Errorf("decode secret create response: %w", err)
	}
	return &sec, nil
}

// UpdateVaultSecret rotates a secret via the management API.
func (c *Client) UpdateVaultSecret(ctx context.Context, name string, fields map[string]string) (*VaultSecret, error) {
	body, err := json.Marshal(map[string]any{
		"fields":     fields,
		"updated_by": "ticloud-cli",
	})
	if err != nil {
		return nil, fmt.Errorf("marshal secret update request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.vaultURL("/secrets/"+url.PathEscape(name)), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var sec VaultSecret
	if err := json.NewDecoder(resp.Body).Decode(&sec); err != nil {
		return nil, fmt.Errorf("decode secret update response: %w", err)
	}
	return &sec, nil
}

// DeleteVaultSecret deletes a secret via the management API.
func (c *Client) DeleteVaultSecret(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.vaultURL("/secrets/"+url.PathEscape(name)), nil)
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return readError(resp)
	}
	return nil
}

// ListVaultSecrets lists secret metadata via the management API.
func (c *Client) ListVaultSecrets(ctx context.Context) ([]VaultSecret, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.vaultURL("/secrets"), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var result struct {
		Secrets []VaultSecret `json:"secrets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode secret list response: %w", err)
	}
	if result.Secrets == nil {
		result.Secrets = []VaultSecret{}
	}
	return result.Secrets, nil
}

// IssueVaultToken issues a scoped capability token via the management API.
func (c *Client) IssueVaultToken(ctx context.Context, agentID, taskID string, scope []string, ttl time.Duration) (*VaultTokenIssueResponse, error) {
	ttlSeconds := int(ttl / time.Second)
	body, err := json.Marshal(map[string]any{
		"agent_id":    agentID,
		"task_id":     taskID,
		"scope":       scope,
		"ttl_seconds": ttlSeconds,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal token issue request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.vaultURL("/tokens"), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var result VaultTokenIssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode token issue response: %w", err)
	}
	return &result, nil
}

// RevokeVaultToken revokes a capability token via the management API.
func (c *Client) RevokeVaultToken(ctx context.Context, tokenID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.vaultURL("/tokens/"+url.PathEscape(tokenID)), nil)
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return readError(resp)
	}
	return nil
}

// QueryVaultAudit queries the audit log via the management API.
func (c *Client) QueryVaultAudit(ctx context.Context, secretName string, limit int) ([]VaultAuditEvent, error) {
	u, err := url.Parse(c.vaultURL("/audit"))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if secretName != "" {
		q.Set("secret", secretName)
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var result struct {
		Events []VaultAuditEvent `json:"events"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode audit response: %w", err)
	}
	if result.Events == nil {
		result.Events = []VaultAuditEvent{}
	}
	return result.Events, nil
}

// ListReadableVaultSecrets enumerates secrets visible to the bearer capability token.
func (c *Client) ListReadableVaultSecrets(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.vaultURL("/read"), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var result struct {
		Secrets []string `json:"secrets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode readable secret list response: %w", err)
	}
	if result.Secrets == nil {
		result.Secrets = []string{}
	}
	return result.Secrets, nil
}

// ReadVaultSecret reads all fields of a secret using the consumption API.
func (c *Client) ReadVaultSecret(ctx context.Context, name string) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.vaultURL("/read/"+url.PathEscape(name)), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return nil, readError(resp)
	}
	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode secret read response: %w", err)
	}
	if result == nil {
		result = map[string]string{}
	}
	return result, nil
}

// ReadVaultSecretField reads a single field via the consumption API.
func (c *Client) ReadVaultSecretField(ctx context.Context, name, field string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.vaultURL("/read/"+url.PathEscape(name)+"/"+url.PathEscape(field)), nil)
	if err != nil {
		return "", err
	}
	resp, err := c.do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return "", readError(resp)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read field response: %w", err)
	}
	return string(data), nil
}
