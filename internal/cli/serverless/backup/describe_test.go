// Copyright 2025 PingCAP, Inc.
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

package backup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/br"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getBackupResult = `{
  "backupId": "289048",
  "clusterId": "10048930788495339885",
  "createTime": "2023-12-15T07:00:00Z",
  "name": "backups/289048"
}
`

type DescribeBackupSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *DescribeBackupSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	suite.mockClient = new(mock.TiDBCloudClient)
	suite.h = &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			return suite.mockClient, nil
		},
		QueryPageSize: 100,
		IOStreams:     iostream.Test(),
	}
}

func (suite *DescribeBackupSuite) TestDescribeBackup() {
	assert := require.New(suite.T())
	ctx := context.Background()

	body := &br.V1beta1Backup{}
	err := json.Unmarshal([]byte(getBackupResult), body)
	assert.Nil(err)
	backupId := "289048"

	suite.mockClient.On("GetBackup", ctx, backupId).
		Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "describe backup success",
			args:         []string{"--backup-id", backupId},
			stdoutString: getBackupResult,
		},
		{
			name: "describe with cluster id should be fail",
			args: []string{"--backup-id", backupId, "--cluster-id", "10048930788495339885"},
			err:  errors.New("unknown flag: --cluster-id"),
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

func TestDescribeBackupSuite(t *testing.T) {
	suite.Run(t, new(DescribeBackupSuite))
}
