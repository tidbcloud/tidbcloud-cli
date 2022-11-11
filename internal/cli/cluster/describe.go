package cluster

import (
	"encoding/json"
	"fmt"

	"tidbcloud-cli/internal/cli/ui"
	"tidbcloud-cli/internal/openapi"
	"tidbcloud-cli/internal/util"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func DescribeCmd() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:     "describe",
		Short:   "Describe a cluster.",
		Aliases: []string{"get"},
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

			params := clusterApi.NewGetClusterParams().
				WithProjectID(inputModel.(ui.TextInputModel).Inputs[toDeletedProjectID].Value()).
				WithClusterID(inputModel.(ui.TextInputModel).Inputs[toDeletedClusterID].Value())
			cluster, err := apiClient.Cluster.GetCluster(params)
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(cluster, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}

	return deleteCmd
}
