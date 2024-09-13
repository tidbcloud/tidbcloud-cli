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

package sqluser

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
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const listResultStr = `{
  "sqlUsers": [
    {
      "authMethod": "mysql_native_password",
      "builtinRole": "role_admin",
      "userName": "4TGJD6zA3NnUiz4.12222"
    },
    {
      "authMethod": "mysql_native_password",
      "builtinRole": "role_readonly",
      "customRoles": [
        "my_role"
      ],
      "userName": "4TGJD6zA3NnUiz4.123"
    }
  ]
}
`

type ListSQLUserSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
	pageSize   int64
}

func (suite *ListSQLUserSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	suite.pageSize = 1
	suite.mockClient = new(mock.TiDBCloudClient)
	suite.h = &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			return suite.mockClient, nil
		},
		QueryPageSize: suite.pageSize,
		IOStreams:     iostream.Test(),
	}
}

func (suite *ListSQLUserSuite) TestListSQLUserArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()
	pageSize := int32(suite.pageSize)
	var pageToken *string

	clusterID := "12345"

	result := &iam.ApiListSqlUsersRsp{}
	err := json.Unmarshal([]byte(listResultStr), result)
	assert.Nil(err)
	suite.mockClient.On("ListSQLUsers", ctx, clusterID, &pageSize, pageToken).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list SQL user success",
			args:         []string{"--cluster-id", clusterID},
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

func TestListSQLUserSuite(t *testing.T) {
	suite.Run(t, new(ListSQLUserSuite))
}
