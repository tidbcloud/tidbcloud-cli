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
	"tidbcloud-cli/internal/cli/project"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func ListCmd(h *internal.Helper) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list <projectID>",
		Short: "List all clusters in a project",
		Example: fmt.Sprintf(`  List all clusters in the project(interactive mode):
  $ %[1]s cluster list

  List the clusters in the project(non-interactive mode):
  $ %[1]s cluster list <projectID> 

  List the clusters in the project with json format:
  $ %[1]s cluster list <projectID> -o json`, config.CliName),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var pID string
			if len(args) == 0 {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				_, projectItems, err := project.RetrieveProjects(h.QueryPageSize, h.Client())
				if err != nil {
					return err
				}
				set := hashset.New()
				for _, item := range projectItems {
					set.Add(item.ID)
				}
				model, err := ui.InitialSelectModel(set.Values(), "Choose the project ID:")
				if err != nil {
					return err
				}
				p := tea.NewProgram(model)
				projectModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if m, _ := projectModel.(ui.SelectModel); m.Interrupted {
					return nil
				}
				pID = projectModel.(ui.SelectModel).Choices[projectModel.(ui.SelectModel).Selected].(string)
			} else {
				pID = args[0]
			}

			total, items, err := retrieveClusters(pID, h.QueryPageSize, h.Client())
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
				res := &clusterApi.ListClustersOfProjectOKBody{
					Items: items,
					Total: &total,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"Name",
					"Status",
					"Version",
					"Region",
					"Type",
				}

				var rows []output.Row
				for _, item := range items {
					t := item.ClusterType
					// Currently serverless is called "DEVELOPER" in the API.
					// For better user experience, we change it to "SERVERLESS".
					// But we still keep the original value in the json result.
					if t == deverloperType {
						t = serverlessType
					}

					rows = append(rows, output.Row{
						*(item.ID),
						item.Name,
						item.Status.ClusterStatus,
						item.Status.TidbVersion,
						item.Region,
						t,
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

	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, "Output format. One of: human, json, default: human")
	return listCmd
}

func retrieveClusters(pID string, pageSize int64, d util.CloudClient) (int64, []*clusterApi.ListClustersOfProjectOKBodyItemsItems0, error) {
	params := clusterApi.NewListClustersOfProjectParams().WithProjectID(pID)
	var total int64 = math.MaxInt64
	var page int64 = 1
	var items []*clusterApi.ListClustersOfProjectOKBodyItemsItems0
	// loop to get all clusters
	for (page-1)*pageSize < total {
		clusters, err := d.ListClustersOfProject(params.WithPage(&page).WithPageSize(&pageSize))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}

		total = *clusters.Payload.Total
		page += 1
		items = append(items, clusters.Payload.Items...)
	}
	return total, items, nil
}
