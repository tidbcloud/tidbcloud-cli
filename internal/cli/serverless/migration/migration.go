package migration

import (
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
)

func MigrationCmd(h *internal.Helper) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "migration",
		Short:   "Manage TiDB Cloud Serverless migrations",
		Aliases: []string{"dm"},
	}

	cmd.AddCommand(CreateCmd(h))
	cmd.AddCommand(DescribeCmd(h))
	cmd.AddCommand(ListCmd(h))
	cmd.AddCommand(DeleteCmd(h))
	cmd.AddCommand(TemplateCmd(h))
	cmd.AddCommand(PauseCmd(h))
	cmd.AddCommand(ResumeCmd(h))

	return cmd
}
