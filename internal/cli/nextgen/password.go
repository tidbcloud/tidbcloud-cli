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
	"io"
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
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

const encryptedPasswordPrefix = "rsa_oaep_sha256:"

type rootPasswordPrompter func() (string, error)

func passwordCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	return passwordCmdWithDependencies(h, plan, selectInstance, promptRootPassword)
}

func passwordCmdWithDependencies(h *internal.Helper, plan planSpec, selectInstance instanceSelector, promptPassword rootPasswordPrompter) *cobra.Command {
	var passwordStdin bool
	command := &cobra.Command{
		Use:         "password",
		Short:       fmt.Sprintf("Reset the root password of a %s instance", plan.displayName),
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Reset the root password interactively:
  $ %[1]s %[2]s password
  $ %[1]s %[2]s password -c <instance-id>

  Reset the root password from a protected file (or pipe from a secret manager):
  $ %[1]s %[2]s password -c <instance-id> --password-stdin < /path/to/password-file`, config.CliName, plan.commandName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if passwordStdin {
				id, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return err
				}
				if id == "" {
					return fmt.Errorf("--%s is required when --password-stdin is specified", flag.ClusterID)
				}
			}
			if !passwordStdin && !h.IOStreams.CanPrompt {
				return fmt.Errorf("the terminal doesn't support prompts; use --password-stdin with --%s", flag.ClusterID)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !passwordStdin {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
			}
			var password string
			var err error
			if passwordStdin {
				password, err = readRootPassword(cmd.InOrStdin())
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

			if !passwordStdin {
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
	command.Flags().BoolVar(&passwordStdin, "password-stdin", false, "Read the new root password from a pipe or file on standard input until EOF; terminal input is rejected. Requires --cluster-id; removes one trailing LF or CRLF. If omitted, prompt with hidden input and confirmation.")
	return command
}

func readRootPassword(reader io.Reader) (string, error) {
	if input, ok := reader.(interface{ Fd() uintptr }); ok {
		if isatty.IsTerminal(input.Fd()) || isatty.IsCygwinTerminal(input.Fd()) {
			return "", fmt.Errorf("--password-stdin requires a pipe or file; omit --password-stdin to enter the password with hidden input")
		}
	}
	// Bound input even for the encrypted password format, whose length varies.
	const maxInput = 1 << 20
	data, err := io.ReadAll(io.LimitReader(reader, maxInput+1))
	if err != nil {
		return "", fmt.Errorf("read password from standard input: %w", err)
	}
	if len(data) > maxInput {
		return "", fmt.Errorf("password input exceeds 1 MiB")
	}
	password := string(data)
	if strings.HasSuffix(password, "\n") {
		password = strings.TrimSuffix(strings.TrimSuffix(password, "\n"), "\r")
	}
	return password, nil
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
