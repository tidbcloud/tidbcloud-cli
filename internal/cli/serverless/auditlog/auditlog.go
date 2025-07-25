// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auditlog

import (
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/auditlog/filter"
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
	auditLoggingCmd.AddCommand(filter.FilterRuleCmd(h))

	return auditLoggingCmd
}
