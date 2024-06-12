// Copyright 2024 PingCAP, Inc.
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

package export

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	exportApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/client/export_service"
)

var wg = sync.WaitGroup{}

const BatchSize = 20

// downloadPool download export files concurrently
type downloadPool struct {
	path        string
	concurrency int
	fileNames   []string
	client      cloud.TiDBCloudClient
	// The size of the batch to request download url
	clusterID string
	exportID  string
	h         *internal.Helper

	jobs    chan *downloadJob
	results chan *downloadResult

	// fileJobs is a list of jobs that will be set to the jobs channel
	fileJobs  []*downloadJob
	batchSize int
}

func NewDownloadPool(h *internal.Helper, files []string, path string,
	concurrency int, exportID, clusterID string, client cloud.TiDBCloudClient) (*downloadPool, error) {
	if len(files) <= 0 {
		return nil, errors.New("no files to download")
	}
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}

	jobBufferSize := 2 * concurrency
	if jobBufferSize > len(files) {
		jobBufferSize = len(files)
	}
	jobs := make(chan *downloadJob, jobBufferSize)
	results := make(chan *downloadResult, len(files))

	return &downloadPool{
		h:           h,
		path:        path,
		concurrency: concurrency,
		fileNames:   files,
		client:      client,
		clusterID:   clusterID,
		exportID:    exportID,
		jobs:        jobs,
		results:     results,
		batchSize:   BatchSize,
	}, nil
}

type downloadJob struct {
	fileName    string
	downloadUrl *string
	size        int
	err         error
}

type downloadResult struct {
	name   string
	err    error
	status ui.JobStatus
}

func (r *downloadResult) GetResult() string {
	if r.status == ui.Succeeded {
		return fmt.Sprintf("%s success", r.name)
	}
	if r.err == nil {
		return fmt.Sprintf("%s %s", r.name, r.status)
	}
	return fmt.Sprintf("%s %s: %s", r.name, r.status, r.err.Error())
}

func (d *downloadPool) Start() error {
	// create the path if not exist
	err := util.CreateFolder(d.path)
	if err != nil {
		return err
	}
	fmt.Fprintf(d.h.IOStreams.Out, color.GreenString("start to download files to %s:\n", d.path))
	// start produce
	go d.produce()
	// start consumers:
	for i := 0; i < d.concurrency; i++ {
		wg.Add(1)
		go d.consume()
	}
	// wait for all consumers to finish
	wg.Wait()
	close(d.results)
	// summarize the download results
	succeededCount := 0
	failedCount := 0
	skippedCount := 0
	downloadResults := make([]*downloadResult, 0)
	for result := range d.results {
		switch result.status {
		case ui.Succeeded:
			succeededCount++
		case ui.Failed:
			failedCount++
		case ui.Skipped:
			skippedCount++
		}
		downloadResults = append(downloadResults, result)
	}
	fmt.Fprintf(d.h.IOStreams.Out, generateDownloadSummary(succeededCount, skippedCount, failedCount))
	index := 0
	for _, f := range downloadResults {
		if f.status != ui.Succeeded {
			index++
			fmt.Fprintf(d.h.IOStreams.Out, "%d.%s\n", index, f.GetResult())
		}
	}
	if failedCount > 0 {
		return fmt.Errorf("%d file(s) failed to download", failedCount)
	}
	return nil
}

func (d *downloadPool) produce() {
	jobSize := len(d.fileNames)
	for i := 0; i < jobSize; i++ {
		// request the next batch when the fileJobs are not enough
		if len(d.fileJobs) < i+1 {
			size := d.batchSize
			if size > len(d.fileNames) {
				size = len(d.fileNames)
			}
			downloadFileNames := d.fileNames[:size]
			d.fileNames = d.fileNames[size:]
			body := exportApi.ExportServiceDownloadExportFilesBody{
				FileNames: downloadFileNames,
			}
			params := exportApi.NewExportServiceDownloadExportFilesParams().WithClusterID(d.clusterID).
				WithExportID(d.exportID).WithBody(body)
			resp, err := d.client.DownloadExportFiles(params)
			if err != nil {
				for _, file := range downloadFileNames {
					d.fileJobs = append(d.fileJobs, &downloadJob{fileName: file, err: err})
				}
				return
			}
			for _, file := range resp.Payload.Files {
				job := &downloadJob{
					fileName:    file.Name,
					downloadUrl: file.DownloadURL,
				}
				d.fileJobs = append(d.fileJobs, job)
			}
		}
		d.jobs <- d.fileJobs[i]
	}
	close(d.jobs)
}

func (d *downloadPool) consume() {
	defer wg.Done()
	for job := range d.jobs {
		func() {
			var err error
			defer func() {
				// record result
				if err != nil {
					if strings.Contains(err.Error(), "file already exists") {
						fmt.Fprintf(d.h.IOStreams.Out, "download %s skipped: %s\n", job.fileName, err.Error())
						d.results <- &downloadResult{name: job.fileName, err: err, status: ui.Skipped}
					} else {
						fmt.Fprintf(d.h.IOStreams.Out, "download %s failed: %s\n", job.fileName, err.Error())
						d.results <- &downloadResult{name: job.fileName, err: err, status: ui.Failed}
					}
				} else {
					fmt.Fprintf(d.h.IOStreams.Out, "download %s succeeded | %s\n", job.fileName, humanize.IBytes(uint64(job.size)))
					d.results <- &downloadResult{name: job.fileName, err: nil, status: ui.Succeeded}
				}
			}()

			// request the url
			if job.err != nil {
				err = job.err
				return
			}
			if job.downloadUrl == nil {
				err = errors.New("empty download url")
				return
			}
			resp, err := util.GetResponse(*job.downloadUrl, os.Getenv(config.DebugEnv) != "")
			if err != nil {
				return
			}
			job.size = int(resp.ContentLength)
			defer resp.Body.Close()

			file, err := util.CreateFile(d.path, job.fileName)
			if err != nil {
				return
			}
			defer file.Close()
			_, err = io.Copy(file, resp.Body)
		}()
	}
}

func generateDownloadSummary(succeededCount, skippedCount, failedCount int) string {
	summaryMessage := fmt.Sprintf("%s %s %s", color.BlueString("download summary:"), color.GreenString("succeeded: %d", succeededCount), color.GreenString("skipped: %d", skippedCount))
	if failedCount > 0 {
		summaryMessage += color.RedString(" failed: %d", failedCount)
	} else {
		summaryMessage += fmt.Sprintf(" failed: %d", failedCount)
	}
	return summaryMessage + "\n"
}
