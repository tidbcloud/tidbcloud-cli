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
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	importOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const listResultStr = `{
  "imports": [
    {
      "clusterId": "12345",
      "completePercent": 100,
      "completeTime": "2024-04-01T06:49:50.000Z",
      "createTime": "2024-04-01T06:39:50.000Z",
      "createdBy": "test",
      "creationDetails": {
        "importOptions": {
          "csvFormat": {
            "backslashEscape": false,
            "delimiter": ",",
            "header": true,
            "notNull": false,
            "null": "\\N",
            "separator": "\"",
            "trimLastSeparator": true
          },
          "fileType": "CSV"
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
      "message": "import success",
      "name": "import-2024-04-01T06:39:50.000Z",
      "state": "COMPLETED",
      "totalSize": "37"
    }
  ],
  "total": 1
}
`

const listResultMultiPageStr = `{
  "imports": [
    {
      "clusterId": "12345",
      "completePercent": 100,
      "completeTime": "2024-04-01T06:49:50.000Z",
      "createTime": "2024-04-01T06:39:50.000Z",
      "createdBy": "test",
      "creationDetails": {
        "importOptions": {
          "csvFormat": {
            "backslashEscape": false,
            "delimiter": ",",
            "header": true,
            "notNull": false,
            "null": "\\N",
            "separator": "\"",
            "trimLastSeparator": true
          },
          "fileType": "CSV"
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
      "message": "import success",
      "name": "import-2024-04-01T06:39:50.000Z",
      "state": "COMPLETED",
      "totalSize": "37"
    },
    {
      "clusterId": "12345",
      "completePercent": 100,
      "completeTime": "2024-04-01T06:49:50.000Z",
      "createTime": "2024-04-01T06:39:50.000Z",
      "createdBy": "test",
      "creationDetails": {
        "importOptions": {
          "csvFormat": {
            "backslashEscape": false,
            "delimiter": ",",
            "header": true,
            "notNull": false,
            "null": "\\N",
            "separator": "\"",
            "trimLastSeparator": true
          },
          "fileType": "CSV"
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
      "message": "import success",
      "name": "import-2024-04-01T06:39:50.000Z",
      "state": "COMPLETED",
      "totalSize": "37"
    }
  ],
  "total": 2
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
	ctx := context.Background()

	body := &importModel.V1beta1ListImportsResp{}
	err := json.Unmarshal([]byte(listResultStr), body)
	assert.Nil(err)
	result := &importOp.ImportServiceListImportsOK{
		Payload: body,
	}
	clusterID := "12345"
	suite.mockClient.On("ListImports", importOp.NewImportServiceListImportsParams().
		WithClusterID(clusterID).WithPage(&page).WithPageSize(&pageSize).WithContext(ctx)).
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
			args:         []string{"--cluster-id", clusterID},
			stdoutString: listResultStr,
		},
		{
			name:         "list imports with output flag",
			args:         []string{"--cluster-id", clusterID, "--output", "json"},
			stdoutString: listResultStr,
		},
		{
			name:         "list imports with output shorthand flag",
			args:         []string{"-c", clusterID, "-o", "json"},
			stdoutString: listResultStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ListCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetContext(ctx)
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
	ctx := context.Background()

	body := &importModel.V1beta1ListImportsResp{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"total": 1`, `"total": 2`)), body)
	assert.Nil(err)
	result := &importOp.ImportServiceListImportsOK{
		Payload: body,
	}
	clusterID := "12345"
	suite.mockClient.On("ListImports", importOp.NewImportServiceListImportsParams().
		WithClusterID(clusterID).WithPage(&pageOne).WithPageSize(&pageSize).WithContext(ctx)).
		Return(result, nil)
	suite.mockClient.On("ListImports", importOp.NewImportServiceListImportsParams().
		WithClusterID(clusterID).WithPage(&pageTwo).WithPageSize(&pageSize).WithContext(ctx)).
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
			args:         []string{"-c", clusterID, "--output", "json"},
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
