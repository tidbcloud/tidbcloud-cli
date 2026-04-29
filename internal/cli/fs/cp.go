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
	"os"
	"path/filepath"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func cpCmd(h *internal.Helper) *cobra.Command {
	var resume bool
	var cmd = &cobra.Command{
		Use:   "cp <src> <dst>",
		Short: "Copy files between local and remote",
		Long: `Copy files between local and remote filesystems.

Source and destination can be:
  - Local paths (any path not starting with :)
  - Remote paths (starting with :/)
  - - for stdin (source) or stdout (destination)

Examples:
  ticloud fs cp local.txt :/remote/           # Upload
  ticloud fs cp :/remote/file.txt .           # Download
  ticloud fs cp :/a/b :/c/d                   # Server-side copy
  cat data.txt | ticloud fs cp - :/remote/    # Upload from stdin`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return suggestInitIfTenantNotFound(err)
			}

			src, dst := args[0], args[1]
			srcRemote := IsRemote(src)
			dstRemote := IsRemote(dst)

			ctx := context.Background()

			switch {
			case !srcRemote && !dstRemote:
				return fmt.Errorf("local-to-local copy not supported")

			case !srcRemote && dstRemote:
				// Upload
				dstPath := ParseRemotePath(dst).Path
				if src == "-" {
					data, err := io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("read stdin: %w", err)
					}
					return client.Write(dstPath, data)
				}
				if strings.HasSuffix(dstPath, "/") {
					dstPath = filepath.Join(dstPath, filepath.Base(src))
				}
				f, err := os.Open(src)
				if err != nil {
					return suggestInitIfTenantNotFound(fmt.Errorf("open %s: %w", src, err))
				}
				defer f.Close()
				stat, err := f.Stat()
				if err != nil {
					return suggestInitIfTenantNotFound(fmt.Errorf("stat %s: %w", src, err))
				}
				if resume {
					return suggestInitIfTenantNotFound(client.ResumeUpload(ctx, dstPath, f, stat.Size(), nil))
				}
				return suggestInitIfTenantNotFound(client.WriteStream(ctx, dstPath, f, stat.Size(), nil))

			case srcRemote && !dstRemote:
				// Download
				srcPath := ParseRemotePath(src).Path
				reader, err := client.ReadStream(ctx, srcPath)
				if err != nil {
					return suggestInitIfTenantNotFound(fmt.Errorf("read %s: %w", src, err))
				}
				defer reader.Close()

				if dst == "-" {
					_, err = io.Copy(h.IOStreams.Out, reader)
					return suggestInitIfTenantNotFound(err)
				}
				f, err := os.Create(dst)
				if err != nil {
					return suggestInitIfTenantNotFound(fmt.Errorf("create %s: %w", dst, err))
				}
				defer f.Close()
				_, err = io.Copy(f, reader)
				return suggestInitIfTenantNotFound(err)

			case srcRemote && dstRemote:
				// Server-side copy
				srcPath := ParseRemotePath(src).Path
				dstPath := ParseRemotePath(dst).Path
				if strings.HasSuffix(dstPath, "/") {
					dstPath = filepath.Join(dstPath, filepath.Base(srcPath))
				}
				return suggestInitIfTenantNotFound(client.Copy(srcPath, dstPath))
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&resume, "resume", false, "Resume interrupted upload")
	return cmd
}
