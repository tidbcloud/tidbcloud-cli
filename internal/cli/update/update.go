// Copyright 2022 PingCAP, Inc.
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

package update

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/service/github"
	"tidbcloud-cli/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func UpdateCmd(h *internal.Helper, ver string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the CLI to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			// When update CLI, we don't need to check the version again after command executes.
			newRelease, err := github.CheckForUpdate(config.Repo, ver, false)
			if err != nil {
				return err
			}
			if newRelease == nil {
				fmt.Fprintln(h.IOStreams.Out, "The CLI is already up to date.")
				return nil
			}

			if h.IOStreams.CanPrompt {
				return CreateAndSpinnerWait(h, newRelease)
			} else {
				return CreateAndWaitReady(h, newRelease)
			}
		},
	}

	return cmd
}

func CreateAndWaitReady(h *internal.Helper, newRelease *github.ReleaseInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	fmt.Fprintf(h.IOStreams.Out, "... Updating the CLI to version %s\n", newRelease.Version)
	c1 := exec.CommandContext(ctx, "curl", "https://raw.githubusercontent.com/tidbcloud/tidbcloud-cli/main/install.sh")
	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("timeout when download the install.sh script")
	}

	out, err := c1.Output()
	if err != nil {
		return errors.Annotate(err, "failed to download the install.sh script")
	}

	_, err = exec.CommandContext(ctx, "/bin/sh", "-c", string(out)).Output() //nolint:gosec
	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("timeout when execute the install.sh script")
	}
	if err != nil {
		return errors.Annotate(err, "execute the install.sh script")
	}

	fmt.Fprintln(h.IOStreams.Out, "Update successfully!")

	return nil
}

func CreateAndSpinnerWait(h *internal.Helper, newRelease *github.ReleaseInfo) error {
	task := func() tea.Msg {
		res := make(chan error)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()
			c1 := exec.CommandContext(ctx, "curl", "https://raw.githubusercontent.com/tidbcloud/tidbcloud-cli/main/install.sh")
			if ctx.Err() == context.DeadlineExceeded {
				res <- errors.New("timeout when download the install.sh script")
			}

			out, err := c1.Output()
			if err != nil {
				res <- errors.Annotate(err, "failed to download the install.sh script")
			}

			_, err = exec.CommandContext(ctx, "/bin/sh", "-c", string(out)).Output() //nolint:gosec
			if ctx.Err() == context.DeadlineExceeded {
				res <- errors.New("timeout when execute the install.sh script")
			}
			if err != nil {
				res <- errors.Annotate(err, "execute the install.sh script")
			}

			res <- nil
		}()

		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				select {
				case err := <-res:
					if err != nil {
						return err
					} else {
						return ui.Result("Update successfully!")
					}
				default:
					// continue
				}
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, fmt.Sprintf("Updating the CLI to version %s", newRelease.Version)))
	createModel, err := p.StartReturningModel()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := createModel.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	} else {
		fmt.Fprintln(h.IOStreams.Out, color.GreenString(m.Output))
	}
	return nil
}
