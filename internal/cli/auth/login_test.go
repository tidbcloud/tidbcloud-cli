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

package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
)

func TestValidateVerificationURL(t *testing.T) {
	for _, value := range []string{
		"https://tidbcloud.com/verify?user_code=abc&scope=openid",
		"https://login.staging.example.test/device?code=a%2Bb%26c",
		"https://localhost:8443/verify",
		"https://[2001:db8::1]:443/verify",
	} {
		t.Run(value, func(t *testing.T) { require.NoError(t, validateVerificationURL(value)) })
	}
	for _, value := range invalidVerificationURLs() {
		t.Run(value, func(t *testing.T) {
			err := validateVerificationURL(value)
			require.ErrorContains(t, err, "invalid authorization URL")
			require.NotContains(t, err.Error(), "code-secret")
		})
	}
}

func invalidVerificationURLs() []string {
	return []string{
		"", "/verify?code=code-secret", "//example.test/verify?code=code-secret",
		"http://example.test/?code=code-secret", "file:///tmp/code-secret", "javascript:code-secret",
		"https:///verify?code=code-secret", "https://user:code-secret@example.test/verify",
		"https://example.test:0/?code=code-secret", "https://example.test:65536/?code=code-secret",
		"https://example.test:bad/?code=code-secret", "https://example.test:/?code=code-secret",
		"https://example.test/\ncode-secret", "https://example.test/\x1bcode-secret",
		"https://example.test/\u0085code-secret", "https://example.test/\u2028code-secret",
		"https://example.test/ code-secret", `https://example.test\code-secret`,
	}
}

type authorizationTransport func(*http.Request) (*http.Response, error)

func (f authorizationTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// A rejected server-provided URL must reach neither terminal mode nor the OS,
// and must stop the flow before token polling or credential persistence.
func TestLoginRejectsUnsafeVerificationURLBeforeOutput(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	oldProfile := config.ActiveProfileName()
	config.SetActiveProfile("authorization-test")
	t.Cleanup(func() { config.SetActiveProfile(oldProfile) })
	originalTransport := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = originalTransport })
	for _, interactive := range []bool{false, true} {
		for _, value := range invalidVerificationURLs() {
			t.Run(fmt.Sprintf("interactive=%t/%s", interactive, value), func(t *testing.T) {
				requests := 0
				http.DefaultTransport = authorizationTransport(func(r *http.Request) (*http.Response, error) {
					requests++
					require.Equal(t, authPath, r.URL.Path)
					require.Equal(t, http.MethodPost, r.Method)
					data, err := json.Marshal(AuthResponse{VerificationURIComplete: value, UserCode: "confirmation-secret"})
					require.NoError(t, err)
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": {"application/json"}},
						Body:       io.NopCloser(strings.NewReader(string(data))),
						Request:    r,
					}, nil
				})
				h := &internal.Helper{IOStreams: iostream.Test()}
				h.IOStreams.CanPrompt = interactive
				cmd := LoginCmd(h)
				cmd.Flags().Bool(flag.Debug, false, "")
				cmd.SilenceErrors, cmd.SilenceUsage = true, true
				cmd.SetArgs([]string{"--insecure-storage"})
				err := cmd.ExecuteContext(t.Context())
				require.ErrorContains(t, err, "invalid authorization URL")
				require.NotContains(t, err.Error(), "code-secret")
				require.Equal(t, 1, requests)
				require.Empty(t, h.IOStreams.Out.(*bytes.Buffer).String())
				require.Empty(t, h.IOStreams.Err.(*bytes.Buffer).String())
			})
		}
	}
}
