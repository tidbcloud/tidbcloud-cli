package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
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
	str := color.New(color.FgYellow).Sprintf("%s %s...press q to quit \n", m.spinner.View(), m.Hint)
	if m.quitting {
		return str + "\n"
	}
	return str
}
