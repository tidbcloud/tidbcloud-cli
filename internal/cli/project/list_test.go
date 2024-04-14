// Copyright 2022 PingCAP, Inc.
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

package project

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	iamApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/client/account"
	iamModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const resultStr = `{
  "projects": [
    {
      "cluster_count": 1,
      "create_timestamp": "1640076859",
      "id": "1372813089189381287",
      "name": "default project",
      "org_id": "1372813089189621285",
      "user_count": 1
    }
  ]
}
`

const resultPageOne = `{
	"projects": [
	  {
		"cluster_count": 1,
		"create_timestamp": "1640076859",
		"id": "1372813089189381287",
		"name": "default project",
		"org_id": "1372813089189621285",
		"user_count": 1
	  }
	],
	"nextPageToken": "next_token"
  }
  `

const resultMultiPageStr = `{
  "projects": [
    {
      "cluster_count": 1,
      "create_timestamp": "1640076859",
      "id": "1372813089189381287",
      "name": "default project",
      "org_id": "1372813089189621285",
      "user_count": 1
    },
    {
      "cluster_count": 1,
      "create_timestamp": "1640076859",
      "id": "1372813089189381287",
      "name": "default project",
      "org_id": "1372813089189621285",
      "user_count": 1
    }
  ]
}
`

type ListProjectSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ListProjectSuite) SetupTest() {
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

func (suite *ListProjectSuite) TestListProjectArgs() {
	assert := require.New(suite.T())

	body := &iamModel.APIListProjectsRsp{}
	err := json.Unmarshal([]byte(resultStr), body)
	assert.Nil(err)
	result := &iamApi.GetV1beta1ProjectsOK{
		Payload: body,
	}
	suite.mockClient.On("ListProjects", iamApi.NewGetV1beta1ProjectsParams().
		WithPageSize(&suite.h.QueryPageSize)).Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list projects with default format(json when without tty)",
			args:         []string{},
			stdoutString: resultStr,
		},
		{
			name:         "list projects with output flag",
			args:         []string{"--output", "json"},
			stdoutString: resultStr,
		},
		{
			name:         "list projects with output shorthand flag",
			args:         []string{"-o", "json"},
			stdoutString: resultStr,
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

func (suite *ListProjectSuite) TestListProjectWithMultiPages() {
	assert := require.New(suite.T())
	suite.h.QueryPageSize = 1
	nextPageToken := "next_token"

	body1 := &iamModel.APIListProjectsRsp{}
	err := json.Unmarshal([]byte(resultPageOne), body1)
	assert.Nil(err)
	resultPageOne := &iamApi.GetV1beta1ProjectsOK{
		Payload: body1,
	}
	body2 := &iamModel.APIListProjectsRsp{}
	err = json.Unmarshal([]byte(resultStr), body2)
	assert.Nil(err)
	resultPageTwo := &iamApi.GetV1beta1ProjectsOK{
		Payload: body2,
	}
	suite.mockClient.On("ListProjects", iamApi.NewGetV1beta1ProjectsParams().
		WithPageSize(&suite.h.QueryPageSize)).Return(resultPageOne, nil)
	suite.mockClient.On("ListProjects", iamApi.NewGetV1beta1ProjectsParams().
		WithPageSize(&suite.h.QueryPageSize).WithPageToken(&nextPageToken)).Return(resultPageTwo, nil)
	tests := []struct {
		name         string
		args         []string
		stdoutString string
		stderrString string
	}{
		{
			name:         "query with multi pages",
			args:         []string{"--output", "json"},
			stdoutString: resultMultiPageStr,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ListCmd(suite.h)
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

func TestListProjectSuite(t *testing.T) {
	suite.Run(t, new(ListProjectSuite))
}
