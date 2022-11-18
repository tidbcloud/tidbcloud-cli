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
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/util"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SetConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *SetConfigSuite) SetupTest() {
	err := os.Setenv("NO_COLOR", "true")
	if err != nil {
		suite.T().Error(err)
	}

	viper.Reset()
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName(".tidbcloud-cli")
	_ = viper.SafeWriteConfig()
	suite.h = &internal.Helper{
		IOStreams: &iostream.IOStreams{
			Out: &bytes.Buffer{},
			Err: &bytes.Buffer{},
		},
		Config: &config.Config{ActiveProfile: "test"},
	}
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

	viper.Set("test.public_key", publicKey)
	viper.Set("test.private_key", privateKey)
	viper.Set("current_profile", profile)
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
			args:         []string{"private_key", newPrivateKey},
			stdoutString: "Set profile `test` property `private_key` to value `TYTYTYYTYT` successfully\n",
		},
		{
			name: "set config with no args",
			args: []string{},
			err:  fmt.Errorf("missing arguments <propertyName, value> \n\nUsage:\n  set <propertyName> <value> [flags]\n\nFlags:\n  -h, --help   help for set\n"),
		},
		{
			name: "set config with unknown property",
			args: []string{"unknown", "value"},
			err:  fmt.Errorf("unrecognized property `unknown`, use `config set --help` to find available properties"),
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

			if err == nil {
				viper.Reset()
				viper.AddConfigPath(".")
				viper.SetConfigType("toml")
				viper.SetConfigName(".tidbcloud-cli")
				err := viper.ReadInConfig()
				assert.Nil(err)
				assert.Equal(profile, viper.GetString("current_profile"))
				assert.Equal(publicKey, viper.GetString(profile+".public_key"))
				assert.Equal(newPrivateKey, viper.GetString(profile+".private_key"))
			}
		})
	}
}

func (suite *SetConfigSuite) TestSetConfigWhenNoActiveProfile() {
	assert := require.New(suite.T())

	suite.h.Config.ActiveProfile = ""

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name: "set config",
			args: []string{"private_key", "value"},
			err:  fmt.Errorf("no profile is configured, please use `config init` to create a profile"),
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
