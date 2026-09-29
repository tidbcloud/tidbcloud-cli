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
	"strings"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/stretchr/testify/require"
)

func TestPasswordCommandResetsPremiumRootPassword(t *testing.T) {
	instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", "10000", api.V1BETA1SERVICEPLAN_PREMIUM)
	id := "tidb-1"
	instance.TidbId = &id
	var captured *api.TidbServiceResetRootPasswordBody
	client := &fakeNextGenClient{
		get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
		password: func(_ context.Context, actualID string, body *api.TidbServiceResetRootPasswordBody) error {
			require.Equal(t, id, actualID)
			captured = body
			return nil
		},
	}
	command := passwordCmd(helperWithNextGenClient(client), premiumPlan)
	command.SetArgs([]string{"--cluster-id", id, "--password", "new-password"})

	require.NoError(t, command.ExecuteContext(context.Background()))
	require.NotNil(t, captured)
	require.Equal(t, "new-password", captured.GetRootPassword())
}

func TestPasswordCommandSupportsInteractivePasswordWithExplicitInstance(t *testing.T) {
	instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", "10000", api.V1BETA1SERVICEPLAN_ESSENTIAL_V2)
	id := "tidb-1"
	instance.TidbId = &id
	client := &fakeNextGenClient{
		get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
		password: func(_ context.Context, actualID string, body *api.TidbServiceResetRootPasswordBody) error {
			require.Equal(t, id, actualID)
			require.Equal(t, "new-password", body.GetRootPassword())
			return nil
		},
	}
	h := helperWithNextGenClient(client)
	h.IOStreams.CanPrompt = true
	command := passwordCmdWithDependencies(
		h,
		essentialV2Plan,
		func(context.Context, cloud.NextGenClient, planSpec, int32) (*api.Nextgenv1beta2Tidb, error) {
			t.Fatal("instance selector must not be called when --cluster-id is set")
			return nil, nil
		},
		func() (string, error) { return "new-password", nil },
	)
	command.SetArgs([]string{"--cluster-id", id})

	require.NoError(t, command.ExecuteContext(context.Background()))
	require.Equal(t, "true", command.Annotations[telemetry.InteractiveMode])
}

func TestPasswordCommandRejectsInvalidPasswordBeforeRequest(t *testing.T) {
	clientCalls := 0
	h := helperWithNextGenClient(&fakeNextGenClient{})
	h.NextGenClient = func() (cloud.NextGenClient, error) {
		clientCalls++
		return &fakeNextGenClient{}, nil
	}
	command := passwordCmd(h, premiumPlan)
	command.SetArgs([]string{"--cluster-id", "tidb-1", "--password", "short"})

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "root password must be between 8 and 64 characters")
	require.Zero(t, clientCalls)
}

func TestPasswordCommandRequiresInstanceForNonInteractiveUse(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "omitted", args: []string{"--password", "new-password"}},
		{name: "empty", args: []string{"--cluster-id=", "--password", "new-password"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientCalls := 0
			h := helperWithNextGenClient(&fakeNextGenClient{})
			h.NextGenClient = func() (cloud.NextGenClient, error) {
				clientCalls++
				return &fakeNextGenClient{}, nil
			}
			command := passwordCmd(h, premiumPlan)
			command.SetArgs(tt.args)

			err := command.ExecuteContext(context.Background())
			require.ErrorContains(t, err, `--cluster-id is required when --password is specified`)
			require.Zero(t, clientCalls)
		})
	}
}

func TestPasswordCommandRejectsCrossPlanInstance(t *testing.T) {
	instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", "10000", api.V1BETA1SERVICEPLAN_ESSENTIAL_V2)
	id := "tidb-1"
	instance.TidbId = &id
	resetCalls := 0
	client := &fakeNextGenClient{
		get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
		password: func(context.Context, string, *api.TidbServiceResetRootPasswordBody) error {
			resetCalls++
			return nil
		},
	}
	command := passwordCmd(helperWithNextGenClient(client), premiumPlan)
	command.SetArgs([]string{"--cluster-id", id, "--password", "new-password"})

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, `only manages "Premium" instances`)
	require.Zero(t, resetCalls)
}

func TestValidateRootPasswordAcceptsEncryptedPayload(t *testing.T) {
	require.NoError(t, validateRootPassword("rsa_oaep_sha256:ciphertext"))
	require.Error(t, validateRootPassword("rsa_oaep_sha256:"))
}

func TestValidateRootPasswordMatchesBackendByteLimit(t *testing.T) {
	for _, tc := range []struct {
		name     string
		password string
		valid    bool
	}{
		{name: "below minimum", password: "1234567"},
		{name: "minimum", password: "12345678", valid: true},
		{name: "maximum", password: strings.Repeat("a", 64), valid: true},
		{name: "above maximum", password: strings.Repeat("a", 65)},
		{name: "multibyte below eight runes", password: "测试密码", valid: true},
		{name: "multibyte maximum", password: strings.Repeat("密", 21) + "a", valid: true},
		{name: "multibyte above maximum", password: strings.Repeat("密", 30)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := validateRootPassword(tc.password)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
