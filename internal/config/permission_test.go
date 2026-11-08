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

package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRestrictPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are not supported on Windows")
	}
	for _, tc := range []struct {
		name  string
		dir   bool
		mode  fs.FileMode
		limit fs.FileMode
		want  fs.FileMode
	}{
		{name: "world-readable file", mode: 0644, limit: 0600, want: 0600},
		{name: "group-writable file", mode: 0660, limit: 0600, want: 0600},
		{name: "already private file", mode: 0600, limit: 0600, want: 0600},
		{name: "read-only file stays read-only", mode: 0400, limit: 0600, want: 0400},
		{name: "world-readable directory", dir: true, mode: 0755, limit: 0700, want: 0700},
		{name: "already private directory", dir: true, mode: 0700, limit: 0700, want: 0700},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "target")
			if tc.dir {
				require.NoError(t, os.Mkdir(path, 0700))
			} else {
				require.NoError(t, os.WriteFile(path, []byte("public-key = 'test'\n"), 0600))
			}
			// Set the mode explicitly so the process umask does not affect it.
			require.NoError(t, os.Chmod(path, tc.mode))

			require.NoError(t, RestrictPermissions(path, tc.limit))

			info, err := os.Stat(path)
			require.NoError(t, err)
			require.Equal(t, tc.want, info.Mode().Perm())
		})
	}
}

func TestRestrictPermissionsMissingPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are not supported on Windows")
	}
	err := RestrictPermissions(filepath.Join(t.TempDir(), "missing"), 0600)
	require.ErrorIs(t, err, fs.ErrNotExist)
}
