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
	"fmt"
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

type GetPrivateLinkConnectionSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *GetPrivateLinkConnectionSuite) SetupTest() {
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

func (suite *GetPrivateLinkConnectionSuite) TestGetPrivateLinkConnectionArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	connectionID := "plc-123"
	connectionType := privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE
	response := &privatelink.PrivateLinkConnection{
		PrivateLinkConnectionId: &connectionID,
		ClusterId:               clusterID,
		DisplayName:             "plc-demo",
		Type:                    connectionType,
		AlicloudEndpointService: privatelink.NewAlicloudEndpointService("cn-hangzhou-endpoint"),
	}

	suite.mockClient.On("GetPrivateLinkConnection", ctx, clusterID, connectionID).Return(response, nil)

	expectedJson, err := json.MarshalIndent(response, "", "  ")
	assert.NoError(err)
	expectedOutput := string(expectedJson) + "\n"

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "get private link connection success",
			args:         []string{"--cluster-id", clusterID, "--private-link-connection-id", connectionID},
			stdoutString: expectedOutput,
		},
		{
			name: "missing required connection id",
			args: []string{"--cluster-id", clusterID},
			err:  fmt.Errorf("required flag(s) \"private-link-connection-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := GetCmd(suite.h)
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

func TestGetPrivateLinkConnectionSuite(t *testing.T) {
	suite.Run(t, new(GetPrivateLinkConnectionSuite))
}
