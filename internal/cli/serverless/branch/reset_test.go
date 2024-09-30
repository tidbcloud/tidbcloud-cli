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

type ResetBranchSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ResetBranchSuite) SetupTest() {
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

func (suite *ResetBranchSuite) TestResetBranchArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	branchID := "12345"
	suite.mockClient.On("ResetBranch", ctx, clusterID, branchID).Return(&branch.Branch{}, nil)

	body := &branch.Branch{}
	err := json.Unmarshal([]byte(getBranchResultStr), body)
	assert.Nil(err)
	suite.mockClient.On("GetBranch", ctx, clusterID, branchID, branch.BRANCHSERVICEGETBRANCHVIEWPARAMETER_BASIC).Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "reset branch success",
			args:         []string{"--cluster-id", clusterID, "--branch-id", branchID, "--force"},
			stdoutString: fmt.Sprintf("... Waiting for branch to be ready\nBranch %s is ready.", branchID),
		},
		{
			name: "reset branch without force",
			args: []string{"--cluster-id", clusterID, "--branch-id", branchID},
			err:  fmt.Errorf("the terminal doesn't support prompt, please run with --force to reset the branch"),
		},
		{
			name: "reset branch without required branch id",
			args: []string{"-c", clusterID, "--force"},
			err:  fmt.Errorf("required flag(s) \"branch-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ResetCmd(suite.h)
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

func TestResetBranchSuite(t *testing.T) {
	suite.Run(t, new(ResetBranchSuite))
}
