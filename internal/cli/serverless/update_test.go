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
	"fmt"
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

func (suite *UpdateClusterSuite) TestUpdateClusterArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	displayName := "update_name"
	c := &cluster.RequiredTheClusterToBeUpdated{
		DisplayName: &displayName,
	}
	mask := "displayName"
	body := &cluster.V1beta1ServerlessServicePartialUpdateClusterBody{
		Cluster:    c,
		UpdateMask: mask,
	}

	clusterID := "12345"
	suite.mockClient.On("PartialUpdateCluster", ctx, clusterID, body).Return(&cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "update displayName success",
			args:         []string{"--cluster-id", clusterID, "--display-name", displayName},
			stdoutString: fmt.Sprintf("cluster %s updated\n", clusterID),
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
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func (suite *UpdateClusterSuite) TestUpdateLabels() {
	assert := require.New(suite.T())
	ctx := context.Background()

	labels := "{\"labels\":\"values\"}"
	mask := "labels"

	labelsMap := make(map[string]string)
	_ = json.Unmarshal([]byte(labels), &labelsMap)
	c := &cluster.RequiredTheClusterToBeUpdated{
		Labels: &labelsMap,
	}
	body := &cluster.V1beta1ServerlessServicePartialUpdateClusterBody{
		Cluster:    c,
		UpdateMask: mask,
	}
	clusterID := "12345"
	suite.mockClient.On("PartialUpdateCluster", ctx, clusterID, body).Return(&cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "update labels success",
			args:         []string{"-c", clusterID, "--labels", "{\"labels\":\"values\"}"},
			stdoutString: fmt.Sprintf("cluster %s updated\n", clusterID),
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
