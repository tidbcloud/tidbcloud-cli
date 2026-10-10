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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type fakeNextGenClient struct {
	listPrivate    func(context.Context, string, int32, string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error)
	privateService func(context.Context, string) (*api.Nextgenv1beta2PrivateLinkService, error)
	create         func(context.Context, *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error)
	list           func(context.Context, api.TidbServiceListTidbsServicePlanParameter, *int32, *string) (*api.V1beta2ListTidbsResponse, error)
	get            func(context.Context, string) (*api.Nextgenv1beta2Tidb, error)
	update         func(context.Context, string, *api.TheTiDBCloudPremiumInstanceToUpdate) (*api.Nextgenv1beta2Tidb, error)
	delete         func(context.Context, string) (*api.Nextgenv1beta2Tidb, error)
	password       func(context.Context, string, *api.TidbServiceResetRootPasswordBody) error
	regions        func(context.Context, api.TidbServiceListTidbsServicePlanParameter, *int32, *string) (*api.V1beta2ListRegionsResponse, error)
	certificate    func(context.Context, string) (*api.V1beta2CaCertificateDownloadUrl, error)
	getPublic      func(context.Context, string) (*api.V1beta2PublicConnectionSetting, error)
	public         func(context.Context, string, *api.V1beta2PublicConnectionSetting) (*api.V1beta2PublicConnectionSetting, error)
	principal      func(context.Context, string) (*api.V1beta2CmekAccessIamPrincipal, error)
	verifyCMEK     func(context.Context, *api.V1beta2CustomerManagedEncryptionKey) (*api.V1beta2VerifyCmekAccessIamPrincipalResponse, error)
}

func (f *fakeNextGenClient) CreateTiDB(ctx context.Context, body *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error) {
	return f.create(ctx, body)
}

func (f *fakeNextGenClient) ListTiDBs(ctx context.Context, plan api.TidbServiceListTidbsServicePlanParameter, size *int32, token *string) (*api.V1beta2ListTidbsResponse, error) {
	return f.list(ctx, plan, size, token)
}

func (f *fakeNextGenClient) GetTiDB(ctx context.Context, id string) (*api.Nextgenv1beta2Tidb, error) {
	return f.get(ctx, id)
}

func (f *fakeNextGenClient) UpdateTiDB(ctx context.Context, id string, body *api.TheTiDBCloudPremiumInstanceToUpdate) (*api.Nextgenv1beta2Tidb, error) {
	return f.update(ctx, id, body)
}

func (f *fakeNextGenClient) DeleteTiDB(ctx context.Context, id string) (*api.Nextgenv1beta2Tidb, error) {
	return f.delete(ctx, id)
}

func (f *fakeNextGenClient) ResetRootPassword(ctx context.Context, id string, body *api.TidbServiceResetRootPasswordBody) error {
	return f.password(ctx, id, body)
}

func (f *fakeNextGenClient) ListRegions(ctx context.Context, plan api.TidbServiceListTidbsServicePlanParameter, size *int32, token *string) (*api.V1beta2ListRegionsResponse, error) {
	return f.regions(ctx, plan, size, token)
}

func (f *fakeNextGenClient) GetCACertificateDownloadURL(ctx context.Context, id string) (*api.V1beta2CaCertificateDownloadUrl, error) {
	return f.certificate(ctx, id)
}

func (f *fakeNextGenClient) GetPublicConnectionSetting(ctx context.Context, id string) (*api.V1beta2PublicConnectionSetting, error) {
	return f.getPublic(ctx, id)
}

func (f *fakeNextGenClient) UpdatePublicConnectionSetting(ctx context.Context, id string, body *api.V1beta2PublicConnectionSetting) (*api.V1beta2PublicConnectionSetting, error) {
	return f.public(ctx, id, body)
}

func (f *fakeNextGenClient) GetCmekAccessIAMPrincipal(ctx context.Context, regionID string) (*api.V1beta2CmekAccessIamPrincipal, error) {
	return f.principal(ctx, regionID)
}

func (f *fakeNextGenClient) VerifyCmekAccessIAMPrincipal(ctx context.Context, key *api.V1beta2CustomerManagedEncryptionKey) (*api.V1beta2VerifyCmekAccessIamPrincipalResponse, error) {
	return f.verifyCMEK(ctx, key)
}

func helperWithNextGenClient(client cloud.NextGenClient) *internal.Helper {
	return &internal.Helper{
		NextGenClient: func() (cloud.NextGenClient, error) { return client, nil },
		QueryPageSize: 10,
		IOStreams:     iostream.Test(),
	}
}

func TestEssentialV2IsCanonicalCommand(t *testing.T) {
	command := EssentialCmd(&internal.Helper{})
	require.Equal(t, "essential-v2", command.Name())
	require.Contains(t, command.Aliases, "essential")
}

func TestCreateBindsPlanEncryptionAndProject(t *testing.T) {
	tests := []struct {
		name string
		plan planSpec
	}{
		{name: "premium", plan: premiumPlan},
		{name: "essential-v2", plan: essentialV2Plan},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectID := "12345"
			var captured *api.Nextgenv1beta2Tidb
			client := &fakeNextGenClient{create: func(_ context.Context, body *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error) {
				captured = body
				id := "tidb-1"
				state := api.V1BETA1CLUSTERSTATE_CREATING
				body.TidbId = &id
				body.State = &state
				return body, nil
			}}
			command := createCmd(helperWithNextGenClient(client), tt.plan)
			command.SetArgs([]string{"--project-id", projectID, "--display-name", "test-instance", "--region", "aws-us-west-2", "--max-rcu", "10000", "--encryption", "default-key", "--output", "json"})
			require.NoError(t, command.ExecuteContext(context.Background()))
			require.Equal(t, tt.plan.servicePlan, captured.ServicePlan)
			require.Equal(t, "10000", captured.GetMaxRcu())
			require.NotNil(t, captured.Labels)
			require.Equal(t, projectID, (*captured.Labels)[projectIDLabel])
			require.Equal(t, projectID, command.Annotations[telemetry.ProjectID])
			require.NotNil(t, captured.DualLayerDataEncryption)
			require.NotNil(t, captured.DualLayerDataEncryption.DefaultKey)
			require.Nil(t, captured.DualLayerDataEncryption.Cmek)
			payload, err := json.Marshal(captured)
			require.NoError(t, err)
			require.Contains(t, string(payload), `"defaultKey":{}`)
			require.Contains(t, string(payload), `"maxRcu":"10000"`)
			require.NotContains(t, string(payload), `"baselineRcu"`)
			require.NotContains(t, string(payload), `"cmek"`)
		})
	}
}

func TestCreateOmitsProjectByDefault(t *testing.T) {
	var captured *api.Nextgenv1beta2Tidb
	client := &fakeNextGenClient{create: func(_ context.Context, body *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error) {
		captured = body
		id := "tidb-1"
		body.TidbId = &id
		return body, nil
	}}
	command := createCmd(helperWithNextGenClient(client), premiumPlan)
	command.SetArgs([]string{"--display-name", "test-instance", "--region", "aws-us-west-2", "--max-rcu", "10000"})

	require.NoError(t, command.ExecuteContext(context.Background()))
	require.Nil(t, captured.Labels)
	_, found := command.Annotations[telemetry.ProjectID]
	require.False(t, found)
}

func TestCreateRejectsInvalidOutputBeforeRequest(t *testing.T) {
	calls := 0
	client := &fakeNextGenClient{create: func(_ context.Context, _ *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error) {
		calls++
		return nil, nil
	}}
	command := createCmd(helperWithNextGenClient(client), premiumPlan)
	command.SetArgs([]string{"--display-name", "test-instance", "--region", "aws-us-west-2", "--max-rcu", "10000", "--output", "xml"})

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "unsupported output format: xml")
	require.Zero(t, calls)
}

func TestUpdateSendsOnlyChangedFields(t *testing.T) {
	id := "tidb-1"
	current := api.NewNextgenv1beta2Tidb("old-name", "aws-us-west-2", api.V1BETA1SERVICEPLAN_PREMIUM)
	current.TidbId = &id
	current.SetMaxRcu("10000")
	var captured *api.TheTiDBCloudPremiumInstanceToUpdate
	client := &fakeNextGenClient{
		get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return current, nil },
		update: func(_ context.Context, _ string, body *api.TheTiDBCloudPremiumInstanceToUpdate) (*api.Nextgenv1beta2Tidb, error) {
			captured = body
			result := *current
			result.DisplayName = body.GetDisplayName()
			return &result, nil
		},
	}
	command := updateCmd(helperWithNextGenClient(client), premiumPlan)
	command.SetArgs([]string{"--cluster-id", id, "--display-name", "new-name"})
	require.NoError(t, command.ExecuteContext(context.Background()))
	require.NotNil(t, captured.DisplayName)
	require.Equal(t, "new-name", *captured.DisplayName)
	require.Nil(t, captured.MaxRcu)
}

func TestUpdateRejectsInvalidInputBeforeRequest(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "output",
			args: []string{"--cluster-id", "tidb-1", "--display-name", "new-name", "--output", "xml"},
			want: "unsupported output format: xml",
		},
		{
			name: "display name",
			args: []string{"--cluster-id", "tidb-1", "--display-name", "bad"},
			want: "display name must be 4-64 characters",
		},
		{
			name: "max rcu",
			args: []string{"--cluster-id", "tidb-1", "--max-rcu", "0"},
			want: "--max-rcu must be greater than zero",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			client := &fakeNextGenClient{get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
				calls++
				return nil, nil
			}}
			command := updateCmd(helperWithNextGenClient(client), premiumPlan)
			command.SetArgs(tt.args)

			err := command.ExecuteContext(context.Background())
			require.ErrorContains(t, err, tt.want)
			require.Zero(t, calls)
		})
	}
}

func TestPlanGuardRejectsCrossPlanOperation(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		otherPlan := premiumPlan
		if plan.servicePlan == premiumPlan.servicePlan {
			otherPlan = essentialV2Plan
		}
		for _, args := range [][]string{
			{"describe", "-c", "tidb-1"},
			{"update", "-c", "tidb-1", "--display-name", "new-name"},
			{"delete", "-c", "tidb-1", "--force"},
			{"shell", "-c", "tidb-1", "--password", "test-password"},
			{"shell", "-c", "tidb-1", "--connection-type", "private-endpoint"},
			{"password", "-c", "tidb-1", "--password", "test-password"},
			{"public-endpoint", "enable", "-c", "tidb-1"},
			{"public-endpoint", "disable", "-c", "tidb-1", "--force"},
		} {
			t.Run(plan.commandName+"/"+args[0], func(t *testing.T) {
				instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", otherPlan.servicePlan)
				instance.TidbId = api.PtrString("tidb-1")
				client := &fakeNextGenClient{get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil }}
				h := helperWithNextGenClient(client)
				h.IOStreams.CanPrompt = true
				command := cmd(h, plan)
				command.SetArgs(args)
				// All other client methods are unset: a mutation or CA request would panic.
				require.ErrorContains(t, command.Execute(), "only manages \""+string(plan.servicePlan)+"\" instances")
			})
		}
	}
}

func TestDeleteRequiresConfirmationOrForce(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, force := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/force=%t", plan.commandName, force), func(t *testing.T) {
				deletes := 0
				client := &fakeNextGenClient{
					get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
						return api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", plan.servicePlan), nil
					},
					delete: func(_ context.Context, id string) (*api.Nextgenv1beta2Tidb, error) {
						require.Equal(t, "tidb-1", id)
						deletes++
						return nil, nil
					},
				}
				command := deleteCmd(helperWithNextGenClient(client), plan)
				args := []string{"-c", "tidb-1"}
				if force {
					args = append(args, "--force")
				}
				command.SetArgs(args)
				err := command.Execute()
				if force {
					require.NoError(t, err)
					require.Equal(t, 1, deletes)
				} else {
					require.ErrorContains(t, err, "terminal doesn't support prompts")
					require.Zero(t, deletes)
				}
			})
		}
	}
}

func TestRetrieveTiDBsUsesPlanAndPaginates(t *testing.T) {
	calls := 0
	client := &fakeNextGenClient{list: func(_ context.Context, plan api.TidbServiceListTidbsServicePlanParameter, size *int32, token *string) (*api.V1beta2ListTidbsResponse, error) {
		require.Equal(t, essentialV2Plan.queryPlan, plan)
		require.Equal(t, int32(10), *size)
		calls++
		instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", api.V1BETA1SERVICEPLAN_ESSENTIAL_V2)
		if calls == 1 {
			next := "page-2"
			return &api.V1beta2ListTidbsResponse{Tidbs: []api.Nextgenv1beta2Tidb{*instance}, NextPageToken: &next}, nil
		}
		require.NotNil(t, token)
		require.Equal(t, "page-2", *token)
		return &api.V1beta2ListTidbsResponse{Tidbs: []api.Nextgenv1beta2Tidb{*instance}}, nil
	}}
	instances, err := retrieveTiDBs(context.Background(), client, essentialV2Plan, 10)
	require.NoError(t, err)
	require.Len(t, instances, 2)
	require.Equal(t, 2, calls)
}

func TestRetrieveTiDBsRejectsRepeatedPageToken(t *testing.T) {
	calls := 0
	client := &fakeNextGenClient{list: func(_ context.Context, _ api.TidbServiceListTidbsServicePlanParameter, _ *int32, _ *string) (*api.V1beta2ListTidbsResponse, error) {
		calls++
		next := "repeated-token"
		return &api.V1beta2ListTidbsResponse{NextPageToken: &next}, nil
	}}

	_, err := retrieveTiDBs(context.Background(), client, premiumPlan, 10)
	require.ErrorContains(t, err, "repeated next page token")
	require.Equal(t, 2, calls)
}

func TestListRejectsInvalidOutputBeforeRequest(t *testing.T) {
	calls := 0
	client := &fakeNextGenClient{list: func(context.Context, api.TidbServiceListTidbsServicePlanParameter, *int32, *string) (*api.V1beta2ListTidbsResponse, error) {
		calls++
		return nil, nil
	}}
	command := listCmd(helperWithNextGenClient(client), premiumPlan)
	command.SetArgs([]string{"--output", "xml"})

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "unsupported output format: xml")
	require.Zero(t, calls)
}

func TestRegionRejectsRepeatedPageToken(t *testing.T) {
	calls := 0
	client := &fakeNextGenClient{regions: func(_ context.Context, _ api.TidbServiceListTidbsServicePlanParameter, _ *int32, _ *string) (*api.V1beta2ListRegionsResponse, error) {
		calls++
		next := "repeated-token"
		return &api.V1beta2ListRegionsResponse{NextPageToken: &next}, nil
	}}
	command := regionCmd(helperWithNextGenClient(client), premiumPlan)

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "repeated next page token")
	require.Equal(t, 2, calls)
}

func TestRegionRejectsInvalidOutputBeforeRequest(t *testing.T) {
	calls := 0
	client := &fakeNextGenClient{regions: func(context.Context, api.TidbServiceListTidbsServicePlanParameter, *int32, *string) (*api.V1beta2ListRegionsResponse, error) {
		calls++
		return nil, nil
	}}
	command := regionCmd(helperWithNextGenClient(client), premiumPlan)
	command.SetArgs([]string{"--output", "xml"})

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "unsupported output format: xml")
	require.Zero(t, calls)
}

func TestShellRequiresClusterIDForNonInteractiveMode(t *testing.T) {
	for _, args := range [][]string{{"--user", "app-user"}, {"--password", "secret"}} {
		h := helperWithNextGenClient(&fakeNextGenClient{})
		h.IOStreams.CanPrompt = true
		command := shellCmd(h, premiumPlan)
		command.SetArgs(args)

		err := command.ExecuteContext(context.Background())
		require.EqualError(t, err, `required flag(s) "cluster-id" not set`)
	}
}

func TestShellInteractiveModeSelectsInstance(t *testing.T) {
	id := "tidb-1"
	instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", api.V1BETA1SERVICEPLAN_PREMIUM)
	instance.TidbId = &id
	state := api.V1BETA1CLUSTERSTATE_CREATING
	instance.State = &state
	client := &fakeNextGenClient{get: func(_ context.Context, actualID string) (*api.Nextgenv1beta2Tidb, error) {
		require.Equal(t, id, actualID)
		return instance, nil
	}}
	h := helperWithNextGenClient(client)
	h.IOStreams.CanPrompt = true
	selectorCalls := 0
	command := shellCmdWithSelector(h, premiumPlan, func(_ context.Context, actualClient cloud.NextGenClient, actualPlan planSpec, pageSize int32) (*api.Nextgenv1beta2Tidb, error) {
		selectorCalls++
		require.Same(t, client, actualClient)
		require.Equal(t, premiumPlan, actualPlan)
		require.Equal(t, int32(h.QueryPageSize), pageSize)
		return instance, nil
	})

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, `instance "tidb-1" is in state "CREATING"`)
	require.Equal(t, 1, selectorCalls)
}

func TestShellClusterIDSkipsInteractiveSelection(t *testing.T) {
	id := "tidb-1"
	instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", api.V1BETA1SERVICEPLAN_PREMIUM)
	instance.TidbId = &id
	state := api.V1BETA1CLUSTERSTATE_CREATING
	instance.State = &state
	client := &fakeNextGenClient{get: func(_ context.Context, actualID string) (*api.Nextgenv1beta2Tidb, error) {
		require.Equal(t, id, actualID)
		return instance, nil
	}}
	h := helperWithNextGenClient(client)
	h.IOStreams.CanPrompt = true
	selectorCalls := 0
	command := shellCmdWithSelector(h, premiumPlan, func(context.Context, cloud.NextGenClient, planSpec, int32) (*api.Nextgenv1beta2Tidb, error) {
		selectorCalls++
		return nil, nil
	})
	command.SetArgs([]string{"--cluster-id", id})

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, `instance "tidb-1" is in state "CREATING"`)
	require.Zero(t, selectorCalls)
}

func TestSelectShellInstanceRejectsEmptyList(t *testing.T) {
	client := &fakeNextGenClient{list: func(_ context.Context, plan api.TidbServiceListTidbsServicePlanParameter, _ *int32, _ *string) (*api.V1beta2ListTidbsResponse, error) {
		require.Equal(t, essentialV2Plan.queryPlan, plan)
		return &api.V1beta2ListTidbsResponse{}, nil
	}}

	_, err := selectInstance(context.Background(), client, essentialV2Plan, 10)
	require.EqualError(t, err, "no available TiDB Cloud Essential V2 instances found")
}

func TestShellInstanceChoiceIncludesProject(t *testing.T) {
	id := "tidb-1"
	instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", api.V1BETA1SERVICEPLAN_PREMIUM)
	instance.TidbId = &id
	state := api.V1BETA1CLUSTERSTATE_ACTIVE
	instance.State = &state
	labels := map[string]string{projectIDLabel: "project-1"}
	instance.Labels = &labels

	require.Equal(t, "test-instance(tidb-1)[ACTIVE][aws-us-west-2][project:project-1]", instanceChoice{instance: instance}.String())
}

func TestResolvePublicEndpointSelectsPublicType(t *testing.T) {
	privateType := api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT
	publicType := api.ENDPOINTCONNECTIONTYPE_PUBLIC
	privateHost, publicHost := "private.example.com", "public.example.com"
	privatePort, publicPort := int32(4000), int32(4001)
	instance := &api.Nextgenv1beta2Tidb{Endpoints: []api.TidbEndpoint{
		{Host: &privateHost, Port: &privatePort, ConnectionType: &privateType},
		{Host: &publicHost, Port: &publicPort, ConnectionType: &publicType},
	}}
	endpoint, warning, err := resolvePublicEndpoint(context.Background(), nil, "tidb-1", instance)
	require.NoError(t, err)
	require.False(t, warning)
	require.Equal(t, publicHost, endpoint.GetHost())
	require.Equal(t, publicPort, endpoint.GetPort())
}

func TestResolveEndpointReportsReachabilityMessage(t *testing.T) {
	publicType := api.ENDPOINTCONNECTIONTYPE_PUBLIC
	host := "public.example.com"
	port := int32(4000)
	reachable := false
	message := "allowlist required"
	detail := api.ConnectionReachabilityDetail{Message: &message}
	instance := &api.Nextgenv1beta2Tidb{Endpoints: []api.TidbEndpoint{{
		Host:           &host,
		Port:           &port,
		ConnectionType: &publicType,
		ConnectionReachability: &api.V1beta2ConnectionReachability{
			Reachable: &reachable,
			Detail:    &detail,
		},
	}}}

	_, _, err := resolvePublicEndpoint(context.Background(), nil, "tidb-1", instance)
	require.EqualError(t, err, "PUBLIC endpoint is not reachable: allowlist required")
}

func TestDownloadCACertificateRequiresHTTPS(t *testing.T) {
	_, err := downloadCACertificate(context.Background(), "http://example.com/ca.pem")
	require.ErrorContains(t, err, "must use HTTPS")
}

func TestAnnotateShellTLSError(t *testing.T) {
	unknownAuthority := x509.UnknownAuthorityError{Cert: &x509.Certificate{}}
	err := annotateShellTLSError(api.ENDPOINTCONNECTIONTYPE_PUBLIC, "sql.example.com", &tls.CertificateVerificationError{Err: unknownAuthority})
	require.ErrorContains(t, err, `the CA certificate returned by the API does not verify PUBLIC endpoint "sql.example.com"`)
	var actual x509.UnknownAuthorityError
	require.ErrorAs(t, err, &actual)

	other := errors.New("connection reset")
	require.Same(t, other, annotateShellTLSError(api.ENDPOINTCONNECTIONTYPE_PUBLIC, "sql.example.com", other))
}

func actionTestRoot() *cobra.Command {
	root := &cobra.Command{Use: "ticloud"}
	root.PersistentFlags().BoolP("debug", "D", false, "")
	root.PersistentFlags().Bool("no-color", false, "")
	root.PersistentFlags().StringP("profile", "P", "", "")
	root.AddCommand(PremiumCmd(&internal.Helper{}), EssentialCmd(&internal.Helper{}))
	return root
}

func TestNormalizeActionArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{name: "premium create", args: []string{"premium", "--create", "--display-name", "test-instance"}, want: []string{"premium", "create", "--display-name", "test-instance"}},
		{name: "essential list", args: []string{"essential", "--list"}, want: []string{"essential", "list"}},
		{name: "essential v2 describe", args: []string{"essential-v2", "--describe", "-c", "tidb-1"}, want: []string{"essential-v2", "describe", "-c", "tidb-1"}},
		{name: "explicit true action", args: []string{"premium", "--list=true"}, want: []string{"premium", "list"}},
		{name: "plan name used as profile", args: []string{"--profile", "premium", "premium", "--list"}, want: []string{"--profile", "premium", "premium", "list"}},
		{name: "profile before action", args: []string{"premium", "--profile", "prod", "--list"}, want: []string{"premium", "list", "--profile", "prod"}},
		{name: "short profile before action", args: []string{"essential", "-P", "prod", "--describe", "-c", "tidb-1"}, want: []string{"essential", "describe", "-P", "prod", "-c", "tidb-1"}},
		{name: "business flags before action", args: []string{"premium", "--output", "json", "--list"}, want: []string{"premium", "list", "--output", "json"}},
		{name: "project before create action", args: []string{"premium", "--project-id", "12345", "--create", "--display-name", "test-instance"}, want: []string{"premium", "create", "--project-id", "12345", "--display-name", "test-instance"}},
		{name: "debug before action", args: []string{"premium", "--debug", "--list"}, want: []string{"premium", "list", "--debug"}},
		{name: "short debug before command", args: []string{"-D=true", "premium", "--list"}, want: []string{"-D=true", "premium", "list"}},
		{name: "grouped root flags", args: []string{"-DPprod", "premium", "--list"}, want: []string{"-DPprod", "premium", "list"}},
		{name: "grouped inherited flags", args: []string{"essential-v2", "-DPprod", "--list"}, want: []string{"essential-v2", "list", "-DPprod"}},
		{name: "action name as profile value", args: []string{"premium", "-P", "--delete", "--list"}, want: []string{"premium", "list", "-P", "--delete"}},
		{name: "action name as display name value", args: []string{"premium", "--display-name", "--delete", "--create"}, want: []string{"premium", "create", "--display-name", "--delete"}},
		{name: "subcommand name as display name value", args: []string{"premium", "-n", "shell", "--create"}, want: []string{"premium", "create", "-n", "shell"}},
		{name: "help before action", args: []string{"premium", "-h", "--create"}, want: []string{"premium", "create", "-h"}},
		{name: "completion", args: []string{"__complete", "premium", "--create", "--"}, want: []string{"__complete", "premium", "create", "--"}},
		{name: "completion after grouped flags", args: []string{"__complete", "-DPprod", "essential", "--create", "--reg"}, want: []string{"__complete", "-DPprod", "essential", "create", "--reg"}},
		{name: "completion without descriptions", args: []string{"__completeNoDesc", "essential", "--list", "--o"}, want: []string{"__completeNoDesc", "essential", "list", "--o"}},
		{name: "missing value left for completion", args: []string{"__complete", "premium", "--create", "--region"}, want: []string{"__complete", "premium", "create", "--region"}},
		{name: "subcommand unchanged", args: []string{"premium", "shell", "-c", "tidb-1"}, want: []string{"premium", "shell", "-c", "tidb-1"}},
		{name: "profile before subcommand unchanged", args: []string{"premium", "--profile", "prod", "shell", "-c", "tidb-1"}, want: []string{"premium", "--profile", "prod", "shell", "-c", "tidb-1"}},
		{name: "short debug before subcommand unchanged", args: []string{"premium", "-D=true", "shell", "-c", "tidb-1", "--password", "--create"}, want: []string{"premium", "-D=true", "shell", "-c", "tidb-1", "--password", "--create"}},
		{name: "attached profile before subcommand unchanged", args: []string{"premium", "-Pprod", "shell", "-c", "tidb-1", "--password", "--create"}, want: []string{"premium", "-Pprod", "shell", "-c", "tidb-1", "--password", "--create"}},
		{name: "password value unchanged", args: []string{"premium", "shell", "-c", "tidb-1", "--password", "--create"}, want: []string{"premium", "shell", "-c", "tidb-1", "--password", "--create"}},
		{name: "password before subcommand unchanged", args: []string{"premium", "--password", "--delete", "shell", "-c", "tidb-1"}, want: []string{"premium", "--password", "--delete", "shell", "-c", "tidb-1"}},
		{name: "terminator used as password value", args: []string{"essential", "--password", "--", "shell", "-c", "tidb-1"}, want: []string{"essential", "--password", "--", "shell", "-c", "tidb-1"}},
		{name: "action after flag terminator unchanged", args: []string{"premium", "-c", "tidb-1", "--force", "--", "--delete"}, want: []string{"premium", "-c", "tidb-1", "--force", "--", "--delete"}},
		{name: "root flag terminator unchanged", args: []string{"--", "premium", "--list"}, want: []string{"--", "premium", "--list"}},
		{name: "unknown flag retained for Cobra", args: []string{"premium", "--list", "--unknown"}, want: []string{"premium", "list", "--unknown"}},
		{name: "unknown flag value not treated as action", args: []string{"premium", "--unknown", "--delete"}, want: []string{"premium", "--unknown", "--delete"}},
		{name: "empty args"},
		{name: "unrelated command values unchanged", args: []string{"serverless", "shell", "-c", "cluster-1", "--user", "premium", "--password", "--list"}, want: []string{"serverless", "shell", "-c", "cluster-1", "--user", "premium", "--password", "--list"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := actionTestRoot()
			original := append([]string(nil), tt.args...)
			actual, err := NormalizeActionArgs(root, tt.args)
			require.NoError(t, err)
			require.Equal(t, tt.want, actual)
			require.Equal(t, original, tt.args)
			require.Zero(t, root.PersistentFlags().NFlag())
			debug, err := root.PersistentFlags().GetBool("debug")
			require.NoError(t, err)
			require.False(t, debug)
			profile, err := root.PersistentFlags().GetString("profile")
			require.NoError(t, err)
			require.Empty(t, profile)
		})
	}
}

func TestNormalizeActionArgsRejectsInvalidSelections(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "false action", args: []string{"premium", "--list=false"}, want: "--list=false does not select an action"},
		{name: "invalid boolean", args: []string{"premium", "--list=invalid"}, want: `invalid value "invalid" for --list`},
		{name: "multiple actions", args: []string{"essential", "--list", "--delete"}, want: "only one action flag may be specified"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NormalizeActionArgs(actionTestRoot(), tt.args)
			require.ErrorContains(t, err, tt.want)
		})
	}
}

func TestActionFlagsAreRegistered(t *testing.T) {
	command := PremiumCmd(helperWithNextGenClient(&fakeNextGenClient{}))
	for _, name := range []string{"create", "list", "describe", "update", "delete"} {
		require.NotNil(t, command.Flags().Lookup(name))
		child, _, err := command.Find([]string{name})
		require.NoError(t, err)
		require.True(t, child.Hidden)
	}
}

func TestActionHelpUsesPublicFlagSyntax(t *testing.T) {
	command := PremiumCmd(helperWithNextGenClient(&fakeNextGenClient{}))
	create, _, err := command.Find([]string{"create"})
	require.NoError(t, err)

	var output bytes.Buffer
	create.SetOut(&output)
	require.NoError(t, create.Help())
	require.Contains(t, output.String(), "ticloud premium --create [flags]")
	require.NotContains(t, output.String(), "ticloud premium create [flags]")
	require.Contains(t, command.Long, "--display-name")
	require.Contains(t, command.Long, "--cluster-id")
}

func (f *fakeNextGenClient) ListPrivateEndpointConnections(ctx context.Context, id string, size int32, token string) (*api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse, error) {
	return f.listPrivate(ctx, id, size, token)
}

func (f *fakeNextGenClient) GetPrivateLinkService(ctx context.Context, id string) (*api.Nextgenv1beta2PrivateLinkService, error) {
	return f.privateService(ctx, id)
}
