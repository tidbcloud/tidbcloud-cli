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
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/util"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UseConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *UseConfigSuite) SetupTest() {
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

	viper.Set("test.public_key", publicKey)
	viper.Set("test.private_key", privateKey)
	viper.Set("current_profile", profile)
	viper.Set("test1.public_key", publicKey)
	viper.Set("test1.private_key", privateKey)
	err := viper.WriteConfig()
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *UseConfigSuite) TearDownTest() {
	err := util.RemoveFile(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *UseConfigSuite) TestUseConfigArgs() {
	assert := require.New(suite.T())

	err := viper.ReadInConfig()
	assert.Nil(err)
	assert.Equal("test", viper.GetString("current_profile"))

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "use config",
			args:         []string{"test1"},
			stdoutString: "Current profile has been changed to test1\n",
		},
		{
			name: "use config with no args",
			args: []string{},
			err:  fmt.Errorf("missing argument <profileName> \n\nUsage:\n  use <profileName> [flags]\n\nExamples:\n  Use the \"test\" profile as the active profile:\n  $ ticloud config use test\n\nFlags:\n  -h, --help   help for use\n"),
		},
		{
			name: "use config with non-existed profile",
			args: []string{"test2"},
			err:  fmt.Errorf("profile test2 not found"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := viper.ReadInConfig()
			assert.Nil(err)

			cmd := UseCmd(suite.h)
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
				assert.Equal("test1", viper.GetString("current_profile"))
			}
		})
	}
}

func TestUseConfigSuite(t *testing.T) {
	suite.Run(t, new(UseConfigSuite))
}
