package config

import (
	"fmt"

	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SetCmd() *cobra.Command {
	var setCmd = &cobra.Command{
		Use:   "set <propertyName> <value>",
		Short: "Configure specific properties of the active profile or global.",
		Long: fmt.Sprintf(`Configure specific properties of the active profile or global.
Available profile properties : %v. 
Available global properties : %v`, prop.ProfileProperties(), prop.GlobalProperties()),
		Args: util.RequiredArgs("propertyName", "value"),
		RunE: func(cmd *cobra.Command, args []string) error {
			propertyName := args[0]
			value := args[1]

			if util.StringInSlice(prop.GlobalProperties(), propertyName) {
				if propertyName == prop.CurProfile {
					setProfile(value)
				} else {
					viper.Set(propertyName, value)
				}

			} else if util.StringInSlice(prop.ProfileProperties(), propertyName) {
				curP := viper.Get(prop.CurProfile)
				if curP == nil {
					return fmt.Errorf("no profile is configured, please use `config init` to create a profile")
				}
				viper.Set(fmt.Sprintf("%s.%s", curP, propertyName), value)
			} else {
				return fmt.Errorf("unrecognized property %s ", propertyName)
			}

			err := viper.WriteConfig()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return setCmd
}
