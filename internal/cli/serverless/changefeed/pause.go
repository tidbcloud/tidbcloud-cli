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

package changefeed

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type PauseOpts struct {
	interactive bool
}

func (c PauseOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ChangefeedID,
	}
}

func (c *PauseOpts) MarkInteractive(cmd *cobra.Command) error {
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

func PauseCmd(h *internal.Helper) *cobra.Command {
	opts := PauseOpts{
		interactive: true,
	}

	var pauseCmd = &cobra.Command{
		Use:   "pause",
		Short: "Pause a changefeed",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Pause a changefeed in interactive mode:
  $ %[1]s serverless changefeed pause

  Pause a changefeed in non-interactive mode:
  $ %[1]s serverless changefeed pause -c <cluster-id> --changefeed-id <changefeed-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, changefeedID string
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

				cf, err := cloud.GetSelectedChangefeed(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				changefeedID = cf.ID
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				changefeedID, err = cmd.Flags().GetString(flag.ChangefeedID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			_, err = d.StopChangefeed(ctx, clusterID, changefeedID)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintf(h.IOStreams.Out, "changefeed %s paused\n", changefeedID)
			return nil
		},
	}

	pauseCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	pauseCmd.Flags().StringP(flag.ChangefeedID, flag.ChangefeedIDShort, "", "The changefeed ID.")
	return pauseCmd
}
