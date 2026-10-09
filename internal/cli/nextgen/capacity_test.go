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
	"fmt"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/rodaine/table"
	"github.com/stretchr/testify/require"
)

func TestCreateStillRequiresMaxRCU(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		t.Run(plan.commandName, func(t *testing.T) {
			h := helperWithNextGenClient(nil)
			h.NextGenClient = func() (cloud.NextGenClient, error) {
				t.Fatal("must reject missing capacity before obtaining a client")
				return nil, nil
			}
			command := createCmd(h, plan)
			command.SetArgs([]string{"--display-name", "test-instance", "--region", "aws-us-west-2"})
			require.ErrorContains(t, command.Execute(), `required flag(s) "max-rcu" not set`)
		})
	}
}

func TestListOptionalMaxRCU(t *testing.T) {
	for _, tc := range []struct{ fields, want string }{
		{`"maxRcu":"10000"`, "10000"},
		{`"maxRcu":null`, "-"},
		{`"capacityMode":"ELASTIC","baselineRcu":"5000"`, "-"},
	} {
		t.Run(tc.fields, func(t *testing.T) {
			var instance api.Nextgenv1beta2Tidb
			require.NoError(t, json.Unmarshal([]byte(`{"displayName":"test-instance","regionId":"aws-us-west-2","servicePlan":"Premium",`+tc.fields+`}`), &instance))
			client := &fakeNextGenClient{list: func(context.Context, api.TidbServiceListTidbsServicePlanParameter, *int32, *string) (*api.V1beta2ListTidbsResponse, error) {
				return &api.V1beta2ListTidbsResponse{Tidbs: []api.Nextgenv1beta2Tidb{instance}}, nil
			}}
			h := helperWithNextGenClient(client)
			h.IOStreams.CanPrompt = true
			previousWriter := table.DefaultWriter
			table.DefaultWriter = h.IOStreams.Out
			t.Cleanup(func() { table.DefaultWriter = previousWriter })
			command := listCmd(h, premiumPlan)
			require.NoError(t, command.Execute())
			out := h.IOStreams.Out.(*bytes.Buffer).String()
			require.Regexp(t, `aws-us-west-2\s+`+tc.want, out)
		})
	}
}

func TestRenameElasticInstancePreservesCapacity(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		t.Run(plan.commandName, func(t *testing.T) {
			var current api.Nextgenv1beta2Tidb
			baseline := "5000"
			if plan.servicePlan == api.V1BETA1SERVICEPLAN_ESSENTIAL_V2 {
				baseline = "2000"
			}
			body := fmt.Sprintf(`{"tidbId":"123","displayName":"old-name","regionId":"aws-us-west-2","servicePlan":%q,"baselineRcu":%q,"capacityMode":"ELASTIC"}`, plan.servicePlan, baseline)
			require.NoError(t, json.Unmarshal([]byte(body), &current))
			client := &fakeNextGenClient{
				get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return &current, nil },
				update: func(_ context.Context, _ string, request *api.TheTiDBCloudPremiumInstanceToUpdate) (*api.Nextgenv1beta2Tidb, error) {
					payload, err := json.Marshal(request)
					require.NoError(t, err)
					require.JSONEq(t, `{"displayName":"new-name"}`, string(payload))
					result := current
					result.DisplayName = request.GetDisplayName()
					return &result, nil
				},
			}
			h := helperWithNextGenClient(client)
			command := updateCmd(h, plan)
			command.SetArgs([]string{"--cluster-id", "123", "--display-name", "new-name", "--output", "json"})
			require.NoError(t, command.Execute())
			out := h.IOStreams.Out.(*bytes.Buffer).String()
			require.NotContains(t, out, `"maxRcu"`)
			require.Contains(t, out, `"baselineRcu": "`+baseline+`"`)
		})
	}
}
