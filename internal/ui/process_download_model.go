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

package ui

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/util"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
)

const (
	processDownloadModelPadding  = 2
	processDownloadModelMaxWidth = 80
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
	url    string
	size   int64
	status JobStatus
	err    error
	pw     *progressWriter
}

func (f *FileJob) GetErrorString() string {
	if f.status == Succeeded {
		return ""
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
	// startedJobs is the jobs that are started (running and finished jobs)
	startedJobs  []*FileJob
	finishedJobs []*FileJob
	// totalNumber is the total number of jobs
	totalNumber int
	lock        sync.Mutex
}

type ui struct {
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
	ui          *ui
	p           *tea.Program
}

type URLMsg struct {
	Name string
	Url  string
	Size int64
}

func (m *ProcessDownloadModel) SetProgram(p *tea.Program) {
	m.p = p
}

func (m *ProcessDownloadModel) GetFinishedJobs() []*FileJob {
	return m.jobInfo.finishedJobs
}

func NewProcessDownloadModel(urls []URLMsg, concurrency int, path string) *ProcessDownloadModel {
	jobs := make(chan *FileJob, len(urls))
	idToJob := make(map[int]*FileJob)
	for i, url := range urls {
		job := &FileJob{id: i, name: url.Name, url: url.Url, size: url.Size}
		idToJob[i] = job
	}
	jobInfo := &JobInfo{
		idToJob:      idToJob,
		startedJobs:  make([]*FileJob, 0, len(urls)),
		finishedJobs: make([]*FileJob, 0),
		totalNumber:  len(urls),
	}
	totalSize := 0
	for _, url := range urls {
		totalSize += int(url.Size)
	}
	return &ProcessDownloadModel{
		jobsCh:      jobs,
		jobInfo:     jobInfo,
		concurrency: concurrency,
		outputPath:  path,
		ui: &ui{
			totalSize: totalSize,
		},
	}
}

func (m *ProcessDownloadModel) Init() tea.Cmd {
	pro := progress.New(progress.WithDefaultGradient())
	m.ui.progress = &pro
	// start consumer goroutine
	for i := 0; i < m.concurrency; i++ {
		go m.consume(m.jobsCh)
	}
	// start produce
	go m.produce(m.jobsCh)
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
		m.ui.progress.Width = msg.Width - processDownloadModelPadding*2 - 4
		if m.ui.progress.Width > processDownloadModelMaxWidth {
			m.ui.progress.Width = processDownloadModelMaxWidth
		}
		return m, nil

	case TickMsg:
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
			cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
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
	viewString += fmt.Sprintf("%s/~%s (%s/s) with ~%d files(s) remaining\n", humanize.IBytes(uint64(m.ui.downloadedSize)),
		humanize.IBytes(uint64(m.ui.totalSize)), humanize.IBytes(uint64(m.ui.speed)), m.jobInfo.totalNumber-len(m.jobInfo.finishedJobs))
	percent := float64(m.ui.downloadedSize) / float64(m.ui.totalSize)
	// workaround: set to 100% when all jobs are finished in case totalSize is not accurate
	if succeededCount == m.jobInfo.totalNumber && percent < 1 {
		percent = 1
	}
	viewString += m.ui.progress.ViewAs(percent) + "\n\n"
	return viewString
}

func (m *ProcessDownloadModel) produce(jobsCh chan<- *FileJob) {
	for _, job := range m.jobInfo.idToJob {
		jobsCh <- job
	}
	close(jobsCh)
}

func (m *ProcessDownloadModel) consume(jobs <-chan *FileJob) {
	for job := range jobs {
		func() {
			// add job to startJobs
			m.jobInfo.lock.Lock()
			m.jobInfo.startedJobs = append(m.jobInfo.startedJobs, job)
			m.jobInfo.lock.Unlock()

			// request the url
			resp, err := util.GetResponse(job.url, os.Getenv(config.DebugEnv) != "")
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

func (m *ProcessDownloadModel) doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		totalDownloaded := 0
		for _, job := range m.jobInfo.startedJobs {
			if job.pw != nil {
				totalDownloaded += job.pw.downloadedSize
			}
		}
		m.ui.speed = (totalDownloaded - m.ui.downloadedSize) * 10
		m.ui.downloadedSize = totalDownloaded
		return TickMsg(t)
	})
}

func (m *ProcessDownloadModel) sendMsg(msg tea.Msg) {
	if m.p != nil {
		m.p.Send(msg)
	}
}
