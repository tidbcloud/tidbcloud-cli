// Copyright 2024 PingCAP, Inc.
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

package serverless

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
	serverlessModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getClusterResultStr = `{
  "annotations": {
    "tidb.cloud/has-set-password": "false"
  },
  "automatedBackupPolicy": {
    "retentionDays": 31,
    "startTime": "15:00"
  },
  "clusterId": "10058425682284910921",
  "createTime": "2023-08-15T10:08:21.911Z",
  "createdBy": "yuhang.shi@pingcap.com",
  "displayName": "Cluster0",
  "endpoints": {
    "private": {
      "aws": {
        "availabilityZone": [
          "use1-az1"
        ],
        "serviceName": "com.amazonaws.vpce.us-east-1.vpce-svc-03342995daf1bc4d4"
      },
      "host": "gateway01-privatelink.us-east-1.prod.aws.tidbcloud.com",
      "port": 4000
    },
    "public": {
      "host": "gateway01.us-east-1.prod.aws.tidbcloud.com",
      "port": 4000
    }
  },
  "labels": {
    "tidb.cloud/organization": "30018",
    "tidb.cloud/project": "163469"
  },
  "name": "clusters/10058425682284910921",
  "region": {
    "displayName": "N. Virginia (us-east-1)",
    "name": "regions/aws-us-east-1",
    "provider": "AWS"
  },
  "spendingLimit": {},
  "state": "ACTIVE",
  "updateTime": "2023-08-15T10:08:21.911Z",
  "usage": {
    "requestUnit": "0"
  },
  "userPrefix": "28cDWcUJJiewaQ7",
  "version": "v6.6.0"
}
`

type DescribeClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *DescribeClusterSuite) SetupTest() {
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

func (suite *DescribeClusterSuite) TestDescribeClusterArgs() {
	assert := require.New(suite.T())

	body := &serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err := json.Unmarshal([]byte(getClusterResultStr), body)
	assert.Nil(err)
	result := &serverlessApi.ServerlessServiceGetClusterOK{
		Payload: body,
	}
	clusterID := "12345"
	suite.mockClient.On("GetCluster", serverlessApi.NewServerlessServiceGetClusterParams().
		WithClusterID(clusterID)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "describe cluster success",
			args:         []string{"--cluster-id", clusterID},
			stdoutString: getClusterResultStr,
		},
		{
			name:         "describe cluster with shorthand flag",
			args:         []string{"-c", clusterID},
			stdoutString: getClusterResultStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DescribeCmd(suite.h)
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

func TestDescribeClusterSuite(t *testing.T) {
	suite.Run(t, new(DescribeClusterSuite))
}
