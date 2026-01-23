// Copyright 2026 PingCAP, Inc.

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

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

const (
	AuthorizedNetworkMask = "endpoints.public.authorizedNetworks"
)

type CreateOpts struct {
	interactive bool
}

var createAuthorizedNetworkField = map[string]int{
	flag.DisplayName:    0,
	flag.StartIPAddress: 1,
	flag.EndIPAddress:   2,
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.StartIPAddress,
		flag.EndIPAddress,
		flag.DisplayName,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.StartIPAddress,
		flag.EndIPAddress,
	}
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var CreateCmd = &cobra.Command{
		Use:         "create",
		Short:       "Create an authorized network",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Create an authorized network in interactive mode:
  $ %[1]s serverless authorized-network create

  Create an authorized network in non-interactive mode:
  $ %[1]s serverless authorized-network create -c <cluster-id> --display-name <display-name> --start-ip-address <start-ip-address> --end-ip-address <end-ip-address>`,
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
			var startIPAddress string
			var endIPAddress string
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

				p := tea.NewProgram(initialCreateInputModel())
				inputModel, err := p.Run()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				displayName = inputModel.(ui.TextInputModel).Inputs[createAuthorizedNetworkField[flag.DisplayName]].Value()
				startIPAddress = inputModel.(ui.TextInputModel).Inputs[createAuthorizedNetworkField[flag.StartIPAddress]].Value()
				endIPAddress = inputModel.(ui.TextInputModel).Inputs[createAuthorizedNetworkField[flag.EndIPAddress]].Value()

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

				startIPAddress, err = cmd.Flags().GetString(flag.StartIPAddress)
				if err != nil {
					return errors.Trace(err)
				}

				endIPAddress, err = cmd.Flags().GetString(flag.EndIPAddress)
				if err != nil {
					return errors.Trace(err)
				}
			}

			authorizedNetwork, err := util.ConvertToAuthorizedNetwork(startIPAddress, endIPAddress, displayName)
			if err != nil {
				return errors.Trace(err)
			}

			existedAuthorizedNetworks, err := cloud.RetrieveAuthorizedNetworks(ctx, clusterID, d)
			if err != nil {
				return errors.Trace(err)
			}

			authorizedNetworks := append(existedAuthorizedNetworks, authorizedNetwork)

			body := &cluster.V1beta1ClusterServicePartialUpdateClusterBody{
				Cluster: &cluster.V1beta1ClusterServicePartialUpdateClusterBodyCluster{
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

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("authorized network %s-%s is created", startIPAddress, endIPAddress))
			if err != nil {
				return err
			}
			return nil
		},
	}

	CreateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	CreateCmd.Flags().StringP(flag.StartIPAddress, "", "", "The start IP address of the authorized network.")
	CreateCmd.Flags().StringP(flag.EndIPAddress, "", "", "The end IP address of the authorized network.")
	CreateCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The name of the authorized network.")

	return CreateCmd
}

func initialCreateInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(createAuthorizedNetworkField)),
	}

	for k, v := range createAuthorizedNetworkField {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 32

		switch k {
		case flag.DisplayName:
			t.Placeholder = "Display Name (optional)"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.StartIPAddress:
			t.Placeholder = "Start IP Address (e.g. 0.0.0.0)"
		case flag.EndIPAddress:
			t.Placeholder = "End IP Address (e.g. 255.255.255.255)"
		}
		m.Inputs[v] = t
	}
	return m
}
