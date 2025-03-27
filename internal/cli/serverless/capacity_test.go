// Copyright 2025 PingCAP, Inc.
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

type CapacitySuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *CapacitySuite) SetupTest() {
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

func (suite *CapacitySuite) TestSetCapacity() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "0"
	var minRcu int64 = 10000
	var maxRcu int64 = 20000
	body := &cluster.V1beta1ServerlessServicePartialUpdateClusterBody{
		Cluster: &cluster.RequiredTheClusterToBeUpdated{
			AutoScaling: &cluster.V1beta1ClusterAutoScaling{},
		},
	}
	body.UpdateMask = CapacityMask
	body.Cluster.AutoScaling.MaxRcu = &maxRcu
	body.Cluster.AutoScaling.MinRcu = &minRcu
	suite.mockClient.On("PartialUpdateCluster", ctx, clusterID, body).Return(&cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "update capacity success",
			args:         []string{"--cluster-id", clusterID, "--min-rcu", fmt.Sprintf("%d", minRcu), "--max-rcu", fmt.Sprintf("%d", maxRcu)},
			stdoutString: fmt.Sprintf("set capacity to [%d, %d] success\n", minRcu, maxRcu),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := CapacityCmd(suite.h)
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

func TestCapacitySuite(t *testing.T) {
	suite.Run(t, new(CapacitySuite))
}
