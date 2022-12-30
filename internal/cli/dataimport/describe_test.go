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

package dataimport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	importOp "tidbcloud-cli/pkg/tidbcloud/import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/import/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getImportResultStr = `{
  "all_completed_tables": [],
  "cluster_id": "120466",
  "completed_percent": 0,
  "completed_tables": 1,
  "created_at": "2022-12-28T08:18:11.000Z",
  "current_tables": [],
  "data_format": "CSV",
  "elapsed_time_seconds": 10,
  "id": "120117",
  "import_create_req": {
    "aws_role_arn": "arn",
    "cluster_id": "120466",
    "csv_format": {
      "delimiter": "^",
      "not_null": true,
      "null": "@",
      "separator": "$"
    },
    "data_format": "CSV",
    "project_id": "0",
    "source_url": "s3"
  },
  "message": "[Lightning:Restore:ErrEncodeKV]encode kv error in file jianxu.test1.000000000001.csv:0 at offset 12: column count mismatch, expected 2, got 1",
  "pending_tables": 0,
  "processed_source_data_size": "0",
  "source_url": "s3",
  "status": "FAILED",
  "total_files": 0,
  "total_size": "39",
  "total_tables_count": 1
}
`

type DescribeImportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *DescribeImportSuite) SetupTest() {
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

func (suite *DescribeImportSuite) TestDescribeImportArgs() {
	assert := require.New(suite.T())

	body := &importModel.OpenapiGetImportResp{}
	err := json.Unmarshal([]byte(getImportResultStr), body)
	assert.Nil(err)
	result := &importOp.GetImportOK{
		Payload: body,
	}

	projectID := "12345"
	clusterID := "12345"
	importID := "12345"
	suite.mockClient.On("GetImport", importOp.NewGetImportParams().
		WithProjectID(projectID).WithClusterID(clusterID).WithID(importID)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "describe import success",
			args:         []string{"--project-id", projectID, "--cluster-id", clusterID, "--import-id", importID},
			stdoutString: getImportResultStr,
		},
		{
			name:         "describe import with shorthand flag",
			args:         []string{"-p", projectID, "-c", clusterID, "--import-id", importID},
			stdoutString: getImportResultStr,
		},
		{
			name: "describe import without required project id",
			args: []string{"-c", clusterID, "--import-id", importID},
			err:  fmt.Errorf("required flag(s) \"project-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DescribeCmd(suite.h)
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

func TestDescribeImportSuite(t *testing.T) {
	suite.Run(t, new(DescribeImportSuite))
}
