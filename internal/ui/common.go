package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type TickMsg time.Time

// finalPause prevent the progress bar from exiting before it reaches 100%.
// See https://github.com/charmbracelet/bubbletea/blob/702b43d6b06287363b72836c88be35d985624a2b/examples/progress-download/tui.go#L23
func finalPause() tea.Cmd {
	return tea.Tick(time.Second*1, func(_ time.Time) tea.Msg {
		return nil
	})
}
