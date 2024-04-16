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
	"encoding/json"
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"

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
		flag.ServerlessAnnotations,
		flag.ServerlessLabels,
	}
}

type mutableField string

const (
	DisplayName mutableField = "displayName"
	Annotations mutableField = "annotations"
	Labels      mutableField = "labels"
)

var mutableFields = []string{
	string(DisplayName),
	string(Labels),
	string(Annotations),
}

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateOpts{
		interactive: true,
	}

	var updateCmd = &cobra.Command{
		Use:         "update",
		Short:       "Update a TiDB Serverless cluster",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Update a TiDB Serverless cluster in interactive mode:
 $ %[1]s serverless update

 Update displayName of a TiDB Serverless cluster in non-interactive mode:
 $ %[1]s serverless update -c <cluster-id> --display-name <new-cluster-mame>
 
 Update labels of a TiDB Serverless cluster in non-interactive mode:
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
				cmd.MarkFlagsMutuallyExclusive(flag.DisplayName, flag.ServerlessAnnotations, flag.ServerlessLabels)
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
			var displayName, labels, annotations string
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
				fidleValue := inputModel.(ui.TextInputModel).Inputs[0].Value()

				switch fieldName {
				case string(DisplayName):
					displayName = fidleValue
				case string(Annotations):
					annotations = fidleValue
				case string(Labels):
					labels = fidleValue
				default:
					return errors.Errorf("invalid field %s", fieldName)
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
				annotations, err = cmd.Flags().GetString(flag.ServerlessAnnotations)
				if err != nil {
					return errors.Trace(err)
				}
			}

			body := &serverlessApi.ServerlessServicePartialUpdateClusterBody{
				Cluster: &serverlessApi.ServerlessServicePartialUpdateClusterParamsBodyCluster{},
			}
			if displayName != "" {
				body.Cluster.DisplayName = displayName
				fieldName = string(DisplayName)
			}
			if labels != "" {
				labelsMap, err := stringToMap(labels)
				if err != nil {
					return errors.Errorf("invalid labels %s", labels)
				}
				body.Cluster.Labels = labelsMap
				fieldName = string(Labels)
			}
			if annotations != "" {
				annotationsMap, err := stringToMap(annotations)
				if err != nil {
					return errors.Errorf("invalid annotations %s", annotations)
				}
				body.Cluster.Annotations = annotationsMap
				fieldName = string(Annotations)
			}
			body.UpdateMask = &fieldName

			params := serverlessApi.NewServerlessServicePartialUpdateClusterParams().WithClusterClusterID(clusterID).WithBody(*body)
			_, err = d.PartialUpdateCluster(params)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("cluster %s updated", clusterID)))
			return nil
		},
	}

	updateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster to be updated.")
	updateCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The new displayName of the cluster to be updated.")
	updateCmd.Flags().String(flag.ServerlessLabels, "", "The labels of the cluster to be added or updated.\nInteractive example: {\"label1\":\"value1\",\"label2\":\"value2\"}\nNonInteractive example: \"{\\\"label1\\\":\\\"value1\\\",\\\"label2\\\":\\\"value2\\\"}\".")
	updateCmd.Flags().String(flag.ServerlessAnnotations, "", "The annotations of the cluster to be added or updated.\nInteractive example: {\"annotation1\":\"value1\",\"annotation2\":\"value2\"}\nNonInteractive example: \"{\\\"annotation1\\\":\\\"value1\\\",\\\"annotation2\\\":\\\"value2\\\"}\".")
	return updateCmd
}

func GetUpdateClusterInput() (tea.Model, error) {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 1),
	}
	t := textinput.New()
	t.Cursor.Style = config.CursorStyle
	t.CharLimit = 64
	t.Placeholder = "New value"
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
