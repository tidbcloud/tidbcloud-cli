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
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/stretchr/testify/require"
)

func TestNextGenClientUsesV1Beta2EndpointPlanAndBearerToken(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "/v1beta2/tidbs", r.URL.Path)
		require.Equal(t, "Premium", r.URL.Query().Get("servicePlan"))
		require.Equal(t, "10", r.URL.Query().Get("pageSize"))
		require.Equal(t, "next", r.URL.Query().Get("pageToken"))
		require.Equal(t, "Bearer oauth-token", r.Header.Get("Authorization"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"tidbs":[]}`)),
			Request:    r,
		}, nil
	})

	client, err := newNextGenClientDelegate(NewTransportWithBearToken(transport, "oauth-token"), "https://example.test")
	require.NoError(t, err)
	size := int32(10)
	token := "next"
	response, err := client.ListTiDBs(context.Background(), api.TIDBSERVICELISTTIDBSSERVICEPLANPARAMETER_PREMIUM, &size, &token)
	require.NoError(t, err)
	require.Empty(t, response.Tidbs)
}

func TestNextGenClientResetsRootPassword(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1beta2/tidbs/tidb-1:resetRootPassword", r.URL.Path)
		var body api.TidbServiceResetRootPasswordBody
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "new-password", body.GetRootPassword())
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{}`)),
			Request:    r,
		}, nil
	})

	client, err := newNextGenClientDelegate(transport, "https://example.test")
	require.NoError(t, err)
	err = client.ResetRootPassword(context.Background(), "tidb-1", api.NewTidbServiceResetRootPasswordBody("new-password"))
	require.NoError(t, err)
}

func TestNextGenClientUpdatesPublicConnectionSetting(t *testing.T) {
	for _, auth := range []string{"bearer", "digest"} {
		for _, enabled := range []bool{true, false} {
			action := "disable"
			if enabled {
				action = "enable"
			}
			t.Run(auth+"/"+action, func(t *testing.T) {
				calls := 0
				transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
					calls++
					require.Equal(t, http.MethodPatch, r.Method)
					require.Equal(t, "/v1beta2/tidbs/tidb-1/publicConnectionSetting", r.URL.Path)
					require.Empty(t, r.URL.RawQuery)
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					expected := `{"enabled":true}`
					if !enabled {
						expected = `{"enabled":false}`
					}
					// Both the initial request and Digest retry must preserve false
					// and omit IP access list fields.
					require.JSONEq(t, expected, string(body))
					response := &http.Response{
						StatusCode: http.StatusOK,
						Status:     "200 OK",
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(expected)),
						Request:    r,
					}
					if auth == "digest" && r.Header.Get("Authorization") == "" {
						response.StatusCode, response.Status = http.StatusUnauthorized, "401 Unauthorized"
						response.Header.Set("WWW-Authenticate", `Digest realm="test", nonce="nonce", algorithm=MD5, qop="auth"`)
						return response, nil
					}
					if auth == "bearer" {
						require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
					} else {
						require.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Digest "))
						require.NotContains(t, r.Header.Get("Authorization"), "test-private-key")
					}
					return response, nil
				})

				var authenticated http.RoundTripper = NewTransportWithBearToken(transport, "test-token")
				if auth == "digest" {
					authenticated = NewTransportWithDigest(transport, "test-public-key", "test-private-key")
				}
				client, err := newNextGenClientDelegate(authenticated, "https://example.test")
				require.NoError(t, err)
				request := api.NewV1beta2PublicConnectionSetting()
				request.SetEnabled(enabled)
				result, err := client.UpdatePublicConnectionSetting(context.Background(), "tidb-1", request)
				require.NoError(t, err)
				require.True(t, result.HasEnabled())
				require.Equal(t, enabled, result.GetEnabled())
				if auth == "digest" {
					require.Equal(t, 2, calls)
				} else {
					require.Equal(t, 1, calls)
				}
			})
		}
	}
}

func TestNextGenCMEKRequestContract(t *testing.T) {
	for _, auth := range []string{"bearer", "digest"} {
		for _, provider := range []string{"aws", "alicloud"} {
			t.Run(auth+"/"+provider, func(t *testing.T) {
				region := provider + "-test-region"
				key := api.NewV1beta2CustomerManagedEncryptionKey(region, api.V1BETA1SERVICEPLAN_PREMIUM)
				keyField, principalField := "awsKms", "awsPrincipal"
				if provider == "aws" {
					key.AwsKms = api.NewCustomerManagedEncryptionKeyAwsKms("test-key-arn")
					key.AwsPrincipal = &api.CmekAccessIamPrincipalAwsCmekPrincipal{AccountId: api.PtrString("123")}
				} else {
					keyField, principalField = "aliyunKms", "aliyunPrincipal"
					key.AliyunKms = api.NewCustomerManagedEncryptionKeyAliyunKms("test-key-arn")
					key.AliyunPrincipal = &api.CmekAccessIamPrincipalAliyunCmekPrincipal{AccountId: api.PtrString("123")}
				}
				var calls []string
				transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
					response := &http.Response{StatusCode: http.StatusOK, Status: "200 OK",
						Header: http.Header{"Content-Type": []string{"application/json"}}, Request: r}
					if auth == "digest" && r.Header.Get("Authorization") == "" {
						response.StatusCode, response.Status = http.StatusUnauthorized, "401 Unauthorized"
						response.Header.Set("WWW-Authenticate", `Digest realm="test", nonce="nonce", algorithm=MD5, qop="auth"`)
						response.Body = io.NopCloser(strings.NewReader(""))
						return response, nil
					}
					if auth == "bearer" {
						require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
					} else {
						require.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Digest "))
						require.NotContains(t, r.Header.Get("Authorization"), "test-private-key")
					}
					calls = append(calls, r.Method+" "+r.URL.Path)
					switch r.URL.Path {
					case "/v1beta2/cmeks:principal":
						require.Equal(t, http.MethodGet, r.Method)
						require.Equal(t, region, r.URL.Query().Get("regionId"))
						require.Equal(t, "Premium", r.URL.Query().Get("servicePlan"))
						response.Body = io.NopCloser(strings.NewReader(`{"servicePlan":"Premium"}`))
					case "/v1beta2/cmeks:verifyPrincipal":
						require.Equal(t, http.MethodPost, r.Method)
						require.Empty(t, r.URL.RawQuery) // No existing tidbId during creation.
						body, err := io.ReadAll(r.Body)
						require.NoError(t, err)
						// The HTTP body is the key itself, not a {"key": ...} wrapper.
						require.JSONEq(t, `{"regionId":"`+region+`","servicePlan":"Premium","`+keyField+`":{"kmsKeyArn":"test-key-arn"},"`+principalField+`":{"accountId":"123"}}`, string(body))
						response.Body = io.NopCloser(strings.NewReader(`{"valid":true}`))
					default:
						t.Fatalf("unexpected path: %s", r.URL.Path)
					}
					return response, nil
				})
				var authenticated http.RoundTripper = NewTransportWithBearToken(transport, "test-token")
				if auth == "digest" {
					authenticated = NewTransportWithDigest(transport, "test-public-key", "test-private-key")
				}
				client, err := newNextGenClientDelegate(authenticated, "https://example.test")
				require.NoError(t, err)
				_, err = client.GetCmekAccessIAMPrincipal(context.Background(), region)
				require.NoError(t, err)
				result, err := client.VerifyCmekAccessIAMPrincipal(context.Background(), key)
				require.NoError(t, err)
				require.True(t, result.GetValid())
				require.Equal(t, []string{"GET /v1beta2/cmeks:principal", "POST /v1beta2/cmeks:verifyPrincipal"}, calls)
			})
		}
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestNextGenClientDoesNotDuplicateVersionPath(t *testing.T) {
	client, err := NewNextGenClientDelegateWithToken("token", "https://example.com/v1beta2/", "")
	require.NoError(t, err)
	require.Equal(t, "https://example.com/v1beta2", client.client.GetConfig().Servers[0].URL)
}

func TestNextGenClientRequiresHTTPS(t *testing.T) {
	_, err := NewNextGenClientDelegateWithToken("token", "http://example.com", "")
	require.ErrorContains(t, err, "must use HTTPS")
}

func TestNextGenClientRejectsFragmentEndpoint(t *testing.T) {
	_, err := NewNextGenClientDelegateWithToken("token", "https://example.com/path#fragment", "")
	require.ErrorContains(t, err, "must not contain a query or fragment")
}

func TestNextGenClientRejectsQueryEndpoint(t *testing.T) {
	_, err := NewNextGenClientDelegateWithToken("token", "https://example.com/path?region=us-west-2", "")
	require.ErrorContains(t, err, "must not contain a query or fragment")
}

func TestNextGenClientRejectsEmptyQueryEndpoint(t *testing.T) {
	_, err := NewNextGenClientDelegateWithToken("token", "https://example.com/path?", "")
	require.ErrorContains(t, err, "must not contain a query or fragment")
}

func TestNextGenClientUsesCustomCACertificate(t *testing.T) {
	requestDetails := make(chan [2]string, 1)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestDetails <- [2]string{r.URL.Path, r.Header.Get("Authorization")}
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, `{"tidbs":[]}`); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	caCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	caCertPath := filepath.Join(t.TempDir(), "ca.pem")
	require.NoError(t, os.WriteFile(caCertPath, caCert, 0o600))

	client, err := NewNextGenClientDelegateWithToken("token", server.URL, caCertPath)
	require.NoError(t, err)
	response, err := client.ListTiDBs(context.Background(), api.TIDBSERVICELISTTIDBSSERVICEPLANPARAMETER_PREMIUM, nil, nil)
	require.NoError(t, err)
	require.Empty(t, response.Tidbs)
	require.Equal(t, [2]string{"/v1beta2/tidbs", "Bearer token"}, <-requestDetails)
}

func TestNextGenClientCustomCARequiresTLS12(t *testing.T) {
	server := httptest.NewTLSServer(nil)
	defer server.Close()

	caCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	caCertPath := filepath.Join(t.TempDir(), "ca.pem")
	require.NoError(t, os.WriteFile(caCertPath, caCert, 0o600))

	roundTripper, err := nextGenTransport(caCertPath)
	require.NoError(t, err)
	transport, ok := roundTripper.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport.TLSClientConfig)
	require.Equal(t, uint16(tls.VersionTLS12), transport.TLSClientConfig.MinVersion)
}

func TestNextGenClientCustomCAPreservesStricterTLSMinimum(t *testing.T) {
	server := httptest.NewTLSServer(nil)
	defer server.Close()
	caCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	caCertPath := filepath.Join(t.TempDir(), "ca.pem")
	require.NoError(t, os.WriteFile(caCertPath, caCert, 0o600))

	originalTransport := http.DefaultTransport
	strictTransport := originalTransport.(*http.Transport).Clone()
	strictTransport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS13}
	http.DefaultTransport = strictTransport
	t.Cleanup(func() { http.DefaultTransport = originalTransport })

	roundTripper, err := nextGenTransport(caCertPath)
	require.NoError(t, err)
	transport, ok := roundTripper.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport.TLSClientConfig)
	require.Equal(t, uint16(tls.VersionTLS13), transport.TLSClientConfig.MinVersion)
}
