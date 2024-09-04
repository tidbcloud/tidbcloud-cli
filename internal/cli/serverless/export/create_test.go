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
	"errors"
	"fmt"
	"os"
	"testing"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
)

type CreateExportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *CreateExportSuite) SetupTest() {
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

func (suite *CreateExportSuite) TestCreateExportToLocal() {
	ctx := context.Background()

	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"

	suite.mockClient.On("CreateExport", ctx, clusterId, getDefaultCreateExportBody()).
		Return(&export.Export{
			ExportId: &exportId,
		}, nil)

	tests := []Test{
		{
			name:         "export all data to local with force",
			args:         []string{"-c", clusterId, "--force"},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
		{
			name: "export all data to local without force",
			args: []string{"-c", clusterId},
			err:  errors.New("the terminal doesn't support prompt, please run with --force to create export"),
		},
	}
	suite.AssertTest(ctx, tests)
}

func (suite *CreateExportSuite) TestCreateExportToS3WithRoleArn() {
	ctx := context.Background()

	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	targetType := export.EXPORTTARGETTYPEENUM_S3
	uri := "s3://fake-bucket/fake-prefix"
	roleArn := "arn:aws:iam::123456789012:role/service-role/AmazonS3FullAccess"

	body := getDefaultCreateExportBody()
	body.Target = &export.ExportTarget{
		Type: &targetType,
		S3: &export.S3Target{
			Uri:      &uri,
			AuthType: export.EXPORTS3AUTHTYPEENUM_ROLE_ARN,
			RoleArn:  &roleArn,
		},
	}
	suite.mockClient.On("CreateExport", ctx, clusterId, body).
		Return(&export.Export{
			ExportId: &exportId,
		}, nil)

	tests := []Test{
		{
			name:         "export all data to s3 using role arn",
			args:         []string{"-c", clusterId, "--target-type", "S3", "--s3.uri", uri, "--s3.role-arn", roleArn, "--force"},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
		{
			name: "export all data to s3 without uri",
			args: []string{"-c", clusterId, "--target-type", "S3", "--s3.role-arn", roleArn, "--force"},
			err:  errors.New("S3 URI is required when target type is S3"),
		},
		{
			name: "export all data to s3 without auth",
			args: []string{"-c", clusterId, "--target-type", "S3", "--s3.uri", uri, "--force"},
			err:  errors.New("missing S3 auth information, require either role arn or access key id and secret access key"),
		},
	}

	suite.AssertTest(ctx, tests)
}

func (suite *CreateExportSuite) TestCreateExportToS3WithAccessKey() {
	ctx := context.Background()

	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	targetType := export.EXPORTTARGETTYPEENUM_S3
	uri := "s3://fake-bucket/fake-prefix"
	accessKeyId := "fake-id"
	secretAccess := "fake-secret"

	body := getDefaultCreateExportBody()
	body.Target = &export.ExportTarget{
		Type: &targetType,
		S3: &export.S3Target{
			Uri:      &uri,
			AuthType: export.EXPORTS3AUTHTYPEENUM_ACCESS_KEY,
			AccessKey: &export.S3TargetAccessKey{
				Id:     accessKeyId,
				Secret: secretAccess,
			},
		},
	}
	suite.mockClient.On("CreateExport", ctx, clusterId, body).
		Return(&export.Export{
			ExportId: &exportId,
		}, nil)

	tests := []Test{
		{
			name:         "export all data to s3 using access key",
			args:         []string{"-c", clusterId, "--target-type", "S3", "--s3.uri", uri, "--s3.access-key-id", accessKeyId, "--s3.secret-access-key", secretAccess, "--force"},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
	}
	suite.AssertTest(ctx, tests)
}

func (suite *CreateExportSuite) TestCreateExportToGCS() {
	ctx := context.Background()

	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	targetType := export.EXPORTTARGETTYPEENUM_GCS
	uri := "s3://fake-bucket/fake-prefix"
	serviceAccountKey := "fake-service-account-key"

	body := getDefaultCreateExportBody()
	body.Target = &export.ExportTarget{
		Type: &targetType,
		Gcs: &export.GCSTarget{
			Uri:               uri,
			AuthType:          export.EXPORTGCSAUTHTYPEENUM_SERVICE_ACCOUNT_KEY,
			ServiceAccountKey: &serviceAccountKey,
		},
	}
	suite.mockClient.On("CreateExport", ctx, clusterId, body).
		Return(&export.Export{
			ExportId: &exportId,
		}, nil)

	tests := []Test{
		{
			name:         "export all data to gcs",
			args:         []string{"-c", clusterId, "--target-type", "GCS", "--gcs.uri", uri, "--gcs.service-account-key", serviceAccountKey, "--force"},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
		{
			name: "export all data to gcs without auth",
			args: []string{"-c", clusterId, "--target-type", "GCS", "--gcs.uri", uri, "--force"},
			err:  errors.New("GCS service account key is required when target type is GCS"),
		},
		{
			name: "export all data to gcs without uri",
			args: []string{"-c", clusterId, "--target-type", "GCS", "--gcs.service-account-key", serviceAccountKey, "--force"},
			err:  errors.New("GCS URI is required when target type is GCS"),
		},
	}
	suite.AssertTest(ctx, tests)
}

func (suite *CreateExportSuite) TestCreateExportToAzure() {
	ctx := context.Background()

	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	targetType := export.EXPORTTARGETTYPEENUM_AZURE_BLOB
	uri := "s3://fake-bucket/fake-prefix"
	sasToken := "fake-sas-token"

	body := getDefaultCreateExportBody()
	body.Target = &export.ExportTarget{
		Type: &targetType,
		AzureBlob: &export.AzureBlobTarget{
			Uri:      uri,
			AuthType: export.EXPORTAZUREBLOBAUTHTYPEENUM_SAS_TOKEN,
			SasToken: &sasToken,
		},
	}
	suite.mockClient.On("CreateExport", ctx, clusterId, body).
		Return(&export.Export{
			ExportId: &exportId,
		}, nil)

	tests := []Test{
		{
			name:         "export all data to azure blob",
			args:         []string{"-c", clusterId, "--target-type", "AZURE_BLOB", "--azblob.uri", uri, "--azblob.sas-token", sasToken, "--force"},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
		{
			name: "export all data to azure blob without auth",
			args: []string{"-c", clusterId, "--target-type", "AZURE_BLOB", "--azblob.uri", uri, "--force"},
			err:  errors.New("Azure Blob SAS token is required when target type is AZURE_BLOB"),
		},
		{
			name: "export all data to azure blob without uri",
			args: []string{"-c", clusterId, "--target-type", "AZURE_BLOB", "--azblob.sas-token", sasToken, "--force"},
			err:  errors.New("Azure Blob URI is required when target type is AZURE_BLOB"),
		},
	}
	suite.AssertTest(ctx, tests)
}

func (suite *CreateExportSuite) TestCreateExportWithSQLFile() {
	ctx := context.Background()

	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	fileType := export.EXPORTFILETYPEENUM_SQL

	body := getDefaultCreateExportBody()
	body.ExportOptions.FileType = &fileType
	body.ExportOptions.CsvFormat = nil

	suite.mockClient.On("CreateExport", ctx, clusterId, body).
		Return(&export.Export{
			ExportId: &exportId,
		}, nil)

	tests := []Test{
		{
			name:         "export all data to sql file",
			args:         []string{"-c", clusterId, "--file-type", "SQL", "--force"},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
	}
	suite.AssertTest(ctx, tests)
}

func (suite *CreateExportSuite) TestCreateExportWithParquetFile() {
	ctx := context.Background()

	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	fileType := export.EXPORTFILETYPEENUM_PARQUET
	parquetCompression := export.EXPORTPARQUETCOMPRESSIONTYPEENUM_ZSTD

	body := getDefaultCreateExportBody()
	body.ExportOptions.FileType = &fileType
	body.ExportOptions.Compression = nil
	body.ExportOptions.CsvFormat = nil
	body.ExportOptions.ParquetFormat = &export.ExportOptionsParquetFormat{
		Compression: &parquetCompression,
	}

	suite.mockClient.On("CreateExport", ctx, clusterId, body).
		Return(&export.Export{
			ExportId: &exportId,
		}, nil)

	tests := []Test{
		{
			name:         "export all data to parquet file",
			args:         []string{"-c", clusterId, "--file-type", "PARQUET", "--force"},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
		{
			name: "export all data to parquet file with compression",
			args: []string{"-c", clusterId, "--file-type", "PARQUET", "--compression", "GZIP", "--force"},
			err:  errors.New("--compression is not supported when file type is parquet, please use --parquet.compression instead"),
		},
	}
	suite.AssertTest(ctx, tests)
}

func (suite *CreateExportSuite) TestCreateExportWithSQLFilter() {
	ctx := context.Background()

	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	sql := "SELECT * FROM db.table WHERE column = 'value'"

	body := getDefaultCreateExportBody()
	body.ExportOptions.Filter = &export.ExportOptionsFilter{
		Sql: &sql,
	}

	suite.mockClient.On("CreateExport", ctx, clusterId, body).
		Return(&export.Export{
			ExportId: &exportId,
		}, nil)

	tests := []Test{
		{
			name:         "export data with sql filter",
			args:         []string{"-c", clusterId, "--sql", sql},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
	}
	suite.AssertTest(ctx, tests)
}

func (suite *CreateExportSuite) TestCreateExportWithTableFilter() {
	ctx := context.Background()

	clusterId := "fake-cluster-id"
	exportId := "fake-export-id"
	where := "column = 'value'"
	pattern1 := "db.t\\.able"
	pattern2 := "db.table"

	body := getDefaultCreateExportBody()
	body.ExportOptions.Filter = &export.ExportOptionsFilter{
		Table: &export.ExportOptionsFilterTable{
			Patterns: []string{pattern1, pattern2},
			Where:    &where,
		},
	}

	suite.mockClient.On("CreateExport", ctx, clusterId, body).
		Return(&export.Export{
			ExportId: &exportId,
		}, nil)

	tests := []Test{
		{
			name:         "export data with table filter",
			args:         []string{"-c", clusterId, "--where", where, "--filter", pattern1, "--filter", pattern2},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
		{
			name:         "export data with table filter2",
			args:         []string{"-c", clusterId, "--where", where, "--filter", fmt.Sprintf("%s,%s", pattern1, pattern2)},
			stdoutString: fmt.Sprintf("export %s is running now\n", exportId),
		},
	}
	suite.AssertTest(ctx, tests)
}

func getDefaultCreateExportBody() *export.ExportServiceCreateExportBody {
	defaultFileType := export.EXPORTFILETYPEENUM_CSV
	defaultTargetType := export.EXPORTTARGETTYPEENUM_LOCAL
	defaultCompression := export.EXPORTCOMPRESSIONTYPEENUM_GZIP
	separatorDefaultValue := ","
	delimiterDefaultValue := "\""
	nullValueDefaultValue := "\\N"
	skipHeaderDefaultValue := false
	return &export.ExportServiceCreateExportBody{
		ExportOptions: &export.ExportOptions{
			FileType:    &defaultFileType,
			Compression: &defaultCompression,
			CsvFormat: &export.ExportOptionsCSVFormat{
				Separator:  &separatorDefaultValue,
				Delimiter:  *export.NewNullableString(&delimiterDefaultValue),
				NullValue:  *export.NewNullableString(&nullValueDefaultValue),
				SkipHeader: &skipHeaderDefaultValue,
			},
		},
		Target: &export.ExportTarget{
			Type: &defaultTargetType,
		},
	}
}

func TestCreateExportSuite(t *testing.T) {
	suite.Run(t, new(CreateExportSuite))
}

func (suite *CreateExportSuite) AssertTest(ctx context.Context, tests []Test) {
	assert := require.New(suite.T())
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := CreateCmd(suite.h)
			cmd.SetContext(ctx)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if err != nil {
				assert.Equal(tt.err.Error(), err.Error())
			}

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
			}
		})
	}
}

type Test struct {
	name         string
	args         []string
	err          error
	stdoutString string
	stderrString string
}
