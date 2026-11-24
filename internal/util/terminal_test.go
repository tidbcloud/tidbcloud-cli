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

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeTerminalText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"plain name", "export-2026.tar.gz", "export-2026.tar.gz"},
		{"OSC8 hyperlink", "file\x1b]8;;http://evil/\x07name.tar.gz", "filename.tar.gz"},
		{"OSC52 clipboard write", "file\x1b]52;c;Zm9v\x07.tar.gz", "file.tar.gz"},
		{"cursor control", "name\r.tar\n.gz\x1b[2K", "name.tar.gz"},
		{"bell and DEL", "a\x07b\x7fc", "abc"},
		{"C1 control", "a\u0085b\u009fc", "abc"},
		{"newline and tab", "line1\nline2\tend", "line1line2end"},
		{"unicode preserved", "导出-2026.tar.gz", "导出-2026.tar.gz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, SanitizeTerminalText(tt.input))
		})
	}
}
