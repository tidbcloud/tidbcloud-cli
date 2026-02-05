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

type ZonesPrivateLinkConnectionSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ZonesPrivateLinkConnectionSuite) SetupTest() {
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

func (suite *ZonesPrivateLinkConnectionSuite) TestZonesPrivateLinkConnectionArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	accountID := "123"
	response := &privatelink.GetAvailabilityZonesResponse{
		AccountId: &accountID,
		AzIds:     []string{"az-1", "az-2"},
	}
	suite.mockClient.On("GetAvailabilityZones", ctx, clusterID).Return(response, nil)

	expectedJson, err := json.MarshalIndent(response, "", "  ")
	assert.NoError(err)
	expectedOutput := string(expectedJson) + "\n"

	tests := []struct {
		name              string
		args              []string
		err               error
		errContains       string
		stdoutString      string
		stderrString      string
		expectMockedCalls bool
	}{
		{
			name:              "zones private link connection success",
			args:              []string{"--cluster-id", clusterID},
			stdoutString:      expectedOutput,
			expectMockedCalls: true,
		},
		{
			name:        "missing required cluster id",
			args:        []string{},
			errContains: "The terminal doesn't support interactive mode, please use non-interactive mode",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ZonesCmd(suite.h)
			cmd.SetContext(ctx)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if tt.errContains != "" {
				assert.ErrorContains(err, tt.errContains)
			} else {
				assert.Equal(tt.err, err)
			}

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.expectMockedCalls {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func TestZonesPrivateLinkConnectionSuite(t *testing.T) {
	suite.Run(t, new(ZonesPrivateLinkConnectionSuite))
}
