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
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenBrowserWindowsAvoidsShell(t *testing.T) {
	if runtime.GOOS != "windows" {
		// Only resolve a fake executable; never start a browser in unit tests.
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "rundll32"), []byte("#!/bin/sh\nexit 0\n"), 0700))
		t.Setenv("PATH", dir)
	}
	url := "https://example.com/verify?code=abc&scope|<^>%"
	cmd := OpenBrowser("windows", url)
	require.NoError(t, cmd.Err)
	require.Equal(t, []string{"url.dll,FileProtocolHandler", url}, cmd.Args[1:])
}

func TestOpenBrowserWindowsReportsMissingExecutable(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	cmd := OpenBrowser("windows", "https://example.com/verify")
	require.Error(t, cmd.Err)
	require.Contains(t, cmd.Err.Error(), "rundll32")
}

func TestOpenBrowserPassesURLUnmodified(t *testing.T) {
	url := "https://example.com/verify?code=abc"
	cmd := OpenBrowser(runtime.GOOS, url)
	require.Equal(t, url, cmd.Args[len(cmd.Args)-1])
}
