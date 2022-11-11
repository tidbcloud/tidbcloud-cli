package config

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	exec "golang.org/x/sys/execabs"
)

const defaultEditor = "vi"

func EditCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "edit",
		Short: "Opens the config file with the default text editor.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := exec.Command(defaultEditor, viper.ConfigFileUsed())
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			return c.Run()
		},
	}

	return listCmd
}
