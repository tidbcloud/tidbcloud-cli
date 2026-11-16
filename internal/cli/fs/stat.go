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

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func statCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "stat <path>",
		Short: "Display file metadata",
		Example: `  ticloud fs stat :/file.txt
  ticloud fs stat :/myfolder/`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			rp := ParseRemotePath(args[0])
			path := rp.Path

			info, err := client.Stat(path)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("stat %s: %w", path, err))
			}

			fmt.Printf("Path:    %s\n", path)
			fmt.Printf("Size:    %s (%d bytes)\n", humanize.IBytes(uint64(info.Size)), info.Size)
			fmt.Printf("Type:    %s\n", map[bool]string{true: "directory", false: "file"}[info.IsDir])
			fmt.Printf("Version: %d\n", info.Revision)
			return nil
		},
	}
}
