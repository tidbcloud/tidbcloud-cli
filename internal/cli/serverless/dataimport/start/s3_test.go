// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	ctx := context.Background()

	clusterID := "12345"
	importID := "imp-asdasd"
	roleArn := "xasdas"
	S3Uri := "S3://xxx"
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
		ImportId:  aws.String(importID),
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
			Type: imp.IMPORTSOURCETYPEENUM_S3,
			S3: &imp.S3Source{
				AuthType: imp.IMPORTS3AUTHTYPEENUM_ROLE_ARN,
				RoleArn:  &roleArn,
				Uri:      S3Uri,
			},
		},
	}

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
			args:         []string{"--source-type", "S3", "--cluster-id", clusterID, "--file-type", string(fileType), "--s3.uri", S3Uri, "--s3.role-arn", roleArn},
			stdoutString: fmt.Sprintf("... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import with unsupported data format",
			args: []string{"--source-type", "S3", "--cluster-id", clusterID, "--file-type", "yaml"},
			err:  fmt.Errorf("file type \"yaml\" is not supported, please use one of [\"CSV\" \"SQL\" \"AURORA_SNAPSHOT\" \"PARQUET\"]"),
		},
		{
			name: "start import without required cluster id",
			args: []string{"--source-type", "S3", "--file-type", string(fileType)},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
		},
		{
			name: "start import without required uri",
			args: []string{"--source-type", "S3", "-c", clusterID, "--file-type", string(fileType), "--s3.role-arn", roleArn},
			err:  fmt.Errorf("empty S3 URI"),
		},
		{
			name: "start import without required rolearn",
			args: []string{"--source-type", "S3", "-c", clusterID, "--file-type", string(fileType), "--s3.uri", S3Uri},
			err:  fmt.Errorf("either role arn or access key id and secret access key must be provided"),
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

func (suite *S3ImportSuite) TestS3ImportArgsAccessKey() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	importID := "imp-asdasd"
	secretId := "xasdas"
	secret := "sadas"
	S3Uri := "S3://xxx"
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
		ImportId:  aws.String(importID),
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
			Type: imp.IMPORTSOURCETYPEENUM_S3,
			S3: &imp.S3Source{
				AuthType: imp.IMPORTS3AUTHTYPEENUM_ACCESS_KEY,
				AccessKey: &imp.S3SourceAccessKey{
					Id:     secretId,
					Secret: secret,
				},
				Uri: S3Uri,
			},
		},
	}

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
			args:         []string{"--source-type", "S3", "--cluster-id", clusterID, "--file-type", string(fileType), "--s3.uri", S3Uri, "--s3.access-key-id", secretId, "--s3.secret-access-key", secret},
			stdoutString: fmt.Sprintf("... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import without required uri",
			args: []string{"--source-type", "S3", "-c", clusterID, "--file-type", string(fileType), "--s3.role-arn", secretId},
			err:  fmt.Errorf("empty S3 URI"),
		},
		{
			name: "start import without required access key id",
			args: []string{"--source-type", "S3", "-c", clusterID, "--file-type", string(fileType), "--s3.uri", S3Uri, "--s3.secret-access-key", secret},
			err:  fmt.Errorf("if any flags in the group [s3.access-key-id s3.secret-access-key] are set they must all be set; missing [s3.access-key-id]"),
		},
		{
			name: "start import without required access key id",
			args: []string{"--source-type", "S3", "-c", clusterID, "--file-type", string(fileType), "--s3.uri", S3Uri, "--s3.access-key-id", secretId},
			err:  fmt.Errorf("if any flags in the group [s3.access-key-id s3.secret-access-key] are set they must all be set; missing [s3.secret-access-key]"),
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

func TestS3ImportSuite(t *testing.T) {
	suite.Run(t, new(S3ImportSuite))
}
