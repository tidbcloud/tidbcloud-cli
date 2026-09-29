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
	"errors"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/stretchr/testify/require"
)

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
	public := testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PUBLIC, "public.example.com", 4000)
	peering := testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_VPC_PEERING, "peering.example.com", 4000)

	instance := &api.Nextgenv1beta2Tidb{Endpoints: []api.TidbEndpoint{public, peering}}
	_, err := resolvePrivateEndpoint(instance, "", nil)
	require.ErrorContains(t, err, "instance has no PRIVATE_ENDPOINT")

	instance.Endpoints = append(instance.Endpoints, testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "", 0))
	_, err = resolvePrivateEndpoint(instance, "", nil)
	require.ErrorContains(t, err, "not ready")

	instance.Endpoints = append(instance.Endpoints, private)
	selected, err := resolvePrivateEndpoint(instance, "", func([]string) (int, error) {
		t.Fatal("one usable address must not prompt")
		return 0, nil
	})
	require.NoError(t, err)
	require.Equal(t, private, *selected)

	instance.Endpoints = append(instance.Endpoints, second)
	_, err = resolvePrivateEndpoint(instance, "", nil)
	require.ErrorContains(t, err, "specify --endpoint")
	require.ErrorContains(t, err, "private.example.com:4000, second.example.com:4001")

	selected, err = resolvePrivateEndpoint(instance, "", func(addresses []string) (int, error) {
		require.Equal(t, []string{"private.example.com:4000", "second.example.com:4001"}, addresses)
		return 1, nil
	})
	require.NoError(t, err)
	require.Equal(t, second, *selected)

	selected, err = resolvePrivateEndpoint(instance, "second.example.com:4001", func([]string) (int, error) {
		t.Fatal("an explicit address must not prompt")
		return 0, nil
	})
	require.NoError(t, err)
	require.Equal(t, second, *selected)

	for _, address := range []string{"public.example.com:4000", "peering.example.com:4000", "other.example.com:4000", "second.example.com:4000"} {
		_, err = resolvePrivateEndpoint(instance, address, nil)
		require.ErrorContains(t, err, "is not a PRIVATE_ENDPOINT address returned for this instance")
	}
	_, err = resolvePrivateEndpoint(instance, "", func([]string) (int, error) { return 0, util.InterruptError })
	require.ErrorIs(t, err, util.InterruptError)
}

func TestPrivateEndpointRejectsExplicitlyUnreachableAddress(t *testing.T) {
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
	_, err := resolvePrivateEndpoint(instance, "blocked.example.com:4000", nil)
	require.EqualError(t, err, "PRIVATE_ENDPOINT endpoint is not reachable: endpoint inactive")
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
			{"foreign private address", []string{"--connection-type", "private-endpoint", "--endpoint", "other.example.com:4000"}, []api.TidbEndpoint{
				testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "private.example.com", 4000),
			}, false},
		} {
			t.Run(plan.commandName+"/"+tt.name, func(t *testing.T) {
				instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", "5000", plan.servicePlan)
				instance.State = api.V1BETA1CLUSTERSTATE_ACTIVE.Ptr()
				instance.Endpoints = tt.endpoints
				caStage := errors.New("CA request reached")
				caCalled := false
				client := &fakeNextGenClient{
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
