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
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/sahilm/fuzzy"
)

var (
	activeDot   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	inactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
)

type SelectModel struct {
	Hint           string
	Choices        []interface{} // items to be selected, which need implement fmt.Stringer interface
	cursor         int           // which item our cursor is pointing at
	Selected       int           // which items are Selected
	Interrupted    bool
	showPagination bool
	showFilter     bool

	// Copy from Choices to show in terminal. When enabling filter,
	// VisibleChoices is part of Choices and is filtered by user input.
	VisibleChoices []interface{}

	helper      help.Model
	keys        keyMap
	Paginator   paginator.Model
	FilterInput textinput.Model

	// Function used to get value string from Choices item.
	// Only use when enabling filter.
	ChoiceValue func(choice interface{}) string
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Select, k.Quit},              // second column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "pgup"),
		key.WithHelp("←/pgup", "prev page"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "pgdown"),
		key.WithHelp("→/pgdown", "next page"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
}

func InitialSelectModel(choices []interface{}, hint string) (*SelectModel, error) {
	if len(choices) == 0 {
		return nil, errors.New("There are no available choices")
	}

	p := buildPaginator(len(choices))
	f := buildFilterInput()
	df := defaultChoicesValueFunc()

	return &SelectModel{
		// Our to-do list is a grocery list
		Choices: choices,
		// init VisibleChoices and make it the same as Choices
		VisibleChoices: choices,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		Selected:       -1,
		Interrupted:    false,
		Hint:           hint,
		showFilter:     false,
		showPagination: false,
		FilterInput:    f,
		Paginator:      p,
		ChoiceValue:    df,
		keys:           keys,
		helper:         help.New(),
	}, nil
}

func buildPaginator(numOfTotalItems int) paginator.Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = numOfTotalItems
	p.ActiveDot = activeDot
	p.InactiveDot = inactiveDot
	return p
}

func buildFilterInput() textinput.Model {
	f := textinput.New()
	f.Placeholder = "Type to filter"
	f.Prompt = " "
	return f
}

func defaultChoicesValueFunc() func(choice interface{}) string {
	return func(choice interface{}) string {
		return fmt.Sprintf("%s", choice)
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
		switch {

		// These keys should exit the program.
		case key.Matches(msg, m.keys.Quit):
			m.Interrupted = true
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case key.Matches(msg, m.keys.Down):
			if m.cursor < m.ItemsOnPage()-1 {
				m.cursor++
			}

		// The "h" and "left" keys skip to the pre page
		case key.Matches(msg, m.keys.Left):
			if m.showPagination {
				m.Paginator.PrevPage()
				m.ResetCursor()
			}

		// The "j" and "right" keys skip to the next page
		case key.Matches(msg, m.keys.Right):
			if m.showPagination {
				m.Paginator.NextPage()
				m.ResetCursor()
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case key.Matches(msg, m.keys.Select):
			m.Selected = m.Index()
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	if m.showFilter {
		// Filter VisibleChoices
		cmd = m.handleFiltering(msg)
	}

	// Return the updated SpinnerModel to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, tea.Batch(cmd)
}

func (m SelectModel) View() string {
	// The header
	s := color.New(color.FgHiGreen).Sprint(m.Hint)

	// Filter hints or input value
	if m.showFilter && m.Selected == -1 {
		s += m.FilterInput.View() + "\n"
	} else {
		s += "\n"
	}

	// Has user selected?
	if m.Selected != -1 {
		checked := "x" // selected!
		cursor := ">"
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, m.VisibleChoices[m.Selected])
		return s
	}

	// Iterate over our VisibleChoices
	var start, end int
	if m.showPagination {
		start, end = m.Paginator.GetSliceBounds(len(m.VisibleChoices))
	} else {
		start = 0
		end = len(m.VisibleChoices)
	}
	for i, choice := range m.VisibleChoices[start:end] {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		if m.Selected == -1 {
			s += fmt.Sprintf("%s [ ] %s\n", cursor, choice)
		}
	}

	// Show paginator dot
	if m.showPagination {
		// Only show when there are more than one page
		if m.Paginator.TotalPages > 1 {
			s += m.Paginator.View() + "\n"
			s += m.helper.View(m.keys)
		}
	}

	// Send the UI for rendering
	return s
}

// GetSelectedItem returns the selected item in VisibleChoices.
func (m SelectModel) GetSelectedItem() interface{} {
	return m.VisibleChoices[m.Selected]
}

func (m *SelectModel) EnableFilter() {
	m.showFilter = true
	m.FilterInput.Focus()
}

// SetChoiceValue set the ChoiceValue function which is used to get value of
// item in filter process. When you EnableFilter, you should set a special
// ChoiceValue function for your own Choices.
func (m *SelectModel) SetChoiceValue(f func(choice interface{}) string) {
	m.ChoiceValue = f
}

// EnablePagination shows pagination in terminal.
// itemsPerPage is used to set the number of items displayed per page.
func (m *SelectModel) EnablePagination(itemsPerPage int) {
	m.showPagination = true
	m.Paginator.PerPage = itemsPerPage
	m.Paginator.SetTotalPages(len(m.VisibleChoices))
}

func (m *SelectModel) ResetCursor() {
	m.cursor = 0
}

// ItemsOnPage returns the numer of items on the current page given the
// total numer of items passed as an argument.
func (m SelectModel) ItemsOnPage() int {
	return m.Paginator.ItemsOnPage(len(m.VisibleChoices))
}

// Index returns the index of the currently selected item as it appears
// in the VisibleChoices.
func (m SelectModel) Index() int {
	return m.cursor + m.Paginator.Page*m.Paginator.PerPage
}

// handleFiltering updates VisiableChoices ,cusor and paginator
// when a user is typing for filter.
func (m *SelectModel) handleFiltering(msg tea.Msg) tea.Cmd {
	// Update the filter text input component
	newFilterInputModel, cmd := m.FilterInput.Update(msg)
	filterChanged := m.FilterInput.Value() != newFilterInputModel.Value()
	m.FilterInput = newFilterInputModel

	// If the filtering input has changed, filtering
	if filterChanged {
		filterItems(m)
		// Update pagination
		m.resetPagination()
	}

	return tea.Batch(cmd)
}

func (m *SelectModel) resetPagination() {
	items := len(m.VisibleChoices)
	if items < 1 { // at least one page
		items = 1
	}
	m.Paginator.SetTotalPages(items)
	m.Paginator.Page = 0
	m.cursor = 0
}

func filterItems(m *SelectModel) {
	if m.FilterInput.Value() == "" {
		m.VisibleChoices = m.Choices
		return
	}
	// get Choices value
	targets := make([]string, len(m.Choices))
	for i, v := range m.Choices {
		targets[i] = m.ChoiceValue(v)
	}
	r := filter(m.FilterInput.Value(), targets)
	// reset VisibleChoices
	ChoicesAfterFilter := make([]interface{}, len(r))
	for i, v := range r {
		ChoicesAfterFilter[i] = m.Choices[v]
	}
	m.VisibleChoices = ChoicesAfterFilter
}

// filter returns index list
func filter(term string, targets []string) []int {
	var ranks = fuzzy.Find(term, targets)
	result := make([]int, len(ranks))
	for i, r := range ranks {
		result[i] = r.Index
	}
	return result
}
