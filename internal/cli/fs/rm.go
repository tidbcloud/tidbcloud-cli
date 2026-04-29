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

func rmCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "rm <path>",
		Short: "Remove a file or directory",
		Example: `  ticloud fs rm :/file.txt
  ticloud fs rm :/empty_folder/`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			path := ParseRemotePath(args[0]).Path

			if err := client.Delete(path); err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("remove %s: %w", path, err))
			}
			return nil
		},
	}
}
