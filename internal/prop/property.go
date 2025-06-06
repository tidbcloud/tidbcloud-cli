// Copyright 2025 PingCAP, Inc.
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

package prop

import (
	"net/url"

	"github.com/juju/errors"
)

const (
	PublicKey          string = "public-key"
	PrivateKey         string = "private-key"
	CurProfile         string = "current-profile"
	ApiUrl             string = "api-url"
	ServerlessEndpoint string = "serverless-endpoint"
	IAMEndpoint        string = "iam-endpoint"
	OAuthEndpoint      string = "oauth-endpoint"
	OAuthClientID      string = "oauth-client-id"
	OAuthClientSecret  string = "oauth-client-secret"
	TelemetryEnabled   string = "telemetry-enabled"

	// shall not be set by user
	TokenExpiredAt string = "token-expired-at"
	TokenType      string = "token-type"
	AccessToken    string = "access-token"
)

func GlobalProperties() []string {
	return []string{CurProfile}
}

func ProfileProperties() []string {
	return []string{PublicKey, PrivateKey, ApiUrl, ServerlessEndpoint, IAMEndpoint, OAuthEndpoint, OAuthClientID, OAuthClientSecret, TelemetryEnabled}
}

func ValidateApiUrl(value string) (*url.URL, error) {
	u, err := url.ParseRequestURI(value)
	if err != nil {
		return nil, errors.Annotate(err, "api url should format as <schema>://<host>")
	}
	return u, nil
}
