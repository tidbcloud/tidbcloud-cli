// Copyright 2025 PingCAP, Inc.
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

package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/ai"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/auth"
	configCmd "github.com/tidbcloud/tidbcloud-cli/internal/cli/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/project"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/upgrade"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/version"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/config/store"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/log"
	"github.com/tidbcloud/tidbcloud-cli/internal/prop"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/aws/s3"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/github"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	ver "github.com/tidbcloud/tidbcloud-cli/internal/version"

	"github.com/fatih/color"
	"github.com/juju/errors"
	logger "github.com/pingcap/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"go.uber.org/zap"
)

var ErrMissingCredentials = errors.New("this action requires authentication")

func Execute(ctx context.Context) {
	h := &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			publicKey, privateKey := config.GetPublicKey(), config.GetPrivateKey()
			apiUrl := config.GetApiUrl()
			// If the user has not set the api url, use the default one.
			if apiUrl == "" {
				apiUrl = cloud.DefaultApiUrl
			}
			serverlessEndpoint := config.GetServerlessEndpoint()
			if serverlessEndpoint == "" {
				serverlessEndpoint = cloud.DefaultServerlessEndpoint
			}
			iamEndpoint := config.GetIAMEndpoint()
			if iamEndpoint == "" {
				iamEndpoint = cloud.DefaultIAMEndpoint
			}

			var delegate cloud.TiDBCloudClient
			if publicKey != "" && privateKey != "" {
				var err error
				delegate, err = cloud.NewClientDelegateWithApiKey(publicKey, privateKey, apiUrl, serverlessEndpoint, iamEndpoint)
				if err != nil {
					return nil, err
				}
			} else {
				err := config.ValidateToken()
				if err != nil {
					logger.Debug("Failed to validate token", zap.Error(err))
					color.Yellow("\nTo log in using your TiDB Cloud username and password, run: \n %[1]s auth login\nTo set credentials using API keys, run: \n %[1]s config set public-key <public-key>\n %[1]s config set private-key <private-key>",
						config.CliName)
					return nil, ErrMissingCredentials
				}
				token, err := config.GetAccessToken()
				if err != nil {
					if errors.Is(err, keyring.ErrNotFound) || errors.Is(err, store.ErrNotSupported) {
						return nil, ErrMissingCredentials
					}
					return nil, err
				}
				delegate, err = cloud.NewClientDelegateWithToken(token, apiUrl, serverlessEndpoint, iamEndpoint)
				if err != nil {
					return nil, err
				}
			}

			return delegate, nil
		},
		Uploader: func(client cloud.TiDBCloudClient) s3.Uploader {
			return s3.NewUploader(client)
		},
		QueryPageSize: internal.DefaultPageSize,
		IOStreams:     iostream.System(),
	}

	rootCmd := RootCmd(h)
	initConfig()

	ctx = telemetry.NewTelemetryContext(ctx)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		telemetry.FinishTrackingCommand(telemetry.TrackOptions{Err: err})
		fmt.Fprint(h.IOStreams.Out, color.RedString("Error: %s\n", err.Error()))
		if errors.Is(err, util.InterruptError) {
			os.Exit(130)
		}
		os.Exit(1)
	}
}

func RootCmd(h *internal.Helper) *cobra.Command {
	updateMessageChan := make(chan *github.ReleaseInfo)
	var debugMode bool

	var rootCmd = &cobra.Command{
		Use:   config.CliName,
		Short: "CLI tool to manage TiDB Cloud",
		Long:  fmt.Sprintf("%s is a CLI library for communicating with TiDB Cloud's API.", config.CliName),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if debugMode {
				log.InitLogger("DEBUG")
				err := os.Setenv(config.DebugEnv, "1")
				if err != nil {
					return err
				}
			} else {
				log.InitLogger("WARN")
			}
			profile, err := cmd.Flags().GetString(flag.Profile)
			if err != nil {
				return errors.Trace(err)
			}
			profile = strings.ToLower(profile)
			if profile != "" {
				err := config.ValidateProfile(profile)
				if err != nil {
					return errors.Trace(err)
				}

				config.SetActiveProfile(profile)
			} else {
				config.SetActiveProfile(viper.GetString(prop.CurProfile))
			}

			telemetry.StartTrackingCommand(cmd, args)

			if shouldCheckNewRelease(cmd) {
				go func() {
					rel, err := checkForUpdate(cmd.Context(), h.IOStreams.CanPrompt)
					if err != nil {
						logger.Debug("Error checking for new release", zap.Error(err))
					}
					updateMessageChan <- rel
				}()
			}

			var flagNoColor = cmd.Flags().Lookup(flag.NoColor)
			if flagNoColor != nil && flagNoColor.Changed {
				color.NoColor = true
			}

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if shouldCheckNewRelease(cmd) {
				newRelease := <-updateMessageChan
				if newRelease != nil {
					fmt.Fprintf(h.IOStreams.Out, "\n%s %s → %s\n",
						color.YellowString("A new version of %s is available:", config.CliName),
						color.CyanString(ver.Version),
						color.CyanString(newRelease.Version))

					if config.IsUnderTiUP {
						fmt.Fprintln(h.IOStreams.Out, color.GreenString("Use `tiup update cloud` to update to the latest version"))
					} else {
						fmt.Fprintln(h.IOStreams.Out, color.GreenString("Use `ticloud update` to update to the latest version"))
					}
				}
			}

			telemetry.FinishTrackingCommand(telemetry.TrackOptions{})
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(auth.AuthCmd(h))
	rootCmd.AddCommand(configCmd.ConfigCmd(h))
	rootCmd.AddCommand(serverless.Cmd(h))
	rootCmd.AddCommand(ai.AICmd(h))
	rootCmd.AddCommand(project.ProjectCmd(h))
	rootCmd.AddCommand(version.VersionCmd(h))
	rootCmd.AddCommand(upgrade.Cmd(h))

	rootCmd.PersistentFlags().BoolVarP(&debugMode, flag.Debug, flag.DebugShort, false, "Enable debug mode")
	rootCmd.PersistentFlags().Bool(flag.NoColor, false, "Disable color output")
	rootCmd.PersistentFlags().StringP(flag.Profile, flag.ProfileShort, "", "Profile to use from your configuration file")
	return rootCmd
}

func shouldCheckNewRelease(cmd *cobra.Command) bool {
	cmdPrefixShouldSkip := []string{
		fmt.Sprintf("%s %s", config.CliName, "update")}
	for _, p := range cmdPrefixShouldSkip {
		if strings.HasPrefix(cmd.CommandPath(), p) {
			return false
		}
	}

	return true
}

func checkForUpdate(ctx context.Context, isTerminal bool) (*github.ReleaseInfo, error) {
	if !isTerminal || ver.Version == ver.DevVersion {
		logger.Debug("Skip checking for new release", zap.Bool("isTerminal", isTerminal), zap.String("version", ver.Version))
		return nil, nil
	}

	return github.CheckForUpdate(ctx, true)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	path := filepath.Join(home, config.HomePath)
	err = os.MkdirAll(path, 0700)
	if err != nil {
		color.Red("Failed to create home directory: %s", err)
		os.Exit(1)
	}

	// Search config in home directory with name ".tidbcloud-cli" (without extension).
	viper.AddConfigPath(path)
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.SetConfigPermissions(0600)
	err = viper.SafeWriteConfig()
	if err != nil {
		var existErr viper.ConfigFileAlreadyExistsError
		if !errors.As(err, &existErr) {
			color.Red("Failed to write config file: %s", err.Error())
			os.Exit(1)
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		color.Red("Failed to read config file: %s", err.Error())
		os.Exit(1)
	}
}
