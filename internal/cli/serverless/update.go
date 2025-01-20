// Copyright 2025 PingCAP, Inc.
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
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
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

type UpdateOpts struct {
	interactive bool
}

func (c UpdateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
		flag.ServerlessLabels,
		flag.PublicEndpointDisabled,
	}
}

type mutableField string

const (
	DisplayName            mutableField = "displayName"
	Labels                 mutableField = "labels"
	PublicEndpointDisabled mutableField = "endpoints.public.disabled"
)

const (
	PublicEndpointDisabledHumanReadable = "disable public endpoint"
)

var mutableFields = []string{
	string(DisplayName),
	string(Labels),
	string(PublicEndpointDisabledHumanReadable),
}

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateOpts{
		interactive: true,
	}

	var updateCmd = &cobra.Command{
		Use:         "update",
		Short:       "Update a TiDB Cloud Serverless cluster",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Update a TiDB Cloud Serverless cluster in interactive mode:
  $ %[1]s serverless update

  Update displayName of a TiDB Cloud Serverless cluster in non-interactive mode:
  $ %[1]s serverless update -c <cluster-id> --display-name <new-cluster-name>
 
  Update labels of a TiDB Cloud Serverless cluster in non-interactive mode:
  $ %[1]s serverless update -c <cluster-id> --labels "{\"label1\":\"value1\"}"`, config.CliName),
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
				cmd.MarkFlagsMutuallyExclusive(flag.DisplayName, flag.ServerlessLabels, flag.PublicEndpointDisabled)
				cmd.MarkFlagsOneRequired(flag.DisplayName, flag.ServerlessLabels, flag.PublicEndpointDisabled)
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
			var fieldName string
			var displayName string
			var labels string
			var publicEndpointDisabled bool
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

				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				fieldName, err = cloud.GetSelectedField(mutableFields)
				if err != nil {
					return err
				}

				if fieldName == PublicEndpointDisabledHumanReadable {
					prompt := &survey.Confirm{
						Message: "Disable the public endpoint of the cluster?",
						Default: false,
					}
					err = survey.AskOne(prompt, &publicEndpointDisabled)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
				} else {
					// variables for input
					inputModel, err := GetUpdateClusterInput(fieldName)
					if err != nil {
						return err
					}

					fieldValue := inputModel.(ui.TextInputModel).Inputs[0].Value()

					switch fieldName {
					case string(DisplayName):
						displayName = fieldValue
					case string(Labels):
						labels = fieldValue
					default:
						return errors.Errorf("invalid field %s", fieldName)
					}
				}

			} else {
				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID = cID
				displayName, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}
				labels, err = cmd.Flags().GetString(flag.ServerlessLabels)
				if err != nil {
					return errors.Trace(err)
				}
				publicEndpointDisabled, err = cmd.Flags().GetBool(flag.PublicEndpointDisabled)
				if err != nil {
					return errors.Trace(err)
				}
			}

			body := &cluster.V1beta1ServerlessServicePartialUpdateClusterBody{
				Cluster: &cluster.RequiredTheClusterToBeUpdated{},
			}
			if displayName != "" {
				body.Cluster.DisplayName = &displayName
				fieldName = string(DisplayName)
			}
			if labels != "" {
				labelsMap, err := stringToMap(labels)
				if err != nil {
					return errors.Errorf("invalid labels %s", labels)
				}
				body.Cluster.Labels = &labelsMap
				fieldName = string(Labels)
			}
			// if fieldName is PublicEndpointDisabled, means this field is changed in Interactive mode
			if cmd.Flags().Changed(flag.PublicEndpointDisabled) || fieldName == string(PublicEndpointDisabledHumanReadable) {
				body.Cluster.Endpoints = &cluster.V1beta1ClusterEndpoints{
					Public: &cluster.EndpointsPublic{
						Disabled: &publicEndpointDisabled,
					},
				}
				fieldName = string(PublicEndpointDisabled)
			}
			body.UpdateMask = fieldName
			_, err = d.PartialUpdateCluster(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("Cluster %s updated", clusterID)))
			return nil
		},
	}

	updateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster to be updated.")
	updateCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The new displayName of the cluster to be updated.")
	updateCmd.Flags().String(flag.ServerlessLabels, "", "The labels of the cluster to be added or updated.\nInteractive example: {\"label1\":\"value1\",\"label2\":\"value2\"}.\nNonInteractive example: \"{\\\"label1\\\":\\\"value1\\\",\\\"label2\\\":\\\"value2\\\"}\".")
	updateCmd.Flags().Bool(flag.PublicEndpointDisabled, false, "Disable the public endpoint of the cluster.")
	return updateCmd
}

func GetUpdateClusterInput(fieldName string) (tea.Model, error) {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 1),
	}
	t := textinput.New()
	t.Cursor.Style = config.CursorStyle
	t.CharLimit = 64
	switch fieldName {
	case string(Labels):
		t.Placeholder = "update labels, e.g. {\"label1\":\"value1\",\"label2\":\"value2\"}"
	default:
		t.Placeholder = "new value"
	}
	t.Focus()
	t.PromptStyle = config.FocusedStyle
	t.TextStyle = config.FocusedStyle
	m.Inputs[0] = t

	p := tea.NewProgram(m)
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}

func stringToMap(s string) (map[string]string, error) {
	if s == "" {
		return nil, nil
	}
	m := make(map[string]string)
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
