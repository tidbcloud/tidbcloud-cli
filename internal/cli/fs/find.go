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
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func findCmd(h *internal.Helper) *cobra.Command {
	var (
		nameGlob   string
		tagFilter  string
		newer      string
		older      string
		sizeFilter string
	)

	var cmd = &cobra.Command{
		Use:   "find [path]",
		Short: "Find files by attributes",
		Long: `Find files matching specified criteria.

Filters:
  --name <glob>       Match filename pattern
  --tag <key=value>   Match tag
  --newer <date>      Modified after date (YYYY-MM-DD)
  --older <date>      Modified before date (YYYY-MM-DD)
  --size <+N|-N>      Size filter (+N = larger than, -N = smaller than)

Examples:
  ticloud fs find :/ --name "*.go"
  ticloud fs find :/ --newer 2024-01-01 --size +1M
  ticloud fs find :/myfolder --tag env=production`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			path := "/"
			if len(args) > 0 {
				path = ParseRemotePath(args[0]).Path
			}

			params := url.Values{}

			if nameGlob != "" {
				params.Set("name", nameGlob)
			}

			if tagFilter != "" {
				parts := strings.SplitN(tagFilter, "=", 2)
				if len(parts) == 2 {
					params.Set("tag_key", parts[0])
					params.Set("tag_value", parts[1])
				}
			}

			if newer != "" {
				t, err := time.Parse("2006-01-02", newer)
				if err != nil {
					return fmt.Errorf("invalid -newer date: %w", err)
				}
				params.Set("modified_after", strconv.FormatInt(t.Unix(), 10))
			}

			if older != "" {
				t, err := time.Parse("2006-01-02", older)
				if err != nil {
					return fmt.Errorf("invalid -older date: %w", err)
				}
				params.Set("modified_before", strconv.FormatInt(t.Unix(), 10))
			}

			if sizeFilter != "" {
				if strings.HasPrefix(sizeFilter, "+") {
					sizeStr := sizeFilter[1:]
					size, err := humanize.ParseBytes(sizeStr)
					if err != nil {
						return fmt.Errorf("invalid -size value: %w", err)
					}
					params.Set("min_size", strconv.FormatUint(size, 10))
				} else if strings.HasPrefix(sizeFilter, "-") {
					sizeStr := sizeFilter[1:]
					size, err := humanize.ParseBytes(sizeStr)
					if err != nil {
						return fmt.Errorf("invalid -size value: %w", err)
					}
					params.Set("max_size", strconv.FormatUint(size, 10))
				}
			}

			results, err := client.Find(path, params)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("find: %w", err))
			}

			for _, r := range results {
				size := humanize.IBytes(uint64(r.SizeBytes))
				fmt.Printf("%s\t%s\n", size, r.Path)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&nameGlob, "name", "", "Match filename pattern")
	cmd.Flags().StringVar(&tagFilter, "tag", "", "Match tag (key=value)")
	cmd.Flags().StringVar(&newer, "newer", "", "Modified after date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&older, "older", "", "Modified before date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&sizeFilter, "size", "", "Size filter (+N or -N, e.g., +1M, -100K)")

	return cmd
}
