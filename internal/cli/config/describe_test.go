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
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DescribeConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *DescribeConfigSuite) SetupTest() {
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

func (suite *DescribeConfigSuite) TearDownTest() {
	err := util.RemoveFile(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *DescribeConfigSuite) TestDescribeConfigArgs() {
	assert := require.New(suite.T())

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "describe config",
			args:         []string{"test"},
			stdoutString: "{\n  \"private-key\": \"SDWIOUEOSDSDC\",\n  \"public-key\": \"SDIWODIJQNDKJQW\"\n}\n",
		},
		{
			name:         "describe config case-insensitive",
			args:         []string{"teSt"},
			stdoutString: "{\n  \"private-key\": \"SDWIOUEOSDSDC\",\n  \"public-key\": \"SDIWODIJQNDKJQW\"\n}\n",
		},
		{
			name: "describe config with no args",
			args: []string{},
			err:  fmt.Errorf("accepts 1 arg(s), received 0"),
		},
		{
			name: "describe config with non-existed profile",
			args: []string{"test1"},
			err:  fmt.Errorf("profile test1 not found"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := viper.ReadInConfig()
			assert.Nil(err)

			cmd := DescribeCmd(suite.h)
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

func (suite *DescribeConfigSuite) TestDescribeConfigWithSpecialCharacters() {
	assert := require.New(suite.T())
	newProfile := "~`!@#$%^&*()_+-={}[]\\|;:,<>/?"
	publicKey := "SDIWODIJQNDKJQW"
	privateKey := "SDWIOUEOSDSDC"

	viper.Set("~`!@#$%^&*()_+-={}[]\\|;:,<>/?.public-key", publicKey)
	viper.Set("~`!@#$%^&*()_+-={}[]\\|;:,<>/?.private-key", privateKey)
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
			name:         "describe active profile",
			args:         []string{"~`!@#$%^&*()_+-={}[]\\|;:,<>/?"},
			stdoutString: "{\n  \"private-key\": \"SDWIOUEOSDSDC\",\n  \"public-key\": \"SDIWODIJQNDKJQW\"\n}\n",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DescribeCmd(suite.h)
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

func TestDescribeConfigSuite(t *testing.T) {
	suite.Run(t, new(DescribeConfigSuite))
}
