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
	"math"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/charmbracelet/bubbles/table"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func ListCmd(h *internal.Helper) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list <projectID>",
		Short: "List all clusters in a project.",
		Args:  util.RequiredArgs("projectID"),
		Example: fmt.Sprintf(`  List the clusters in the project:
  $ %[1]s cluster list <projectID> 

  List the clusters in the project with json format:
  $ %[1]s cluster list <projectID> -o json`, config.CliName),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			d := h.Client()
			pID := args[0]

			params := clusterApi.NewListClustersOfProjectParams().WithProjectID(pID)
			var total int64 = math.MaxInt64
			var page int64 = 1
			var pageSize = h.QueryPageSize
			var items []*clusterApi.ListClustersOfProjectOKBodyItemsItems0
			// loop to get all clusters
			for (page-1)*pageSize < total {
				clusters, err := d.ListClustersOfProject(params.WithPage(&page).WithPageSize(&pageSize))
				if err != nil {
					return errors.Trace(err)
				}

				total = *clusters.Payload.Total
				page += 1
				items = append(items, clusters.Payload.Items...)
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			// for terminal which can prompt, humanFormat is the default format.
			// for other terminals, json format is the default format.
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				res := &clusterApi.ListClustersOfProjectOKBody{
					Items: items,
					Total: &total,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				// for human format, we print the table with brief information.
				color.New(color.BgYellow).Fprintf(h.IOStreams.Out, "  For detailed information, please output with json format.")
				columns := []table.Column{
					{Title: "ID", Width: 20},
					{Title: "Name", Width: 20},
					{Title: "Status", Width: 10},
					{Title: "Version", Width: 10},
					{Title: "Region", Width: 10},
					{Title: "Type", Width: 10},
				}

				var rows []table.Row
				for _, item := range items {
					rows = append(rows, table.Row{
						*(item.ID),
						item.Name,
						item.Status.ClusterStatus,
						item.Status.TidbVersion,
						item.Region,
						item.ClusterType,
					})
				}

				err := output.PrintHumanTable(columns, rows)
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

	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, "Output format. One of: human, json, default: human")
	return listCmd
}
