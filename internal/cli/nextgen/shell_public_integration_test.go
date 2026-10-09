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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/gohxs/readline"
	"github.com/stretchr/testify/require"
	"github.com/xo/usql/env"
	"github.com/xo/usql/text"
)

// Run against an isolated local MySQL with a certificate valid for 127.0.0.1
// but not localhost. No TiDB Cloud credentials or resources are used.
func TestShellPublicDNSRecoverySQL(t *testing.T) {
	address := os.Getenv("TICLOUD_TEST_MYSQL_ADDR")
	if address == "" {
		t.Skip("set TICLOUD_TEST_MYSQL_ADDR and TICLOUD_TEST_MYSQL_CA for local SQL integration")
	}
	host, port, err := net.SplitHostPort(address)
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", host, "this fixture must only connect to local MySQL")
	portNumber, err := strconv.ParseInt(port, 10, 32)
	require.NoError(t, err)
	ca, err := os.ReadFile(os.Getenv("TICLOUD_TEST_MYSQL_CA"))
	require.NoError(t, err)
	user := os.Getenv("TICLOUD_TEST_MYSQL_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("TICLOUD_TEST_MYSQL_PASSWORD")
	oldPrompt := env.Get("PROMPT1")
	t.Cleanup(func() { _ = env.Set("PROMPT1", oldPrompt) })
	t.Setenv(text.CommandUpper()+"_HISTORY", filepath.Join(t.TempDir(), "sql-history"))

	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, auth := range []string{"bearer", "digest"} {
			for _, mode := range []string{"valid", "wrong CA", "wrong hostname"} {
				t.Run(plan.commandName+"/"+auth+"/"+mode, func(t *testing.T) {
					endpoint := publicDNSFailureEndpoint(t)
					endpoint.SetHost(host)
					endpoint.SetPort(int32(portNumber))
					if auth == "digest" {
						endpoint.ConnectionReachability.Detail.SetMessage("DNS_PROBE_FAILED")
					}
					if mode == "wrong hostname" {
						endpoint.SetHost("localhost")
					}
					instance := api.NewNextgenv1beta2Tidb("test", "aws-us-west-2", plan.servicePlan)
					instance.State = api.V1BETA1CLUSTERSTATE_ACTIVE.Ptr()
					instance.Endpoints = []api.TidbEndpoint{endpoint}
					var settingCalls atomic.Int32
					var server *httptest.Server
					server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Path == "/ca" {
							if mode == "wrong CA" {
								_, _ = w.Write(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw}))
							} else {
								_, _ = w.Write(ca)
							}
							return
						}
						if auth == "digest" && !strings.HasPrefix(r.Header.Get("Authorization"), "Digest ") {
							w.Header().Set("WWW-Authenticate", `Digest realm="test", nonce="nonce", algorithm=MD5, qop="auth"`)
							w.WriteHeader(http.StatusUnauthorized)
							return
						}
						if auth == "bearer" && r.Header.Get("Authorization") != "Bearer test-token" {
							w.WriteHeader(http.StatusUnauthorized)
							return
						}
						w.Header().Set("Content-Type", "application/json")
						switch r.URL.Path {
						case "/v1beta2/tidbs/tidb-1":
							_ = json.NewEncoder(w).Encode(instance)
						case "/v1beta2/tidbs/tidb-1/publicConnectionSetting":
							settingCalls.Add(1)
							_, _ = fmt.Fprint(w, `{"enabled":true}`)
						case "/v1beta2/tidbs/tidb-1/caCertificateUrl":
							_ = json.NewEncoder(w).Encode(map[string]string{"uri": server.URL + "/ca"})
						default:
							http.NotFound(w, r)
						}
					}))
					defer server.Close()
					apiCA := filepath.Join(t.TempDir(), "api-ca.pem")
					require.NoError(t, os.WriteFile(apiCA, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw}), 0600))
					// Trust only the local HTTPS API for downloading the SQL CA;
					// SQL verification still uses the CA returned by that API.
					roots := x509.NewCertPool()
					roots.AddCert(server.Certificate())
					transport := http.DefaultTransport.(*http.Transport).Clone()
					transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12, RootCAs: roots}
					oldTransport := http.DefaultTransport
					http.DefaultTransport = transport
					defer func() { http.DefaultTransport = oldTransport; transport.CloseIdleConnections() }()
					var client *cloud.NextGenClientDelegate
					if auth == "bearer" {
						client, err = cloud.NewNextGenClientDelegateWithToken("test-token", server.URL, apiCA)
					} else {
						client, err = cloud.NewNextGenClientDelegateWithAPIKey("test-key", "test-secret", server.URL, apiCA)
					}
					require.NoError(t, err)
					input, err := os.CreateTemp(t.TempDir(), "sql-input")
					require.NoError(t, err)
					defer input.Close()
					query := "SELECT 'dns-recovery-ok' AS probe_token, CURRENT_USER() AS sql_user;\n"
					if plan.commandName == "essential-v2" && auth == "bearer" && mode == "valid" {
						// The initial 30-second deadline must not limit SQL queries.
						query = "SELECT SLEEP(31) AS waited, 'dns-recovery-ok' AS probe_token, CURRENT_USER() AS sql_user;\n"
					}
					_, err = input.WriteString(query + "SHOW STATUS LIKE 'Ssl_cipher';\n")
					require.NoError(t, err)
					_, err = input.Seek(0, 0)
					require.NoError(t, err)
					output, err := os.CreateTemp(t.TempDir(), "sql-output")
					require.NoError(t, err)
					defer output.Close()
					oldIn, oldOut := os.Stdin, os.Stdout
					os.Stdin, os.Stdout = input, output
					defer func() { os.Stdin, os.Stdout = oldIn, oldOut }()
					oldReadIn, oldReadOut, oldReadErr := readline.Stdin, readline.Stdout, readline.Stderr
					readline.Stdin, readline.Stdout, readline.Stderr = input, output, output
					defer func() { readline.Stdin, readline.Stdout, readline.Stderr = oldReadIn, oldReadOut, oldReadErr }()
					h := helperWithNextGenClient(client)
					h.IOStreams.CanPrompt = true
					cmd := shellCmd(h, plan)
					cmd.SilenceErrors, cmd.SilenceUsage = true, true
					cmd.SetArgs([]string{"-c", "tidb-1", "-u", user, "--password", password})
					err = cmd.Execute()
					require.EqualValues(t, 1, settingCalls.Load())
					require.Equal(t, 1, bytes.Count(h.IOStreams.Err.(*bytes.Buffer).Bytes(), []byte("Warning:")))
					switch mode {
					case "wrong CA":
						var target x509.UnknownAuthorityError
						require.ErrorAs(t, err, &target)
					case "wrong hostname":
						var target x509.HostnameError
						require.ErrorAs(t, err, &target)
					default:
						require.NoError(t, err)
						data, readErr := os.ReadFile(output.Name())
						require.NoError(t, readErr)
						require.Contains(t, string(data), "dns-recovery-ok")
						require.Contains(t, string(data), user+"@")
						require.Regexp(t, `TLS_AES_|ECDHE-`, string(data))
					}
				})
			}
		}
	}
}
