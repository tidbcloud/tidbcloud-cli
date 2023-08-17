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

package cluster

import (
	"fmt"
	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	"github.com/fatih/color"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	serverlessApi "tidbcloud-cli/pkg/tidbcloud/serverless/client/serverless_service"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type UpdateOpts struct {
	interactive bool
}

func (c UpdateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.UpdateField,
		flag.UpdateValue,
	}
}

type mutableField string

const (
	DisplayName mutableField = "displayName"
)

var mutableFields = []string{
	string(DisplayName),
}

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateOpts{
		interactive: true,
	}

	var updateCmd = &cobra.Command{
		Use:         "update",
		Short:       "Update a cluster",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Update a cluster in interactive mode:
 $ %[1]s cluster update

 Update a cluster in non-interactive mode:
 $ %[1]s cluster update -c <cluster-id> --field <fieldName> --value <newValue>`, config.CliName),
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

			var clusterID string
			var fieldName string
			var newValue string
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

				fieldName, err = cloud.GetSelectedField(mutableFields)
				if err != nil {
					return err
				}

				// variables for input
				inputModel, err := GetUpdateClusterInput()
				if err != nil {
					return err
				}
				newValue = inputModel.(ui.TextInputModel).Inputs[0].Value()
			} else {
				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID = cID
				fieldName, err = cmd.Flags().GetString(flag.UpdateField)
				if err != nil {
					return errors.Trace(err)
				}
				newValue, err = cmd.Flags().GetString(flag.UpdateValue)
				if err != nil {
					return errors.Trace(err)
				}
			}
			body, err := generateUpdateBody(fieldName, newValue)
			if err != nil {
				return err
			}
			params := serverlessApi.NewServerlessServicePartialUpdateClusterParams().WithClusterClusterID(clusterID).WithBody(*body)
			_, err = d.PartialUpdateCluster(params)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("cluster %s updated", clusterID)))
			return nil
		},
	}

	updateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster to be deleted")
	updateCmd.Flags().String(flag.UpdateField, "", "The field you want to update. Support [\"displayName\"] now")
	updateCmd.Flags().String(flag.UpdateValue, "", "The value you want to update of the field, e.g. \"newName\"")
	updateCmd.MarkFlagsRequiredTogether(flag.UpdateField, flag.UpdateValue)

	return updateCmd
}

func GetUpdateClusterInput() (tea.Model, error) {

	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 1),
	}
	t := textinput.New()
	t.CursorStyle = config.CursorStyle
	t.CharLimit = 64
	t.Placeholder = "New value"
	t.Focus()
	t.PromptStyle = config.FocusedStyle
	t.TextStyle = config.FocusedStyle
	m.Inputs[0] = t

	p := tea.NewProgram(m)
	inputModel, err := p.StartReturningModel()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}

func generateUpdateBody(field, value string) (*serverlessApi.ServerlessServicePartialUpdateClusterBody, error) {
	body := &serverlessApi.ServerlessServicePartialUpdateClusterBody{
		Cluster:    &serverlessApi.ServerlessServicePartialUpdateClusterParamsBodyCluster{},
		UpdateMask: field,
	}

	switch field {
	case string(DisplayName):
		body.Cluster.DisplayName = value
		return body, nil
	}
	return nil, fmt.Errorf("unsupported update field %s", field)
}
