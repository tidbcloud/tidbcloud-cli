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

package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/util"

	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const listResultStr = `{
  "items": [
    {
      "cloud_provider": "AWS",
      "cluster_type": "DEVELOPER",
      "config": {
        "components": {
          "tidb": {
            "node_quantity": 1,
            "node_size": "Shared0"
          },
          "tiflash": {
            "node_quantity": 1,
            "node_size": "Shared0",
            "storage_size_gib": 1
          },
          "tikv": {
            "node_quantity": 1,
            "node_size": "Shared0",
            "storage_size_gib": 1
          }
        },
        "port": 4000
      },
      "create_timestamp": "1668508515",
      "id": "1379661944635994072",
      "name": "sdfds",
      "project_id": "1372813089189381287",
      "region": "us-east-1",
      "status": {
        "cluster_status": "AVAILABLE",
        "connection_strings": {
          "default_user": "28cDWcUJJiewaQ7.root",
          "standard": {
            "host": "gateway01.us-east-1.prod.aws.tidbcloud.com",
            "port": 4000
          }
        },
        "node_map": {
          "tidb": [],
          "tiflash": [],
          "tikv": []
        },
        "tidb_version": "v6.3.0"
      }
    }
  ],
  "total": 1
}
`

const listResultMultiPageStr = `{
  "items": [
    {
      "cloud_provider": "AWS",
      "cluster_type": "DEVELOPER",
      "config": {
        "components": {
          "tidb": {
            "node_quantity": 1,
            "node_size": "Shared0"
          },
          "tiflash": {
            "node_quantity": 1,
            "node_size": "Shared0",
            "storage_size_gib": 1
          },
          "tikv": {
            "node_quantity": 1,
            "node_size": "Shared0",
            "storage_size_gib": 1
          }
        },
        "port": 4000
      },
      "create_timestamp": "1668508515",
      "id": "1379661944635994072",
      "name": "sdfds",
      "project_id": "1372813089189381287",
      "region": "us-east-1",
      "status": {
        "cluster_status": "AVAILABLE",
        "connection_strings": {
          "default_user": "28cDWcUJJiewaQ7.root",
          "standard": {
            "host": "gateway01.us-east-1.prod.aws.tidbcloud.com",
            "port": 4000
          }
        },
        "node_map": {
          "tidb": [],
          "tiflash": [],
          "tikv": []
        },
        "tidb_version": "v6.3.0"
      }
    },
    {
      "cloud_provider": "AWS",
      "cluster_type": "DEVELOPER",
      "config": {
        "components": {
          "tidb": {
            "node_quantity": 1,
            "node_size": "Shared0"
          },
          "tiflash": {
            "node_quantity": 1,
            "node_size": "Shared0",
            "storage_size_gib": 1
          },
          "tikv": {
            "node_quantity": 1,
            "node_size": "Shared0",
            "storage_size_gib": 1
          }
        },
        "port": 4000
      },
      "create_timestamp": "1668508515",
      "id": "1379661944635994072",
      "name": "sdfds",
      "project_id": "1372813089189381287",
      "region": "us-east-1",
      "status": {
        "cluster_status": "AVAILABLE",
        "connection_strings": {
          "default_user": "28cDWcUJJiewaQ7.root",
          "standard": {
            "host": "gateway01.us-east-1.prod.aws.tidbcloud.com",
            "port": 4000
          }
        },
        "node_map": {
          "tidb": [],
          "tiflash": [],
          "tikv": []
        },
        "tidb_version": "v6.3.0"
      }
    }
  ],
  "total": 2
}
`

type ListClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.ApiClient
}

func (suite *ListClusterSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	var pageSize int64 = 10
	suite.mockClient = new(mock.ApiClient)
	suite.h = &internal.Helper{
		Client: func() util.CloudClient {
			return suite.mockClient
		},
		QueryPageSize: pageSize,
		IOStreams:     iostream.Test(),
	}
}

func (suite *ListClusterSuite) TestListClusterArgs() {
	assert := require.New(suite.T())
	var page int64 = 1

	body := &cluster.ListClustersOfProjectOKBody{}
	err := json.Unmarshal([]byte(listResultStr), body)
	assert.Nil(err)
	result := &cluster.ListClustersOfProjectOK{
		Payload: body,
	}
	projectID := "12345"
	suite.mockClient.On("ListClustersOfProject", cluster.NewListClustersOfProjectParams().
		WithProjectID(projectID).WithPage(&page).WithPageSize(&suite.h.QueryPageSize)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "list clusters with default format(json when without tty)",
			args:         []string{projectID},
			stdoutString: listResultStr,
		},
		{
			name:         "list clusters with output flag",
			args:         []string{projectID, "--output", "json"},
			stdoutString: listResultStr,
		},
		{
			name:         "list clusters with output shorthand flag",
			args:         []string{projectID, "-o", "json"},
			stdoutString: listResultStr,
		},
		{
			name: "list clusters without required project id",
			args: []string{"-o", "json"},
			err:  fmt.Errorf("missing argument <projectID> \n\nUsage:\n  list <projectID> [flags]\n\nAliases:\n  list, ls\n\nExamples:\n  List the clusters in the project:\n  $ ticloud cluster list <projectID> \n\n  List the clusters in the project with json format:\n  $ ticloud cluster list <projectID> -o json\n\nFlags:\n  -h, --help            help for list\n  -o, --output string   Output format. One of: human, json, default: human (default \"human\")\n"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := ListCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
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

func (suite *ListClusterSuite) TestListClusterWithMultiPages() {
	assert := require.New(suite.T())
	var pageOne int64 = 1
	var pageTwo int64 = 2
	suite.h.QueryPageSize = 1

	body := &cluster.ListClustersOfProjectOKBody{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(listResultStr, `"total": 1`, `"total": 2`)), body)
	assert.Nil(err)
	result := &cluster.ListClustersOfProjectOK{
		Payload: body,
	}
	projectID := "12345"
	suite.mockClient.On("ListClustersOfProject", cluster.NewListClustersOfProjectParams().
		WithProjectID(projectID).WithPage(&pageOne).WithPageSize(&suite.h.QueryPageSize)).
		Return(result, nil)
	suite.mockClient.On("ListClustersOfProject", cluster.NewListClustersOfProjectParams().
		WithProjectID(projectID).WithPage(&pageTwo).WithPageSize(&suite.h.QueryPageSize)).
		Return(result, nil)
	cmd := ListCmd(suite.h)

	tests := []struct {
		name         string
		args         []string
		stdoutString string
		stderrString string
	}{
		{
			name:         "query with multi pages",
			args:         []string{projectID, "--output", "json"},
			stdoutString: listResultMultiPageStr,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err = cmd.Execute()
			assert.Nil(err)

			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
			suite.mockClient.AssertExpectations(suite.T())
		})
	}
}

func TestListClusterSuite(t *testing.T) {
	suite.Run(t, new(ListClusterSuite))
}
