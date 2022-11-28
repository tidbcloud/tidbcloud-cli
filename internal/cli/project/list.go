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

package project

import (
	"fmt"
	"math"
	"strconv"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/util"

	projectApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"
	"github.com/charmbracelet/bubbles/table"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func ListCmd(h *internal.Helper) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all accessible projects.",
		Aliases: []string{"ls"},
		Example: fmt.Sprintf(`  List the projects:
  $ %[1]s project list

  List the projects with json format:
  $ %[1]s project list -o json`, config.CliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			total, items, err := RetrieveProjects(h.QueryPageSize, h.Client())
			if err != nil {
				return err
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				res := projectApi.ListProjectsOKBody{
					Items: items,
					Total: &total,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []table.Column{
					{Title: "ID", Width: 20},
					{Title: "Name", Width: 20},
					{Title: "ClusterCount", Width: 15},
					{Title: "UserCount", Width: 10},
					{Title: "OrgID", Width: 20},
				}

				var rows []table.Row
				for _, item := range items {
					rows = append(rows, table.Row{
						item.ID,
						item.Name,
						strconv.FormatInt(item.ClusterCount, 10),
						strconv.FormatInt(item.UserCount, 10),
						item.OrgID,
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

	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, "Output format. One of: json|human, default: human")
	return listCmd
}

func RetrieveProjects(size int64, d util.CloudClient) (int64, []*projectApi.ListProjectsOKBodyItemsItems0, error) {
	params := projectApi.NewListProjectsParams()
	var total int64 = math.MaxInt64
	var page int64 = 1
	var pageSize = size
	var items []*projectApi.ListProjectsOKBodyItemsItems0
	for (page-1)*pageSize < total {
		projects, err := d.ListProjects(params.WithPage(&page).WithPageSize(&pageSize))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}

		total = *projects.Payload.Total
		page += 1
		items = append(items, projects.Payload.Items...)
	}
	return total, items, nil
}
