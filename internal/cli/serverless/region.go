// Copyright 2023 PingCAP, Inc.
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
	"fmt"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func RegionCmd(h *internal.Helper) *cobra.Command {
	var regionCmd = &cobra.Command{
		Use:         "region",
		Short:       "List all available regions for TiDB Serverless",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  List all available regions for TiDB Serverless:
  $ %[1]s serverless region`, config.CliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			regions, err := d.ListProviderRegions(serverlessApi.NewServerlessServiceListRegionsParams())
			if err != nil {
				return errors.Trace(err)
			}

			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				err = output.PrintJson(h.IOStreams.Out, regions.Payload.Regions)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"Name",
					"DisplayName",
					"Provider",
				}

				var rows []output.Row
				for _, item := range regions.Payload.Regions {
					rows = append(rows, output.Row{
						*item.Name,
						item.DisplayName,
						string(*item.Provider),
					})
				}
				err = output.PrintHumanTable(h.IOStreams.Out, columns, rows)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				return fmt.Errorf("unsupported output format: %s", format)
			}
			return nil
		},
	}
	regionCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return regionCmd
}
