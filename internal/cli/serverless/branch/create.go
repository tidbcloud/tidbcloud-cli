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
	"context"
	"fmt"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/branch"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

var createBranchField = map[string]int{
	flag.DisplayName:     0,
	flag.ParentTimestamp: 1,
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
		flag.DisplayName,
		flag.ClusterID,
		flag.ParentTimestamp,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
	}
}

func (c *CreateOpts) MarkInteractive(cmd *cobra.Command) error {
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
		for _, fn := range c.RequiredFlags() {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a branch",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a branch in interactive mode:
  $ %[1]s serverless branch create

  Create a branch in non-interactive mode:
  $ %[1]s serverless branch create --cluster-id <cluster-id> --display-name <branch-name> --parent-id <parent-id>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.MarkInteractive(cmd)
			if err != nil {
				return errors.Trace(err)
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
			var parentID string
			var parentTimestampStr string
			var parentTimestamp time.Time
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
				clusterId = cluster.ID

				parentID, err = cloud.GetSelectedParentID(ctx, cluster, h.QueryPageSize, d)
				if err != nil {
					return err
				}

				// variables for input
				inputModel, err := GetCreateBranchInput()
				if err != nil {
					return err
				}
				branchName = inputModel.(ui.TextInputModel).Inputs[createBranchField[flag.DisplayName]].Value()
				if len(branchName) == 0 {
					return errors.New("branch name is required")
				}
				parentTimestampStr = inputModel.(ui.TextInputModel).Inputs[createBranchField[flag.ParentTimestamp]].Value()
			} else {
				// non-interactive mode, get values from flags
				var err error
				branchName, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}
				clusterId, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				parentID, err = cmd.Flags().GetString(flag.ParentID)
				if err != nil {
					return errors.Trace(err)
				}
				parentTimestampStr, err = cmd.Flags().GetString(flag.ParentTimestamp)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if len(parentTimestampStr) != 0 {
				parentTimestamp, err = time.Parse(time.RFC3339, parentTimestampStr)
				if err != nil {
					return errors.New("Invalid parent timestamp format, please use RFC3339 format")
				}
			}

			body := branch.NewBranch(branchName)
			if parentID != "" {
				body.SetParentId(parentID)
			}
			if !parentTimestamp.IsZero() {
				body.SetParentTimestamp(parentTimestamp)
			}
			if h.IOStreams.CanPrompt {
				err := CreateAndSpinnerWait(ctx, h, d, clusterId, body)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				err := CreateAndWaitReady(ctx, h, d, clusterId, body)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	createCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The displayName of the branch to be created.")
	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster, in which the branch will be created.")
	createCmd.Flags().StringP(flag.ParentID, "", "", "The ID of the branch parent, default is cluster id.")
	createCmd.Flags().StringP(flag.ParentTimestamp, "", "", "The timestamp of the parent branch, default is current time. (RFC3339 format, e.g., 2024-01-01T00:00:00Z)")
	return createCmd
}

func CreateAndWaitReady(ctx context.Context, h *internal.Helper, d cloud.TiDBCloudClient, clusterId string, body *branch.Branch) error {
	b, err := d.CreateBranch(ctx, clusterId, body)
	if err != nil {
		return errors.Trace(err)
	}
	newBranchID := *b.BranchId

	fmt.Fprintln(h.IOStreams.Out, "... Waiting for branch to be ready")
	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()
	timer := time.After(WaitTimeout)
	for {
		select {
		case <-timer:
			return errors.New(fmt.Sprintf("Timeout waiting for branch %s to be ready, please check status on dashboard.", newBranchID))
		case <-ticker.C:
			b, err := d.GetBranch(ctx, clusterId, newBranchID)
			if err != nil {
				return errors.Trace(err)
			}
			if *b.State == branch.BRANCHSTATE_ACTIVE {
				fmt.Fprint(h.IOStreams.Out, color.GreenString("Branch %s is ready.", newBranchID))
				return nil
			}
		}
	}
}

func CreateAndSpinnerWait(ctx context.Context, h *internal.Helper, d cloud.TiDBCloudClient, clusterId string, body *branch.Branch) error {
	// use spinner to indicate that the cluster is being created
	task := func() tea.Msg {
		b, err := d.CreateBranch(ctx, clusterId, body)
		if err != nil {
			return errors.Trace(err)
		}
		newBranchID := *b.BranchId

		ticker := time.NewTicker(WaitInterval)
		defer ticker.Stop()
		timer := time.After(WaitTimeout)
		for {
			select {
			case <-timer:
				return ui.Result(fmt.Sprintf("Timeout waiting for branch %s to be ready, please check status on dashboard.", newBranchID))
			case <-ticker.C:
				b, err := d.GetBranch(ctx, clusterId, newBranchID)
				if err != nil {
					return errors.Trace(err)
				}
				if *b.State == branch.BRANCHSTATE_ACTIVE {
					return ui.Result(fmt.Sprintf("Branch %s is ready.", newBranchID))
				}
			case <-ctx.Done():
				return util.InterruptError
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Waiting for branch to be ready"))
	createModel, err := p.Run()
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
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 64

		switch k {
		case flag.DisplayName:
			t.Placeholder = "Display Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.ParentTimestamp:
			timestampExample := time.Now().Format(time.RFC3339)
			t.Placeholder = fmt.Sprintf("Parent Timestamp (optional, e.g., %s)", timestampExample)
		}
		m.Inputs[v] = t
	}
	return m
}

func GetCreateBranchInput() (tea.Model, error) {
	p := tea.NewProgram(initialCreateBranchInputModel())
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}
