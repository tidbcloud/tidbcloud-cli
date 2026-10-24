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
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"
)

func rcuAPIError(status int, message string) *cloud.NextGenAPIError {
	return &cloud.NextGenAPIError{Status: status, RequestID: "request-1", Method: "POST", Path: "/v1beta2/tidbs",
		Response: map[string]interface{}{"code": 400, "message": message}}
}

func TestRCUErrorHintOnlyForMinimumRejection(t *testing.T) {
	for _, message := range []string{
		"Min RCU must be less than or equal to max RCU.",
		"min RCU must be less than or equal to max RCU",
		"  MIN RCU MUST BE LESS THAN OR EQUAL TO MAX RCU.  ",
	} {
		for _, input := range []int64{100, 20000, 40000} {
			err := rcuAPIError(400, message)
			wrapped := explainRCUError(err, input, premiumPlan)
			require.ErrorIs(t, wrapped, err)
			require.Contains(t, wrapped.Error(), fmt.Sprintf("--max-rcu=%d", input))
			require.Contains(t, wrapped.Error(), "default MAX_RCU minimum for TiDB Cloud Premium is 5000 RCU")
			require.Contains(t, wrapped.Error(), "server did not return the actual minimum")
			require.Contains(t, wrapped.Error(), "organization-specific limits may differ")
			require.NotContains(t, wrapped.Error(), "max-rcu / 4")
			require.NotContains(t, wrapped.Error(), "default calculation")
			require.Contains(t, wrapped.Error(), "no usable default")
			text, ok := marshalJSONError(wrapped)
			require.True(t, ok)
			var result map[string]map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(text), &result))
			require.Equal(t, message, result["error"]["message"])
			require.Contains(t, result["error"]["hint"], "organization-specific limits may differ")

		}
	}

	for _, err := range []error{
		errors.New("Min RCU must be less than or equal to max RCU."),
		rcuAPIError(403, "Min RCU must be less than or equal to max RCU."),
		rcuAPIError(400, "Other validation error"),
		rcuAPIError(400, "Min RCU must be less than or equal to max RCU. More detail"),
		rcuAPIError(400, "Min RCU must be less than or equal to max RCU.."),
		rcuAPIError(403, "min RCU must be less than or equal to max RCU"),
		rcuAPIError(400, "Minimum RCU is 2000"),
	} {
		require.Same(t, err, explainRCUError(err, 100, premiumPlan))
	}
}

func TestPlanCreateAndUpdateRCUErrors(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, action := range []string{"create", "update"} {
			t.Run(plan.commandName+"/"+action, func(t *testing.T) {
				apiErr := rcuAPIError(400, "Min RCU must be less than or equal to max RCU.")
				calls := 0
				client := &fakeNextGenClient{
					create: func(context.Context, *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error) {
						calls++
						return nil, apiErr
					},
					get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
						return api.NewNextgenv1beta2Tidb("test", "aws-us-west-2", plan.servicePlan), nil
					},
					update: func(context.Context, string, *api.TheTiDBCloudPremiumInstanceToUpdate) (*api.Nextgenv1beta2Tidb, error) {
						calls++
						return nil, apiErr
					},
				}
				root := cmd(helperWithNextGenClient(client), plan)
				root.SilenceErrors, root.SilenceUsage = true, true
				args := []string{action, "--max-rcu", "100", "-ojson"}
				if action == "create" {
					args = append(args, "--display-name", "test", "--region", "aws-us-west-2")
				} else {
					args = append(args, "-c", "test-id")
				}
				root.SetArgs(args)
				command, err := root.ExecuteC()
				require.ErrorIs(t, err, apiErr)
				require.Equal(t, 1, calls, "no retry or extra minimum lookup")
				text, ok := JSONError(command, args, err)
				require.True(t, ok)
				var result map[string]map[string]interface{}
				require.NoError(t, json.Unmarshal([]byte(text), &result))
				require.Equal(t, float64(400), result["error"]["status"])
				require.Equal(t, "request-1", result["error"]["requestId"])
				require.Equal(t, "Min RCU must be less than or equal to max RCU.", result["error"]["message"])
				require.Contains(t, result["error"]["hint"], "organization-specific limits may differ")
				require.Contains(t, result["error"]["hint"], "server did not return the actual minimum")
				expectedMinimum := "5000"
				if plan.commandName == "essential-v2" {
					expectedMinimum = "2000"
				}
				require.Contains(t, result["error"]["hint"], "is "+expectedMinimum+" RCU")
				require.Contains(t, err.Error(), "is "+expectedMinimum+" RCU")
				require.NotContains(t, text, "max-rcu / 4")
				require.NotContains(t, text, "default calculation")
			})
		}
	}
}

func TestJSONErrorValidationAndScope(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, args := range [][]string{
			{"create", "-o", "json"}, // missing required flags
			{"create", "--max-rcu", "not-a-number", "--output=json"},
			{"update", "-c", "id", "-ojson"}, // PreRunE
			{"create", "--unknown-flag", "-ojson"},
		} {
			root := cmd(helperWithNextGenClient(nil), plan)
			root.SilenceErrors, root.SilenceUsage = true, true
			root.SetArgs(args)
			command, err := root.ExecuteC()
			require.Error(t, err)
			text, ok := JSONError(command, args, err)
			require.True(t, ok, "%v", args)
			require.True(t, json.Valid([]byte(text)))
			require.NotContains(t, text, "Error:")
		}
	}
	root := &cobra.Command{Use: "ticloud"}
	root.AddCommand(cmd(helperWithNextGenClient(nil), premiumPlan))
	legacy := &cobra.Command{Use: "serverless"}
	legacy.Flags().StringP("output", "o", "human", "")
	root.AddCommand(legacy)
	for _, args := range [][]string{{"serverless", "-ojson"}, {"premium", "shell"}, {"premium", "password"}, {"premium", "create", "-o", "human"}} {
		command, _, err := root.Find(args)
		require.NoError(t, err)
		_, ok := JSONError(command, args, errors.New("failure"))
		require.False(t, ok, "%v", args)
	}
}

func TestJSONOutputSelectionBeforeExecution(t *testing.T) {
	root := &cobra.Command{Use: "ticloud"}
	root.PersistentFlags().StringP("profile", "P", "", "")
	root.AddCommand(PremiumCmd(helperWithNextGenClient(nil)), EssentialCmd(helperWithNextGenClient(nil)))
	for _, args := range [][]string{
		{"-P", "test", "premium", "--create", "--max-rcu", "100", "-ojson"},
		{"essential", "--output=json", "--create"},
		{"premium", "public-endpoint", "enable", "-o", "json"},
		{"premium", "cmek", "principal", "--output=json"},
	} {
		normalized, err := NormalizeActionArgs(root, args)
		require.NoError(t, err)
		command, _, err := root.Find(normalized)
		require.NoError(t, err)
		require.True(t, UsesJSONOutput(command, normalized), "%v", args)
		require.Equal(t, "human", command.Flags().Lookup("output").Value.String(), "selection must not mutate flags")
	}
}

func TestJSONActionNormalizationErrors(t *testing.T) {
	for _, args := range [][]string{
		{"premium", "--create", "--list", "-ojson"},
		{"essential", "--create=false", "--output=json"},
	} {
		root := &cobra.Command{Use: "ticloud"}
		root.AddCommand(PremiumCmd(helperWithNextGenClient(nil)), EssentialCmd(helperWithNextGenClient(nil)))
		_, err := NormalizeActionArgs(root, args)
		require.Error(t, err)
		command, _, _ := root.Find(args)
		text, ok := JSONError(command, args, err)
		require.True(t, ok)
		require.True(t, json.Valid([]byte(text)))
	}
}

func TestJSONTransportErrorRedaction(t *testing.T) {
	original := &url.Error{Op: "Post", URL: "https://user:secret-pass@example.com/api?token=secret-token", Err: errors.New("timeout")}
	err := fmt.Errorf("change may have been applied; check status before retrying: %w", original)
	text, ok := marshalJSONError(err)
	require.True(t, ok)
	require.True(t, json.Valid([]byte(text)))
	require.NotContains(t, text, "secret-pass")
	require.NotContains(t, text, "secret-token")
	require.Contains(t, text, "check status before retrying")
	require.Contains(t, original.Error(), "secret-token", "do not mutate errors used by shell/legacy callers")
}

func TestRCUHelpUsesPlanMinimum(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		t.Run(plan.commandName, func(t *testing.T) {
			h := helperWithNextGenClient(nil)
			for _, help := range []string{cmd(h, plan).Long, createCmd(h, plan).UsageString()} {
				expectedMinimum := "5000"
				if plan.commandName == "essential-v2" {
					expectedMinimum = "2000"
				}
				require.Contains(t, help, "Default minimum: "+expectedMinimum+" RCU")
				require.Contains(t, help, "organization-specific limits may differ")
				require.Contains(t, help, "No CLI default")

				require.NotContains(t, help, "max-rcu / 4")
				require.NotContains(t, help, "default calculation")
			}
		})
	}
}

func TestPrivateEndpointDefaultJSONErrors(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, action := range []string{"service", "list", "bind", "describe", "unbind"} {
			for _, format := range []string{"", "json", "human"} {
				t.Run(plan.commandName+"/"+action+"/"+format, func(t *testing.T) {
					root := cmd(helperWithNextGenClient(nil), plan)
					args := []string{"private-endpoint", action, "-c", "invalid-id"}
					if format != "" {
						args = append(args, "-o", format)
					}
					selected, _, err := root.Find(args)
					require.NoError(t, err)
					require.Equal(t, format != "human", UsesJSONOutput(selected, args))
					require.Equal(t, "json", selected.Flags().Lookup("output").Value.String(), "must not mutate flags before execution")
					root.SilenceErrors, root.SilenceUsage = true, true
					root.SetArgs(args)
					executed, err := root.ExecuteC()
					require.Error(t, err)
					text, ok := JSONError(executed, args, err)
					require.Equal(t, format != "human", ok)
					if ok {
						require.True(t, json.Valid([]byte(text)))
					}
				})
			}
		}
	}
}

func TestExistingCommandsKeepDefaultErrorFormat(t *testing.T) {
	root := &cobra.Command{Use: "ticloud"}
	root.AddCommand(PremiumCmd(helperWithNextGenClient(nil)), EssentialCmd(helperWithNextGenClient(nil)))
	legacy := &cobra.Command{Use: "serverless"}
	legacy.Flags().StringP("output", "o", "json", "")
	root.AddCommand(legacy)
	for _, args := range [][]string{{"premium", "describe"}, {"premium", "shell"}, {"premium", "password"}, {"premium", "public-endpoint", "enable"}, {"essential-v2", "describe"}, {"serverless"}} {
		selected, _, err := root.Find(args)
		require.NoError(t, err)
		require.False(t, UsesJSONOutput(selected, args), "%v", args)
	}
}

func TestShellHelpWarnsAboutPasswordArguments(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		command := shellCmd(helperWithNextGenClient(nil), plan)
		require.Contains(t, command.Flags().Lookup("password").Usage, "shell history")
		require.Contains(t, command.Flags().Lookup("password").Usage, "hidden prompt")
	}
}
