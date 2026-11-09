// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSetStateEntryPermissions(t *testing.T) {
	for _, existing := range []bool{false, true} {
		name := "new directory"
		if existing {
			name = "existing private directory"
		}
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), ".ticloud")
			if existing {
				require.NoError(t, os.Mkdir(dir, 0700))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "config"), []byte("dummy-config"), 0600))
			}
			file := filepath.Join(dir, "state.yml")
			checked := time.Now().UTC().Truncate(time.Second)
			require.NoError(t, setStateEntry(file, checked, ReleaseInfo{Version: "v1.0.0"}))
			state, err := getStateEntry(file)
			require.NoError(t, err)
			require.Equal(t, "v1.0.0", state.LatestRelease.Version)
			require.True(t, checked.Equal(state.CheckedForUpdateAt))
			if runtime.GOOS == "windows" {
				return
			}
			info, err := os.Stat(dir)
			require.NoError(t, err)
			require.Zero(t, info.Mode().Perm() & ^os.FileMode(0700))
			info, err = os.Stat(file)
			require.NoError(t, err)
			require.Zero(t, info.Mode().Perm() & ^os.FileMode(0600))
			if existing {
				info, err = os.Stat(filepath.Join(dir, "config"))
				require.NoError(t, err)
				require.Zero(t, info.Mode().Perm() & ^os.FileMode(0600))
			}
		})
	}
}
