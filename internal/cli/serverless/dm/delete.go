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

package dm

import (
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type DeleteOpts struct {
	interactive bool
}

func (c DeleteOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.TaskID,
	}
}

func DeleteCmd(h *internal.Helper) *cobra.Command {
	var force bool
	opts := DeleteOpts{
		interactive: true,
	}

	var deleteCmd = &cobra.Command{
		Use:         "delete",
		Short:       "Delete a DM task",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Delete a DM task in interactive mode:
  $ %[1]s serverless dm delete

  Delete a DM task in non-interactive mode:
  $ %[1]s serverless dm delete --cluster-id <cluster-id> --task-id <task-id>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// mark required flags in non-interactive mode
			if !opts.interactive {
				for _, fn := range flags {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var clusterID, taskID string
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
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

				task, err := cloud.GetSelectedDMTask(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				taskID = *task.Id
			} else {
				// non-interactive mode
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				taskID = cmd.Flag(flag.TaskID).Value.String()
			}

			cmd.Annotations[telemetry.ClusterID] = clusterID

			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to delete the DM task")
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
					return errors.New("incorrect confirm string entered, skipping DM task deletion")
				}
			}

			_, err = d.DeleteTask(ctx, clusterID, taskID)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintf(h.IOStreams.Out, "DM task %s is being deleted.\n", taskID)
			return nil
		},
	}

	deleteCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	deleteCmd.Flags().String(flag.TaskID, "", "Task ID.")
	deleteCmd.Flags().BoolVar(&force, flag.Force, false, "Delete without confirmation.")
	return deleteCmd
}

const confirmed = "delete"
