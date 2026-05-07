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
	"text/tabwriter"

	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func lsCmd(h *internal.Helper) *cobra.Command {
	var long bool
	var cmd = &cobra.Command{
		Use:   "ls [path]",
		Short: "List directory contents",
		Example: `  ticloud fs ls :/
  ticloud fs ls -l :/myfolder`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			path := "/"
			if len(args) > 0 {
				rp := ParseRemotePath(args[0])
				path = rp.Path
			}
			if path == "" {
				path = "/"
			}

			entries, err := client.List(path)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("list %s: %w", path, err))
			}

			if long {
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				for _, e := range entries {
					typeChar := "-"
					if e.IsDir {
						typeChar = "d"
					}
					fmt.Fprintf(w, "%s\t%s\t%s\n", typeChar, humanize.IBytes(uint64(e.Size)), e.Name)
				}
				return w.Flush()
			}

			for _, e := range entries {
				name := e.Name
				if e.IsDir {
					name += "/"
				}
				fmt.Println(name)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&long, "long", "l", false, "Use long listing format")
	return cmd
}
