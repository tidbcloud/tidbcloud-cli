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
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/github"

	"github.com/stretchr/testify/require"
)

func TestInstallerChecksum(t *testing.T) {
	const script = "#!/bin/sh\nexit 0\n"
	const checksum = "306c6ca7407560340797866e077e053627ad409277d1b9da58106fce4cf717cb"
	const entry = checksum + "  install.sh\n"
	tests := []struct {
		name     string
		script   string
		manifest string
		wantErr  string
	}{
		{name: "matching checksum", script: script, manifest: entry},
		{name: "binary checksum marker", script: script, manifest: checksum + " *install.sh\n"},
		{name: "uppercase checksum", script: script, manifest: strings.ToUpper(checksum) + "  install.sh\n"},
		{name: "other assets", script: script, manifest: checksum + "  ticloud.tar.gz\n" + entry},
		{name: "changed script", script: script + "\n", manifest: entry, wantErr: "checksum verification failed for install.sh"},
		{name: "incorrect checksum", script: script, manifest: strings.Repeat("0", 64) + "  install.sh\n", wantErr: "checksum verification failed for install.sh"},
		{name: "empty script", manifest: entry, wantErr: "checksum verification failed for install.sh"},
		{name: "missing checksum", script: script, manifest: checksum + "  other.sh\n", wantErr: "missing checksum for install.sh"},
		{name: "duplicate checksum", script: script, manifest: entry + entry, wantErr: "duplicate checksum for install.sh"},
		{name: "malformed checksum", script: script, manifest: "invalid  install.sh\n", wantErr: "invalid checksum for install.sh"},
		{name: "short checksum", script: script, manifest: "0123  install.sh\n", wantErr: "invalid checksum for install.sh"},
		{name: "empty manifest", script: script, wantErr: "missing checksum for install.sh"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyInstallerChecksum([]byte(tt.script), []byte(tt.manifest))
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
	err := verifyInstallerChecksum([]byte(script), bytes.Repeat([]byte("x"), 65537))
	require.ErrorContains(t, err, "failed to read checksums.txt")
}

func TestUpgradeDownload(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the shell installer supports macOS and Linux")
	}
	tests := []struct {
		name      string
		version   string
		manifest  string
		failAsset string
		wantErr   string
	}{
		{name: "selected release", version: "v1.2.3"},
		{name: "version without prefix", version: "1.2.3"},
		{name: "newer release", version: "v2.3.4"},
		{name: "prerelease", version: "v2.3.4-beta.1"},
		{name: "checksum mismatch", version: "v1.2.3", manifest: "mismatch", wantErr: "checksum verification failed for install.sh"},
		{name: "missing checksum", version: "v1.2.3", manifest: "missing", wantErr: "missing checksum for install.sh"},
		{name: "installer download fails", version: "v1.2.3", failAsset: "install.sh", wantErr: "failed to download install.sh for v1.2.3"},
		{name: "manifest download fails", version: "v1.2.3", failAsset: "checksums.txt", wantErr: "failed to download checksums.txt for v1.2.3"},
		{name: "latest is not a release", version: "latest", wantErr: "invalid release version"},
		{name: "empty version", wantErr: "invalid release version"},
		{name: "invalid version", version: "v1.2.3/other", wantErr: "invalid release version"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const script = "printf '%s' \"$TICLOUD_VERSION\" > \"$INSTALLER_TEST_MARKER\"\n"
			manifest := fmt.Sprintf("%x  install.sh\n", sha256.Sum256([]byte(script)))
			switch tt.manifest {
			case "mismatch":
				manifest = strings.Repeat("0", 64) + "  install.sh\n"
			case "missing":
				manifest = ""
			}
			marker, requests, argsPath := mockInstallerDownloads(t, script, manifest, tt.failAsset)
			t.Setenv("TICLOUD_VERSION", "0.0.1")
			h := &internal.Helper{IOStreams: iostream.Test()}
			err := updateAndWaitReady(context.Background(), h, &github.ReleaseInfo{Version: tt.version})
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				require.NoFileExists(t, marker)
			} else {
				require.NoError(t, err)
				installedVersion, err := os.ReadFile(marker)
				require.NoError(t, err)
				require.Equal(t, strings.TrimPrefix(tt.version, "v"), string(installedVersion))
			}
			require.Equal(t, fmt.Sprintf("... Upgrading the CLI to version %s\n", tt.version), h.IOStreams.Out.(*bytes.Buffer).String())
			require.Empty(t, h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.wantErr == "invalid release version" {
				require.NoFileExists(t, requests)
				return
			}
			base := releaseDownloadURL + "/v" + strings.TrimPrefix(tt.version, "v")
			expectedRequests := base + "/install.sh\n"
			expectedArgs := "-q\n--proto\n=https\n--proto-redir\n=https\n-fsSL\n" + base + "/install.sh\n"
			if tt.failAsset != "install.sh" {
				expectedRequests += base + "/checksums.txt\n"
				expectedArgs += "-q\n--proto\n=https\n--proto-redir\n=https\n-fsSL\n" + base + "/checksums.txt\n"
			}
			requestLog, err := os.ReadFile(requests)
			require.NoError(t, err)
			require.Equal(t, expectedRequests, string(requestLog))
			args, err := os.ReadFile(argsPath)
			require.NoError(t, err)
			require.Equal(t, expectedArgs, string(args))
		})
	}
}

func TestUpgradeStopsWhenVersionLookupFails(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the shell installer supports macOS and Linux")
	}
	previous := config.IsUnderTiUP
	config.IsUnderTiUP = false
	t.Cleanup(func() { config.IsUnderTiUP = previous })
	marker, requests, _ := mockInstallerDownloads(t, "", "", "latest-version")
	h := &internal.Helper{IOStreams: iostream.Test()}
	cmd := Cmd(h)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	err := cmd.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "failed to check for CLI updates")
	require.Empty(t, h.IOStreams.Out.(*bytes.Buffer).String())
	require.NoFileExists(t, marker)
	requestLog, err := os.ReadFile(requests)
	require.NoError(t, err)
	require.Equal(t, "https://raw.githubusercontent.com/tidbcloud/tidbcloud-cli/main/latest-version\n", string(requestLog))
}

func mockInstallerDownloads(t *testing.T, script, manifest, failAsset string) (marker, requests, argsPath string) {
	t.Helper()
	dir := t.TempDir()
	marker = filepath.Join(dir, "installed")
	requests = filepath.Join(dir, "requests")
	argsPath = filepath.Join(dir, "args")
	curl := `#!/bin/sh
printf '%s\n' "$@" >> "$INSTALLER_TEST_ARGS"
for argument; do url=$argument; done
printf '%s\n' "$url" >> "$INSTALLER_TEST_REQUESTS"
asset=${url##*/}
[ "$asset" != "$INSTALLER_TEST_FAIL_ASSET" ] || exit 22
case "$asset" in
    install.sh) printf '%s' "$INSTALLER_TEST_SCRIPT" ;;
    checksums.txt) printf '%s' "$INSTALLER_TEST_MANIFEST" ;;
    *) exit 22 ;;
esac
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "curl"), []byte(curl), 0700))
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("INSTALLER_TEST_ARGS", argsPath)
	t.Setenv("INSTALLER_TEST_REQUESTS", requests)
	t.Setenv("INSTALLER_TEST_SCRIPT", script)
	t.Setenv("INSTALLER_TEST_MANIFEST", manifest)
	t.Setenv("INSTALLER_TEST_MARKER", marker)
	t.Setenv("INSTALLER_TEST_FAIL_ASSET", failAsset)
	return marker, requests, argsPath
}
