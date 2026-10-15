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

package scripts_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstallChecksums(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("install.sh supports macOS and Linux")
	}
	script, err := filepath.Abs("../install.sh")
	require.NoError(t, err)
	source, err := os.ReadFile(script)
	require.NoError(t, err)
	match := regexp.MustCompile(`(?m)^version="([^"]+)"$`).FindSubmatch(source)
	require.Len(t, match, 2)
	defaultVersion := string(match[1])

	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)
	binary := []byte("fixture binary\n")
	require.NoError(t, tw.WriteHeader(&tar.Header{Name: "ticloud", Mode: 0755, Size: int64(len(binary))}))
	_, err = tw.Write(binary)
	require.NoError(t, err)
	require.NoError(t, tw.Close())
	require.NoError(t, gz.Close())
	digest := fmt.Sprintf("%x", sha256.Sum256(archive.Bytes()))

	tests := []struct {
		name      string
		mirror    string
		version   string
		checksum  string
		failFetch bool
		wantError string
	}{
		{name: "matching checksum"},
		{name: "binary checksum marker", checksum: "binary"},
		{name: "HTTPS mirror", mirror: "https://mirror.example.com/ticloud/"},
		{name: "selected version", version: "v1.2.3"},
		{name: "checksum mismatch", checksum: "mismatch", wantError: "Checksum verification failed"},
		{name: "missing checksum", checksum: "missing", wantError: "Expected exactly one checksum"},
		{name: "duplicate checksum", checksum: "duplicate", wantError: "Expected exactly one checksum"},
		{name: "checksum download fails", failFetch: true, wantError: "Failed to download and/or extract"},
		{name: "HTTP mirror", mirror: "http://mirror.example.com", wantError: "must use HTTPS"},
		{name: "empty mirror host", mirror: "https:///ticloud", wantError: "must have a host"},
		{name: "mirror credentials", mirror: "https://user@mirror.example.com", wantError: "without credentials"},
		{name: "mirror query", mirror: "https://mirror.example.com?file=", wantError: "must not contain"},
		{name: "mirror fragment", mirror: "https://mirror.example.com#file", wantError: "must not contain"},
		{name: "mirror whitespace", mirror: "https://mirror.example.com/a b", wantError: "must not contain"},
		{name: "mirror backslash", mirror: `https://mirror.example.com\path`, wantError: "must not contain"},
	}
	for _, tool := range []string{"sha256sum", "shasum"} {
		t.Run(tool, func(t *testing.T) {
			toolPath, err := exec.LookPath(tool)
			if err != nil {
				t.Skipf("%s is unavailable", tool)
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					dir := t.TempDir()
					commands := filepath.Join(dir, "commands")
					binDir := filepath.Join(dir, "installation", "bin")
					require.NoError(t, os.MkdirAll(commands, 0700))
					require.NoError(t, os.MkdirAll(binDir, 0700))
					for _, name := range []string{"tar", "gzip", "awk", "mktemp", "cp", "rm", "mkdir", "chmod"} {
						path, err := exec.LookPath(name)
						require.NoError(t, err)
						require.NoError(t, os.Symlink(path, filepath.Join(commands, name)))
					}
					require.NoError(t, os.Symlink(toolPath, filepath.Join(commands, tool)))
					require.NoError(t, os.WriteFile(filepath.Join(commands, "uname"), []byte("#!/bin/sh\ncase \"$1\" in -s) echo Linux ;; -m) echo x86_64 ;; esac\n"), 0700))
					curl := `#!/bin/sh
while [ "$#" -gt 0 ]; do
    case "$1" in
        https://*) url=$1 ;;
        -o) shift; destination=$1 ;;
    esac
    shift
done
printf '%s\n' "$url" >> "$INSTALL_TEST_REQUESTS"
case "$url" in
    */checksums.txt)
        [ "$INSTALL_TEST_FAIL_FETCH" != true ] || exit 22
        cp "$INSTALL_TEST_DIR/checksums.txt" "$destination"
        ;;
    */*.tar.gz) cp "$INSTALL_TEST_DIR/archive.tar.gz" "$destination" ;;
    *) exit 22 ;;
esac
`
					require.NoError(t, os.WriteFile(filepath.Join(commands, "curl"), []byte(curl), 0700))
					require.NoError(t, os.WriteFile(filepath.Join(dir, "archive.tar.gz"), archive.Bytes(), 0600))
					version := defaultVersion
					if tt.version != "" {
						version = strings.TrimPrefix(tt.version, "v")
					}
					filename := "ticloud_" + version + "_linux_amd64.tar.gz"
					entry := digest + "  " + filename + "\n"
					switch tt.checksum {
					case "binary":
						entry = digest + " *" + filename + "\n"
					case "mismatch":
						entry = strings.Repeat("0", 64) + "  " + filename + "\n"
					case "missing":
						entry = digest + "  other.tar.gz\n"
					case "duplicate":
						entry += entry
					}
					require.NoError(t, os.WriteFile(filepath.Join(dir, "checksums.txt"), []byte(entry), 0600))
					installed := filepath.Join(binDir, "ticloud")
					require.NoError(t, os.WriteFile(installed, []byte("existing binary\n"), 0700))
					requests := filepath.Join(dir, "requests")
					cmd := exec.Command("/bin/sh", script)
					cmd.Env = append(os.Environ(),
						"PATH="+binDir+string(os.PathListSeparator)+commands,
						"TICLOUD_HOME="+filepath.Dir(binDir),
						"TICLOUD_MIRRORS="+tt.mirror,
						"TICLOUD_VERSION="+tt.version,
						"INSTALL_TEST_DIR="+dir,
						"INSTALL_TEST_REQUESTS="+requests,
						fmt.Sprintf("INSTALL_TEST_FAIL_FETCH=%t", tt.failFetch),
					)
					output, err := cmd.CombinedOutput()
					got, readErr := os.ReadFile(installed)
					require.NoError(t, readErr)
					if tt.wantError != "" {
						require.Error(t, err, string(output))
						require.Contains(t, string(output), tt.wantError)
						require.Equal(t, "existing binary\n", string(got))
						return
					}
					require.NoError(t, err, string(output))
					require.Equal(t, binary, got)
					requestLog, err := os.ReadFile(requests)
					require.NoError(t, err)
					repo := "https://github.com/tidbcloud/tidbcloud-cli/releases/download"
					archiveRepo := repo
					if tt.mirror != "" {
						archiveRepo = strings.TrimSuffix(tt.mirror, "/")
					}
					require.Equal(t, archiveRepo+"/v"+version+"/"+filename+"\n"+repo+"/v"+version+"/checksums.txt\n", string(requestLog))
				})
			}
		})
	}
}

func TestPrepareInstaller(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the release installer is prepared with /bin/sh")
	}
	repo, err := filepath.Abs("..")
	require.NoError(t, err)
	sourcePath := filepath.Join(repo, "install.sh")
	source, err := os.ReadFile(sourcePath)
	require.NoError(t, err)
	for _, version := range []string{"1.2.3", "v2.3.4-beta.1"} {
		t.Run(version, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "release", "install.sh")
			cmd := exec.Command("/bin/sh", "scripts/prepare-installer.sh", version, output)
			cmd.Dir = repo
			log, err := cmd.CombinedOutput()
			require.NoError(t, err, string(log))
			installer, err := os.ReadFile(output)
			require.NoError(t, err)
			require.Contains(t, string(installer), "\nversion=\""+strings.TrimPrefix(version, "v")+"\"\n")
			require.Contains(t, string(installer), "version=${TICLOUD_VERSION#v}")
			currentSource, err := os.ReadFile(sourcePath)
			require.NoError(t, err)
			require.Equal(t, source, currentSource)
		})
	}
	for _, version := range []string{"", "latest", "1.2.3/other", "1.2.3\nother"} {
		t.Run("invalid "+version, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "install.sh")
			cmd := exec.Command("/bin/sh", "scripts/prepare-installer.sh", version, output)
			cmd.Dir = repo
			log, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(log), "Invalid release version")
			require.NoFileExists(t, output)
		})
	}
}
