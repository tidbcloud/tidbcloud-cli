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
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
)

const (
	processDownloadModelPadding  = 2
	processDownloadModelMaxWidth = 80
	MaxBatchSize                 = 100
)

type ResultMsg struct {
	id     int
	err    error
	status JobStatus
}

type JobStatus string

const (
	Succeeded JobStatus = "succeeded"
	Failed    JobStatus = "failed"
	Skipped   JobStatus = "skipped"
)

type FileJob struct {
	id     int
	name   string
	url    *string
	status JobStatus
	err    error
	pw     *progressWriter
}

func (f *FileJob) GetResult() string {
	if f.status == Succeeded {
		return fmt.Sprintf("%s success", f.name)
	}
	if f.err == nil {
		return fmt.Sprintf("%s %s", f.name, f.status)
	}
	return fmt.Sprintf("%s %s: %s", f.name, f.status, f.err.Error())
}

func (f *FileJob) GetStatus() JobStatus {
	return f.status
}

type JobInfo struct {
	idToJob map[int]*FileJob
	// startedJobs is the jobs that are started (running and finished jobs), used to calculate the downloaded size
	startedJobs []*FileJob
	// finishedJobs is the jobs that are finished, used to display the ui
	finishedJobs []*FileJob
	// totalNumber is the total number of jobs
	totalNumber int
	lock        sync.Mutex
	// fileNames is the name of the files that need to be downloaded
	fileNames []string
}

type progressBar struct {
	progress       *progress.Model
	downloadedSize int
	totalSize      int
	speed          int
}

type ProcessDownloadModel struct {
	concurrency int
	jobsCh      chan *FileJob
	jobInfo     *JobInfo
	Interrupted bool
	outputPath  string
	progressBar *progressBar
	p           *tea.Program

	exportID  string
	clusterID string
	client    cloud.TiDBCloudClient
}

func (m *ProcessDownloadModel) SetProgram(p *tea.Program) {
	m.p = p
}

func (m *ProcessDownloadModel) GetFinishedJobs() []*FileJob {
	return m.jobInfo.finishedJobs
}

func NewProcessDownloadModel(concurrency int, path string,
	exportID, clusterID string, client cloud.TiDBCloudClient, totalSize int, fileNames []string) *ProcessDownloadModel {
	count := len(fileNames)
	jobInfo := &JobInfo{
		idToJob:      make(map[int]*FileJob),
		startedJobs:  make([]*FileJob, 0, count),
		finishedJobs: make([]*FileJob, 0),
		totalNumber:  count,
		fileNames:    fileNames,
	}
	jobBufferSize := concurrency
	if jobBufferSize > count {
		jobBufferSize = count
	}
	return &ProcessDownloadModel{
		jobsCh:      make(chan *FileJob, jobBufferSize),
		jobInfo:     jobInfo,
		concurrency: concurrency,
		outputPath:  path,
		progressBar: &progressBar{
			totalSize: totalSize,
		},
		exportID:  exportID,
		clusterID: clusterID,
		client:    client,
	}
}

func (m *ProcessDownloadModel) Init() tea.Cmd {
	pro := progress.New(progress.WithDefaultGradient())
	m.progressBar.progress = &pro
	// start produce
	go m.produce()
	// start consumer goroutine
	for i := 0; i < m.concurrency; i++ {
		go m.consume(m.jobsCh)
	}
	return m.doTick()
}

func (m *ProcessDownloadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			m.Interrupted = true
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progressBar.progress.Width = msg.Width - processDownloadModelPadding*2 - 4
		if m.progressBar.progress.Width > processDownloadModelMaxWidth {
			m.progressBar.progress.Width = processDownloadModelMaxWidth
		}
		return m, nil

	case ui.TickMsg:
		return m, m.doTick()

	case ResultMsg:
		var cmds []tea.Cmd
		// handle result msg
		f, ok := m.jobInfo.idToJob[msg.id]
		if !ok {
			return m, nil
		}
		// update job status
		f.status = msg.status
		f.err = msg.err
		// increase count
		m.jobInfo.lock.Lock()
		m.jobInfo.finishedJobs = append(m.jobInfo.finishedJobs, f)
		if len(m.jobInfo.finishedJobs) >= m.jobInfo.totalNumber {
			// stop when all jobs are finished
			cmds = append(cmds, tea.Sequence(ui.FinalPause(), tea.Quit))
		}
		m.jobInfo.lock.Unlock()
		return m, tea.Batch(cmds...)
	default:
		return m, nil
	}
}

func (m *ProcessDownloadModel) View() string {
	viewString := ""
	// print finished jobs
	succeededCount := 0
	for _, f := range m.jobInfo.finishedJobs {
		if f.status == Succeeded {
			succeededCount++
			viewString += fmt.Sprintf("download %s succeeded\n", f.name)
		} else if f.status == Failed {
			var errMsg string
			if f.err != nil {
				errMsg = f.err.Error()
			}
			viewString += fmt.Sprintf("download %s failed: %s\n", f.name, errMsg)
		} else if f.status == Skipped {
			var errMsg string
			if f.err != nil {
				errMsg = f.err.Error()
			}
			viewString += fmt.Sprintf("download %s skipped: %s\n", f.name, errMsg)
		}
	}
	// print process bar
	viewString += fmt.Sprintf("%s/~%s (%s/s) with ~%d files(s) remaining\n", humanize.IBytes(uint64(m.progressBar.downloadedSize)),
		humanize.IBytes(uint64(m.progressBar.totalSize)), humanize.IBytes(uint64(m.progressBar.speed)), m.jobInfo.totalNumber-len(m.jobInfo.finishedJobs))
	percent := float64(m.progressBar.downloadedSize) / float64(m.progressBar.totalSize)
	// workaround: set to 100% when all jobs are finished in case totalSize is not accurate
	if succeededCount == m.jobInfo.totalNumber && percent < 1 {
		percent = 1
	}
	viewString += m.progressBar.progress.ViewAs(percent) + "\n\n"
	return viewString
}

func (m *ProcessDownloadModel) GetBatchSize() int {
	batchSize := 2 * m.concurrency
	if batchSize > MaxBatchSize {
		batchSize = MaxBatchSize
	}
	return batchSize
}

func (m *ProcessDownloadModel) produce() {
	batchSize := m.GetBatchSize()
	jobId := 0
	ctx := context.Background()
	for len(m.jobInfo.fileNames) > 0 {
		// request the next batch
		if batchSize > len(m.jobInfo.fileNames) {
			batchSize = len(m.jobInfo.fileNames)
		}
		downloadFileNames := m.jobInfo.fileNames[:batchSize]
		m.jobInfo.fileNames = m.jobInfo.fileNames[batchSize:]
		body := &export.ExportServiceDownloadExportFilesBody{
			FileNames: downloadFileNames,
		}
		resp, err := m.client.DownloadExportFiles(ctx, m.clusterID, m.exportID, body)
		if err != nil {
			for _, name := range downloadFileNames {
				jobId++
				job := &FileJob{id: jobId, name: name, err: err}
				m.jobInfo.idToJob[jobId] = job
				m.jobsCh <- job
			}
			continue
		}
		for _, file := range resp.Files {
			jobId++
			job := &FileJob{
				id:   jobId,
				name: *file.Name,
				url:  file.Url,
			}
			m.jobInfo.idToJob[jobId] = job
			m.jobsCh <- job
		}
	}
	close(m.jobsCh)
}
func (m *ProcessDownloadModel) consume(jobs <-chan *FileJob) {
	for job := range jobs {
		func() {
			// add job to startJobs
			m.jobInfo.lock.Lock()
			m.jobInfo.startedJobs = append(m.jobInfo.startedJobs, job)
			m.jobInfo.lock.Unlock()

			// check job
			if job.err != nil {
				m.sendMsg(ResultMsg{job.id, job.err, Failed})
				return
			}
			if job.url == nil {
				m.sendMsg(ResultMsg{job.id, errors.New("empty download url"), Failed})
				return
			}

			// request the url
			resp, err := util.GetResponse(*job.url, os.Getenv(config.DebugEnv) != "")
			if err != nil {
				m.sendMsg(ResultMsg{job.id, err, Failed})
				return
			}
			defer resp.Body.Close()

			// create file
			file, err := util.CreateFile(m.outputPath, job.name)
			if err != nil {
				if strings.Contains(err.Error(), "file already exists") {
					m.sendMsg(ResultMsg{job.id, err, Skipped})
				} else {
					m.sendMsg(ResultMsg{job.id, err, Failed})
				}
				return
			}
			defer file.Close()

			// create progress writer
			pw := &progressWriter{
				id:     job.id,
				file:   file,
				reader: resp.Body,
				onResult: func(id int, err error, status JobStatus) {
					m.sendMsg(ResultMsg{id: id, err: err, status: status})
				},
			}
			job.pw = pw
			pw.Start()
		}()
	}
}

// doTick update the download size and speed every 100ms
func (m *ProcessDownloadModel) doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		totalDownloaded := 0
		for _, job := range m.jobInfo.startedJobs {
			if job.pw != nil {
				totalDownloaded += job.pw.downloadedSize
			}
		}
		m.progressBar.speed = (totalDownloaded - m.progressBar.downloadedSize) * 10
		m.progressBar.downloadedSize = totalDownloaded
		return ui.TickMsg(t)
	})
}

func (m *ProcessDownloadModel) sendMsg(msg tea.Msg) {
	if m.p != nil {
		m.p.Send(msg)
	}
}

type progressWriter struct {
	id             int
	downloadedSize int
	file           *os.File
	reader         io.Reader
	onResult       func(int, error, JobStatus)
}

func (pw *progressWriter) Read(p []byte) (n int, err error) {
	n, err = pw.reader.Read(p)
	if err == nil || err == io.EOF {
		pw.downloadedSize += n
	}
	return
}

func (pw *progressWriter) Start() {
	_, err := io.Copy(pw.file, pw)
	if err != nil {
		pw.onResult(pw.id, err, Failed)
	} else {
		pw.onResult(pw.id, nil, Succeeded)
	}
}
