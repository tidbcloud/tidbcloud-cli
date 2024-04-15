// Copyright 2022 PingCAP, Inc.
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

package dataimport

import (
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/util"
	importOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type CancelOpts struct {
	interactive bool
}

func (c CancelOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ImportID,
	}
}

func CancelCmd(h *internal.Helper) *cobra.Command {
	var force bool
	opts := CancelOpts{
		interactive: true,
	}

	var cancelCmd = &cobra.Command{
		Use:         "cancel",
		Short:       "Cancel a data import task",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Cancel an import task in interactive mode:
  $ %[1]s serverless import cancel

  Cancel an import task in non-interactive mode:
  $ %[1]s serverless import cancel --cluster-id <cluster-id> --import-id <import-id>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// mark required flags in non-interactive mode
			if !opts.interactive {
				for _, fn := range flags {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var clusterID, importID string
			d, err := h.Client()
			if err != nil {
				return err
			}

			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}

				cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				// Only task status is pending or importing can be canceled.
				selectedImport, err := cloud.GetSelectedImport(cmd.Context(), clusterID, h.QueryPageSize, d, []importModel.V1beta1ImportState{
					importModel.V1beta1ImportStatePREPARING,
					importModel.V1beta1ImportStateIMPORTING,
				})
				if err != nil {
					return err
				}
				importID = selectedImport.ID
			} else {
				// non-interactive mode
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				importID = cmd.Flag(flag.ImportID).Value.String()
			}

			cmd.Annotations[telemetry.ClusterID] = clusterID

			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to cancel the import task")
				}

				confirmationMessage := fmt.Sprintf("%s %s %s", color.BlueString("Please type"), color.HiBlueString(config.Confirmed), color.BlueString("to confirm:"))

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

				if userInput != config.Confirmed {
					return errors.New("incorrect confirm string entered, skipping import cancellation")
				}
			}

			params := importOp.NewImportServiceCancelImportParams().WithClusterID(clusterID).WithID(importID).WithContext(cmd.Context())
			_, err = d.CancelImport(params)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s has been canceled.", importID))
			return nil
		},
	}

	cancelCmd.Flags().BoolVar(&force, flag.Force, false, "Cancel an import task without confirmation.")
	cancelCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	cancelCmd.Flags().String(flag.ImportID, "", "The ID of import task.")
	return cancelCmd
}
