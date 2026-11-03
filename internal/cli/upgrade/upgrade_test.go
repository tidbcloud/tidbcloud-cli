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

package upgrade

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
)

func TestUpgradeDoesNotRunDownloads(t *testing.T) {
	previous := config.IsUnderTiUP
	config.IsUnderTiUP = false
	t.Cleanup(func() { config.IsUnderTiUP = previous })
	dir := t.TempDir()
	marker := filepath.Join(dir, "curl-ran")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "curl"), []byte("#!/bin/sh\ntouch \"$UPGRADE_TEST_MARKER\"\n"), 0700))
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("UPGRADE_TEST_MARKER", marker)
	for _, prompt := range []bool{false, true} {
		h := &internal.Helper{IOStreams: iostream.Test()}
		h.IOStreams.CanPrompt = prompt
		cmd := Cmd(h)
		cmd.SilenceErrors, cmd.SilenceUsage = true, true
		require.ErrorContains(t, cmd.Execute(), "automatic upgrade is unavailable")
		require.NoFileExists(t, marker)
	}
}

func TestUpgradeUnderTiUP(t *testing.T) {
	previous := config.IsUnderTiUP
	config.IsUnderTiUP = true
	t.Cleanup(func() { config.IsUnderTiUP = previous })
	cmd := Cmd(&internal.Helper{IOStreams: iostream.Test()})
	cmd.SilenceErrors, cmd.SilenceUsage = true, true
	require.ErrorContains(t, cmd.Execute(), "tiup update cloud")
}
