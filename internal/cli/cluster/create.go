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
	"fmt"
	"time"

	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/openapi"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type createClusterField int

const (
	clusterProjectIDIdx createClusterField = iota
	clusterNameIdx
	passwordIdx
)

type CreateServerlessOpts struct {
	serverlessProviders []*clusterApi.ListProviderRegionsOKBodyItemsItems0
}

func CreateCmd() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create one cluster in the specified project.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() != 0 {
				err := cmd.MarkFlagRequired(flag.ClusterName)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.ClusterType)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.CloudProvider)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.Region)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.RootPassword)
				if err != nil {
					return err
				}
				err = cmd.MarkFlagRequired(flag.ProjectID)
				if err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKey, privateKey := util.GetAccessKeys()
			apiClient := openapi.NewApiClient(publicKey, privateKey)

			var clusterName string
			var clusterType string
			var cloudProvider string
			var region string
			var rootPassword string
			var projectID string
			if cmd.Flags().NFlag() == 0 {
				regions, err := apiClient.Cluster.ListProviderRegions(nil)
				if err != nil {
					return err
				}

				opts := CreateServerlessOpts{}
				for i, item := range regions.Payload.Items {
					if item.ClusterType == "DEVELOPER" {
						opts.serverlessProviders = append(opts.serverlessProviders, regions.Payload.Items[i])
					}
				}

				p := tea.NewProgram(ui.InitialSelectModel([]interface{}{"DEVELOPER"}, "Choose the cluster type:"))
				typeModel, err := p.StartReturningModel()
				if err != nil {
					return err
				}
				if m, _ := typeModel.(ui.SelectModel); m.Interrupted {
					return nil
				}
				clusterType = typeModel.(ui.SelectModel).Choices[typeModel.(ui.SelectModel).Selected].(string)

				set := hashset.New()
				for _, provider := range opts.serverlessProviders {
					set.Add(provider.CloudProvider)
				}
				p = tea.NewProgram(ui.InitialSelectModel(set.Values(), "Choose the cloud provider:"))
				providerModel, err := p.StartReturningModel()
				if err != nil {
					return err
				}
				if m, _ := providerModel.(ui.SelectModel); m.Interrupted {
					return nil
				}
				cloudProvider = providerModel.(ui.SelectModel).Choices[providerModel.(ui.SelectModel).Selected].(string)

				set = hashset.New()
				for _, provider := range opts.serverlessProviders {
					if provider.CloudProvider == providerModel.(ui.SelectModel).Choices[providerModel.(ui.SelectModel).Selected] {
						set.Add(provider.Region)
					}
				}
				p = tea.NewProgram(ui.InitialSelectModel(set.Values(), "Choose the cloud region:"))
				regionModel, err := p.StartReturningModel()
				if err != nil {
					return err
				}
				if m, _ := regionModel.(ui.SelectModel); m.Interrupted {
					return nil
				}
				region = regionModel.(ui.SelectModel).Choices[regionModel.(ui.SelectModel).Selected].(string)

				p = tea.NewProgram(initialCreateInputModel())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return err
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return nil
				}

				clusterName = inputModel.(ui.TextInputModel).Inputs[clusterNameIdx].Value()
				rootPassword = inputModel.(ui.TextInputModel).Inputs[passwordIdx].Value()
				projectID = inputModel.(ui.TextInputModel).Inputs[clusterProjectIDIdx].Value()
			} else {
				cName, err := cmd.Flags().GetString(flag.ClusterName)
				if err != nil {
					return err
				}
				clusterName = cName

				clusterType, err = cmd.Flags().GetString(flag.ClusterType)
				if err != nil {
					return err
				}
				cloudProvider, err = cmd.Flags().GetString(flag.CloudProvider)
				if err != nil {
					return err
				}
				region, err = cmd.Flags().GetString(flag.Region)
				if err != nil {
					return err
				}
				rootPassword, err = cmd.Flags().GetString(flag.RootPassword)
				if err != nil {
					return err
				}
				projectID, err = cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return err
				}
			}

			clusterDefBody := &clusterApi.CreateClusterBody{}
			err := clusterDefBody.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"name": "%s",
			"cluster_type": %s,
			"cloud_provider": %s,
			"region": %s,
			"config" : {
				"root_password": "%s",
				"ip_access_list": [
					{
						"CIDR": "0.0.0.0/0",
						"description": "Allow All"
					}
				]
			}
			}`, clusterName, clusterType, cloudProvider, region, rootPassword)))
			if err != nil {
				return err
			}

			task := func() tea.Msg {
				createClusterResult, err := apiClient.Cluster.CreateCluster(clusterApi.NewCreateClusterParams().WithProjectID(projectID).WithBody(*clusterDefBody))
				if err != nil {
					return err
				}
				newClusterID := *createClusterResult.GetPayload().ID

				ticker := time.NewTicker(1 * time.Second)
				for {
					select {
					case <-time.After(2 * time.Minute):
						return ui.Result("Timeout waiting for cluster to be ready, please check status on dashboard.")
					case <-ticker.C:
						clusterResult, err := apiClient.Cluster.GetCluster(clusterApi.NewGetClusterParams().
							WithClusterID(newClusterID).
							WithProjectID(projectID))
						if err != nil {
							return err
						}
						s := clusterResult.GetPayload().Status.ClusterStatus
						if s == "AVAILABLE" {
							return ui.Result(fmt.Sprintf("Cluster %s is ready.", newClusterID))
						}
					}
				}
			}
			p := tea.NewProgram(ui.InitialSpinnerModel(task, "Waiting for cluster to be ready..."))
			createModel, err := p.StartReturningModel()
			if err != nil {
				return err
			}
			if m, _ := createModel.(ui.SpinnerModel); m.Err != nil {
				color.Red(m.Err.Error())
			} else {
				color.Green(m.Output)
			}
			return nil
		},
	}

	createCmd.Flags().String(flag.ClusterName, "", "Name of the cluster to de created")
	createCmd.Flags().String(flag.ClusterType, "", "Cluster type, only support serverless now")
	createCmd.Flags().String(flag.CloudProvider, "", "Cloud provider, e.g. AWS")
	createCmd.Flags().StringP(flag.Region, flag.RegionShort, "", "Cloud region")
	createCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The ID of the project, in which the cluster will be created")
	createCmd.Flags().String(flag.RootPassword, "", "The root password of the cluster")
	return createCmd
}

func initialCreateInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 64
		f := createClusterField(i)

		switch f {
		case clusterProjectIDIdx:
			t.Placeholder = "Project ID"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.Validate = func(s string) error {
				if len(s) == 0 {
					return errors.New("project ID is required")
				}
				return nil
			}
		case clusterNameIdx:
			t.Placeholder = "Cluster Name"
			t.Validate = func(s string) error {
				if len(s) == 0 {
					return errors.New("cluster Name is required")
				}
				return nil
			}
		case passwordIdx:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.Validate = func(s string) error {
				if len(s) == 0 {
					return errors.New("password is required")
				}
				return nil
			}
		}

		m.Inputs[i] = t
	}

	return m
}
