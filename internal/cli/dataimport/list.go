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

package dataimport

import (
	"fmt"
	"strconv"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	importModel "tidbcloud-cli/pkg/tidbcloud/import/models"

	"github.com/dustin/go-humanize"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type ListOpts struct {
	interactive bool
}

func (c ListOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
	}
}

func ListCmd(h *internal.Helper) *cobra.Command {
	opts := ListOpts{
		interactive: true,
	}

	var listCmd = &cobra.Command{
		Use:         "list",
		Short:       "List data import tasks",
		Aliases:     []string{"ls"},
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  List import tasks in interactive mode:
  $ %[1]s import list

  List import tasks in non-interactive mode:
  $ %[1]s import list --project-id <project-id> --cluster-id <cluster-id>
  
  List the clusters in the project with json format:
  $ %[1]s import list --project-id <project-id> --cluster-id <cluster-id> --output json`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// mark required flags in non-interactive mode
			if !opts.interactive {
				for _, fn := range flags {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectID, clusterID string
			d, err := h.Client()
			if err != nil {
				return err
			}

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
				projectID = project.ID

				cluster, err := cloud.GetSelectedCluster(projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID
			} else {
				// non-interactive mode
				projectID = cmd.Flag(flag.ProjectID).Value.String()
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
			}

			cmd.Annotations[telemetry.ProjectID] = projectID

			total, importTasks, err := cloud.RetrieveImports(projectID, clusterID, h.QueryPageSize, d)
			if err != nil {
				return err
			}
			totalStr := strconv.FormatUint(total, 10)

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			// for terminal which can prompt, humanFormat is the default format.
			// for other terminals, json format is the default format.
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				res := &importModel.OpenapiListImportsResp{
					Imports: importTasks,
					Total:   &totalStr,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"Type",
					"Status",
					"CreatedAt",
					"Source",
					"DataFormat",
					"Size",
				}

				var rows []output.Row
				for _, item := range importTasks {
					var source string
					if item.CreationDetails.Type != nil && *item.CreationDetails.Type == importModel.CreateImportReqImportTypeS3 {
						source = item.CreationDetails.SourceURL
					} else {
						source = item.CreationDetails.FileName
					}

					rows = append(rows, output.Row{
						item.ID,
						string(*item.CreationDetails.Type),
						string(*item.Status),
						item.CreatedAt.String(),
						source,
						string(*item.DataFormat),
						convertToStoreSize(*item.TotalSize),
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

	listCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "Project ID")
	listCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, "Output format. One of[\"human\", \"json\"]. For the complete result, please use json format.")
	return listCmd
}

func convertToStoreSize(totalSize string) string {
	size, err := strconv.ParseUint(totalSize, 10, 64)
	if err != nil {
		return "NaN"
	}
	return humanize.IBytes(size)
}
