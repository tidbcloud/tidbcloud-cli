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

package private_link_connection

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const listResultStr = `{
  "privateLinkConnections": [
    {
      "awsEndpointService": {
        "name": "com.amazonaws.vpce.us-east-1.vpce-svc-123"
      },
      "clusterId": "12345",
      "createTime": "2024-01-01T00:00:00Z",
      "displayName": "plc-test",
      "privateLinkConnectionId": "plc-12345",
      "state": "ACTIVE",
      "type": "AWS_ENDPOINT_SERVICE"
    }
  ],
  "totalSize": 1
}
`

const listResultMultiPageStr = `{
  "privateLinkConnections": [
    {
      "awsEndpointService": {
        "name": "com.amazonaws.vpce.us-east-1.vpce-svc-123"
      },
      "clusterId": "12345",
      "createTime": "2024-01-01T00:00:00Z",
      "displayName": "plc-test",
      "privateLinkConnectionId": "plc-12345",
      "state": "ACTIVE",
      "type": "AWS_ENDPOINT_SERVICE"
    },
    {
      "alicloudEndpointService": {
        "name": "privatelink.alicloud.example"
      },
      "clusterId": "12345",
      "createTime": "2024-01-01T00:00:00Z",
      "displayName": "plc-test-2",
      "privateLinkConnectionId": "plc-67890",
      "state": "ACTIVE",
      "type": "ALICLOUD_ENDPOINT_SERVICE"
    }
  ],
  "totalSize": 2
}
`

const listResultSecondPageStr = `{
  "privateLinkConnections": [
    {
      "alicloudEndpointService": {
        "name": "privatelink.alicloud.example"
      },
      "clusterId": "12345",
      "createTime": "2024-01-01T00:00:00Z",
      "displayName": "plc-test-2",
      "privateLinkConnectionId": "plc-67890",
      "state": "ACTIVE",
      "type": "ALICLOUD_ENDPOINT_SERVICE"
    }
  ],
  "totalSize": 2
}
`

type ListPrivateLinkConnectionsSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ListPrivateLinkConnectionsSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	var pageSize int64 = 10
	suite.mockClient = new(mock.TiDBCloudClient)
	suite.h = &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			return suite.mockClient, nil
		},
		QueryPageSize: pageSize,
		IOStreams:     iostream.Test(),
	}
}

func (suite *ListPrivateLinkConnectionsSuite) TestListPrivateLinkConnectionsArgs() {
	assert := require.New(suite.T())
	pageSize := int32(suite.h.QueryPageSize)
	ctx := context.Background()

	body := &privatelink.ListPrivateLinkConnectionsResponse{}
	err := json.Unmarshal([]byte(listResultStr), body)
	assert.Nil(err)
	clusterID := "12345"
	suite.mockClient.On("ListPrivateLinkConnections", ctx, clusterID, &pageSize, (*string)(nil)).Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list private link connections with default format(json when without tty)",
			args:         []string{"--cluster-id", clusterID},
			stdoutString: listResultStr,
		},
		{
			name:         "list private link connections with output flag",
			args:         []string{"--cluster-id", clusterID, "--output", "json"},
			stdoutString: listResultStr,
		},
		{
			name:         "list private link connections with output shorthand flag",
			args:         []string{"--cluster-id", clusterID, "-o", "json"},
			stdoutString: listResultStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ListCmd(suite.h)
			cmd.SetContext(ctx)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func (suite *ListPrivateLinkConnectionsSuite) TestListPrivateLinkConnectionsWithMultiPages() {
	assert := require.New(suite.T())
	ctx := context.Background()
	pageSize := int32(suite.h.QueryPageSize)
	pageToken := "2"
	body := &privatelink.ListPrivateLinkConnectionsResponse{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"totalSize": 1`, `"totalSize": 2`)), body)
	assert.Nil(err)
	body.NextPageToken = &pageToken

	clusterID := "12345"
	suite.mockClient.On("ListPrivateLinkConnections", ctx, clusterID, &pageSize, (*string)(nil)).Return(body, nil)

	body2 := &privatelink.ListPrivateLinkConnectionsResponse{}
	err = json.Unmarshal([]byte(listResultSecondPageStr), body2)
	assert.Nil(err)
	suite.mockClient.On("ListPrivateLinkConnections", ctx, clusterID, &pageSize, &pageToken).Return(body2, nil)

	cmd := ListCmd(suite.h)
	tests := []struct {
		name         string
		args         []string
		stdoutString string
		stderrString string
	}{
		{
			name:         "query with multi pages",
			args:         []string{"--cluster-id", clusterID, "--output", "json"},
			stdoutString: listResultMultiPageStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd.SetContext(ctx)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Nil(err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			suite.mockClient.AssertExpectations(suite.T())
		})
	}
}

func TestListPrivateLinkConnectionsSuite(t *testing.T) {
	suite.Run(t, new(ListPrivateLinkConnectionsSuite))
}
