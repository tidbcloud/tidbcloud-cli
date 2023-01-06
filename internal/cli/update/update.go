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
	"bytes"
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
			// If is managed by TiUP, we should disable the update command since binpath is different.
			if h.IsUnderTiUP {
				return errors.New("the CLI is managed by TiUP, please update it by `tiup update cloud`")
			}

			// When update CLI, we don't need to check the version again after command executes.
			newRelease, err := github.CheckForUpdate(config.Repo, ver, false)
			// If we can't get the latest version, we should update the CLI assuming it's not the latest version.
			if err != nil {
				newRelease = &github.ReleaseInfo{
					Version: "latest",
				}
			}

			if newRelease == nil {
				fmt.Fprintln(h.IOStreams.Out, "The CLI is already up to date.")
				return nil
			}

			if h.IOStreams.CanPrompt {
				return updateAndSpinnerWait(h, newRelease)
			} else {
				return updateAndWaitReady(h, newRelease)
			}
		},
	}

	return cmd
}

func updateAndWaitReady(h *internal.Helper, newRelease *github.ReleaseInfo) error {
	fmt.Fprintf(h.IOStreams.Out, "... Updating the CLI to version %s\n", newRelease.Version)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	c1 := exec.CommandContext(ctx, "curl", "https://raw.githubusercontent.com/tidbcloud/tidbcloud-cli/main/install.sh") //nolint:gosec
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c1.Stdout = &stdout
	c1.Stderr = &stderr

	err := c1.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("timeout when download the install.sh script")
	}
	if err != nil {
		return errors.Annotate(err, stderr.String())
	}

	c2 := exec.CommandContext(ctx, "/bin/sh", "-c", stdout.String()) //nolint:gosec
	stderr = bytes.Buffer{}
	c2.Stderr = &stderr

	err = c2.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("timeout when execute the install.sh script")
	}
	if err != nil {
		return errors.Annotate(err, stderr.String())
	}

	return nil
}

func updateAndSpinnerWait(h *internal.Helper, newRelease *github.ReleaseInfo) error {
	task := func() tea.Msg {
		res := make(chan error)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()
			out, err := exec.CommandContext(ctx, "curl", "https://raw.githubusercontent.com/tidbcloud/tidbcloud-cli/main/install.sh").Output()
			if ctx.Err() == context.DeadlineExceeded {
				res <- errors.New("timeout when download the install.sh script")
				return
			}
			if err != nil {
				res <- errors.Annotate(err, string(out))
				return
			}

			out1, err := exec.CommandContext(ctx, "/bin/sh", "-c", string(out)).Output() //nolint:gosec
			if ctx.Err() == context.DeadlineExceeded {
				res <- errors.New("timeout when execute the install.sh script")
				return
			}
			if err != nil {
				res <- errors.Annotate(err, string(out1))
				return
			}

			res <- nil
		}()

		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case err := <-res:
				if err != nil {
					return err
				} else {
					return ui.Result("Update successfully!")
				}
			case <-ticker.C:
				// continue
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
