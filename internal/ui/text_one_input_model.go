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

package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
)

type TextOneInputModel struct {
	Prompt      string
	Input       textinput.Model
	Err         error
	Interrupted bool
}

func (m TextOneInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextOneInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Interrupted = true
			return m, tea.Quit
		}

	case errMsg:
		m.Err = msg
		return m, nil
	}

	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m TextOneInputModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n",
		m.Prompt,
		m.Input.View(),
		helpMessageStyle("Press Enter to submit (esc to quit)"),
	)
}

// InitialOneInputModel runs an interactive single-line input and returns the model and any error.
// view is the prompt line shown above the input. placeholder is used as the input placeholder text.
func InitialOneInputModel(prompt, placeholder string) (TextOneInputModel, error) {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.PromptStyle = config.FocusedStyle
	ti.TextStyle = config.FocusedStyle

	p := tea.NewProgram(TextOneInputModel{Prompt: prompt, Input: ti})
	model, err := p.Run()
	finalModel := model.(TextOneInputModel)
	if err != nil {
		return finalModel, errors.Trace(err)
	}
	if finalModel.Interrupted {
		return finalModel, util.InterruptError
	}
	if finalModel.Err != nil {
		return finalModel, finalModel.Err
	}
	return finalModel, nil
}
