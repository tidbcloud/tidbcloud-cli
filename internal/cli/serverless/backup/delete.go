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
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/util"
	brApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

const confirmed = "yes"

type DeleteOpts struct {
	interactive bool
}

func (c DeleteOpts) NonInteractiveFlags() []string {
	return []string{
		flag.BackupID,
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

func DeleteCmd(h *internal.Helper) *cobra.Command {
	opts := DeleteOpts{
		interactive: true,
	}

	var force bool
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a backup",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Delete a backup in interactive mode:
  $ %[1]s serverless backup delete

  Delete a backup in non-interactive mode:
  $ %[1]s serverless backup delete --backup-id <backup-id>`, config.CliName),
		Aliases: []string{"rm"},
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
				backup, err := cloud.GetSelectedServerlessBackup(ctx, cluster.ID, int32(h.QueryPageSize), d)
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

			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to delete the backup")
				}

				confirmationMessage := fmt.Sprintf("%s %s %s", color.BlueString("Please type"), color.HiBlueString(confirmed), color.BlueString("to confirm:"))

				prompt := &survey.Input{
					Message: confirmationMessage,
				}

				var userInput string
				err := survey.AskOne(prompt, &userInput)
				if err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					} else {
						return err
					}
				}

				if userInput != confirmed {
					return errors.New("incorrect confirm string entered, skipping backup deletion")
				}
			}

			params := brApi.NewBackupRestoreServiceDeleteBackupParams().WithBackupID(backupID).WithContext(ctx)
			_, err = d.DeleteBackup(params)
			if err != nil {
				return errors.Trace(err)
			}
			// print success for delete branch is a sync operation
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("backup %s deleted", backupID))
			return nil
		},
	}

	deleteCmd.Flags().BoolVar(&force, flag.Force, false, "Delete a backup without confirmation.")
	deleteCmd.Flags().String(flag.BackupID, "", "The ID of the backup to be deleted.")

	return deleteCmd
}
