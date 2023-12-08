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
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"

	serverlessApi "tidbcloud-cli/pkg/tidbcloud/serverless/client/serverless_service"
	serverlessModel "tidbcloud-cli/pkg/tidbcloud/serverless/models"

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

	projectID := "12345"
	clusterID := "12345"
	clusterName := "test"
	regionName := "regions/aws-us-west-1"
	cloudProvider := "AWS"
	region := "us-west-1"
	v1Cluster := &serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster{
		DisplayName: &clusterName,
		Region: &serverlessModel.TidbCloudOpenApiserverlessv1beta1Region{
			Name: &regionName,
		},
		Labels: map[string]string{"tidb.cloud/project": projectID},
	}

	body := &serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err := json.Unmarshal([]byte(getClusterResultStr), body)
	assert.Nil(err)
	res := &serverlessApi.ServerlessServiceGetClusterOK{
		Payload: body,
	}

	suite.mockClient.On("CreateCluster", serverlessApi.NewServerlessServiceCreateClusterParams().
		WithCluster(v1Cluster)).
		Return(&serverlessApi.ServerlessServiceCreateClusterOK{
			Payload: &serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster{
				ClusterID: clusterID,
			},
		}, nil)
	suite.mockClient.On("GetCluster", serverlessApi.NewServerlessServiceGetClusterParams().
		WithClusterID(clusterID)).
		Return(res, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "create cluster success",
			args:         []string{"--project-id", projectID, "--cluster-name", clusterName, "--cloud-provider", cloudProvider, "--region", region},
			stdoutString: "... Waiting for cluster to be ready\nCluster 12345 is ready.",
		},
		{
			name:         "create cluster with shorthand flag",
			args:         []string{"-p", projectID, "--cluster-name", clusterName, "--cloud-provider", cloudProvider, "-r", region},
			stdoutString: "... Waiting for cluster to be ready\nCluster 12345 is ready.",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := CreateCmd(suite.h)
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

	clusterID := "12345"
	clusterName := "test"
	regionName := "regions/aws-us-west-1"
	cloudProvider := "AWS"
	region := "us-west-1"

	v1ClusterWithoutProject := &serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster{
		DisplayName: &clusterName,
		Region: &serverlessModel.TidbCloudOpenApiserverlessv1beta1Region{
			Name: &regionName,
		},
	}

	body := &serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err := json.Unmarshal([]byte(getClusterResultStr), body)
	assert.Nil(err)
	res := &serverlessApi.ServerlessServiceGetClusterOK{
		Payload: body,
	}

	suite.mockClient.On("CreateCluster", serverlessApi.NewServerlessServiceCreateClusterParams().
		WithCluster(v1ClusterWithoutProject)).
		Return(&serverlessApi.ServerlessServiceCreateClusterOK{
			Payload: &serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster{
				ClusterID: clusterID,
			},
		}, nil)
	suite.mockClient.On("GetCluster", serverlessApi.NewServerlessServiceGetClusterParams().
		WithClusterID(clusterID)).
		Return(res, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "without project id",
			args:         []string{"--cluster-name", clusterName, "--cloud-provider", cloudProvider, "-r", region},
			stdoutString: "... Waiting for cluster to be ready\nCluster 12345 is ready.",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := CreateCmd(suite.h)
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
