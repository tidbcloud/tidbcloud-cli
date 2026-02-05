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

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type DescribeOpts struct {
	interactive bool
}

func (c DescribeOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkConnectionID,
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
	if !c.interactive {
		for _, fn := range flags {
			if err := cmd.MarkFlagRequired(fn); err != nil {
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
		Short:   "Describe a private link connection",
		Aliases: []string{"get"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Get a private link connection in interactive mode:
  $ %[1]s serverless private-link-connection describe

  Get a private link connection in non-interactive mode:
  $ %[1]s serverless private-link-connection describe -c <cluster-id> --private-link-connection-id <connection-id>`,
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

			var clusterID, connectionID string
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

				connection, err := cloud.GetSelectedPrivateLinkConnection(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				connectionID = connection.ID
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				connectionID, err = cmd.Flags().GetString(flag.PrivateLinkConnectionID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			connection, err := d.GetPrivateLinkConnection(ctx, clusterID, connectionID)
			if err != nil {
				return errors.Trace(err)
			}

			err = output.PrintJson(h.IOStreams.Out, connection)
			return errors.Trace(err)
		},
	}

	describeCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	describeCmd.Flags().StringP(flag.PrivateLinkConnectionID, "", "", "The private link connection ID.")
	return describeCmd
}
