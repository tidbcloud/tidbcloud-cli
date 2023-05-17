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

package start

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/mock"
	"tidbcloud-cli/internal/service/cloud"
	connectInfoService "tidbcloud-cli/pkg/tidbcloud/connect_info/client/connect_info_service"
	connectInfoModel "tidbcloud-cli/pkg/tidbcloud/connect_info/models"

	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	mockTool "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const getConnectInfoResultStr = `{
    "ca_path": {
        "Alpine": "/etc/ssl/cert.pem",
        "Arch": "/etc/ssl/certs/ca-certificates.crt",
        "CentOS": "/etc/pki/tls/certs/ca-bundle.crt",
        "Debian": "/etc/ssl/certs/ca-certificates.crt",
        "Fedora": "/etc/pki/tls/certs/ca-bundle.crt",
        "OpenSUSE": "/etc/ssl/ca-bundle.pem",
        "RedHat": "/etc/pki/tls/certs/ca-bundle.crt",
        "Ubuntu": "/etc/ssl/certs/ca-certificates.crt",
        "Windows": "<path_to_ca_cert>",
        "macOS": "/etc/ssl/cert.pem",
        "others": "<path_to_ca_cert>"
    },
    "client_data": [
        {
            "id": "mysql_cli",
            "display_name": "MySQL CLI",
            "language": "shell",
            "content": [
                {
                    "type": "serverless",
                    "comment": "",
                    "connection_string": "mysql -u '${username}' -h ${host} -P ${port} -D test --ssl-mode=VERIFY_IDENTITY --ssl-ca=${ca_root_path} -p${password}",
                    "connection_example": ""
                },
                {
                    "type": "dedicated",
                    "comment": "# version >= 5.7.11",
                    "connection_string": "mysql -u '${username}' -h ${host} -P ${port} -D test --ssl-mode=VERIFY_IDENTITY --ssl-ca=ca.pem -p${password}",
                    "connection_example": ""
                }
            ]
        },
        {
            "id": "mysqlclient",
            "display_name": "mysqlclient",
            "language": "python",
            "content": [
                {
                    "type": "serverless",
                    "comment": "# Be sure to replace the parameters in the following connection string.\n# Requires mysqlclient package ('pip3 install mysqlclient'). Please check https://pypi.org/project/mysqlclient/ for install guide.",
                    "connection_string": "host=\"${host}\", \nuser=\"${username}\", \npassword=\"${password}\", \nport=${port}, \ndatabase=\"test\",  \nssl={\"ca\": \"${ca_root_path}\"}",
                    "connection_example": "import MySQLdb\n\nconnection = MySQLdb.connect(\n  host=\"${host}\",\n  port=${port},\n  user=\"${username}\",\n  password=\"${password}\",\n  database=\"test\",\n  ssl={\n    \"ca\": \"${ca_root_path}\"\n  }\n)\n\nwith connection:\n  with connection.cursor() as cursor:\n    cursor.execute(\"SELECT DATABASE();\")\n    m = cursor.fetchone()\n    print(m[0])"
                },
                {
                    "type": "dedicated",
                    "comment": "# Be sure to replace the parameters in the following connection string.\n# Requires mysqlclient package ('pip3 install mysqlclient'). Please check https://pypi.org/project/mysqlclient/ for install guide.",
                    "connection_string": "host=\"${host}\", \nuser=\"${username}\", \npassword=\"${password}\", \nport=${port}, \ndatabase=\"test\",  \nssl={\"ca\": \"ca.pem\"}",
                    "connection_example": "import MySQLdb\n\nconnection = MySQLdb.connect(\n  host=\"${host}\",\n  port=${port},\n  user=\"${username}\",\n  password=\"${password}\",\n  database=\"test\",\n  ssl={\n    \"ca\": \"ca.pem\"\n  }\n)\n\nwith connection:\n  with connection.cursor() as cursor:\n    cursor.execute(\"SELECT DATABASE();\")\n    m = cursor.fetchone()\n    print(m[0])"
                }
            ]
        },
        {
            "id": "general",
            "display_name": "General",
            "language": "",
            "content": [
                {
                    "type": "serverless",
                    "comment": "",
                    "connection_string": "Host: ${host}\nUsername: ${username}\nPort: ${port}\nPassword: ${password}",
                    "connection_example": ""
                },
                {
                    "type": "dedicated",
                    "comment": "",
                    "connection_string": "Host: ${host}\nUsername: ${username}\nPort: ${port}\nPassword: ${password}",
                    "connection_example": ""
                }
            ]
		}
    ]
}
`

type MySQLImportSuite struct {
	suite.Suite
	h          *internal.Helper
	mockClient *mock.TiDBCloudClient
	mockHelper *mock.MySQLHelper
}

func (suite *MySQLImportSuite) SetupTest() {
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		suite.T().Error(err)
	}

	suite.mockClient = new(mock.TiDBCloudClient)
	suite.mockHelper = new(mock.MySQLHelper)
	suite.h = &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			return suite.mockClient, nil
		},
		MySQLHelper: suite.mockHelper,
		IOStreams:   iostream.Test(),
	}
}

func (suite *MySQLImportSuite) TestMySQLImportArgs() {
	assert := require.New(suite.T())
	cachePath := "/tmp/test.sql"
	suite.mockHelper.On("GenerateSqlCachePath").Return(cachePath)
	suite.mockHelper.On("CheckMySQLClient").Return(nil)

	sourceHost := "127.0.0.1"
	sourcePort := "3306"
	sourceDatabase := "test"
	sourceTable := "table"
	sourceUser := "root"
	sourcePassword := "passwd"
	password := "passwd"
	projectID := "1234"
	clusterID := "4321"
	database := "mysql"

	dumpArgs := []string{
		"mysqldump",
		"-h", sourceHost,
		"-P", sourcePort,
		"-u", sourceUser,
		fmt.Sprintf("--password=%s", sourcePassword),
		"--skip-add-drop-table",
		"--skip-add-locks",
		"--skip-triggers",
		"-r",
		cachePath,
		sourceDatabase,
		sourceTable,
	}

	var caPath string
	if runtime.GOOS != "darwin" {
		suite.mockHelper.On("DownloadCaFile", mockTool.Anything).Return(nil)
		home, _ := os.UserHomeDir()
		caPath = filepath.Join(home, config.HomePath, "isrgrootx1.pem")
	} else {
		caPath = "/etc/ssl/cert.pem"
	}
	suite.mockHelper.On("DumpFromMySQL", mockTool.Anything, dumpArgs, cachePath).Return(nil)

	targetHost := "9.9.9.9"
	targetPort := int32(4000)
	targetUser := "root"

	importArgs := []string{
		"mysql",
		"-u",
		targetUser,
		"-h",
		targetHost,
		"-P",
		fmt.Sprint(targetPort),
		"-D",
		database,
		"--ssl-mode=VERIFY_IDENTITY",
		"--ssl-ca=" + caPath,
		"-p" + password,
	}
	suite.mockHelper.On("ImportToServerless", mockTool.Anything, importArgs, cachePath).Return(nil)
	suite.mockClient.On("GetCluster", cluster.NewGetClusterParams().
		WithProjectID(projectID).WithClusterID(clusterID)).Return(&cluster.GetClusterOK{
		Payload: &cluster.GetClusterOKBody{
			ClusterType: "DEVELOPER",
			Status: &cluster.GetClusterOKBodyStatus{
				ConnectionStrings: &cluster.GetClusterOKBodyStatusConnectionStrings{
					DefaultUser: targetUser,
					Standard: &cluster.GetClusterOKBodyStatusConnectionStringsStandard{
						Host: targetHost,
						Port: targetPort,
					},
				},
			},
		},
	}, nil)
	connectInfoBody := &connectInfoModel.ConnectInfo{}
	err := json.Unmarshal([]byte(getConnectInfoResultStr), connectInfoBody)
	assert.Nil(err)
	suite.mockClient.On("GetConnectInfo", connectInfoService.NewGetInfoParams()).
		Return(&connectInfoService.GetInfoOK{
			Payload: connectInfoBody,
		}, nil)

	tests := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "start import success",
			args: []string{"--project-id", projectID, "--cluster-id", clusterID, "--source-host", sourceHost, "--source-port", sourcePort, "--source-database", sourceDatabase, "--source-table", sourceTable, "--source-user", sourceUser, "--source-password", sourcePassword, "--password", password, "--database", database},
		},
		{
			name: "start import without required project-id flag",
			args: []string{"--cluster-id", clusterID, "--source-host", sourceHost, "--source-port", sourcePort, "--source-database", sourceDatabase, "--source-table", sourceTable, "--source-user", sourceUser, "--source-password", sourcePassword, "--password", password, "--database", database},
			err:  fmt.Errorf("required flag(s) \"project-id\" not set"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := MySQLCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			assert.Equal(tt.err, err)

			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
				suite.mockHelper.AssertExpectations(suite.T())
			}
		})
	}
}

func (suite *MySQLImportSuite) TestMySQLImportWithoutCreateTable() {
	assert := require.New(suite.T())
	cachePath := "/tmp/test.sql"
	suite.mockHelper.On("GenerateSqlCachePath").Return(cachePath)
	suite.mockHelper.On("CheckMySQLClient").Return(nil)

	sourceHost := "127.0.0.1"
	sourcePort := "3306"
	sourceDatabase := "test"
	sourceTable := "table"
	sourceUser := "root"
	sourcePassword := "passwd"
	password := "passwd"
	projectID := "1234"
	clusterID := "4321"
	database := "mysql"

	dumpArgs := []string{
		"mysqldump",
		"-h", sourceHost,
		"-P", sourcePort,
		"-u", sourceUser,
		fmt.Sprintf("--password=%s", sourcePassword),
		"--skip-add-drop-table",
		"--skip-add-locks",
		"--skip-triggers",
		"-r",
		cachePath,
		sourceDatabase,
		sourceTable,
	}

	var caPath string
	if runtime.GOOS != "darwin" {
		suite.mockHelper.On("DownloadCaFile", mockTool.Anything).Return(nil)
		home, _ := os.UserHomeDir()
		caPath = filepath.Join(home, config.HomePath, "isrgrootx1.pem")
	} else {
		caPath = "/etc/ssl/cert.pem"
	}
	suite.mockHelper.On("DumpFromMySQL", mockTool.Anything, dumpArgs, cachePath).Return(nil)

	targetHost := "9.9.9.9"
	targetPort := int32(4000)
	targetUser := "root"

	importArgs := []string{
		"mysql",
		"-u",
		targetUser,
		"-h",
		targetHost,
		"-P",
		fmt.Sprint(targetPort),
		"-D",
		database,
		"--ssl-mode=VERIFY_IDENTITY",
		"--ssl-ca=" + caPath,
		"-p" + password,
	}
	suite.mockHelper.On("ImportToServerless", mockTool.Anything, importArgs, cachePath).Return(nil)
	suite.mockClient.On("GetCluster", cluster.NewGetClusterParams().
		WithProjectID(projectID).WithClusterID(clusterID)).Return(&cluster.GetClusterOK{
		Payload: &cluster.GetClusterOKBody{
			ClusterType: "DEVELOPER",
			Status: &cluster.GetClusterOKBodyStatus{
				ConnectionStrings: &cluster.GetClusterOKBodyStatusConnectionStrings{
					DefaultUser: "root",
					Standard: &cluster.GetClusterOKBodyStatusConnectionStringsStandard{
						Host: targetHost,
						Port: targetPort,
					},
				},
			},
		},
	}, nil)
	connectInfoBody := &connectInfoModel.ConnectInfo{}
	err := json.Unmarshal([]byte(getConnectInfoResultStr), connectInfoBody)
	assert.Nil(err)
	suite.mockClient.On("GetConnectInfo", connectInfoService.NewGetInfoParams()).
		Return(&connectInfoService.GetInfoOK{
			Payload: connectInfoBody,
		}, nil)

	tests := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "start import success",
			args: []string{"--project-id", projectID, "--cluster-id", clusterID, "--source-host", sourceHost, "--source-port", sourcePort, "--source-database", sourceDatabase, "--source-table", sourceTable, "--source-user", sourceUser, "--source-password", sourcePassword, "--password", password, "--database", database, "--user", targetUser},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := MySQLCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			assert.Equal(tt.err, err)

			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
				suite.mockHelper.AssertExpectations(suite.T())
			}
		})
	}
}

func (suite *MySQLImportSuite) TestMySQLImportWithSpecificUser() {
	assert := require.New(suite.T())
	cachePath := "/tmp/test.sql"
	suite.mockHelper.On("GenerateSqlCachePath").Return(cachePath)
	suite.mockHelper.On("CheckMySQLClient").Return(nil)

	sourceHost := "127.0.0.1"
	sourcePort := "3306"
	sourceDatabase := "test"
	sourceTable := "table"
	sourceUser := "root"
	sourcePassword := "passwd"
	password := "passwd"
	projectID := "1234"
	clusterID := "4321"
	database := "mysql"

	dumpArgs := []string{
		"mysqldump",
		"-h", sourceHost,
		"-P", sourcePort,
		"-u", sourceUser,
		fmt.Sprintf("--password=%s", sourcePassword),
		"--skip-add-drop-table",
		"--skip-add-locks",
		"--skip-triggers",
		"--no-create-info",
		"-r",
		cachePath,
		sourceDatabase,
		sourceTable,
	}
	var caPath string
	if runtime.GOOS != "darwin" {
		suite.mockHelper.On("DownloadCaFile", mockTool.Anything).Return(nil)
		home, _ := os.UserHomeDir()
		caPath = filepath.Join(home, config.HomePath, "isrgrootx1.pem")
	} else {
		caPath = "/etc/ssl/cert.pem"
	}
	suite.mockHelper.On("DumpFromMySQL", mockTool.Anything, dumpArgs, cachePath).Return(nil)

	targetHost := "9.9.9.9"
	targetPort := int32(4000)
	targetUser := "root"

	importArgs := []string{
		"mysql",
		"-u",
		targetUser,
		"-h",
		targetHost,
		"-P",
		fmt.Sprint(targetPort),
		"-D",
		database,
		"--ssl-mode=VERIFY_IDENTITY",
		"--ssl-ca=" + caPath,
		"-p" + password,
	}
	suite.mockHelper.On("ImportToServerless", mockTool.Anything, importArgs, cachePath).Return(nil)
	suite.mockClient.On("GetCluster", cluster.NewGetClusterParams().
		WithProjectID(projectID).WithClusterID(clusterID)).Return(&cluster.GetClusterOK{
		Payload: &cluster.GetClusterOKBody{
			ClusterType: "DEVELOPER",
			Status: &cluster.GetClusterOKBodyStatus{
				ConnectionStrings: &cluster.GetClusterOKBodyStatusConnectionStrings{
					DefaultUser: targetUser,
					Standard: &cluster.GetClusterOKBodyStatusConnectionStringsStandard{
						Host: targetHost,
						Port: targetPort,
					},
				},
			},
		},
	}, nil)
	connectInfoBody := &connectInfoModel.ConnectInfo{}
	err := json.Unmarshal([]byte(getConnectInfoResultStr), connectInfoBody)
	assert.Nil(err)
	suite.mockClient.On("GetConnectInfo", connectInfoService.NewGetInfoParams()).
		Return(&connectInfoService.GetInfoOK{
			Payload: connectInfoBody,
		}, nil)

	tests := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "start import success",
			args: []string{"--project-id", projectID, "--cluster-id", clusterID, "--source-host", sourceHost, "--source-port", sourcePort, "--source-database", sourceDatabase, "--source-table", sourceTable, "--source-user", sourceUser, "--source-password", sourcePassword, "--password", password, "--database", database, "--skip-create-table"},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cmd := MySQLCmd(suite.h)
			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			assert.Equal(tt.err, err)

			if tt.err == nil {
				suite.mockClient.AssertExpectations(suite.T())
				suite.mockHelper.AssertExpectations(suite.T())
			}
		})
	}
}

func TestMySQLImportSuite(t *testing.T) {
	suite.Run(t, new(MySQLImportSuite))
}
