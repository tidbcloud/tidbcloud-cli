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
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type downloadTestTransport func(*http.Request) (*http.Response, error)

func (f downloadTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetResponseDestinations(t *testing.T) {
	const origin = "https://storage.example/start?X-Amz-Signature=test%2Bsignature"
	const final = "https://storage.example/file"
	const body = "synthetic-export-data"
	for _, status := range []int{301, 302, 303, 307, 308} {
		for _, tc := range []struct {
			name         string
			initial      string
			redirects    map[string]string
			wantRequests int
			wantError    bool
		}{
			{name: "HTTPS", initial: origin, wantRequests: 1},
			{name: "HTTP", initial: "http://storage.example/file", wantError: true},
			{name: "other scheme", initial: "ftp://storage.example/file", wantError: true},
			{name: "missing host", initial: "https:///file", wantError: true},
			{name: "HTTP loopback", initial: "http://127.0.0.1:8080/file", wantError: true},
			{name: "relative redirect", initial: origin, redirects: map[string]string{origin: "/file"}, wantRequests: 2},
			{name: "same origin", initial: origin, redirects: map[string]string{origin: final}, wantRequests: 2},
			{name: "equivalent origin", initial: origin, redirects: map[string]string{origin: "https://STORAGE.example:443/file"}, wantRequests: 2},
			{name: "cross host", initial: origin, redirects: map[string]string{origin: "https://attacker.example/file"}, wantRequests: 1, wantError: true},
			{name: "cross port", initial: origin, redirects: map[string]string{origin: "https://storage.example:8443/file"}, wantRequests: 1, wantError: true},
			{name: "downgrade", initial: origin, redirects: map[string]string{origin: "http://storage.example/file"}, wantRequests: 1, wantError: true},
			{name: "loopback redirect", initial: origin, redirects: map[string]string{origin: "http://127.0.0.1:8080/file"}, wantRequests: 1, wantError: true},
			{name: "link local redirect", initial: origin, redirects: map[string]string{origin: "http://169.254.169.254/file"}, wantRequests: 1, wantError: true},
			{name: "second hop", initial: origin, redirects: map[string]string{origin: final, final: "https://attacker.example/file"}, wantRequests: 2, wantError: true},
			{name: "redirect loop", initial: origin, redirects: map[string]string{origin: origin}, wantRequests: 10, wantError: true},
		} {
			t.Run(fmt.Sprintf("%d/%s", status, tc.name), func(t *testing.T) {
				var reached []string
				original := http.DefaultTransport
				// Replace network I/O, retaining the production HTTPS transport
				// and redirect policy in GetResponse. No network is used.
				http.DefaultTransport = downloadTestTransport(func(req *http.Request) (*http.Response, error) {
					reached = append(reached, req.URL.String())
					require.Equal(t, http.MethodGet, req.Method)
					header := http.Header{}
					code := http.StatusOK
					if location, ok := tc.redirects[req.URL.String()]; ok {
						header.Set("Location", location)
						code = status
					}
					return &http.Response{
						StatusCode: code, Header: header, ContentLength: int64(len(body)),
						Body: io.NopCloser(strings.NewReader(body)), Request: req,
					}, nil
				})
				t.Cleanup(func() { http.DefaultTransport = original })
				resp, err := GetResponse(tc.initial, false)
				require.Len(t, reached, tc.wantRequests)
				if len(reached) > 0 {
					require.Equal(t, tc.initial, reached[0])
				}
				if tc.wantError {
					require.Error(t, err)
					require.Nil(t, resp)
					return
				}
				require.NoError(t, err)
				defer resp.Body.Close()
				got, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				require.Equal(t, body, string(got))
			})
		}
	}
}

func TestCreateFileAllowsRelativeDestinations(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "normal file name",
			fileName: "a.sql.gz",
		},
		{
			name:     "subdirectory",
			fileName: filepath.Join("folder", "a.sql.gz"),
		},
		{
			name:     "nested subdirectory",
			fileName: filepath.Join("a", "b", "c.sql.gz"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			parent := filepath.Dir(filepath.Join(base, tt.fileName))
			if err := os.MkdirAll(parent, 0755); err != nil {
				t.Fatalf("MkdirAll() error = %v", err)
			}

			file, err := CreateFile(base, tt.fileName)
			if err != nil {
				t.Fatalf("CreateFile() error = %v", err)
			}
			file.Close()

			if _, err := os.Stat(filepath.Join(base, tt.fileName)); err != nil {
				t.Fatalf("expected file inside base: %v", err)
			}
		})
	}
}

func TestCreateFileRejectsUnsupportedDestinations(t *testing.T) {
	base := t.TempDir()

	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "path outside output path",
			fileName: filepath.Join("..", "..", "tmp", "pwned"),
		},
		{
			name:     "absolute path",
			fileName: filepath.Join(string(os.PathSeparator), "tmp", "pwned"),
		},
		{
			name:     "empty file name",
			fileName: "",
		},
		{
			name:     "current directory",
			fileName: ".",
		},
		{
			name:     "parent directory",
			fileName: "..",
		},
		{
			name:     "cleans to current directory",
			fileName: filepath.Join("a", ".."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := CreateFile(base, tt.fileName)
			if err == nil {
				file.Close()
				t.Fatalf("CreateFile() succeeded for %q", tt.fileName)
			}
		})
	}
}

func TestCreateFileDoesNotOverwriteExistingDestination(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "a.sql.gz")
	if err := os.WriteFile(path, []byte("existing"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	file, err := CreateFile(base, "a.sql.gz")
	if err == nil {
		file.Close()
		t.Fatalf("CreateFile() succeeded for existing destination")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "existing" {
		t.Fatalf("existing destination was overwritten: %q", string(content))
	}
}

func TestCreateFileUsesCurrentDirectoryWhenBaseIsEmpty(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	base := t.TempDir()
	if err := os.Chdir(base); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	file, err := CreateFile("", "a.sql.gz")
	if err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}
	file.Close()

	if _, err := os.Stat(filepath.Join(base, "a.sql.gz")); err != nil {
		t.Fatalf("expected file in current directory: %v", err)
	}
}

func TestCreateFileRejectsSymlinkedParentOutsideBase(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink behavior requires additional privileges on Windows")
	}

	base := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(base, "link")); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}

	file, err := CreateFile(base, filepath.Join("link", "export.sql.gz"))
	if err == nil {
		file.Close()
		t.Fatalf("CreateFile() succeeded through symlinked parent")
	}
	if _, statErr := os.Stat(filepath.Join(outside, "export.sql.gz")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no file outside base, stat error = %v", statErr)
	}
}
