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

package chatbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/pkg/tidbcloud/pingchat/client/operations"
	"tidbcloud-cli/pkg/tidbcloud/pingchat/models"

	mockTool "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ChatBotSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
}

func (suite *ChatBotSuite) SetupTest() {
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

func (suite *ChatBotSuite) TestChatBotArgs() {
	assert := require.New(suite.T())

	links := []*models.PingchatLink{
		{
			Link: "https://tidbcloud.com",
		},
	}
	chatOK := operations.ChatOK{Payload: &models.PingchatChatResponse{
		Content: "hello",
		Links:   links,
	}}

	res, _ := json.MarshalIndent(chatOK, "", "  ")

	suite.mockClient.On("Chat", mockTool.Anything).
		Return(&chatOK, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "chat success",
			args:         []string{"--question", "hello"},
			stdoutString: string(res) + "\n",
		},
		{
			name:         "chat succes with shorthand flag",
			args:         []string{"-q", "hello"},
			stdoutString: string(res) + "\n",
		},
		{
			name: "with unknown flag",
			args: []string{"--cluster-name", "test"},
			err:  fmt.Errorf("unknown flag: --cluster-name"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ChatBotCmd(suite.h)
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

func TestChatBotSuite(t *testing.T) {
	suite.Run(t, new(ChatBotSuite))
}
