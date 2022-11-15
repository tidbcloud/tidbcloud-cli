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

	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/openapi"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/util"

	projectApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"
	"github.com/charmbracelet/bubbles/table"

	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all accessible projects.",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKey, privateKey := util.GetAccessKeys()
			apiClient := openapi.NewApiClient(publicKey, privateKey)
			params := projectApi.NewListProjectsParams()

			var total int64 = math.MaxInt64
			var page int64 = 1
			var pageSize int64 = defaultPageSize
			var items []*projectApi.ListProjectsOKBodyItemsItems0
			for (page-1)*pageSize < total {
				projects, err := apiClient.Project.ListProjects(params.WithPage(&page).WithPageSize(&pageSize))
				if err != nil {
					return err
				}

				total = *projects.Payload.Total
				page += 1
				items = append(items, projects.Payload.Items...)
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return err
			}
			if format == output.JsonFormat {
				err := output.PrintJson(items)
				if err != nil {
					return err
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
					return err
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
