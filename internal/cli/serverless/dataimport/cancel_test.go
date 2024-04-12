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

package dataimport

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	importOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CancelImportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *CancelImportSuite) SetupTest() {
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

func (suite *CancelImportSuite) TestCancelImportArgs() {
	assert := require.New(suite.T())
	ctx := context.Background()
	result := &importOp.ImportServiceCancelImportOK{}
	clusterID := "12345"
	importID := "imp-asdasd"
	suite.mockClient.On("CancelImport", importOp.NewImportServiceCancelImportParams().
		WithClusterID(clusterID).WithID(importID).WithContext(ctx)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "cancel import with default format(json when without tty)",
			args:         []string{"--cluster-id", clusterID, "--import-id", importID, "--force"},
			stdoutString: fmt.Sprintf("Import task %s has been canceled.\n", importID),
		},
		{
			name:         "cancel import with output shorthand flag",
			args:         []string{"-c", clusterID, "--import-id", importID, "--force"},
			stdoutString: fmt.Sprintf("Import task %s has been canceled.\n", importID),
		},
		{
			name: "cancel import without required cluster id",
			args: []string{"--import-id", importID, "--force"},
			err:  fmt.Errorf("required flag(s) \"cluster-id\" not set"),
		},
		{
			name: "cancel import without force flag",
			args: []string{"-c", clusterID, "--import-id", importID},
			err:  fmt.Errorf("the terminal doesn't support prompt, please run with --force to cancel the import task"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := CancelCmd(suite.h)
			cmd.SetContext(ctx)
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

func TestCancelImportSuite(t *testing.T) {
	suite.Run(t, new(CancelImportSuite))
}
