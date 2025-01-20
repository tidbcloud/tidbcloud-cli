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

package config

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
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
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.\nCurrent profile has been changed to test\n",
		},
		{
			name:         "create config case-insensitive",
			args:         []string{"--profile-name", "teSt1", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.\nCurrent profile has been changed to test1\n",
		},
		{
			name: "create config with 1 arg",
			args: []string{"arg1"},
			err:  fmt.Errorf(`unknown command "arg1" for "create"`),
		},
		{
			name:         "create config with special characters",
			args:         []string{"--profile-name", "~`!@#$%^&*()_+-={}[]\\|;:,<>/?", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.\nCurrent profile has been changed to ~`!@#$%^&*()_+-={}[]\\|;:,<>/?\n",
		},
		{
			name:         "create config with special character '",
			args:         []string{"--profile-name", "'", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.\nCurrent profile has been changed to '\n",
		},
		{
			name:         "create config with special character \"",
			args:         []string{"--profile-name", "\"", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.\nCurrent profile has been changed to \"\n",
		},
		{
			name:         "create config with both \" and '",
			args:         []string{"--profile-name", "'\"", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.\n",
			err:          fmt.Errorf("profile name cannot contain both single and double quotes"),
		}, {
			name:         "create config with invalid characters .",
			args:         []string{"--profile-name", ".", "--public-key", publicKey, "--private-key", privateKey},
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.\n",
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

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
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
			stdoutString: "Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.\n",
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

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
		})
	}
}

func TestCreateConfigSuite(t *testing.T) {
	suite.Run(t, new(CreateConfigSuite))
}
