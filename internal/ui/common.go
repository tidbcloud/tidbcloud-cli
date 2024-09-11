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
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

// FinalPause prevent the progress bar from exiting before it reaches 100%.
// See https://github.com/charmbracelet/bubbletea/blob/702b43d6b06287363b72836c88be35d985624a2b/examples/progress-download/tui.go#L23
func FinalPause() tea.Cmd {
	return tea.Tick(time.Second*1, func(_ time.Time) tea.Msg {
		return nil
	})
}
