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
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/icholy/digest"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func redirectTestClient(t *testing.T, auth string, inner http.RoundTripper) *http.Client {
	t.Helper()
	var rt http.RoundTripper
	if auth == "bearer" {
		rt = NewTransportWithBearToken(inner, "test-token")
	} else {
		dt := NewDigestTransport("test-public-key", "test-private-key").(*digest.Transport)
		dt.Transport = inner
		rt = dt
	}
	branch, _, _, _, _, iam, _, _, _, _, err := NewApiClient(rt, DefaultServerlessEndpoint, DefaultIAMEndpoint)
	require.NoError(t, err)
	require.Same(t, branch.GetConfig().HTTPClient, iam.GetConfig().HTTPClient)
	return branch.GetConfig().HTTPClient
}

func testHTTPResponse(r *http.Request, status int, header http.Header) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     header,
		Body:       http.NoBody,
		Request:    r,
	}
}

func TestAPIRedirectRejectsOtherOrigins(t *testing.T) {
	for _, auth := range []string{"bearer", "digest"} {
		for _, origin := range []string{DefaultServerlessEndpoint, DefaultIAMEndpoint} {
			for _, target := range []string{
				"https://attacker.example/target",
				origin + ".attacker.example/target",
				origin + ":8443/target",
				strings.Replace(origin, "https:", "http:", 1) + "/target",
			} {
				for _, status := range []int{301, 302, 303, 307, 308} {
					t.Run(auth+"/"+origin+"/"+target+"/"+strconv.Itoa(status), func(t *testing.T) {
						var targetRequests, targetCredentials int
						client := redirectTestClient(t, auth, roundTripFunc(func(r *http.Request) (*http.Response, error) {
							if r.URL.String() == origin+"/start" {
								return testHTTPResponse(r, status, http.Header{"Location": {"/middle"}}), nil
							}
							if r.URL.String() == origin+"/middle" {
								return testHTTPResponse(r, status, http.Header{"Location": {target}}), nil
							}
							targetRequests++
							if r.Header.Get("Authorization") != "" {
								targetCredentials++
								return testHTTPResponse(r, http.StatusOK, nil), nil
							}
							// A malicious target can issue its own Digest challenge.
							return testHTTPResponse(r, http.StatusUnauthorized, http.Header{
								"Www-Authenticate": {`Digest realm="attacker", nonce="test-nonce", qop="auth", algorithm=MD5`},
							}), nil
						}))
						resp, err := client.Get(origin + "/start")
						if resp != nil {
							resp.Body.Close()
						}
						require.Zero(t, targetCredentials, "redirect target received credentials")
						require.Zero(t, targetRequests, "redirect target received a request")
						require.ErrorContains(t, err, "API redirect must remain on the same HTTPS origin")
					})
				}
			}
		}
	}
}

func TestAPIRedirectAllowsSameOrigin(t *testing.T) {
	for _, auth := range []string{"bearer", "digest"} {
		for _, target := range []string{"/target", "https://SERVERLESS.tidbapi.com:443/target"} {
			t.Run(auth+"/"+target, func(t *testing.T) {
				var authorizedPaths []string
				client := redirectTestClient(t, auth, roundTripFunc(func(r *http.Request) (*http.Response, error) {
					if auth == "digest" && r.Header.Get("Authorization") == "" {
						return testHTTPResponse(r, http.StatusUnauthorized, http.Header{
							"Www-Authenticate": {`Digest realm="api", nonce="test-nonce", qop="auth", algorithm=MD5`},
						}), nil
					}
					require.NotEmpty(t, r.Header.Get("Authorization"))
					authorizedPaths = append(authorizedPaths, r.URL.Path)
					if r.URL.Path == "/start" {
						return testHTTPResponse(r, http.StatusTemporaryRedirect, http.Header{"Location": {target}}), nil
					}
					return testHTTPResponse(r, http.StatusOK, nil), nil
				}))
				resp, err := client.Get(DefaultServerlessEndpoint + "/start")
				require.NoError(t, err)
				defer resp.Body.Close()
				require.Equal(t, http.StatusOK, resp.StatusCode)
				require.Equal(t, []string{"/start", "/target"}, authorizedPaths)
			})
		}
	}
}

func TestAPIRedirectLimit(t *testing.T) {
	requests := 0
	client := redirectTestClient(t, "bearer", roundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests++
		return testHTTPResponse(r, http.StatusFound, http.Header{"Location": {"/loop"}}), nil
	}))
	resp, err := client.Get(DefaultServerlessEndpoint + "/loop")
	if resp != nil {
		resp.Body.Close()
	}
	require.ErrorContains(t, err, "stopped after 10 redirects")
	require.Equal(t, 10, requests)
}

func TestAPIErrorBodyRedacted(t *testing.T) {
	for _, body := range []string{
		`{"password":"body-secret","message":"invalid request","nested":{"access_token":"token-secret"},"url":"https://example.com?X-Amz-Signature=url-secret"}`,
		"unstructured body-secret",
	} {
		resp := &http.Response{
			Status: "400 Bad Request",
			Header: http.Header{"X-Debug-Trace-Id": []string{"trace-123"}},
			Body:   io.NopCloser(strings.NewReader(body)),
		}
		err := parseError(errors.New("SDK may duplicate body-secret"), resp)
		for _, secret := range []string{"body-secret", "token-secret", "url-secret"} {
			require.NotContains(t, err.Error(), secret)
		}
		require.Contains(t, err.Error(), "400 Bad Request")
		require.Contains(t, err.Error(), "trace-123")
	}
}
