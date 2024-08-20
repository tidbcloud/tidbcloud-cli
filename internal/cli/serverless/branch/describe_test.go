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

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	branchApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client/branch_service"
	branchModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getBranchResultStr = `{
  "annotations": {
    "tidb.cloud/has-set-password": "false"
  },
  "branchId": "bran-fgwdnpasmrahnh5iozqawnmijq",
  "clusterId": "10202848322613926203",
  "createTime": "2023-12-11T09:41:44.000Z",
  "createdBy": "yuhang.shi@pingcap.com",
  "displayName": "t2",
  "endpoints": {
    "public": {
      "host": "gateway01.us-east-1.dev.shared.aws.tidbcloud.com",
      "port": 4000
    }
  },
  "name": "clusters/10202848322613926203/branches/bran-fgwdnpasmrahnh5iozqawnmijq",
  "parentId": "10202848322613926203",
  "parentTimestamp": "0001-01-01T00:00:00.000Z",
  "state": "ACTIVE",
  "updateTime": "2023-12-11T09:44:05.000Z",
  "usage": {
    "requestUnit": "0",
    "rowStorage": 951526
  },
  "userPrefix": "yxfrrVaa55wvBKE"
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
	ctx := context.Background()

	body := &branchModel.V1beta1Branch{}
	err := json.Unmarshal([]byte(getBranchResultStr), body)
	assert.Nil(err)
	result := &branchApi.BranchServiceGetBranchOK{
		Payload: body,
	}
	clusterID := "10202848322613926203"
	branchID := "bran-fgwdnpasmrahnh5iozqawnmijq"
	suite.mockClient.On("GetBranch", branchApi.NewBranchServiceGetBranchParams().
		WithBranchID(branchID).WithClusterID(clusterID).WithContext(ctx)).
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
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DescribeCmd(suite.h)
			cmd.SetContext(ctx)
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
