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
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"
	"testing"
)

func TestPrivateEndpointBind(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		t.Run(plan.commandName, func(t *testing.T) {
			called := false
			client := &fakeNextGenClient{get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
				v := &api.Nextgenv1beta2Tidb{}
				v.SetServicePlan(plan.servicePlan)
				v.SetRegionId("aws-us-west-2")
				return v, nil
			}, bindPrivate: func(_ context.Context, id string, b *api.Nextgenv1beta2PrivateEndpointConnection) (*api.Nextgenv1beta2PrivateEndpointConnection, error) {
				called = true
				require.Equal(t, "123", id)
				require.Equal(t, "aws-us-west-2", b.RegionId)
				require.Equal(t, "vpce-example", b.EndpointId)
				return b, nil
			}}
			cmd := privateEndpointCmd(helperWithNextGenClient(client), plan)
			cmd.SetArgs([]string{"bind", "-c", "123", "--endpoint-id", "vpce-example", "-o", "json"})
			require.NoError(t, cmd.Execute())
			require.True(t, called)
		})
	}
}

func TestPrivateEndpointRejectsInvalidInput(t *testing.T) {
	for _, args := range [][]string{{"bind", "-c", "123"}, {"describe", "-c", "123", "--connection-id", "../x"}, {"unbind", "-c", "123", "--connection-id", "c1"}, {"nonsense"}} {
		cmd := privateEndpointCmd(helperWithNextGenClient(&fakeNextGenClient{}), premiumPlan)
		cmd.SetArgs(args)
		require.Error(t, cmd.Execute())
	}
}
func TestPrivateEndpointPagination(t *testing.T) {
	for _, repeat := range []bool{false, true} {
		t.Run(fmt.Sprint(repeat), func(t *testing.T) {
			calls := 0
			client := &fakeNextGenClient{get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
				v := &api.Nextgenv1beta2Tidb{}
				v.SetServicePlan(premiumPlan.servicePlan)
				return v, nil
			}, listPrivate: func(_ context.Context, id string, size int32, token string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
				calls++
				p := &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{}
				if calls == 1 {
					require.Empty(t, token)
					p.SetNextPageToken("next")
				} else {
					require.Equal(t, "next", token)
					if repeat {
						p.SetNextPageToken("next")
					}
				}
				return p, nil
			}}
			cmd := privateEndpointCmd(helperWithNextGenClient(client), premiumPlan)
			cmd.SetArgs([]string{"list", "-c", "123"})
			err := cmd.Execute()
			if repeat {
				require.ErrorContains(t, err, "repeated")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, 2, calls)
		})
	}
}
func TestPrivateEndpointUnbindGuards(t *testing.T) {
	for _, scenario := range []string{"foreign", "cancel", "interrupt", "wrapped-interrupt", "force", "confirm"} {
		t.Run(scenario, func(t *testing.T) {
			deleted := false
			confirmed := false
			client := &fakeNextGenClient{get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
				v := &api.Nextgenv1beta2Tidb{}
				v.SetServicePlan(premiumPlan.servicePlan)
				return v, nil
			}, getPrivate: func(context.Context, string, string) (*api.Nextgenv1beta2PrivateEndpointConnection, error) {
				v := &api.Nextgenv1beta2PrivateEndpointConnection{}
				v.SetTidbId("123")
				if scenario == "foreign" {
					v.SetTidbId("456")
				}
				return v, nil
			}, unbindPrivate: func(context.Context, string, string) (*api.Nextgenv1beta2PrivateEndpointConnection, error) {
				deleted = true
				return &api.Nextgenv1beta2PrivateEndpointConnection{}, nil
			}}
			h := helperWithNextGenClient(client)
			h.IOStreams.CanPrompt = true
			cmd := privateEndpointActionCmd(h, premiumPlan, "unbind", func(id, c string) error {
				confirmed = true
				require.Equal(t, "123", id)
				require.Equal(t, "c1", c)
				if scenario == "interrupt" {
					return terminal.InterruptErr
				}
				if scenario == "wrapped-interrupt" {
					return fmt.Errorf("prompt: %w", terminal.InterruptErr)
				}
				if scenario == "cancel" {
					return fmt.Errorf("cancelled")
				}
				return nil
			})
			args := []string{"-c", "123", "--connection-id", "c1"}
			if scenario == "force" {
				args = append(args, "--force")
			}
			cmd.SetArgs(args)
			err := cmd.Execute()
			interrupted := scenario == "interrupt" || scenario == "wrapped-interrupt"
			if interrupted {
				require.ErrorIs(t, err, util.InterruptError)
			}
			if scenario == "foreign" || scenario == "cancel" || interrupted {
				require.Error(t, err)
				require.False(t, deleted)
			} else {
				require.NoError(t, err)
				require.True(t, deleted)
			}
			require.Equal(t, scenario == "cancel" || scenario == "confirm" || interrupted, confirmed)
		})
	}
}

func TestPrivateEndpointPlanAndEmptyResponse(t *testing.T) {
	for _, wrongPlan := range []bool{true, false} {
		t.Run(fmt.Sprint(wrongPlan), func(t *testing.T) {
			client := &fakeNextGenClient{get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
				v := &api.Nextgenv1beta2Tidb{}
				v.SetServicePlan(premiumPlan.servicePlan)
				if wrongPlan {
					v.SetServicePlan(essentialV2Plan.servicePlan)
				}
				return v, nil
			}, privateService: func(context.Context, string) (*api.Nextgenv1beta2PrivateLinkService, error) { return nil, nil }}
			cmd := privateEndpointCmd(helperWithNextGenClient(client), premiumPlan)
			cmd.SetArgs([]string{"service", "-c", "123"})
			require.Error(t, cmd.Execute())
		})
	}
}
