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
	"strings"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestCreateRequiresExplicitCapacity(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		t.Run(plan.commandName, func(t *testing.T) {
			h := helperWithNextGenClient(nil)
			h.NextGenClient = func() (cloud.NextGenClient, error) {
				t.Fatal("must reject missing capacity before obtaining a client")
				return nil, nil
			}
			command := createCmd(h, plan)
			command.SetArgs([]string{"--display-name", "test-instance", "--region", "aws-us-west-2"})
			require.ErrorContains(t, command.Execute(), `at least one of the flags in the group [max-rcu baseline-rcu] is required`)
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
			require.Contains(t, out, "CapacityMode")
			require.Contains(t, out, "BaselineRCU")
			if strings.Contains(tc.fields, "ELASTIC") {
				require.Regexp(t, `ELASTIC\s+5000`, out)
			}
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

func TestElasticCapacityRequests(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, action := range []string{"create", "update"} {
			for _, capacity := range []string{"max-rcu", "baseline-rcu"} {
				t.Run(plan.commandName+"/"+action+"/"+capacity, func(t *testing.T) {
					field := "maxRcu"
					other := "baselineRcu"
					if capacity == "baseline-rcu" {
						field, other = other, field
					}
					baseline := "5000"
					if plan.servicePlan == api.V1BETA1SERVICEPLAN_ESSENTIAL_V2 {
						baseline = "2000"
					}
					calls := 0
					check := func(payload interface{}) {
						calls++
						data, err := json.Marshal(payload)
						require.NoError(t, err)
						var body map[string]interface{}
						require.NoError(t, json.Unmarshal(data, &body))
						require.Equal(t, baseline, body[field])
						require.NotContains(t, body, other)
						require.NotContains(t, body, "capacityMode")
						require.NotContains(t, body, "minRcu")
					}
					current := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", plan.servicePlan)
					// Exercise both directions of mode switches; only the selected field goes on the wire.
					if capacity == "max-rcu" {
						current.SetBaselineRcu(baseline)
					} else {
						current.SetMaxRcu("20000")
					}
					client := &fakeNextGenClient{
						get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return current, nil },
						create: func(_ context.Context, b *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error) {
							check(b)
							return current, nil
						},
						update: func(_ context.Context, _ string, b *api.TheTiDBCloudPremiumInstanceToUpdate) (*api.Nextgenv1beta2Tidb, error) {
							check(b)
							return current, nil
						},
					}
					root := &cobra.Command{Use: "ticloud"}
					root.AddCommand(cmd(helperWithNextGenClient(client), plan))
					args := []string{plan.commandName, "--" + capacity, baseline, "--" + action, "-ojson"}
					if action == "create" {
						args = append(args, "--display-name", "test-instance", "--region", "aws-us-west-2")
					} else {
						args = append(args, "-c", "123")
					}
					normalized, err := NormalizeActionArgs(root, args)
					require.NoError(t, err)
					root.SetArgs(normalized)
					require.NoError(t, root.Execute())
					require.Equal(t, 1, calls)
				})
			}
		}
	}
}

func TestInvalidCapacityRejectedBeforeClient(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, action := range []string{"create", "update"} {
			for _, capacity := range [][]string{
				{"--baseline-rcu", "0"}, {"--baseline-rcu", "-1"}, {"--baseline-rcu", "invalid"},
				{"--baseline-rcu", "9223372036854775808"},
				{"--baseline-rcu", "5000", "--max-rcu", "20000"},
				{"--baseline-rcu", "0", "--max-rcu", "20000"},
			} {
				t.Run(fmt.Sprint(plan.commandName, action, capacity), func(t *testing.T) {
					h := helperWithNextGenClient(nil)
					h.NextGenClient = func() (cloud.NextGenClient, error) { t.Fatal("invalid request reached client"); return nil, nil }
					root := cmd(h, plan)
					args := append([]string{action}, capacity...)
					if action == "create" {
						args = append(args, "--display-name", "test-instance", "--region", "aws-us-west-2")
					} else {
						args = append(args, "-c", "123")
					}
					args = append(args, "-ojson")
					root.SetArgs(args)
					command, err := root.ExecuteC()
					require.Error(t, err)
					result, ok := JSONError(command, args, err)
					require.True(t, ok)
					require.True(t, json.Valid([]byte(result)))
				})
			}
		}
	}
}
