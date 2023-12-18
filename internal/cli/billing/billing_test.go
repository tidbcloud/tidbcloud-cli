// Copyright 2023 PingCAP, Inc.
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

package billing

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	biApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/billing/client/billing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BillingSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *BillingSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	suite.mockClient = new(mock.TiDBCloudClient)
	suite.h = &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			return suite.mockClient, nil
		},
		IOStreams: iostream.Test(),
	}
}

func (suite *BillingSuite) TestBilling() {
	assert := require.New(suite.T())

	month := "2023-08"
	invalidMonth := "202-08"
	suite.mockClient.On("GetBillsBilledMonth", biApi.NewGetBillsBilledMonthParams().
		WithBilledMonth(month)).
		Return(&biApi.GetBillsBilledMonthOK{}, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "get bill success",
			args:         []string{"--month", month},
			stdoutString: "null\n",
		},
		{
			name: "get bill with invalid month",
			args: []string{"--month", invalidMonth},
			err:  fmt.Errorf("invalid month format: %s, should be YYYY-MM", invalidMonth),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := Cmd(suite.h)
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

func TestBillingSuite(t *testing.T) {
	suite.Run(t, new(BillingSuite))
}
