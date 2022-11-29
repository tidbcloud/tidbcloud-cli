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
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func DescribeCmd(h *internal.Helper) *cobra.Command {
	var describeCmd = &cobra.Command{
		Use:     "describe",
		Short:   "Describe a cluster",
		Aliases: []string{"get"},
		Example: fmt.Sprintf(`  Get the cluster info in interactive mode:
  $ %[1]s cluster describe

  Get the cluster info in non-interactive mode:
  $ %[1]s cluster describe -p <project-id> -c <cluster-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// mark required flags in non-interactive mode
			if cmd.Flags().NFlag() != 0 {
				err := cmd.MarkFlagRequired(flag.ProjectID)
				if err != nil {
					return errors.Trace(err)
				}
				err = cmd.MarkFlagRequired(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d := h.Client()

			var projectID string
			var clusterID string
			if cmd.Flags().NFlag() == 0 {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, h.Client())
				if err != nil {
					return err
				}
				projectID = project.ID

				cluster, err := cloud.GetSelectedCluster(projectID, h.QueryPageSize, h.Client())
				if err != nil {
					return err
				}
				clusterID = cluster.ID
			} else {
				// non-interactive mode, get values from flags
				pID, err := cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return errors.Trace(err)
				}

				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				projectID = pID
				clusterID = cID
			}

			params := clusterApi.NewGetClusterParams().
				WithProjectID(projectID).
				WithClusterID(clusterID)
			cluster, err := d.GetCluster(params)
			if err != nil {
				return errors.Trace(err)
			}

			v, err := json.MarshalIndent(cluster.Payload, "", "  ")
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, string(v))
			return nil
		},
	}

	describeCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The project ID of the cluster to be deleted.")
	describeCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster to be deleted.")
	return describeCmd
}
