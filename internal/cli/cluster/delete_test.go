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
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"

	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DeleteClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.ApiClient
}

func (suite *DeleteClusterSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	var pageSize int64 = 10
	suite.mockClient = new(mock.ApiClient)
	suite.h = &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			return suite.mockClient, nil
		},
		QueryPageSize: pageSize,
		IOStreams:     iostream.Test(),
	}
}

func (suite *DeleteClusterSuite) TestDeleteClusterArgs() {
	assert := require.New(suite.T())

	projectID := "12345"
	clusterID := "12345"
	suite.mockClient.On("DeleteCluster", cluster.NewDeleteClusterParams().
		WithProjectID(projectID).WithClusterID(clusterID)).
		Return(&cluster.DeleteClusterOK{}, nil)
	suite.mockClient.On("GetCluster", cluster.NewGetClusterParams().
		WithProjectID(projectID).WithClusterID(clusterID)).
		Return(nil, &cluster.GetClusterNotFound{})

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "delete cluster success",
			args:         []string{"--project-id", projectID, "--cluster-id", clusterID, "--force"},
			stdoutString: "cluster deleted\n",
		},
		{
			name: "delete cluster without force",
			args: []string{"--project-id", projectID, "--cluster-id", clusterID},
			err:  fmt.Errorf("the terminal doesn't support prompt, please run with --force to delete the cluster"),
		},
		{
			name:         "delete cluster with output flag",
			args:         []string{"-p", projectID, "-c", clusterID, "--force"},
			stdoutString: "cluster deleted\n",
		},
		{
			name: "delete cluster without required project id",
			args: []string{"-c", clusterID, "--force"},
			err:  fmt.Errorf("required flag(s) \"project-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DeleteCmd(suite.h)
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

func TestDeleteClusterSuite(t *testing.T) {
	suite.Run(t, new(DeleteClusterSuite))
}
