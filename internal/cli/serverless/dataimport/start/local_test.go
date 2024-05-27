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
  "clusterId": "12345",
  "completePercent": 100,
  "createTime": "2024-04-01T06:39:50.000Z",
  "completeTime": "2024-04-01T06:49:50.000Z",
  "creationDetails": {
    "importOptions": {
      "fileType": "CSV",
      "csvFormat": {
        "backslashEscape": true,
        "delimiter": "\"",
        "header": true,
        "null": "\\N",
        "separator": ","
      }
    },
    "source": {
      "local": {
        "fileName": "a.csv",
       	"targetDatabase": "test",
        "targetTable": "test"
      },
      "type": "LOCAL"
    }
  },
  "id": "%s",
  "name": "import-2024-04-01T06:39:50.000Z",
  "state": "COMPLETED",
  "totalSize": "37",
  "createdBy": "test",
  "message": "import success"
}
`, importID)), body)
	assert.Nil(err)
	result := &importOp.ImportServiceCreateImportOK{
		Payload: body,
	}

	fileType := "CSV"
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
    "importOptions": {
      "fileType": "%s",
      "csvFormat": {
        "backslashEscape": true,
        "delimiter": "\"",
        "header": true,
		"notNull": false,
        "null": "\\N",
        "separator": ",",
		"trimLastSeparator": false
      }
    },
    "source": {
      "local": {
       	"targetDatabase": "%s",
        "targetTable": "%s",
		"uploadID": "%s"
      },
      "type": "LOCAL"
    }
  }`, fileType, targetDatabase, targetTable, uploadID)))
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
			args:         []string{"--source-type", "LOCAL", "--local.file-path", fileName, "--cluster-id", clusterID, "--file-type", fileType, "--local.target-database", targetDatabase, "--local.target-table", targetTable},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import with unsupported data format",
			args: []string{"--source-type", "LOCAL", "--local.file-path", fileName, "--cluster-id", clusterID, "--file-type", "yaml", "--local.target-database", targetDatabase, "--local.target-table", targetTable},
			err:  fmt.Errorf("file type \"yaml\" is not supported, please use one of [\"CSV\"]"),
		},
		{
			name:         "start import with shorthand flag",
			args:         []string{"--source-type", "LOCAL", "--local.file-path", fileName, "-c", clusterID, "--file-type", fileType, "--local.target-database", targetDatabase, "--local.target-table", targetTable},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import without required cluster id",
			args: []string{"--source-type", "LOCAL", "--local.file-path", fileName, "--file-type", fileType, "--local.target-database", targetDatabase, "--local.target-table", targetTable},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
		},
		{
			name: "start import without required file path",
			args: []string{"--source-type", "LOCAL", "-c", clusterID, "--file-type", fileType, "--local.target-database", targetDatabase, "--local.target-table", targetTable},
			err:  fmt.Errorf("required flag(s) \"local.file-path\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := StartCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			cmd.SetContext(ctx)
			err = cmd.Execute()
			if err != nil {
				assert.NotNil(tt.err)
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

	fileType := "CSV"
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
    "importOptions": {
      "fileType": "%s",
      "csvFormat": {
        "backslashEscape": false,
        "delimiter": ",",
        "header": true,
		"notNull": false,
        "null": "\\N",
        "separator": "\"",
		"trimLastSeparator": true
      }
    },
    "source": {
      "local": {
       	"targetDatabase": "%s",
        "targetTable": "%s",
		"uploadID": "%s"
      },
      "type": "LOCAL"
    }
  }`, fileType, targetDatabase, targetTable, uploadID)))
	assert.Nil(err)

	body := &importModel.V1beta1Import{}
	err = json.Unmarshal([]byte(fmt.Sprintf(`{
  "clusterId": "12345",
  "completePercent": 100,
  "createTime": "2024-04-01T06:39:50.000Z",
  "completeTime": "2024-04-01T06:49:50.000Z",
  "creationDetails": {
    "importOptions": {
      "fileType": "CSV",
      "csvFormat": {
        "backslashEscape": false,
        "delimiter": ",",
        "header": true,
        "null": "\\N",
        "separator": "\"",
		"trimLastSeparator": true
      }
    },
    "source": {
      "local": {
        "fileName": "a.csv",
       	"targetDatabase": "test",
        "targetTable": "test"
      },
      "type": "LOCAL"
    }
  },
  "id": "%s",
  "name": "import-2024-04-01T06:39:50.000Z",
  "state": "COMPLETED",
  "totalSize": "37",
  "createdBy": "test",
  "message": "import success"
}
`, importID)), body)
	assert.Nil(err)
	result := &importOp.ImportServiceCreateImportOK{
		Payload: body,
	}

	suite.mockClient.On("CreateImport", importOp.NewImportServiceCreateImportParams().
		WithClusterID(clusterID).WithBody(reqBody).WithContext(ctx)).
		Return(result, nil)
	suite.mockUploader.On("SetConcurrency", 5).Return(nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "start import success",
			args:         []string{"--source-type", "LOCAL", "--local.file-path", fileName, "--cluster-id", clusterID, "--file-type", fileType, "--local.target-database", targetDatabase, "--local.target-table", targetTable, "--csv.separator", "\"", "--csv.delimiter", ",", "--csv.backslash-escape=false", "--csv.trim-last-separator=true"},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := StartCmd(suite.h)
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

func TestLocalImportSuite(t *testing.T) {
	suite.Run(t, new(LocalImportSuite))
}
