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
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/prop"

	"github.com/spf13/cobra"
)

const (
	tokenTypeHint = "access_token"
	grantType     = "urn:ietf:params:oauth:grant-type:device_code"

	authPath     = "/v1/device_authorization"
	accessPath   = "/v1/token"
	revokePath   = "/v1/revoke"
	userInfoPath = "/v1/userinfo"
	orgPath      = "/v1beta1/org"

	errSlowDown = "slow_down"
	errPending  = "authorization_pending"
)

func AuthCmd(h *internal.Helper) *cobra.Command {
	var authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Login and logout via TiDB Cloud API",
	}

	authCmd.AddCommand(LoginCmd(h))
	authCmd.AddCommand(LogoutCmd(h))
	authCmd.AddCommand(WhoamiCmd(h))
	return authCmd
}

// Validate each outgoing request, including redirects, so endpoints loaded
// directly from the config cannot send credentials over HTTP.
type httpsTransport struct {
	inner http.RoundTripper
}

func (t httpsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if _, err := prop.ValidateApiUrl(req.URL.String()); err != nil {
		return nil, err
	}
	return t.inner.RoundTrip(req)
}

func newAuthClient() *resty.Client {
	return resty.New().SetTransport(httpsTransport{inner: http.DefaultTransport})
}
