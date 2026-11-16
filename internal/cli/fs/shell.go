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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func shellCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Interactive filesystem shell",
		Long: `Start an interactive shell for filesystem operations.

Commands:
  cd <path>      Change directory
  pwd            Print current directory
  ls [path]      List directory
  cat <path>     Display file contents
  cp <src> <dst> Copy files (remote only)
  mkdir <path>   Create directory
  mv <old> <new> Move/rename
  rm <path>      Remove file
  sql <query>    Execute SQL query
  stat <path>    Show metadata
  help           Show help
  exit           Exit shell

The prompt shows the current remote directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			cwd := "/"
			ctx := context.Background()
			scanner := bufio.NewScanner(h.IOStreams.In)

			fmt.Fprintln(h.IOStreams.Out, "TiDB Cloud FS Shell. Type 'help' for commands, 'exit' to quit.")

			for {
				fmt.Fprint(h.IOStreams.Out, "ticloud:fs> ")
				if !scanner.Scan() {
					break
				}

				line := scanner.Text()
				parts := strings.Fields(line)
				if len(parts) == 0 {
					continue
				}

				cmd := parts[0]
				args := parts[1:]

				switch cmd {
				case "exit", "quit":
					return nil

				case "help":
					fmt.Fprintln(h.IOStreams.Out, `Commands:
  cd <path>      Change directory
  pwd            Print current directory
  ls [path]      List directory
  cat <path>     Display file contents
  cp <src> <dst> Copy files (remote only)
  mkdir <path>   Create directory
  mv <old> <new> Move/rename
  rm <path>      Remove file
  sql <query>    Execute SQL query
  stat <path>    Show metadata
  help           Show this help
  exit           Exit shell`)

				case "pwd":
					fmt.Fprintln(h.IOStreams.Out, cwd)

				case "cd":
					if len(args) == 0 {
						cwd = "/"
					} else {
						cwd = resolvePath(cwd, args[0])
					}

				case "ls":
					path := cwd
					if len(args) > 0 {
						path = resolvePath(cwd, args[0])
					}
					entries, err := client.List(path)
					if err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", suggestInitIfTenantNotFound(fmt.Errorf("list: %w", err)))
						continue
					}
					for _, e := range entries {
						name := e.Name
						if e.IsDir {
							name += "/"
						}
						fmt.Fprintln(h.IOStreams.Out, name)
					}

				case "cat":
					if len(args) < 1 {
						fmt.Fprintln(h.IOStreams.Err, "Usage: cat <path>")
						continue
					}
					path := resolvePath(cwd, args[0])
					reader, err := client.ReadStream(ctx, path)
					if err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", suggestInitIfTenantNotFound(fmt.Errorf("read: %w", err)))
						continue
					}
					_, err = io.Copy(h.IOStreams.Out, reader)
					reader.Close()
					if err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", err)
					}

				case "cp":
					if len(args) < 2 {
						fmt.Fprintln(h.IOStreams.Err, "Usage: cp <src> <dst>")
						continue
					}
					src, dst := resolvePath(cwd, args[0]), resolvePath(cwd, args[1])
					if err := client.Copy(src, dst); err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", suggestInitIfTenantNotFound(fmt.Errorf("copy: %w", err)))
					}

				case "mkdir":
					if len(args) < 1 {
						fmt.Fprintln(h.IOStreams.Err, "Usage: mkdir <path>")
						continue
					}
					path := resolvePath(cwd, args[0])
					if err := client.Mkdir(path); err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", suggestInitIfTenantNotFound(fmt.Errorf("mkdir: %w", err)))
					}

				case "mv":
					if len(args) < 2 {
						fmt.Fprintln(h.IOStreams.Err, "Usage: mv <old> <new>")
						continue
					}
					old, newPath := resolvePath(cwd, args[0]), resolvePath(cwd, args[1])
					if err := client.Rename(old, newPath); err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", suggestInitIfTenantNotFound(fmt.Errorf("mv: %w", err)))
					}

				case "rm":
					if len(args) < 1 {
						fmt.Fprintln(h.IOStreams.Err, "Usage: rm <path>")
						continue
					}
					path := resolvePath(cwd, args[0])
					if err := client.Delete(path); err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", suggestInitIfTenantNotFound(fmt.Errorf("rm: %w", err)))
					}

				case "sql":
					if len(args) < 1 {
						fmt.Fprintln(h.IOStreams.Err, "Usage: sql <query>")
						continue
					}
					query := strings.Join(args, " ")
					rows, err := client.SQL(query)
					if err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", suggestInitIfTenantNotFound(fmt.Errorf("sql: %w", err)))
						continue
					}
					enc := json.NewEncoder(h.IOStreams.Out)
					enc.SetIndent("", "  ")
					if err := enc.Encode(rows); err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", err)
					}

				case "stat":
					if len(args) < 1 {
						fmt.Fprintln(h.IOStreams.Err, "Usage: stat <path>")
						continue
					}
					path := resolvePath(cwd, args[0])
					info, err := client.Stat(path)
					if err != nil {
						fmt.Fprintf(h.IOStreams.Err, "Error: %v\n", suggestInitIfTenantNotFound(fmt.Errorf("stat: %w", err)))
						continue
					}
					fmt.Fprintf(h.IOStreams.Out, "Size: %d, Dir: %v, Rev: %d\n", info.Size, info.IsDir, info.Revision)

				default:
					fmt.Fprintf(h.IOStreams.Err, "Unknown command: %s\n", cmd)
				}
			}

			return nil
		},
	}
}

func resolvePath(cwd, p string) string {
	if strings.HasPrefix(p, "/") {
		return p
	}
	return filepath.Join(cwd, p)
}
