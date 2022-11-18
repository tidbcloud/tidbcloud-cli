package cli

import (
	"bytes"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/iostream"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RootCmdSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *RootCmdSuite) SetupTest() {
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
		Config: &config.Config{
			ActiveProfile: "",
		},
	}
}

func (suite *RootCmdSuite) TestFlagProfile() {
	assert := require.New(suite.T())

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
			args:          []string{"config", "set", "private_key", privateKey3},
			stdoutString:  "Set profile `test` property `private_key` to value `324OPIFO2423423DFO` successfully\n",
			propertyKey:   "test.private_key",
			propertyValue: "324OPIFO2423423DFO",
		},
		{
			name:          "test flag --profile",
			args:          []string{"config", "set", "private_key", privateKey1, "--profile", "test1"},
			stdoutString:  "Set profile `test1` property `private_key` to value `SAJKGDUYAKGD` successfully\n",
			propertyKey:   "test1.private_key",
			propertyValue: "SAJKGDUYAKGD",
		},
		{
			name:          "test flag -P",
			args:          []string{"config", "set", "private_key", privateKey2, "-P", "test1"},
			stdoutString:  "Set profile `test1` property `private_key` to value `{OPIFOPIDFO` successfully\n",
			propertyKey:   "test1.private_key",
			propertyValue: "{OPIFOPIDFO",
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
			assert.Equal("test", viper.GetString("current_profile"))
			assert.Equal(tt.propertyValue, viper.GetString(tt.propertyKey))
		})
	}
}

func TestRootCmdSuite(t *testing.T) {
	suite.Run(t, new(RootCmdSuite))
}
