// Copyright 2025 PingCAP, Inc.
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
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

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

type DescribeMigrationSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *DescribeMigrationSuite) SetupTest() {
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

func (suite *DescribeMigrationSuite) TestDescribeMigrations() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "c123"
	migrationID := "m456"
	mode := pkgmigration.TASKMODE_ALL
	state := pkgmigration.MIGRATIONSTATE_RUNNING
	createdAt := time.Now()

	mig := &pkgmigration.Migration{
		MigrationId: aws.String(migrationID),
		DisplayName: aws.String("test-migration"),
		Mode:        &mode,
		State:       &state,
		CreateTime:  &createdAt,
	}

	suite.mockClient.On(
		"GetMigration",
		ctx,
		clusterID,
		migrationID,
	).Return(mig, nil)

	resultJSON, err := json.MarshalIndent(mig, "", "  ")
	assert.Nil(err)
	expected := string(resultJSON) + "\n"

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "describe migration in non-interactive mode",
			args:         []string{"--cluster-id", clusterID, "--migration-id", migrationID},
			stdoutString: expected,
		},
		{
			name: "describe migration without required flags in non-interactive terminal",
			args: []string{},
			err:  fmt.Errorf("The terminal doesn't support interactive mode, please use non-interactive mode"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DescribeCmd(suite.h)
			cmd.SetContext(ctx)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if tt.err != nil {
				assert.EqualError(err, tt.err.Error())
			} else {
				assert.NoError(err)
			}

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			} else {
				suite.mockClient.AssertNotCalled(suite.T(), "GetMigration", mockTool.Anything)
			}
		})
	}
}

func TestDescribeMigrationSuite(t *testing.T) {
	suite.Run(t, new(DescribeMigrationSuite))
}
