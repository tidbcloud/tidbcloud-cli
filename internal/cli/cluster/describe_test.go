package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/util"

	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getClusterResultStr = `{
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
`

type DescribeClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.ApiClient
}

func (suite *DescribeClusterSuite) SetupTest() {
	var pageSize int64 = 10

	suite.mockClient = new(mock.ApiClient)
	suite.h = &internal.Helper{
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

func (suite *DescribeClusterSuite) TestDescribeClusterArgs() {
	assert := require.New(suite.T())

	body := &cluster.GetClusterOKBody{}
	err := json.Unmarshal([]byte(getClusterResultStr), body)
	assert.Nil(err)
	result := &cluster.GetClusterOK{
		Payload: body,
	}
	projectID := "12345"
	clusterID := "12345"
	suite.mockClient.On("GetCluster", cluster.NewGetClusterParams().
		WithProjectID(projectID).WithClusterID(clusterID)).
		Return(result, nil)

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "describe cluster success",
			args:         []string{"--project-id", projectID, "--cluster-id", clusterID},
			stdoutString: getClusterResultStr,
			stderrString: "",
		},
		{
			name:         "describe cluster with shorthand flag",
			args:         []string{"-p", projectID, "-c", clusterID},
			stdoutString: getClusterResultStr,
			stderrString: "",
		},
		{
			name:         "describe cluster without required project id",
			args:         []string{"-c", clusterID},
			err:          fmt.Errorf("required flag(s) \"project-id\" not set"),
			stdoutString: "",
			stderrString: "",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DescribeCmd(suite.h)
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

func TestDescribeClusterSuite(t *testing.T) {
	suite.Run(t, new(DescribeClusterSuite))
}
