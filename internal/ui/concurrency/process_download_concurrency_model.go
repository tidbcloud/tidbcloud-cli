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
	"strings"
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

type FileJob struct {
	id      int
	name    string
	url     string
	size    int64
	process progress.Model
	err     error
}

type JobInfo struct {
	idToJob       map[int]*FileJob
	viewJobs      []*FileJob
	finishedCount int
	total         int
	lock          sync.Mutex
}

type Model struct {
	concurrency int
	jobsCh      chan *FileJob
	jobInfo     *JobInfo
	width       int
	once        sync.Once
	onProgress  func(int, float64)
	onError     func(int, error)
	Interrupted bool
	outputPath  string
}

type URLMsg struct {
	Name string
	Url  string
	Size int64
}

func NewModel(urls []URLMsg, onProgress func(int, float64), onError func(int, error), concurrency int, path string) Model {
	jobs := make(chan *FileJob, len(urls))
	idToJob := make(map[int]*FileJob)

	for i, url := range urls {
		job := &FileJob{id: i, name: url.Name, url: url.Url, size: url.Size}
		idToJob[i] = job
	}

	jobInfo := &JobInfo{
		idToJob:       idToJob,
		viewJobs:      make([]*FileJob, 0, len(urls)),
		finishedCount: 0,
		total:         len(urls),
	}
	return Model{
		jobsCh:      jobs,
		jobInfo:     jobInfo,
		concurrency: concurrency,
		onProgress:  onProgress,
		onError:     onError,
		outputPath:  path,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// InitPool use pointer receiver to ensure m.once works
func (m *Model) InitPool() {
	m.once.Do(func() {
		// start consumer goroutine
		for i := 0; i < m.concurrency; i++ {
			go m.consume(m.jobsCh)
		}
		// start produce
		for _, job := range m.jobInfo.idToJob {
			go m.produce(job, m.jobsCh)
		}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			m.Interrupted = true
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.InitPool()

		return m, nil

	case ProgressErrMsg:
		// handle error msg
		f, ok := m.jobInfo.idToJob[msg.id]
		if ok {
			f.err = msg.err
			m.jobInfo.lock.Lock()
			m.jobInfo.finishedCount++
			m.jobInfo.lock.Unlock()
			if m.jobInfo.finishedCount >= m.jobInfo.total {
				var cmds []tea.Cmd
				cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
				return m, tea.Batch(cmds...)
			}
		}
		return m, nil

	case ProgressMsg:
		var cmds []tea.Cmd

		f, ok := m.jobInfo.idToJob[msg.id]
		if ok {
			cmds = append(cmds, f.process.SetPercent(msg.percent))
		}

		if msg.percent >= 1.0 {
			m.jobInfo.lock.Lock()
			m.jobInfo.finishedCount++
			m.jobInfo.lock.Unlock()
			if m.jobInfo.finishedCount >= m.jobInfo.total {
				cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
			}
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

func (m Model) View() string {
	viewString := "\n"
	pad := strings.Repeat(" ", padding)
	for _, f := range m.jobInfo.viewJobs {
		if f.process.Percent() >= 1.0 {
			viewString += fmt.Sprintf("download %s succeeded\n", f.name)
		} else if f.err != nil {
			viewString += fmt.Sprintf("download %s failed: %s\n", f.name, f.err.Error())
		} else {
			viewString += fmt.Sprintf("downloading %s | %s\n", f.name, humanize.IBytes(uint64(f.size)))
		}
		viewString += pad + f.process.View() + "\n\n"
	}

	return viewString
}

func (m *Model) produce(job *FileJob, jobsCh chan<- *FileJob) {
	jobsCh <- job
}

func (m *Model) consume(jobs <-chan *FileJob) {
	for job := range jobs {
		func() {
			// create progress bar before download
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
			resp, err := util.GetResponse(job.url, os.Getenv(config.DebugEnv) != "")
			if err != nil {
				m.onError(job.id, err)
				return
			}
			defer resp.Body.Close()

			// create file
			file, err := util.CreateFile(m.outputPath, job.name)
			if err != nil {
				m.onError(job.id, err)
				return
			}
			defer file.Close()

			// create progress writer
			pw := &progressConcurrencyWriter{
				id:     job.id,
				total:  int(resp.ContentLength),
				file:   file,
				reader: resp.Body,
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

// finalPause prevent the progress bar from exiting before it reaches 100%.
// See https://github.com/charmbracelet/bubbletea/blob/702b43d6b06287363b72836c88be35d985624a2b/examples/progress-download/tui.go#L23
func finalPause() tea.Cmd {
	return tea.Tick(time.Second*1, func(_ time.Time) tea.Msg {
		return nil
	})
}
