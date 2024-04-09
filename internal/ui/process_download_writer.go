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
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var p *tea.Program

type ProgressWriter struct {
	Total      int
	Downloaded int
	File       *os.File
	Reader     io.Reader
	OnProgress func(float64)
}

func (pw *ProgressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	_, err := io.Copy(pw.File, io.TeeReader(pw.Reader, pw))
	if err != nil {
		if p != nil {
			p.Send(ProgressErrMsg{err})
		}
	}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	pw.Downloaded += len(p)
	if pw.Total > 0 && pw.OnProgress != nil {
		pw.OnProgress(float64(pw.Downloaded) / float64(pw.Total))
	}
	return len(p), nil
}
