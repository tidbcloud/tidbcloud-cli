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

package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/stretchr/testify/require"
)

func TestNextGenCapacityResponseRoundTrip(t *testing.T) {
	for _, plan := range []api.TidbServiceListTidbsServicePlanParameter{
		api.TIDBSERVICELISTTIDBSSERVICEPLANPARAMETER_PREMIUM,
		api.TIDBSERVICELISTTIDBSSERVICEPLANPARAMETER_ESSENTIAL_V2,
	} {
		t.Run(string(plan), func(t *testing.T) {
			baseline := "5000"
			if plan == api.TIDBSERVICELISTTIDBSSERVICEPLANPARAMETER_ESSENTIAL_V2 {
				baseline = "2000"
			}
			var instances []string
			for _, tc := range []struct {
				name, fields string
			}{
				{"legacy", `"minRcu":"5000","maxRcu":"20000"`},
				{"max-rcu", `"minRcu":"5000","maxRcu":"20000","capacityMode":"MAX_RCU"`},
				{"elastic", fmt.Sprintf(`"baselineRcu":%q,"capacityMode":"ELASTIC"`, baseline)},
				{"null", fmt.Sprintf(`"minRcu":null,"maxRcu":null,"baselineRcu":%q,"capacityMode":"ELASTIC"`, baseline)},
				{"unspecified", `"capacityMode":"CAPACITY_MODE_UNSPECIFIED"`},
				{"future", `"capacityMode":"FUTURE_MODE"`},
			} {
				body := fmt.Sprintf(`{"tidbId":"123","displayName":"test-instance","regionId":"aws-us-west-2","servicePlan":%q,%s}`, plan, tc.fields)
				instances = append(instances, body)
				t.Run(tc.name, func(t *testing.T) {
					transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
						require.Equal(t, http.MethodGet, r.Method)
						require.Equal(t, "/v1beta2/tidbs/123", r.URL.Path)
						return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
					})
					client, err := newNextGenClientDelegate(transport, "https://example.test")
					require.NoError(t, err)
					result, err := client.GetTiDB(context.Background(), "123")
					require.NoError(t, err)
					output, err := json.Marshal(result)
					require.NoError(t, err)
					require.JSONEq(t, body, string(output))
				})
			}
			t.Run("mixed-list", func(t *testing.T) {
				body := `{"tidbs":[` + strings.Join(instances, ",") + `]}`
				transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
					require.Equal(t, http.MethodGet, r.Method)
					require.Equal(t, "/v1beta2/tidbs", r.URL.Path)
					require.Equal(t, string(plan), r.URL.Query().Get("servicePlan"))
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
				})
				client, err := newNextGenClientDelegate(transport, "https://example.test")
				require.NoError(t, err)
				result, err := client.ListTiDBs(context.Background(), plan, nil, nil)
				require.NoError(t, err)
				require.Len(t, result.Tidbs, len(instances))
				output, err := json.Marshal(result)
				require.NoError(t, err)
				require.JSONEq(t, body, string(output))
			})
		})
	}
}
