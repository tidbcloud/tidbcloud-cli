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

package cloud_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"
)

func TestRetrievePrivateLinkConnectionsPagination(t *testing.T) {
	second, third, empty := "page-2", "page-3", ""
	for _, tc := range []struct {
		name       string
		nextTokens []*string
		wantIDs    []string
		wantError  bool
	}{
		{name: "single page", nextTokens: []*string{nil}, wantIDs: []string{"connection-1"}},
		{name: "empty final token", nextTokens: []*string{&empty}, wantIDs: []string{"connection-1"}},
		{name: "two pages", nextTokens: []*string{&second, nil}, wantIDs: []string{"connection-1", "connection-2"}},
		{name: "three pages", nextTokens: []*string{&second, &third, nil}, wantIDs: []string{"connection-1", "connection-2", "connection-3"}},
		{name: "repeated token", nextTokens: []*string{&second, &second}, wantError: true},
		{name: "token cycle", nextTokens: []*string{&second, &third, &second}, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := mock.NewTiDBCloudClient(t)
			ctx := t.Context()
			pageSize := int32(1)
			for i, nextToken := range tc.nextTokens {
				var requestToken *string
				if i > 0 {
					requestToken = tc.nextTokens[i-1]
				}
				id := fmt.Sprintf("connection-%d", i+1)
				response := &privatelink.ListPrivateLinkConnectionsResponse{
					PrivateLinkConnections: []privatelink.PrivateLinkConnection{{
						PrivateLinkConnectionId: &id,
						ClusterId:               "cluster-1",
						DisplayName:             id,
						Type:                    privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE,
					}},
					NextPageToken: nextToken,
				}
				client.On("ListPrivateLinkConnections", ctx, "cluster-1", &pageSize, requestToken).
					Return(response, nil).Once()
			}
			connections, err := cloud.RetrievePrivateLinkConnections(ctx, "cluster-1", pageSize, client)
			if tc.wantError {
				require.ErrorContains(t, err, "repeated next page token")
				require.Nil(t, connections)
				return
			}
			require.NoError(t, err)
			ids := make([]string, 0, len(connections))
			for _, connection := range connections {
				ids = append(ids, connection.GetPrivateLinkConnectionId())
			}
			require.Equal(t, tc.wantIDs, ids)
		})
	}
}
