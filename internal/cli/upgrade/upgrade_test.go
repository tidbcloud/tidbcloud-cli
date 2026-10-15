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
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/github"

	"github.com/stretchr/testify/require"
)

func TestInstallerChecksum(t *testing.T) {
	const script = "#!/bin/sh\nexit 0\n"
	const checksum = "306c6ca7407560340797866e077e053627ad409277d1b9da58106fce4cf717cb"
	tests := []struct {
		name     string
		script   string
		checksum string
		wantErr  bool
	}{
		{name: "matching checksum", script: script, checksum: checksum},
		{name: "changed script", script: script + "\n", checksum: checksum, wantErr: true},
		{name: "incorrect checksum", script: script, checksum: strings.Repeat("0", 64), wantErr: true},
		{name: "empty script", checksum: checksum, wantErr: true},
		{name: "malformed checksum", script: script, checksum: "invalid", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyInstallerChecksum([]byte(tt.script), tt.checksum)
			if tt.wantErr {
				require.EqualError(t, err, "checksum verification failed for install.sh")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpgradeDownload(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the shell installer supports macOS and Linux")
	}
	tests := []struct {
		name     string
		exitCode string
		wantErr  string
	}{
		{name: "checksum mismatch", exitCode: "0", wantErr: "checksum verification failed for install.sh"},
		{name: "download fails", exitCode: "22", wantErr: "exit status 22"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			argsPath := filepath.Join(dir, "args")
			marker := filepath.Join(dir, "installed")
			curl := `#!/bin/sh
printf '%s\n' "$@" > "$INSTALLER_TEST_ARGS"
printf '%s' "$INSTALLER_TEST_SCRIPT"
exit "$INSTALLER_TEST_EXIT_CODE"
`
			require.NoError(t, os.WriteFile(filepath.Join(dir, "curl"), []byte(curl), 0700))
			t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
			t.Setenv("INSTALLER_TEST_ARGS", argsPath)
			t.Setenv("INSTALLER_TEST_SCRIPT", "printf installed > \"$INSTALLER_TEST_MARKER\"\n")
			t.Setenv("INSTALLER_TEST_MARKER", marker)
			t.Setenv("INSTALLER_TEST_EXIT_CODE", tt.exitCode)
			h := &internal.Helper{IOStreams: iostream.Test()}
			err := updateAndWaitReady(context.Background(), h, &github.ReleaseInfo{Version: "v1.2.3"})
			require.EqualError(t, err, tt.wantErr)
			require.NoFileExists(t, marker)
			require.Equal(t, "... Upgrading the CLI to version v1.2.3\n", h.IOStreams.Out.(*bytes.Buffer).String())
			require.Empty(t, h.IOStreams.Err.(*bytes.Buffer).String())
			args, err := os.ReadFile(argsPath)
			require.NoError(t, err)
			require.Equal(t, "-q\n--proto\n=https\n--proto-redir\n=https\n-fsSL\n"+installerURL+"\n", string(args))
			require.Regexp(t, regexp.MustCompile(`^https://raw\.githubusercontent\.com/tidbcloud/tidbcloud-cli/[0-9a-f]{40}/install\.sh$`), installerURL)
		})
	}
}
