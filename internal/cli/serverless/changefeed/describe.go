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

package changefeed

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
		flag.ChangefeedID,
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
		Short:   "Describe a changefeed",
		Aliases: []string{"get"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Get a changefeed in interactive mode:
  $ %[1]s serverless changefeed describe

  Get a changefeed in non-interactive mode:
  $ %[1]s serverless changefeed describe -c <cluster-id> --changefeed-id <changefeed-id>`, config.CliName),
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

			if changefeedID == "" {
				return errors.New("changefeed-id is required")
			}

			changefeed, err := d.GetConnector(ctx, clusterID, changefeedID)
			if err != nil {
				return errors.Trace(err)
			}

			err = output.PrintJson(h.IOStreams.Out, changefeed)
			return errors.Trace(err)
		},
	}

	describeCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	describeCmd.Flags().StringP(flag.ChangefeedID, flag.ChangefeedIDShort, "", "The changefeed ID.")
	return describeCmd
}
