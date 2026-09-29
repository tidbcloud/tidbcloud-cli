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

package prop

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/juju/errors"
)

const (
	PublicKey          string = "public-key"
	PrivateKey         string = "private-key"
	CurProfile         string = "current-profile"
	ServerlessEndpoint string = "serverless-endpoint"
	NextGenEndpoint    string = "nextgen-endpoint"
	NextGenCACertPath  string = "nextgen-ca-cert-path"
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
	return []string{PublicKey, PrivateKey, ServerlessEndpoint, NextGenEndpoint, NextGenCACertPath, IAMEndpoint, OAuthEndpoint, OAuthClientID, OAuthClientSecret, TelemetryEnabled}
}

func ValidateApiUrl(value string) (*url.URL, error) {
	u, err := url.ParseRequestURI(value)
	if err != nil {
		return nil, errors.Annotate(err, "api url should format as <schema>://<host>")
	}
	return u, nil
}

func ValidateNextGenApiUrl(value string) (*url.URL, error) {
	if strings.Contains(value, "#") {
		return nil, fmt.Errorf("the API URL must not contain a query or fragment")
	}
	u, err := ValidateApiUrl(value)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("the API URL must use HTTPS")
	}
	if u.Host == "" {
		return nil, fmt.Errorf("the API URL must include a host")
	}
	if u.ForceQuery || u.RawQuery != "" || u.Fragment != "" {
		return nil, fmt.Errorf("the API URL must not contain a query or fragment")
	}
	return u, nil
}
