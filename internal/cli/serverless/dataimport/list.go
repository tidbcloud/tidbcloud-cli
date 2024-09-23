// Copyright 2024 PingCAP, Inc.
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
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"

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
		Short:       "List data import tasks",
		Aliases:     []string{"ls"},
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  List import tasks in interactive mode:
  $ %[1]s serverless import list

  List import tasks in non-interactive mode:
  $ %[1]s serverless import list --cluster-id <cluster-id>
  
  List the clusters in the project with json format:
  $ %[1]s serverless import list --cluster-id <cluster-id> --output json`,
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

			total, importTasks, err := cloud.RetrieveImports(cmd.Context(), clusterID, h.QueryPageSize, d)
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
				res := &imp.ListImportsResp{
					Imports:   importTasks,
					TotalSize: &total,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"SourceType",
					"State",
					"CreateTime",
					"Source",
					"FileType",
					"Size",
				}

				var rows []output.Row
				for _, item := range importTasks {
					var source, st, ft string
					if item.CreationDetails != nil && item.CreationDetails.Source != nil {
						st = string(item.CreationDetails.Source.Type)
						if item.CreationDetails.Source.Type == imp.IMPORTSOURCETYPEENUM_LOCAL {
							source = *item.CreationDetails.Source.Local.FileName
						} else if item.CreationDetails.Source.Type == imp.IMPORTSOURCETYPEENUM_S3 {
							source = item.CreationDetails.Source.S3.Uri
						} else if item.CreationDetails.Source.Type == imp.IMPORTSOURCETYPEENUM_GCS {
							source = item.CreationDetails.Source.Gcs.Uri
						} else if item.CreationDetails.Source.Type == imp.IMPORTSOURCETYPEENUM_AZURE_BLOB {
							source = item.CreationDetails.Source.AzureBlob.Uri
						}
					}
					if item.CreationDetails != nil && item.CreationDetails.ImportOptions != nil {
						ft = string(item.CreationDetails.ImportOptions.FileType)
					}
					rows = append(rows, output.Row{
						*item.ImportId,
						st,
						string(*item.State),
						item.CreateTime.Format(time.RFC3339),
						source,
						ft,
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

	listCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return listCmd
}

func convertToStoreSize(totalSize string) string {
	size, err := strconv.ParseUint(totalSize, 10, 64)
	if err != nil {
		return "NaN"
	}
	return humanize.IBytes(size)
}
