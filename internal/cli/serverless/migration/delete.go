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

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
)

type DeleteOpts struct {
	interactive bool
}

func (c DeleteOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.MigrationTaskID,
	}
}

func (c *DeleteOpts) MarkInteractive(cmd *cobra.Command) error {
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

func DeleteCmd(h *internal.Helper) *cobra.Command {
	opts := DeleteOpts{interactive: true}
	var force bool

	var cmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete a migration task",
		Aliases: []string{"rm"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Delete a migration task in interactive mode:
	  $ %[1]s serverless migration delete

	  Delete a migration task in non-interactive mode:
	  $ %[1]s serverless migration delete -c <cluster-id> --migration-id <task-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, taskID string
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
				task, err := cloud.GetSelectedMigrationTask(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				taskID = task.ID
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				taskID, err = cmd.Flags().GetString(flag.MigrationTaskID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if !force {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support prompt, please run with --force to delete the migration task")
				}
				prompt := &survey.Input{
					Message: fmt.Sprintf("%s %s %s", color.BlueString("Please type"), color.HiBlueString("yes"), color.BlueString("to confirm:")),
				}
				var confirmation string
				if err := survey.AskOne(prompt, &confirmation); err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					}
					return err
				}
				if confirmation != "yes" {
					return errors.New("Incorrect confirm string entered, skipping migration task deletion")
				}
			}

			if _, err := d.DeleteMigration(ctx, clusterID, taskID); err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("migration task %s deleted", taskID))
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, flag.Force, false, "Delete without confirmation.")
	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID that owns the migration task.")
	cmd.Flags().StringP(flag.MigrationTaskID, flag.MigrationTaskIDShort, "", "ID of the migration task to delete.")
	return cmd
}
