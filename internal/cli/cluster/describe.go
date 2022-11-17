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
	"encoding/json"
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/ui"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func DescribeCmd(h *internal.Helper) *cobra.Command {
	var describeCmd = &cobra.Command{
		Use:     "describe",
		Short:   "Describe a cluster.",
		Aliases: []string{"get"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// mark required flags in non-interactive mode
			if cmd.Flags().NFlag() != 0 {
				err := cmd.MarkFlagRequired(flag.ProjectID)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.ClusterID)
				if err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d := h.Client()

			var projectID string
			var clusterID string
			if cmd.Flags().NFlag() == 0 {
				// interactive mode
				p := tea.NewProgram(initialClusterIdentifies())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return err
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return nil
				}

				projectID = inputModel.(ui.TextInputModel).Inputs[projectIDIdx].Value()
				clusterID = inputModel.(ui.TextInputModel).Inputs[clusterIDIdx].Value()
			} else {
				// non-interactive mode, get values from flags
				pID, err := cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return err
				}

				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return err
				}
				projectID = pID
				clusterID = cID
			}

			params := clusterApi.NewGetClusterParams().
				WithProjectID(projectID).
				WithClusterID(clusterID)
			cluster, err := d.GetCluster(params)
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(cluster.Payload, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(h.IOStreams.Out, string(v))
			return nil
		},
	}

	describeCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The project ID of the cluster to be deleted.")
	describeCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster to be deleted.")
	return describeCmd
}
