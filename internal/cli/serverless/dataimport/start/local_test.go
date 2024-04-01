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
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	importOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"

	mockTool "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const fileName = "test.csv"

type LocalImportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *LocalImportSuite) SetupTest() {
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

	importID := "12345"
	body := &importModel.V1beta1Import{
		ID: &importID,
	}
	result := &importOp.CreateImportOK{
		Payload: body,
	}

	dataFormat := "CSV"
	targetDatabase := "test"
	targetTable := "test"
	projectID := "12345"
	clusterID := "12345"
	size := "5"
	fN := fileName
	uploadUrl := "http://test.com"
	uploadRes := &importOp.GenerateUploadURLOK{
		Payload: &importModel.OpenapiGenerateUploadURLResq{
			NewFileName: &fN,
			UploadURL:   &uploadUrl,
		},
	}
	suite.mockClient.On("GenerateUploadURL", importOp.NewGenerateUploadURLParams().WithProjectID(projectID).WithClusterID(clusterID).WithBody(importOp.GenerateUploadURLBody{
		ContentLength: &size,
		FileName:      &fN,
	})).Return(uploadRes, nil)
	suite.mockClient.On("PreSignedUrlUpload", &uploadUrl, mockTool.Anything, mockTool.Anything).Return(nil)

	reqBody := importOp.CreateImportBody{}
	err := reqBody.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"type": "LOCAL",
			"data_format": "%s",
			"file_name": "%s",
			"csv_format": {
                "separator": ",",
				"delimiter": "\"",
				"header": true,
				"backslash_escape": true,
				"null": "\\N",
				"trim_last_separator": false,
				"not_null": false
			},
			"target_table": {
				"schema": "%s",
				"table": "%s"
			}}`, dataFormat, fileName, targetDatabase, targetTable)))
	assert.Nil(err)

	suite.mockClient.On("CreateImport", importOp.NewCreateImportParams().
		WithProjectID(projectID).WithClusterID(clusterID).WithBody(reqBody)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
	}{
		{
			name:         "start import success",
			args:         []string{fileName, "--project-id", projectID, "--cluster-id", clusterID, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import with unsupported data format",
			args: []string{fileName, "--project-id", projectID, "--cluster-id", clusterID, "--data-format", "yaml", "--target-database", targetDatabase, "--target-table", targetTable},
			err:  fmt.Errorf("data format yaml is not supported, please use one of [\"CSV\"]"),
		},
		{
			name:         "start import with shorthand flag",
			args:         []string{fileName, "-p", projectID, "-c", clusterID, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable},
			stdoutString: fmt.Sprintf("... Uploading file\nFile has been uploaded\n... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import without required project id",
			args: []string{fileName, "-c", clusterID, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable},
			err:  fmt.Errorf("required flag(s) \"project-id\" not set"),
		},
		{
			name: "start import without required file path",
			args: []string{"-p", projectID, "-c", clusterID, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable},
			err:  fmt.Errorf("missing argument <file-path>"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := LocalCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
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

	importID := "12345"
	body := &importModel.OpenapiCreateImportResp{
		ID: &importID,
	}
	result := &importOp.CreateImportOK{
		Payload: body,
	}

	dataFormat := "CSV"
	targetDatabase := "test"
	targetTable := "test"
	projectID := "12345"
	clusterID := "12345"
	size := "5"
	fN := fileName
	uploadUrl := "http://test.com"
	uploadRes := &importOp.GenerateUploadURLOK{
		Payload: &importModel.OpenapiGenerateUploadURLResq{
			NewFileName: &fN,
			UploadURL:   &uploadUrl,
		},
	}
	suite.mockClient.On("GenerateUploadURL", importOp.NewGenerateUploadURLParams().WithProjectID(projectID).WithClusterID(clusterID).WithBody(importOp.GenerateUploadURLBody{
		ContentLength: &size,
		FileName:      &fN,
	})).Return(uploadRes, nil)
	suite.mockClient.On("PreSignedUrlUpload", &uploadUrl, mockTool.Anything, mockTool.Anything).Return(nil)

	reqBody := importOp.CreateImportBody{}
	err := reqBody.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"type": "LOCAL",
			"data_format": "%s",
			"file_name": "%s",
			"csv_format": {
                "separator": "\"",
				"delimiter": ",",
				"header": true,
				"backslash_escape": false,
				"null": "\\N",
				"trim_last_separator": true,
				"not_null": false
			},
			"target_table": {
				"schema": "%s",
				"table": "%s"
			}}`, dataFormat, fileName, targetDatabase, targetTable)))
	assert.Nil(err)

	suite.mockClient.On("CreateImport", importOp.NewCreateImportParams().
		WithProjectID(projectID).WithClusterID(clusterID).WithBody(reqBody)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "start import success",
			args:         []string{fileName, "--project-id", projectID, "--cluster-id", clusterID, "--data-format", dataFormat, "--target-database", targetDatabase, "--target-table", targetTable, "--separator", "\"", "--delimiter", ",", "--backslash-escape=false", "--trim-last-separator=true"},
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
