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
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Exercise the real Makefile cleanup and restoration steps in isolation. Stub
// the external generator and formatters; no Java, downloads or source edits.
func TestNextGenGenerationRestoresExtensions(t *testing.T) {
	makePath, err := exec.LookPath("make")
	if err != nil {
		t.Skip("make is unavailable")
	}
	root, err := filepath.Abs("..")
	require.NoError(t, err)
	workspace := t.TempDir()
	write := func(name string, data []byte, mode os.FileMode) {
		path := filepath.Join(workspace, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
		require.NoError(t, os.WriteFile(path, data, mode))
	}
	makefile, err := os.ReadFile(filepath.Join(root, "Makefile"))
	require.NoError(t, err)
	write("Makefile", makefile, 0644)
	extension := "tools/openapi-generator/go/extensions/nextgen/private_endpoint_connection.go.tmpl"
	source, err := os.ReadFile(filepath.Join(root, extension))
	require.NoError(t, err)
	write(extension, source, 0644)
	target := "pkg/tidbcloud/v1beta2/nextgen/private_endpoint_connection.go"
	current, err := os.ReadFile(filepath.Join(root, target))
	require.NoError(t, err)
	require.Equal(t, string(source), string(current), "extension source and SDK copy must stay synchronized")
	write(target, []byte("stale file must be replaced"), 0644)
	write("pkg/tidbcloud/v1beta2/nextgen/stale.go", []byte("stale"), 0644)
	write("bin/npx", []byte("#!/bin/sh\nmkdir -p ../../pkg/tidbcloud/v1beta2/nextgen\n"), 0755)
	for _, tool := range []string{"go", "goimports"} {
		write("bin/"+tool, []byte("#!/bin/sh\nexit 0\n"), 0755)
	}
	command := exec.Command(makePath, "-o", "install-openapi-generator", "generate-v1beta2-nextgen-client")
	command.Dir = workspace
	command.Env = append(os.Environ(), "PATH="+filepath.Join(workspace, "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
	output, err := command.CombinedOutput()
	require.NoError(t, err, "%s", output)
	restored, err := os.ReadFile(filepath.Join(workspace, target))
	require.NoError(t, err)
	require.Equal(t, string(source), string(restored))
	_, err = os.Stat(filepath.Join(workspace, "pkg/tidbcloud/v1beta2/nextgen/stale.go"))
	require.True(t, os.IsNotExist(err))
}
