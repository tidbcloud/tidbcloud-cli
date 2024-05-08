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

package ui_concurrency

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/util"
)

const (
	padding  = 2
	maxWidth = 80
)

type ResultMsg struct {
	id  int
	err error
}

type jobStatus int

const (
	pending jobStatus = iota
	running
	succeed
	failed
)

type FileJob struct {
	id     int
	name   string
	url    string
	size   int64
	status jobStatus
	err    error
	pw     *progressConcurrencyWriter
}

func (f *FileJob) GetErrorString() string {
	return fmt.Sprintf("%s: %s", f.name, f.err.Error())
}

type JobInfo struct {
	idToJob map[int]*FileJob
	// viewJobs is the jobs that are start (running and finished jobs)
	startJobs    []*FileJob
	finishedJobs []*FileJob
	// count is the number of finished jobs
	count int
	// total is the total number of jobs
	total int
	lock  sync.Mutex
}

type ui struct {
	progress     *progress.Model
	downloadSize int
	totalSize    int
	width        int
	speed        int
}

type Model struct {
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

func (m *Model) SetProgram(p *tea.Program) {
	m.p = p
}

func (m *Model) GetFailedJobs() []*FileJob {
	result := make([]*FileJob, 0)
	for _, job := range m.jobInfo.finishedJobs {
		if job.status == failed {
			result = append(result, job)
		}
	}
	return result
}

func NewModel(urls []URLMsg, concurrency int, path string) *Model {
	jobs := make(chan *FileJob, len(urls))
	idToJob := make(map[int]*FileJob)
	for i, url := range urls {
		job := &FileJob{id: i, name: url.Name, url: url.Url, size: url.Size}
		idToJob[i] = job
	}
	jobInfo := &JobInfo{
		idToJob:      idToJob,
		startJobs:    make([]*FileJob, 0, len(urls)),
		finishedJobs: make([]*FileJob, 0),
		count:        0,
		total:        len(urls),
	}
	totalSize := 0
	for _, url := range urls {
		totalSize += int(url.Size)
	}
	return &Model{
		jobsCh:      jobs,
		jobInfo:     jobInfo,
		concurrency: concurrency,
		outputPath:  path,
		ui: &ui{
			totalSize: totalSize,
		},
	}
}

func (m *Model) Init() tea.Cmd {
	pro := progress.New(progress.WithDefaultGradient())
	m.ui.progress = &pro
	// start consumer goroutine
	for i := 0; i < m.concurrency; i++ {
		go m.consume(m.jobsCh)
	}
	// start produce
	for _, job := range m.jobInfo.idToJob {
		go m.produce(job, m.jobsCh)
	}
	return m.doTick()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			m.Interrupted = true
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.ui.progress.Width = msg.Width - padding*2 - 4
		if m.ui.progress.Width > maxWidth {
			m.ui.progress.Width = maxWidth
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
		if msg.err != nil {
			f.err = msg.err
			f.status = failed
		} else {
			f.status = succeed
		}
		// increase count
		m.jobInfo.lock.Lock()
		m.jobInfo.count++
		m.jobInfo.finishedJobs = append(m.jobInfo.finishedJobs, f)
		m.jobInfo.lock.Unlock()
		if m.jobInfo.count >= m.jobInfo.total {
			// stop when all jobs are finished
			cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
		}
		return m, tea.Batch(cmds...)
	default:
		return m, nil
	}
}

func (m *Model) View() string {
	viewString := "\n"
	// print finished jobs
	succeedCount := 0
	testcount := 0
	for _, f := range m.jobInfo.finishedJobs {
		testcount++
		if f.status == succeed {
			succeedCount++
			viewString += fmt.Sprintf("%d: download %s succeeded\n", testcount, f.name)
		} else if f.status == failed {
			var errMsg string
			if f.err != nil {
				errMsg = f.err.Error()
			}
			viewString += fmt.Sprintf("%d: download %s failed: %s\n", testcount, f.name, errMsg)
		}
	}
	// print process bar
	viewString += fmt.Sprintf("%s/～%s ｜ (%s/s)\n", humanize.IBytes(uint64(m.ui.downloadSize)),
		humanize.IBytes(uint64(m.ui.totalSize)), humanize.IBytes(uint64(m.ui.speed)))
	percent := float64(m.ui.downloadSize) / float64(m.ui.totalSize)
	// workaround: set to 100% when all jobs are finished in case totalSize is not accurate
	if succeedCount == m.jobInfo.total && percent < 1 {
		percent = 1
	}
	viewString += m.ui.progress.ViewAs(percent) + "\n\n"
	return viewString
}

func (m *Model) produce(job *FileJob, jobsCh chan<- *FileJob) {
	jobsCh <- job
}

func (m *Model) consume(jobs <-chan *FileJob) {
	for job := range jobs {
		func() {
			// add job to viewJobs
			m.jobInfo.lock.Lock()
			m.jobInfo.startJobs = append(m.jobInfo.startJobs, job)
			m.jobInfo.lock.Unlock()

			// request the url
			resp, err := util.GetResponse(job.url, os.Getenv(config.DebugEnv) != "")
			if err != nil {
				m.sendMsg(ResultMsg{job.id, err})
				return
			}
			defer resp.Body.Close()

			// create file
			file, err := util.CreateFile(m.outputPath, job.name)
			if err != nil {
				m.sendMsg(ResultMsg{job.id, err})
				return
			}
			defer file.Close()

			// create progress writer
			pw := &progressConcurrencyWriter{
				id:     job.id,
				total:  int(resp.ContentLength),
				file:   file,
				reader: resp.Body,
				onResult: func(id int, err error) {
					m.sendMsg(ResultMsg{id: id, err: err})
				},
			}
			job.pw = pw
			pw.Start()
		}()
	}
}

// finalPause prevent the progress bar from exiting before it reaches 100%.
// See https://github.com/charmbracelet/bubbletea/blob/702b43d6b06287363b72836c88be35d985624a2b/examples/progress-download/tui.go#L23
func finalPause() tea.Cmd {
	return tea.Tick(time.Second*1, func(_ time.Time) tea.Msg {
		return nil
	})
}

type TickMsg time.Time

func (m *Model) doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		count := 0
		for _, job := range m.jobInfo.startJobs {
			if job.pw != nil {
				count += job.pw.downloaded
			}
		}
		m.ui.speed = (count - m.ui.downloadSize) * 10
		m.ui.downloadSize = count
		return TickMsg(t)
	})
}

func (m *Model) sendMsg(msg tea.Msg) {
	if m.p != nil {
		m.p.Send(msg)
	}
}
