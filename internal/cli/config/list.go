// Copyright 2026 PingCAP, Inc.
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

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/prop"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCmd(h *internal.Helper) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all profiles",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  List all available profiles:
  $ %[1]s config list`, config.CliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := config.GetAllProfiles()
			if err != nil {
				return errors.Trace(err)
			}
			curP := viper.Get(prop.CurProfile)

			fmt.Fprintln(h.IOStreams.Out, "Profile Name")
			for _, key := range profiles {
				if key == curP {
					fmt.Fprintln(h.IOStreams.Out, color.GreenString(key+"\t(active)"))
				} else {
					fmt.Fprintln(h.IOStreams.Out, color.GreenString(key))
				}
			}
			return nil
		},
	}

	return listCmd
}
