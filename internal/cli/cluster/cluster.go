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
	"tidbcloud-cli/internal/util"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	defaultPageSize = 100
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle.Copy()
)

func ClusterCmd() *cobra.Command {
	var clusterCmd = &cobra.Command{
		Use:               "cluster",
		Short:             "Manage clusters for your project.",
		PersistentPreRunE: util.CheckAuth(),
	}

	clusterCmd.AddCommand(CreateCmd())
	clusterCmd.AddCommand(DeleteCmd())
	clusterCmd.AddCommand(ListCmd())
	clusterCmd.AddCommand(DescribeCmd())
	return clusterCmd
}
