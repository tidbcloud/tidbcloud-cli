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
	"fmt"
	"os"
	"testing"

	mockTool "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type ResumeMigrationSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ResumeMigrationSuite) SetupTest() {
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

func (suite *ResumeMigrationSuite) TestResumeMigrations() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "c123"
	migrationID := "m456"

	suite.mockClient.On(
		"ResumeMigration",
		ctx,
		clusterID,
		migrationID,
	).Return(nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "resume migration in non-interactive mode",
			args:         []string{"--cluster-id", clusterID, "--migration-id", migrationID},
			stdoutString: fmt.Sprintf("migration %s resumed\n", migrationID),
		},
		{
			name: "resume migration without required flags in non-interactive terminal",
			args: []string{},
			err:  fmt.Errorf("The terminal doesn't support interactive mode, please use non-interactive mode"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ResumeCmd(suite.h)
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
				suite.mockClient.AssertNotCalled(suite.T(), "ResumeMigration", mockTool.Anything)
			}
		})
	}
}

func TestResumeMigrationSuite(t *testing.T) {
	suite.Run(t, new(ResumeMigrationSuite))
}
