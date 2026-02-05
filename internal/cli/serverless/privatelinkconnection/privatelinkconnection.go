// Copyright 2026 PingCAP, Inc.
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

package privatelinkconnection

import (
	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func Cmd(h *internal.Helper) *cobra.Command {
	var privateLinkConnectionCmd = &cobra.Command{
		Use:     "private-link-connection",
		Short:   "Manage TiDB Cloud Serverless private link connections",
		Aliases: []string{"plc"},
	}

	privateLinkConnectionCmd.AddCommand(CreateCmd(h))
	privateLinkConnectionCmd.AddCommand(GetCmd(h))
	privateLinkConnectionCmd.AddCommand(ListCmd(h))
	privateLinkConnectionCmd.AddCommand(DeleteCmd(h))
	privateLinkConnectionCmd.AddCommand(ZonesCmd(h))
	return privateLinkConnectionCmd
}
