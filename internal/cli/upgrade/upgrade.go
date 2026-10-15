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
	"bufio"
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
	"github.com/hashicorp/go-version"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

const releaseDownloadURL = "https://github.com/tidbcloud/tidbcloud-cli/releases/download"

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
			if err != nil {
				return errors.Annotate(err, "failed to check for CLI updates")
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
	if _, err := version.NewSemver(releaseVersion); err != nil {
		return errors.Annotate(err, "invalid release version")
	}
	releaseTag := "v" + strings.TrimPrefix(releaseVersion, "v")

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	script, err := downloadReleaseAsset(ctx, releaseTag, "install.sh")
	if err != nil {
		return err
	}
	checksums, err := downloadReleaseAsset(ctx, releaseTag, "checksums.txt")
	if err != nil {
		return err
	}
	if err := verifyInstallerChecksum(script, checksums); err != nil {
		return err
	}

	c2 := exec.CommandContext(ctx, "/bin/sh", "-c", string(script)) //nolint:gosec // The installer checksum is verified before execution.
	c2.Env = append(os.Environ(), "TICLOUD_VERSION="+strings.TrimPrefix(releaseTag, "v"))
	var stderr bytes.Buffer
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

func downloadReleaseAsset(ctx context.Context, releaseTag, asset string) ([]byte, error) {
	assetURL := releaseDownloadURL + "/" + releaseTag + "/" + asset
	cmd := exec.CommandContext(ctx, "curl", "-q", "--proto", "=https", "--proto-redir", "=https", "-fsSL", assetURL) //nolint:gosec // The release version is validated before constructing the URL.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		} else if stderr.Len() != 0 {
			err = errors.Annotate(err, strings.TrimSpace(stderr.String()))
		}
		return nil, errors.Annotatef(err, "failed to download %s for %s", asset, releaseTag)
	}
	return stdout.Bytes(), nil
}

func verifyInstallerChecksum(script, checksums []byte) error {
	var expectedChecksum []byte
	scanner := bufio.NewScanner(bytes.NewReader(checksums))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 || strings.TrimPrefix(fields[1], "*") != "install.sh" {
			continue
		}
		if expectedChecksum != nil {
			return errors.New("duplicate checksum for install.sh")
		}
		var err error
		expectedChecksum, err = hex.DecodeString(fields[0])
		if err != nil || len(expectedChecksum) != sha256.Size {
			return errors.New("invalid checksum for install.sh")
		}
	}
	if err := scanner.Err(); err != nil {
		return errors.Annotate(err, "failed to read checksums.txt")
	}
	if expectedChecksum == nil {
		return errors.New("missing checksum for install.sh")
	}
	checksum := sha256.Sum256(script)
	if !bytes.Equal(checksum[:], expectedChecksum) {
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
