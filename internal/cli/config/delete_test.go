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

type DeleteConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *DeleteConfigSuite) SetupTest() {
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
}

func (suite *DeleteConfigSuite) TearDownTest() {
	err := util.RemoveFile(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *DeleteConfigSuite) TestDeleteConfigArgs() {
	assert := require.New(suite.T())

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "delete config",
			args:         []string{"test", "--force"},
			stdoutString: "Profile test deleted successfully\n",
		},
		{
			name: "delete config without force",
			args: []string{"test"},
			err:  fmt.Errorf("the terminal doesn't support prompt, please run with --force to delete the profile"),
		},
		{
			name: "delete config with no args",
			args: []string{"--force"},
			err:  fmt.Errorf("accepts 1 arg(s), received 0"),
		},
		{
			name: "delete config with non-existed profile",
			args: []string{"test1", "--force"},
			err:  fmt.Errorf("profile `test1` does not exist"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := viper.ReadInConfig()
			assert.Nil(err)

			cmd := DeleteCmd(suite.h)
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
				assert.Equal("", viper.GetString("current-profile"))
				assert.Equal("", viper.GetString(tt.args[0]+".public-key"))
				assert.Equal("", viper.GetString(tt.args[0]+".private-key"))
				assert.Equal("", viper.GetString(tt.args[0]))
			}
		})
	}
}

func (suite *DeleteConfigSuite) TestDeleteConfigWithActiveProfile() {
	assert := require.New(suite.T())
	newProfile := "newTest"
	publicKey := "SDIWODIJQNDKJQW"
	privateKey := "SDWIOUEOSDSDC"

	viper.Set("newTest.public-key", publicKey)
	viper.Set("newTest.private-key", privateKey)
	viper.Set("current-profile", newProfile)

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
			name:         "delete active profile",
			args:         []string{"newtest", "--force"},
			stdoutString: "Profile newtest deleted successfully\n",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DeleteCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())

			viper.Reset()
			viper.AddConfigPath(".")
			viper.SetConfigType("toml")
			viper.SetConfigName(".tidbcloud-cli")
			err = viper.ReadInConfig()
			assert.Nil(err)
			assert.Equal("test", viper.GetString("current-profile"))
			assert.Equal("", viper.GetString(tt.args[0]+".public-key"))
			assert.Equal("", viper.GetString(tt.args[0]+".private-key"))
			assert.Equal("", viper.GetString(tt.args[0]))
		})
	}
}

func (suite *DeleteConfigSuite) TestDeleteConfigWithSpecialCharactersSingleQuote() {
	assert := require.New(suite.T())
	newProfile := "'~`!'@#$%^&*()_+-={}[]\\|;:,<>/?'"
	publicKey := "SDIWODIJQNDKJQW"
	privateKey := "SDWIOUEOSDSDC"

	viper.Set(fmt.Sprintf("%s.public-key", newProfile), publicKey)
	viper.Set(fmt.Sprintf("%s.private-key", newProfile), privateKey)
	viper.Set("current-profile", newProfile)

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
			name:         "delete active profile",
			args:         []string{newProfile, "--force"},
			stdoutString: "Profile '~`!'@#$%^&*()_+-={}[]\\|;:,<>/?' deleted successfully\n",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DeleteCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())

			viper.Reset()
			viper.AddConfigPath(".")
			viper.SetConfigType("toml")
			viper.SetConfigName(".tidbcloud-cli")
			err = viper.ReadInConfig()
			assert.Nil(err)
			assert.Equal("test", viper.GetString("current-profile"))
			assert.Equal("", viper.GetString(tt.args[0]+".public-key"))
			assert.Equal("", viper.GetString(tt.args[0]+".private-key"))
			assert.Equal("", viper.GetString(tt.args[0]))
		})
	}
}

func (suite *DeleteConfigSuite) TestDeleteConfigWithSpecialCharactersDoubleQuote() {
	assert := require.New(suite.T())
	newProfile := "\"~`!\"@#$%^&*()_+-={}[]\\|;:,<>/?\""
	publicKey := "SDIWODIJQNDKJQW"
	privateKey := "SDWIOUEOSDSDC"

	viper.Set(fmt.Sprintf("%s.public-key", newProfile), publicKey)
	viper.Set(fmt.Sprintf("%s.private-key", newProfile), privateKey)
	viper.Set("current-profile", newProfile)

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
			name:         "delete active profile",
			args:         []string{newProfile, "--force"},
			stdoutString: "Profile \"~`!\"@#$%^&*()_+-={}[]\\|;:,<>/?\" deleted successfully\n",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DeleteCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())

			viper.Reset()
			viper.AddConfigPath(".")
			viper.SetConfigType("toml")
			viper.SetConfigName(".tidbcloud-cli")
			err = viper.ReadInConfig()
			assert.Nil(err)
			assert.Equal("test", viper.GetString("current-profile"))
			assert.Equal("", viper.GetString(tt.args[0]+".public-key"))
			assert.Equal("", viper.GetString(tt.args[0]+".private-key"))
			assert.Equal("", viper.GetString(tt.args[0]))
		})
	}
}

func TestDeleteConfigSuite(t *testing.T) {
	suite.Run(t, new(DeleteConfigSuite))
}
