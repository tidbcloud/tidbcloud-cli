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
	"os"
	"strings"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/branch"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const listResultStr = `{
  "branches": [
    {
      "branchId": "bran-wscjvwen2jajdjiy7hawcebxke",
      "clusterId": "10202848322613926203",
      "createTime": "2023-12-12T10:17:15Z",
      "createdBy": "apikey-MCTGR3Jv",
      "displayName": "t",
      "name": "clusters/10202848322613926203/branches/bran-wscjvwen2jajdjiy7hawcebxke",
      "parentId": "10202848322613926203",
      "state": "ACTIVE",
      "updateTime": "2023-12-12T10:18:24Z"
    }
  ],
  "totalSize": 1
}
`

const listResultMultiPageStr = `{
  "branches": [
    {
      "branchId": "bran-wscjvwen2jajdjiy7hawcebxke",
      "clusterId": "10202848322613926203",
      "createTime": "2023-12-12T10:17:15Z",
      "createdBy": "apikey-MCTGR3Jv",
      "displayName": "t",
      "name": "clusters/10202848322613926203/branches/bran-wscjvwen2jajdjiy7hawcebxke",
      "parentId": "10202848322613926203",
      "state": "ACTIVE",
      "updateTime": "2023-12-12T10:18:24Z"
    },
    {
      "branchId": "bran-wscjvwen2jajdjiy7hawcebxke",
      "clusterId": "10202848322613926203",
      "createTime": "2023-12-12T10:17:15Z",
      "createdBy": "apikey-MCTGR3Jv",
      "displayName": "t",
      "name": "clusters/10202848322613926203/branches/bran-wscjvwen2jajdjiy7hawcebxke",
      "parentId": "10202848322613926203",
      "state": "ACTIVE",
      "updateTime": "2023-12-12T10:18:24Z"
    }
  ],
  "totalSize": 2
}
`

type ListBranchesSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ListBranchesSuite) SetupTest() {
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

func (suite *ListBranchesSuite) TestListBranchesArgs() {
	assert := require.New(suite.T())
	pageSize := int32(suite.h.QueryPageSize)
	ctx := context.Background()

	body := &branch.ListBranchesResponse{}
	err := json.Unmarshal([]byte(listResultStr), body)
	assert.Nil(err)
	clusterID := "12345"
	suite.mockClient.On("ListBranches", ctx, clusterID, &pageSize, (*string)(nil)).Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list branches with default format(json when without tty)",
			args:         []string{"--cluster-id", clusterID},
			stdoutString: listResultStr,
		},
		{
			name:         "list branches with output flag",
			args:         []string{"--cluster-id", clusterID, "--output", "json"},
			stdoutString: listResultStr,
		},
		{
			name:         "list branches with output shorthand flag",
			args:         []string{"--cluster-id", clusterID, "-o", "json"},
			stdoutString: listResultStr,
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
			assert.Equal(tt.err, err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func (suite *ListBranchesSuite) TestListBranchesWithMultiPages() {
	assert := require.New(suite.T())
	ctx := context.Background()
	pageSize := int32(suite.h.QueryPageSize)
	pageToken := "2"
	body := &branch.ListBranchesResponse{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"totalSize": 1`, `"totalSize": 2`)), body)
	assert.Nil(err)
	body.NextPageToken = &pageToken

	clusterID := "12345"
	suite.mockClient.On("ListBranches", ctx, clusterID, &pageSize, (*string)(nil)).Return(body, nil)

	body2 := &branch.ListBranchesResponse{}
	err = json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"totalSize": 1`, `"totalSize": 2`)), body2)
	assert.Nil(err)
	suite.mockClient.On("ListBranches", ctx, clusterID, &pageSize, &pageToken).Return(body2, nil)

	cmd := ListCmd(suite.h)
	tests := []struct {
		name         string
		args         []string
		stdoutString string
		stderrString string
	}{
		{
			name:         "query with multi pages",
			args:         []string{"--cluster-id", clusterID, "--output", "json"},
			stdoutString: listResultMultiPageStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
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

func TestListBranchSuite(t *testing.T) {
	suite.Run(t, new(ListBranchesSuite))
}
