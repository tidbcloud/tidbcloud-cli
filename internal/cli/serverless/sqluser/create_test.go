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
	"fmt"
	"os"
	"strings"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

	"github.com/juju/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getSQLUserResultStr = `{
	"authMethod": "mysql_native_password",
	"builtinRole": "role_admin",
	"customRoles": ["my_role"],
	"userName": "4TGJD6zA3Nn2333.test"
}`

const getClusterResultStr = `{
	"clusterId": "12345",
	"userPrefix": "4TGJD6zA3Nn2333",
	"displayName": "test",
	"region": {
    	"displayName": "Singapore (ap-southeast-1)",
        "name": "regions/aws-ap-southeast-1",
        "provider": "aws"
      }
}`

type CreateSQLUserSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *CreateSQLUserSuite) SetupTest() {
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

func (suite *CreateSQLUserSuite) TestCreateSQLUserArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	userName := "test"
	builtinRole := util.ADMIN_ROLE
	customRole := []string{"my_role"}
	password := "123"
	userNamePrefix := "4TGJD6zA3Nn2333"
	fullUserName := fmt.Sprintf("%s.%s", userNamePrefix, userName)
	customRoleStr := strings.Join(customRole, ",")
	roleStr := fmt.Sprintf("%s,%s", builtinRole, customRoleStr)

	authMethod := util.MYSQLNATIVEPASSWORD
	autoPrefix := DefaultAutoPrefix
	createSQLUserBody := &iam.ApiCreateSqlUserReq{
		UserName:    &userName,
		BuiltinRole: &builtinRole,
		CustomRoles: customRole,
		Password:    &password,
		AuthMethod:  &authMethod,
		AutoPrefix:  &autoPrefix,
	}
	result := &iam.ApiSqlUser{}
	err := json.Unmarshal([]byte(getSQLUserResultStr), result)
	assert.Nil(err)

	suite.mockClient.On("CreateSQLUser", ctx, clusterID, createSQLUserBody).
		Return(result, nil)

	res := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err = json.Unmarshal([]byte(getClusterResultStr), res)
	assert.Nil(err)
	suite.mockClient.On("GetCluster", ctx, clusterID).Return(res, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "delete SQL user success",
			args:         []string{"--cluster-id", clusterID, "--user", userName, "--password", password, "--role", roleStr},
			stdoutString: fmt.Sprintf("SQL user %s is created\n", fullUserName),
		},
		{
			name: "multi built-in roles",
			args: []string{"--cluster-id", clusterID, "--user", userName, "--password", password, "--role", fmt.Sprintf("%s,%s", util.ADMIN_ROLE, util.READONLY_ROLE)},
			err:  errors.New("only one built-in role is allowed"),
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

			if err != nil {
				assert.NotNil(tt.err)
				assert.Contains(err.Error(), tt.err.Error())
			}
			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())

			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func TestCreateSQLUserSuite(t *testing.T) {
	suite.Run(t, new(CreateSQLUserSuite))
}
