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

package config

import (
	"fmt"
	"io"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UseCmd(h *internal.Helper) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "use <profile-name>",
		Short: "Use the specified profile as the active profile",
		Example: fmt.Sprintf(`  Use the "test" profile as the active profile:
  $ %[1]s config use test`, config.CliName),
		Args: util.RequiredArgs("profile-name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			err := SetProfile(h.IOStreams.Out, profileName)
			if err != nil {
				return err
			}
			return nil
		},
	}

	return listCmd
}

// SetProfile sets the specified profile as the active profile if profile exist.
// If not, return error.
func SetProfile(out io.Writer, profileName string) error {
	err := config.ValidateProfile(profileName)
	if err != nil {
		return err
	}

	viper.Set(prop.CurProfile, profileName)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	fgGreen := color.New(color.FgGreen).SprintFunc()
	hiGreen := color.New(color.FgHiCyan).SprintFunc()
	fmt.Fprintf(out, "%s %s\n", fgGreen("Current profile has been changed to"), hiGreen(profileName))
	return nil
}
