package config

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UseConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *UseConfigSuite) SetupTest() {
	err := os.Setenv("NO_COLOR", "true")
	if err != nil {
		suite.T().Error(err)
	}

	fs := afero.NewOsFs()
	err = fs.Remove(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}

	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName(".tidbcloud-cli")
	_ = viper.SafeWriteConfig()
	suite.h = &internal.Helper{
		IOStreams: &iostream.IOStreams{
			Out: &bytes.Buffer{},
			Err: &bytes.Buffer{},
		},
	}

	profile := "test"
	publicKey := "SDIWODIJQNDKJQW"
	privateKey := "SDWIOUEOSDSDC"

	viper.Set("test.public_key", publicKey)
	viper.Set("test.private_key", privateKey)
	viper.Set("current_profile", profile)
	viper.Set("test1.public_key", publicKey)
	viper.Set("test1.private_key", privateKey)
	err = viper.WriteConfig()
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
			err:  fmt.Errorf("missing argument <profileName> \n\nUsage:\n  use <profileName> [flags]\n\nFlags:\n  -h, --help   help for use\n"),
		},
		{
			name:         "use config with non-existed profile",
			args:         []string{"test2"},
			err:          fmt.Errorf("profile test2 not found"),
			stdoutString: "",
			stderrString: "",
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
