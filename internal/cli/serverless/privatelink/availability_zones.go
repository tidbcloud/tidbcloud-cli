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

package privatelink

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

type AvailabilityZonesOpts struct {
	interactive bool
}

func (c AvailabilityZonesOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func (c *AvailabilityZonesOpts) MarkInteractive(cmd *cobra.Command) error {
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

func AvailabilityZonesCmd(h *internal.Helper) *cobra.Command {
	opts := AvailabilityZonesOpts{
		interactive: true,
	}

	var availabilityZonesCmd = &cobra.Command{
		Use:   "availability-zones",
		Short: "Get account and availability zones information",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Get availability zones in interactive mode:
  $ %[1]s serverless private-link-connection availability-zones

  Get availability zones in non-interactive mode:
  $ %[1]s serverless private-link-connection availability-zones -c <cluster-id>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
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

			resp, err := d.GetPrivateLinkConnectionAvailabilityZones(ctx, clusterID)
			if err != nil {
				return errors.Trace(err)
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				err = output.PrintJson(h.IOStreams.Out, resp)
				if err != nil {
					return errors.Trace(err)
				}
				return nil
			}
			if format != output.HumanFormat {
				return fmt.Errorf("unsupported output format: %s", format)
			}

			columns := []output.Column{"AccountID", "AvailabilityZones"}
			availabilityZones := ""
			if resp.AzIds != nil {
				availabilityZones = strings.Join(resp.AzIds, ", ")
			}
			rows := []output.Row{{valueOrEmptyString(resp.AccountId), availabilityZones}}
			err = output.PrintHumanTable(h.IOStreams.Out, columns, rows)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
	}

	availabilityZonesCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	availabilityZonesCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return availabilityZonesCmd
}

func valueOrEmptyString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
