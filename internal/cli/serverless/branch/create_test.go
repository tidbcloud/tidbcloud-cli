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

package branch

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
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/branch"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreateBranchSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *CreateBranchSuite) SetupTest() {
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

func (suite *CreateBranchSuite) TestCreateBranchArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	branchName := "test"
	branchId := "12345"

	createBranchBody := &branch.Branch{
		DisplayName: branchName,
	}
	suite.mockClient.On("CreateBranch", ctx, clusterID, createBranchBody).
		Return(&branch.Branch{
			BranchId: &branchId,
		}, nil)

	body := &branch.Branch{}
	err := json.Unmarshal([]byte(getBranchResultStr), body)
	assert.Nil(err)
	suite.mockClient.On("GetBranch", ctx, clusterID, branchId, branch.BRANCHSERVICEGETBRANCHVIEWPARAMETER_BASIC).Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "create branch success",
			args:         []string{"--cluster-id", clusterID, "--display-name", branchName},
			stdoutString: fmt.Sprintf("... Waiting for branch to be ready\nBranch %s is ready.", branchId),
		},
		{
			name:         "create branch with shorthand flag",
			args:         []string{"-c", clusterID, "--display-name", branchName},
			stdoutString: fmt.Sprintf("... Waiting for branch to be ready\nBranch %s is ready.", branchId),
		},
		{
			name: "without required project id",
			args: []string{"--display-name", branchName},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
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

func TestCreateBranchSuite(t *testing.T) {
	suite.Run(t, new(CreateBranchSuite))
}
