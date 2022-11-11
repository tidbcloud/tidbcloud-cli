package project

import (
	"tidbcloud-cli/internal/util"

	"github.com/spf13/cobra"
)

func ProjectCmd() *cobra.Command {
	var projectCmd = &cobra.Command{
		Use:               "project",
		Short:             "Manage projects.",
		PersistentPreRunE: util.CheckAuth,
	}

	projectCmd.AddCommand(ListCmd())
	return projectCmd
}
