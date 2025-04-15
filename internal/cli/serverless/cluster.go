// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package serverless

import (
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/auditlog"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/authorizednetwork"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/branch"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/dataimport"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/export"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/sqluser"

	"github.com/spf13/cobra"
)

func Cmd(h *internal.Helper) *cobra.Command {
	var serverlessCmd = &cobra.Command{
		Use:     "serverless",
		Short:   "Manage TiDB Cloud Serverless clusters",
		Aliases: []string{"s"},
	}

	serverlessCmd.AddCommand(ListCmd(h))
	serverlessCmd.AddCommand(CreateCmd(h))
	serverlessCmd.AddCommand(DescribeCmd(h))
	serverlessCmd.AddCommand(UpdateCmd(h))
	serverlessCmd.AddCommand(DeleteCmd(h))
	serverlessCmd.AddCommand(branch.Cmd(h))
	serverlessCmd.AddCommand(ShellCmd(h))
	serverlessCmd.AddCommand(sqluser.SQLUserCmd(h))
	//serverlessCmd.AddCommand(backup.Cmd(h))
	//serverlessCmd.AddCommand(RestoreCmd(h))
	serverlessCmd.AddCommand(dataimport.ImportCmd(h))
	serverlessCmd.AddCommand(export.Cmd(h))
	serverlessCmd.AddCommand(SpendingLimitCmd(h))
	serverlessCmd.AddCommand(RegionCmd(h))
	serverlessCmd.AddCommand(auditlog.AuditLoggingCmd(h))
	serverlessCmd.AddCommand(CapacityCmd(h))
	serverlessCmd.AddCommand(authorizednetwork.AuthorizedNetworkCmd(h))

	return serverlessCmd
}
