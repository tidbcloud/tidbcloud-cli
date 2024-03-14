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
	"strings"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/pingcap/log"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/juju/errors"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeleteCmd(h *internal.Helper) *cobra.Command {
	var force bool

	var deleteCmd = &cobra.Command{
		Use:   "delete <profile-name>",
		Short: "Delete a profile",
		Example: fmt.Sprintf(`  Delete the profile configuration:
  $ %[1]s config delete <profile-name>`, config.CliName),
		Aliases: []string{"rm"},
		Args:    util.RequiredArgs("profile-name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Configuration needs to be deleted from toml, as viper doesn't support this yet.
			// FIXME :: change when https://github.com/spf13/viper/pull/519 is merged.
			curProfileName := strings.ToLower(args[0])

			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to delete the profile")
				}

				confirmationMessage := fmt.Sprintf("%s %s %s", color.BlueString("Please type"), color.HiBlueString(config.Confirmed), color.BlueString("to confirm:"))
				prompt := &survey.Input{
					Message: confirmationMessage,
				}
				var userInput string
				err := survey.AskOne(prompt, &userInput)
				if err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					} else {
						return err
					}
				}

				if userInput != config.Confirmed {
					return errors.New("incorrect confirm string entered, skipping profile deletion")
				}
			}

			settings := viper.AllSettings()
			t, err := toml.TreeFromMap(settings)
			if err != nil {
				return errors.Trace(err)
			}

			err = t.Delete(curProfileName)
			if err != nil {
				return errors.Trace(err)
			}

			// If the deleting profile is the current profile, set the current profile to another profile
			curP := t.Get(prop.CurProfile)
			curPString, ok := curP.(string)
			if !ok {
				log.Debug("Failed to get current profile", zap.Any("current profile", curP))
				curPString = ""
			}
			if strings.EqualFold(curPString, curProfileName) {
				profiles, err := config.GetAllProfiles()
				if err != nil {
					return errors.Trace(err)
				}

				newP := ""
				for _, profile := range profiles {
					if !strings.EqualFold(profile, curProfileName) {
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
			file, err := fs.OpenFile(viper.ConfigFileUsed(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
			if err != nil {
				return errors.Trace(err)
			}

			defer file.Close()

			s := t.String()
			_, err = file.WriteString(s)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Profile %s deleted successfully", curProfileName))
			return nil
		},
	}

	deleteCmd.Flags().BoolVar(&force, flag.Force, false, "Delete a profile without confirmation")
	return deleteCmd
}
