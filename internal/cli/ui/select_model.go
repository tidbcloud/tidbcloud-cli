package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type SelectModel struct {
	Hint        string
	Choices     []interface{} // items on the to-do list
	cursor      int           // which to-do list item our cursor is pointing at
	Selected    int           // which to-do items are Selected
	Interrupted bool
}

func InitialSelectModel(choices []interface{}, hint string) SelectModel {
	return SelectModel{
		// Our to-do list is a grocery list
		Choices: choices,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		Selected:    -1,
		Interrupted: false,
		Hint:        hint,
	}
}

func (m SelectModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			m.Interrupted = true
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.Choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter":
			m.Selected = m.cursor
			return m, tea.Quit
		}
	}

	// Return the updated SpinnerModel to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m SelectModel) View() string {
	// The header
	s := color.New(color.FgHiGreen).Sprint(m.Hint) + "\n"

	// Iterate over our choices
	for i, choice := range m.Choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if i == m.Selected {
			checked = "x" // selected!
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
			break
		}

		// Render the row
		if m.Selected == -1 {
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	}

	// Send the UI for rendering
	return s
}
