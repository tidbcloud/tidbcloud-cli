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
	"errors"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestPublicEndpointEnableAndDisable(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, enabled := range []bool{true, false} {
			action := "disable"
			if enabled {
				action = "enable"
			}
			t.Run(plan.commandName+"/"+action, func(t *testing.T) {
				id := "tidb-1"
				instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", "5000", plan.servicePlan)
				instance.TidbId = &id
				calls := 0
				client := &fakeNextGenClient{
					get: func(_ context.Context, actualID string) (*api.Nextgenv1beta2Tidb, error) {
						require.Equal(t, id, actualID)
						return instance, nil
					},
					public: func(_ context.Context, actualID string, body *api.V1beta2PublicConnectionSetting) (*api.V1beta2PublicConnectionSetting, error) {
						require.Equal(t, id, actualID)
						require.True(t, body.HasEnabled())
						require.Equal(t, enabled, body.GetEnabled())
						payload, err := body.MarshalJSON()
						require.NoError(t, err)
						expected := `{"enabled":true}`
						if !enabled {
							expected = `{"enabled":false}`
						}
						require.JSONEq(t, expected, string(payload))
						calls++
						// The response includes settings preserved by the server.
						body.IpAccessList = []api.V1beta2PublicConnectionSettingIpAccessList{*api.NewV1beta2PublicConnectionSettingIpAccessList("192.0.2.0/24")}
						return body, nil
					},
				}
				h := helperWithNextGenClient(client)
				command := &cobra.Command{Use: "ticloud"}
				command.AddCommand(PremiumCmd(h), EssentialCmd(h))
				args := []string{plan.commandName, "public-endpoint", action, "--cluster-id", id, "--output", "json"}
				if !enabled {
					args = append(args, "--force")
					if plan.commandName == "essential-v2" {
						args[0] = "essential" // Preserve the existing alias.
					}
				}
				args, err := NormalizeActionArgs(args)
				require.NoError(t, err)
				command.SetArgs(args)

				require.NoError(t, command.ExecuteContext(context.Background()))
				require.Equal(t, 1, calls)
				out := h.IOStreams.Out.(*bytes.Buffer).String()
				expected := `{"enabled":true,"ipAccessList":[{"cidrNotation":"192.0.2.0/24"}]}`
				if !enabled {
					expected = `{"enabled":false,"ipAccessList":[{"cidrNotation":"192.0.2.0/24"}]}`
				}
				require.JSONEq(t, expected, out)
			})
		}
	}
}

func TestPublicEndpointInteractiveSelectionReadsCurrentInstance(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, changedPlan := range []bool{false, true} {
			name := plan.commandName + "/matching plan"
			if changedPlan {
				name = plan.commandName + "/changed plan"
			}
			t.Run(name, func(t *testing.T) {
				instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", "5000", plan.servicePlan)
				instance.SetTidbId("tidb-1")
				var calls []string
				client := &fakeNextGenClient{
					get: func(_ context.Context, id string) (*api.Nextgenv1beta2Tidb, error) {
						require.Equal(t, "tidb-1", id)
						calls = append(calls, "get")
						current := *instance
						if changedPlan {
							current.ServicePlan = api.V1BETA1SERVICEPLAN_STARTER
						}
						return &current, nil
					},
					public: func(_ context.Context, id string, body *api.V1beta2PublicConnectionSetting) (*api.V1beta2PublicConnectionSetting, error) {
						require.Equal(t, "tidb-1", id)
						calls = append(calls, "patch")
						return body, nil
					},
				}
				h := helperWithNextGenClient(client)
				h.IOStreams.CanPrompt = true
				command := publicEndpointActionCmd(h, plan, true,
					func(_ context.Context, actualClient cloud.NextGenClient, actualPlan planSpec, pageSize int32) (*api.Nextgenv1beta2Tidb, error) {
						require.Same(t, client, actualClient)
						require.Equal(t, plan, actualPlan)
						require.Equal(t, int32(h.QueryPageSize), pageSize)
						calls = append(calls, "select")
						return instance, nil
					}, nil)
				err := command.ExecuteContext(context.Background())
				if changedPlan {
					require.ErrorContains(t, err, "only manages")
					require.Equal(t, []string{"select", "get"}, calls)
				} else {
					require.NoError(t, err)
					require.Equal(t, []string{"select", "get", "patch"}, calls)
					require.Contains(t, h.IOStreams.Out.(*bytes.Buffer).String(), "Enable request accepted")
				}
			})
		}
	}
}

func TestPublicEndpointRejectsInvalidInputBeforeRequest(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, tc := range []struct {
			name      string
			enabled   bool
			canPrompt bool
			args      []string
			want      string
		}{
			{name: "missing ID", enabled: true, want: "use --cluster-id to specify the instance"},
			{name: "empty explicit ID", enabled: true, canPrompt: true, args: []string{"-c", ""}, want: "--cluster-id must not be empty"},
			{name: "disable needs confirmation", args: []string{"-c", "tidb-1"}, want: "use --force to disable the public endpoint"},
			{name: "enable output", enabled: true, args: []string{"-c", "tidb-1", "-o", "xml"}, want: "unsupported output format: xml"},
			{name: "disable output", args: []string{"-c", "tidb-1", "--force", "-o", "xml"}, want: "unsupported output format: xml"},
		} {
			t.Run(plan.commandName+"/"+tc.name, func(t *testing.T) {
				h := helperWithNextGenClient(&fakeNextGenClient{})
				h.IOStreams.CanPrompt = tc.canPrompt
				h.NextGenClient = func() (cloud.NextGenClient, error) {
					t.Fatal("invalid input must be rejected before obtaining credentials/client")
					return nil, nil
				}
				command := publicEndpointActionCmd(h, plan, tc.enabled, nil, nil)
				command.SetArgs(tc.args)
				require.ErrorContains(t, command.ExecuteContext(context.Background()), tc.want)
			})
		}
	}
}

func TestPublicEndpointDisableConfirmation(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, tc := range []struct {
			name  string
			force bool
			err   error
		}{
			{name: "confirmed"},
			{name: "declined", err: errors.New("incorrect confirmation")},
			{name: "interrupted", err: util.InterruptError},
			{name: "force", force: true},
		} {
			t.Run(plan.commandName+"/"+tc.name, func(t *testing.T) {
				client := &fakeNextGenClient{
					get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
						return api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", "5000", plan.servicePlan), nil
					},
				}
				patches, confirmations := 0, 0
				client.public = func(_ context.Context, id string, body *api.V1beta2PublicConnectionSetting) (*api.V1beta2PublicConnectionSetting, error) {
					patches++
					require.Equal(t, "tidb-1", id)
					require.False(t, body.GetEnabled())
					return body, nil
				}
				h := helperWithNextGenClient(client)
				h.IOStreams.CanPrompt = true
				command := publicEndpointActionCmd(h, plan, false, nil, func(id string) error {
					confirmations++
					require.Equal(t, "tidb-1", id)
					return tc.err
				})
				args := []string{"-c", "tidb-1"}
				if tc.force {
					args = append(args, "--force")
				}
				command.SetArgs(args)
				err := command.ExecuteContext(context.Background())
				if tc.err != nil {
					require.ErrorIs(t, err, tc.err)
					require.Zero(t, patches)
					require.Empty(t, h.IOStreams.Out.(*bytes.Buffer).String())
				} else {
					require.NoError(t, err)
					require.Equal(t, 1, patches)
				}
				if tc.force {
					require.Zero(t, confirmations)
				} else {
					require.Equal(t, 1, confirmations)
				}
			})
		}
	}
}

func TestPublicEndpointStopsOnErrors(t *testing.T) {
	failure := errors.New("request failed")
	for _, stage := range []string{"selection", "get", "update", "nil response"} {
		t.Run(stage, func(t *testing.T) {
			patches := 0
			instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", "5000", premiumPlan.servicePlan)
			instance.SetTidbId("tidb-1")
			client := &fakeNextGenClient{
				get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
					require.NotEqual(t, "selection", stage)
					if stage == "get" {
						return nil, failure
					}
					return instance, nil
				},
				public: func(context.Context, string, *api.V1beta2PublicConnectionSetting) (*api.V1beta2PublicConnectionSetting, error) {
					patches++
					if stage == "nil response" {
						return nil, nil
					}
					return nil, failure
				},
			}
			h := helperWithNextGenClient(client)
			h.IOStreams.CanPrompt = true
			command := publicEndpointActionCmd(h, premiumPlan, true,
				func(context.Context, cloud.NextGenClient, planSpec, int32) (*api.Nextgenv1beta2Tidb, error) {
					if stage == "selection" {
						return nil, util.InterruptError
					}
					return instance, nil
				}, nil)
			err := command.ExecuteContext(context.Background())
			switch stage {
			case "selection":
				require.ErrorIs(t, err, util.InterruptError)
				require.Zero(t, patches)
			case "get":
				require.ErrorIs(t, err, failure)
				require.Zero(t, patches)
			case "update":
				require.ErrorIs(t, err, failure)
				require.Equal(t, 1, patches)
			case "nil response":
				require.ErrorContains(t, err, "empty public connection setting")
				require.Equal(t, 1, patches)
			}
			require.Empty(t, h.IOStreams.Out.(*bytes.Buffer).String())
		})
	}
}
