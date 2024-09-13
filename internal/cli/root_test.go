// Copyright 2024 PingCAP, Inc.
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

package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RootCmdSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *RootCmdSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName(".tidbcloud-cli")
	_ = viper.SafeWriteConfig()
	suite.h = &internal.Helper{
		IOStreams:     iostream.Test(),
		QueryPageSize: 10,
	}
}

func (suite *RootCmdSuite) TearDownTest() {
	err := util.RemoveFile(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *RootCmdSuite) TestFlagProfile() {
	assert := require.New(suite.T())

	profile := "test"
	publicKey := "SDIWODIJQNDKJQW"
	privateKey := "SDWIOUEOSDSDC"

	viper.Set("test.public-key", publicKey)
	viper.Set("test.private-key", privateKey)
	viper.Set("current-profile", profile)
	viper.Set("test1.public-key", publicKey)
	viper.Set("test1.private-key", privateKey)
	err := viper.WriteConfig()
	if err != nil {
		suite.T().Error(err)
	}

	privateKey1 := "SAJKGDUYAKGD"
	privateKey2 := "{OPIFOPIDFO"
	privateKey3 := "324OPIFO2423423DFO"

	tests := []struct {
		name          string
		args          []string
		err           error
		stdoutString  string
		stderrString  string
		propertyKey   string
		propertyValue string
	}{
		{
			name:          "test without flag profile",
			args:          []string{"config", "set", "private-key", privateKey3},
			stdoutString:  "Set profile `test` property `private-key` to value `324OPIFO2423423DFO` successfully\n",
			propertyKey:   "test.private-key",
			propertyValue: "324OPIFO2423423DFO",
		},
		{
			name:          "test flag --profile",
			args:          []string{"config", "set", "private-key", privateKey1, "--profile", "test1"},
			stdoutString:  "Set profile `test1` property `private-key` to value `SAJKGDUYAKGD` successfully\n",
			propertyKey:   "test1.private-key",
			propertyValue: "SAJKGDUYAKGD",
		},
		{
			name:          "test flag -P",
			args:          []string{"config", "set", "private-key", privateKey2, "-P", "test1"},
			stdoutString:  "Set profile `test1` property `private-key` to value `{OPIFOPIDFO` successfully\n",
			propertyKey:   "test1.private-key",
			propertyValue: "{OPIFOPIDFO",
		},
		{
			name:          "test flag -P case-insensitive",
			args:          []string{"config", "set", "private-key", "SADASDIDFO", "-P", "tESt1"},
			stdoutString:  "Set profile `test1` property `private-key` to value `SADASDIDFO` successfully\n",
			propertyKey:   "test1.private-key",
			propertyValue: "SADASDIDFO",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := RootCmd(suite.h)
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
			assert.Equal(tt.propertyValue, viper.GetString(tt.propertyKey))
		})
	}
}

func TestRootCmdSuite(t *testing.T) {
	suite.Run(t, new(RootCmdSuite))
}
