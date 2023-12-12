// Copyright 2023 PingCAP, Inc.
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
	"fmt"
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

type SpendingLimitSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *SpendingLimitSuite) SetupTest() {
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

func (suite *SpendingLimitSuite) TestSetSpendingLimit() {
	assert := require.New(suite.T())

	clusterID := "0"
	monthly := 10
	cluster := &serverlessApi.ServerlessServicePartialUpdateClusterParamsBodyCluster{
		SpendingLimit: &serverlessModel.ClusterSpendingLimit{
			Monthly: int32(monthly),
		},
	}
	body := serverlessApi.ServerlessServicePartialUpdateClusterBody{
		Cluster:    cluster,
		UpdateMask: &SpendingLimitMonthlyMask,
	}

	suite.mockClient.On("PartialUpdateCluster", serverlessApi.NewServerlessServicePartialUpdateClusterParams().
		WithClusterClusterID(clusterID).WithBody(body)).
		Return(&serverlessApi.ServerlessServicePartialUpdateClusterOK{}, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "update displayName success",
			args:         []string{"--cluster-id", clusterID, "--monthly", fmt.Sprintf("%d", monthly)},
			stdoutString: fmt.Sprintf("set spending limit to %d cents success\n", monthly),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := SpendingLimitCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func TestSpendingLimitSuite(t *testing.T) {
	suite.Run(t, new(SpendingLimitSuite))
}
