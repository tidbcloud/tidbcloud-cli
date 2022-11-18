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

type DeleteConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *DeleteConfigSuite) SetupTest() {
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
	err = viper.WriteConfig()
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
			args:         []string{"test"},
			stdoutString: "Profile test deleted successfully\n",
		},
		{
			name: "delete config with no args",
			args: []string{},
			err:  fmt.Errorf("missing argument <profileName> \n\nUsage:\n  delete <profileName> [flags]\n\nAliases:\n  delete, rm\n\nFlags:\n  -h, --help   help for delete\n"),
		},
		{
			name:         "delete config with non-existed profile",
			args:         []string{"test1"},
			stdoutString: "Profile test1 deleted successfully\n",
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
				assert.Equal("", viper.GetString("current_profile"))
				assert.Equal("", viper.GetString(tt.args[0]+".public_key"))
				assert.Equal("", viper.GetString(tt.args[0]+".private_key"))
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

	viper.Set("test.public_key", publicKey)
	viper.Set("test.private_key", privateKey)
	viper.Set("current_profile", newProfile)

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
			args:         []string{"test"},
			stdoutString: "Profile test deleted successfully\n",
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

			viper.Reset()
			viper.AddConfigPath(".")
			viper.SetConfigType("toml")
			viper.SetConfigName(".tidbcloud-cli")
			err = viper.ReadInConfig()
			assert.Nil(err)
			assert.Equal(newProfile, viper.GetString("current_profile"))
			assert.Equal("", viper.GetString(tt.args[0]+".public_key"))
			assert.Equal("", viper.GetString(tt.args[0]+".private_key"))
			assert.Equal("", viper.GetString(tt.args[0]))
		})
	}
}

func TestDeleteConfigSuite(t *testing.T) {
	suite.Run(t, new(DeleteConfigSuite))
}
