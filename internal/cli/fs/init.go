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

	"github.com/tidbcloud/tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func initCmd(h *internal.Helper) *cobra.Command {
	var (
		user     string
		password string
	)

	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize FS for the associated database",
		Long: `Provision the FS tenant for the configured database.

This command must be run once before any file operations when using a new
TiDB Cloud cluster or TiDB Zero instance.

Examples:
  ticloud fs init --user admin --password secret`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			if user == "" || password == "" {
				return fmt.Errorf("--user and --password are required")
			}

			result, err := client.Provision(context.Background(), user, password)
			if err != nil {
				return fmt.Errorf("provision failed: %w", err)
			}

			fmt.Fprintf(h.IOStreams.Out, "FS initialized for tenant %s (status: %s)\n", result.TenantID, result.Status)
			return nil
		},
	}

	cmd.Flags().StringVar(&user, "user", "", "FS admin user")
	cmd.Flags().StringVar(&password, "password", "", "FS admin password")
	_ = cmd.MarkFlagRequired("user")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
