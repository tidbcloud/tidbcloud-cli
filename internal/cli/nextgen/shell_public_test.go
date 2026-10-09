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
	"fmt"
	"testing"
	"time"

	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/stretchr/testify/require"
)

func publicDNSFailureEndpoint(t *testing.T) api.TidbEndpoint {
	t.Helper()
	var endpoint api.TidbEndpoint
	// Keep the server's exact JSON shape: absent bools must not become false.
	require.NoError(t, json.Unmarshal([]byte(`{"connectionType":"PUBLIC","host":"sql.example.com","port":4000,"connectionReachability":{"reachable":false,"detail":{"serviceActive":true,"endpointActive":true,"dnsReachable":false,"message":"DNS_NOT_REACHABLE"}}}`), &endpoint))
	return endpoint
}

func TestShellPublicDNSRecoveryBoundaries(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, tc := range []struct {
			name    string
			mutate  func(*api.TidbEndpoint)
			wantGet bool
			wantCA  bool
			private bool
		}{
			{name: "DNS not reachable", wantGet: true, wantCA: true},
			{name: "DNS probe failed", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.SetMessage("DNS_PROBE_FAILED") }, wantGet: true, wantCA: true},
			{name: "unknown reason", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.SetMessage("allowlist required") }},
			{name: "reason with suffix", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.SetMessage("DNS_NOT_REACHABLE: example") }},
			{name: "missing reason", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.Message = nil }},
			{name: "service inactive", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.SetServiceActive(false) }},
			{name: "missing service state", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.ServiceActive = nil }},
			{name: "endpoint inactive", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.SetEndpointActive(false) }},
			{name: "missing endpoint state", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.EndpointActive = nil }},
			{name: "DNS reachable", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.SetDnsReachable(true) }},
			{name: "missing DNS state", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail.DnsReachable = nil }},
			{name: "missing detail", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Detail = nil }},
			{name: "missing reachability", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability = nil }, wantCA: true},
			{name: "missing reachable", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.Reachable = nil }, wantCA: true},
			{name: "reachable", mutate: func(e *api.TidbEndpoint) { e.ConnectionReachability.SetReachable(true) }, wantCA: true},
			{name: "empty host", mutate: func(e *api.TidbEndpoint) { e.SetHost("") }},
			{name: "host with URL", mutate: func(e *api.TidbEndpoint) { e.SetHost("https://sql.example.com") }},
			{name: "invalid port", mutate: func(e *api.TidbEndpoint) { e.SetPort(65536) }},
			{name: "zero port", mutate: func(e *api.TidbEndpoint) { e.SetPort(0) }},
			{name: "negative port", mutate: func(e *api.TidbEndpoint) { e.SetPort(-1) }},
			{name: "private without ready connection remains blocked", private: true, mutate: func(e *api.TidbEndpoint) { e.SetConnectionType(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT) }},
		} {
			t.Run(plan.commandName+"/"+tc.name, func(t *testing.T) {
				endpoint := publicDNSFailureEndpoint(t)
				if tc.mutate != nil {
					tc.mutate(&endpoint)
				}
				instance := api.NewNextgenv1beta2Tidb("test", "aws-us-west-2", plan.servicePlan)
				instance.State = api.V1BETA1CLUSTERSTATE_ACTIVE.Ptr()
				instance.Endpoints = []api.TidbEndpoint{endpoint}
				var calls []string
				caStop := errors.New("CA stage")
				client := &fakeNextGenClient{
					listPrivate: func(context.Context, string, int32, string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
						require.True(t, tc.private)
						return &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{}, nil
					},
					get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
					getPublic: func(ctx context.Context, id string) (*api.V1beta2PublicConnectionSetting, error) {
						calls = append(calls, "setting")
						require.Equal(t, "tidb-1", id)
						deadline, ok := ctx.Deadline()
						require.True(t, ok)
						require.InDelta(t, 30, time.Until(deadline).Seconds(), 1)
						setting := api.NewV1beta2PublicConnectionSetting()
						setting.SetEnabled(true)
						return setting, nil
					},
					certificate: func(context.Context, string) (*api.V1beta2CaCertificateDownloadUrl, error) {
						calls = append(calls, "CA")
						return nil, caStop
					},
				}
				h := helperWithNextGenClient(client)
				h.IOStreams.CanPrompt = true
				cmd := shellCmd(h, plan)
				args := []string{"-c", "tidb-1", "--password", "test-password"}
				if tc.private {
					args = append(args, "--connection-type", "private-endpoint")
				}
				cmd.SetArgs(args)
				err := cmd.Execute()
				require.Error(t, err)
				var expected []string
				if tc.wantGet {
					expected = append(expected, "setting")
					require.Equal(t, 1, bytes.Count(h.IOStreams.Err.(*bytes.Buffer).Bytes(), []byte("Warning:")))
				} else {
					require.Empty(t, h.IOStreams.Err.(*bytes.Buffer).String())
				}
				if tc.wantCA {
					expected = append(expected, "CA")
					require.ErrorIs(t, err, caStop)
				}
				require.Equal(t, expected, calls)
				require.NotContains(t, h.IOStreams.Out.(*bytes.Buffer).String(), "Warning:")
			})
		}
	}
}

func TestPublicDNSRecoveryRequiresExplicitEnabled(t *testing.T) {
	for _, payload := range []string{`{"enabled":true}`, `{"enabled":false}`, `{"enabled":null}`, `{}`, `null`} {
		t.Run(payload, func(t *testing.T) {
			var setting *api.V1beta2PublicConnectionSetting
			require.NoError(t, json.Unmarshal([]byte(payload), &setting))
			client := &fakeNextGenClient{getPublic: func(context.Context, string) (*api.V1beta2PublicConnectionSetting, error) { return setting, nil }}
			instance := &api.Nextgenv1beta2Tidb{Endpoints: []api.TidbEndpoint{publicDNSFailureEndpoint(t)}}
			endpoint, warning, err := resolvePublicEndpoint(context.Background(), client, "tidb-1", instance)
			if payload == `{"enabled":true}` {
				require.NoError(t, err)
				require.True(t, warning)
				require.Same(t, &instance.Endpoints[0], endpoint)
			} else {
				require.Error(t, err)
				require.Nil(t, endpoint)
				require.False(t, warning)
			}
		})
	}
	for _, requestErr := range []error{errors.New("401 Unauthorized"), errors.New("403 Forbidden"), errors.New("404 Not Found"), context.DeadlineExceeded, context.Canceled} {
		t.Run(requestErr.Error(), func(t *testing.T) {
			client := &fakeNextGenClient{getPublic: func(context.Context, string) (*api.V1beta2PublicConnectionSetting, error) { return nil, requestErr }}
			instance := &api.Nextgenv1beta2Tidb{Endpoints: []api.TidbEndpoint{publicDNSFailureEndpoint(t)}}
			endpoint, warning, err := resolvePublicEndpoint(context.Background(), client, "tidb-1", instance)
			require.ErrorIs(t, err, requestErr)
			require.ErrorContains(t, err, "cannot confirm PUBLIC access is enabled")
			require.Nil(t, endpoint)
			require.False(t, warning)
		})
	}
}

func TestShellPublicDNSRecoveryRetainsPlanAndStateGuards(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, wrongPlan := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s/wrong-plan=%t", plan.commandName, wrongPlan), func(t *testing.T) {
				instance := api.NewNextgenv1beta2Tidb("test", "aws-us-west-2", plan.servicePlan)
				instance.State = api.V1BETA1CLUSTERSTATE_CREATING.Ptr()
				if wrongPlan {
					instance.State = api.V1BETA1CLUSTERSTATE_ACTIVE.Ptr()
					instance.ServicePlan = api.V1BETA1SERVICEPLAN_STARTER
				}
				instance.Endpoints = []api.TidbEndpoint{publicDNSFailureEndpoint(t)}
				// Any request beyond GetTiDB panics, including the new setting GET.
				h := helperWithNextGenClient(&fakeNextGenClient{get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil }})
				h.IOStreams.CanPrompt = true
				cmd := shellCmd(h, plan)
				cmd.SetArgs([]string{"-c", "tidb-1", "--password", "test-password"})
				require.Error(t, cmd.Execute())
				require.Empty(t, h.IOStreams.Err.(*bytes.Buffer).String())
			})
		}
	}
}
