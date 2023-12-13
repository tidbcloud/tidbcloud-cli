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

package serverless

import (
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"

	serverlessModel "tidbcloud-cli/pkg/tidbcloud/serverless/models"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type ListOpts struct {
	interactive bool
}

func (c ListOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ProjectID,
	}
}

func ListCmd(h *internal.Helper) *cobra.Command {
	opts := ListOpts{
		interactive: true,
	}

	var listCmd = &cobra.Command{
		Use:         "list",
		Short:       "List all serverless clusters in a project",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  List all serverless clusters in interactive mode):
 $ %[1]s serverless list

 List the serverless clusters in non-interactive mode:
 $ %[1]s serverless list -p <project-id>

 List the serverless clusters in non-interactive mode:
 $ %[1]s serverless list -p <project-id> -o json`, config.CliName),
		Aliases: []string{"ls"},
		PreRun: func(cmd *cobra.Command, args []string) {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
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
				// non-interactive mode does not need projectID
				pID, err = cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			cmd.Annotations[telemetry.ProjectID] = pID

			total, items, err := cloud.RetrieveClusters(pID, h.QueryPageSize, d)
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
				res := &serverlessModel.TidbCloudOpenApiserverlessv1beta1ListClustersResponse{
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
	listCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The ID of the project")
	return listCmd
}
