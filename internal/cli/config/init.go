package config

import (
	"errors"
	"fmt"

	"tidbcloud-cli/internal/cli/ui"
	"tidbcloud-cli/internal/prop"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			color.Green("Check the https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management for more information about how to create API keys.")

			p := tea.NewProgram(initialDeletionInputModel())
			inputModel, err := p.StartReturningModel()
			if err != nil {
				return err
			}
			if inputModel.(ui.TextInputModel).Interrupted {
				return nil
			}

			inputs := inputModel.(ui.TextInputModel).Inputs
			pName := inputs[profileNameIdx].Value()
			pKey := inputs[publicKeyIdx].Value()
			sKey := inputs[privateKeyIdx].Value()

			profiles, err := GetAllProfiles()
			if err != nil {
				return err
			}
			if util.StringInSlice(profiles, pName) {
				return errors.New("profile already exists, use `config set` to update")
			}

			viper.Set(fmt.Sprintf("%s.%s", pName, prop.PublicKey), pKey)
			viper.Set(fmt.Sprintf("%s.%s", pName, prop.PrivateKey), sKey)
			viper.Set(prop.CurProfile, pName)
			err = viper.WriteConfig()
			if err != nil {
				return err
			}

			fgGreen := color.New(color.FgGreen).SprintFunc()
			hiGreen := color.New(color.FgHiGreen, color.BgWhite).SprintFunc()
			fmt.Printf("%s %s\n", fgGreen("Current profile has been changed to"), hiGreen(pName))
			return nil
		},
	}

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
