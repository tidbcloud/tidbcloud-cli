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

	"tidbcloud-cli/internal/cli/cluster"
	"tidbcloud-cli/internal/cli/config"
	"tidbcloud-cli/internal/cli/project"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cliName = "tidbcloud-cli"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          cliName,
	Short:        "CLI tool to manage TiDB Cloud",
	Long:         fmt.Sprintf("%s is a CLI library for communicating with TiDB Cloud's API.", cliName),
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) {
	rootCmd.AddCommand(cluster.ClusterCmd())
	rootCmd.AddCommand(config.ConfigCmd())
	rootCmd.AddCommand(project.ProjectCmd())

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
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
