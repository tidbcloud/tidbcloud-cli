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

package branch

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/branch"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type ResetOpts struct {
	interactive bool
}

func (c ResetOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.BranchID,
	}
}

func (c *ResetOpts) MarkInteractive(cmd *cobra.Command) error {
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

func ResetCmd(h *internal.Helper) *cobra.Command {
	opts := ResetOpts{
		interactive: true,
	}

	var force bool
	var resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Reset a branch to its parent's latest state",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Reset a branch in interactive mode:
  $ %[1]s serverless branch reset

  Reset a branch in non-interactive mode:
  $ %[1]s serverless branch reset -c <cluster-id> -b <branch-id>`, config.CliName),
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
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to reset the branch")
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
					return errors.New("incorrect confirm string entered, skipping branch reset")
				}
			}

			// print success for reset branch is a sync operation
			if h.IOStreams.CanPrompt {
				err := ResetAndSpinnerWait(ctx, h, d, clusterID, branchID)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				err := ResetAndWaitReady(ctx, h, d, clusterID, branchID)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	resetCmd.Flags().BoolVar(&force, flag.Force, false, "Reset a branch without confirmation.")
	resetCmd.Flags().StringP(flag.BranchID, flag.BranchIDShort, "", "The ID of the branch to be reset.")
	resetCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID of the branch to be reset.")

	return resetCmd
}

func ResetAndWaitReady(ctx context.Context, h *internal.Helper, d cloud.TiDBCloudClient, clusterID string, branchID string) error {
	_, err := d.ResetBranch(ctx, clusterID, branchID)
	if err != nil {
		return errors.Trace(err)
	}

	fmt.Fprintln(h.IOStreams.Out, "... Waiting for branch to be ready")
	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()
	timer := time.After(WaitTimeout)
	for {
		select {
		case <-timer:
			return errors.New(fmt.Sprintf("Timeout waiting for branch %s to be ready, please check status on dashboard.", branchID))
		case <-ticker.C:
			b, err := d.GetBranch(ctx, clusterID, branchID, branch.BRANCHSERVICEGETBRANCHVIEWPARAMETER_BASIC)
			if err != nil {
				return errors.Trace(err)
			}
			if *b.State == branch.BRANCHSTATE_ACTIVE {
				fmt.Fprint(h.IOStreams.Out, color.GreenString("Branch %s is ready.", branchID))
				return nil
			}
		}
	}
}

func ResetAndSpinnerWait(ctx context.Context, h *internal.Helper, d cloud.TiDBCloudClient, clusterID string, branchID string) error {
	// use spinner to indicate that the cluster is resetting
	task := func() tea.Msg {
		_, err := d.ResetBranch(ctx, clusterID, branchID)
		if err != nil {
			return errors.Trace(err)
		}

		ticker := time.NewTicker(WaitInterval)
		defer ticker.Stop()
		timer := time.After(WaitTimeout)
		for {
			select {
			case <-timer:
				return ui.Result(fmt.Sprintf("Timeout waiting for branch %s to be ready, please check status on dashboard.", branchID))
			case <-ticker.C:
				b, err := d.GetBranch(ctx, clusterID, branchID, branch.BRANCHSERVICEGETBRANCHVIEWPARAMETER_BASIC)
				if err != nil {
					return errors.Trace(err)
				}
				if *b.State == branch.BRANCHSTATE_ACTIVE {
					return ui.Result(fmt.Sprintf("Branch %s is ready.", branchID))
				}
			case <-ctx.Done():
				return util.InterruptError
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Waiting for branch to be ready"))
	resetModel, err := p.Run()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := resetModel.(ui.SpinnerModel); m.Interrupted {
		return util.InterruptError
	}
	if m, _ := resetModel.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	} else {
		fmt.Fprintln(h.IOStreams.Out, color.GreenString(m.Output))
	}
	return nil
}
