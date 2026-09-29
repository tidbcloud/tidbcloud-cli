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
	"strconv"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func updateCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: fmt.Sprintf("Update a %s instance", plan.displayName),
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed(flag.DisplayName) && !cmd.Flags().Changed(flag.MaxRCU) {
				return fmt.Errorf("at least one of --%s or --%s must be specified", flag.DisplayName, flag.MaxRCU)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := validatedOutputFormat(cmd)
			if err != nil {
				return err
			}

			id, err := cmd.Flags().GetString(flag.ClusterID)
			if err != nil {
				return err
			}
			body := api.NewTheTiDBCloudPremiumInstanceToUpdate()
			if cmd.Flags().Changed(flag.DisplayName) {
				displayName, err := cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return err
				}
				if err := validateDisplayName(displayName); err != nil {
					return err
				}
				body.DisplayName = &displayName
			}
			if cmd.Flags().Changed(flag.MaxRCU) {
				maxRCU, err := cmd.Flags().GetInt64(flag.MaxRCU)
				if err != nil {
					return err
				}
				if maxRCU <= 0 {
					return fmt.Errorf("--%s must be greater than zero", flag.MaxRCU)
				}
				value := strconv.FormatInt(maxRCU, 10)
				body.MaxRcu = &value
			}

			client, err := h.NextGenClient()
			if err != nil {
				return err
			}
			current, err := client.GetTiDB(cmd.Context(), id)
			if err != nil {
				return err
			}
			if err := plan.validate(current); err != nil {
				return err
			}

			instance, err := client.UpdateTiDB(cmd.Context(), id, body)
			if err != nil {
				return err
			}
			if err := plan.validate(instance); err != nil {
				return err
			}
			if format == output.JsonFormat {
				return output.PrintJson(h.IOStreams.Out, instance)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Update request accepted for instance %s", id))
			return nil
		},
	}
	command.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the instance.")
	command.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The new display name of the instance.")
	command.Flags().Int64(flag.MaxRCU, 0, "The new maximum number of Request Capacity Units (RCUs).")
	command.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	_ = command.MarkFlagRequired(flag.ClusterID)
	return command
}
