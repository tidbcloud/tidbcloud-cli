// Copyright 2023 PingCAP, Inc.
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

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	branchApi "tidbcloud-cli/pkg/tidbcloud/branch/client/branch_service"
	branchModel "tidbcloud-cli/pkg/tidbcloud/branch/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

var createBranchField = map[string]int{
	flag.BranchName: 0,
}

const (
	WaitInterval = 5 * time.Second
	WaitTimeout  = 5 * time.Minute
)

type CreateOpts struct {
	interactive bool
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.BranchName,
		flag.ClusterID,
	}
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a branch in a specified serverless cluster",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a branch in interactive mode:
  $ %[1]s branch create

  Create a branch in non-interactive mode:
  $ %[1]s branch create --cluster-id <cluster-id> --branch-name <branch-name>`,
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
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}

			var branchName string
			var clusterId string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}

				cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterId = cluster.ID

				// variables for input
				inputModel, err := GetCreateBranchInput()
				if err != nil {
					return err
				}
				branchName = inputModel.(ui.TextInputModel).Inputs[createBranchField[flag.BranchName]].Value()
				if len(branchName) == 0 {
					return errors.New("branch name is required")
				}
			} else {
				// non-interactive mode, get values from flags
				var err error
				branchName, err = cmd.Flags().GetString(flag.BranchName)
				if err != nil {
					return errors.Trace(err)
				}
				clusterId, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			params := branchApi.NewCreateBranchParams().WithClusterID(clusterId).WithBody(branchApi.CreateBranchBody{
				DisplayName: &branchName,
			})

			if h.IOStreams.CanPrompt {
				err := CreateAndSpinnerWait(ctx, d, params, h)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				err := CreateAndWaitReady(h, d, params)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	createCmd.Flags().String(flag.BranchName, "", "Name of the branch to de created")
	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster, in which the branch will be created")
	return createCmd
}

func CreateAndWaitReady(h *internal.Helper, d cloud.TiDBCloudClient, params *branchApi.CreateBranchParams) error {
	createBranchResult, err := d.CreateBranch(params)
	if err != nil {
		return errors.Trace(err)
	}
	newBranchID := *createBranchResult.GetPayload().ID

	fmt.Fprintln(h.IOStreams.Out, "... Waiting for branch to be ready")
	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()
	timer := time.After(WaitTimeout)
	for {
		select {
		case <-timer:
			return errors.New(fmt.Sprintf("Timeout waiting for branch %s to be ready, please check status on dashboard.", newBranchID))
		case <-ticker.C:
			clusterResult, err := d.GetBranch(branchApi.NewGetBranchParams().
				WithClusterID(params.ClusterID).
				WithBranchID(newBranchID))
			if err != nil {
				return errors.Trace(err)
			}
			s := clusterResult.GetPayload().State
			if *s == branchModel.OpenapiBranchStateREADY {
				fmt.Fprint(h.IOStreams.Out, color.GreenString("Branch %s is ready.", newBranchID))
				return nil
			}
		}
	}
}

func CreateAndSpinnerWait(ctx context.Context, d cloud.TiDBCloudClient, params *branchApi.CreateBranchParams, h *internal.Helper) error {
	// use spinner to indicate that the cluster is being created
	task := func() tea.Msg {
		createBranchResult, err := d.CreateBranch(params)
		if err != nil {
			return errors.Trace(err)
		}
		newBranchID := *createBranchResult.GetPayload().ID

		ticker := time.NewTicker(WaitInterval)
		defer ticker.Stop()
		timer := time.After(WaitTimeout)
		for {
			select {
			case <-timer:
				return ui.Result(fmt.Sprintf("Timeout waiting for branch %s to be ready, please check status on dashboard.", newBranchID))
			case <-ticker.C:
				clusterResult, err := d.GetBranch(branchApi.NewGetBranchParams().
					WithClusterID(params.ClusterID).
					WithBranchID(newBranchID))
				if err != nil {
					return errors.Trace(err)
				}
				s := clusterResult.GetPayload().State
				if *s == branchModel.OpenapiBranchStateREADY {
					return ui.Result(fmt.Sprintf("Branch %s is ready.", newBranchID))
				}
			case <-ctx.Done():
				return util.InterruptError
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Waiting for branch to be ready"))
	createModel, err := p.StartReturningModel()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := createModel.(ui.SpinnerModel); m.Interrupted {
		return util.InterruptError
	}
	if m, _ := createModel.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	} else {
		fmt.Fprintln(h.IOStreams.Out, color.GreenString(m.Output))
	}
	return nil
}

func initialCreateBranchInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(createBranchField)),
	}

	for k, v := range createBranchField {
		t := textinput.New()
		t.CursorStyle = config.CursorStyle
		t.CharLimit = 64

		switch k {
		case flag.BranchName:
			t.Placeholder = "Branch Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		}

		m.Inputs[v] = t
	}

	return m
}

func GetCreateBranchInput() (tea.Model, error) {
	p := tea.NewProgram(initialCreateBranchInputModel())
	inputModel, err := p.StartReturningModel()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}
