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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

type ProgressMsg struct {
	id      int
	percent float64
}

func NewProgressMsg(id int, percent float64) ProgressMsg {
	return ProgressMsg{id: id, percent: percent}
}

type ProgressErrMsg struct {
	id  int
	err error
}

func NewProgressErrMsg(id int, err error) ProgressErrMsg {
	return ProgressErrMsg{id: id, err: err}
}

func finalPause() tea.Cmd {
	return tea.Tick(time.Second*2, func(_ time.Time) tea.Msg {
		return nil
	})
}

type FileJob struct {
	id            int
	path          string
	name          string
	url           string
	process       progress.Model
	reader        io.ReadCloser
	file          *os.File
	ContentLength int64
	err           error
}

type JobInfo struct {
	idToJob     map[int]*FileJob
	pendingJobs []*FileJob
	viewJobs    []*FileJob
	count       int
	lock        sync.Mutex
}

type model struct {
	concurrency int
	jobsCh      chan *FileJob
	jobInfo     *JobInfo
	width       int
	once        sync.Once
	onProgress  func(int, float64)
	onError     func(int, error)
}

type URLMsg struct {
	Name string
	Path string
	Url  string
}

func NewModel(urls []URLMsg, onProgress func(int, float64), onError func(int, error), concurrency int) model {
	jobs := make(chan *FileJob, len(urls))
	idToJob := make(map[int]*FileJob)
	pendingJobs := make([]*FileJob, 0, len(urls))

	for i, url := range urls {
		job := &FileJob{id: i, name: url.Name, url: url.Url, path: url.Path}
		idToJob[i] = job
		pendingJobs = append(pendingJobs, job)
	}

	jobInfo := &JobInfo{
		idToJob:     idToJob,
		pendingJobs: pendingJobs,
		viewJobs:    make([]*FileJob, 0, len(urls)),
		count:       0,
	}
	return model{
		jobsCh:      jobs,
		jobInfo:     jobInfo,
		concurrency: concurrency,
		onProgress:  onProgress,
		onError:     onError,
	}
}

func NewDefaultModel(urls []URLMsg, onProgress func(int, float64), onError func(int, error)) model {
	return NewModel(urls, onProgress, onError, 2)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) InitPool() {
	m.once.Do(func() {
		// start consumer goroutine
		for i := 0; i < m.concurrency; i++ {
			go m.consume(m.jobsCh)
		}
		// start produce
		for _, job := range m.jobInfo.pendingJobs {
			go m.produce(job, m.jobsCh)
		}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.InitPool()

		return m, nil

	case ProgressErrMsg:
		f, ok := m.jobInfo.idToJob[msg.id]
		if ok {
			f.err = msg.err
			m.jobInfo.lock.Lock()
			m.jobInfo.count++
			m.jobInfo.lock.Unlock()
			if m.jobInfo.count >= len(m.jobInfo.idToJob) {
				var cmds []tea.Cmd
				cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
				return m, tea.Batch(cmds...)
			}
		}
		return m, nil

	case ProgressMsg:
		var cmds []tea.Cmd

		if msg.percent >= 1.0 {
			m.jobInfo.lock.Lock()
			m.jobInfo.count++
			m.jobInfo.lock.Unlock()
			if m.jobInfo.count >= len(m.jobInfo.idToJob) {
				cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
			}
		}

		f, ok := m.jobInfo.idToJob[msg.id]
		if ok {
			cmds = append(cmds, f.process.SetPercent(msg.percent))
		}
		return m, tea.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		// we can not find which one send the FrameMsg, so we just update all in viewJobs
		for _, f := range m.jobInfo.viewJobs {
			progressModel, cmd := f.process.Update(msg)
			if cmd != nil {
				f.process = progressModel.(progress.Model)
				return m, cmd
			}
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m model) View() string {
	viewString := "\n"
	pad := strings.Repeat(" ", padding)
	for _, f := range m.jobInfo.viewJobs {
		if f.process.Percent() >= 1.0 {
			viewString += "download " + f.name + " success" + "\n"
		} else if f.err != nil {
			viewString += "download " + f.name + " failed: " + f.err.Error() + "\n"
		} else {
			viewString += "downloading " + f.name + "\n"
		}
		viewString += pad + f.process.View() + "\n\n"
	}
	viewString += pad + helpStyle("Press ctrl+c key to quit")

	return viewString
}

func (m *model) produce(job *FileJob, jobsCh chan<- *FileJob) {
	jobsCh <- job
}

func (m *model) consume(jobs <-chan *FileJob) {
	for job := range jobs {
		func() {
			pro := progress.New(progress.WithDefaultGradient())
			pro.Width = m.width - padding*2 - 4
			if pro.Width > maxWidth {
				pro.Width = maxWidth
			}
			job.process = pro

			// add job to viewJobs
			m.jobInfo.lock.Lock()
			m.jobInfo.viewJobs = append(m.jobInfo.viewJobs, job)
			m.jobInfo.lock.Unlock()

			// request the url
			resp, err := getResponse(job.url)
			if err != nil {
				m.onError(job.id, err)
				return
			}
			job.reader = resp.Body
			defer job.reader.Close()

			if resp.ContentLength <= 0 {
				m.onError(job.id, err)
				return
			}
			job.ContentLength = resp.ContentLength

			// skip if the file exists
			if _, err := os.Stat(job.path + "/" + job.name); err == nil {
				m.onError(job.id, errors.New("file already exists"))
				return
			}
			file, err := os.Create(job.path + "/" + job.name)
			if err != nil {
				m.onError(job.id, err)
				return
			}
			job.file = file
			defer job.file.Close()

			pw := &progressConcurrencyWriter{
				id:     job.id,
				total:  int(job.ContentLength),
				file:   job.file,
				reader: job.reader,
				onProgress: func(id int, ratio float64) {
					m.onProgress(id, ratio)
				},
				onError: func(id int, err error) {
					m.onError(id, err)
				},
			}
			pw.Start()
		}()
	}
}

func getResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url) // nolint:gosec
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receiving status of %d for url: %s", resp.StatusCode, url)
	}
	return resp, nil
}
