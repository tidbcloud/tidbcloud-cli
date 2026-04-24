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
)

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
