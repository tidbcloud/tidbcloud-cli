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

type ListConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *ListConfigSuite) SetupTest() {
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

func (suite *ListConfigSuite) TearDownTest() {
	err := util.RemoveFile(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *ListConfigSuite) TestListConfigArgs() {
	assert := require.New(suite.T())

	profile := "test"
	publicKey := "SDIWODIJQNDKJQW"
	privateKey := "SDWIOUEOSDSDC"
	viper.Set("test.public-key", publicKey)
	viper.Set("test.private-key", privateKey)
	viper.Set("current-profile", profile)
	viper.Set("test1.public-key", publicKey)
	viper.Set("test1.private-key", privateKey)
	viper.Set("'~`!@#$%'^&*()_+-={}[]\\|;:,<>/?'.public-key", publicKey)
	viper.Set("'~`!@#$%'^&*()_+-={}[]\\|;:,<>/?'.private-key", privateKey)
	viper.Set("\"~`!@#$%^&*\"()_+-={}[]\\|;:,<>/?\".public-key", publicKey)
	viper.Set("\"~`!@#$%^&*\"()_+-={}[]\\|;:,<>/?\".private-key", privateKey)
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
			name:         "list config",
			args:         []string{},
			stdoutString: "Profile Name\n\"~`!@#$%^&*\"()_+-={}[]\\|;:,<>/?\"\n'~`!@#$%'^&*()_+-={}[]\\|;:,<>/?'\ntest\t(active)\ntest1\n",
		},
		{
			name: "list config with 1 arg",
			args: []string{"arg1"},
			err:  fmt.Errorf(`unknown command "arg1" for "list"`),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := viper.ReadInConfig()
			assert.Nil(err)

			cmd := ListCmd(suite.h)
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

func TestListConfigSuite(t *testing.T) {
	suite.Run(t, new(ListConfigSuite))
}
