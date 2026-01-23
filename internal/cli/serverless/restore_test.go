// Copyright 2026 PingCAP, Inc.
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

package serverless

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/mock"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/br"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RestoreClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *RestoreClusterSuite) SetupTest() {
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

func (suite *RestoreClusterSuite) TestRestoreClusterSnapshot() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "10048930788495339885"
	backupID := "289048"

	body := &br.V1beta1RestoreRequest{
		Snapshot: &br.RestoreRequestSnapshot{
			BackupId: &backupID,
		},
	}
	suite.mockClient.On("Restore", ctx, body).
		Return(&br.V1beta1RestoreResponse{ClusterId: clusterID}, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "restore success",
			args:         []string{"--backup-id", backupID},
			stdoutString: fmt.Sprintf("restore to clsuter %s, use \"ticloud serverless get -c %s\" to check the restore process\n", clusterID, clusterID),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := RestoreCmd(suite.h)
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

func (suite *RestoreClusterSuite) TestRestoreClusterPointInTime() {
	assert := require.New(suite.T())
	ctx := context.Background()

	clusterID := "10048930788495339885"
	backupTimeStr := "2023-12-15T07:00:00.000Z"

	backupTime, err := time.Parse(time.RFC3339, backupTimeStr)
	assert.Nil(err)
	body := &br.V1beta1RestoreRequest{
		PointInTime: &br.RestoreRequestPointInTime{
			BackupTime: &backupTime,
			ClusterId:  &clusterID,
		},
	}
	suite.mockClient.On("Restore", ctx, body).
		Return(&br.V1beta1RestoreResponse{ClusterId: clusterID}, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "restore success",
			args:         []string{"--cluster-id", clusterID, "--backup-time", backupTimeStr},
			stdoutString: fmt.Sprintf("restore to clsuter %s, use \"ticloud serverless get -c %s\" to check the restore process\n", clusterID, clusterID),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := RestoreCmd(suite.h)
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

func (suite *RestoreClusterSuite) TestRestoreClusterInvalidFlag() {
	assert := require.New(suite.T())

	clusterID := "10048930788495339885"
	backupTimeStr := "2023-12-15T07:00:00.000Z"
	backupID := "289048"

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name: "restore with backupId and cluster-id",
			args: []string{"--backup-id", backupID, "--cluster-id", clusterID},
			err:  errors.New("if any flags in the group [backup-id cluster-id] are set none of the others can be; [backup-id cluster-id] were all set"),
		},
		{
			name: "restore with backupId and backup-time",
			args: []string{"--backup-id", backupID, "--backup-time", backupTimeStr},
			err:  errors.New("if any flags in the group [backup-id backup-time] are set none of the others can be; [backup-id backup-time] were all set"),
		},
		{
			name: "point-in-time restore without backup-time",
			args: []string{"--cluster-id", clusterID},
			err:  errors.New("if any flags in the group [cluster-id backup-time] are set they must all be set; missing [backup-time]"),
		},
		{
			name: "point-in-time restore without cluster-id",
			args: []string{"--backup-time", backupTimeStr},
			err:  errors.New("if any flags in the group [cluster-id backup-time] are set they must all be set; missing [cluster-id]"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := RestoreCmd(suite.h)
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

func TestRestoreClusterSuite(t *testing.T) {
	suite.Run(t, new(RestoreClusterSuite))
}
