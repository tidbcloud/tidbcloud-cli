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
	"strings"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const listResultStr = `{
  "imports": [
    {
      "all_completed_tables": [],
      "cluster_id": "121026",
      "completed_percent": 100,
      "completed_tables": 1,
      "created_at": "2023-01-10T10:32:01.000Z",
      "creation_details": {
        "cluster_id": "121026",
        "csv_format": {
          "backslash_escape": true,
          "delimiter": "\"",
          "header": true,
          "null": "\n",
          "separator": ","
        },
        "data_format": "CSV",
        "file_name": "a.csv",
        "project_id": "0",
        "target_table": {
          "schema": "test",
          "table": "yxxxx"
        },
        "type": "LOCAL"
      },
      "current_tables": [],
      "data_format": "CSV",
      "elapsed_time_seconds": 35,
      "id": "120295",
      "message": "",
      "pending_tables": 0,
      "post_import_completed_percent": 100,
      "processed_source_data_size": "36",
      "status": "COMPLETED",
      "total_files": 0,
      "total_size": "36",
      "total_tables_count": 1
    }
  ],
  "total": "1"
}
`

const listResultMultiPageStr = `{
  "imports": [
    {
      "all_completed_tables": [],
      "cluster_id": "121026",
      "completed_percent": 100,
      "completed_tables": 1,
      "created_at": "2023-01-10T10:32:01.000Z",
      "creation_details": {
        "cluster_id": "121026",
        "csv_format": {
          "backslash_escape": true,
          "delimiter": "\"",
          "header": true,
          "null": "\n",
          "separator": ","
        },
        "data_format": "CSV",
        "file_name": "a.csv",
        "project_id": "0",
        "target_table": {
          "schema": "test",
          "table": "yxxxx"
        },
        "type": "LOCAL"
      },
      "current_tables": [],
      "data_format": "CSV",
      "elapsed_time_seconds": 35,
      "id": "120295",
      "message": "",
      "pending_tables": 0,
      "post_import_completed_percent": 100,
      "processed_source_data_size": "36",
      "status": "COMPLETED",
      "total_files": 0,
      "total_size": "36",
      "total_tables_count": 1
    },
    {
      "all_completed_tables": [],
      "cluster_id": "121026",
      "completed_percent": 100,
      "completed_tables": 1,
      "created_at": "2023-01-10T10:32:01.000Z",
      "creation_details": {
        "cluster_id": "121026",
        "csv_format": {
          "backslash_escape": true,
          "delimiter": "\"",
          "header": true,
          "null": "\n",
          "separator": ","
        },
        "data_format": "CSV",
        "file_name": "a.csv",
        "project_id": "0",
        "target_table": {
          "schema": "test",
          "table": "yxxxx"
        },
        "type": "LOCAL"
      },
      "current_tables": [],
      "data_format": "CSV",
      "elapsed_time_seconds": 35,
      "id": "120295",
      "message": "",
      "pending_tables": 0,
      "post_import_completed_percent": 100,
      "processed_source_data_size": "36",
      "status": "COMPLETED",
      "total_files": 0,
      "total_size": "36",
      "total_tables_count": 1
    }
  ],
  "total": "2"
}
`

type ListImportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ListImportSuite) SetupTest() {
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

func (suite *ListImportSuite) TestListImportArgs() {
	assert := require.New(suite.T())
	var page int32 = 1
	var pageSize = int32(suite.h.QueryPageSize)

	body := &import_operations.
	err := json.Unmarshal([]byte(listResultStr), body)
	assert.Nil(err)
	result := &importOp.ListImportsOK{
		Payload: body,
	}
	projectID := "12345"
	clusterID := "12345"
	suite.mockClient.On("ListImportTasks", importOp.NewListImportsParams().
		WithProjectID(projectID).WithClusterID(clusterID).WithPage(&page).WithPageSize(&pageSize)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list imports with default format(json when without tty)",
			args:         []string{"--project-id", projectID, "--cluster-id", clusterID},
			stdoutString: listResultStr,
		},
		{
			name:         "list imports with output flag",
			args:         []string{"--project-id", projectID, "--cluster-id", clusterID, "--output", "json"},
			stdoutString: listResultStr,
		},
		{
			name:         "list imports with output shorthand flag",
			args:         []string{"-p", projectID, "-c", clusterID, "-o", "json"},
			stdoutString: listResultStr,
		},
		{
			name: "list imports without required project id",
			args: []string{"-c", clusterID},
			err:  fmt.Errorf("required flag(s) \"project-id\" not set"),
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

func (suite *ListImportSuite) TestListImportWithMultiPages() {
	assert := require.New(suite.T())
	var pageOne int32 = 1
	var pageTwo int32 = 2
	suite.h.QueryPageSize = 1
	var pageSize = int32(suite.h.QueryPageSize)

	body := &importModel.OpenapiListImportsResp{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"total": "1"`, `"total": "2"`)), body)
	assert.Nil(err)
	result := &importOp.ListImportsOK{
		Payload: body,
	}
	projectID := "12345"
	clusterID := "12345"
	suite.mockClient.On("ListImportTasks", importOp.NewListImportsParams().
		WithProjectID(projectID).WithClusterID(clusterID).WithPage(&pageOne).WithPageSize(&pageSize)).
		Return(result, nil)
	suite.mockClient.On("ListImportTasks", importOp.NewListImportsParams().
		WithProjectID(projectID).WithClusterID(clusterID).WithPage(&pageTwo).WithPageSize(&pageSize)).
		Return(result, nil)
	cmd := ListCmd(suite.h)

	tests := []struct {
		name         string
		args         []string
		stdoutString string
		stderrString string
	}{
		{
			name:         "query with multi pages",
			args:         []string{"-p", projectID, "-c", clusterID, "--output", "json"},
			stdoutString: listResultMultiPageStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
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

func TestListImportSuite(t *testing.T) {
	suite.Run(t, new(ListImportSuite))
}
