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
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"
)

func TestDescribeJSONOutputFlag(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, outputArgs := range [][]string{nil, {"-o", "json"}, {"-ojson"}, {"--output=json"}} {
			client := &fakeNextGenClient{get: func(_ context.Context, id string) (*api.Nextgenv1beta2Tidb, error) {
				require.Equal(t, "123", id)
				return &api.Nextgenv1beta2Tidb{ServicePlan: plan.servicePlan, DisplayName: "test"}, nil
			}}
			h := helperWithNextGenClient(client)
			var out bytes.Buffer
			h.IOStreams.Out = &out
			root := &cobra.Command{Use: "ticloud"}
			root.AddCommand(cmd(h, plan))
			args := append([]string{plan.commandName, "--describe", "-c", "123"}, outputArgs...)
			normalized, err := NormalizeActionArgs(root, args)
			require.NoError(t, err)
			root.SetArgs(normalized)
			require.NoError(t, root.Execute())
			require.True(t, json.Valid(out.Bytes()))
			require.Contains(t, out.String(), "test")
		}
		command := describeCmd(helperWithNextGenClient(nil), plan)
		command.SetArgs([]string{"-c", "123", "-o", "yaml"})
		require.ErrorContains(t, command.Execute(), "describe supports json only")
	}
}
