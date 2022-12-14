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
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/service/github"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute(ctx context.Context, ver, commit, buildDate string) {
	c := &config.Config{
		ActiveProfile: "",
	}

	var binpath string
	if exepath, err := os.Executable(); err == nil {
		binpath = exepath
	}

	isUnderTiUP := util.IsUnderTiUP(binpath)
	config.SetCliName(isUnderTiUP)

	h := &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			publicKey, privateKey := util.GetAccessKeys(c.ActiveProfile)
			apiUrl := util.GetApiUrl(c.ActiveProfile)
			// If the user has not set the api url, use the default one.
			if apiUrl == "" {
				apiUrl = cloud.DefaultApiUrl
			}
			delegate, err := cloud.NewClientDelegate(publicKey, privateKey, apiUrl, ver)
			if err != nil {
				return nil, err
			}
			return delegate, nil
		},
		QueryPageSize: internal.DefaultPageSize,
		IOStreams:     iostream.System(),
		Config:        c,
		IsUnderTiUP:   isUnderTiUP,
	}

	rootCmd := RootCmd(h, ver, commit, buildDate)
	initConfig()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintf(h.IOStreams.Out, color.RedString("Error: %s\n", err.Error()))
		os.Exit(1)
	}
}

func RootCmd(h *internal.Helper, ver, commit, buildDate string) *cobra.Command {
	updateMessageChan := make(chan *github.ReleaseInfo)
	var rootCmd = &cobra.Command{
		Use:   config.CliName,
		Short: "CLI tool to manage TiDB Cloud",
		Long:  fmt.Sprintf("%s is a CLI library for communicating with TiDB Cloud's API.", config.CliName),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if shouldCheckNewRelease(cmd) {
				go func() {
					rel, _ := checkForUpdate(ver, h.IOStreams.CanPrompt)
					updateMessageChan <- rel
				}()
			}

			var flagNoColor = cmd.Flags().Lookup(flag.NoColor)
			if flagNoColor != nil && flagNoColor.Changed {
				color.NoColor = true
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

				h.Config.ActiveProfile = profile
			} else {
				h.Config.ActiveProfile = viper.GetString(prop.CurProfile)
			}

			if shouldCheckAuth(cmd) {
				err := util.CheckAuth(h.Config.ActiveProfile)
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
					fmt.Fprintf(h.IOStreams.Out, fmt.Sprintf("\n%s %s ??? %s\n",
						color.YellowString("A new version of %s is available:", config.CliName),
						color.CyanString(ver),
						color.CyanString(newRelease.Version)))

					if h.IsUnderTiUP {
						fmt.Fprintln(h.IOStreams.Out, color.GreenString("Use `tiup update cloud` to update to the latest version"))
					} else {
						fmt.Fprintln(h.IOStreams.Out, color.GreenString("Use `ticloud update` to update to the latest version"))
					}
				}
			}
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(cluster.ClusterCmd(h))
	rootCmd.AddCommand(configCmd.ConfigCmd(h))
	rootCmd.AddCommand(project.ProjectCmd(h))
	rootCmd.AddCommand(version.VersionCmd(h, ver, commit, buildDate))
	rootCmd.AddCommand(update.UpdateCmd(h, ver))
	rootCmd.AddCommand(dataimport.ImportCmd(h))

	rootCmd.PersistentFlags().Bool(flag.NoColor, false, "Disable color output")
	rootCmd.PersistentFlags().StringP(flag.Profile, flag.ProfileShort, "", "Profile to use from your configuration file.")
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

func checkForUpdate(currentVersion string, isTerminal bool) (*github.ReleaseInfo, error) {
	if !isTerminal || currentVersion == config.DevVersion {
		return nil, nil
	}

	return github.CheckForUpdate(config.Repo, currentVersion, true)
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
