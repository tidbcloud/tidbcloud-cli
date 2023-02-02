package telemetry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/internal/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EventSuite struct {
	suite.Suite
}

func (suite *EventSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	viper.Reset()
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName(".tidbcloud-cli")
	_ = viper.SafeWriteConfig()
}

func (suite *EventSuite) TearDownTest() {
	err := util.RemoveFile(".tidbcloud-cli.toml")
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *EventSuite) TestWithCommandPath() {
	testCmd := &cobra.Command{
		Use: "test",
	}
	rootCmd := &cobra.Command{
		Use: "root",
	}
	rootCmd.AddCommand(testCmd)

	e := newEvent(withCommandPath(testCmd))

	a := require.New(suite.T())
	a.Equal("root-test", e.Properties["command"])
}

func (suite *EventSuite) TestWithInteractive() {
	testCmd := &cobra.Command{
		Use: "test",
	}

	e := newEvent(WithInteractive(testCmd))
	a := require.New(suite.T())
	a.Equal(false, e.Properties["interactive"])

	testCmd = &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Annotations = map[string]string{"interactive": "true"}
		},
	}

	testCmd.SetArgs([]string{"t"})
	cmd, _ := testCmd.ExecuteContextC(NewTelemetryContext(context.Background()))
	e = newEvent(WithInteractive(cmd))
	a.Equal(true, e.Properties["interactive"])
}

func (suite *EventSuite) TestWithCommandPathAndAlias() {
	rootCmd := &cobra.Command{
		Use: "root",
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:     "test",
		Aliases: []string{"t"},
	})
	rootCmd.SetArgs([]string{"t"})
	calledCmd, _ := rootCmd.ExecuteContextC(NewTelemetryContext(context.Background()))

	e := newEvent(withCommandPath(calledCmd))

	a := require.New(suite.T())
	a.Equal("root-test", e.Properties["command"])
	a.Equal("t", e.Properties["alias"])
}

func (suite *EventSuite) TestWithDuration() {
	cmd := &cobra.Command{
		Use: "test-command",
		Run: func(cmd *cobra.Command, args []string) {
			time.Sleep(10 * time.Millisecond)
		},
	}

	cmd.SetArgs([]string{})
	_ = cmd.ExecuteContext(NewTelemetryContext(context.Background()))
	e := newEvent(withDuration(cmd))

	a := require.New(suite.T())
	a.GreaterOrEqual(e.Properties["duration"], int64(10))
}

func (suite *EventSuite) TestWithFlags() {
	cmd := &cobra.Command{
		Use: "test-command",
		Run: func(cmd *cobra.Command, args []string) {
			time.Sleep(10 * time.Millisecond)
		},
	}
	_ = cmd.Flags().Bool("test", false, "")
	_ = cmd.Flags().Bool("test2", false, "")
	_ = cmd.ParseFlags([]string{"--test"})

	e := newEvent(withFlags(cmd))

	a := require.New(suite.T())
	a.Equal(e.Properties["flags"], []string{"test"})
}

func (suite *EventSuite) TestWithVersion() {
	version.Version = "vTest"
	version.Commit = "sha-test"

	e := newEvent(withVersion())

	a := require.New(suite.T())
	a.Equal(e.Properties["version"], "vTest")
	a.Equal(e.Properties["git_commit"], "sha-test")
}

func (suite *EventSuite) TestWithOS() {
	e := newEvent(withOS())

	a := require.New(suite.T())
	a.Equal(e.Properties["os"], runtime.GOOS)
	a.Equal(e.Properties["arch"], runtime.GOARCH)
}

func (suite *EventSuite) TestWithAuthMethod_apiKey() {
	profileName := config.ActiveProfileName()
	viper.Set(fmt.Sprintf("%s.%s", profileName, prop.PublicKey), "test-public-key")
	viper.Set(fmt.Sprintf("%s.%s", profileName, prop.PrivateKey), "test-private-key")

	e := newEvent(withAuthMethod())

	a := require.New(suite.T())
	a.Equal(e.Properties["auth_method"], "api_key")
}

func (suite *EventSuite) TestWithProjectID_Flag() {
	cmd := &cobra.Command{
		Use: "test-command",
		Run: func(cmd *cobra.Command, args []string) {
			time.Sleep(10 * time.Millisecond)
		},
	}

	const projectID = "test"
	var p string
	cmd.Flags().StringVarP(&p, flag.ProjectID, "", "", "")
	_ = cmd.ParseFlags([]string{"--" + flag.ProjectID, projectID})

	e := newEvent(withProjectID(cmd))

	a := require.New(suite.T())
	a.Equal(projectID, e.Properties["project_id"])
}

func (suite *EventSuite) TestWithProjectID_NoFlag() {
	cmd := &cobra.Command{
		Use: "test-command",
		Run: func(cmd *cobra.Command, args []string) {
			time.Sleep(10 * time.Millisecond)
		},
	}

	e := newEvent(withProjectID(cmd))

	a := require.New(suite.T())
	_, ok := e.Properties["project_id"]
	a.Equal(false, ok)
}

func (suite *EventSuite) TestWithError() {
	e := newEvent(withError(errors.New("test")))

	a := require.New(suite.T())
	a.Equal("ERROR", e.Properties["result"])
	a.Equal("test", e.Properties["error"])
}

func (suite *EventSuite) TestWithHelpCommand() {
	testCmd := &cobra.Command{
		Use: "test",
	}
	rootCmd := &cobra.Command{
		Use: "root",
	}
	rootCmd.AddCommand(testCmd)
	rootCmd.InitDefaultHelpCmd()
	helpCmd := rootCmd.Commands()[0]

	args := []string{"test"}

	e := newEvent(withHelpCommand(helpCmd, args))

	assert.Equal(suite.T(), "root-test", e.Properties["help_command"])
}

func (suite *EventSuite) TestWithHelpCommand_NotFound() {
	testCmd := &cobra.Command{
		Use: "test",
	}
	rootCmd := &cobra.Command{
		Use: "root",
	}
	rootCmd.AddCommand(testCmd)
	rootCmd.InitDefaultHelpCmd()
	helpCmd := rootCmd.Commands()[0]

	args := []string{"test2"}

	e := newEvent(withHelpCommand(helpCmd, args))

	_, ok := e.Properties["help_command"]
	assert.False(suite.T(), ok)
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventSuite))
}
