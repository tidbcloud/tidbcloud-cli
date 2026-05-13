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

package security

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRawHTTPDumpUsageIsCentralized(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test file")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
	allowed := filepath.Join(repoRoot, "pkg/tidbcloud/redact/redact.go")
	needles := []string{"httputil." + "DumpRequestOut", "httputil." + "DumpResponse"}

	for _, dir := range []string{"internal", "pkg/tidbcloud"} {
		root := filepath.Join(repoRoot, dir)
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for _, needle := range needles {
				if strings.Contains(string(content), needle) && path != allowed {
					t.Fatalf("raw HTTP dump %q must go through redaction helper, found in %s", needle, path)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
