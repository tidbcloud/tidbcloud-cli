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

package export

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CancelExportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *CancelExportSuite) SetupTest() {
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

func (suite *DescribeExportSuite) TestCancelExportArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	body := &export.Export{}
	err := json.Unmarshal([]byte(getExportResp), body)
	assert.Nil(err)
	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	suite.mockClient.On("CancelExport", ctx, clusterId, exportId).Return(body, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "cancel export success",
			args:         []string{"-c", clusterId, "-e", exportId, "--force"},
			stdoutString: "export fake-export-id canceled\n",
		},
		{
			name: "cancel export without force",
			args: []string{"-c", clusterId, "-e", exportId},
			err:  fmt.Errorf("the terminal doesn't support prompt, please run with --force to cancel the export"),
		},
		{
			name: "cancel export without required cluster id",
			args: []string{"-e", exportId},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := CancelCmd(suite.h)
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

func TestCancelExportSuite(t *testing.T) {
	suite.Run(t, new(CancelExportSuite))
}
