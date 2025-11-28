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

	"github.com/AlecAivazis/survey/v2"
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

func (o DeleteOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkConnectionID,
		flag.Force,
	}
}

func (o *DeleteOpts) MarkInteractive(cmd *cobra.Command) error {
	o.interactive = true
	for _, fn := range o.NonInteractiveFlags() {
		if f := cmd.Flags().Lookup(fn); f != nil && f.Changed {
			o.interactive = false
		}
	}
	if !o.interactive {
		// require cluster-id and plc-id in non-interactive mode
		for _, fn := range []string{flag.ClusterID, flag.PrivateLinkConnectionID} {
			if err := cmd.MarkFlagRequired(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteCmd(h *internal.Helper) *cobra.Command {
	opts := &DeleteOpts{interactive: true}
	var force bool

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "Delete a private link connection",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Delete a private link connection (interactive):
  $ %[1]s serverless private-link-connection delete

  Delete a private link connection (non-interactive):
  $ %[1]s serverless private-link-connection delete -c <cluster-id> -p <private-link-connection-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, plcID string
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
				privatelink, err := cloud.GetSelectedPrivateLinkConnection(ctx, cluster.ID, int32(h.QueryPageSize), d)
				if err != nil {
					return err
				}
				plcID = privatelink.ID
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				plcID, err = cmd.Flags().GetString(flag.PrivateLinkConnectionID)
				if err != nil {
					return errors.Trace(err)
				}
				force, err = cmd.Flags().GetBool(flag.Force)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if plcID == "" {
				return errors.New("private link connection id is required")
			}

			if !force {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support prompt, run with --force to skip confirmation")
				}
				var confirm string
				if err := survey.AskOne(&survey.Input{Message: "Type 'yes' to confirm deletion:"}, &confirm); err != nil {
					return err
				}
				if confirm != "yes" {
					return errors.New("deletion cancelled")
				}
			}

			_, err = d.DeletePrivateLinkConnection(ctx, clusterID, plcID)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Private link connection deleted"))
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	cmd.Flags().String(flag.PrivateLinkConnectionID, "", "The private link connection ID.")
	cmd.Flags().BoolVar(&force, flag.Force, false, "Delete without confirmation.")
	return cmd
}
