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

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type UpdateOpts struct {
	interactive bool
}

var updateAuthorizedNetworkField = map[string]int{
	flag.NewDisplayName:    0,
	flag.NewStartIPAddress: 1,
	flag.NewEndIPAddress:   2,
}

func (c UpdateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.NewDisplayName,
		flag.NewStartIPAddress,
		flag.NewEndIPAddress,
		flag.StartIPAddress,
		flag.EndIPAddress,
	}
}

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateOpts{
		interactive: true,
	}

	var UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an authorized network",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Update an authorized network in interactive mode:
  $ %[1]s serverless authorized-network update

  Update an authorized network in non-interactive mode:
  $ %[1]s serverless authorized-network update -c <cluster-id> --start-ip-address <start-ip-address> --end-ip-address <end-ip-address> --new-start-ip-address <new-start-ip-address> --new-end-ip-address <new-end-ip-address>`,
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
				err := cmd.MarkFlagRequired(flag.ClusterID)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.StartIPAddress)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.EndIPAddress)
				if err != nil {
					return err
				}
				cmd.MarkFlagsOneRequired(flag.NewStartIPAddress, flag.NewEndIPAddress, flag.NewDisplayName)
				cmd.MarkFlagsRequiredTogether(flag.NewStartIPAddress, flag.NewEndIPAddress)
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
			var startIPAddress string
			var endIPAddress string
			var newStartIPAddress string
			var newEndIPAddress string
			var newDisplayName string
			if opts.interactive {
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

				startIPAddress, endIPAddress, err = cloud.GetSelectedAuthorizedNetwork(ctx, clusterID, d)
				if err != nil {
					return err
				}

				// variables for input
				p := tea.NewProgram(initialUpdateInputModel())
				inputModel, err := p.Run()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				newDisplayName = inputModel.(ui.TextInputModel).Inputs[updateAuthorizedNetworkField[flag.NewDisplayName]].Value()
				newStartIPAddress = inputModel.(ui.TextInputModel).Inputs[updateAuthorizedNetworkField[flag.NewStartIPAddress]].Value()
				newEndIPAddress = inputModel.(ui.TextInputModel).Inputs[updateAuthorizedNetworkField[flag.NewEndIPAddress]].Value()

				if (newStartIPAddress == "" && newEndIPAddress != "") || (newStartIPAddress != "" && newEndIPAddress == "") {
					return errors.New("both new start IP address and new end IP address must be provided")
				}
				if newStartIPAddress == "" && newDisplayName == "" {
					return errors.New("at least one of new display name, new start IP address and new end IP address must be provided")
				}
			} else {
				// non-interactive mode doesn't need projectID
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}

				newDisplayName, err = cmd.Flags().GetString(flag.NewDisplayName)
				if err != nil {
					return errors.Trace(err)
				}

				startIPAddress, err = cmd.Flags().GetString(flag.StartIPAddress)
				if err != nil {
					return errors.Trace(err)
				}

				endIPAddress, err = cmd.Flags().GetString(flag.EndIPAddress)
				if err != nil {
					return errors.Trace(err)
				}

				newStartIPAddress, err = cmd.Flags().GetString(flag.NewStartIPAddress)
				if err != nil {
					return errors.Trace(err)
				}

				newEndIPAddress, err = cmd.Flags().GetString(flag.NewEndIPAddress)
				if err != nil {
					return errors.Trace(err)
				}
			}

			existedAuthorizedNetworks, err := cloud.RetrieveAuthorizedNetworks(ctx, clusterID, d)
			if err != nil {
				return errors.Trace(err)
			}

			findTarget := false
			for i, v := range existedAuthorizedNetworks {
				if v.StartIpAddress == startIPAddress && v.EndIpAddress == endIPAddress {
					findTarget = true
					existedAuthorizedNetworks = slices.Delete(existedAuthorizedNetworks, i, i+1)
					if newDisplayName == "" {
						newDisplayName = v.DisplayName
					}
					if newStartIPAddress == "" {
						newStartIPAddress = v.StartIpAddress
						newEndIPAddress = v.EndIpAddress
					}
					break
				}
			}

			if !findTarget {
				return errors.New(fmt.Sprintf("authorized network %s-%s not found", startIPAddress, endIPAddress))
			}

			newAuthorizedNetwork, err := util.ConvertToAuthorizedNetwork(newStartIPAddress, newEndIPAddress, newDisplayName)
			if err != nil {
				return errors.Trace(err)
			}

			authorizedNetworks := append(existedAuthorizedNetworks, newAuthorizedNetwork)

			body := &cluster.V1beta1ServerlessServicePartialUpdateClusterBody{
				Cluster: &cluster.RequiredTheClusterToBeUpdated{
					Endpoints: &cluster.V1beta1ClusterEndpoints{
						Public: &cluster.EndpointsPublic{
							AuthorizedNetworks: authorizedNetworks,
						},
					},
				},
			}
			body.UpdateMask = AuthorizedNetworkMask

			_, err = d.PartialUpdateCluster(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("authorized network is updated"))
			if err != nil {
				return err
			}
			return nil
		},
	}

	UpdateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	UpdateCmd.Flags().StringP(flag.StartIPAddress, "", "", "The start IP address of the authorized network.")
	UpdateCmd.Flags().StringP(flag.EndIPAddress, "", "", "The end IP address of the authorized network.")
	UpdateCmd.Flags().StringP(flag.NewDisplayName, "", "", "The new display name of the authorized network.")
	UpdateCmd.Flags().StringP(flag.NewStartIPAddress, "", "", "The new start IP address of the authorized network.")
	UpdateCmd.Flags().StringP(flag.NewEndIPAddress, "", "", "The new end IP address of the authorized network.")

	return UpdateCmd
}

func initialUpdateInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(updateAuthorizedNetworkField)),
	}

	for k, v := range updateAuthorizedNetworkField {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 32

		switch k {
		case flag.NewDisplayName:
			t.Placeholder = "New Display Name (optional)"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.NewStartIPAddress:
			t.Placeholder = "New Start IP Address(optional, e.g. 0.0.0.0)"
		case flag.NewEndIPAddress:
			t.Placeholder = "New End IP Address(optional, e.g. 255.255.255.255)"
		}
		m.Inputs[v] = t
	}
	return m
}
