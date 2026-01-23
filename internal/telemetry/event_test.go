// Copyright 2026 PingCAP, Inc.
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

package telemetry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/prop"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/internal/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	a.Equal(false, e.Properties[InteractiveMode])

	testCmd = &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Annotations = map[string]string{InteractiveMode: "true"}
		},
	}

	testCmd.SetArgs([]string{"t"})
	cmd, _ := testCmd.ExecuteContextC(NewTelemetryContext(context.Background()))
	e = newEvent(WithInteractive(cmd))
	a.Equal(true, e.Properties[InteractiveMode])
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

func (suite *EventSuite) TestWithProjectID() {
	const projectID = "test"
	cmd := &cobra.Command{
		Use:         "test-command",
		Annotations: make(map[string]string),
		Run: func(cmd *cobra.Command, args []string) {
			time.Sleep(10 * time.Millisecond)
			cmd.Annotations[ProjectID] = projectID
		},
	}
	cmd.SetArgs([]string{})
	_ = cmd.ExecuteContext(NewTelemetryContext(context.Background()))

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

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventSuite))
}
