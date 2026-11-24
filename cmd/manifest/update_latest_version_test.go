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

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Exercise the actual file update, including the release-script validation,
// without GitHub access or changes to the checkout's installer.
func TestUpdateManifest(t *testing.T) {
	root, err := filepath.Abs("../..")
	require.NoError(t, err)
	installer := "#!/bin/sh\nversion=\"0.1.0\"\necho \"$version\"\n"
	for _, tag := range []string{
		"v1", "v1.2.3", "v2.3.4-rc.1+dev", "v1.0.0-nightly-20261123",
		"", "v", "1.2.3", "latest", "vfoo", "v-1.2.3", "v1.2.3/other",
		"v1.2.3;id", "v1.2.3\"", "v1.2.3$(id)", "v1.2.3`id`", "v1.2.3\nother", "v1.2.3 other",
	} {
		t.Run(tag, func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)
			require.NoError(t, os.WriteFile("install.sh", []byte(installer), 0600))
			require.NoError(t, os.WriteFile("latest-version", []byte("v0.1.0"), 0600))
			// The first four entries are existing supported release formats.
			valid := tag == "v1" || tag == "v1.2.3" || tag == "v2.3.4-rc.1+dev" || tag == "v1.0.0-nightly-20261123"
			err := updateManifest(tag)
			gotInstaller, readErr := os.ReadFile("install.sh")
			require.NoError(t, readErr)
			latest, readErr := os.ReadFile("latest-version")
			require.NoError(t, readErr)
			if !valid {
				require.Error(t, err)
				require.Equal(t, installer, string(gotInstaller))
				require.Equal(t, "v0.1.0", string(latest))
			} else {
				require.NoError(t, err)
				require.Equal(t, tag, string(latest))
				require.Equal(t, strings.Replace(installer, "0.1.0", strings.TrimPrefix(tag, "v"), 1), string(gotInstaller))
			}
			if runtime.GOOS == "windows" {
				return
			}
			output := filepath.Join(dir, "release", "install.sh")
			command := exec.Command("/bin/sh", filepath.Join(root, "scripts/prepare-installer.sh"), strings.TrimPrefix(tag, "v"), output)
			log, prepareErr := command.CombinedOutput()
			// The release script also accepts a bare version; GitHub tags
			// intentionally require the leading v.
			if valid || tag == "1.2.3" {
				require.NoError(t, prepareErr, string(log))
			} else {
				require.Error(t, prepareErr, string(log))
				require.NoFileExists(t, output)
			}
		})
	}
}

func TestUpdateManifestRejectsUnexpectedInstaller(t *testing.T) {
	for _, installer := range []string{"#!/bin/sh\necho missing\n", "version=\"1\"\nversion=\"2\"\n"} {
		t.Run(installer, func(t *testing.T) {
			t.Chdir(t.TempDir())
			require.NoError(t, os.WriteFile("install.sh", []byte(installer), 0600))
			require.NoError(t, os.WriteFile("latest-version", []byte("v0.1.0"), 0600))
			require.ErrorContains(t, updateManifest("v1.2.3"), "exactly one version assignment")
			latest, err := os.ReadFile("latest-version")
			require.NoError(t, err)
			require.Equal(t, "v0.1.0", string(latest))
			got, err := os.ReadFile("install.sh")
			require.NoError(t, err)
			require.Equal(t, installer, string(got))
		})
	}
}
