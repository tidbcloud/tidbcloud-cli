package cluster

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/util"

	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const resultStr = `{
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
}`

const resultMultiPageStr = `{
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
}`

type ListClusterSuite struct {
	suite.Suite
	h          *util.Helper
	mockClient *mock.ApiClient
}

func (suite *ListClusterSuite) SetupTest() {
	var pageSize int64 = 10

	suite.mockClient = new(mock.ApiClient)
	suite.h = &util.Helper{
		Client: func() util.CloudClient {
			return suite.mockClient
		},
		QueryPageSize: pageSize,
		IOStreams: &iostream.IOStreams{
			Out: &bytes.Buffer{},
			Err: &bytes.Buffer{},
		},
	}
}

func (suite *ListClusterSuite) TestListClusterArgs() {
	assert := require.New(suite.T())
	var page int64 = 1

	body := &cluster.ListClustersOfProjectOKBody{}
	err := json.Unmarshal([]byte(resultStr), body)
	assert.Nil(err)
	result := &cluster.ListClustersOfProjectOK{
		Payload: body,
	}
	projectID := "12345"
	suite.mockClient.On("ListClustersOfProject", cluster.NewListClustersOfProjectParams().
		WithProjectID(projectID).WithPage(&page).WithPageSize(&suite.h.QueryPageSize)).
		Return(result, nil)
	cmd := ListCmd(suite.h)

	tests := []struct {
		name         string
		args         []string
		stdoutString string
		stderrString string
	}{
		{
			name:         "print json with output arg",
			args:         []string{projectID, "--output", "json"},
			stdoutString: resultStr,
			stderrString: "",
		},
		{
			name:         "print json with output shorthand arg",
			args:         []string{projectID, "-o", "json"},
			stdoutString: resultStr,
			stderrString: "",
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

func (suite *ListClusterSuite) TestListClusterWithMultiPages() {
	assert := require.New(suite.T())
	var pageOne int64 = 1
	var pageTwo int64 = 2
	suite.h.QueryPageSize = 1

	body := &cluster.ListClustersOfProjectOKBody{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(resultStr, `"total": 1`, `"total": 2`)), body)
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
			stdoutString: resultMultiPageStr,
			stderrString: "",
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
