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

type DescribeConfigSuite struct {
	suite.Suite
	h *internal.Helper
}

func (suite *DescribeConfigSuite) SetupTest() {
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
			stdoutString: "{\n  \"private_key\": \"SDWIOUEOSDSDC\",\n  \"public_key\": \"SDIWODIJQNDKJQW\"\n}\n",
		},
		{
			name: "describe config with no args",
			args: []string{},
			err:  fmt.Errorf("missing argument <profileName> \n\nUsage:\n  describe <profileName> [flags]\n\nAliases:\n  describe, get\n\nFlags:\n  -h, --help   help for describe\n"),
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

func TestDescribeConfigSuite(t *testing.T) {
	suite.Run(t, new(DescribeConfigSuite))
}
