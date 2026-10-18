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
)

type trackingReadCloser struct {
	io.Reader
	closed bool
}

func (r *trackingReadCloser) Close() error {
	r.closed = true
	return nil
}

func TestRequestForDebugDumpRedactsAuthorizationWithoutMutatingRequest(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "https://example.com", nil)
	require.NoError(t, err)
	originalBody := &trackingReadCloser{Reader: strings.NewReader("payload")}
	request.Body = originalBody
	request.ContentLength = int64(len("payload"))
	request.Header.Set("Authorization", "Bearer secret")

	redacted, err := requestForDebugDump(request)
	require.NoError(t, err)
	require.True(t, originalBody.closed)
	require.Equal(t, "<redacted>", redacted.Header.Get("Authorization"))
	require.Equal(t, "Bearer secret", request.Header.Get("Authorization"))
	redactedBody, err := io.ReadAll(redacted.Body)
	require.NoError(t, err)
	require.Equal(t, "payload", string(redactedBody))
	originalPayload, err := io.ReadAll(request.Body)
	require.NoError(t, err)
	require.Equal(t, "payload", string(originalPayload))
}

func TestRequestForDebugDumpRedactsRootPasswordWithoutMutatingRequest(t *testing.T) {
	const body = `{"rootPassword":"super-secret"}`
	request, err := http.NewRequest(http.MethodPost, "https://example.com/v1beta2/tidbs/123:resetRootPassword", strings.NewReader(body))
	require.NoError(t, err)

	redacted, err := requestForDebugDump(request)
	require.NoError(t, err)
	redactedBody, err := io.ReadAll(redacted.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"rootPassword":"<redacted>"}`, string(redactedBody))
	require.NotContains(t, string(redactedBody), "super-secret")
	originalBody, err := io.ReadAll(request.Body)
	require.NoError(t, err)
	require.JSONEq(t, body, string(originalBody))
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
	require.EqualError(t, parsed, `[GET /v1beta2/regions][undefined response type][<trace_id>] unexpected response content type "text/html; charset=utf-8" (56-byte body omitted)`)
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
