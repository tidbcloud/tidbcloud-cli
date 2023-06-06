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
	"os"
	"strings"
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

const listResultStr = `{
  "branches": [
    {
      "cluster_id": "3478958",
      "create_time": "2023-06-05T05:05:33.000Z",
      "delete_time": null,
      "display_name": "ru-test",
      "id": "branch-sm4ee5usauqfgsftywrapd",
      "name": "clusters/3478958/branches/branch-sm4ee5usauqfgsftywrapd",
      "parent_id": "3478958",
      "state": "READY",
      "update_time": "2023-06-05T05:06:41.000Z"
    }
  ],
  "total": 1
}
`

const listResultMultiPageStr = `{
  "branches": [
    {
      "cluster_id": "3478958",
      "create_time": "2023-06-05T05:05:33.000Z",
      "delete_time": null,
      "display_name": "ru-test",
      "id": "branch-sm4ee5usauqfgsftywrapd",
      "name": "clusters/3478958/branches/branch-sm4ee5usauqfgsftywrapd",
      "parent_id": "3478958",
      "state": "READY",
      "update_time": "2023-06-05T05:06:41.000Z"
    },
    {
      "cluster_id": "3478958",
      "create_time": "2023-06-05T05:05:33.000Z",
      "delete_time": null,
      "display_name": "ru-test",
      "id": "branch-sm4ee5usauqfgsftywrapd",
      "name": "clusters/3478958/branches/branch-sm4ee5usauqfgsftywrapd",
      "parent_id": "3478958",
      "state": "READY",
      "update_time": "2023-06-05T05:06:41.000Z"
    }
  ],
  "total": 2
}
`

type ListBranchSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ListBranchSuite) SetupTest() {
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

func (suite *ListBranchSuite) TestListBranchesArgs() {
	assert := require.New(suite.T())
	var page int64 = 1

	body := &branchModel.OpenapiListBranchesResp{}
	err := json.Unmarshal([]byte(listResultStr), body)
	assert.Nil(err)
	result := &branchApi.ListBranchesOK{
		Payload: body,
	}
	clusterID := "12345"
	suite.mockClient.On("ListBranches", branchApi.NewListBranchesParams().
		WithClusterID(clusterID).WithPageToken(&page).WithPageSize(&suite.h.QueryPageSize)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list branches with default format(json when without tty)",
			args:         []string{clusterID},
			stdoutString: listResultStr,
		},
		{
			name:         "list branches with output flag",
			args:         []string{clusterID, "--output", "json"},
			stdoutString: listResultStr,
		},
		{
			name:         "list branches with output shorthand flag",
			args:         []string{clusterID, "-o", "json"},
			stdoutString: listResultStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ListCmd(suite.h)
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

func (suite *ListBranchSuite) TestListBranchesWithMultiPages() {
	assert := require.New(suite.T())
	var pageOne int64 = 1
	var pageTwo int64 = 2
	suite.h.QueryPageSize = 1

	body := &branchModel.OpenapiListBranchesResp{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"total": 1`, `"total": 2`)), body)
	assert.Nil(err)
	result := &branchApi.ListBranchesOK{
		Payload: body,
	}
	clusterID := "12345"
	suite.mockClient.On("ListBranches", branchApi.NewListBranchesParams().
		WithClusterID(clusterID).WithPageToken(&pageOne).WithPageSize(&suite.h.QueryPageSize)).
		Return(result, nil)
	suite.mockClient.On("ListBranches", branchApi.NewListBranchesParams().
		WithClusterID(clusterID).WithPageToken(&pageTwo).WithPageSize(&suite.h.QueryPageSize)).
		Return(result, nil)
	cmd := ListCmd(suite.h)

	tests := []struct {
		name         string
		args         []string
		stdoutString string
		stderrString string
	}{
		{
			name:         "query with multi pages",
			args:         []string{clusterID, "--output", "json"},
			stdoutString: listResultMultiPageStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
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

func TestListBranchSuite(t *testing.T) {
	suite.Run(t, new(ListBranchSuite))
}
