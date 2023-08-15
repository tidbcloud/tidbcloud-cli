// Copyright 2022 PingCAP, Inc.
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

package cluster

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/serverless/client/serverless_service"
	serverlessModel "tidbcloud-cli/pkg/tidbcloud/serverless/models"
)

const listResultStr = `{
  "clusters": [
    {
      "clusterId": "3779024",
      "createTime": "2023-07-04T01:56:12.000Z",
      "createdBy": "yuhang.shi@pingcap.com",
      "displayName": "Cluster0",
      "endpoints": {
        "privateEndpoint": {
          "aws": {
            "availabilityZone": [
              "use1-az1"
            ],
            "serviceName": "com.amazonaws.vpce.us-east-1.vpce-svc-03342995daf1bc4d4"
          },
          "host": "gateway01-privatelink.us-east-1.prod.aws.tidbcloud.com",
          "port": 4000
        },
        "publicEndpoint": {
          "host": "gateway01.us-east-1.prod.aws.tidbcloud.com",
          "port": 4000
        }
      },
      "labels": {
        "tidb.cloud/organization": "30018",
        "tidb.cloud/project": "163469"
      },
      "name": "clusters/3779024",
      "region": {
        "displayName": "N. Virginia (us-east-1)",
        "name": "regions/aws-us-east-1",
        "provider": "AWS"
      },
      "state": "ACTIVE",
      "updateTime": "2023-08-03T09:08:07.753Z",
      "userPrefix": "4FNu72xBpLXjFnC",
      "version": "v6.6.0"
    }
  ],
  "totalSize": 1
}
`

const listResultMultiPageStr = `{
  "clusters": [
    {
      "clusterId": "3779024",
      "createTime": "2023-07-04T01:56:12.000Z",
      "createdBy": "yuhang.shi@pingcap.com",
      "displayName": "Cluster0",
      "endpoints": {
        "privateEndpoint": {
          "aws": {
            "availabilityZone": [
              "use1-az1"
            ],
            "serviceName": "com.amazonaws.vpce.us-east-1.vpce-svc-03342995daf1bc4d4"
          },
          "host": "gateway01-privatelink.us-east-1.prod.aws.tidbcloud.com",
          "port": 4000
        },
        "publicEndpoint": {
          "host": "gateway01.us-east-1.prod.aws.tidbcloud.com",
          "port": 4000
        }
      },
      "labels": {
        "tidb.cloud/organization": "30018",
        "tidb.cloud/project": "163469"
      },
      "name": "clusters/3779024",
      "region": {
        "displayName": "N. Virginia (us-east-1)",
        "name": "regions/aws-us-east-1",
        "provider": "AWS"
      },
      "state": "ACTIVE",
      "updateTime": "2023-08-03T09:08:07.753Z",
      "userPrefix": "4FNu72xBpLXjFnC",
      "version": "v6.6.0"
    },
    {
      "clusterId": "3779024",
      "createTime": "2023-07-04T01:56:12.000Z",
      "createdBy": "yuhang.shi@pingcap.com",
      "displayName": "Cluster0",
      "endpoints": {
        "privateEndpoint": {
          "aws": {
            "availabilityZone": [
              "use1-az1"
            ],
            "serviceName": "com.amazonaws.vpce.us-east-1.vpce-svc-03342995daf1bc4d4"
          },
          "host": "gateway01-privatelink.us-east-1.prod.aws.tidbcloud.com",
          "port": 4000
        },
        "publicEndpoint": {
          "host": "gateway01.us-east-1.prod.aws.tidbcloud.com",
          "port": 4000
        }
      },
      "labels": {
        "tidb.cloud/organization": "30018",
        "tidb.cloud/project": "163469"
      },
      "name": "clusters/3779024",
      "region": {
        "displayName": "N. Virginia (us-east-1)",
        "name": "regions/aws-us-east-1",
        "provider": "AWS"
      },
      "state": "ACTIVE",
      "updateTime": "2023-08-03T09:08:07.753Z",
      "userPrefix": "4FNu72xBpLXjFnC",
      "version": "v6.6.0"
    }
  ],
  "totalSize": 2
}
`

type ListClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ListClusterSuite) SetupTest() {
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

func (suite *ListClusterSuite) TestListClusterArgs() {
	assert := require.New(suite.T())
	var page string = "1"

	body := &serverlessModel.V1ListClustersResponse{}
	err := json.Unmarshal([]byte(listResultStr), body)
	assert.Nil(err)
	result := &serverlessApi.ServerlessServiceListClustersOK{
		Payload: body,
	}
	projectID := "12345"
	pageSize := int32(suite.h.QueryPageSize)
	suite.mockClient.On("ListClustersOfProject", serverlessApi.NewServerlessServiceListClustersParams().
		WithProjectID(&projectID).WithPageToken(&page).WithPageSize(&pageSize)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list clusters with default format(json when without tty)",
			args:         []string{projectID},
			stdoutString: listResultStr,
		},
		{
			name:         "list clusters with output flag",
			args:         []string{projectID, "--output", "json"},
			stdoutString: listResultStr,
		},
		{
			name:         "list clusters with output shorthand flag",
			args:         []string{projectID, "-o", "json"},
			stdoutString: listResultStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ListCmd(suite.h)
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

func (suite *ListClusterSuite) TestListClusterWithMultiPages() {
	assert := require.New(suite.T())
	var pageOne = "1"
	var pageTwo = "2"
	suite.h.QueryPageSize = 1
	pageSize := int32(suite.h.QueryPageSize)

	body := &serverlessModel.V1ListClustersResponse{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"totalSize": 1`, `"totalSize": 2`)), body)
	assert.Nil(err)
	result := &serverlessApi.ServerlessServiceListClustersOK{
		Payload: body,
	}
	projectID := "12345"
	suite.mockClient.On("ListClustersOfProject", serverlessApi.NewServerlessServiceListClustersParams().
		WithProjectID(&projectID).WithPageToken(&pageOne).WithPageSize(&pageSize)).
		Return(result, nil)
	suite.mockClient.On("ListClustersOfProject", serverlessApi.NewServerlessServiceListClustersParams().
		WithProjectID(&projectID).WithPageToken(&pageTwo).WithPageSize(&pageSize)).
		Return(result, nil)
	cmd := ListCmd(suite.h)

	tests := []struct {
		name         string
		args         []string
		stdoutString string
		stderrString string
	}{
		{
			name:         "query with multi pages",
			args:         []string{projectID, "--output", "json"},
			stdoutString: listResultMultiPageStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
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

func TestListClusterSuite(t *testing.T) {
	suite.Run(t, new(ListClusterSuite))
}
