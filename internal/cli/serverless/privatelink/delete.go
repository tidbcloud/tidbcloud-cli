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

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type DeleteOpts struct {
	interactive bool
}

func (c DeleteOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkConnectionID,
	}
}

func (c *DeleteOpts) MarkInteractive(cmd *cobra.Command) error {
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

func DeleteCmd(h *internal.Helper) *cobra.Command {
	opts := DeleteOpts{
		interactive: true,
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a private link connection",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Delete a private link connection in interactive mode:
  $ %[1]s serverless private-link-connection delete

  Delete a private link connection in non-interactive mode:
  $ %[1]s serverless private-link-connection delete -c <cluster-id> --private-link-connection-id <connection-id>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}

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

			_, err = d.DeletePrivateLinkConnection(ctx, clusterID, connectionID)
			if err != nil {
				return errors.Trace(err)
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("private link connection %s is deleted", connectionID))
			if err != nil {
				return err
			}
			return nil
		},
	}

	deleteCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	deleteCmd.Flags().StringP(flag.PrivateLinkConnectionID, "", "", "The private link connection ID.")
	return deleteCmd
}
