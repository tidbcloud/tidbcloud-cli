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

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SetCmd(h *internal.Helper) *cobra.Command {
	var setCmd = &cobra.Command{
		Use:   "set <propertyName> <value>",
		Short: "Configure specific properties of the active profile",
		Long: fmt.Sprintf(`Configure specific properties of the active profile.
Available properties : %v.

If use -P flag, the config in the specific profile will be set.
If not, the config in the active profile will be set`, prop.ProfileProperties()),
		Example: fmt.Sprintf(`  Set the value of the public-key in active profile:
  $ %[1]s config set public-key <public-key>

  Set the value of the public-key in the specific profile "test":
  $ %[1]s config set public-key <public-key> -P test`, config.CliName),
		Args: util.RequiredArgs("propertyName", "value"),
		RunE: func(cmd *cobra.Command, args []string) error {
			propertyName := args[0]
			value := args[1]

			var res string
			if util.StringInSlice(prop.ProfileProperties(), propertyName) {
				curP := h.Config.ActiveProfile
				if curP == "" {
					return fmt.Errorf("no profile is configured, please use `config create` to create a profile")
				}
				viper.Set(fmt.Sprintf("%s.%s", curP, propertyName), value)
				res = fmt.Sprintf("Set profile `%s` property `%s` to value `%s` successfully", curP, propertyName, value)
			} else {
				return fmt.Errorf("unrecognized property `%s`, use `config set --help` to find available properties", propertyName)
			}

			err := viper.WriteConfig()
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString(res))
			return nil
		},
	}

	return setCmd
}
