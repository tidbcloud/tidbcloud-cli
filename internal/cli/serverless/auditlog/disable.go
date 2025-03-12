// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package auditlog

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
)

type DisableOpts struct {
	interactive bool
}

func (c DisableOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func (c *DisableOpts) MarkInteractive(cmd *cobra.Command) error {
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

func DisableCmd(h *internal.Helper) *cobra.Command {
	opts := DisableOpts{
		interactive: true,
	}

	var disableCmd = &cobra.Command{
		Use:   "disable",
		Short: "disable the database audit logging",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  disable the database audit logging in interactive mode:
  $ %[1]s serverless audit-log disable

  disable the database audit logging in non-interactive mode:
  $ %[1]s serverless audit-log disable -c <cluster-id> `, config.CliName),
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

			enable := false
			body := &cluster.V1beta1ServerlessServicePartialUpdateClusterBody{
				Cluster: &cluster.RequiredTheClusterToBeUpdated{
					AuditLogConfig: &cluster.V1beta1ClusterAuditLogConfig{
						Enabled: *cluster.NewNullableBool(&enable),
					},
				},
				UpdateMask: "auditLogConfig",
			}
			_, err = d.PartialUpdateCluster(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("cluster %s database audit logging disabled", clusterID)))
			return nil
		},
	}

	disableCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	return disableCmd
}
