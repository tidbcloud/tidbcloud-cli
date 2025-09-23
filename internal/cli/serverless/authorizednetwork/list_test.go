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

package authorizednetwork

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const listResultStr = `[
  {
    "displayName": "123456",
    "endIpAddress": "1.2.2.2",
    "startIpAddress": "0.0.0.0"
  },
  {
    "displayName": "456789",
    "endIpAddress": "3.3.3.3",
    "startIpAddress": "0.0.0.0"
  }
]
`

const listClusterResultStr = `{
    "name": "clusters/123456",
    "clusterId": "123456",
    "displayName": "test",
    "region": {
        "name": "regions/aws-us-west-2",
        "regionId": "us-west-2",
        "cloudProvider": "aws",
        "displayName": "Oregon (us-west-2)",
        "provider": "aws"
    },
    "endpoints": {
        "public": {
            "host": "gateway01.us-west-2.prod.aws.tidbcloud.com",
            "port": 4000,
            "disabled": false,
            "authorizedNetworks": [
                {
					"startIpAddress": "0.0.0.0",
					"endIpAddress": "1.2.2.2",
					"displayName": "123456"
				},
				{
					"startIpAddress": "0.0.0.0",
					"endIpAddress": "3.3.3.3",
					"displayName": "456789"
				}
            ]
        }
    }
}`

type ListAuthorizedNetworkSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
	pageSize   int64
}

func (suite *ListAuthorizedNetworkSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	suite.pageSize = 1
	suite.mockClient = new(mock.TiDBCloudClient)
	suite.h = &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			return suite.mockClient, nil
		},
		QueryPageSize: suite.pageSize,
		IOStreams:     iostream.Test(),
	}
}

func (suite *ListAuthorizedNetworkSuite) TestListAuthorizedNetworkArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"

	result := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err := json.Unmarshal([]byte(listClusterResultStr), result)
	assert.Nil(err)

	suite.mockClient.On("GetCluster", ctx, clusterID, cluster.CLUSTERSERVICEGETCLUSTERVIEWPARAMETER_BASIC).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list authorized networks success",
			args:         []string{"--cluster-id", clusterID},
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

func TestListAuthorizedNetworkSuite(t *testing.T) {
	suite.Run(t, new(ListAuthorizedNetworkSuite))
}
