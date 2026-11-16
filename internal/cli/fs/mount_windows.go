//go:build windows

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

	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func mountCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "mount <mountpoint>",
		Short: "Mount remote filesystem via FUSE",
		Long:  `Mount the remote filesystem as a local FUSE mount.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("FUSE mount is not supported on Windows")
		},
	}
}

func umountCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "umount <mountpoint>",
		Short: "Unmount FUSE filesystem",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("FUSE umount is not supported on Windows")
		},
	}
}
