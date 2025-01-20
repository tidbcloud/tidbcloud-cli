// Copyright 2025 PingCAP, Inc.
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

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
)

type TextAreaModel struct {
	Textarea    textarea.Model
	Err         error
	Interrupted bool
}

func InitialTextAreaModel(placeholder string) (TextAreaModel, error) {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.Focus()
	ta.SetWidth(80)
	ta.SetHeight(20)
	ta.ShowLineNumbers = false
	ta.CharLimit = 0

	p := tea.NewProgram(TextAreaModel{Textarea: ta})
	model, err := p.Run()
	finalModel := model.(TextAreaModel)
	if err != nil {
		return finalModel, err
	}
	if finalModel.Interrupted {
		return finalModel, util.InterruptError
	}
	if finalModel.Err != nil {
		return finalModel, finalModel.Err
	}
	return finalModel, nil
}

func (m TextAreaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m TextAreaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.Textarea.Focused() {
				m.Textarea.Blur()
			}
		case tea.KeyCtrlS:
			return m, tea.Quit
		case tea.KeyCtrlC:
			m.Interrupted = true
			return m, tea.Quit
		default:
			if !m.Textarea.Focused() {
				cmd = m.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.Err = msg
		return m, nil
	}

	m.Textarea, cmd = m.Textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m TextAreaModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n",
		m.Textarea.View(),
		helpMessageStyle("Press Ctrl+S to save and quit"),
	)
}
