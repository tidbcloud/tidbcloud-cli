package cluster

import (
	"encoding/json"
	"fmt"

	"tidbcloud-cli/internal/openapi"
	"tidbcloud-cli/internal/util"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:     "list <projectID>",
		Short:   "List all clusters in a project.",
		Args:    util.RequiredArgs("projectID"),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKey, privateKey := util.GetAccessKeys()
			apiClient := openapi.NewApiClient(publicKey, privateKey)
			pID := args[0]

			params := clusterApi.NewListClustersOfProjectParams().WithProjectID(pID)

			cluster, err := apiClient.Cluster.ListClustersOfProject(params)
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
