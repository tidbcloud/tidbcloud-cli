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
	"fmt"
	"os"
	"testing"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/aws/s3"
	"tidbcloud-cli/internal/service/cloud"

	imp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/import"

	"github.com/aws/aws-sdk-go-v2/aws"
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

	clusterID := "12345"
	uploadID := "upl-sadads"
	importID := "imp-asdasd"
	targetDatabase := "test"
	targetTable := "test"
	t := time.Now()
	fileType := imp.IMPORTFILETYPEENUM_CSV
	csvFormat := &imp.CSVFormat{
		BackslashEscape:   *imp.NewNullableBool(aws.Bool(true)),
		Delimiter:         *imp.NewNullableString(aws.String("\"")),
		Header:            *imp.NewNullableBool(aws.Bool(true)),
		Null:              *imp.NewNullableString(aws.String("\\N")),
		NotNull:           *imp.NewNullableBool(aws.Bool(false)),
		Separator:         aws.String(","),
		TrimLastSeparator: *imp.NewNullableBool(aws.Bool(false)),
	}
	i := imp.Import{
		ClusterId:       aws.String(clusterID),
		CompletePercent: aws.Int64(100),
		CompleteTime:    *imp.NewNullableTime(&t),
		CreateTime:      &t,
		CreatedBy:       aws.String("test"),
		CreationDetails: &imp.CreationDetails{
			ImportOptions: &imp.ImportOptions{
				FileType:  fileType,
				CsvFormat: csvFormat,
			},
		},
		Id:        aws.String(importID),
		Message:   aws.String("import success"),
		Name:      aws.String("import-2024-04-01T06:39:50.000Z"),
		State:     (*imp.ImportStateEnum)(aws.String("COMPLETED")),
		TotalSize: aws.String("37"),
	}
	body := &imp.ImportServiceCreateImportBody{
		ImportOptions: imp.ImportOptions{
			FileType:  fileType,
			CsvFormat: csvFormat,
		},
		Source: imp.ImportSource{
			Type: "LOCAL",
			Local: &imp.LocalSource{
				UploadId:       uploadID,
				TargetDatabase: targetDatabase,
				TargetTable:    targetTable,
			},
		},
	}

	suite.mockUploader.On("Upload", ctx, mockTool.MatchedBy(func(keys *s3.PutObjectInput) bool {
		assert.Equal(fileName, *keys.FileName)
		assert.Equal(targetDatabase, *keys.DatabaseName)
		assert.Equal(targetTable, *keys.TableName)
		assert.Equal(clusterID, keys.ClusterID)
		return true
	})).Return(uploadID, nil)
	suite.mockUploader.On("SetConcurrency", 5).Return(nil)
	suite.mockUploader.On("SetPartSize", int64(5*1024*1024)).Return(nil)

	suite.mockClient.On("CreateImport", ctx, clusterID, body).
		Return(&i, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
	}{
		{
			name:         "start import success",
			args:         []string{"--source-type", "LOCAL", "--local.file-path", fileName, "--cluster-id", clusterID, "--file-type", string(fileType), "--local.target-database", targetDatabase, "--local.target-table", targetTable},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import with unsupported data format",
			args: []string{"--source-type", "LOCAL", "--local.file-path", fileName, "--cluster-id", clusterID, "--file-type", "yaml", "--local.target-database", targetDatabase, "--local.target-table", targetTable},
			err:  fmt.Errorf("file type \"yaml\" is not supported, please use one of [\"CSV\"]"),
		},
		{
			name:         "start import with shorthand flag",
			args:         []string{"--source-type", "LOCAL", "--local.file-path", fileName, "-c", clusterID, "--file-type", string(fileType), "--local.target-database", targetDatabase, "--local.target-table", targetTable},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import without required cluster id",
			args: []string{"--source-type", "LOCAL", "--local.file-path", fileName, "--file-type", string(fileType), "--local.target-database", targetDatabase, "--local.target-table", targetTable},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
		},
		{
			name: "start import without required file path",
			args: []string{"--source-type", "LOCAL", "-c", clusterID, "--file-type", string(fileType), "--local.target-database", targetDatabase, "--local.target-table", targetTable},
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
			err := cmd.Execute()
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

	clusterID := "12345"
	uploadID := "upl-sadads"
	importID := "imp-asdasd"
	targetDatabase := "test"
	targetTable := "test"
	t := time.Now()
	fileType := imp.IMPORTFILETYPEENUM_CSV
	csvFormat := &imp.CSVFormat{
		BackslashEscape:   *imp.NewNullableBool(aws.Bool(false)),
		Delimiter:         *imp.NewNullableString(aws.String("i")),
		Header:            *imp.NewNullableBool(aws.Bool(false)),
		Null:              *imp.NewNullableString(aws.String("null")),
		NotNull:           *imp.NewNullableBool(aws.Bool(false)),
		Separator:         aws.String("sep"),
		TrimLastSeparator: *imp.NewNullableBool(aws.Bool(true)),
	}
	i := imp.Import{
		ClusterId:       aws.String(clusterID),
		CompletePercent: aws.Int64(100),
		CompleteTime:    *imp.NewNullableTime(&t),
		CreateTime:      &t,
		CreatedBy:       aws.String("test"),
		CreationDetails: &imp.CreationDetails{
			ImportOptions: &imp.ImportOptions{
				FileType:  fileType,
				CsvFormat: csvFormat,
			},
		},
		Id:        aws.String(importID),
		Message:   aws.String("import success"),
		Name:      aws.String("import-2024-04-01T06:39:50.000Z"),
		State:     (*imp.ImportStateEnum)(aws.String("COMPLETED")),
		TotalSize: aws.String("37"),
	}
	body := &imp.ImportServiceCreateImportBody{
		ImportOptions: imp.ImportOptions{
			FileType:  fileType,
			CsvFormat: csvFormat,
		},
		Source: imp.ImportSource{
			Type: "LOCAL",
			Local: &imp.LocalSource{
				UploadId:       uploadID,
				TargetDatabase: targetDatabase,
				TargetTable:    targetTable,
			},
		},
	}

	suite.mockUploader.On("Upload", ctx, mockTool.MatchedBy(func(keys *s3.PutObjectInput) bool {
		assert.Equal(fileName, *keys.FileName)
		assert.Equal(targetDatabase, *keys.DatabaseName)
		assert.Equal(targetTable, *keys.TableName)
		assert.Equal(clusterID, keys.ClusterID)
		return true
	})).Return(uploadID, nil)
	suite.mockUploader.On("SetConcurrency", 5).Return(nil)
	suite.mockUploader.On("SetPartSize", int64(5*1024*1024)).Return(nil)

	suite.mockClient.On("CreateImport", ctx, clusterID, body).
		Return(&i, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name: "start import success",
			args: []string{"--source-type", "LOCAL", "--local.file-path", fileName, "--cluster-id", clusterID, "--file-type", string(fileType), "--local.target-database", targetDatabase, "--local.target-table", targetTable,
				"--csv.separator", "sep", "--csv.delimiter", "i", "--csv.backslash-escape=false", "--csv.trim-last-separator",
				"--csv.skip-header=true", "--csv.null-value", `null`},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := StartCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
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

func TestLocalImportSuite(t *testing.T) {
	suite.Run(t, new(LocalImportSuite))
}
