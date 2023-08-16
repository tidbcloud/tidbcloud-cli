// Copyright 2022 PingCAP, Inc.
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

package cluster

import (
	"tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func ClusterCmd(h *internal.Helper) *cobra.Command {
	var clusterCmd = &cobra.Command{
		Use:   "cluster",
		Short: "Manage your clusters",
	}

	clusterCmd.AddCommand(CreateCmd(h))
	clusterCmd.AddCommand(DeleteCmd(h))
	clusterCmd.AddCommand(ListCmd(h))
	clusterCmd.AddCommand(DescribeCmd(h))
	clusterCmd.AddCommand(ConnectInfoCmd(h))
	clusterCmd.AddCommand(UpdateCmd(h))
	return clusterCmd
}
