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

type DeleteAuthorizedNetworkSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *DeleteAuthorizedNetworkSuite) SetupTest() {
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

func (suite *DeleteAuthorizedNetworkSuite) TestDeleteAuthorizedNetworkArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "123456"

	startIPAddress := "0.0.0.0"
	endIPAddress := "1.2.2.2"
	wrongIPAddress := "0.0.0.22"

	result := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err := json.Unmarshal([]byte(getClusterResultStr), result)
	assert.Nil(err)

	c := &cluster.V1beta1ClusterServicePartialUpdateClusterBodyCluster{
		Endpoints: &cluster.V1beta1ClusterEndpoints{
			Public: &cluster.EndpointsPublic{
				AuthorizedNetworks: []cluster.EndpointsPublicAuthorizedNetwork{},
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
			name:         "delete authorized network success",
			args:         []string{"--cluster-id", clusterID, "--start-ip-address", startIPAddress, "--end-ip-address", endIPAddress, "--force"},
			stdoutString: fmt.Sprintf("authorized network %s-%s is deleted\n", startIPAddress, endIPAddress),
		},
		{
			name: "ip range does not exist",
			args: []string{"--cluster-id", clusterID, "--start-ip-address", startIPAddress, "--end-ip-address", wrongIPAddress, "--force"},
			err:  errors.New(fmt.Sprintf("authorized network %s-%s not found", startIPAddress, wrongIPAddress)),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DeleteCmd(suite.h)
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

func TestDeleteAuthorizedNetworkSuite(t *testing.T) {
	suite.Run(t, new(DeleteAuthorizedNetworkSuite))
}
