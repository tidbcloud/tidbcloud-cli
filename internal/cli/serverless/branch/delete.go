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

package branch

import (
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

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
		flag.ClusterID,
		flag.BranchID,
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
		Short: "Delete a branch",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Delete a branch in interactive mode:
  $ %[1]s serverless branch delete

  Delete a branch in non-interactive mode:
  $ %[1]s serverless branch delete -c <cluster-id> -b <branch-id>`, config.CliName),
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

			var clusterID string
			var branchID string
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

				branch, err := cloud.GetSelectedBranch(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				branchID = branch.ID
			} else {
				// non-interactive mode, get values from flags
				bID, err := cmd.Flags().GetString(flag.BranchID)
				if err != nil {
					return errors.Trace(err)
				}

				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				branchID = bID
				clusterID = cID
			}

			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to delete the branch")
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
					return errors.New("incorrect confirm string entered, skipping branch deletion")
				}
			}

			_, err = d.DeleteBranch(ctx, clusterID, branchID)
			if err != nil {
				return errors.Trace(err)
			}
			// print success for delete branch is a sync operation
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("branch %s deleted", branchID))
			return nil
		},
	}

	deleteCmd.Flags().BoolVar(&force, flag.Force, false, "Delete a branch without confirmation.")
	deleteCmd.Flags().StringP(flag.BranchID, flag.BranchIDShort, "", "The ID of the branch to be deleted.")
	deleteCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID of the branch to be deleted.")

	return deleteCmd
}
