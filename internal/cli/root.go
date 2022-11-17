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
	"tidbcloud-cli/internal/cli/project"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cliName = "ticloud"
)

func Execute(ctx context.Context) {
	c := &config.Config{
		ActiveProfile: "",
	}

	h := &internal.Helper{
		Client: func() util.CloudClient {
			publicKey, privateKey := util.GetAccessKeys(c.ActiveProfile)
			return util.NewClientDelegate(publicKey, privateKey)
		},
		QueryPageSize: internal.DefaultPageSize,
		IOStreams:     iostream.System(),
		Config:        c,
	}

	var rootCmd = &cobra.Command{
		Use:   cliName,
		Short: "CLI tool to manage TiDB Cloud",
		Long:  fmt.Sprintf("%s is a CLI library for communicating with TiDB Cloud's API.", cliName),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			profile, err := cmd.Flags().GetString(flag.Profile)
			if err != nil {
				return err
			}
			if profile != "" {
				err := config.ValidateProfile(profile)
				if err != nil {
					return err
				}

				h.Config.ActiveProfile = profile
			} else {
				h.Config.ActiveProfile = viper.GetString(prop.CurProfile)
			}

			if shouldCheckAuth(cmd) {
				err := util.CheckAuth(h.Config.ActiveProfile)
				if err != nil {
					return err
				}
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(cluster.ClusterCmd(h))
	rootCmd.AddCommand(configCmd.ConfigCmd(h))
	rootCmd.AddCommand(project.ProjectCmd(h))

	rootCmd.PersistentFlags().StringP(flag.Profile, flag.ProfileShort, "", "Profile to use from your configuration file.")

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintf(h.IOStreams.Out, color.RedString("Error: %s\n", err.Error()))
		os.Exit(1)
	}
}

func shouldCheckAuth(cmd *cobra.Command) bool {
	cmdPrefixShouldSkip := []string{
		fmt.Sprintf("%s %s", cliName, "config"),
		fmt.Sprintf("%s %s", cliName, "help"),
		fmt.Sprintf("%s %s", cliName, "completion")}
	for _, p := range cmdPrefixShouldSkip {
		if strings.HasPrefix(cmd.CommandPath(), p) {
			return false
		}
	}

	return true
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".tidbcloud-cli" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("toml")
	viper.SetConfigName(".tidbcloud-cli")
	_ = viper.SafeWriteConfig()
	err = viper.ReadInConfig()
	if err != nil {
		color.Red("Failed to read config file: %s", err.Error())
		os.Exit(1)
	}
}
