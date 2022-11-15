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
	"errors"
	"fmt"

	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle.Copy()
)

type initConfigField int

const (
	profileNameIdx initConfigField = iota
	publicKeyIdx
	privateKeyIdx
)

func InitCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Configure a profile to store access settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() != 0 {
				err := cmd.MarkFlagRequired(flag.ProfileName)
				if err != nil {
					return err
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
			color.Green("Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.")

			var profileName string
			var publicKey string
			var privateKey string
			if cmd.Flags().NFlag() == 0 {
				p := tea.NewProgram(initialDeletionInputModel())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return err
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
					return err
				}
				profileName = pName

				pKey, err := cmd.Flags().GetString(flag.PublicKey)
				if err != nil {
					return err
				}
				publicKey = pKey

				priKey, err := cmd.Flags().GetString(flag.PrivateKey)
				if err != nil {
					return err
				}
				privateKey = priKey
			}

			profiles, err := GetAllProfiles()
			if err != nil {
				return err
			}
			if util.StringInSlice(profiles, profileName) {
				return errors.New("profile already exists, use `config set` to update")
			}

			viper.Set(fmt.Sprintf("%s.%s", profileName, prop.PublicKey), publicKey)
			viper.Set(fmt.Sprintf("%s.%s", profileName, prop.PrivateKey), privateKey)
			viper.Set(prop.CurProfile, profileName)
			err = viper.WriteConfig()
			if err != nil {
				return err
			}

			fgGreen := color.New(color.FgGreen).SprintFunc()
			hiGreen := color.New(color.FgHiGreen, color.BgWhite).SprintFunc()
			fmt.Printf("%s %s\n", fgGreen("Current profile has been changed to"), hiGreen(profileName))
			return nil
		},
	}

	initCmd.Flags().StringP(flag.ProfileName, flag.ProfileNameShort, "", "the name of the profile")
	initCmd.Flags().String(flag.PublicKey, "", "the public key of the TiDB Cloud API")
	initCmd.Flags().String(flag.PrivateKey, "", "the private key of the TiDB Cloud API")
	return initCmd
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
		f := initConfigField(i)

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
