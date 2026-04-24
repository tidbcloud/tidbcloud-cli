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
	"strings"
	"testing"
)

func TestCreateFileAllowsPathInsideBase(t *testing.T) {
	base := t.TempDir()

	file, err := CreateFile(base, "export.sql.gz")
	if err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}
	file.Close()

	if _, err := os.Stat(filepath.Join(base, "export.sql.gz")); err != nil {
		t.Fatalf("expected file inside base: %v", err)
	}
}

func TestCreateFileAllowsExistingSubdirectoryInsideBase(t *testing.T) {
	base := t.TempDir()
	if err := os.Mkdir(filepath.Join(base, "schema"), 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	file, err := CreateFile(base, filepath.Join("schema", "export.sql.gz"))
	if err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}
	file.Close()

	if _, err := os.Stat(filepath.Join(base, "schema", "export.sql.gz")); err != nil {
		t.Fatalf("expected nested file inside base: %v", err)
	}
}

func TestCreateFileRejectsPathOutsideBase(t *testing.T) {
	base := t.TempDir()
	outside := filepath.Join(base, "..", "outside.sql.gz")

	tests := []string{
		"",
		".",
		"..",
		filepath.Join("..", "outside.sql.gz"),
		filepath.Join("nested", "..", "..", "outside.sql.gz"),
		outside,
		"C:\\Temp\\outside.sql.gz",
		"\\\\server\\share\\outside.sql.gz",
		"bad\x00name.sql.gz",
		"bad\nname.sql.gz",
	}

	for _, tt := range tests {
		t.Run(strings.ReplaceAll(tt, string(os.PathSeparator), "_"), func(t *testing.T) {
			file, err := CreateFile(base, tt)
			if err == nil {
				file.Close()
				t.Fatalf("CreateFile() succeeded for %q", tt)
			}
		})
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
