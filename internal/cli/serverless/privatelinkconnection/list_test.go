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

package privatelinkconnection

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
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ListPrivateLinkConnectionSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ListPrivateLinkConnectionSuite) SetupTest() {
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

func (suite *ListPrivateLinkConnectionSuite) TestListPrivateLinkConnectionsArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()
	pageSize := int32(suite.h.QueryPageSize)

	clusterID := "12345"
	connectionID := "plc-123"
	connection := privatelink.PrivateLinkConnection{
		PrivateLinkConnectionId: &connectionID,
		ClusterId:               clusterID,
		DisplayName:             "plc-demo",
		Type:                    privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE,
		AwsEndpointService:      privatelink.NewAwsEndpointService("com.amazonaws.vpce.svc-123"),
	}
	totalSize := int64(1)
	body := &privatelink.ListPrivateLinkConnectionsResponse{
		PrivateLinkConnections: []privatelink.PrivateLinkConnection{connection},
		TotalSize:              &totalSize,
	}

	suite.mockClient.On(
		"ListPrivateLinkConnections",
		ctx,
		clusterID,
		&pageSize,
		(*string)(nil),
		(*privatelink.PrivateLinkConnectionServiceListPrivateLinkConnectionsStateParameter)(nil),
	).Return(body, nil)

	expectedJson, err := json.MarshalIndent(body, "", "  ")
	assert.NoError(err)
	expectedOutput := string(expectedJson) + "\n"

	tests := []struct {
		name         string
		args         []string
		stdoutString string
		stderrString string
	}{
		{
			name:         "list private link connections",
			args:         []string{"--cluster-id", clusterID},
			stdoutString: expectedOutput,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ListCmd(suite.h)
			cmd.SetContext(ctx)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Nil(err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			suite.mockClient.AssertExpectations(suite.T())
		})
	}
}

func TestListPrivateLinkConnectionSuite(t *testing.T) {
	suite.Run(t, new(ListPrivateLinkConnectionSuite))
}
