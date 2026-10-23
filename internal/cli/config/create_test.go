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

package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreateConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *CreateConfigSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	viper.Reset()
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName(".tidbcloud-cli")
	_ = viper.SafeWriteConfig()
	suite.h = &internal.Helper{
		IOStreams: iostream.Test(),
	}
}

func (suite *CreateConfigSuite) TearDownTest() {
	err := util.RemoveFile(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *CreateConfigSuite) TestCreateConfigArgs() {
	assert := require.New(suite.T())
	profile := "test"
	publicKey := "SDIWODIJQNDKJQW"
	privateKey := "SDWIOUEOSDSDC"

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "create config",
			args:         []string{"--profile-name", profile, "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Current profile has been changed to test",
		},
		{
			name:         "create config case-insensitive",
			args:         []string{"--profile-name", "teSt1", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Current profile has been changed to test1",
		},
		{
			name: "create config with 1 arg",
			args: []string{"arg1"},
			err:  fmt.Errorf(`unknown command "arg1" for "create"`),
		},
		{
			name:         "create config with special characters",
			args:         []string{"--profile-name", "~`!@#$%^&*()_+-={}[]\\|;:,<>/?", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Current profile has been changed to ~`!@#$%^&*()_+-={}[]\\|;:,<>/?",
		},
		{
			name:         "create config with special character '",
			args:         []string{"--profile-name", "'", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Current profile has been changed to '",
		},
		{
			name:         "create config with special character \"",
			args:         []string{"--profile-name", "\"", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Current profile has been changed to \"",
		},
		{
			name:         "create config with both \" and '",
			args:         []string{"--profile-name", "'\"", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.",
			err:          fmt.Errorf("profile name cannot contain both single and double quotes"),
		}, {
			name:         "create config with invalid characters .",
			args:         []string{"--profile-name", ".", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.",
			err:          fmt.Errorf("profile name cannot contain periods"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := viper.ReadInConfig()
			assert.Nil(err)

			cmd := CreateCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Equal(tt.err, err)

			assert.Contains(suite.h.IOStreams.Out.(*bytes.Buffer).String(), tt.stdoutString)
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
		})
	}
}

func (suite *CreateConfigSuite) TestCreateConfigWithExistedProfile() {
	assert := require.New(suite.T())
	profile := "test"
	publicKey := "SDIWODIJQNDKJQW"
	privateKey := "SDWIOUEOSDSDC"

	viper.Set("test.public-key", publicKey)
	viper.Set("test.private-key", privateKey)
	viper.Set("current-profile", profile)
	err := viper.WriteConfig()
	if err != nil {
		suite.T().Error(err)
	}

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "create config with existed profile",
			args:         []string{"--profile-name", profile, "--public-key", publicKey, "--private-key", privateKey},
			err:          fmt.Errorf("profile test already exists, use `config set` to modify"),
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := viper.ReadInConfig()
			assert.Nil(err)

			cmd := CreateCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Equal(tt.err, err)

			assert.Contains(suite.h.IOStreams.Out.(*bytes.Buffer).String(), tt.stdoutString)
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
		})
	}
}

func TestCreateConfigSuite(t *testing.T) {
	suite.Run(t, new(CreateConfigSuite))
}

func TestCreateEnvironment(t *testing.T) {
	for _, tt := range []struct {
		name, environment, api, oauth, clientID, iam string
	}{
		{"default", "", "https://cloud.tidbapi.com", config.OAuthEndpoint, config.ClientID, "https://iam.tidbapi.com"},
		{"production", "prod", "https://cloud.tidbapi.com", config.OAuthEndpoint, config.ClientID, "https://iam.tidbapi.com"},
		{"staging", "staging", "https://cloud.staging.tidbapi.com", config.StagingOAuthEndpoint, config.StagingClientID, "https://iam.staging.tidbapi.com"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			configPath := filepath.Join(t.TempDir(), "config.toml")
			viper.SetConfigFile(configPath)
			viper.Set("existing.public-key", "preserved")
			viper.Set("existing.iam-endpoint", "https://custom.example.com")
			require.NoError(t, viper.WriteConfig())
			cmd := CreateCmd(&internal.Helper{IOStreams: iostream.Test()})
			args := []string{"--profile-name", "new-profile"}
			if tt.environment != "" {
				args = append(args, "--environment", tt.environment)
			}
			cmd.SetArgs(args)
			require.NoError(t, cmd.Execute())
			viper.Reset()
			viper.SetConfigFile(configPath)
			require.NoError(t, viper.ReadInConfig())
			require.Equal(t, tt.api, viper.GetString("new-profile.cloud-api-endpoint"))
			require.Equal(t, tt.oauth, viper.GetString("new-profile.oauth-endpoint"))
			require.Equal(t, tt.iam, viper.GetString("new-profile.iam-endpoint"))
			require.Equal(t, "https://custom.example.com", viper.GetString("existing.iam-endpoint"))
			require.Empty(t, viper.GetString("new-profile.oauth-client-id"))
			require.Equal(t, "preserved", viper.GetString("existing.public-key"))
			previous := config.ActiveProfileName()
			config.SetActiveProfile("new-profile")
			require.Equal(t, tt.clientID, config.GetOAuthClientID())
			if previous != "" {
				config.SetActiveProfile(previous)
			}
		})
	}
}

func TestCreateInvalidEnvironment(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	configPath := filepath.Join(t.TempDir(), "config.toml")
	viper.SetConfigFile(configPath)
	cmd := CreateCmd(&internal.Helper{IOStreams: iostream.Test()})
	cmd.SetArgs([]string{"--profile-name", "bad", "--environment", "typo"})
	require.ErrorContains(t, cmd.Execute(), "choose prod or staging")
	_, err := os.Stat(configPath)
	require.True(t, os.IsNotExist(err))
}
