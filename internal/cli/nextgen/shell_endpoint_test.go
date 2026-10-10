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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/stretchr/testify/require"
)

func TestShellEndpointReachabilityFromJSON(t *testing.T) {
	for _, connectionType := range []api.EndpointConnectionType{api.ENDPOINTCONNECTIONTYPE_PUBLIC, api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, api.ENDPOINTCONNECTIONTYPE_VPC_PEERING} {
		for _, tc := range []struct {
			name, fields string
			wantError    bool
		}{
			{name: "unknown"},
			{name: "reachable", fields: `,"connectionReachability":{"reachable":true}`},
			{name: "unreachable", fields: `,"connectionReachability":{"reachable":false}`, wantError: true},
		} {
			t.Run(string(connectionType)+"/"+tc.name, func(t *testing.T) {
				var endpoint api.TidbEndpoint
				payload := fmt.Sprintf(`{"host":"sql.example.com","port":4000,"connectionType":%q%s}`, connectionType, tc.fields)
				require.NoError(t, json.Unmarshal([]byte(payload), &endpoint))
				instance := &api.Nextgenv1beta2Tidb{Endpoints: []api.TidbEndpoint{endpoint}}
				var selected *api.TidbEndpoint
				var err error
				if connectionType == api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT {
					selected, err = resolvePrivateEndpoint(instance.Endpoints, "", nil)
				} else if connectionType == api.ENDPOINTCONNECTIONTYPE_PUBLIC {
					selected, _, err = resolvePublicEndpoint(context.Background(), nil, "tidb-1", instance)
				} else {
					selected, err = validateEndpoint(&endpoint)
				}
				if tc.wantError {
					require.ErrorContains(t, err, "endpoint is not reachable")
					return
				}
				require.NoError(t, err)
				require.Equal(t, "sql.example.com", selected.GetHost())
			})
		}
	}
}

func TestParseEndpointAddress(t *testing.T) {
	for _, address := range []string{"private.example.com:4000", "[2001:db8::1]:4000"} {
		parsed, err := parseEndpointAddress(address)
		require.NoError(t, err)
		require.Equal(t, address, parsed)
	}
	parsed, err := parseEndpointAddress("private.example.com:04000")
	require.NoError(t, err)
	require.Equal(t, "private.example.com:4000", parsed)
	for _, address := range []string{"", "host", ":4000", "host:", "host:0", "host:65536", "host:-1", "host:mysql", "https://host:4000", "host/path:4000", "user@host:4000", "host?query:4000", "host#fragment:4000", "host :4000"} {
		t.Run(address, func(t *testing.T) {
			_, err := parseEndpointAddress(address)
			require.Error(t, err)
		})
	}
}

func TestShellRejectsNetworkOptionsBeforeClient(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, args := range [][]string{
			{"--connection-type", "vpc-peering"},
			{"--connection-type", ""},
			{"--endpoint", "host:4000"},
			{"--connection-type", "public", "--endpoint", "host:4000"},
			{"--connection-type", "private-endpoint", "--endpoint", ""},
			{"--connection-type", "private-endpoint", "--endpoint", "https://host:4000"},
			{"--connection-type", "private-endpoint", "--endpoint", "host:65536"},
			{"--connection-type", "private-endpoint", "--password", "secret"},
			{"--connection-type", "private-endpoint", "-u", "user"},
		} {
			h := helperWithNextGenClient(&fakeNextGenClient{})
			h.IOStreams.CanPrompt = true
			h.NextGenClient = func() (cloud.NextGenClient, error) {
				t.Fatal("invalid options must be rejected before obtaining a client")
				return nil, nil
			}
			command := shellCmd(h, plan)
			command.SetArgs(args)
			require.Error(t, command.Execute())
		}
	}
}

func TestShellNetworkFlagsPreserveInteractionMode(t *testing.T) {
	for _, tt := range []struct {
		args        []string
		interactive bool
	}{
		{nil, true},
		{[]string{"--connection-type", "private-endpoint"}, true},
		{[]string{"--connection-type", "private-endpoint", "--endpoint", "private.example.com:4000"}, true},
		{[]string{"-c", "tidb-1", "--connection-type", "private-endpoint"}, false},
		{[]string{"--password", "secret", "--connection-type", "private-endpoint"}, false},
		{[]string{"-u", "user", "--connection-type", "private-endpoint"}, false},
	} {
		command := shellCmd(helperWithNextGenClient(&fakeNextGenClient{}), premiumPlan)
		require.NoError(t, command.ParseFlags(tt.args))
		require.Equal(t, tt.interactive, shellIsInteractive(command))
	}
}

func testShellEndpoint(kind api.EndpointConnectionType, host string, port int32) api.TidbEndpoint {
	return api.TidbEndpoint{ConnectionType: &kind, Host: &host, Port: &port}
}

func TestResolvePrivateEndpoint(t *testing.T) {
	private := testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "private.example.com", 4000)
	second := testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "second.example.com", 4001)
	instance := &api.Nextgenv1beta2Tidb{}
	_, err := resolvePrivateEndpoint(instance.Endpoints, "", nil)
	require.ErrorContains(t, err, "instance has no PRIVATE_ENDPOINT")

	instance.Endpoints = append(instance.Endpoints, private)
	selected, err := resolvePrivateEndpoint(instance.Endpoints, "", func([]string) (int, error) {
		t.Fatal("one usable address must not prompt")
		return 0, nil
	})
	require.NoError(t, err)
	require.Equal(t, private, *selected)

	instance.Endpoints = append(instance.Endpoints, second)
	_, err = resolvePrivateEndpoint(instance.Endpoints, "", nil)
	require.ErrorContains(t, err, "specify --endpoint")
	require.ErrorContains(t, err, "private.example.com:4000, second.example.com:4001")

	selected, err = resolvePrivateEndpoint(instance.Endpoints, "", func(addresses []string) (int, error) {
		require.Equal(t, []string{"private.example.com:4000", "second.example.com:4001"}, addresses)
		return 1, nil
	})
	require.NoError(t, err)
	require.Equal(t, second, *selected)

	selected, err = resolvePrivateEndpoint(instance.Endpoints, "second.example.com:4001", func([]string) (int, error) {
		t.Fatal("an explicit address must not prompt")
		return 0, nil
	})
	require.NoError(t, err)
	require.Equal(t, second, *selected)

	for _, address := range []string{"public.example.com:4000", "peering.example.com:4000", "other.example.com:4000", "second.example.com:4000"} {
		selected, err = resolvePrivateEndpoint(instance.Endpoints, address, nil)
		require.NoError(t, err)
		require.Equal(t, address, net.JoinHostPort(selected.GetHost(), strconv.Itoa(int(selected.GetPort()))))
	}
	_, err = resolvePrivateEndpoint(instance.Endpoints, "", func([]string) (int, error) { return 0, util.InterruptError })
	require.ErrorIs(t, err, util.InterruptError)
}

func TestManualPrivateEndpointIgnoresDiscoveryReachability(t *testing.T) {
	blocked := testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "blocked.example.com", 4000)
	blocked.ConnectionReachability = &api.V1beta2ConnectionReachability{
		Reachable: api.PtrBool(false),
		Detail:    &api.ConnectionReachabilityDetail{Message: api.PtrString("endpoint inactive")},
	}
	instance := &api.Nextgenv1beta2Tidb{Endpoints: []api.TidbEndpoint{
		blocked,
		testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "other.example.com", 4000),
		testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PUBLIC, "public.example.com", 4000),
	}}
	selected, err := resolvePrivateEndpoint(instance.Endpoints, "blocked.example.com:4000", nil)
	require.NoError(t, err)
	require.Equal(t, "blocked.example.com", selected.GetHost())
	require.Nil(t, selected.ConnectionReachability)
}

func TestShellConnectionTypeRoutesBeforeCARequest(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, tt := range []struct {
			name      string
			args      []string
			endpoints []api.TidbEndpoint
			wantCA    bool
		}{
			{"public default", nil, []api.TidbEndpoint{testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PUBLIC, "public.example.com", 4000)}, true},
			{"public never falls back", nil, []api.TidbEndpoint{testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "private.example.com", 4000)}, false},
			{"private", []string{"--connection-type", "private-endpoint"}, []api.TidbEndpoint{testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "private.example.com", 4000)}, true},
			{"private never falls back", []string{"--connection-type", "private-endpoint"}, []api.TidbEndpoint{testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PUBLIC, "public.example.com", 4000)}, false},
			{"private is not peering", []string{"--connection-type", "private-endpoint"}, []api.TidbEndpoint{testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_VPC_PEERING, "private.example.com", 4000)}, false},
			{"multiple private addresses", []string{"--connection-type", "private-endpoint"}, []api.TidbEndpoint{
				testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "one.example.com", 4000),
				testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "two.example.com", 4000),
			}, false},
			{"explicit private address", []string{"--connection-type", "private-endpoint", "--endpoint", "two.example.com:4000"}, []api.TidbEndpoint{
				testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "one.example.com", 4000),
				testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "two.example.com", 4000),
			}, true},
			{"manual address not in discovery", []string{"--connection-type", "private-endpoint", "--endpoint", "other.example.com:4000"}, []api.TidbEndpoint{
				testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "private.example.com", 4000),
			}, true},
		} {
			t.Run(plan.commandName+"/"+tt.name, func(t *testing.T) {
				instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", plan.servicePlan)
				instance.State = api.V1BETA1CLUSTERSTATE_ACTIVE.Ptr()
				instance.Endpoints = tt.endpoints
				if len(tt.args) > 0 {
					// Old GetTidb responses may report false by default; private discovery
					// must use the dedicated API even when the detail has a stale address.
					instance.Endpoints = []api.TidbEndpoint{{Host: api.PtrString("stale.example.com"), Port: api.PtrInt32(4000),
						ConnectionType:         api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT.Ptr(),
						ConnectionReachability: &api.V1beta2ConnectionReachability{Reachable: api.PtrBool(false)}}}
				}
				caStage := errors.New("CA request reached")
				caCalled := false
				client := &fakeNextGenClient{
					listPrivate: func(_ context.Context, id string, _ int32, _ string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
						require.NotEmpty(t, tt.args, "public shell must not query private addresses")
						require.NotContains(t, tt.args, "--endpoint", "manual endpoint must skip discovery")
						require.Equal(t, "tidb-1", id)
						result := &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{}
						for _, endpoint := range tt.endpoints {
							if endpoint.GetConnectionType() == api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT {
								result.PrivateEndpointConnections = append(result.PrivateEndpointConnections, activePrivateConnection(endpoint.GetHost(), strconv.Itoa(int(endpoint.GetPort()))))
							}
						}
						return result, nil
					},
					get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
					certificate: func(context.Context, string) (*api.V1beta2CaCertificateDownloadUrl, error) {
						caCalled = true
						return nil, caStage
					},
				}
				h := helperWithNextGenClient(client)
				h.IOStreams.CanPrompt = true
				command := shellCmd(h, plan)
				command.SetArgs(append([]string{"-c", "tidb-1"}, tt.args...))
				err := command.Execute()
				require.Error(t, err)
				require.Equal(t, tt.wantCA, caCalled)
				if tt.wantCA {
					require.ErrorIs(t, err, caStage)
				}
			})
		}
	}
}

func TestAnnotatePrivateShellTLSError(t *testing.T) {
	unknownAuthority := x509.UnknownAuthorityError{Cert: &x509.Certificate{}}
	err := annotateShellTLSError(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "private.example.com", &tls.CertificateVerificationError{Err: unknownAuthority})
	require.ErrorContains(t, err, `does not verify PRIVATE_ENDPOINT endpoint "private.example.com"`)
	var actual x509.UnknownAuthorityError
	require.ErrorAs(t, err, &actual)
}
