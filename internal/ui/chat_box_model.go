// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

import (
	"fmt"
	"strings"
	"time"

	"tidbcloud-cli/internal/util"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Role string

const (
	RoleUser Role = "user"
	RoleBot  Role = "bot"

	LoadingPrompt = "AI generating"
)

type ChatMessage struct {
	Role    Role
	Content string
}

type ChatBoxModel struct {
	Interrupted bool
	Err         error

	botName  string
	chatLog  []ChatMessage
	textarea textarea.Model
	viewport viewport.Model

	tickNumber  int
	ready       bool
	isLoading   bool
	sendMsgFunc func(messages []ChatMessage) bubbletea.Msg
}

type TickMsg time.Time

type EndSendingMsg struct {
	Err error
	Msg ChatMessage
}

func InitialChatBoxModel(sendMsgFunc func(messages []ChatMessage) bubbletea.Msg, botName string) ChatBoxModel {
	var chatLog []ChatMessage

	ta := textarea.New()
	ta.Placeholder = "Type your message here..."
	ta.Focus()
	// no limit
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.Prompt = "┃ "
	ta.SetHeight(5)

	return ChatBoxModel{
		sendMsgFunc: sendMsgFunc,
		chatLog:     chatLog,
		textarea:    ta,
		isLoading:   false,
		botName:     botName,
		tickNumber:  0,
	}
}

func (m ChatBoxModel) Init() bubbletea.Cmd {
	return bubbletea.Batch(
		tick(), // tick
	)
}

func tick() bubbletea.Cmd {
	return bubbletea.Every(time.Second/2, func(t time.Time) bubbletea.Msg {
		return TickMsg(t)
	})
}

func (m ChatBoxModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	var (
		taCmd   bubbletea.Cmd
		vpCmd   bubbletea.Cmd
		nextCmd bubbletea.Cmd
		vpCmd2  bubbletea.Cmd
	)

	m.viewport, vpCmd = m.viewport.Update(msg)
	m.textarea, taCmd = m.textarea.Update(msg)

	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		switch msg.Type {
		case bubbletea.KeyCtrlC:
			m.Interrupted = true
			return m, bubbletea.Quit
		case bubbletea.KeyEsc:
			return m, bubbletea.Quit
		case bubbletea.KeyCtrlS:
			inputMessage := m.textarea.Value()
			// trim whitespace
			inputMessage = strings.TrimSpace(inputMessage)

			if len(inputMessage) > 0 {
				m.chatLog = append(m.chatLog, ChatMessage{
					Role:    RoleUser,
					Content: inputMessage,
				}, ChatMessage{
					Role:    RoleBot,
					Content: LoadingPrompt,
				})
				m.viewport.SetContent(m.RenderChatLog())
				m.textarea.Reset()
				m.isLoading = true
				m.viewport.GotoBottom()

				nextCmd = m.sendMessage(inputMessage)
			}
		}
	case bubbletea.WindowSizeMsg:
		verticalMarginHeight := m.textarea.Height() + 2

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		m.textarea.SetWidth(msg.Width)
		m.viewport.SetContent(m.RenderChatLog())
	case TickMsg:
		if m.isLoading {
			m.tickNumber = (m.tickNumber + 1) % 4
			m.chatLog[len(m.chatLog)-1] = ChatMessage{
				Role:    RoleBot,
				Content: fmt.Sprintf("%s%s", LoadingPrompt, generateDots(m.tickNumber%4)),
			}
			m.viewport.SetContent(m.RenderChatLog())
		}
		nextCmd = tick()
	case EndSendingMsg:
		m.isLoading = false
		if msg.Err != nil {
			m.Err = msg.Err
			return m, bubbletea.Quit
		}
		m.chatLog[len(m.chatLog)-1] = msg.Msg
		m.viewport.SetContent(m.RenderChatLog())
		m.viewport.GotoBottom()
	}

	return m, bubbletea.Batch(taCmd, vpCmd, nextCmd, vpCmd2)
}

func (m ChatBoxModel) sendMessage(prompt string) bubbletea.Cmd {
	return func() bubbletea.Msg {
		msg := ChatMessage{
			Content: prompt,
			Role:    RoleUser,
		}
		return m.sendMsgFunc([]ChatMessage{
			msg,
		})
	}
}

var chatUserTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Align(lipgloss.Left).Width(6).Render
var chatAITextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Align(lipgloss.Left).Width(20).Render
var helpMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Align(lipgloss.Left).Render

func (m ChatBoxModel) RenderChatLog() string {
	var maxWidth = m.viewport.Width - 6
	r, _ := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
		// wrap output at specific width (default is 80)
		glamour.WithWordWrap(maxWidth),
		glamour.WithPreservedNewLines(),
	)

	chatLogString := ""
	for _, message := range m.chatLog {
		// Due to a bug in the formatting of Chinese characters (see https://github.com/charmbracelet/glamour/pull/249),
		// glamour cannot correctly word-wrap Chinese characters. Therefore, we need to wrap the string before rendering it.
		s := util.String(message.Content, maxWidth)
		out, _ := r.Render(s)

		var who string
		if message.Role == RoleUser {
			who = chatUserTextStyle("You:")
		} else {
			who = chatAITextStyle(m.botName + ":")
		}

		chatLogString += fmt.Sprintf("%s\n%s", who, out)
	}

	return chatLogString
}

func (m ChatBoxModel) View() string {
	helpMessage := helpMessageStyle("Press Ctrl+S to send message (esc to quit)")

	if m.isLoading {
		return m.viewport.View()
	} else {
		return fmt.Sprintf(
			"%s\n\n%s\n%s",
			m.viewport.View(),
			m.textarea.View(),
			helpMessage,
		)
	}
}

func generateDots(n int) string {
	return strings.Repeat(".", n)
}
