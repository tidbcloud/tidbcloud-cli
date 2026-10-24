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
	"github.com/stretchr/testify/require"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNextGenUpdateBaselineOnWire(t *testing.T) {
	calls := 0
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/v1beta2/tidbs/123", r.URL.Path)
		data, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"baselineRcu":"12500"}`, string(data))
		var decoded api.TheTiDBCloudPremiumInstanceToUpdate
		require.NoError(t, json.Unmarshal(data, &decoded))
		require.Equal(t, "12500", decoded.GetBaselineRcu())
		require.NotContains(t, decoded.AdditionalProperties, "baselineRcu")
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"displayName":"test-instance","regionId":"aws-us-west-2","servicePlan":"Premium","capacityMode":"ELASTIC","baselineRcu":"12500"}`)), Request: r}, nil
	})
	client, err := newNextGenClientDelegate(transport, "https://example.test")
	require.NoError(t, err)
	body := api.NewTheTiDBCloudPremiumInstanceToUpdate()
	body.SetBaselineRcu("12500")
	result, err := client.UpdateTiDB(context.Background(), "123", body)
	require.NoError(t, err)
	require.Equal(t, "12500", result.GetBaselineRcu())
	require.Equal(t, api.TIDBCAPACITYMODE_ELASTIC, result.GetCapacityMode())
	require.Equal(t, 1, calls)
}
