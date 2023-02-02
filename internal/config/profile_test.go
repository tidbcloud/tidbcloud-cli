package config

import (
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/internal/version"

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
	assert.True(TelemetryEnabled(), "telemetry should be enabled by default in release version")

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
