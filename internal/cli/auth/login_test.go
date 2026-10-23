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
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
)

func TestValidateLoginEnvironment(t *testing.T) {
	for _, tt := range []struct {
		name, api, oauth, want string
	}{
		{"production", "", config.OAuthEndpoint, ""},
		{"staging", "https://cloud.staging.tidbapi.com", "https://oauth.staging.tidbcloud.com", ""},
		{"staging API production OAuth", "https://cloud.staging.tidbapi.com/v1beta2/", config.OAuthEndpoint, config.StagingOAuthEndpoint},
		{"Serverless staging without Cloud API endpoint", "", config.StagingOAuthEndpoint, ""},
		{"explicit production API staging OAuth", "https://cloud.tidbapi.com", "https://oauth.staging.tidbcloud.com", config.OAuthEndpoint},
		{"case and port", "https://CLOUD.STAGING.TIDBAPI.COM:443/", config.OAuthEndpoint, config.StagingOAuthEndpoint},
		{"custom API", "https://gateway.example.com", config.OAuthEndpoint, ""},
		{"custom OAuth", "https://cloud.staging.tidbapi.com", "https://auth.example.com", ""},
		{"host suffix is not known environment", "https://cloud.staging.tidbapi.com.example.com", config.OAuthEndpoint, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLoginEnvironment(tt.api, tt.oauth, "selected-profile")
			if tt.want == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, "OAuth environment mismatch")
			require.ErrorContains(t, err, "config set oauth-endpoint "+tt.want)
			require.ErrorContains(t, err, `-P "selected-profile"`)
			require.NotContains(t, err.Error(), "%!")
		})
	}
}

func TestLoginRejectsMixedEnvironmentBeforeSideEffects(t *testing.T) {
	previousProfile := config.ActiveProfileName()
	viper.Reset()
	t.Cleanup(func() {
		viper.Reset()
		config.SetActiveProfile(previousProfile)
	})
	config.SetActiveProfile("mixed")
	viper.Set("mixed.cloud-api-endpoint", "https://cloud.staging.tidbapi.com")
	// API keys would prompt for confirmation if validation did not run first.
	viper.Set("mixed.public-key", "test")
	viper.Set("mixed.private-key", "test")
	// No IO streams: a prompt, browser action, or credential access must not occur.
	cmd := LoginCmd(&internal.Helper{})
	err := cmd.RunE(cmd, nil)
	require.ErrorContains(t, err, "OAuth environment mismatch")
	require.ErrorContains(t, err, "config set oauth-endpoint https://oauth.staging.tidbcloud.com")
}
