package config

import (
	"fmt"
	"sort"

	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "list all profiles",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := GetAllProfiles()
			if err != nil {
				return err
			}
			curP := viper.Get(prop.CurProfile)

			fmt.Println("Profile Name")
			for _, key := range profiles {
				if key == curP {
					color.Green(key + "\t*")
				} else {
					color.Green(key)
				}
			}
			return nil
		},
	}

	return listCmd
}

func GetAllProfiles() ([]string, error) {
	s := viper.AllSettings()
	keys := make([]string, 0, len(s))
	for k := range s {
		if !util.StringInSlice(prop.GlobalProperties(), k) {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)
	return keys, nil
}
