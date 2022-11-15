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
	"errors"
	"time"

	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/openapi"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type deleteClusterField int

const (
	projectIDIdx deleteClusterField = iota
	clusterIDIdx
)

func DeleteCmd() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete a cluster from your project.",
		Aliases: []string{"rm"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() != 0 {
				err := cmd.MarkFlagRequired(flag.ProjectID)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.ClusterID)
				if err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKey, privateKey := util.GetAccessKeys()
			apiClient := openapi.NewApiClient(publicKey, privateKey)

			var projectID string
			var clusterID string
			if cmd.Flags().NFlag() == 0 {
				p := tea.NewProgram(initialClusterIdentifies())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return err
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return nil
				}

				projectID = inputModel.(ui.TextInputModel).Inputs[projectIDIdx].Value()
				clusterID = inputModel.(ui.TextInputModel).Inputs[clusterIDIdx].Value()
			} else {
				pID, err := cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return err
				}

				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return err
				}
				projectID = pID
				clusterID = cID
			}

			params := clusterApi.NewDeleteClusterParams().
				WithProjectID(projectID).
				WithClusterID(clusterID)
			_, err := apiClient.Cluster.DeleteCluster(params)
			if err != nil {
				return err
			}

			ticker := time.NewTicker(1 * time.Second)
			for {
				select {
				case <-time.After(2 * time.Minute):
					return errors.New("timeout waiting for deleting cluster, please check status on dashboard")
				case <-ticker.C:
					_, err := apiClient.Cluster.GetCluster(clusterApi.NewGetClusterParams().
						WithClusterID(clusterID).
						WithProjectID(projectID))
					if err != nil {
						if _, ok := err.(*clusterApi.GetClusterNotFound); ok {
							color.Green("cluster deleted")
							return nil
						}
						return err
					}
				}
			}
		},
	}

	deleteCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The project ID of the cluster to be deleted.")
	deleteCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster to be deleted.")
	return deleteCmd
}

func initialClusterIdentifies() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 64
		f := deleteClusterField(i)

		switch f {
		case projectIDIdx:
			t.Placeholder = "Project ID"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case clusterIDIdx:
			t.Placeholder = "Cluster ID"
		}

		m.Inputs[i] = t
	}

	return m
}
