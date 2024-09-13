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

package dataimport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"

	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

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

	clusterID := "12345"
	importID := "imp-qwert"
	t := time.Now()
	i := imp.Import{
		ClusterId:       aws.String(clusterID),
		CompletePercent: aws.Int64(100),
		CompleteTime:    *imp.NewNullableTime(&t),
		CreateTime:      &t,
		CreatedBy:       aws.String("test"),
		CreationDetails: &imp.CreationDetails{
			ImportOptions: &imp.ImportOptions{
				FileType:  "CSV",
				CsvFormat: imp.NewCSVFormat(),
			},
		},
		Id:        aws.String(importID),
		Message:   aws.String("import success"),
		Name:      aws.String("import-2024-04-01T06:39:50.000Z"),
		State:     (*imp.ImportStateEnum)(aws.String("COMPLETED")),
		TotalSize: aws.String("37"),
	}
	suite.mockClient.On("GetImport", ctx, clusterID, importID).
		Return(&i, nil)

	getImportResultStr, err := json.MarshalIndent(i, "", "  ")
	assert.Nil(err)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "describe import success",
			args:         []string{"--cluster-id", clusterID, "--import-id", importID},
			stdoutString: string(getImportResultStr) + "\n",
		},
		{
			name:         "describe import with shorthand flag",
			args:         []string{"-c", clusterID, "--import-id", importID},
			stdoutString: string(getImportResultStr) + "\n",
		},
		{
			name: "describe import without required cluster id",
			args: []string{"--import-id", importID},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DescribeCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			cmd.SetContext(ctx)
			err := cmd.Execute()
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
