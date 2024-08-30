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
	iamClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
	serverlessModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/models"

	"github.com/juju/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UpdateSQLUserSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *UpdateSQLUserSuite) SetupTest() {
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

func (suite *UpdateSQLUserSuite) TestUpdateSQLUserArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	userName := "test"
	builtinRole := util.ADMIN_ROLE
	customRole := []string{"my_role", "my_role2"}
	password := "123"
	userNamePrefix := "4TGJD6zA3Nn2333"
	nonExistRole := "non-exist-role"
	fullUserName := fmt.Sprintf("%s.%s", userNamePrefix, userName)
	customRoleStr := strings.Join(customRole, ",")
	roleStr := fmt.Sprintf("%s,%s", builtinRole, customRoleStr)

	result := &iamClient.ApiSqlUser{}
	err := json.Unmarshal([]byte(getSQLUserResultStr), result)
	assert.Nil(err)

	suite.mockClient.On("GetSQLUser", ctx, clusterID, fullUserName).
		Return(result, nil)

	clusterBody := &serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster{}
	err = json.Unmarshal([]byte(getClusterResultStr), clusterBody)
	assert.Nil(err)
	res := &serverlessApi.ServerlessServiceGetClusterOK{
		Payload: clusterBody,
	}
	suite.mockClient.On("GetCluster", serverlessApi.NewServerlessServiceGetClusterParams().
		WithClusterID(clusterID).WithContext(ctx)).
		Return(res, nil)

	updateBody := &iamClient.ApiUpdateSqlUserReq{
		BuiltinRole: &builtinRole,
		CustomRoles: customRole,
		Password:    &password,
	}

	suite.mockClient.On("UpdateSQLUser", ctx, clusterID, fullUserName, updateBody).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "update SQL user success",
			args:         []string{"--cluster-id", clusterID, "--user", userName, "--password", password, "--role", roleStr},
			stdoutString: fmt.Sprintf("SQL user %s is updated\n", fullUserName),
		},
		{
			name: "no SQL user attribute",
			args: []string{"--cluster-id", clusterID, "--user", userName},
			err:  errors.New("at least one of the flags in the group [password role add-role delete-role] is required"),
		},
		{
			name: "multi role flags",
			args: []string{"--cluster-id", clusterID, "--user", userName, "--role", roleStr, "--delete-role", roleStr},
			err:  errors.New("if any flags in the group [role add-role delete-role] are set none of the others can be; [delete-role role] were all set"),
		},
		{
			name: "delete non-exist role",
			args: []string{"--cluster-id", clusterID, "--user", userName, "--delete-role", nonExistRole},
			err:  errors.New(fmt.Sprintf("role %s doesn't exist in the SQL user", nonExistRole)),
		},
		{
			name: "add built-in role to user with built-in role",
			args: []string{"--cluster-id", clusterID, "--user", userName, "--add-role", builtinRole},
			err:  errors.New("built-in role already exists in the SQL user"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := UpdateCmd(suite.h)
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

func TestUpdateSQLUserSuite(t *testing.T) {
	suite.Run(t, new(UpdateSQLUserSuite))
}
