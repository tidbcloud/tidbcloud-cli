package config

import (
	"fmt"

	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UseCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "use <profileName>",
		Short: "Use the specified profile.",
		Args:  util.RequiredArgs("profileName"),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			err := setProfile(profileName)
			if err != nil {
				return err
			}
			return nil
		},
	}

	return listCmd
}

func setProfile(profileName string) error {
	profiles, err := GetAllProfiles()
	if err != nil {
		return err
	}

	if !util.StringInSlice(profiles, profileName) {
		return fmt.Errorf("profile %s not found", profileName)
	}

	viper.Set(prop.CurProfile, profileName)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	fgGreen := color.New(color.FgGreen).SprintFunc()
	hiGreen := color.New(color.FgHiGreen, color.BgWhite).SprintFunc()
	fmt.Printf("%s %s\n", fgGreen("Current profile has been changed to"), hiGreen(profileName))
	return nil
}
