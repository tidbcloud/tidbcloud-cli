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

type UpdateAuthorizedNetworkSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *UpdateAuthorizedNetworkSuite) SetupTest() {
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

func (suite *UpdateAuthorizedNetworkSuite) TestUpdateAuthorizedNetworkArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "123456"

	displayName := "test"
	startIPAddress := "0.0.0.0"
	endIPAddress := "1.2.2.2"

	result := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err := json.Unmarshal([]byte(getClusterResultStr), result)
	assert.Nil(err)

	c := &cluster.V1beta1ClusterServicePartialUpdateClusterBodyCluster{
		Endpoints: &cluster.V1beta1ClusterEndpoints{
			Public: &cluster.EndpointsPublic{
				AuthorizedNetworks: []cluster.EndpointsPublicAuthorizedNetwork{
					{
						StartIpAddress: startIPAddress,
						EndIpAddress:   endIPAddress,
						DisplayName:    displayName,
					},
				},
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
			name:         "update authorized network success",
			args:         []string{"--cluster-id", clusterID, "--start-ip-address", startIPAddress, "--end-ip-address", endIPAddress, "--new-display-name", displayName},
			stdoutString: "authorized network is updated\n",
		},
		{
			name: "does not set ip",
			args: []string{"--cluster-id", clusterID, "--start-ip-address", startIPAddress, "--new-display-name", displayName},
			err:  errors.New("required flag(s) \"end-ip-address\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := UpdateCmd(suite.h)
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

func TestUpdateAuthorizedNetworkSuite(t *testing.T) {
	suite.Run(t, new(UpdateAuthorizedNetworkSuite))
}
