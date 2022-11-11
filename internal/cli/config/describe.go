package config

import (
	"encoding/json"
	"fmt"

	"tidbcloud-cli/internal/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DescribeCmd() *cobra.Command {
	describeCmd := &cobra.Command{
		Use:     "describe <profileName>",
		Aliases: []string{"get"},
		Short:   "Return a specific profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			profiles, err := GetAllProfiles()
			if err != nil {
				return err
			}

			if !util.StringInSlice(profiles, name) {
				return fmt.Errorf("profile %s not found", name)
			}

			value := viper.Get(name)
			v, err := json.MarshalIndent(value, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}

	return describeCmd
}
