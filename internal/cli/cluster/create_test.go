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

type CreateClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.ApiClient
}

func (suite *CreateClusterSuite) SetupTest() {
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

func (suite *CreateClusterSuite) TestCreateClusterArgs() {
	assert := require.New(suite.T())

	projectID := "12345"
	clusterID := "12345"
	clusterName := "test"
	clusterType := "SERVERLESS"
	cloudProvider := "AWS"
	region := "us-west-1"
	rootPassword := "123456"
	clusterDefBody := &cluster.CreateClusterBody{}

	err := clusterDefBody.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"name": "%s",
			"cluster_type": "%s",
			"cloud_provider": "%s",
			"region": "%s",
			"config" : {
				"root_password": "%s",
				"ip_access_list": [
					{
						"CIDR": "0.0.0.0/0",
						"description": "Allow All"
					}
				]
			}
			}`, clusterName, "DEVELOPER", cloudProvider, region, rootPassword)))
	assert.Nil(err)

	body := &cluster.GetClusterOKBody{}
	err = json.Unmarshal([]byte(getClusterResultStr), body)
	assert.Nil(err)
	res := &cluster.GetClusterOK{
		Payload: body,
	}

	suite.mockClient.On("CreateCluster", cluster.NewCreateClusterParams().
		WithProjectID(projectID).WithBody(*clusterDefBody)).
		Return(&cluster.CreateClusterOK{
			Payload: &cluster.CreateClusterOKBody{
				ID: &clusterID,
			},
		}, nil)
	suite.mockClient.On("GetCluster", cluster.NewGetClusterParams().
		WithProjectID(projectID).WithClusterID(clusterID)).
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
			args:         []string{"--project-id", projectID, "--cluster-name", clusterName, "--cluster-type", clusterType, "--cloud-provider", cloudProvider, "--region", region, "--root-password", rootPassword},
			stdoutString: "... Waiting for cluster to be ready\nCluster 12345 is ready.",
		},
		{
			name:         "create cluster with shorthand flag",
			args:         []string{"-p", projectID, "--cluster-name", clusterName, "--cluster-type", clusterType, "--cloud-provider", cloudProvider, "-r", region, "--root-password", rootPassword},
			stdoutString: "... Waiting for cluster to be ready\nCluster 12345 is ready.",
		},
		{
			name: "without required project id",
			args: []string{"--cluster-name", clusterName, "--cluster-type", clusterType, "--cloud-provider", cloudProvider, "-r", region, "--root-password", rootPassword},
			err:  fmt.Errorf("required flag(s) \"project-id\" not set"),
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
