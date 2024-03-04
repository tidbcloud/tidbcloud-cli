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

package export

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	exportApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/client/export_service"
)

type DescribeOpts struct {
	interactive bool
}

func (c DescribeOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ExportID,
	}
}

func (c *DescribeOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := c.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	// Mark required flags
	if !c.interactive {
		for _, fn := range flags {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DescribeCmd(h *internal.Helper) *cobra.Command {
	opts := DescribeOpts{
		interactive: true,
	}

	var describeCmd = &cobra.Command{
		Use:     "describe",
		Short:   "Describe an export",
		Aliases: []string{"get"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Get an export in interactive mode:
  $ %[1]s serverless export describe

  Get an export in non-interactive mode:
  $ %[1]s serverless export describe -c <cluster-id> -e <export-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}

			var exportID string
			var clusterID string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				export, err := cloud.GetSelectedExport(clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				exportID = export.ID
			} else {
				// non-interactive mode, get values from flags
				exportID, err = cmd.Flags().GetString(flag.ExportID)
				if err != nil {
					return errors.Trace(err)
				}

				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			params := exportApi.NewExportServiceGetExportParams().
				WithClusterID(clusterID).WithExportID(exportID)
			export, err := d.GetExport(params)
			if err != nil {
				return errors.Trace(err)
			}

			v, err := JSONMarshalWithoutEscape(export.Payload)
			if err != nil {
				return errors.Trace(err)
			}
			var dst = new(bytes.Buffer)
			err = json.Indent(dst, v, "", "  ")
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, dst.String())
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
	}

	describeCmd.Flags().StringP(flag.ExportID, flag.ExportIDShort, "", "The ID of the export to be described")
	describeCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID of the export to be described")
	describeCmd.MarkFlagsRequiredTogether(flag.ExportID, flag.ClusterID)
	return describeCmd
}

func JSONMarshalWithoutEscape(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
