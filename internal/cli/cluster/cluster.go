package cluster

import (
	"tidbcloud-cli/internal/util"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle.Copy()
)

func ClusterCmd() *cobra.Command {
	var clusterCmd = &cobra.Command{
		Use:               "cluster",
		Short:             "Manage clusters for your project.",
		PersistentPreRunE: util.CheckAuth,
	}

	clusterCmd.AddCommand(CreateCmd())
	clusterCmd.AddCommand(DeleteCmd())
	clusterCmd.AddCommand(ListCmd())
	clusterCmd.AddCommand(DescribeCmd())
	return clusterCmd
}
