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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"
)

func activePrivateConnection(host, port string) api.Nextgenv1beta2PrivateEndpointConnection {
	return api.Nextgenv1beta2PrivateEndpointConnection{
		Host: api.PtrString(host), Port: api.PtrString(port),
		EndpointState:           api.NEXTGENV1BETA2PRIVATEENDPOINTCONNECTIONENDPOINTSTATE_ACTIVE.Ptr(),
		PrivateLinkServiceState: api.NEXTGENV1BETA2PRIVATELINKSERVICESTATE_ACTIVE.Ptr(),
	}
}

func TestLoadPrivateEndpointsPaginationAndReadiness(t *testing.T) {
	pending := activePrivateConnection("pending.example.com", "4000")
	pending.SetEndpointState(api.NEXTGENV1BETA2PRIVATEENDPOINTCONNECTIONENDPOINTSTATE_PENDING)
	inactiveService := activePrivateConnection("inactive.example.com", "4000")
	inactiveService.PrivateLinkServiceState = nil
	calls := 0
	client := &fakeNextGenClient{listPrivate: func(_ context.Context, id string, size int32, token string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
		calls++
		require.Equal(t, "tidb-1", id)
		require.EqualValues(t, 2, size)
		if calls == 1 {
			require.Empty(t, token)
			return &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{
				PrivateEndpointConnections: []api.Nextgenv1beta2PrivateEndpointConnection{
					activePrivateConnection("one.example.com", "04000"), pending, inactiveService,
					activePrivateConnection("bad.example.com", "65536"), activePrivateConnection("", "4000"),
					activePrivateConnection("https://invalid.example.com", "4000"), activePrivateConnection("zero.example.com", "0"),
				}, NextPageToken: api.PtrString("page-2"),
			}, nil
		}
		require.Equal(t, "page-2", token)
		return &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{
			PrivateEndpointConnections: []api.Nextgenv1beta2PrivateEndpointConnection{
				activePrivateConnection("one.example.com", "4000"), activePrivateConnection("2001:db8::1", "4001"),
			},
		}, nil
	}}
	// AWS must also prefer the list, avoiding the service getter's multi-service restriction.
	endpoints, err := loadPrivateEndpoints(context.Background(), client, "tidb-1", api.REGIONCLOUDPROVIDER_AWS, 2)
	require.NoError(t, err)
	require.Equal(t, 2, calls)
	require.Len(t, endpoints, 2)
	require.Equal(t, "one.example.com", endpoints[0].GetHost())
	require.EqualValues(t, 4000, endpoints[0].GetPort())
	require.Nil(t, endpoints[0].ConnectionReachability, "service state does not prove network reachability")
	selected, err := resolvePrivateEndpoint(endpoints, "[2001:db8::1]:4001", nil)
	require.NoError(t, err)
	require.Equal(t, "2001:db8::1", selected.GetHost())
}

func TestLoadPrivateEndpointsAWSServiceFallback(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		provider                 api.RegionCloudProvider
		records                  bool
		mutate                   func(*api.Nextgenv1beta2PrivateLinkService)
		serviceError             error
		wantService, wantSuccess bool
	}{
		{name: "AWS empty list", provider: api.REGIONCLOUDPROVIDER_AWS, wantService: true, wantSuccess: true},
		{name: "pending connection blocks fallback", provider: api.REGIONCLOUDPROVIDER_AWS, records: true},
		{name: "GCP requires connection", provider: api.REGIONCLOUDPROVIDER_GCP},
		{name: "unknown provider requires connection"},
		{name: "inactive service", provider: api.REGIONCLOUDPROVIDER_AWS, wantService: true, mutate: func(s *api.Nextgenv1beta2PrivateLinkService) { s.State = nil }},
		{name: "non AWS service DNS is not a host", provider: api.REGIONCLOUDPROVIDER_AWS, wantService: true, mutate: func(s *api.Nextgenv1beta2PrivateLinkService) { s.CloudProvider = api.REGIONCLOUDPROVIDER_GCP.Ptr() }},
		{name: "invalid service address", provider: api.REGIONCLOUDPROVIDER_AWS, wantService: true, mutate: func(s *api.Nextgenv1beta2PrivateLinkService) { s.ServicePort = api.PtrInt32(0) }},
		{name: "missing DNS", provider: api.REGIONCLOUDPROVIDER_AWS, wantService: true, mutate: func(s *api.Nextgenv1beta2PrivateLinkService) { s.ServiceDnsName = nil }},
		{name: "service permission denied", provider: api.REGIONCLOUDPROVIDER_AWS, wantService: true, serviceError: errors.New("403 forbidden")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			serviceCalled := false
			client := &fakeNextGenClient{
				listPrivate: func(context.Context, string, int32, string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
					result := &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{}
					if tc.records {
						result.PrivateEndpointConnections = []api.Nextgenv1beta2PrivateEndpointConnection{{}}
					}
					return result, nil
				},
				privateService: func(_ context.Context, id string) (*api.Nextgenv1beta2PrivateLinkService, error) {
					serviceCalled = true
					require.Equal(t, "tidb-1", id)
					service := &api.Nextgenv1beta2PrivateLinkService{
						CloudProvider:  api.REGIONCLOUDPROVIDER_AWS.Ptr(),
						State:          api.NEXTGENV1BETA2PRIVATELINKSERVICESTATE_ACTIVE.Ptr(),
						ServiceDnsName: api.PtrString("service.example.com"), ServicePort: api.PtrInt32(4000),
					}
					if tc.mutate != nil {
						tc.mutate(service)
					}
					return service, tc.serviceError
				},
			}
			endpoints, err := loadPrivateEndpoints(context.Background(), client, "tidb-1", tc.provider, 10)
			require.Equal(t, tc.wantService, serviceCalled)
			if tc.wantSuccess {
				require.NoError(t, err)
				require.Len(t, endpoints, 1)
				require.Equal(t, "service.example.com", endpoints[0].GetHost())
			} else {
				require.Error(t, err)
				require.Nil(t, endpoints)
				if tc.serviceError != nil {
					require.ErrorIs(t, err, tc.serviceError)
				}
			}
		})
	}
}

func TestLoadPrivateEndpointsFailsClosed(t *testing.T) {
	denied := errors.New("403 forbidden")
	for _, tc := range []struct {
		name     string
		failPage int
		repeat   bool
		nilPage  bool
	}{
		{name: "first page denied", failPage: 1},
		{name: "later page denied", failPage: 2},
		{name: "repeated token", repeat: true},
		{name: "nil response", nilPage: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			client := &fakeNextGenClient{listPrivate: func(context.Context, string, int32, string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
				calls++
				if calls == tc.failPage {
					return nil, denied
				}
				if tc.nilPage {
					return nil, nil
				}
				return &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{
					PrivateEndpointConnections: []api.Nextgenv1beta2PrivateEndpointConnection{activePrivateConnection("one.example.com", "4000")},
					NextPageToken:              api.PtrString("same-token"),
				}, nil
			}}
			endpoints, err := loadPrivateEndpoints(context.Background(), client, "tidb-1", api.REGIONCLOUDPROVIDER_AWS, 10)
			require.Error(t, err)
			require.Nil(t, endpoints, "do not use a partial list or bypass an error via service lookup")
			if tc.failPage > 0 {
				require.ErrorIs(t, err, denied)
			}
			if tc.repeat {
				require.ErrorContains(t, err, "repeated page token")
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := loadPrivateEndpoints(ctx, &fakeNextGenClient{}, "tidb-1", api.REGIONCLOUDPROVIDER_AWS, 10)
	require.ErrorIs(t, err, context.Canceled)
}

func TestPrivateShellListDenialStopsBeforeCA(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		t.Run(plan.commandName, func(t *testing.T) {
			denied := errors.New("403 forbidden")
			instance := api.NewNextgenv1beta2Tidb("test", "aws-us-west-2", plan.servicePlan)
			instance.State = api.V1BETA1CLUSTERSTATE_ACTIVE.Ptr()
			instance.CloudProvider = api.REGIONCLOUDPROVIDER_AWS.Ptr()
			instance.Endpoints = []api.TidbEndpoint{testShellEndpoint(api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, "old.example.com", 4000)}
			client := &fakeNextGenClient{
				get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
				listPrivate: func(context.Context, string, int32, string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
					return nil, denied
				},
			}
			h := helperWithNextGenClient(client)
			h.IOStreams.CanPrompt = true
			cmd := shellCmd(h, plan)
			cmd.SetArgs([]string{"-c", "tidb-1", "--connection-type", "private-endpoint"})
			require.ErrorIs(t, cmd.Execute(), denied)
		})
	}
}

func TestPrivateEndpointDiscoverySharesAndReleasesDeadline(t *testing.T) {
	parent := context.Background()
	var lookupCtx context.Context
	var deadline time.Time
	calls := 0
	checkContext := func(ctx context.Context) {
		actual, ok := ctx.Deadline()
		require.True(t, ok)
		if lookupCtx == nil {
			lookupCtx, deadline = ctx, actual
			require.InDelta(t, 30, time.Until(deadline).Seconds(), 1)
		} else {
			require.Equal(t, deadline, actual, "pagination and service fallback must share one budget")
		}
	}
	client := &fakeNextGenClient{
		listPrivate: func(ctx context.Context, _ string, _ int32, token string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
			checkContext(ctx)
			calls++
			if calls == 1 {
				require.Empty(t, token)
				return &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{NextPageToken: api.PtrString("page-2")}, nil
			}
			require.Equal(t, "page-2", token)
			return &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{}, nil
		},
		privateService: func(ctx context.Context, _ string) (*api.Nextgenv1beta2PrivateLinkService, error) {
			checkContext(ctx)
			return &api.Nextgenv1beta2PrivateLinkService{
				CloudProvider:  api.REGIONCLOUDPROVIDER_AWS.Ptr(),
				State:          api.NEXTGENV1BETA2PRIVATELINKSERVICESTATE_ACTIVE.Ptr(),
				ServiceDnsName: api.PtrString("sql.example.com"), ServicePort: api.PtrInt32(4000),
			}, nil
		},
	}
	endpoints, err := loadPrivateEndpoints(parent, client, "tidb-1", api.REGIONCLOUDPROVIDER_AWS, 10)
	require.NoError(t, err)
	require.Len(t, endpoints, 1)
	require.Equal(t, 2, calls)
	require.ErrorIs(t, lookupCtx.Err(), context.Canceled, "release the discovery context before prompting or SQL")
	require.NoError(t, parent.Err(), "do not cancel the command or SQL session")
}

func TestPrivateEndpointDiscoveryCancelsInFlightRequests(t *testing.T) {
	for _, stage := range []string{"first page", "later page", "AWS service"} {
		t.Run(stage, func(t *testing.T) {
			parent, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			waitForDeadline := func(ctx context.Context) error {
				deadline, ok := ctx.Deadline()
				require.True(t, ok)
				parentDeadline, _ := parent.Deadline()
				require.Equal(t, parentDeadline, deadline, "respect an earlier caller deadline")
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(2 * time.Second):
					t.Fatal("request was not canceled")
					return nil
				}
			}
			calls := 0
			client := &fakeNextGenClient{
				listPrivate: func(ctx context.Context, _ string, _ int32, _ string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
					calls++
					if stage == "first page" || (stage == "later page" && calls == 2) {
						return nil, waitForDeadline(ctx)
					}
					result := &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{}
					if stage == "later page" {
						result.PrivateEndpointConnections = []api.Nextgenv1beta2PrivateEndpointConnection{activePrivateConnection("sql.example.com", "4000")}
						result.NextPageToken = api.PtrString("page-2")
					}
					return result, nil
				},
				privateService: func(ctx context.Context, _ string) (*api.Nextgenv1beta2PrivateLinkService, error) {
					require.Equal(t, "AWS service", stage)
					return nil, waitForDeadline(ctx)
				},
			}
			endpoints, err := loadPrivateEndpoints(parent, client, "tidb-1", api.REGIONCLOUDPROVIDER_AWS, 10)
			require.ErrorIs(t, err, context.DeadlineExceeded)
			require.Nil(t, endpoints, "do not connect using a partial result after timeout")
		})
	}
}

func TestManualPrivateEndpointSkipsDiscoveryButKeepsInstanceChecks(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, failure := range []string{"get denied", "wrong plan", "inactive", "CA reached"} {
			t.Run(plan.commandName+"/"+failure, func(t *testing.T) {
				denied := errors.New("403 forbidden")
				caStop := errors.New("CA reached")
				instance := api.NewNextgenv1beta2Tidb("test", "aws-us-west-2", plan.servicePlan)
				instance.State = api.V1BETA1CLUSTERSTATE_ACTIVE.Ptr()
				instance.CloudProvider = api.REGIONCLOUDPROVIDER_AWS.Ptr()
				if failure == "wrong plan" {
					instance.ServicePlan = api.V1BETA1SERVICEPLAN_STARTER
				}
				if failure == "inactive" {
					instance.State = api.V1BETA1CLUSTERSTATE_CREATING.Ptr()
				}
				caCalled := false
				client := &fakeNextGenClient{
					get: func(_ context.Context, id string) (*api.Nextgenv1beta2Tidb, error) {
						require.Equal(t, "tidb-1", id)
						if failure == "get denied" {
							return nil, denied
						}
						return instance, nil
					},
					listPrivate: func(context.Context, string, int32, string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
						t.Fatal("explicit endpoint must skip listing, including when no connections are registered")
						return nil, nil
					},
					privateService: func(context.Context, string) (*api.Nextgenv1beta2PrivateLinkService, error) {
						t.Fatal("explicit endpoint must not hit the multi-service restriction")
						return nil, nil
					},
					certificate: func(ctx context.Context, id string) (*api.V1beta2CaCertificateDownloadUrl, error) {
						caCalled = true
						require.Equal(t, "tidb-1", id)
						require.NoError(t, ctx.Err())
						_, hasDeadline := ctx.Deadline()
						require.False(t, hasDeadline, "discovery deadline must not leak into later stages")
						return nil, caStop
					},
				}
				h := helperWithNextGenClient(client)
				h.IOStreams.CanPrompt = true
				cmd := shellCmd(h, plan)
				cmd.SetArgs([]string{"-c", "tidb-1", "--connection-type", "private-endpoint", "--endpoint", "manual.example.com:4000"})
				err := cmd.Execute()
				require.Error(t, err)
				require.Equal(t, failure == "CA reached", caCalled)
				if failure == "CA reached" {
					require.ErrorIs(t, err, caStop)
				}
				if failure == "get denied" {
					require.ErrorIs(t, err, denied)
				}
			})
		}
	}
}

func TestManualPrivateEndpointParsesWithoutDiscovery(t *testing.T) {
	for _, tc := range []struct {
		name, address, host string
		port                int32
	}{
		{name: "hostname", address: "manual.example.com:04000", host: "manual.example.com", port: 4000},
		{name: "IPv6", address: "[2001:db8::1]:4001", host: "2001:db8::1", port: 4001},
	} {
		t.Run(tc.name, func(t *testing.T) {
			endpoint, err := resolvePrivateEndpoint(nil, tc.address, func([]string) (int, error) {
				t.Fatal("explicit endpoint must not prompt")
				return 0, nil
			})
			require.NoError(t, err)
			require.Equal(t, tc.host, endpoint.GetHost())
			require.Equal(t, tc.port, endpoint.GetPort())
			require.Equal(t, api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT, endpoint.GetConnectionType())
		})
	}
	_, err := resolvePrivateEndpoint(nil, "invalid:65536", nil)
	require.Error(t, err)
}
