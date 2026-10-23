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

package cloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextGenStructuredErrorPreservesMetadataAndRedactsBody(t *testing.T) {
	request, err := http.NewRequest("POST", "https://example.com/v1beta2/tidbs?token=do-not-print", nil)
	require.NoError(t, err)
	body := `{"code":400,"message":"Min RCU must be less than or equal to max RCU.","details":[{"requestId":"trace-123","rootPassword":"secret-password"}],"token":"secret-token","large":1372813089454840200}`
	response := &http.Response{StatusCode: 400, Request: request, Header: http.Header{"X-Debug-Trace-Id": {"trace-123"}}, Body: io.NopCloser(strings.NewReader(body))}
	parsed := parseNextGenError(fmt.Errorf("400 Bad Request"), response)
	var structured *NextGenAPIError
	require.True(t, errors.As(parsed, &structured))
	require.Equal(t, "trace-123", structured.RequestID)
	require.Equal(t, 400, structured.Status)
	require.Equal(t, "/v1beta2/tidbs", structured.Path)
	require.Equal(t, "POST", structured.Method)
	encoded, err := json.Marshal(structured)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "secret-password")
	require.NotContains(t, string(encoded), "secret-token")
	require.NotContains(t, string(encoded), "do-not-print")
	require.Contains(t, string(encoded), "1372813089454840200")
	// Terminal and telemetry both consume Error(), not the structured JSON.
	require.NotContains(t, parsed.Error(), "secret-password")
	require.NotContains(t, parsed.Error(), "secret-token")
	require.Contains(t, parsed.Error(), "[POST /v1beta2/tidbs][400 Bad Request][trace-123]")
	require.Contains(t, parsed.Error(), "Min RCU")
	require.Contains(t, parsed.Error(), "1372813089454840200")
}

func TestNextGenStructuredErrorNonJSONDoesNotEchoBody(t *testing.T) {
	response := &http.Response{StatusCode: 502, Header: http.Header{}, Body: io.NopCloser(strings.NewReader("proxy secret"))}
	var structured *NextGenAPIError
	require.True(t, errors.As(parseNextGenError(errors.New("bad gateway"), response), &structured))
	require.Nil(t, structured.Response)
	require.NotContains(t, structured.Error(), "proxy secret")
	data, err := json.Marshal(structured)
	require.NoError(t, err)
	require.NotContains(t, string(data), "proxy secret")
}
