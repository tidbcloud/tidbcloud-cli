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
	"os"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"

	"github.com/juju/errors"
	"github.com/pelletier/go-toml"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeleteCmd(h *internal.Helper) *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete <profileName>",
		Short: "Delete a profile",
		Example: fmt.Sprintf(`  Delete the profile configuration:
  $ %[1]s config delete <profileName>`, config.CliName),
		Aliases: []string{"rm"},
		Args:    util.RequiredArgs("profileName"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Configuration needs to be deleted from toml, as viper doesn't support this yet.
			// FIXME :: change when https://github.com/spf13/viper/pull/519 is merged.
			profileName := args[0]

			settings := viper.AllSettings()
			t, err := toml.TreeFromMap(settings)
			if err != nil {
				return errors.Trace(err)
			}

			err = t.Delete(profileName)
			if err != nil {
				return errors.Trace(err)
			}

			// If the deleting profile is the current profile, set the current profile to another profile
			curP := t.Get(prop.CurProfile)
			if curP == profileName {
				profiles, err := config.GetAllProfiles()
				if err != nil {
					return errors.Trace(err)
				}

				newP := ""
				for _, profile := range profiles {
					if profile != profileName {
						newP = profile
						break
					}
				}
				if newP == "" {
					// If there is no other profile, unset current profile
					err = t.Delete(prop.CurProfile)
					if err != nil {
						return errors.Trace(err)
					}
				} else {
					t.Set(prop.CurProfile, newP)
				}
			}

			fs := afero.NewOsFs()
			file, err := fs.OpenFile(viper.ConfigFileUsed(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				return errors.Trace(err)
			}

			defer file.Close()

			s := t.String()
			_, err = file.WriteString(s)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Profile %s deleted successfully", profileName))
			return nil
		},
	}

	return deleteCmd
}
