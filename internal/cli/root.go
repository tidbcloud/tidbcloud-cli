// Copyright 2022 PingCAP, Inc.
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
	"strings"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/cli/cluster"
	configCmd "tidbcloud-cli/internal/cli/config"
	"tidbcloud-cli/internal/cli/dataimport"
	"tidbcloud-cli/internal/cli/project"
	"tidbcloud-cli/internal/cli/update"
	"tidbcloud-cli/internal/cli/version"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/log"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/service/github"
	"tidbcloud-cli/internal/telemetry"
	ver "tidbcloud-cli/internal/version"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute(ctx context.Context) {
	h := &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			publicKey, privateKey := config.GetPublicKey(), config.GetPrivateKey()
			apiUrl := config.GetApiUrl()
			// If the user has not set the api url, use the default one.
			if apiUrl == "" {
				apiUrl = cloud.DefaultApiUrl
			}
			delegate, err := cloud.NewClientDelegate(publicKey, privateKey, apiUrl)
			if err != nil {
				return nil, err
			}
			return delegate, nil
		},
		QueryPageSize: internal.DefaultPageSize,
		IOStreams:     iostream.System(),
	}

	rootCmd := RootCmd(h)
	initConfig()

	ctx = telemetry.NewTelemetryContext(ctx)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		telemetry.FinishTrackingCommand(telemetry.TrackOptions{Err: err})
		fmt.Fprintf(h.IOStreams.Out, color.RedString("Error: %s\n", err.Error()))
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
			} else {
				log.InitLogger("WARN")
			}
			profile, err := cmd.Flags().GetString(flag.Profile)
			if err != nil {
				return errors.Trace(err)
			}
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
					rel, _ := checkForUpdate(h.IOStreams.CanPrompt)
					updateMessageChan <- rel
				}()
			}

			var flagNoColor = cmd.Flags().Lookup(flag.NoColor)
			if flagNoColor != nil && flagNoColor.Changed {
				color.NoColor = true
			}

			if shouldCheckAuth(cmd) {
				err := config.CheckAuth()
				if err != nil {
					return errors.Trace(err)
				}
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if shouldCheckNewRelease(cmd) {
				newRelease := <-updateMessageChan
				if newRelease != nil {
					fmt.Fprintf(h.IOStreams.Out, fmt.Sprintf("\n%s %s → %s\n",
						color.YellowString("A new version of %s is available:", config.CliName),
						color.CyanString(ver.Version),
						color.CyanString(newRelease.Version)))

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

	rootCmd.AddCommand(cluster.ClusterCmd(h))
	rootCmd.AddCommand(configCmd.ConfigCmd(h))
	rootCmd.AddCommand(project.ProjectCmd(h))
	rootCmd.AddCommand(version.VersionCmd(h))
	rootCmd.AddCommand(update.UpdateCmd(h))
	rootCmd.AddCommand(dataimport.ImportCmd(h))

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

func shouldCheckAuth(cmd *cobra.Command) bool {
	cmdPrefixShouldSkip := []string{
		fmt.Sprintf("%s %s", config.CliName, "config"),
		fmt.Sprintf("%s %s", config.CliName, "help"),
		fmt.Sprintf("%s %s", config.CliName, "completion"),
		fmt.Sprintf("%s %s", config.CliName, "version"),
		fmt.Sprintf("%s %s", config.CliName, "update")}
	for _, p := range cmdPrefixShouldSkip {
		if strings.HasPrefix(cmd.CommandPath(), p) {
			return false
		}
	}

	return true
}

func checkForUpdate(isTerminal bool) (*github.ReleaseInfo, error) {
	if !isTerminal || ver.Version == ver.DevVersion {
		return nil, nil
	}

	return github.CheckForUpdate(config.Repo, true)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	path := home + "/" + config.HomePath
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
	_ = viper.SafeWriteConfig()

	// After version 0.1.2, we replace underscore with hyphen in properties.
	// In order to keep backward compatibility, we need to replace the old names to the new ones.
	data, _ := os.ReadFile(path + "/config.toml")
	newData := strings.Replace(string(data), "public_key", "public-key", -1)
	newData = strings.Replace(newData, "private_key", "private-key", -1)
	newData = strings.Replace(newData, "current_profile", "current-profile", -1)
	_ = os.WriteFile(path+"/config.toml", []byte(newData), 0600)

	err = viper.ReadInConfig()
	if err != nil {
		color.Red("Failed to read config file: %s", err.Error())
		os.Exit(1)
	}
}
