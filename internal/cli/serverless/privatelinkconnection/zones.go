// Copyright 2026 PingCAP, Inc.
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

package privatelinkconnection

import (
	"fmt"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type ZonesOpts struct {
	interactive bool
}

func (c ZonesOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func (c *ZonesOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := c.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	if !c.interactive {
		for _, fn := range flags {
			if err := cmd.MarkFlagRequired(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func ZonesCmd(h *internal.Helper) *cobra.Command {
	opts := ZonesOpts{
		interactive: true,
	}

	var zonesCmd = &cobra.Command{
		Use:   "zones",
		Short: "Get private link connection availability zones",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Get private link connection availability zones in interactive mode:
  $ %[1]s serverless private-link-connection zones

  Get private link connection availability zones in non-interactive mode:
  $ %[1]s serverless private-link-connection zones -c <cluster-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.MarkInteractive(cmd); err != nil {
				return errors.Trace(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}
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
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			res, err := d.GetAvailabilityZones(ctx, clusterID)
			if err != nil {
				return errors.Trace(err)
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				return errors.Trace(output.PrintJson(h.IOStreams.Out, res))
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"AccountID",
					"AvailabilityZoneIDs",
				}
				accountID := ""
				if res.AccountId != nil {
					accountID = *res.AccountId
				}
				rows := []output.Row{
					{
						accountID,
						strings.Join(res.AzIds, ","),
					},
				}
				return errors.Trace(output.PrintHumanTable(h.IOStreams.Out, columns, rows))
			}
			return fmt.Errorf("unsupported output format: %s", strings.ToLower(format))
		},
	}

	zonesCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	zonesCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return zonesCmd
}
