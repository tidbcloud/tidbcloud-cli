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

	serverlessApi "tidbcloud-cli/pkg/tidbcloud/serverless/client/serverless_service"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UpdateClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *UpdateClusterSuite) SetupTest() {
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

func (suite *DescribeClusterSuite) TestUpdateClusterArgs() {
	assert := require.New(suite.T())

	displayName := "update_name"
	cluster := &serverlessApi.ServerlessServicePartialUpdateClusterParamsBodyCluster{
		DisplayName: displayName,
	}
	mask := "displayName"
	body := serverlessApi.ServerlessServicePartialUpdateClusterBody{
		Cluster:    cluster,
		UpdateMask: &mask,
	}

	clusterID := "12345"
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
			args:         []string{"--cluster-id", clusterID, "--field", "displayName", "--value", displayName},
			stdoutString: fmt.Sprintf("cluster %s updated\n", clusterID),
		},
		{
			name: "update unsupported field",
			args: []string{"-c", clusterID, "--field", "state", "--value", "running"},
			err:  fmt.Errorf("unsupported update field state"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := UpdateCmd(suite.h)
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

func TestUpdateClusterSuite(t *testing.T) {
	suite.Run(t, new(UpdateClusterSuite))
}
