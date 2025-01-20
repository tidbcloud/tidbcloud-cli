// Copyright 2025 PingCAP, Inc.
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

package export

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
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const listResultStr = `{
  "exports": [
    {
      "clusterId": "10289717998856001017",
      "completeTime": "2024-09-04T04:45:41Z",
      "createTime": "2024-09-04T04:45:27Z",
      "createdBy": "apikey-MCTGR3Jv",
      "displayName": "SNAPSHOT_2024-09-04T04:45:23Z",
      "expireTime": "2024-09-06T04:45:41Z",
      "exportId": "exp-q6j5hwy7vzhhfhlx3ilxqti3ay",
      "exportOptions": {
        "compression": "GZIP",
        "database": "*",
        "fileType": "SQL",
        "table": "*"
      },
      "name": "clusters/10289717998856001017/exports/exp-q6j5hwy7vzhhfhlx3ilxqti3ay",
      "snapshotTime": "2024-09-04T04:45:23.189Z",
      "state": "SUCCEEDED",
      "target": {
        "type": "LOCAL"
      },
      "updateTime": "2024-09-04T04:46:44Z"
    }
  ],
  "totalSize": 1
}
`

const listResultMultiPageStr = `{
  "exports": [
    {
      "clusterId": "10289717998856001017",
      "completeTime": "2024-09-04T04:45:41Z",
      "createTime": "2024-09-04T04:45:27Z",
      "createdBy": "apikey-MCTGR3Jv",
      "displayName": "SNAPSHOT_2024-09-04T04:45:23Z",
      "expireTime": "2024-09-06T04:45:41Z",
      "exportId": "exp-q6j5hwy7vzhhfhlx3ilxqti3ay",
      "exportOptions": {
        "compression": "GZIP",
        "database": "*",
        "fileType": "SQL",
        "table": "*"
      },
      "name": "clusters/10289717998856001017/exports/exp-q6j5hwy7vzhhfhlx3ilxqti3ay",
      "snapshotTime": "2024-09-04T04:45:23.189Z",
      "state": "SUCCEEDED",
      "target": {
        "type": "LOCAL"
      },
      "updateTime": "2024-09-04T04:46:44Z"
    },
    {
      "clusterId": "10289717998856001017",
      "completeTime": "2024-09-04T04:45:41Z",
      "createTime": "2024-09-04T04:45:27Z",
      "createdBy": "apikey-MCTGR3Jv",
      "displayName": "SNAPSHOT_2024-09-04T04:45:23Z",
      "expireTime": "2024-09-06T04:45:41Z",
      "exportId": "exp-q6j5hwy7vzhhfhlx3ilxqti3ay",
      "exportOptions": {
        "compression": "GZIP",
        "database": "*",
        "fileType": "SQL",
        "table": "*"
      },
      "name": "clusters/10289717998856001017/exports/exp-q6j5hwy7vzhhfhlx3ilxqti3ay",
      "snapshotTime": "2024-09-04T04:45:23.189Z",
      "state": "SUCCEEDED",
      "target": {
        "type": "LOCAL"
      },
      "updateTime": "2024-09-04T04:46:44Z"
    }
  ],
  "totalSize": 2
}
`

type ListExportsSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ListExportsSuite) SetupTest() {
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

func (suite *ListExportsSuite) TestListExportsArgs() {
	assert := require.New(suite.T())
	pageSize := int32(suite.h.QueryPageSize)
	orderBy := "create_time desc"
	ctx := context.Background()

	body := &export.ListExportsResponse{}
	err := json.Unmarshal([]byte(listResultStr), body)
	assert.Nil(err)
	clusterID := "fake-cluster-id"
	suite.mockClient.On("ListExports", ctx, clusterID, &pageSize, (*string)(nil), &orderBy).Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list exports with default format(json when without tty)",
			args:         []string{"--cluster-id", clusterID},
			stdoutString: listResultStr,
		},
		{
			name:         "list exports with output flag",
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

func (suite *ListExportsSuite) TestListExportsWithMultiPages() {
	assert := require.New(suite.T())
	ctx := context.Background()
	// mock first page
	pageSize := int32(suite.h.QueryPageSize)
	pageToken := "2"
	orderBy := "create_time desc"
	body := &export.ListExportsResponse{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"totalSize": 1`, `"totalSize": 2`)), body)
	assert.Nil(err)
	body.NextPageToken = &pageToken

	clusterID := "fake-cluster-id"
	suite.mockClient.On("ListExports", ctx, clusterID, &pageSize, (*string)(nil), &orderBy).Return(body, nil)

	body2 := &export.ListExportsResponse{}
	err = json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"totalSize": 1`, `"totalSize": 2`)), body2)
	assert.Nil(err)
	suite.mockClient.On("ListExports", ctx, clusterID, &pageSize, &pageToken, &orderBy).Return(body2, nil)

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

func TestListExportSuite(t *testing.T) {
	suite.Run(t, new(ListExportsSuite))
}
