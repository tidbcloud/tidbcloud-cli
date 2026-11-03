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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
)

func TestHTTPRequestsBlocked(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP endpoint received a request")
	}))
	defer server.Close()
	client := newAuthClient()
	for _, path := range []string{authPath, accessPath, revokePath, userInfoPath, orgPath} {
		_, err := client.R().SetAuthToken("test-token").SetBody("test-secret").Post(server.URL + path)
		require.ErrorContains(t, err, "api url must use HTTPS")
	}
}

func TestConfiguredHTTPEndpointsBlocked(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	previous := config.ActiveProfileName()
	t.Cleanup(func() { config.SetActiveProfile(previous) })
	config.SetActiveProfile("test")
	viper.Set("test.oauth-endpoint", "http://127.0.0.1:1")
	viper.Set("test.iam-endpoint", "http://127.0.0.1:1")
	require.ErrorContains(t, getUserInfo(context.Background(), newAuthClient(), "test-token").Error, "api url must use HTTPS")
	require.ErrorContains(t, getOrgInfo(context.Background(), newAuthClient(), "test-token").Error, "api url must use HTTPS")
}

func TestHTTPSRedirectToHTTPBlocked(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP redirect target received a request")
	}))
	defer target.Close()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer server.Close()
	client := newAuthClient().SetTransport(httpsTransport{inner: server.Client().Transport})
	_, err := client.R().SetAuthToken("test-token").Get(server.URL)
	require.ErrorContains(t, err, "api url must use HTTPS")
}
