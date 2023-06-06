// Copyright 2023 PingCAP, Inc.
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

package branch

import (
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/service/cloud"
	branchModel "tidbcloud-cli/pkg/tidbcloud/branch/models"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type ListOpts struct {
	interactive bool
}

func ListCmd(h *internal.Helper) *cobra.Command {
	opts := ListOpts{
		interactive: true,
	}

	var listCmd = &cobra.Command{
		Use:   "list <cluster-id>",
		Short: "List branches in the specified cluster",
		Example: fmt.Sprintf(`  List all branches in the cluster(interactive mode):
  $ %[1]s branch list

  List the branches in the cluster(non-interactive mode):
  $ %[1]s branch list <cluster-id> 

  List the branches in the cluster with json format:
  $ %[1]s branch list <cluster-id> -o json`, config.CliName),
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				opts.interactive = false
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				cluster, err := cloud.GetSelectedClusterWithoutProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID
			} else {
				clusterID = args[0]
			}

			total, items, err := cloud.RetrieveBranches(clusterID, h.QueryPageSize, d)
			if err != nil {
				return err
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			// for terminal which can prompt, humanFormat is the default format.
			// for other terminals, json format is the default format.
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				res := &branchModel.OpenapiListBranchesResp{
					Branches: items,
					Total:    &total,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"DisplayName",
					"ParentID",
					"State",
					"CreateTime",
					"UpdateTime",
				}

				var rows []output.Row
				for _, item := range items {
					rows = append(rows, output.Row{
						*item.ID,
						*item.DisplayName,
						*item.ParentID,
						string(*item.State),
						item.CreateTime.String(),
						item.UpdateTime.String(),
					})
				}

				err := output.PrintHumanTable(h.IOStreams.Out, columns, rows)
				if err != nil {
					return errors.Trace(err)
				}
				return nil
			} else {
				return fmt.Errorf("unsupported output format: %s", format)
			}
			return nil
		},
	}

	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, "Output format. One of: human, json. For the complete result, please use json format.")
	return listCmd
}
