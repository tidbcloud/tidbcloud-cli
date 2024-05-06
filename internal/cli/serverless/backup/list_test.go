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

package backup

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
	brApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"
	brModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const listResultStr = `{
  "backups": [
    {
      "backupId": "289048",
      "clusterId": "10048930788495339885",
      "createTime": "2023-12-15T07:00:00.000Z",
      "name": "backups/289048"
    }
  ],
  "totalSize": 1
}
`

const listResultMultiPageStr = `{
  "backups": [
    {
      "backupId": "289048",
      "clusterId": "10048930788495339885",
      "createTime": "2023-12-15T07:00:00.000Z",
      "name": "backups/289048"
    },
    {
      "backupId": "289048",
      "clusterId": "10048930788495339885",
      "createTime": "2023-12-15T07:00:00.000Z",
      "name": "backups/289048"
    }
  ],
  "totalSize": 2
}
`

type ListBackupSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
	pageSize   int64
}

func (suite *ListBackupSuite) SetupTest() {
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

func (suite *ListBackupSuite) TestListBackups() {
	assert := require.New(suite.T())
	ctx := context.Background()

	body := &brModel.V1beta1ListBackupsResponse{}
	err := json.Unmarshal([]byte(listResultStr), body)
	assert.Nil(err)
	result := &brApi.BackupRestoreServiceListBackupsOK{
		Payload: body,
	}
	pageSize := int32(suite.pageSize)
	clusterID := "10048930788495339885"
	suite.mockClient.On("ListBackups", brApi.NewBackupRestoreServiceListBackupsParams().
		WithClusterID(clusterID).WithPageSize(&pageSize).WithContext(ctx)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list backups with default format(json when without tty)",
			args:         []string{"--cluster-id", clusterID},
			stdoutString: listResultStr,
		},
		{
			name:         "list backups with output flag",
			args:         []string{"--cluster-id", clusterID, "--output", "json"},
			stdoutString: listResultStr,
		},
		{
			name:         "list backups with shorthand flag",
			args:         []string{"-c", clusterID, "-o", "json"},
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

func (suite *ListBackupSuite) TestListBackupsWithMultiPages() {
	assert := require.New(suite.T())
	
	ctx := context.Background()

	pageSize := int32(suite.pageSize)
	pageToken := "2"
	clusterID := "10048930788495339885"

	body := &brModel.V1beta1ListBackupsResponse{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"total": 1`, `"total": 2`)), body)
	assert.Nil(err)
	body.NextPageToken = pageToken
	suite.mockClient.On("ListBackups", brApi.NewBackupRestoreServiceListBackupsParams().
		WithClusterID(clusterID).WithPageSize(&pageSize).WithContext(ctx)).
		Return(&brApi.BackupRestoreServiceListBackupsOK{Payload: body}, nil)

	body2 := &brModel.V1beta1ListBackupsResponse{}
	err = json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"total": 1`, `"total": 2`)), body2)
	assert.Nil(err)
	suite.mockClient.On("ListBackups", brApi.NewBackupRestoreServiceListBackupsParams().
		WithClusterID(clusterID).WithPageToken(&pageToken).WithPageSize(&pageSize).WithContext(ctx)).
		Return(&brApi.BackupRestoreServiceListBackupsOK{Payload: body2}, nil)

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
		cmd := ListCmd(suite.h)
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

func TestListBackupSuite(t *testing.T) {
	suite.Run(t, new(ListBackupSuite))
}
