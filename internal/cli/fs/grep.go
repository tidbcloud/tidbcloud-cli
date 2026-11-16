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

func grepCmd(h *internal.Helper) *cobra.Command {
	var limit int
	var cmd = &cobra.Command{
		Use:   "grep <pattern> [path]",
		Short: "Search file contents",
		Long: `Search for a pattern in file contents across the filesystem.

Results are ranked by relevance score. Only text files are searched.

Examples:
  ticloud fs grep "function main" :/
  ticloud fs grep "TODO" :/myproject --limit 20`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			pattern := args[0]
			path := "/"
			if len(args) > 1 {
				path = ParseRemotePath(args[1]).Path
			}

			results, err := client.Grep(pattern, path, limit)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("grep: %w", err))
			}

			for _, r := range results {
				if r.Score != nil {
					fmt.Printf("%.3f %s\n", *r.Score, r.Path)
				} else {
					fmt.Println(r.Path)
				}
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 0, "Maximum number of results")
	return cmd
}
