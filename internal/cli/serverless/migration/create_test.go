// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package migration

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	mockTool "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	pkgmigration "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/migration"
)

type CreateMigrationSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *CreateMigrationSuite) SetupTest() {
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

func (suite *CreateMigrationSuite) TestCreateMigration() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "c123"
	migrationID := "mig-1"
	displayName := "my-migration"

	configPath := suite.writeTempConfig(validMigrationConfig())

	suite.mockClient.On(
		"CreateMigration",
		ctx,
		clusterID,
		mockTool.MatchedBy(func(body *pkgmigration.MigrationServiceCreateMigrationBody) bool {
			return body != nil &&
				body.DisplayName == displayName &&
				body.Mode == pkgmigration.TASKMODE_ALL &&
				len(body.Sources) == 1 &&
				body.Target.User == "migration_user"
		}),
	).Return(&pkgmigration.Migration{MigrationId: aws.String(migrationID)}, nil)

	cmd := CreateCmd(suite.h)
	cmd.SetContext(ctx)
	suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
	suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
	cmd.SetArgs([]string{"--cluster-id", clusterID, "--display-name", displayName, "--config-file", configPath})

	err := cmd.Execute()
	assert.NoError(err)
	assert.Equal("migration "+displayName+"("+migrationID+") created\n", suite.h.IOStreams.Out.(*bytes.Buffer).String())
	assert.Equal("", suite.h.IOStreams.Err.(*bytes.Buffer).String())
	suite.mockClient.AssertExpectations(suite.T())
}

func (suite *CreateMigrationSuite) TestCreateMigrationInvalidInputs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	validPath := suite.writeTempConfig(validMigrationConfig())
	blankPath := suite.writeTempConfig(" ")
	invalidJSONPath := suite.writeTempConfig("{invalid")
	invalidModePath := suite.writeTempConfig(`{ "mode": "invalid", "target": {"user":"u","password":"p"}, "sources": [{"sourceType":"MYSQL","connProfile":{"connType":"PUBLIC","host":"h","port":3306,"user":"u","password":"p"}}] }`)

	tests := []struct {
		name        string
		args        []string
		errContains string
	}{
		{
			name:        "empty display name",
			args:        []string{"--cluster-id", "c1", "--display-name", "   ", "--config-file", validPath},
			errContains: "display name is required",
		},
		{
			name:        "empty config path",
			args:        []string{"--cluster-id", "c1", "--display-name", "name", "--config-file", blankPath},
			errContains: "migration config is required",
		},
		{
			name:        "invalid json",
			args:        []string{"--cluster-id", "c1", "--display-name", "name", "--config-file", invalidJSONPath},
			errContains: "invalid migration definition JSON",
		},
		{
			name:        "invalid mode",
			args:        []string{"--cluster-id", "c1", "--display-name", "name", "--config-file", invalidModePath},
			errContains: "invalid mode",
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

			assert.Error(err)
			assert.Contains(err.Error(), tt.errContains)
			suite.mockClient.AssertNotCalled(suite.T(), "CreateMigration", mockTool.Anything)
		})
	}
}

func (suite *CreateMigrationSuite) writeTempConfig(content string) string {
	f, err := os.CreateTemp("", "migration-config.json")
	suite.Require().NoError(err)
	suite.Require().NoError(os.WriteFile(f.Name(), []byte(content), 0o600))
	return f.Name()
}

func validMigrationConfig() string {
	return `{
  "mode": "ALL",
  "target": {
    "user": "migration_user",
    "password": "Passw0rd!"
  },
  "sources": [
    {
      "sourceType": "MYSQL",
      "connProfile": {
        "connType": "PUBLIC",
        "host": "10.0.0.8",
        "port": 3306,
        "user": "dm_sync_user",
        "password": "Passw0rd!"
      }
    }
  ]
}`
}

func TestCreateMigrationSuite(t *testing.T) {
	suite.Run(t, new(CreateMigrationSuite))
}
