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
	"net/url"
	"os"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	"github.com/juju/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SetConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *SetConfigSuite) SetupTest() {
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
	config.SetActiveProfile("test")
}

func (suite *SetConfigSuite) TearDownTest() {
	err := util.RemoveFile(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *SetConfigSuite) TestSetConfigArgs() {
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
	newPrivateKey := "TYTYTYYTYT"

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "set config",
			args:         []string{"private-key", newPrivateKey},
			stdoutString: "Set profile `test` property `private-key` to value `***` successfully\n",
		},
		{
			name: "set config with no args",
			args: []string{},
			err:  fmt.Errorf("accepts 2 arg(s), received 0"),
		},
		{
			name: "set config with unknown property",
			args: []string{"unknown", "value"},
			err:  fmt.Errorf("unrecognized property `unknown`, use `config set --help` to find available properties"),
		},
		{
			name: "set config with unknown property",
			args: []string{"serverless-endpoint", "baidu.com"},
			err: errors.Annotate(&url.Error{
				Op:  "parse",
				URL: "baidu.com",
				Err: fmt.Errorf("invalid URI for request"),
			}, "api url should format as <schema>://<host>"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := viper.ReadInConfig()
			assert.Nil(err)

			cmd := SetCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			if err != nil {
				assert.EqualError(tt.err, err.Error())
			} else {
				assert.Equal(tt.err, err)
			}

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())

			if err == nil {
				viper.Reset()
				viper.AddConfigPath(".")
				viper.SetConfigType("toml")
				viper.SetConfigName(".tidbcloud-cli")
				err := viper.ReadInConfig()
				assert.Nil(err)
				assert.Equal(profile, viper.GetString("current-profile"))
				assert.Equal(publicKey, viper.GetString(profile+".public-key"))
				assert.Equal(newPrivateKey, viper.GetString(profile+".private-key"))
			}
		})
	}
}

func (suite *SetConfigSuite) TestSetConfigWhenNoActiveProfile() {
	assert := require.New(suite.T())

	config.SetActiveProfile("")

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "set config",
			args:         []string{"private-key", "value"},
			stdoutString: "Set profile `default` property `private-key` to value `***` successfully\n",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := viper.ReadInConfig()
			assert.Nil(err)

			cmd := SetCmd(suite.h)
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

func (suite *SetConfigSuite) TestSetCloudAPIProperties() {
	for _, tt := range []struct {
		property string
		value    string
		get      func() string
	}{
		{"cloud-api-endpoint", "https://gateway.example.com", config.GetCloudAPIEndpoint},
		{"cloud-api-ca-cert-path", "/tmp/management-ca.pem", config.GetCloudAPICACertPath},
	} {
		cmd := SetCmd(suite.h)
		cmd.SetArgs([]string{tt.property, tt.value})
		require.NoError(suite.T(), cmd.Execute())
		require.Equal(suite.T(), tt.value, tt.get())
	}

	// Read the persisted profile without Viper's in-memory overrides.
	viper.Reset()
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName(".tidbcloud-cli")
	require.NoError(suite.T(), viper.ReadInConfig())
	require.Equal(suite.T(), "https://gateway.example.com", config.GetCloudAPIEndpoint())
	require.Equal(suite.T(), "/tmp/management-ca.pem", config.GetCloudAPICACertPath())
	config.SetActiveProfile("other")
	require.Empty(suite.T(), config.GetCloudAPIEndpoint())
	require.Empty(suite.T(), config.GetCloudAPICACertPath())
	config.SetActiveProfile("test")
}

func (suite *SetConfigSuite) TestSetCloudAPIEndpointValidation() {
	tests := []struct {
		name      string
		endpoint  string
		errString string
	}{
		{
			name:      "reject HTTP endpoint",
			endpoint:  "http://example.com",
			errString: "the API URL must use HTTPS",
		},
		{
			name:      "reject endpoint with query",
			endpoint:  "https://example.com?region=us-west-2",
			errString: "the API URL must not contain a query or fragment",
		},
		{
			name:      "reject endpoint with empty query",
			endpoint:  "https://example.com?",
			errString: "the API URL must not contain a query or fragment",
		},
		{
			name:      "reject endpoint with fragment",
			endpoint:  "https://example.com#fragment",
			errString: "the API URL must not contain a query or fragment",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := SetCmd(suite.h)
			cmd.SetArgs([]string{"cloud-api-endpoint", tt.endpoint})
			err := cmd.Execute()
			require.EqualError(t, err, tt.errString)
			require.Empty(t, viper.GetString("test.cloud-api-endpoint"))
		})
	}
}

func TestSetConfigSuite(t *testing.T) {
	suite.Run(t, new(SetConfigSuite))
}
