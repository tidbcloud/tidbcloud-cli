// Copyright 2024 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/prop"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/internal/version"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProfileSuite struct {
	suite.Suite
}

func (suite *ProfileSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	viper.Reset()
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName(".tidbcloud-cli")
	_ = viper.SafeWriteConfig()
}

func (suite *ProfileSuite) TearDownTest() {
	err := util.RemoveFile(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *ProfileSuite) TestTelemetryDefault() {
	assert := require.New(suite.T())
	version.Version = "v1.0.0"
	assert.False(TelemetryEnabled(), "telemetry should be disabled by default in release version")

	version.Version = "dev"
	assert.False(TelemetryEnabled(), "telemetry should be disabled in dev version")
}

func (suite *ProfileSuite) TestTelemetrySetFalse() {
	assert := require.New(suite.T())
	viper.Set(fmt.Sprintf("%s.%s", ActiveProfileName(), prop.TelemetryEnabled), false)
	version.Version = "v1.0.0"
	assert.False(TelemetryEnabled(), "telemetry should be disabled by config in release version")

	version.Version = "dev"
	assert.False(TelemetryEnabled(), "telemetry should be disabled in dev version")
}

func (suite *ProfileSuite) TestTelemetrySetTrue() {
	assert := require.New(suite.T())
	viper.Set(fmt.Sprintf("%s.%s", ActiveProfileName(), prop.TelemetryEnabled), true)
	version.Version = "v1.0.0"
	assert.True(TelemetryEnabled(), "telemetry should be enabled by config in release version")

	version.Version = "dev"
	assert.False(TelemetryEnabled(), "telemetry should be disabled in dev version")
}

func TestProfile(t *testing.T) {
	suite.Run(t, new(ProfileSuite))
}
