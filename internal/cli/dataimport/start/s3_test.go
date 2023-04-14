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

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type S3ImportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *S3ImportSuite) SetupTest() {
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

func (suite *S3ImportSuite) TestS3ImportArgs() {
	assert := require.New(suite.T())

	importID := "12345"
	body := &importModel.OpenapiCreateImportResp{
		ID: &importID,
	}
	result := &importOp.CreateImportOK{
		Payload: body,
	}

	awsRoleArn := "aws"
	dataFormat := "CSV"
	sourceUrl := "s3://test"
	reqBody := importOp.CreateImportBody{}
	err := reqBody.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"aws_role_arn": "%s",
			"data_format": "%s",
			"source_url": "%s",
			"type": "S3",
			"csv_format": {
                "separator": ",",
				"delimiter": "\"",
				"header": true,
				"backslash_escape": true,
				"null": "\\N",
				"trim_last_separator": false,
				"not_null": false
			}
			}`, awsRoleArn, dataFormat, sourceUrl)))
	assert.Nil(err)

	projectID := "12345"
	clusterID := "12345"
	suite.mockClient.On("CreateImportTask", importOp.NewCreateImportParams().
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
			args:         []string{"--project-id", projectID, "--cluster-id", clusterID, "--aws-role-arn", awsRoleArn, "--data-format", dataFormat, "--source-url", sourceUrl},
			stdoutString: fmt.Sprintf("... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import with unsupported data format",
			args: []string{"--project-id", projectID, "--cluster-id", clusterID, "--aws-role-arn", awsRoleArn, "--data-format", "yaml", "--source-url", sourceUrl},
			err:  fmt.Errorf("data format yaml is not supported, please use one of [\"CSV\" \"SqlFile\" \"Parquet\" \"AuroraSnapshot\"]"),
		},
		{
			name:         "start import with shorthand flag",
			args:         []string{"-p", projectID, "-c", clusterID, "--aws-role-arn", awsRoleArn, "--data-format", dataFormat, "--source-url", sourceUrl},
			stdoutString: fmt.Sprintf("... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import without required project id",
			args: []string{"-c", clusterID, "--aws-role-arn", awsRoleArn, "--data-format", dataFormat, "--source-url", sourceUrl},
			err:  fmt.Errorf("required flag(s) \"project-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := StartCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			tt.args = append([]string{"s3"}, tt.args...)
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

func (suite *S3ImportSuite) TestS3ImportCSVFormat() {
	assert := require.New(suite.T())

	importID := "12345"
	body := &importModel.OpenapiCreateImportResp{
		ID: &importID,
	}
	result := &importOp.CreateImportOK{
		Payload: body,
	}

	awsRoleArn := "aws"
	dataFormat := "CSV"
	sourceUrl := "s3://test"
	reqBody := importOp.CreateImportBody{}
	err := reqBody.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"aws_role_arn": "%s",
			"data_format": "%s",
			"source_url": "%s",
			"type": "S3",
			"csv_format": {
                "separator": "\"",
				"delimiter": ",",
				"header": true,
				"backslash_escape": false,
				"null": "\\N",
				"trim_last_separator": true,
				"not_null": false
			}
			}`, awsRoleArn, dataFormat, sourceUrl)))
	assert.Nil(err)

	projectID := "12345"
	clusterID := "12345"
	suite.mockClient.On("CreateImportTask", importOp.NewCreateImportParams().
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
			name:         "start import with custom format",
			args:         []string{"--project-id", projectID, "--cluster-id", clusterID, "--aws-role-arn", awsRoleArn, "--data-format", dataFormat, "--source-url", sourceUrl, "--separator", "\"", "--delimiter", ",", "--backslash-escape=false", "--trim-last-separator=true"},
			stdoutString: fmt.Sprintf("... Starting the import task\nImport task %s started.\n", importID),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := StartCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			tt.args = append([]string{"s3"}, tt.args...)
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

func TestS3ImportSuite(t *testing.T) {
	suite.Run(t, new(S3ImportSuite))
}
