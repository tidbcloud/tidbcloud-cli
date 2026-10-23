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
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	configCmd "github.com/tidbcloud/tidbcloud-cli/internal/cli/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"

	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"
)

type projectLookupClient struct {
	cloud.TiDBCloudClient
	list func(context.Context, *int32, *string) (*iam.ApiListProjectsRsp, error)
}

func (c projectLookupClient) ListProjects(ctx context.Context, size *int32, token *string) (*iam.ApiListProjectsRsp, error) {
	return c.list(ctx, size, token)
}

func TestCreateProjectName(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, tc := range []struct {
			name            string
			firstName       string
			secondName      string
			listError       bool
			secondPageError bool
			wantError       string
		}{
			{name: "match on second page", firstName: "other", secondName: "my-project"},
			{name: "duplicate across pages", firstName: "my-project", secondName: "my-project", wantError: "multiple accessible projects"},
			{name: "case sensitive", firstName: "My-project", secondName: "other", wantError: "no accessible project"},
			{name: "failure after match", firstName: "my-project", secondPageError: true, wantError: "permission denied"},
			{name: "list failure", listError: true, wantError: "permission denied"},
		} {
			t.Run(plan.commandName+"/"+tc.name, func(t *testing.T) {
				created := false
				h := helperWithNextGenClient(&fakeNextGenClient{create: func(_ context.Context, body *api.Nextgenv1beta2Tidb) (*api.Nextgenv1beta2Tidb, error) {
					created = true
					require.Equal(t, "456", (*body.Labels)[projectIDLabel])
					require.Equal(t, plan.servicePlan, body.ServicePlan)
					return body, nil
				}})
				calls := 0
				h.Client = func() (cloud.TiDBCloudClient, error) {
					return projectLookupClient{list: func(_ context.Context, size *int32, token *string) (*iam.ApiListProjectsRsp, error) {
						calls++
						require.Equal(t, int32(h.QueryPageSize), *size)
						if tc.listError || (tc.secondPageError && calls == 2) {
							return nil, errors.New("permission denied")
						}
						id, name := "123", tc.firstName
						var next *string
						if calls == 1 {
							require.Nil(t, token)
							value := "page-2"
							next = &value
						} else {
							require.Equal(t, "page-2", *token)
							id, name = "456", tc.secondName
						}
						return &iam.ApiListProjectsRsp{Projects: []iam.ApiProject{{Id: &id, Name: &name}}, NextPageToken: next}, nil
					}}, nil
				}
				command := createCmd(h, plan)
				command.SetArgs([]string{"--project-name", "my-project", "--display-name", "test-instance", "--region", "aws-us-west-2", "--max-rcu", "10000", "--output", "json"})
				err := command.ExecuteContext(context.Background())
				if tc.wantError != "" {
					require.ErrorContains(t, err, tc.wantError)
					require.False(t, created)
				} else {
					require.NoError(t, err)
					require.True(t, created)
					require.Equal(t, "456", command.Annotations[telemetry.ProjectID])
				}
				if !tc.listError {
					require.Equal(t, 2, calls)
				}
			})
		}
	}
}

func TestCreateProjectNameRejectsInvalidFlagsBeforeClient(t *testing.T) {
	for _, args := range [][]string{
		{"--project-name", ""},
		{"--project-name", "my-project", "--project-id", "123"},
		{"--project-name", "my-project", "--project-id", ""},
	} {
		command := createCmd(&internal.Helper{}, premiumPlan)
		command.SetArgs(append([]string{"--display-name", "test-instance", "--region", "aws-us-west-2", "--max-rcu", "10000"}, args...))
		require.Error(t, command.ExecuteContext(context.Background()))
	}
}

// Exercise persisted profile endpoints with real SDK clients. The transport
// intercepts every request, so this test never contacts either environment.
type projectRoutingTransport func(*http.Request) (*http.Response, error)

func (f projectRoutingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestCreateProjectNameUsesStagingProfileEndpoints(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		t.Run(plan.commandName, func(t *testing.T) {
			previousProfile := config.ActiveProfileName()
			previousTransport := http.DefaultTransport
			viper.Reset()
			t.Cleanup(func() {
				http.DefaultTransport = previousTransport
				viper.Reset()
				config.SetActiveProfile(previousProfile)
			})
			configPath := filepath.Join(t.TempDir(), "config.toml")
			viper.SetConfigFile(configPath)
			require.NoError(t, viper.WriteConfig())
			setup := configCmd.CreateCmd(&internal.Helper{IOStreams: iostream.Test()})
			setup.SetArgs([]string{"--profile-name", "staging-test", "--environment", "staging"})
			require.NoError(t, setup.Execute())
			viper.Reset()
			viper.SetConfigFile(configPath)
			require.NoError(t, viper.ReadInConfig())
			config.SetActiveProfile("staging-test")

			var requests []string
			http.DefaultTransport = projectRoutingTransport(func(r *http.Request) (*http.Response, error) {
				requests = append(requests, r.Method+" "+r.URL.Host+r.URL.Path)
				require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
				var body string
				switch r.URL.Host {
				case "iam.staging.tidbapi.com":
					require.Equal(t, http.MethodGet, r.Method)
					require.Equal(t, "/v1beta1/projects", r.URL.Path)
					body = `{"projects":[{"id":"456","name":"my-project"}]}`
				case "cloud.staging.tidbapi.com":
					require.Equal(t, http.MethodPost, r.Method)
					require.Equal(t, "/v1beta2/tidbs", r.URL.Path)
					payload, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					var instance api.Nextgenv1beta2Tidb
					require.NoError(t, json.Unmarshal(payload, &instance))
					require.Equal(t, "456", instance.GetLabels()[projectIDLabel])
					require.Equal(t, plan.servicePlan, instance.ServicePlan)
					body = string(payload)
				default:
					t.Fatalf("unexpected endpoint: %s", r.URL.Host)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})
			h := &internal.Helper{
				IOStreams: iostream.Test(), QueryPageSize: 10,
				Client: func() (cloud.TiDBCloudClient, error) {
					return cloud.NewClientDelegateWithToken("test-token", cloud.DefaultServerlessEndpoint, config.GetIAMEndpoint())
				},
				NextGenClient: func() (cloud.NextGenClient, error) {
					return cloud.NewNextGenClientDelegateWithToken("test-token", config.GetCloudAPIEndpoint(), "")
				},
			}
			command := createCmd(h, plan)
			command.SetArgs([]string{"--project-name", "my-project", "--display-name", "test-instance", "--region", "aws-us-west-2", "--max-rcu", "10000", "--output", "json"})
			require.NoError(t, command.ExecuteContext(context.Background()))
			require.Equal(t, []string{"GET iam.staging.tidbapi.com/v1beta1/projects", "POST cloud.staging.tidbapi.com/v1beta2/tidbs"}, requests)
		})
	}
}
