// Copyright 2024 PingCAP, Inc.
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

package export

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/util"
)

const confirmed = "yes"

type CancelOpts struct {
	interactive bool
}

func (c CancelOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ExportID,
	}
}

func (c *CancelOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := c.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	// Mark required flags
	if !c.interactive {
		for _, fn := range flags {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CancelCmd(h *internal.Helper) *cobra.Command {
	opts := CancelOpts{
		interactive: true,
	}

	var force bool
	var cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel an export task",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Cancel an export in interactive mode:
  $ %[1]s serverless export cancel

  Cancel an export in non-interactive mode:
  $ %[1]s serverless export cancel -c <cluster-id> -e <export-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID string
			var exportID string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				export, err := cloud.GetSelectedExport(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				exportID = export.ID
			} else {
				// non-interactive mode, get values from flags
				exportID, err = cmd.Flags().GetString(flag.ExportID)
				if err != nil {
					return errors.Trace(err)
				}

				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to cancel the export")
				}

				confirmationMessage := fmt.Sprintf("%s %s %s", color.BlueString("Please type"), color.HiBlueString(confirmed), color.BlueString("to confirm:"))

				prompt := &survey.Input{
					Message: confirmationMessage,
				}

				var userInput string
				err := survey.AskOne(prompt, &userInput)
				if err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					} else {
						return err
					}
				}

				if userInput != confirmed {
					return errors.New("incorrect confirm string entered, skipping export cancellation")
				}
			}

			_, err = d.CancelExport(ctx, clusterID, exportID)
			if err != nil {
				return errors.Trace(err)
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("export %s canceled", exportID))
			if err != nil {
				return err
			}
			return nil
		},
	}

	cancelCmd.Flags().BoolVar(&force, flag.Force, false, "Cancel an export without confirmation.")
	cancelCmd.Flags().StringP(flag.ExportID, flag.ExportIDShort, "", "The ID of the export to be canceled.")
	cancelCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID of the export to be canceled.")

	return cancelCmd
}
