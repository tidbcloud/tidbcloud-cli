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

package serverless

import (
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

const confirmed = "yes"

type DeleteOpts struct {
	interactive bool
}

func (c DeleteOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func DeleteCmd(h *internal.Helper) *cobra.Command {
	opts := DeleteOpts{
		interactive: true,
	}

	var force bool
	var deleteCmd = &cobra.Command{
		Use:         "delete",
		Short:       "Delete a TiDB Cloud Serverless cluster",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Delete a TiDB Cloud Serverless cluster in interactive mode:
 $ %[1]s serverless delete

 Delete a TiDB Cloud Serverless cluster in non-interactive mode:
 $ %[1]s serverless delete -c <cluster-id>`, config.CliName),
		Aliases: []string{"rm"},
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
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID string
			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				projectID := project.ID

				c, err := cloud.GetSelectedCluster(ctx, projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = c.ID
			} else {
				// non-interactive mode doesn't need projectID
				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID = cID
			}

			if !force {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support prompt, please run with --force to delete the cluster")
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
					return errors.New("Incorrect confirm string entered, skipping database deletion")
				}
			}

			c, err := d.DeleteCluster(ctx, clusterID)
			if err != nil {
				return errors.Trace(err)
			}
			if *c.State == cluster.COMMONV1BETA1CLUSTERSTATE_DELETING || *c.State == cluster.COMMONV1BETA1CLUSTERSTATE_DELETED {
				fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("cluster %s deleted", clusterID)))
				return nil
			} else {
				return errors.New(fmt.Sprintf("Delete cluster %s failed, please check status on dashboard", clusterID))
			}
		},
	}

	deleteCmd.Flags().BoolVar(&force, flag.Force, false, "Delete a cluster without confirmation.")
	deleteCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster to be deleted.")
	return deleteCmd
}
