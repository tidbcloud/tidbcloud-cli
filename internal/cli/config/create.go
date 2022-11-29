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
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle.Copy()
)

type createConfigField int

const (
	profileNameIdx createConfigField = iota
	publicKeyIdx
	privateKeyIdx
)

func CreateCmd(h *internal.Helper) *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Configure a profile to store access settings",
		Example: fmt.Sprintf(`  To configure the tool to work with TiDB Cloud in interactive mode:
  $ %[1]s config create

  To configure the tool to work with TiDB Cloud in non-interactive mode::
  $ %[1]s config create --profile-name <profile-name> --public-key <public-key> --private-key <private-key>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// mark required flags in non-interactive mode
			if cmd.Flags().NFlag() != 0 {
				err := cmd.MarkFlagRequired(flag.ProfileName)
				if err != nil {
					return errors.Trace(err)
				}
				err = cmd.MarkFlagRequired(flag.PublicKey)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.PrivateKey)
				if err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys."))

			var profileName string
			var publicKey string
			var privateKey string
			if cmd.Flags().NFlag() == 0 {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				p := tea.NewProgram(initialDeletionInputModel())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return nil
				}

				inputs := inputModel.(ui.TextInputModel).Inputs
				profileName = inputs[profileNameIdx].Value()
				publicKey = inputs[publicKeyIdx].Value()
				privateKey = inputs[privateKeyIdx].Value()
			} else {
				pName, err := cmd.Flags().GetString(flag.ProfileName)
				if err != nil {
					return errors.Trace(err)
				}
				profileName = pName

				pKey, err := cmd.Flags().GetString(flag.PublicKey)
				if err != nil {
					return errors.Trace(err)
				}
				publicKey = pKey

				priKey, err := cmd.Flags().GetString(flag.PrivateKey)
				if err != nil {
					return errors.Trace(err)
				}
				privateKey = priKey
			}

			profiles, err := config.GetAllProfiles()
			if err != nil {
				return errors.Trace(err)
			}
			if util.StringInSlice(profiles, profileName) {
				return fmt.Errorf("profile %s already exists, use `config set` to modify", profileName)
			}

			viper.Set(fmt.Sprintf("%s.%s", profileName, prop.PublicKey), publicKey)
			viper.Set(fmt.Sprintf("%s.%s", profileName, prop.PrivateKey), privateKey)
			viper.Set(prop.CurProfile, profileName)
			err = viper.WriteConfig()
			if err != nil {
				return errors.Trace(err)
			}

			fgGreen := color.New(color.FgGreen).SprintFunc()
			hiGreen := color.New(color.FgHiCyan).SprintFunc()
			fmt.Fprintf(h.IOStreams.Out, "%s %s\n", fgGreen("Current profile has been changed to"), hiGreen(profileName))
			return nil
		},
	}

	createCmd.Flags().String(flag.ProfileName, "", "the name of the profile to be created")
	createCmd.Flags().String(flag.PublicKey, "", "the public key of the TiDB Cloud API")
	createCmd.Flags().String(flag.PrivateKey, "", "the private key of the TiDB Cloud API")
	return createCmd
}

func initialDeletionInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 64
		f := createConfigField(i)

		switch f {
		case profileNameIdx:
			t.Placeholder = "Profile Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.Validate = func(s string) error {
				if len(s) == 0 {
					return errors.New("profile name is required")
				}
				return nil
			}
		case publicKeyIdx:
			t.Placeholder = "Public Key"
			t.CharLimit = 128
		case privateKeyIdx:
			t.Placeholder = "Private Key"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.CharLimit = 128
		}

		m.Inputs[i] = t
	}

	return m
}
