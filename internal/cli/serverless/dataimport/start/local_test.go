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

package start

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/aws/s3"
	"tidbcloud-cli/internal/service/cloud"
	importOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"

	mockTool "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	fileName = "test.csv"
)

type LocalImportSuite struct {
	suite.Suite
	h            *internal.Helper
	mockClient   *mock.TiDBCloudClient
	mockUploader *mock.Uploader
}

func (suite *LocalImportSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	var pageSize int64 = 10
	suite.mockClient = new(mock.TiDBCloudClient)
	suite.mockUploader = new(mock.Uploader)
	suite.h = &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			return suite.mockClient, nil
		},
		Uploader: func(client cloud.TiDBCloudClient) s3.Uploader {
			return suite.mockUploader
		},
		QueryPageSize: pageSize,
		IOStreams:     iostream.Test(),
	}

	err := os.WriteFile(fileName, []byte("1,2,3"), 0644)
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *LocalImportSuite) TearDownTest() {
	err := os.Remove(fileName)
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *LocalImportSuite) TestLocalImportArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	uploadID := "upl-sadads"
	importID := "imp-asdasd"
	body := &importModel.V1beta1Import{}
	err := json.Unmarshal([]byte(fmt.Sprintf(`{
  "allCompletedTables": [
    {
      "result": "SUCCESS",
      "tableName": "test.ttt"
    }
  ],
  "clusterId": "12345",
  "completedPercent": 100,
  "completedTables": 1,
  "createTime": "2024-04-01T06:39:50.000Z",
  "creationDetails": {
    "clusterId": "12345",
    "dataFormat": "CSV",
    "importOptions": {
      "csvFormat": {
        "backslashEscape": true,
        "delimiter": "\"",
        "header": true,
        "null": "\\N",
        "separator": ","
      }
    },
    "target": {
      "local": {
        "fileName": "a.csv",
        "targetTable": {
          "schema": "test",
          "table": "yxxxx"
        }
      },
      "type": "LOCAL"
    },
    "type": "LOCAL"
  },
  "currentTables": [],
  "dataFormat": "CSV",
  "elapsedTimeSeconds": 14,
  "id": "%s",
  "postImportCompletedPercent": 100,
  "processedSourceDataSize": "37",
  "state": "COMPLETED",
  "totalSize": "37",
  "totalTablesCount": 1
}
`, importID)), body)
	assert.Nil(err)
	result := &importOp.ImportServiceCreateImportOK{
		Payload: body,
	}

	dataFormat := "CSV"
	targetDatabase := "test"
	targetTable := "test"
	clusterID := "12345"

	suite.mockUploader.On("Upload", ctx, mockTool.MatchedBy(func(keys *s3.PutObjectInput) bool {
		assert.Equal(fileName, *keys.FileName)
		assert.Equal(targetDatabase, *keys.DatabaseName)
		assert.Equal(targetTable, *keys.TableName)
		assert.Equal(clusterID, *keys.ClusterID)
		return true
	})).Return(uploadID, nil)
	suite.mockUploader.On("SetConcurrency", 5).Return(nil)
	suite.mockUploader.On("SetPartSize", int64(5*1024*1024)).Return(nil)

	reqBody := importOp.ImportServiceCreateImportBody{}
	err = reqBody.UnmarshalBinary([]byte(fmt.Sprintf(`{
    "clusterId": "12345",
    "dataFormat": "%s",
    "importOptions": {
      "csvFormat": {
        "backslashEscape": true,
        "delimiter": "\"",
        "header": true,
        "null": "\\N",
        "separator": ","
      }
    },
    "target": {
      "local": {
        "fileName": "%s",
        "targetTable": {
          "schema": "%s",
          "table": "%s"
        },
		"uploadID": "%s"
      },
      "type": "LOCAL"
    },
    "type": "LOCAL"
  }`, dataFormat, fileName, targetDatabase, targetTable, uploadID)))
	assert.Nil(err)

	suite.mockClient.On("CreateImport", importOp.NewImportServiceCreateImportParams().
		WithClusterID(clusterID).WithBody(reqBody).WithContext(ctx)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
	}{
		{
			name:         "start import success",
			args:         []string{fileName, "--cluster-id", clusterID, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import with unsupported data format",
			args: []string{fileName, "--cluster-id", clusterID, "--data-format", "yaml", "--target-database", targetDatabase, "--target-table", targetTable},
			err:  fmt.Errorf("data format yaml is not supported, please use one of [\"CSV\"]"),
		},
		{
			name:         "start import with shorthand flag",
			args:         []string{fileName, "-c", clusterID, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import without required cluster id",
			args: []string{fileName, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
		},
		{
			name: "start import without required file path",
			args: []string{"-c", clusterID, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable},
			err:  fmt.Errorf("missing argument <file-path>"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := LocalCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			cmd.SetContext(ctx)
			err = cmd.Execute()
			if err != nil {
				assert.Contains(err.Error(), tt.err.Error())
			}
			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

func (suite *LocalImportSuite) TestLocalImportCSVFormat() {
	assert := require.New(suite.T())
	ctx := context.Background()

	uploadID := "upl-sadads"
	importID := "imp-asdasd"

	dataFormat := "CSV"
	targetDatabase := "test"
	targetTable := "test"
	clusterID := "12345"
	suite.mockUploader.On("Upload", ctx, mockTool.MatchedBy(func(keys *s3.PutObjectInput) bool {
		assert.Equal(fileName, *keys.FileName)
		assert.Equal(targetDatabase, *keys.DatabaseName)
		assert.Equal(targetTable, *keys.TableName)
		assert.Equal(clusterID, *keys.ClusterID)
		return true
	})).Return(uploadID, nil)

	reqBody := importOp.ImportServiceCreateImportBody{}
	err := reqBody.UnmarshalBinary([]byte(fmt.Sprintf(`{
    "clusterId": "12345",
    "dataFormat": "%s",
    "importOptions": {
      "csvFormat": {
        "backslashEscape": false,
        "delimiter": ",",
        "header": true,
        "null": "\\N",
        "separator": "\"",
		"trimLastSeparator": true
      }
    },
    "target": {
      "local": {
        "fileName": "%s",
        "targetTable": {
          "schema": "%s",
          "table": "%s"
        },
		"uploadID": "%s"
      },
      "type": "LOCAL"
    },
    "type": "LOCAL"
  }`, dataFormat, fileName, targetDatabase, targetTable, uploadID)))
	assert.Nil(err)

	body := &importModel.V1beta1Import{}
	err = json.Unmarshal([]byte(fmt.Sprintf(`{
  "allCompletedTables": [
    {
      "result": "SUCCESS",
      "tableName": "test.ttt"
    }
  ],
  "clusterId": "12345",
  "completedPercent": 100,
  "completedTables": 1,
  "createTime": "2024-04-01T06:39:50.000Z",
  "creationDetails": {
    "clusterId": "12345",
    "dataFormat": "CSV",
    "importOptions": {
      "csvFormat": {
        "backslashEscape": true,
        "delimiter": "\"",
        "header": true,
        "null": "\\N",
        "separator": ","
      }
    },
    "target": {
      "local": {
        "fileName": "a.csv",
        "targetTable": {
          "schema": "test",
          "table": "yxxxx"
        }
      },
      "type": "LOCAL"
    },
    "type": "LOCAL"
  },
  "currentTables": [],
  "dataFormat": "CSV",
  "elapsedTimeSeconds": 14,
  "id": "%s",
  "postImportCompletedPercent": 100,
  "processedSourceDataSize": "37",
  "state": "COMPLETED",
  "totalSize": "37",
  "totalTablesCount": 1
}
`, importID)), body)
	result := &importOp.ImportServiceCreateImportOK{
		Payload: body,
	}

	suite.mockClient.On("CreateImport", importOp.NewImportServiceCreateImportParams().
		WithClusterID(clusterID).WithBody(reqBody).WithContext(ctx)).
		Return(result, nil)
	suite.mockUploader.On("SetConcurrency", 5).Return(nil)
	suite.mockUploader.On("SetPartSize", int64(5*1024*1024)).Return(nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "start import success",
			args:         []string{fileName, "--cluster-id", clusterID, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable, "--separator", "\"", "--delimiter", ",", "--backslash-escape=false", "--trim-last-separator=true"},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := StartCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			tt.args = append([]string{"local"}, tt.args...)
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

func TestLocalImportSuite(t *testing.T) {
	suite.Run(t, new(LocalImportSuite))
}
