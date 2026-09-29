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
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func deleteCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	var force bool
	command := &cobra.Command{
		Use:   "delete",
		Short: fmt.Sprintf("Delete a %s instance", plan.displayName),
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

			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompts; use --%s to delete the instance", flag.Force)
				}
				var answer string
				prompt := &survey.Input{Message: fmt.Sprintf("Type %s to confirm deletion:", color.HiBlueString("yes"))}
				if err := survey.AskOne(prompt, &answer); err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					}
					return err
				}
				if answer != "yes" {
					return fmt.Errorf("incorrect confirmation string; skipping instance deletion")
				}
			}

			if _, err := client.DeleteTiDB(cmd.Context(), id); err != nil {
				return err
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Delete request accepted for instance %s", id))
			return nil
		},
	}
	command.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the instance to delete.")
	command.Flags().BoolVar(&force, flag.Force, false, "Delete the instance without confirmation.")
	_ = command.MarkFlagRequired(flag.ClusterID)
	return command
}
