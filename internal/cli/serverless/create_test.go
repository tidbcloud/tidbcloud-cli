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

type CreateClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *CreateClusterSuite) SetupTest() {
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

func (suite *CreateClusterSuite) TestCreateClusterArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	projectID := "12345"
	clusterID := "12345"
	clusterName := "test"
	regionName := "regions/aws-us-west-1"
	v1Cluster := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{
		DisplayName: clusterName,
		Region: cluster.Commonv1beta1Region{
			Name: &regionName,
		},
		Labels: &map[string]string{"tidb.cloud/project": projectID},
	}

	body := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err := json.Unmarshal([]byte(getClusterResultStr), body)
	assert.Nil(err)

	suite.mockClient.On("CreateCluster", ctx, v1Cluster).
		Return(&cluster.TidbCloudOpenApiserverlessv1beta1Cluster{
			ClusterId: &clusterID,
		}, nil)
	suite.mockClient.On("GetCluster", ctx, clusterID).Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "create cluster success",
			args:         []string{"--project-id", projectID, "--display-name", clusterName, "--region", regionName},
			stdoutString: "... Waiting for cluster to be ready\nCluster 12345 is ready.",
		},
		{
			name:         "create cluster with shorthand flag",
			args:         []string{"-p", projectID, "-n", clusterName, "-r", regionName},
			stdoutString: "... Waiting for cluster to be ready\nCluster 12345 is ready.",
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
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func (suite *CreateClusterSuite) TestCreateClusterWithoutProject() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	clusterName := "test"
	regionName := "regions/aws-us-west-1"

	v1ClusterWithoutProject := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{
		DisplayName: clusterName,
		Region: cluster.Commonv1beta1Region{
			Name: &regionName,
		},
	}

	body := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err := json.Unmarshal([]byte(getClusterResultStr), body)
	assert.Nil(err)

	suite.mockClient.On("CreateCluster", ctx, v1ClusterWithoutProject).
		Return(&cluster.TidbCloudOpenApiserverlessv1beta1Cluster{
			ClusterId: &clusterID,
		}, nil)
	suite.mockClient.On("GetCluster", ctx, clusterID).Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "without project id",
			args:         []string{"--display-name", clusterName, "-r", regionName},
			stdoutString: "... Waiting for cluster to be ready\nCluster 12345 is ready.",
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
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func TestCreateClusterSuite(t *testing.T) {
	suite.Run(t, new(CreateClusterSuite))
}
