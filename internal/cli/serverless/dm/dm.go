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

package dm

import (
	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func DMCmd(h *internal.Helper) *cobra.Command {
	var dmCmd = &cobra.Command{
		Use:   "dm",
		Short: "Manage TiDB Cloud Serverless data migration",
	}

	dmCmd.AddCommand(ListCmd(h))
	dmCmd.AddCommand(DescribeCmd(h))
	dmCmd.AddCommand(DeleteCmd(h))
	dmCmd.AddCommand(CreateCmd(h))
	dmCmd.AddCommand(PrecheckCmd(h))
	dmCmd.AddCommand(GetPrecheckCmd(h))
	dmCmd.AddCommand(CancelPrecheckCmd(h))
	dmCmd.AddCommand(OperateCmd(h))
	return dmCmd
}
