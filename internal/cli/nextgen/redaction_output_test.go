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

package nextgen

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type redactionTransport func(*http.Request) (*http.Response, error)

func (f redactionTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestNextGenErrorsSafeForTerminalAndJSON(t *testing.T) {
	oldTransport := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = oldTransport })
	http.DefaultTransport = redactionTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "ca.example.test" {
			if r.Context().Err() != nil {
				return nil, r.Context().Err()
			}
			return nil, context.DeadlineExceeded
		}
		return &http.Response{StatusCode: 400, Status: "400 Bad Request", Header: http.Header{"Content-Type": {"application/json"}},
			Body: io.NopCloser(strings.NewReader(`{"message":"invalid request","rootPassword":"secret-password","details":[{"token":"secret-token"}]}`)), Request: r}, nil
	})
	client, err := cloud.NewNextGenClientDelegateWithToken("test-token", "https://api.example.test", "")
	require.NoError(t, err)
	_, apiErr := client.GetTiDB(context.Background(), "123")
	_, caErr := downloadCACertificate(context.Background(), "https://ca.example.test/ca?X-Amz-Signature=secret-signature")
	require.ErrorIs(t, caErr, context.DeadlineExceeded)
	require.Contains(t, caErr.Error(), "timed out")
	for _, sourceErr := range []error{apiErr, caErr} {
		require.Error(t, sourceErr)
		terminal := "Error: " + sourceErr.Error()
		encoded, ok := marshalJSONError(sourceErr)
		require.True(t, ok)
		require.True(t, json.Valid([]byte(encoded)))
		for _, text := range []string{terminal, encoded} {
			for _, secret := range []string{"secret-password", "secret-token", "secret-signature"} {
				require.NotContains(t, text, secret)
			}
		}
	}
	_, malformed := downloadCACertificate(context.Background(), "https://ca.example.test/%zz?token=secret-token")
	require.Error(t, malformed)
	require.NotContains(t, malformed.Error(), "secret-token")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, canceled := downloadCACertificate(ctx, "https://ca.example.test/ca?token=secret-token")
	require.Error(t, canceled)
	// The transport's cause remains accessible despite safe display text.
	require.ErrorIs(t, canceled, context.Canceled)
	require.Contains(t, canceled.Error(), "request canceled")
}

// The SDK consumes the initial read error; its error parser sees EOF afterward.
type failedResponseBody struct {
	cause error
	read  bool
}

func (b *failedResponseBody) Read([]byte) (int, error) {
	if !b.read {
		b.read = true
		return 0, b.cause
	}
	return 0, io.EOF
}
func (b *failedResponseBody) Close() error { return nil }

func TestJSONOutputPreservesLocalAPIFailure(t *testing.T) {
	oldTransport := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = oldTransport })
	cases := []struct {
		name   string
		status int
		cause  error
		reason string
	}{
		{"read timeout after 200", 200, context.DeadlineExceeded, "context deadline exceeded"},
		{"read timeout after 400", 400, context.DeadlineExceeded, "context deadline exceeded"},
		{"read canceled", 200, context.Canceled, "context canceled"},
	}
	for _, status := range []int{301, 302, 303, 307, 308} {
		cases = append(cases, struct {
			name   string
			status int
			cause  error
			reason string
		}{
			http.StatusText(status), status, nil, "Cloud API redirects are not allowed; check the cloud-api-endpoint configuration",
		})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			http.DefaultTransport = redactionTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				header := http.Header{"Content-Type": {"application/json"}, "X-Debug-Trace-Id": {"trace-local-failure"}}
				var body io.ReadCloser = http.NoBody
				if tc.cause != nil {
					body = &failedResponseBody{cause: tc.cause}
				} else {
					header.Set("Location", "https://other.example/collect?X-Amz-Signature=secret-signature")
				}
				return &http.Response{StatusCode: tc.status, Status: http.StatusText(tc.status), Header: header, Body: body, Request: r}, nil
			})
			client, err := cloud.NewNextGenClientDelegateWithToken("secret-token", "https://api.example.test", "")
			require.NoError(t, err)
			_, err = client.GetTiDB(context.Background(), "123")
			require.Error(t, err)
			require.Equal(t, 1, calls, "no retry or redirect request")
			require.Contains(t, err.Error(), tc.reason)
			if tc.cause != nil {
				require.ErrorIs(t, err, tc.cause)
			}
			for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
				root := cmd(helperWithNextGenClient(nil), plan)
				command, _, findErr := root.Find([]string{"describe"})
				require.NoError(t, findErr)
				text, ok := JSONError(command, []string{"describe", "-c", "123", "-o", "json"}, err)
				require.True(t, ok)
				var result struct {
					Error struct {
						Message   string `json:"message"`
						Status    int    `json:"status"`
						RequestID string `json:"requestId"`
						Method    string `json:"method"`
						Path      string `json:"path"`
					} `json:"error"`
				}
				require.NoError(t, json.Unmarshal([]byte(text), &result))
				require.Equal(t, tc.reason, result.Error.Message)
				require.Equal(t, tc.status, result.Error.Status)
				require.Equal(t, "trace-local-failure", result.Error.RequestID)
				require.Equal(t, "GET", result.Error.Method)
				require.Equal(t, "/v1beta2/tidbs/123", result.Error.Path)
				require.NotContains(t, text, "secret-signature")
				require.NotContains(t, text, "secret-token")
			}
		})
	}
}

func TestMalformedSuccessResponsesHaveSafeFailureMessage(t *testing.T) {
	oldTransport := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = oldTransport })
	for _, tc := range []struct{ name, contentType, body string }{
		{"invalid JSON", "application/json", `{"rootPassword":"secret-password"`},
		{"wrong field type", "application/json", `{"displayName":123,"regionId":"aws-us-west-2","servicePlan":"Premium","rootPassword":"secret-password"}`},
		{"missing required fields", "application/json", `{"rootPassword":"secret-password"}`},
		{"wrong content type", "text/html", "<html>secret-proxy-body</html>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			http.DefaultTransport = redactionTransport(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Status: "200 OK",
					Header: http.Header{"Content-Type": {tc.contentType}, "X-Debug-Trace-Id": {"trace-decode"}},
					Body:   io.NopCloser(strings.NewReader(tc.body)), Request: r}, nil
			})
			client, err := cloud.NewNextGenClientDelegateWithToken("test-token", "https://api.example.test", "")
			require.NoError(t, err)
			_, err = client.GetTiDB(context.Background(), "123")
			require.ErrorContains(t, err, "failed to read or decode the API response")
			for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
				root := cmd(helperWithNextGenClient(nil), plan)
				command, _, findErr := root.Find([]string{"describe"})
				require.NoError(t, findErr)
				text, ok := JSONError(command, []string{"describe", "-c", "123", "-o", "json"}, err)
				require.True(t, ok)
				var result map[string]map[string]interface{}
				require.NoError(t, json.Unmarshal([]byte(text), &result))
				require.Equal(t, "failed to read or decode the API response", result["error"]["message"])
				require.Equal(t, float64(200), result["error"]["status"])
				require.Equal(t, "trace-decode", result["error"]["requestId"])
				require.Equal(t, "/v1beta2/tidbs/123", result["error"]["path"])
				require.NotContains(t, result["error"], "response")
				for _, out := range []string{err.Error(), text} {
					require.NotContains(t, out, "secret-")
					require.NotContains(t, out, tc.body)
				}
			}
		})
	}
}
