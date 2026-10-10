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

package nextgen

import (
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"

	"github.com/spf13/cobra"
)

func describeCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	command := &cobra.Command{
		Use:   "describe",
		Short: fmt.Sprintf("Describe a %s instance", plan.displayName),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cmd.Flags().GetString(flag.ClusterID)
			if err != nil {
				return err
			}
			client, err := h.NextGenClient()
			if err != nil {
				return err
			}
			instance, err := client.GetTiDB(cmd.Context(), id)
			if err != nil {
				return err
			}
			if err := plan.validate(instance); err != nil {
				return err
			}
			return output.PrintJson(h.IOStreams.Out, instance)
		},
	}
	command.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the instance.")
	_ = command.MarkFlagRequired(flag.ClusterID)
	return command
}
