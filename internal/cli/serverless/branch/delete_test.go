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
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	branchApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client/branch_service"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DeleteBranchSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *DeleteBranchSuite) SetupTest() {
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

func (suite *DeleteBranchSuite) TestDeleteBranchArgs() {
	assert := require.New(suite.T())

	clusterID := "12345"
	branchID := "12345"
	suite.mockClient.On("DeleteBranch", branchApi.NewBranchServiceDeleteBranchParams().
		WithBranchID(branchID).WithClusterID(clusterID)).
		Return(&branchApi.BranchServiceDeleteBranchOK{}, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "delete branch success",
			args:         []string{"--cluster-id", clusterID, "--branch-id", branchID, "--force"},
			stdoutString: fmt.Sprintf("branch %s deleted\n", branchID),
		},
		{
			name: "delete branch without force",
			args: []string{"--cluster-id", clusterID, "--branch-id", branchID},
			err:  fmt.Errorf("the terminal doesn't support prompt, please run with --force to delete the branch"),
		},
		{
			name:         "delete branch with simple flag",
			args:         []string{"-c", clusterID, "-b", branchID, "--force"},
			stdoutString: fmt.Sprintf("branch %s deleted\n", branchID),
		},
		{
			name: "delete branch without required branch id",
			args: []string{"-c", clusterID, "--force"},
			err:  fmt.Errorf("required flag(s) \"branch-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DeleteCmd(suite.h)
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

func TestDeleteBranchSuite(t *testing.T) {
	suite.Run(t, new(DeleteBranchSuite))
}
