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
	"encoding/json"
	"fmt"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/serverless/client/serverless_service"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
)

func RegionsCmd(h *internal.Helper) *cobra.Command {
	var regionsCmd = &cobra.Command{
		Use:         "regions",
		Short:       "List all available regions for serverless cluster",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  List all available regions for serverless cluster:
 $ %[1]s serverless regions`, config.CliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}

			regions, err := d.ListProviderRegions(serverlessApi.NewServerlessServiceListRegionsParams())
			if err != nil {
				return errors.Trace(err)
			}
			v, err := json.MarshalIndent(regions.Payload, "", "  ")
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, string(v))
			return nil
		},
	}
	return regionsCmd
}
