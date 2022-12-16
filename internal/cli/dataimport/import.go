package dataimport

import (
	"tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func ImportCmd(h *internal.Helper) *cobra.Command {
	var importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import data into TiDB Cloud cluster",
	}

	importCmd.AddCommand(ListCmd(h))
	importCmd.AddCommand(CancelCmd(h))
	importCmd.AddCommand(StartCmd(h))
	importCmd.AddCommand(DescribeCmd(h))
	return importCmd
}
