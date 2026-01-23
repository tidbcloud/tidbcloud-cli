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

package authorizednetwork

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/juju/errors"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getClusterResultStr = `{
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
                }
            ]
        }
    }
}`

type CreateAuthorizedNetworkSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *CreateAuthorizedNetworkSuite) SetupTest() {
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

func (suite *CreateAuthorizedNetworkSuite) TestCreateAuthorizedNetworkArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "123456"

	displayName := "test"
	startIPAddress := "0.0.0.0"
	endIPAddress := "1.1.1.1"
	wrongIPAddress := "0.0.0.256"

	result := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err := json.Unmarshal([]byte(getClusterResultStr), result)
	assert.Nil(err)

	c := &cluster.V1beta1ClusterServicePartialUpdateClusterBodyCluster{
		Endpoints: &cluster.V1beta1ClusterEndpoints{
			Public: &cluster.EndpointsPublic{
				AuthorizedNetworks: append(result.Endpoints.Public.AuthorizedNetworks, cluster.EndpointsPublicAuthorizedNetwork{
					StartIpAddress: startIPAddress,
					EndIpAddress:   endIPAddress,
					DisplayName:    displayName,
				}),
			},
		},
	}

	body := &cluster.V1beta1ClusterServicePartialUpdateClusterBody{
		Cluster:    c,
		UpdateMask: AuthorizedNetworkMask,
	}

	suite.mockClient.On("GetCluster", ctx, clusterID, cluster.CLUSTERSERVICEGETCLUSTERVIEWPARAMETER_BASIC).
		Return(result, nil)
	suite.mockClient.On("PartialUpdateCluster", ctx, clusterID, body).
		Return(&cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}, nil)

	assert.Nil(err)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "create authorized network success",
			args:         []string{"--cluster-id", clusterID, "--display-name", displayName, "--start-ip-address", startIPAddress, "--end-ip-address", endIPAddress},
			stdoutString: fmt.Sprintf("authorized network %s-%s is created\n", startIPAddress, endIPAddress),
		},
		{
			name: "wrong ip range",
			args: []string{"--cluster-id", clusterID, "--display-name", displayName, "--start-ip-address", startIPAddress, "--end-ip-address", wrongIPAddress},
			err:  errors.New(fmt.Sprintf("invalid IPv4 address: %s", wrongIPAddress)),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := CreateCmd(suite.h)
			cmd.SetContext(ctx)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if err != nil {
				assert.NotNil(tt.err)
				assert.Contains(err.Error(), tt.err.Error())
			}
			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())

			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func TestCreateAuthorizedNetworkSuite(t *testing.T) {
	suite.Run(t, new(CreateAuthorizedNetworkSuite))
}
