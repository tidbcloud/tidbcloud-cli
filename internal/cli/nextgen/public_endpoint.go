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
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func publicEndpointCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	command := &cobra.Command{
		Use:   "public-endpoint",
		Short: fmt.Sprintf("Manage the public endpoint of a %s instance", plan.displayName),
	}
	command.AddCommand(publicEndpointActionCmd(h, plan, true, selectInstance, confirmDisablePublicEndpoint))
	command.AddCommand(publicEndpointActionCmd(h, plan, false, selectInstance, confirmDisablePublicEndpoint))
	return command
}

func publicEndpointActionCmd(h *internal.Helper, plan planSpec, enabled bool, selectInstance instanceSelector, confirmDisable func(string) error) *cobra.Command {
	action := "enable"
	actionTitle := "Enable"
	if !enabled {
		action = "disable"
		actionTitle = "Disable"
	}
	var force bool
	command := &cobra.Command{
		Use:   action,
		Short: fmt.Sprintf("%s the public endpoint of a %s instance", actionTitle, plan.displayName),
		Long: fmt.Sprintf("%s the public endpoint of a %s instance, preserving its existing IP access list. No IP access list entries are added. "+
			"The request is asynchronous; use --describe to check endpoint readiness after it is accepted.", actionTitle, plan.displayName),
		Args: cobra.NoArgs,
		Example: fmt.Sprintf(`  $ %[1]s %[2]s public-endpoint %[3]s
  $ %[1]s %[2]s public-endpoint %[3]s -c <instance-id>`, config.CliName, plan.commandName, action),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := validatedOutputFormat(cmd)
			if err != nil {
				return err
			}
			id, err := cmd.Flags().GetString(flag.ClusterID)
			if err != nil {
				return err
			}
			if id == "" && cmd.Flags().Changed(flag.ClusterID) {
				return fmt.Errorf("--%s must not be empty", flag.ClusterID)
			}
			if id == "" && !h.IOStreams.CanPrompt {
				return fmt.Errorf("the terminal doesn't support prompts; use --%s to specify the instance", flag.ClusterID)
			}
			if !enabled && !force && !h.IOStreams.CanPrompt {
				return fmt.Errorf("the terminal doesn't support prompts; use --%s to disable the public endpoint", flag.Force)
			}
			client, err := h.NextGenClient()
			if err != nil {
				return err
			}
			if id == "" {
				selected, err := selectInstance(cmd.Context(), client, plan, int32(h.QueryPageSize))
				if err != nil {
					return err
				}
				id = selected.GetTidbId()
			}
			instance, err := client.GetTiDB(cmd.Context(), id)
			if err != nil {
				return err
			}
			if err := plan.validate(instance); err != nil {
				return err
			}

			if !enabled && !force {
				if err := confirmDisable(id); err != nil {
					return err
				}
			}

			body := api.NewV1beta2PublicConnectionSetting()
			body.SetEnabled(enabled)
			setting, err := client.UpdatePublicConnectionSetting(cmd.Context(), id, body)
			if err != nil {
				return err
			}
			if setting == nil {
				return fmt.Errorf("the API returned an empty public connection setting")
			}
			if format == output.JsonFormat {
				return output.PrintJson(h.IOStreams.Out, setting)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("%s request accepted for the public endpoint of instance %s", actionTitle, id))
			return nil
		},
	}
	command.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the instance. If omitted, select an instance interactively.")
	command.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	if !enabled {
		command.Long += " Disabling the public endpoint interrupts public connections."
		command.Flags().BoolVar(&force, flag.Force, false, "Disable the public endpoint without confirmation.")
		command.Example += fmt.Sprintf("\n  $ %s %s public-endpoint disable -c <instance-id> --force", config.CliName, plan.commandName)
	}
	return command
}

func confirmDisablePublicEndpoint(id string) error {
	var answer string
	prompt := &survey.Input{Message: fmt.Sprintf("Disabling the public endpoint of instance %s interrupts public connections. Type %s to confirm:", id, color.HiBlueString("yes"))}
	if err := survey.AskOne(prompt, &answer); err != nil {
		if err == terminal.InterruptErr {
			return util.InterruptError
		}
		return err
	}
	if answer != "yes" {
		return fmt.Errorf("incorrect confirmation string; skipping public endpoint disablement")
	}
	return nil
}
