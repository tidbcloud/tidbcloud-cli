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
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func sqlCmd(h *internal.Helper) *cobra.Command {
	var query, file string
	var cmd = &cobra.Command{
		Use:   "sql",
		Short: "Execute a SQL query",
		Long: `Execute a SQL query against the FS backend database.

You can provide the query directly with -q or from a file with -f.`,
		Example: `  ticloud fs sql -q "SELECT 1"
  ticloud fs sql -f query.sql`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			if query == "" && file == "" {
				return fmt.Errorf("either -q or -f is required")
			}
			if query != "" && file != "" {
				return fmt.Errorf("-q and -f are mutually exclusive")
			}

			if file != "" {
				data, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("read file: %w", err)
				}
				query = string(data)
			}

			rows, err := client.SQL(query)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("sql: %w", err))
			}

			enc := json.NewEncoder(h.IOStreams.Out)
			enc.SetIndent("", "  ")
			return enc.Encode(rows)
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "SQL query string")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a SQL file")

	return cmd
}
