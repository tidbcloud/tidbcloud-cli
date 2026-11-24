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
	"regexp"
	"strings"
)

// Terminal escape sequences are stripped whole so that neither the sequence
// nor its payload (for example an OSC8 hyperlink URL) reaches the output.
var (
	oscSequence = regexp.MustCompile(`\x1b\][^\x07\x1b]*(\x07|\x1b\\)`)
	csiSequence = regexp.MustCompile(`\x1b[[\][0-9;?]*[ -/]*[@-~]`)
)

// SanitizeTerminalText removes ANSI escape sequences and C0/C1 control
// characters from server-provided text before printing, preventing terminal
// escape injection (OSC8 hyperlinks, OSC52 clipboard writes, cursor control)
// through file names or other API-controlled strings. The injection happens
// even when the value is later rejected, so sanitize at every output site.
func SanitizeTerminalText(s string) string {
	s = oscSequence.ReplaceAllString(s, "")
	s = csiSequence.ReplaceAllString(s, "")
	return strings.Map(func(r rune) rune {
		if r < 0x20 || (r >= 0x7f && r <= 0x9f) {
			return -1
		}
		return r
	}, s)
}
