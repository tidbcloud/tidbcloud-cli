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
	flag.DisplayName: 0,
	flag.IPRange:     1,
}

func (c UpdateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.TargetIPRange,
		flag.IPRange,
		flag.DisplayName,
	}
}

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateOpts{
		interactive: true,
	}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an authorized network",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Update an authorized network in interactive mode:
  $ %[1]s serverless authorized-network update

  Update an authorized network in non-interactive mode:
  $ %[1]s serverless authorized-network update -c <cluster-id> --ip-range <ip-range> --display-name <display-name>`, config.CliName),
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
				err = cmd.MarkFlagRequired(flag.TargetIPRange)
				if err != nil {
					return err
				}
				cmd.MarkFlagsOneRequired(flag.IPRange, flag.DisplayName)
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
			var targetIPRange string
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

				targetIPRange, err = cloud.GetSelectedAuthorizedNetwork(ctx, clusterID, d)
				if err != nil {
					return err
				}

				// variables for input
				fmt.Fprintln(h.IOStreams.Out, color.BlueString("Please input the following options"))

				p := tea.NewProgram(initialUpdateInputModel())
				inputModel, err := p.Run()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				displayName = inputModel.(ui.TextInputModel).Inputs[createAuthorizedNetworkField[flag.DisplayName]].Value()
				ipRange = inputModel.(ui.TextInputModel).Inputs[createAuthorizedNetworkField[flag.IPRange]].Value()
			} else {
				// non-interactive mode doesn't need projectID
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}

				displayName, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}

				ipRange, err = cmd.Flags().GetString(flag.IPRange)
				if err != nil {
					return errors.Trace(err)
				}

				targetIPRange, err = cmd.Flags().GetString(flag.TargetIPRange)
				if err != nil {
					return errors.Trace(err)
				}
			}

			authorizedNetwork, err := util.ConvertToAuthorizedNetwork(ipRange, displayName)
			if err != nil {
				return errors.Trace(err)
			}

			targetAuthorizedNetwork, err := util.ConvertToAuthorizedNetwork(targetIPRange, "")
			if err != nil {
				return errors.Trace(err)
			}

			existedAuthorizedNetworks, err := cloud.RetrieveAuthorizedNetworks(ctx, clusterID, d)
			if err != nil {
				return errors.Trace(err)
			}

			for i, v := range existedAuthorizedNetworks {
				if v.StartIpAddress == targetAuthorizedNetwork.StartIpAddress && v.EndIpAddress == targetAuthorizedNetwork.EndIpAddress {
					existedAuthorizedNetworks = slices.Delete(existedAuthorizedNetworks, i, i+1)
					break
				}
			}

			authorizedNetworks := append(existedAuthorizedNetworks, authorizedNetwork)

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

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("authorized network %s is updated", displayName))
			if err != nil {
				return err
			}
			return nil
		},
	}

	updateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	updateCmd.Flags().StringP(flag.IPRange, "", "", "The new IP range of the authorized network.")
	updateCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The name of the authorized network.")
	updateCmd.Flags().StringP(flag.TargetIPRange, "", "", "The IP range of the authorized network to be updated.")

	return updateCmd
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
		case flag.DisplayName:
			t.Placeholder = "Display Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.IPRange:
			ipRangeExample := "0.0.0.0-255.255.255.255"
			t.Placeholder = fmt.Sprintf("IP Range (e.g., %s)", ipRangeExample)
		}
		m.Inputs[v] = t
	}
	return m
}
