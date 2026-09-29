// Copyright 2026 PingCAP, Inc.
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

package nextgen

import (
	"context"
	"fmt"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/spf13/cobra"
)

func listCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: fmt.Sprintf("List %s instances", plan.displayName),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := validatedOutputFormat(cmd)
			if err != nil {
				return err
			}
			client, err := h.NextGenClient()
			if err != nil {
				return err
			}
			instances, err := retrieveTiDBs(cmd.Context(), client, plan, int32(h.QueryPageSize))
			if err != nil {
				return err
			}
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				totalSize := int32(len(instances))
				return output.PrintJson(h.IOStreams.Out, &api.V1beta2ListTidbsResponse{Tidbs: instances, TotalSize: &totalSize})
			}

			columns := []output.Column{"ID", "DisplayName", "State", "Region", "MaxRCU", "CreateTime"}
			rows := make([]output.Row, 0, len(instances))
			for i := range instances {
				instance := &instances[i]
				createTime := ""
				if instance.CreateTime != nil {
					createTime = instance.CreateTime.Format(time.RFC3339)
				}
				rows = append(rows, output.Row{
					instance.GetTidbId(), instance.DisplayName, string(instance.GetState()), instance.RegionId, instance.MaxRcu, createTime,
				})
			}
			return output.PrintHumanTable(h.IOStreams.Out, columns, rows)
		},
	}
	command.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return command
}

func retrieveTiDBs(ctx context.Context, client cloud.NextGenClient, plan planSpec, size int32) ([]api.Nextgenv1beta2Tidb, error) {
	var instances []api.Nextgenv1beta2Tidb
	var token *string
	seenTokens := make(map[string]struct{})
	for {
		response, err := client.ListTiDBs(ctx, plan.queryPlan, &size, token)
		if err != nil {
			return nil, err
		}
		if response == nil {
			return nil, fmt.Errorf("the API returned an empty list response")
		}
		for i := range response.Tidbs {
			if err := plan.validate(&response.Tidbs[i]); err != nil {
				return nil, err
			}
		}
		instances = append(instances, response.Tidbs...)
		if response.GetNextPageToken() == "" {
			return instances, nil
		}
		nextToken := response.GetNextPageToken()
		if _, ok := seenTokens[nextToken]; ok {
			return nil, fmt.Errorf("the API returned a repeated next page token")
		}
		seenTokens[nextToken] = struct{}{}
		token = &nextToken
	}
}
