package cluster

import (
	"bytes"
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

type DeleteClusterSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.ApiClient
}

func (suite *DeleteClusterSuite) SetupTest() {
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

func (suite *DeleteClusterSuite) TestDeleteClusterArgs() {
	assert := require.New(suite.T())

	projectID := "12345"
	clusterID := "12345"
	suite.mockClient.On("DeleteCluster", cluster.NewDeleteClusterParams().
		WithProjectID(projectID).WithClusterID(clusterID)).
		Return(&cluster.DeleteClusterOK{}, nil)
	suite.mockClient.On("GetCluster", cluster.NewGetClusterParams().
		WithProjectID(projectID).WithClusterID(clusterID)).
		Return(nil, &cluster.GetClusterNotFound{})

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "delete cluster success",
			args:         []string{"--project-id", projectID, "--cluster-id", clusterID},
			stdoutString: "cluster deleted",
			stderrString: "",
		},
		{
			name:         "delete cluster with output flag",
			args:         []string{"-p", projectID, "-c", clusterID},
			stdoutString: "cluster deleted",
			stderrString: "",
		},
		{
			name:         "delete cluster without required project id",
			args:         []string{"-c", clusterID},
			err:          fmt.Errorf("required flag(s) \"project-id\" not set"),
			stdoutString: "",
			stderrString: "",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := DeleteCmd(suite.h)
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

func TestDeleteClusterSuite(t *testing.T) {
	suite.Run(t, new(DeleteClusterSuite))
}
