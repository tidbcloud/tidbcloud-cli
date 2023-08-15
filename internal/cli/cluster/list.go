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
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	serverlessModel "tidbcloud-cli/pkg/tidbcloud/serverless/models"
)

type ListOpts struct {
	interactive bool
}

func ListCmd(h *internal.Helper) *cobra.Command {
	opts := ListOpts{
		interactive: true,
	}

	var listCmd = &cobra.Command{
		Use:         "list <project-id>",
		Short:       "List all clusters in a project",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  List all clusters in the project(interactive mode):
 $ %[1]s cluster list

 List the clusters in the project(non-interactive mode):
 $ %[1]s cluster list <project-id>

 List the clusters in the project with json format:
 $ %[1]s cluster list <project-id> -o json`, config.CliName),
		Aliases: []string{"ls"},
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

			var pID string
			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				pID = project.ID
			} else {
				pID = args[0]
			}

			cmd.Annotations[telemetry.ProjectID] = pID

			total, items, err := cloud.RetrieveClusters(pID, int32(h.QueryPageSize), d)
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
				res := &serverlessModel.V1ListClustersResponse{
					Clusters:  items,
					TotalSize: total,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"DisplayName",
					"State",
					"Version",
					"Cloud",
					"Region",
					"Type",
				}

				var rows []output.Row
				for _, item := range items {
					rows = append(rows, output.Row{
						item.ClusterID,
						*item.DisplayName,
						string(*item.State),
						item.Version,
						string(*item.Region.Provider),
						item.Region.DisplayName,
						serverlessType,
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

	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, "Output format, One of [\"human\" \"json\"], for the complete result, please use json format")
	return listCmd
}
