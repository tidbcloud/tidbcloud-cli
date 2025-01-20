// Copyright 2025 PingCAP, Inc.
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
	"testing"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"

	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"

	"github.com/aws/aws-sdk-go-v2/aws"
	mockTool "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.einride.tech/aip/pagination"
)

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
	var pageSize = int32(suite.h.QueryPageSize)
	ctx := context.Background()
	t := time.Now()
	orderBy := "create_time desc"
	clusterID := "12345"

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
		Id:        aws.String("imp-asdasd"),
		Message:   aws.String("import success"),
		Name:      aws.String("import-2024-04-01T06:39:50.000Z"),
		State:     (*imp.ImportStateEnum)(aws.String("COMPLETED")),
		TotalSize: aws.String("37"),
	}

	resp := &imp.ListImportsResp{
		Imports:   []imp.Import{i},
		TotalSize: aws.Int64(1),
	}

	suite.mockClient.On("ListImports", ctx, clusterID, &pageSize, mockTool.Anything, &orderBy).
		Return(resp, nil)
	j, err := json.MarshalIndent(resp, "", "  ")
	assert.Nil(err)
	listResultStr := string(j) + "\n"

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

func (suite *ListImportSuite) TestListImportWithMultiPages() {
	assert := require.New(suite.T())
	suite.h.QueryPageSize = 1
	var pageSize = int32(suite.h.QueryPageSize)
	ctx := context.Background()
	clusterID := "12345"
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
		Id:        aws.String("imp-asdasd"),
		Message:   aws.String("import success"),
		Name:      aws.String("import-2024-04-01T06:39:50.000Z"),
		State:     (*imp.ImportStateEnum)(aws.String("COMPLETED")),
		TotalSize: aws.String("37"),
	}
	nextPageToken := pagination.PageToken{
		Offset: 1,
	}
	resp := &imp.ListImportsResp{
		Imports:       []imp.Import{i},
		TotalSize:     aws.Int64(2),
		NextPageToken: aws.String(nextPageToken.String()),
	}

	orderBy := "create_time desc"
	suite.mockClient.On("ListImports", ctx, clusterID, &pageSize, mockTool.MatchedBy(func(pageToken *string) bool {
		return pageToken == nil
	}), &orderBy).
		Return(resp, nil)

	resp2 := &imp.ListImportsResp{
		Imports:       []imp.Import{i},
		TotalSize:     aws.Int64(2),
		NextPageToken: nil,
	}
	suite.mockClient.On("ListImports", ctx, clusterID, &pageSize, aws.String(nextPageToken.String()), &orderBy).
		Return(resp2, nil)

	resp3 := &imp.ListImportsResp{
		Imports:   []imp.Import{i, i},
		TotalSize: aws.Int64(2),
	}
	listResultMultiPageStr, err := json.MarshalIndent(resp3, "", "  ")
	assert.Nil(err)

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
			stdoutString: string(listResultMultiPageStr) + "\n",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
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
