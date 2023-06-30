// Copyright 2023 PingCAP, Inc.
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
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	branchApi "tidbcloud-cli/pkg/tidbcloud/branch/client/branch_service"
	branchModel "tidbcloud-cli/pkg/tidbcloud/branch/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getBranchResultStr = `{
  "annotations": {},
  "cluster_id": "3478958",
  "create_time": "2023-06-05T05:05:33.000Z",
  "delete_time": null,
  "display_name": "ru-test",
  "endpoints": {
    "public_endpoint": {
      "host": "gateway01.us-east-1.dev.shared.aws.tidbcloud.com",
      "port": 4000
    }
  },
  "id": "branch-sm4ee5usauqfgsftywrapd",
  "name": "clusters/3478958/branches/branch-sm4ee5usauqfgsftywrapd",
  "parent_id": "3478958",
  "state": "READY",
  "update_time": "2023-06-05T05:06:41.000Z",
  "usages": {
    "column_storage": "0",
    "request_unit": "1300000",
    "row_storage": "261260139"
  },
  "user_prefix": "49dDUPpoxGXdsY9"
}
`

type DescribeBranchSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *DescribeBranchSuite) SetupTest() {
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

func (suite *DescribeBranchSuite) TestDescribeBranchArgs() {
	assert := require.New(suite.T())

	body := &branchModel.OpenapiBranch{}
	err := json.Unmarshal([]byte(getBranchResultStr), body)
	assert.Nil(err)
	result := &branchApi.GetBranchOK{
		Payload: body,
	}
	clusterID := "12345"
	branchID := "12345"
	suite.mockClient.On("GetBranch", branchApi.NewGetBranchParams().
		WithBranchID(branchID).WithClusterID(clusterID)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "describe branch success",
			args:         []string{"--cluster-id", clusterID, "--branch-id", branchID},
			stdoutString: getBranchResultStr,
		},
		{
			name:         "describe branch with shorthand flag",
			args:         []string{"-c", clusterID, "-b", branchID},
			stdoutString: getBranchResultStr,
		},
		{
			name: "describe branch without required cluster id",
			args: []string{"-b", branchID},
			err:  fmt.Errorf("if any flags in the group [branch-id cluster-id] are set they must all be set; missing [cluster-id]"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DescribeCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func TestDescribeBranchSuite(t *testing.T) {
	suite.Run(t, new(DescribeBranchSuite))
}
