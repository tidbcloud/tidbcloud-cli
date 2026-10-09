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
	"errors"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/stretchr/testify/require"
)

func TestCreateCMEK(t *testing.T) {
	for _, tc := range []struct {
		region, arn, principal, key string
	}{
		{"aws-us-west-2", "arn:aws:kms:us-west-2:123456789012:key/test", `{"servicePlan":"Premium","aws":{"accountId":"123456789012","externalId":"org-1"}}`, `"awsKms":{"kmsKeyArn":"arn:aws:kms:us-west-2:123456789012:key/test"},"awsPrincipal":{"accountId":"123456789012","externalId":"org-1"}`},
		{"regions/alicloud-cn-hongkong", "acs:kms:cn-hongkong:123456789012:key/test", `{"servicePlan":"Premium","aliyun":{"accountId":"123456789012"}}`, `"aliyunKms":{"kmsKeyArn":"acs:kms:cn-hongkong:123456789012:key/test"},"aliyunPrincipal":{"accountId":"123456789012"}`},
	} {
		t.Run(tc.region, func(t *testing.T) {
			var calls []string
			var verified *api.V1beta2CustomerManagedEncryptionKey
			client := &fakeNextGenClient{
				principal: func(_ context.Context, region string) (*api.V1beta2CmekAccessIamPrincipal, error) {
					calls = append(calls, "principal")
					require.Equal(t, tc.region, region)
					var result api.V1beta2CmekAccessIamPrincipal
					require.NoError(t, json.Unmarshal([]byte(tc.principal), &result))
					return &result, nil
				},
				verifyCMEK: func(_ context.Context, key *api.V1beta2CustomerManagedEncryptionKey) (*api.V1beta2VerifyCmekAccessIamPrincipalResponse, error) {
					calls = append(calls, "verify")
					verified = key
					payload, err := json.Marshal(key)
					require.NoError(t, err)
					require.JSONEq(t, `{"regionId":"`+tc.region+`","servicePlan":"Premium",`+tc.key+`}`, string(payload))
					return &api.V1beta2VerifyCmekAccessIamPrincipalResponse{Valid: api.PtrBool(true)}, nil
				},
				create: func(_ context.Context, body *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error) {
					calls = append(calls, "create")
					require.Same(t, verified, body.DualLayerDataEncryption.Cmek)
					require.Nil(t, body.DualLayerDataEncryption.DefaultKey)
					require.Nil(t, body.DualLayerDataEncryption.NoEncryption)
					require.Equal(t, "4078", (*body.Labels)[projectIDLabel])
					return body, nil
				},
			}
			cmd := PremiumCmd(helperWithNextGenClient(client))
			args, err := NormalizeActionArgs(actionTestRoot(), []string{"premium", "--create", "--display-name", "cmek-test", "--region", tc.region, "--max-rcu", "5000", "--project-id", "4078", "--encryption", "cmek", "--cmek-key-arn", tc.arn, "-o", "json"})
			require.NoError(t, err)
			cmd.SetArgs(args[1:])
			require.NoError(t, cmd.Execute())
			require.Equal(t, []string{"principal", "verify", "create"}, calls)
		})
	}
}

func TestCMEKInputRejectedBeforeClient(t *testing.T) {
	for _, tc := range []struct {
		name string
		plan planSpec
		args []string
		want string
	}{
		{"missing ARN", premiumPlan, []string{"--encryption", "cmek"}, "--cmek-key-arn is required"},
		{"conflicting ARN", premiumPlan, []string{"--cmek-key-arn", "key"}, "requires --encryption cmek"},
		{"empty conflicting ARN", premiumPlan, []string{"--encryption", "default-key", "--cmek-key-arn="}, "requires --encryption cmek"},
		{"unsupported plan", essentialV2Plan, []string{"--encryption", "cmek"}, "only supported for Premium"},
		{"unsupported provider", premiumPlan, []string{"--encryption", "cmek", "--cmek-key-arn", "key", "--region", "gcp-us-central1"}, "AWS or Alibaba Cloud region"},
		{"invalid output", premiumPlan, []string{"--encryption", "cmek", "--cmek-key-arn", "key", "-o", "xml"}, "unsupported output format"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := createCmd(&internal.Helper{}, tc.plan) // A client lookup would panic.
			cmd.SetArgs(append([]string{"--display-name", "cmek-test", "--region", "aws-us-west-2", "--max-rcu", "5000"}, tc.args...))
			require.ErrorContains(t, cmd.Execute(), tc.want)
		})
	}
}

func TestCMEKFailureStopsCreate(t *testing.T) {
	for _, failure := range []string{"principal error", "missing principal", "verify error", "not valid"} {
		t.Run(failure, func(t *testing.T) {
			sentinel := errors.New(failure)
			client := &fakeNextGenClient{
				principal: func(context.Context, string) (*api.V1beta2CmekAccessIamPrincipal, error) {
					if failure == "principal error" {
						return nil, sentinel
					}
					if failure == "missing principal" {
						return &api.V1beta2CmekAccessIamPrincipal{}, nil
					}
					return &api.V1beta2CmekAccessIamPrincipal{Aws: &api.CmekAccessIamPrincipalAwsCmekPrincipal{AccountId: api.PtrString("123")}}, nil
				},
				verifyCMEK: func(context.Context, *api.V1beta2CustomerManagedEncryptionKey) (*api.V1beta2VerifyCmekAccessIamPrincipalResponse, error) {
					if failure == "verify error" {
						return nil, sentinel
					}
					return &api.V1beta2VerifyCmekAccessIamPrincipalResponse{FailReason: api.PtrString("key policy denies access")}, nil
				},
			}
			cmd := createCmd(helperWithNextGenClient(client), premiumPlan)
			cmd.SetArgs([]string{"--display-name", "cmek-test", "--region", "aws-us-west-2", "--max-rcu", "5000", "--encryption", "cmek", "--cmek-key-arn", "key"})
			err := cmd.Execute() // An unexpected create request would panic.
			if failure == "principal error" || failure == "verify error" {
				require.ErrorIs(t, err, sentinel)
			} else if failure == "not valid" {
				require.ErrorContains(t, err, "key policy denies access")
			} else {
				require.ErrorContains(t, err, "no AWS CMEK principal")
			}
		})
	}
}

func TestCreateExistingEncryptionUnchanged(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, encryption := range []string{"", "none", "default-key"} {
			t.Run(plan.commandName+"/"+encryption, func(t *testing.T) {
				calls := 0
				client := &fakeNextGenClient{create: func(_ context.Context, body *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error) {
					calls++
					payload, err := json.Marshal(body.DualLayerDataEncryption)
					require.NoError(t, err)
					want := `{"noEncryption":{}}`
					if encryption == "default-key" {
						want = `{"defaultKey":{}}`
					}
					require.JSONEq(t, want, string(payload))
					require.Nil(t, body.Labels)
					return body, nil
				}}
				cmd := createCmd(helperWithNextGenClient(client), plan)
				args := []string{"--display-name", "plain-test", "--region", "aws-us-west-2", "--max-rcu", "5000"}
				if encryption != "" {
					args = append(args, "--encryption", encryption)
				}
				cmd.SetArgs(args)
				require.NoError(t, cmd.Execute()) // CMEK methods are unset and must never be used.
				require.Equal(t, 1, calls)
			})
		}
	}
}

func TestCMEKPrincipalCommand(t *testing.T) {
	streams := iostream.Test()
	client := &fakeNextGenClient{principal: func(_ context.Context, region string) (*api.V1beta2CmekAccessIamPrincipal, error) {
		require.Equal(t, "aws-us-west-2", region)
		return &api.V1beta2CmekAccessIamPrincipal{ServicePlan: api.V1BETA1SERVICEPLAN_PREMIUM, Aws: &api.CmekAccessIamPrincipalAwsCmekPrincipal{AccountId: api.PtrString("123"), ExternalId: api.PtrString("org-1")}}, nil
	}}
	h := helperWithNextGenClient(client)
	h.IOStreams = streams
	cmd := PremiumCmd(h)
	cmd.SetArgs([]string{"cmek", "principal", "--region", "aws-us-west-2", "-o", "json"})
	require.NoError(t, cmd.Execute())
	require.Contains(t, streams.Out.(*bytes.Buffer).String(), `"externalId": "org-1"`)
	for _, child := range EssentialCmd(h).Commands() {
		require.NotEqual(t, "cmek", child.Name())
	}
	require.NotContains(t, EssentialCmd(h).Long, "cmek")
}
