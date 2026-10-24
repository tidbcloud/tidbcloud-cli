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
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"
)

func TestPrivateEndpointRequests(t *testing.T) {
	for _, auth := range []string{"oauth", "api-key"} {
		for _, method := range []string{"POST", "GET", "DELETE"} {
			t.Run(auth+method, func(t *testing.T) {
				calls := 0
				transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
					calls++
					require.Equal(t, method, r.Method)
					path := "/v1beta2/tidbs/123/privateEndpointConnections"
					if method != "POST" {
						path += "/connection-1"
					}
					require.Equal(t, path, r.URL.Path)
					if method == "POST" {
						b, e := io.ReadAll(r.Body)
						require.NoError(t, e)
						require.JSONEq(t, `{"regionId":"aws-us-west-2","endpointId":"vpce-example"}`, string(b))
					}
					resp := &http.Response{StatusCode: 200, Status: "200 OK", Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"tidbId":"123","regionId":"aws-us-west-2","endpointId":"vpce-example"}`)), Request: r}
					if auth == "api-key" && r.Header.Get("Authorization") == "" {
						resp.StatusCode = 401
						resp.Status = "401 Unauthorized"
						resp.Header.Set("WWW-Authenticate", `Digest realm="test", nonce="nonce", algorithm=MD5, qop="auth"`)
						return resp, nil
					}
					if auth == "oauth" {
						require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
					} else {
						require.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Digest "))
					}
					return resp, nil
				})
				var rt http.RoundTripper = NewTransportWithBearToken(transport, "test-token")
				if auth == "api-key" {
					rt = NewTransportWithDigest(transport, "public", "private")
				}
				c, e := newNextGenClientDelegate(rt, "https://example.test")
				require.NoError(t, e)
				var result *api.Nextgenv1beta2PrivateEndpointConnection
				switch method {
				case "POST":
					result, e = c.CreatePrivateEndpointConnection(context.Background(), "123", api.NewNextgenv1beta2PrivateEndpointConnection("aws-us-west-2", "vpce-example"))
				case "GET":
					result, e = c.GetPrivateEndpointConnection(context.Background(), "123", "connection-1")
				case "DELETE":
					result, e = c.DeletePrivateEndpointConnection(context.Background(), "123", "connection-1")
				}
				require.NoError(t, e)
				require.Equal(t, "123", result.GetTidbId())
				if auth == "api-key" {
					require.Equal(t, 2, calls)
				} else {
					require.Equal(t, 1, calls)
				}
			})
		}
	}
}
func TestPrivateEndpointMutationTimeout(t *testing.T) {
	for _, action := range []string{"bind", "unbind"} {
		t.Run(action, func(t *testing.T) {
			calls := 0
			c, e := newNextGenClientDelegate(roundTripFunc(func(*http.Request) (*http.Response, error) { calls++; return nil, context.DeadlineExceeded }), "https://example.test")
			require.NoError(t, e)
			if action == "bind" {
				_, e = c.CreatePrivateEndpointConnection(context.Background(), "123", api.NewNextgenv1beta2PrivateEndpointConnection("r", "e"))
			} else {
				_, e = c.DeletePrivateEndpointConnection(context.Background(), "123", "c")
			}
			require.ErrorIs(t, e, context.DeadlineExceeded)
			require.ErrorContains(t, e, "Before retrying")
			require.Equal(t, 1, calls)
		})
	}
}

func TestPrivateEndpointErrorsAndRedirects(t *testing.T) {
	for _, status := range []int{400, 307} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			calls := 0
			c, e := newNextGenClientDelegate(roundTripFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				require.Equal(t, "example.test", r.URL.Host)
				return &http.Response{StatusCode: status, Status: http.StatusText(status), Header: http.Header{"Content-Type": {"application/json"}, "Location": {"https://other.test/"}}, Body: io.NopCloser(strings.NewReader(`{"message":"failed","password":"secret-value"}`)), Request: r}, nil
			}), "https://example.test")
			require.NoError(t, e)
			_, e = c.CreatePrivateEndpointConnection(context.Background(), "123", api.NewNextgenv1beta2PrivateEndpointConnection("r", "e"))
			require.Error(t, e)
			require.NotContains(t, e.Error(), "secret-value")
			require.Equal(t, 1, calls)
			if status == 307 {
				require.ErrorContains(t, e, "redirects are not allowed")
			}
		})
	}
}
