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

package export

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
)

func TestDownloadResultGetResultSanitized(t *testing.T) {
	// A server-provided file name must not inject terminal escape sequences,
	// and a download error must not echo the pre-signed URL signature.
	name := "dump\x1b]8;;http://evil/\x07.tar.gz"
	presigned := fmt.Errorf(`Get "https://bucket.example/dump.tar.gz?X-Amz-Signature=abc123&X-Amz-Credential=cred": dial tcp: lookup no-such-host: no such host`)

	failed := &downloadResult{name: name, err: presigned, status: Failed}
	text := failed.GetResult()
	require.NotContains(t, text, "\x1b")
	require.NotContains(t, text, "abc123")
	require.NotContains(t, text, "cred")
	require.Contains(t, text, "failed: Get")

	skipped := &downloadResult{name: name, err: errors.New("file already exists"), status: Skipped}
	require.NotContains(t, skipped.GetResult(), "\x1b")
	require.Contains(t, skipped.GetResult(), "skipped: file already exists")

	succeeded := &downloadResult{name: name, status: Succeeded}
	require.Equal(t, "dump.tar.gz success", succeeded.GetResult())
}

func TestDownloadDiagnosticsAcrossModes(t *testing.T) {
	const filename = "dump\r.csv"
	downloadErr := errors.New(`could not fetch "https://storage.example/export?X-Amz-Signature=TEST_SIGNATURE"`)

	t.Run("interactive progress and summary", func(t *testing.T) {
		job := &FileJob{name: filename, err: downloadErr, status: Failed}
		model := NewProcessDownloadModel(1, t.TempDir(), 1, []string{filename}, nil)
		bar := progress.New()
		model.progressBar.progress = &bar
		model.jobInfo.finishedJobs = []*FileJob{job}
		for _, output := range []string{model.View(), job.GetResult()} {
			require.Contains(t, output, "dump.csv failed")
			require.NotContains(t, output, filename)
			require.NotContains(t, output, "TEST_SIGNATURE")
		}
	})

	t.Run("non-interactive progress and summary", func(t *testing.T) {
		streams := iostream.Test()
		pool, err := NewDownloadPool(&internal.Helper{IOStreams: streams}, t.TempDir(), 1, []string{filename},
			func(context.Context, []string) (map[string]*string, error) { return nil, downloadErr })
		require.NoError(t, err)
		require.ErrorContains(t, pool.Start(), "1 file(s) failed to download")
		output := streams.Out.(*bytes.Buffer).String()
		require.Contains(t, output, "download dump.csv failed")
		require.Contains(t, output, "1.dump.csv failed")
		require.NotContains(t, output, filename)
		require.NotContains(t, output, "TEST_SIGNATURE")
	})
}
