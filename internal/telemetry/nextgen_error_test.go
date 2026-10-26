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

package telemetry

import (
	"context"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type nextgenErrorTransport func(*http.Request) (*http.Response, error)

func (f nextgenErrorTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestNextGenAPIErrorTelemetryRedaction(t *testing.T) {
	old := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = old })
	for _, status := range []int{400, 200} {
		for _, body := range []string{
			`{"message":"invalid request","rootPassword":"secret-password","details":[{"token":"secret-token"}]}`,
			"proxy secret-password secret-token",
		} {
			http.DefaultTransport = nextgenErrorTransport(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: status, Status: http.StatusText(status), Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})
			client, err := cloud.NewNextGenClientDelegateWithToken("test-token", "https://example.test", "")
			require.NoError(t, err)
			_, err = client.GetTiDB(context.Background(), "123")
			require.Error(t, err)
			event := newEvent(withError(err))
			require.Equal(t, "ERROR", event.Properties["result"])
			require.NotContains(t, event.Properties["error"], "secret-password")
			require.NotContains(t, event.Properties["error"], "secret-token")
			if status == 400 {
				require.Contains(t, event.Properties["error"], "400 Bad Request")
			} else {
				require.Contains(t, event.Properties["error"], "failed to read or decode the API response")
			}
		}
	}
}

// Safe wrappers must be recorded using their display text, never their raw cause.
// CA download errors use this contract to retain timeout identity without URLs.
type safeTelemetryError struct{ cause error }

func (e safeTelemetryError) Error() string { return "download CA certificate: request timed out" }
func (e safeTelemetryError) Unwrap() error { return e.cause }
func TestTelemetryDoesNotUnwrapSafeDisplayError(t *testing.T) {
	err := safeTelemetryError{&neturl.Error{Op: "Get", URL: "https://example.test/ca?X-Amz-Signature=secret-signature", Err: context.DeadlineExceeded}}
	event := newEvent(withError(err))
	require.Equal(t, err.Error(), event.Properties["error"])
	require.NotContains(t, event.Properties["error"], "secret-signature")
}
