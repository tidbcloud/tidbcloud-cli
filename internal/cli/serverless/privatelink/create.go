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

package privatelink

import (
	"fmt"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	pl "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type CreateOpts struct {
	interactive bool
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
		flag.PrivateLinkConnectionType,
		flag.EndpointServiceName,
		flag.EndpointServiceRegion,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
		flag.PrivateLinkConnectionType,
		flag.EndpointServiceName,
	}
}

func (c *CreateOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := c.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	if !c.interactive {
		for _, fn := range c.RequiredFlags() {
			if err := cmd.MarkFlagRequired(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a private link connection",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a private link connection in interactive mode:
  $ %[1]s serverless private-link-connection create

  Create a private link connection in non-interactive mode:
  $ %[1]s serverless private-link-connection create --cluster-id <cluster-id> --display-name <display-name> --type <type> --endpoint-service-name <service-name>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID string
			var displayName string
			var endpointServiceName string
			var endpointServiceRegion string
			var connectionTypeRaw string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}

				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				connectionTypeRaw, err = cloud.GetSelectedField([]string{
					string(pl.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE),
					string(pl.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE),
				}, "Choose the private link connection type:")
				if err != nil {
					return err
				}

				connectionType, err := normalizePrivateLinkConnectionType(connectionTypeRaw)
				if err != nil {
					return err
				}
				includeRegion := connectionType == pl.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE

				inputModel, err := getCreateInputModel(includeRegion)
				if err != nil {
					return err
				}
				displayName = strings.TrimSpace(inputModel.Inputs[0].Value())
				endpointServiceName = strings.TrimSpace(inputModel.Inputs[1].Value())
				if includeRegion {
					endpointServiceRegion = strings.TrimSpace(inputModel.Inputs[2].Value())
				}
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				displayName, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}
				endpointServiceName, err = cmd.Flags().GetString(flag.EndpointServiceName)
				if err != nil {
					return errors.Trace(err)
				}
				endpointServiceRegion, err = cmd.Flags().GetString(flag.EndpointServiceRegion)
				if err != nil {
					return errors.Trace(err)
				}
				connectionTypeRaw, err = cmd.Flags().GetString(flag.PrivateLinkConnectionType)
				if err != nil {
					return errors.Trace(err)
				}
			}

			displayName = strings.TrimSpace(displayName)
			if displayName == "" {
				return errors.New("display name is required")
			}
			endpointServiceName = strings.TrimSpace(endpointServiceName)
			if endpointServiceName == "" {
				return errors.New("endpoint service name is required")
			}

			connectionType, err := normalizePrivateLinkConnectionType(connectionTypeRaw)
			if err != nil {
				return err
			}
			if endpointServiceRegion != "" && connectionType == pl.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE {
				return errors.New("endpoint service region is only supported for AWS endpoint service")
			}

			privateLinkConnection := pl.NewPrivateLinkConnection(clusterID, displayName, connectionType)
			if connectionType == pl.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE {
				awsService := pl.NewAwsEndpointService(endpointServiceName)
				if strings.TrimSpace(endpointServiceRegion) != "" {
					awsService.SetRegion(strings.TrimSpace(endpointServiceRegion))
				}
				privateLinkConnection.SetAwsEndpointService(*awsService)
			} else {
				aliService := pl.NewAlicloudEndpointService(endpointServiceName)
				privateLinkConnection.SetAlicloudEndpointService(*aliService)
			}

			body := pl.NewPrivateLinkConnectionServiceCreatePrivateLinkConnectionBody(*privateLinkConnection)
			_, err = d.CreatePrivateLinkConnection(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("private link connection %s is created", displayName))
			if err != nil {
				return err
			}
			return nil
		},
	}

	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	createCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The display name of the private link connection.")
	createCmd.Flags().StringP(flag.PrivateLinkConnectionType, "", "", "The type of the private link connection, one of [\"AWS_ENDPOINT_SERVICE\" \"ALICLOUD_ENDPOINT_SERVICE\"].")
	createCmd.Flags().StringP(flag.EndpointServiceName, "", "", "The endpoint service name.")
	createCmd.Flags().StringP(flag.EndpointServiceRegion, "", "", "The endpoint service region (AWS only).")
	return createCmd
}

func getCreateInputModel(includeRegion bool) (ui.TextInputModel, error) {
	fieldCount := 2
	if includeRegion {
		fieldCount = 3
	}
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, fieldCount),
	}

	for i := 0; i < fieldCount; i++ {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 64
		switch i {
		case 0:
			t.Placeholder = "Display Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case 1:
			t.Placeholder = "Endpoint Service Name"
		case 2:
			t.Placeholder = "Endpoint Service Region (optional)"
		}
		m.Inputs[i] = t
	}

	p := tea.NewProgram(m)
	inputModel, err := p.Run()
	if err != nil {
		return ui.TextInputModel{}, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return ui.TextInputModel{}, util.InterruptError
	}
	return inputModel.(ui.TextInputModel), nil
}

func normalizePrivateLinkConnectionType(input string) (pl.PrivateLinkConnectionTypeEnum, error) {
	value := strings.TrimSpace(strings.ToUpper(input))
	switch value {
	case "AWS", string(pl.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE):
		return pl.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE, nil
	case "ALICLOUD", string(pl.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE):
		return pl.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE, nil
	default:
		return "", fmt.Errorf("invalid private link connection type: %s", input)
	}
}
