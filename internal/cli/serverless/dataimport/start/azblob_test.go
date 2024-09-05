// Copyright 2024 PingCAP, Inc.
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

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	imp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/import"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

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

type AzblobImportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *AzblobImportSuite) SetupTest() {
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

func (suite *AzblobImportSuite) TestAzblobImportArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "12345"
	importID := "imp-asdasd"
	sasToken := "xasdas"
	azureUri := "azure://xxx"
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
			Type: imp.IMPORTSOURCETYPEENUM_AZURE_BLOB,
			AzureBlob: &imp.AzureBlobSource{
				AuthType: imp.IMPORTAZUREBLOBAUTHTYPEENUM_SAS_TOKEN,
				SasToken: &sasToken,
				Uri:      azureUri,
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
			args:         []string{"--source-type", "AZURE_BLOB", "--cluster-id", clusterID, "--file-type", string(fileType), "--azblob.uri", azureUri, "--azblob.sas-token", sasToken},
			stdoutString: fmt.Sprintf("... Starting the import task\nImport task %s started.\n", importID),
		},
		{
			name: "start import with unsupported data format",
			args: []string{"--source-type", "AZURE_BLOB", "--cluster-id", clusterID, "--file-type", "yaml"},
			err:  fmt.Errorf("file type \"yaml\" is not supported, please use one of [\"CSV\" \"PARQUET\" \"SQL\" \"AURORA_SNAPSHOT\"]"),
		},
		{
			name: "start import without required cluster id",
			args: []string{"--source-type", "AZURE_BLOB", "--file-type", string(fileType)},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
		},
		{
			name: "start import without required uri",
			args: []string{"--source-type", "AZURE_BLOB", "-c", clusterID, "--file-type", string(fileType), "--azblob.sas-token", sasToken},
			err:  fmt.Errorf("empty Azure Blob URI"),
		},
		{
			name: "start import without required sas token",
			args: []string{"--source-type", "AZURE_BLOB", "-c", clusterID, "--file-type", string(fileType), "--azblob.uri", azureUri},
			err:  fmt.Errorf("empty Azure Blob SAS token"),
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

func TestAzblobImportSuite(t *testing.T) {
	suite.Run(t, new(AzblobImportSuite))
}
