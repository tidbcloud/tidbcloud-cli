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
	"os"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/pingcap/log"
	"go.uber.org/zap"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type SpinnerModel struct {
	Hint        string
	spinner     spinner.Model
	Interrupted bool
	Err         error
	Output      string
	Task        tea.Cmd
}

type Result string

func InitialSpinnerModel(task tea.Cmd, hint string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("32"))
	return SpinnerModel{spinner: s, Task: task, Output: "", Hint: hint}
}

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.Task)
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				log.Debug("failed to find current process when interrupted in spinner", zap.Error(err))
			}
			err = p.Signal(syscall.SIGINT)
			if err != nil {
				log.Debug("failed to send SIGINT to current process when interrupted in spinner", zap.Error(err))
			}
			m.Interrupted = true
			return m, nil
		default:
			return m, nil
		}

	case errMsg:
		m.Err = msg
		return m, tea.Quit

	case Result:
		m.Output = string(msg)
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m SpinnerModel) View() string {
	str := color.New(color.FgYellow).Sprintf("%s %s", m.spinner.View(), color.BlueString(m.Hint)) + "\n"
	if m.Interrupted {
		return str + "\n"
	}
	return str
}
