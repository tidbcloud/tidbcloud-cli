package auditlog

import (
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
)

func AuditLoggingCmd(h *internal.Helper) *cobra.Command {
	var auditLoggingCmd = &cobra.Command{
		Use:     "audit-log",
		Short:   "Manage TiDB Cloud Serverless database audit logging",
		Aliases: []string{"al"},
	}

	auditLoggingCmd.AddCommand(DownloadCmd(h))
	auditLoggingCmd.AddCommand(DescribeCmd(h))
	auditLoggingCmd.AddCommand(ConfigCmd(h))
	auditLoggingCmd.AddCommand(EnableCmd(h))
	auditLoggingCmd.AddCommand(DisableCmd(h))

	return auditLoggingCmd
}
