// Copyright 2024 PingCAP, Inc.
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
	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/cli/serverless/branch"
	"tidbcloud-cli/internal/cli/serverless/dataimport"
	"tidbcloud-cli/internal/cli/serverless/export"
	"tidbcloud-cli/internal/cli/serverless/sqluser"

	"github.com/spf13/cobra"
)

func Cmd(h *internal.Helper) *cobra.Command {
	var serverlessCmd = &cobra.Command{
		Use:     "serverless",
		Short:   "Manage TiDB Serverless clusters",
		Aliases: []string{"s"},
	}

	serverlessCmd.AddCommand(ListCmd(h))
	serverlessCmd.AddCommand(CreateCmd(h))
	serverlessCmd.AddCommand(DescribeCmd(h))
	serverlessCmd.AddCommand(UpdateCmd(h))
	serverlessCmd.AddCommand(DeleteCmd(h))
	serverlessCmd.AddCommand(branch.Cmd(h))
	serverlessCmd.AddCommand(ShellCmd(h))
	//serverlessCmd.AddCommand(backup.Cmd(h))
	//serverlessCmd.AddCommand(RestoreCmd(h))
	serverlessCmd.AddCommand(dataimport.ImportCmd(h))
	serverlessCmd.AddCommand(export.Cmd(h))
	serverlessCmd.AddCommand(SpendingLimitCmd(h))
	serverlessCmd.AddCommand(RegionCmd(h))

	serverlessCmd.AddCommand(sqluser.SQLUserCmd(h))
	return serverlessCmd
}
