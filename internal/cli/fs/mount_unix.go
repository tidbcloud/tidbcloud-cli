//go:build !windows

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

package fs

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	fusepkg "github.com/tidbcloud/tidbcloud-cli/internal/service/fs/fuse"

	"github.com/spf13/cobra"
)

func mountCmd(h *internal.Helper) *cobra.Command {
	var (
		readOnly   bool
		debug      bool
		allowOther bool
	)

	var cmd = &cobra.Command{
		Use:   "mount <mountpoint>",
		Short: "Mount remote filesystem via FUSE",
		Long: `Mount the remote filesystem as a local FUSE mount.

WARNING: This feature requires FUSE support and may need additional setup.
On Linux, you need fusermount or fusermount3 installed.
On macOS, you need macFUSE installed.

Examples:
  ticloud fs mount /mnt/tidbcloud
  ticloud fs mount /mnt/tidbcloud --read-only --debug`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mountpoint := args[0]

			// Check if mountpoint exists
			if _, err := os.Stat(mountpoint); os.IsNotExist(err) {
				return fmt.Errorf("mountpoint %s does not exist", mountpoint)
			}

			// Check if already mounted
			if isMounted(mountpoint) {
				return fmt.Errorf("%s is already mounted", mountpoint)
			}

			client, err := newClient(cmd)
			if err != nil {
				return suggestInitIfTenantNotFound(err)
			}

			// For now, only support read-only mode (MVP)
			// Full read-write support will be added later
			if !readOnly {
				fmt.Fprintln(h.IOStreams.Out, "Note: Mounting in read-only mode (MVP). Write support coming soon.")
				readOnly = true
			}

			opts := &fusepkg.MountOptions{
				Client:     client,
				MountPoint: mountpoint,
				ReadOnly:   readOnly,
				Debug:      debug,
				AllowOther: allowOther,
			}

			return suggestInitIfTenantNotFound(fusepkg.Mount(opts))
		},
	}

	cmd.Flags().BoolVarP(&readOnly, "read-only", "r", true, "Mount as read-only (default: true for MVP)")
	cmd.Flags().BoolVar(&debug, "debug", false, "Enable FUSE debug logging")
	cmd.Flags().BoolVar(&allowOther, "allow-other", false, "Allow other users to access mount")

	return cmd
}

func umountCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:     "umount <mountpoint>",
		Short:   "Unmount FUSE filesystem",
		Example: `  ticloud fs umount /mnt/tidbcloud`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mountpoint := args[0]

			var umountCmd *exec.Cmd
			if runtime.GOOS == "darwin" {
				umountCmd = exec.Command("umount", mountpoint)
			} else {
				// Try fusermount3 first, then fusermount, then umount
				if _, err := exec.LookPath("fusermount3"); err == nil {
					umountCmd = exec.Command("fusermount3", "-u", mountpoint)
				} else if _, err := exec.LookPath("fusermount"); err == nil {
					umountCmd = exec.Command("fusermount", "-u", mountpoint)
				} else {
					umountCmd = exec.Command("umount", mountpoint)
				}
			}

			output, err := umountCmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("umount failed: %s", string(output))
			}

			fmt.Fprintf(h.IOStreams.Out, "Unmounted %s\n", mountpoint)
			return nil
		},
	}
}

func isMounted(mountpoint string) bool {
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), mountpoint)
}
