package cluster

import (
	"errors"
	"time"

	"tidbcloud-cli/internal/cli/ui"
	"tidbcloud-cli/internal/openapi"
	"tidbcloud-cli/internal/util"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type deleteClusterField int

const (
	toDeletedProjectID deleteClusterField = iota
	toDeletedClusterID
)

func DeleteCmd() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete a cluster from your project.",
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKey, privateKey := util.GetAccessKeys()
			apiClient := openapi.NewApiClient(publicKey, privateKey)

			p := tea.NewProgram(initialClusterIdentifies())
			inputModel, err := p.StartReturningModel()
			if err != nil {
				return err
			}
			if inputModel.(ui.TextInputModel).Interrupted {
				return nil
			}

			params := clusterApi.NewDeleteClusterParams().
				WithProjectID(inputModel.(ui.TextInputModel).Inputs[toDeletedProjectID].Value()).
				WithClusterID(inputModel.(ui.TextInputModel).Inputs[toDeletedClusterID].Value())
			_, err = apiClient.Cluster.DeleteCluster(params)
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
						WithClusterID(inputModel.(ui.TextInputModel).Inputs[toDeletedClusterID].Value()).
						WithProjectID(inputModel.(ui.TextInputModel).Inputs[toDeletedProjectID].Value()))
					if err != nil {
						if _, ok := err.(*clusterApi.GetClusterNotFound); ok {
							color.Green("cluster deleted")
							return nil
						}
						return err
					}
				}
			}
			return nil
		},
	}

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
		case toDeletedProjectID:
			t.Placeholder = "Project ID"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case toDeletedClusterID:
			t.Placeholder = "Cluster ID"
		}

		m.Inputs[i] = t
	}

	return m
}
