// Copyright 2026 PingCAP, Inc.
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

package upgrade

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/github"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/pingcap/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	installerURL    = "https://raw.githubusercontent.com/tidbcloud/tidbcloud-cli/40c0e2c7ef8b1a2ec0993b6935c01494777f64c5/install.sh"
	installerSHA256 = "d981d28e244778af2fe3cb63003ae22dd357a0d56e02bf7edbdbea7dd1288f18"
)

func Cmd(h *internal.Helper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade the CLI to the latest version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			// If is managed by TiUP, we should disable the upgrade command since binpath is different.
			if config.IsUnderTiUP {
				return errors.New("the CLI is managed by TiUP, please upgrade it by `tiup update cloud`")
			}

			// When upgrade CLI, we don't need to check the version again after command executes.
			newRelease, err := github.CheckForUpdate(ctx, false)
			// If we can't get the latest version, we should upgrade the CLI assuming it's not the latest version.
			if err != nil {
				log.Debug("Failed to check for upgrade, upgrade the CLI assuming it's not the latest version", zap.Error(err))
				newRelease = &github.ReleaseInfo{
					Version: "latest",
				}
			}

			if newRelease == nil {
				fmt.Fprintln(h.IOStreams.Out, "The CLI is already up to date.")
				return nil
			}

			if h.IOStreams.CanPrompt {
				return updateAndSpinnerWait(ctx, h, newRelease)
			} else {
				return updateAndWaitReady(ctx, h, newRelease)
			}
		},
	}

	return cmd
}

func updateAndWaitReady(ctx context.Context, h *internal.Helper, newRelease *github.ReleaseInfo) error {
	fmt.Fprintf(h.IOStreams.Out, "... Upgrading the CLI to version %s\n", newRelease.Version)
	return downloadAndRunInstaller(ctx, newRelease.Version)
}

func downloadAndRunInstaller(ctx context.Context, releaseVersion string) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	c1 := exec.CommandContext(ctx, "curl", "-q", "--proto", "=https", "--proto-redir", "=https", "-fsSL", installerURL)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c1.Stdout = &stdout
	c1.Stderr = &stderr

	err := c1.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("timeout when download the install.sh script")
	}
	if err != nil {
		fmt.Println(stderr.String())
		return err
	}
	if err := verifyInstallerChecksum(stdout.Bytes(), installerSHA256); err != nil {
		return err
	}

	c2 := exec.CommandContext(ctx, "/bin/sh", "-c", stdout.String()) //nolint:gosec // The installer checksum is verified before execution.
	version := ""
	if releaseVersion != "latest" {
		version = strings.TrimPrefix(releaseVersion, "v")
	}
	c2.Env = append(os.Environ(), "TICLOUD_VERSION="+version)
	stderr = bytes.Buffer{}
	c2.Stderr = &stderr

	err = c2.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("timeout when execute the install.sh script")
	}
	if err != nil {
		fmt.Println(stderr.String())
		return err
	}

	return nil
}

func verifyInstallerChecksum(script []byte, expectedChecksum string) error {
	checksum := sha256.Sum256(script)
	if hex.EncodeToString(checksum[:]) != expectedChecksum {
		return errors.New("checksum verification failed for install.sh")
	}
	return nil
}

func updateAndSpinnerWait(ctx context.Context, h *internal.Helper, newRelease *github.ReleaseInfo) error {
	task := func() tea.Msg {
		res := make(chan error, 1)

		go func() {
			res <- downloadAndRunInstaller(ctx, newRelease.Version)
		}()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case err := <-res:
				if err != nil {
					return err
				} else {
					return ui.Result("Upgrade successfully!")
				}
			case <-ctx.Done():
				return util.InterruptError
			case <-ticker.C:
				// continue
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, fmt.Sprintf("Upgrading the CLI to version %s", newRelease.Version)))
	model, err := p.Run()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := model.(ui.SpinnerModel); m.Interrupted {
		return util.InterruptError
	}
	if m, _ := model.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	} else {
		fmt.Fprintln(h.IOStreams.Out, color.GreenString(m.Output))
	}
	return nil
}
