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

func mvCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "mv <old> <new>",
		Short: "Move or rename a file or directory",
		Long: `Move or rename a file or directory in the remote filesystem.

This is a metadata-only operation with zero data transfer cost.

Examples:
  ticloud fs mv :/oldname :/newname
  ticloud fs mv :/folder/file.txt :/other/file.txt`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			oldPath := ParseRemotePath(args[0]).Path
			newPath := ParseRemotePath(args[1]).Path

			if err := client.Rename(oldPath, newPath); err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("move %s to %s: %w", oldPath, newPath, err))
			}
			return nil
		},
	}
}
