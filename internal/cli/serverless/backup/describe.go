// Copyright 2024 PingCAP, Inc.
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

package backup

import (
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/service/cloud"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type DescribeOpts struct {
	interactive bool
}

func (c DescribeOpts) NonInteractiveFlags() []string {
	return []string{
		flag.BackupID,
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
		Short:   "Describe a backup",
		Aliases: []string{"get"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Get the backup in interactive mode:
  $ %[1]s serverless backup describe

  Get the backup in non-interactive mode:
  $ %[1]s serverless backup describe --backup-id <backup-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.MarkInteractive(cmd)
			if err != nil {
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

			var backupID string
			var clusterID string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				backup, err := cloud.GetSelectedServerlessBackup(ctx, clusterID, int32(h.QueryPageSize), d)
				if err != nil {
					return err
				}
				backupID = backup.ID
			} else {
				// non-interactive mode, get values from flags
				backupID, err = cmd.Flags().GetString(flag.BackupID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if err != nil {
				return errors.Trace(err)
			}
			backup, err := d.GetBackup(ctx, backupID)
			if err != nil {
				return errors.Trace(err)
			}
			err = output.PrintJson(h.IOStreams.Out, backup)
			if err != nil {
				return errors.Trace(err)
			}

			return nil
		},
	}

	describeCmd.Flags().String(flag.BackupID, "", "The ID of the backup to be described.")
	return describeCmd
}
