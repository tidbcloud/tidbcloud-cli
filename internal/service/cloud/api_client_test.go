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
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal/redact"
)

type trackingReadCloser struct {
	io.Reader
	closed bool
}

func (r *trackingReadCloser) Close() error {
	r.closed = true
	return nil
}

func TestDebugTransportUsesSharedRedaction(t *testing.T) {
	t.Setenv("TICLOUD_DEBUG", "true")
	transport, ok := NewDebugTransport(http.DefaultTransport).(*redact.DebugTransport)
	require.True(t, ok)
	require.True(t, transport.Debug)
	const body = `{"rootPassword":"super-secret"}`
	request, err := http.NewRequest(http.MethodPost, "https://example.com/v1beta2/tidbs/123:resetRootPassword", strings.NewReader(body))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer secret")
	dump := redact.DumpRequest(request)
	require.NotContains(t, dump, "super-secret")
	require.NotContains(t, dump, "Bearer secret")
	require.Equal(t, "Bearer secret", request.Header.Get("Authorization"))
	payload, err := io.ReadAll(request.Body)
	require.NoError(t, err)
	require.JSONEq(t, body, string(payload))
}

func TestParseNextGenErrorOmitsHTMLResponseBody(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "https://example.com/v1beta2/regions", nil)
	require.NoError(t, err)
	body := "<!doctype html><html><body>sensitive-noise</body></html>"
	responseBody := &trackingReadCloser{Reader: strings.NewReader(body)}
	response := &http.Response{
		Header:  http.Header{"Content-Type": []string{"text/html; charset=utf-8"}},
		Body:    responseBody,
		Request: request,
	}

	parsed := parseNextGenError(fmt.Errorf("undefined response type"), response)
	require.Contains(t, parsed.Error(), "[GET /v1beta2/regions]")
	require.Contains(t, parsed.Error(), "body omitted")
	require.NotContains(t, parsed.Error(), "sensitive-noise")
	require.True(t, responseBody.closed)
}

func TestParseErrorPreservesHTMLResponseBody(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "https://example.com/v1beta1/clusters", nil)
	require.NoError(t, err)
	body := "<!doctype html><html><body>legacy-error</body></html>"
	responseBody := &trackingReadCloser{Reader: strings.NewReader(body)}
	response := &http.Response{
		Header:  http.Header{"Content-Type": []string{"text/html; charset=utf-8"}},
		Body:    responseBody,
		Request: request,
	}

	parsed := parseError(fmt.Errorf("undefined response type"), response)
	require.EqualError(t, parsed, `[GET /v1beta1/clusters][undefined response type][<trace_id>] <!doctype html><html><body>legacy-error</body></html>`)
	require.True(t, responseBody.closed)
}
