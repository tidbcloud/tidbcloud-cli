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

package changefeed

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cdc"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type PrivateLinkCreateOpts struct {
	interactive bool
}

func (c PrivateLinkCreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkServiceName,
		flag.ChangefeedType,
	}
}

func (c PrivateLinkCreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkServiceName,
		flag.ChangefeedType,
	}
}

func (c *PrivateLinkCreateOpts) MarkInteractive(cmd *cobra.Command) error {
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

func PrivateLinkCreateCmd(h *internal.Helper) *cobra.Command {
	opts := PrivateLinkCreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a changefeed private link endpoint",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a changefeed private link endpoint in interactive mode:
  $ %[1]s serverless changefeed private-link create

  Create a changefeed private link endpoint in non-interactive mode:
  $ %[1]s serverless changefeed private-link create -c <cluster-id> --private-link-service-name <service-name> --type <type>`,
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
			var privateLinkServiceName string
			var changefeedTypeRaw string
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

				changefeedTypeRaw, err = cloud.GetSelectedField([]string{
					string(cdc.CHANGEFEEDTYPEENUM_KAFKA),
					string(cdc.CHANGEFEEDTYPEENUM_MYSQL),
				}, "Choose the changefeed type:")
				if err != nil {
					return err
				}

				privateLinkServiceName, err = promptPrivateLinkServiceName()
				if err != nil {
					return err
				}
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				privateLinkServiceName, err = cmd.Flags().GetString(flag.PrivateLinkServiceName)
				if err != nil {
					return errors.Trace(err)
				}
				changefeedTypeRaw, err = cmd.Flags().GetString(flag.ChangefeedType)
				if err != nil {
					return errors.Trace(err)
				}
			}

			privateLinkServiceName = strings.TrimSpace(privateLinkServiceName)
			if privateLinkServiceName == "" {
				return errors.New("private link service name is required")
			}

			changefeedType, err := normalizeChangefeedType(changefeedTypeRaw)
			if err != nil {
				return err
			}

			body := cdc.NewChangefeedServiceCreatePrivateLinkEndpointBody(privateLinkServiceName, changefeedType)
			_, err = d.CreatePrivateLinkEndpoint(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("private link endpoint %s is created", privateLinkServiceName))
			if err != nil {
				return err
			}
			return nil
		},
	}

	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	createCmd.Flags().StringP(flag.PrivateLinkServiceName, "", "", "The private link service name.")
	createCmd.Flags().StringP(flag.ChangefeedType, "", "", "The changefeed type, one of [\"KAFKA\" \"MYSQL\"].")
	return createCmd
}

func promptPrivateLinkServiceName() (string, error) {
	model := ui.TextInputModel{
		Inputs: make([]textinput.Model, 1),
	}
	t := textinput.New()
	t.Cursor.Style = config.CursorStyle
	t.CharLimit = 64
	t.Placeholder = "Private Link Service Name"
	t.Focus()
	t.PromptStyle = config.FocusedStyle
	t.TextStyle = config.FocusedStyle
	model.Inputs[0] = t

	p := tea.NewProgram(model)
	inputModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return "", util.InterruptError
	}
	return inputModel.(ui.TextInputModel).Inputs[0].Value(), nil
}

func normalizeChangefeedType(input string) (cdc.ChangefeedTypeEnum, error) {
	value := strings.TrimSpace(strings.ToUpper(input))
	switch value {
	case string(cdc.CHANGEFEEDTYPEENUM_KAFKA):
		return cdc.CHANGEFEEDTYPEENUM_KAFKA, nil
	case string(cdc.CHANGEFEEDTYPEENUM_MYSQL):
		return cdc.CHANGEFEEDTYPEENUM_MYSQL, nil
	default:
		return "", fmt.Errorf("invalid changefeed type: %s", input)
	}
}
