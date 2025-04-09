// Copyright 2025 PingCAP, Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authorizednetwork

import (
	"fmt"
	"slices"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

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
		flag.IPRange,
	}
}

func (c DeleteOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.IPRange,
	}
}

func DeleteCmd(h *internal.Helper) *cobra.Command {
	opts := DeleteOpts{
		interactive: true,
	}

	var force bool
	var DeleteCmd = &cobra.Command{
		Use:         "delete",
		Short:       "Delete an authorized network",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Delete an authorized network in interactive mode:
  $ %[1]s serverless authorized-network delete

  Delete an authorized network in non-interactive mode:
  $ %[1]s serverless authorized-network delete -c <cluster-id> --ip-range <ip-range>`,
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
				for _, fn := range opts.RequiredFlags() {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID string
			var displayName string
			var ipRange string
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

				cluster, err := cloud.GetSelectedCluster(ctx, projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				// variables for input
				fmt.Fprintln(h.IOStreams.Out, color.BlueString("Please input the following options"))

				ipRange, err = cloud.GetSelectedAuthorizedNetwork(ctx, clusterID, d)
				if err != nil {
					return err
				}
			} else {
				// non-interactive mode doesn't need projectID
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}

				ipRange, err = cmd.Flags().GetString(flag.IPRange)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to delete the authorized network")
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
					return errors.New("incorrect confirm string entered, skipping authorized network deletion")
				}
			}

			authorizedNetwork, err := util.ConvertToAuthorizedNetwork(ipRange, displayName)
			if err != nil {
				return errors.Trace(err)
			}

			existedAuthorizedNetworks, err := cloud.RetrieveAuthorizedNetworks(ctx, clusterID, d)
			if err != nil {
				return errors.Trace(err)
			}

			for i, v := range existedAuthorizedNetworks {
				if v.StartIpAddress == authorizedNetwork.StartIpAddress && v.EndIpAddress == authorizedNetwork.EndIpAddress {
					existedAuthorizedNetworks = slices.Delete(existedAuthorizedNetworks, i, i+1)
					break
				}
			}
			body := &cluster.V1beta1ServerlessServicePartialUpdateClusterBody{
				Cluster: &cluster.RequiredTheClusterToBeUpdated{
					Endpoints: &cluster.V1beta1ClusterEndpoints{
						Public: &cluster.EndpointsPublic{
							AuthorizedNetworks: existedAuthorizedNetworks,
						},
					},
				},
			}
			body.UpdateMask = AuthorizedNetworkMask

			_, err = d.PartialUpdateCluster(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("authorized network %s is deleted", ipRange))
			if err != nil {
				return err
			}
			return nil

		},
	}

	DeleteCmd.Flags().BoolVar(&force, flag.Force, false, "Delete an authorized network without confirmation.")
	DeleteCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	DeleteCmd.Flags().StringP(flag.IPRange, "", "", "The IP range of the authorized network.")

	return DeleteCmd
}
