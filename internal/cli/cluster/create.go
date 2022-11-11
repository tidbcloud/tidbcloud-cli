package cluster

import (
	"errors"
	"fmt"
	"time"

	"tidbcloud-cli/internal/cli/ui"
	"tidbcloud-cli/internal/openapi"
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
	projectID createClusterField = iota
	clusterName
	password
)

type CreateServerlessOpts struct {
	serverlessProviders []*clusterApi.ListProviderRegionsOKBodyItemsItems0
}

type CreateClusterResule struct {
	ui.SpinnerModel
	clusterID string
}

func CreateCmd() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create one cluster in the specified project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKey, privateKey := util.GetAccessKeys()
			apiClient := openapi.NewApiClient(publicKey, privateKey)
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

			p = tea.NewProgram(initialCreateInputModel())
			inputModel, err := p.StartReturningModel()
			if err != nil {
				return err
			}
			if inputModel.(ui.TextInputModel).Interrupted {
				return nil
			}

			clusterDefBody := &clusterApi.CreateClusterBody{}
			err = clusterDefBody.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"name": "%s",
			"cluster_type": "DEVELOPER",
			"cloud_provider": "AWS",
			"region": "us-west-2",
			"config" : {
				"root_password": "%s",
				"ip_access_list": [
					{
						"CIDR": "0.0.0.0/0",
						"description": "Allow All"
					}
				]
			}
			}`, inputModel.(ui.TextInputModel).Inputs[clusterName].Value(), inputModel.(ui.TextInputModel).Inputs[password].Value())))
			task := func() tea.Msg {
				createClusterResult, err := apiClient.Cluster.CreateCluster(clusterApi.NewCreateClusterParams().WithProjectID("1372813089189381287").WithBody(*clusterDefBody))
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
							WithProjectID(inputModel.(ui.TextInputModel).Inputs[projectID].Value()))
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
			p = tea.NewProgram(ui.InitialSpinnerModel(task, "Waiting for cluster to be ready..."))
			createModel, err := p.StartReturningModel()
			if err := p.Start(); err != nil {
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
		case projectID:
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
		case clusterName:
			t.Placeholder = "Cluster Name"
			t.Validate = func(s string) error {
				if len(s) == 0 {
					return errors.New("cluster Name is required")
				}
				return nil
			}
		case password:
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
