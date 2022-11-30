// Copyright 2022 PingCAP, Inc.
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
	"github.com/fatih/color"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type SpinnerModel struct {
	Hint     string
	spinner  spinner.Model
	quitting bool
	Err      error
	Output   string
	Task     tea.Cmd
}

type Result string

func InitialSpinnerModel(task tea.Cmd, hint string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Points
	//s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return SpinnerModel{spinner: s, Task: task, Output: "", Hint: hint}
}

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.Task)
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
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
	if m.Err != nil {
		return m.Err.Error()
	}
	if m.Output != "" {
		return m.Output
	}
	str := color.New(color.FgYellow).Sprintf("%s %s.\n", m.spinner.View(), m.Hint)
	if m.quitting {
		return str + "\n"
	}
	return str
}
