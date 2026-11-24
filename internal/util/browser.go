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
	exec "golang.org/x/sys/execabs"
)

// OpenBrowser opens a web browser at the specified url.
func OpenBrowser(goos, url string) *exec.Cmd {
	exe := "open"
	var args []string
	switch goos {
	case "darwin":
		args = append(args, url)
	case "windows":
		// rundll32 with FileProtocolHandler opens the URL through the shell
		// without a command interpreter, so URL metacharacters such as |, <,
		// >, ^ and % cannot be interpreted as shell syntax.
		exe = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler", url)
	default:
		exe = linuxExe()
		args = append(args, url)
	}

	cmd := exec.Command(exe, args...)
	return cmd
}

func linuxExe() string {
	exe := "xdg-open"

	_, err := exec.LookPath(exe)
	if err != nil {
		_, err := exec.LookPath("wslview")
		if err == nil {
			exe = "wslview"
		}
	}

	return exe
}
