package project

import (
	"encoding/json"
	"fmt"

	"tidbcloud-cli/internal/openapi"
	"tidbcloud-cli/internal/util"

	projectApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"

	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all accessible projects.",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKey, privateKey := util.GetAccessKeys()
			apiClient := openapi.NewApiClient(publicKey, privateKey)

			params := projectApi.NewListProjectsParams()

			cluster, err := apiClient.Project.ListProjects(params)
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
