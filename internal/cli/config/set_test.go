// Copyright 2022 PingCAP, Inc.
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

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/util"

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
			stdoutString: "Set profile `test` property `private-key` to value `TYTYTYYTYT` successfully\n",
		},
		{
			name: "set config with no args",
			args: []string{},
			err:  fmt.Errorf("missing arguments <property-name, value> \n\nUsage:\n  set <property-name> <value> [flags]\n\nExamples:\n  Set the value of the public-key in active profile:\n  $ ticloud config set public-key <public-key>\n\n  Set the value of the public-key in the specific profile \"test\":\n  $ ticloud config set public-key <public-key> -P test\n\nFlags:\n  -h, --help   help for set\n"),
		},
		{
			name: "set config with unknown property",
			args: []string{"unknown", "value"},
			err:  fmt.Errorf("unrecognized property `unknown`, use `config set --help` to find available properties"),
		},
		{
			name: "set config with unknown property",
			args: []string{"api-url", "baidu.com"},
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
			stdoutString: "Set profile `default` property `private-key` to value `value` successfully\n",
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

func TestSetConfigSuite(t *testing.T) {
	suite.Run(t, new(SetConfigSuite))
}
