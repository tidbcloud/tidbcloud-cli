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
	"context"
	"fmt"
	"io"

	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func catCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "cat <path>",
		Short: "Display file contents",
		Example: `  ticloud fs cat :/file.txt
  ticloud fs cat :/myfolder/data.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			rp := ParseRemotePath(args[0])
			path := rp.Path

			ctx := context.Background()
			reader, err := client.ReadStream(ctx, path)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("read %s: %w", path, err))
			}
			defer reader.Close()

			_, err = io.Copy(h.IOStreams.Out, reader)
			return suggestInitIfTenantNotFound(err)
		},
	}
}
