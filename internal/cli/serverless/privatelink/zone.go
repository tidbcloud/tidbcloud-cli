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

package privatelink

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type GetZonesOpts struct {
	interactive bool
}

func (o GetZonesOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func (o *GetZonesOpts) MarkInteractive(cmd *cobra.Command) error {
	o.interactive = true
	for _, fn := range o.NonInteractiveFlags() {
		if f := cmd.Flags().Lookup(fn); f != nil && f.Changed {
			o.interactive = false
		}
	}
	if !o.interactive {
		for _, fn := range o.NonInteractiveFlags() {
			if err := cmd.MarkFlagRequired(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetZonesCmd(h *internal.Helper) *cobra.Command {
	opts := &GetZonesOpts{interactive: true}

	cmd := &cobra.Command{
		Use:   "zones",
		Short: "Get account and availability zones for private link connection",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Get availability zones (interactive):
  $ %[1]s serverless private-link-connection get-zones

  Get availability zones (non-interactive):
  $ %[1]s serverless private-link-connection zones -c <cluster-id>`, config.CliName),
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
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			resp, err := d.GetPrivateLinkAvailabilityZones(ctx, clusterID)
			if err != nil {
				return errors.Trace(err)
			}
			return output.PrintJson(h.IOStreams.Out, resp)
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	return cmd
}
