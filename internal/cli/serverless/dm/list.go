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

package dm

import (
	"fmt"
	"strconv"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"

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
	}
}

func ListCmd(h *internal.Helper) *cobra.Command {
	opts := ListOpts{
		interactive: true,
	}

	var listCmd = &cobra.Command{
		Use:         "list",
		Short:       "List DM tasks",
		Aliases:     []string{"ls"},
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  List DM tasks in interactive mode:
  $ %[1]s serverless dm list

  List DM tasks in non-interactive mode:
  $ %[1]s serverless dm list --cluster-id <cluster-id>
  
  List the DM tasks in the cluster with json format:
  $ %[1]s serverless dm list --cluster-id <cluster-id> --output json`,
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
			var clusterID string
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"

				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}

				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID
			} else {
				// non-interactive mode
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
			}

			cmd.Annotations[telemetry.ClusterID] = clusterID

			tasks, err := d.ListTasks(ctx, clusterID, nil, nil, nil)
			if err != nil {
				return errors.Trace(err)
			}

			// get the format for output
			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			// print the result
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				err := output.PrintJson(h.IOStreams.Out, tasks)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"Name",
					"Mode",
					"SubTasks",
					"User",
					"CreateTime",
				}

				var rows []output.Row
				if tasks.Tasks != nil {
					for _, task := range tasks.Tasks {
						createTimeStr := ""
						if task.CreateTime != nil {
							createTimeStr = humanize.Time(*task.CreateTime)
						}

						subTaskCount := ""
						if task.SubTasks != nil {
							subTaskCount = strconv.Itoa(len(task.SubTasks))
						}

						modeStr := ""
						if task.Mode != nil {
							modeStr = string(*task.Mode)
						}

						idStr := ""
						if task.Id != nil {
							idStr = *task.Id
						}

						nameStr := ""
						if task.Name != nil {
							nameStr = *task.Name
						}

						userStr := ""
						if task.ClusterUser != nil {
							userStr = *task.ClusterUser
						}

						rows = append(rows, output.Row{
							idStr,
							nameStr,
							modeStr,
							subTaskCount,
							userStr,
							createTimeStr,
						})
					}
				}

				err := output.PrintHumanTable(h.IOStreams.Out, columns, rows)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				return fmt.Errorf("unsupported output format: %s", format)
			}

			return nil
		},
	}

	listCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return listCmd
}
