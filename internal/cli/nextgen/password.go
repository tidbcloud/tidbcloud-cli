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
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const encryptedPasswordPrefix = "rsa_oaep_sha256:"

type rootPasswordPrompter func() (string, error)

func passwordCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	return passwordCmdWithDependencies(h, plan, selectInstance, promptRootPassword)
}

func passwordCmdWithDependencies(h *internal.Helper, plan planSpec, selectInstance instanceSelector, promptPassword rootPasswordPrompter) *cobra.Command {
	command := &cobra.Command{
		Use:         "password",
		Short:       fmt.Sprintf("Reset the root password of a %s instance", plan.displayName),
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Reset the root password interactively:
  $ %[1]s %[2]s password
  $ %[1]s %[2]s password -c <instance-id>

  Reset the root password non-interactively:
  $ %[1]s %[2]s password -c <instance-id> --password <password>`, config.CliName, plan.commandName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed(flag.Password) {
				id, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return err
				}
				if id == "" {
					return fmt.Errorf("--%s is required when --%s is specified", flag.ClusterID, flag.Password)
				}
			}
			if !cmd.Flags().Changed(flag.Password) && !h.IOStreams.CanPrompt {
				return fmt.Errorf("the terminal doesn't support prompts; use --%s with --%s", flag.Password, flag.ClusterID)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed(flag.Password) {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
			}
			var password string
			var err error
			if cmd.Flags().Changed(flag.Password) {
				password, err = cmd.Flags().GetString(flag.Password)
				if err != nil {
					return err
				}
				if err := validateRootPassword(password); err != nil {
					return err
				}
			}

			client, err := h.NextGenClient()
			if err != nil {
				return err
			}
			id, err := cmd.Flags().GetString(flag.ClusterID)
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

			if !cmd.Flags().Changed(flag.Password) {
				password, err = promptPassword()
				if err != nil {
					return err
				}
				if err := validateRootPassword(password); err != nil {
					return err
				}
			}

			body := api.NewTidbServiceResetRootPasswordBody(password)
			if err := client.ResetRootPassword(cmd.Context(), id, body); err != nil {
				return err
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Root password reset for instance %s", id))
			return nil
		},
	}
	command.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the instance. If omitted, select an instance interactively.")
	command.Flags().String(flag.Password, "", "The new root password. Prefer interactive input to avoid shell history exposure.")
	return command
}

func promptRootPassword() (string, error) {
	var password string
	if err := survey.AskOne(&survey.Password{Message: "Enter the new root password:"}, &password); err != nil {
		if err == terminal.InterruptErr {
			return "", util.InterruptError
		}
		return "", err
	}
	var confirmation string
	if err := survey.AskOne(&survey.Password{Message: "Confirm the new root password:"}, &confirmation); err != nil {
		if err == terminal.InterruptErr {
			return "", util.InterruptError
		}
		return "", err
	}
	if password != confirmation {
		return "", fmt.Errorf("password confirmation does not match")
	}
	return password, nil
}

func validateRootPassword(password string) error {
	if strings.HasPrefix(password, encryptedPasswordPrefix) {
		if len(password) == len(encryptedPasswordPrefix) {
			return fmt.Errorf("encrypted root password payload must not be empty")
		}
		return nil
	}
	// Match the byte-length limit enforced by mgmt and Global.
	length := len(password)
	if length < 8 || length > 64 {
		return fmt.Errorf("root password must be between 8 and 64 bytes")
	}
	return nil
}
