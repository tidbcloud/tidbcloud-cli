package ui

import (
	"fmt"
	"strings"
	"time"

	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/pkg/tidbcloud/pingchat/models"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type ChatMessage struct {
	Role    string
	Content string
}

type ChatBoxModel struct {
	Interrupted bool
	Err         error

	chatLog  []ChatMessage
	textarea textarea.Model
	viewport viewport.Model

	ready       bool
	isLoading   bool
	sendMsgFunc func(messages []ChatMessage) bubbletea.Msg
}

type TickMsg time.Time

type EndSendingMsg struct {
	Err error
	Msg ChatMessage
}

func InitialChatBoxModel(sendMsgFunc func(messages []ChatMessage) bubbletea.Msg) ChatBoxModel {
	var chatLog []ChatMessage

	ta := textarea.New()
	ta.Placeholder = "Type your message here..."
	ta.Focus()
	ta.CharLimit = 500
	ta.ShowLineNumbers = false
	ta.Prompt = "┃ "
	ta.SetHeight(5)

	return ChatBoxModel{
		sendMsgFunc: sendMsgFunc,
		chatLog:     chatLog,
		textarea:    ta,
		isLoading:   false,
	}
}

func (m ChatBoxModel) Init() bubbletea.Cmd {
	return bubbletea.Batch(
		textarea.Blink,
	)
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
					Role:    "user",
					Content: inputMessage,
				}, ChatMessage{
					Role:    "ai",
					Content: "Thinking...",
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
	case ChatMessage:
		m.chatLog = append(m.chatLog, msg)
		m.viewport.SetContent(m.RenderChatLog())
		m.viewport.GotoBottom()
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
			Role:    models.PingchatChatMessageRoleUser,
		}
		return m.sendMsgFunc([]ChatMessage{
			msg,
		})
	}
}

var chatUserTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Align(lipgloss.Left).Width(6).Render
var chatAITextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Align(lipgloss.Left).Width(12).Render
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
		s := util.String(message.Content, maxWidth)
		out, _ := r.Render(s)

		var who string
		if message.Role == "user" {
			who = chatUserTextStyle("You:")
		} else {
			who = chatAITextStyle("TiDB Bot:")
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
