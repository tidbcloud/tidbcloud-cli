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

package redact

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnosticText(t *testing.T) {
	uploadErr := fmt.Errorf("upload failed: %w", &url.Error{
		Op:  "Put",
		URL: "https://storage.example/upload?partNumber=1&X-Amz-Signature=signature-secret&X-Amz-Credential=credential-secret&X-Amz-Security-Token=session-secret",
		Err: errors.New("connection reset"),
	})
	message := Text(uploadErr.Error())
	for _, secret := range []string{"signature-secret", "credential-secret", "session-secret"} {
		require.NotContains(t, message, secret)
	}
	require.Contains(t, message, "storage.example/upload")
	require.Contains(t, message, "partNumber=1")
	require.Contains(t, message, "connection reset")
	require.NotContains(t, Text("failed tidb://user:sql-secret@example.com:4000"), "sql-secret")
	require.NotContains(t, Text("failed https://example.com/%xx?token=bad-secret"), "bad-secret")
	require.Equal(t, "connection refused", Text("connection refused"))
}
