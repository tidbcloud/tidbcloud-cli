// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package migration

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
		flag.MigrationID,
	}
}

func (c *DescribeOpts) MarkInteractive(cmd *cobra.Command) error {
	for _, fn := range c.NonInteractiveFlags() {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	if !c.interactive {
		for _, fn := range c.NonInteractiveFlags() {
			if err := cmd.MarkFlagRequired(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func DescribeCmd(h *internal.Helper) *cobra.Command {
	opts := DescribeOpts{interactive: true}

	var cmd = &cobra.Command{
		Use:     "describe",
		Short:   "Describe a migration",
		Aliases: []string{"get"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Describe a migration in interactive mode:
  $ %[1]s serverless migration describe

  Describe a migration in non-interactive mode:
  $ %[1]s serverless migration describe -c <cluster-id> --migration-id <migration-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, migrationID string
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

				migration, err := cloud.GetSelectedMigration(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				migrationID = migration.ID
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				migrationID, err = cmd.Flags().GetString(flag.MigrationID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			resp, err := d.GetMigration(ctx, clusterID, migrationID)
			if err != nil {
				return errors.Trace(err)
			}

			return errors.Trace(output.PrintJson(h.IOStreams.Out, resp))
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID that owns the migration.")
	cmd.Flags().StringP(flag.MigrationID, flag.MigrationIDShort, "", "ID of the migration to describe.")
	return cmd
}
