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

package private_link_connection

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
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DeletePrivateLinkConnectionSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *DeletePrivateLinkConnectionSuite) SetupTest() {
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

func (suite *DeletePrivateLinkConnectionSuite) TestDeletePrivateLinkConnectionArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	plcID := "plc-12345"
	suite.mockClient.On("DeletePrivateLinkConnection", ctx, clusterID, plcID).Return(&privatelink.PrivateLinkConnection{}, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "delete private link connection success",
			args:         []string{"--cluster-id", clusterID, "--private-link-connection-id", plcID, "--force"},
			stdoutString: fmt.Sprintf("private link connection %s deleted\n", plcID),
		},
		{
			name: "delete private link connection without force",
			args: []string{"--cluster-id", clusterID, "--private-link-connection-id", plcID},
			err:  fmt.Errorf("the terminal doesn't support prompt, please run with --force to delete the private link connection"),
		},
		{
			name:         "delete private link connection with shorthand flag",
			args:         []string{"-c", clusterID, "--private-link-connection-id", plcID, "--force"},
			stdoutString: fmt.Sprintf("private link connection %s deleted\n", plcID),
		},
		{
			name: "delete private link connection without required id",
			args: []string{"-c", clusterID, "--force"},
			err:  fmt.Errorf("required flag(s) \"private-link-connection-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DeleteCmd(suite.h)
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

func TestDeletePrivateLinkConnectionSuite(t *testing.T) {
	suite.Run(t, new(DeletePrivateLinkConnectionSuite))
}
