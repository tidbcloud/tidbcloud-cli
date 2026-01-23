// Copyright 2026 PingCAP, Inc.
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

package export

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getExportResp = `{
  "clusterId": "fake-cluster-id",
  "completeTime": "2024-09-02T11:31:13Z",
  "createTime": "2024-09-02T11:31:01Z",
  "createdBy": "apikey-MCTGR3Jv",
  "displayName": "SNAPSHOT_2024-09-02T11:30:57Z",
  "expireTime": "2024-09-04T11:31:13Z",
  "exportId": "fake-export-id",
  "exportOptions": {
    "compression": "GZIP",
    "csvFormat": {
      "delimiter": "\"",
      "nullValue": "\\N",
      "separator": ",",
      "skipHeader": false
    },
    "database": "",
    "fileType": "CSV",
    "filter": {
      "table": {
        "patterns": [
          "test.t1"
        ],
        "where": ""
      }
    },
    "table": ""
  },
  "name": "clusters/fake-cluster-id/exports/fake-export-id",
  "snapshotTime": "2024-09-02T11:30:57.571Z",
  "state": "SUCCEEDED",
  "target": {
    "type": "LOCAL"
  },
  "updateTime": "2024-09-02T11:31:39Z"
}
`

type DescribeExportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *DescribeExportSuite) SetupTest() {
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

func (suite *DescribeExportSuite) TestDescribeExportArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	body := &export.Export{}
	err := json.Unmarshal([]byte(getExportResp), body)
	assert.Nil(err)
	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	suite.mockClient.On("GetExport", ctx, clusterId, exportId).Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "describe export success",
			args:         []string{"--cluster-id", clusterId, "--export-id", exportId},
			stdoutString: getExportResp,
		},
		{
			name:         "describe export with shorthand flag",
			args:         []string{"-c", clusterId, "-e", exportId},
			stdoutString: getExportResp,
		},
		{
			name: "describe export without required cluster id",
			args: []string{"-e", exportId},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
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

func TestDescribeExportSuite(t *testing.T) {
	suite.Run(t, new(DescribeExportSuite))
}
