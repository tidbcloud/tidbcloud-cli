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
	"fmt"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/spf13/cobra"
)

func regionCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	command := &cobra.Command{
		Use:   "region",
		Short: fmt.Sprintf("List regions available for %s", plan.displayName),
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
			var regions []api.V1beta1Region
			var token *string
			seenTokens := make(map[string]struct{})
			size := int32(h.QueryPageSize)
			for {
				response, err := client.ListRegions(cmd.Context(), plan.queryPlan, &size, token)
				if err != nil {
					return err
				}
				if response == nil {
					return fmt.Errorf("the API returned an empty region response")
				}
				regions = append(regions, response.Regions...)
				if response.GetNextPageToken() == "" {
					break
				}
				nextToken := response.GetNextPageToken()
				if _, ok := seenTokens[nextToken]; ok {
					return fmt.Errorf("the API returned a repeated next page token")
				}
				seenTokens[nextToken] = struct{}{}
				token = &nextToken
			}

			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				totalSize := int32(len(regions))
				return output.PrintJson(h.IOStreams.Out, &api.V1beta2ListRegionsResponse{Regions: regions, TotalSize: &totalSize})
			}

			columns := []output.Column{"RegionID", "DisplayName", "CloudProvider", "ServicePlans"}
			rows := make([]output.Row, 0, len(regions))
			for i := range regions {
				region := &regions[i]
				plans := make([]string, 0, len(region.ServicePlans))
				for _, servicePlan := range region.ServicePlans {
					plans = append(plans, string(servicePlan))
				}
				rows = append(rows, output.Row{region.GetRegionId(), region.GetDisplayName(), string(region.GetCloudProvider()), strings.Join(plans, ",")})
			}
			return output.PrintHumanTable(h.IOStreams.Out, columns, rows)
		},
	}
	command.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return command
}
