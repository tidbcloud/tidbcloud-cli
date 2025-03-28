// Copyright 2025 PingCAP, Inc.
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
	"strconv"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"

	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func ListCmd(h *internal.Helper) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all accessible projects",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  List the projects:
  $ %[1]s project list

  List the projects with json format:
  $ %[1]s project list -o json`, config.CliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			_, items, err := cloud.RetrieveProjects(cmd.Context(), h.QueryPageSize, d)
			if err != nil {
				return err
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				res := iam.ApiListProjectsRsp{
					Projects: items,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"Name",
					"ClusterCount",
					"UserCount",
					"OrgID",
				}

				var rows []output.Row
				for _, item := range items {
					rows = append(rows, output.Row{
						*item.Id,
						*item.Name,
						strconv.FormatInt(int64(*item.ClusterCount), 10),
						strconv.FormatInt(int64(*item.UserCount), 10),
						*item.OrgId,
					})
				}

				err := output.PrintHumanTable(h.IOStreams.Out, columns, rows)
				// for human format, we print the table with brief information.
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

	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return listCmd
}
