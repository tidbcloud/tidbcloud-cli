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
	"fmt"
	"os"
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

const getImportResultStr = `{
  "allCompletedTables": [
    {
      "result": "SUCCESS",
      "tableName": "test.ttt"
    }
  ],
  "clusterId": "12345",
  "completedPercent": 100,
  "completedTables": 1,
  "createdAt": "2024-04-01T06:39:50.000Z",
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
  "id": "imp-qwert",
  "postImportCompletedPercent": 100,
  "processedSourceDataSize": "37",
  "status": "COMPLETED",
  "totalSize": "37",
  "totalTablesCount": 1
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
	ctx := context.Background()
	body := &importModel.V1beta1Import{}
	err := json.Unmarshal([]byte(getImportResultStr), body)
	assert.Nil(err)
	result := &importOp.ImportServiceGetImportOK{
		Payload: body,
	}

	projectID := "12345"
	clusterID := "12345"
	importID := "imp-qwert"
	suite.mockClient.On("GetImport", importOp.NewImportServiceGetImportParams().
		WithClusterID(clusterID).WithID(importID).WithContext(ctx)).
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
			cmd.SetContext(ctx)
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
