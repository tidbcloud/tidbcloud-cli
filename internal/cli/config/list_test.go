package config

import (
	"bytes"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ListConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *ListConfigSuite) SetupTest() {
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
}

func (suite *ListConfigSuite) TestListConfigArgs() {
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
			stdoutString: "Profile Name\ntest\t*\ntest1\n",
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
