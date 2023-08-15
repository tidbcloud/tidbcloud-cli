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

package branch

//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"os"
//	"testing"
//
//	"tidbcloud-cli/internal"
//	"tidbcloud-cli/internal/config"
//	"tidbcloud-cli/internal/iostream"
//	"tidbcloud-cli/internal/mock"
//	"tidbcloud-cli/internal/service/cloud"
//	branchApi "tidbcloud-cli/pkg/tidbcloud/branch/client/branch_service"
//	branchModel "tidbcloud-cli/pkg/tidbcloud/branch/models"
//	connectInfoService "tidbcloud-cli/pkg/tidbcloud/connect_info/client/connect_info_service"
//	connectInfoModel "tidbcloud-cli/pkg/tidbcloud/connect_info/models"
//
//	"github.com/juju/errors"
//	"github.com/stretchr/testify/require"
//	"github.com/stretchr/testify/suite"
//)
//
//const getConnectInfoResultStr = `{
//    "ca_path": {
//        "Alpine": "/etc/ssl/cert.pem",
//        "Arch": "/etc/ssl/certs/ca-certificates.crt",
//        "CentOS": "/etc/pki/tls/certs/ca-bundle.crt",
//        "Debian": "/etc/ssl/certs/ca-certificates.crt",
//        "Fedora": "/etc/pki/tls/certs/ca-bundle.crt",
//        "OpenSUSE": "/etc/ssl/ca-bundle.pem",
//        "RedHat": "/etc/pki/tls/certs/ca-bundle.crt",
//        "Ubuntu": "/etc/ssl/certs/ca-certificates.crt",
//        "Windows": "<path_to_ca_cert>",
//        "macOS": "/etc/ssl/cert.pem",
//        "others": "<path_to_ca_cert>"
//    },
//    "client_data": [
//        {
//            "id": "mysql_cli",
//            "display_name": "MySQL CLI",
//            "language": "shell",
//            "content": [
//                {
//                    "type": "serverless",
//                    "comment": "",
//                    "connection_string": "mysql -u '${username}' -h ${host} -P ${port} -D test --ssl-mode=VERIFY_IDENTITY --ssl-ca=${ca_root_path} -p${password}",
//                    "connection_example": ""
//                },
//                {
//                    "type": "dedicated",
//                    "comment": "# version >= 5.7.11",
//                    "connection_string": "mysql -u '${username}' -h ${host} -P ${port} -D test --ssl-mode=VERIFY_IDENTITY --ssl-ca=ca.pem -p${password}",
//                    "connection_example": ""
//                }
//            ]
//        },
//        {
//            "id": "mysqlclient",
//            "display_name": "mysqlclient",
//            "language": "python",
//            "content": [
//                {
//                    "type": "serverless",
//                    "comment": "# Be sure to replace the parameters in the following connection string.\n# Requires mysqlclient package ('pip3 install mysqlclient'). Please check https://pypi.org/project/mysqlclient/ for install guide.",
//                    "connection_string": "host=\"${host}\", \nuser=\"${username}\", \npassword=\"${password}\", \nport=${port}, \ndatabase=\"test\",  \nssl={\"ca\": \"${ca_root_path}\"}",
//                    "connection_example": "import MySQLdb\n\nconnection = MySQLdb.connect(\n  host=\"${host}\",\n  port=${port},\n  user=\"${username}\",\n  password=\"${password}\",\n  database=\"test\",\n  ssl={\n    \"ca\": \"${ca_root_path}\"\n  }\n)\n\nwith connection:\n  with connection.cursor() as cursor:\n    cursor.execute(\"SELECT DATABASE();\")\n    m = cursor.fetchone()\n    print(m[0])"
//                },
//                {
//                    "type": "dedicated",
//                    "comment": "# Be sure to replace the parameters in the following connection string.\n# Requires mysqlclient package ('pip3 install mysqlclient'). Please check https://pypi.org/project/mysqlclient/ for install guide.",
//                    "connection_string": "host=\"${host}\", \nuser=\"${username}\", \npassword=\"${password}\", \nport=${port}, \ndatabase=\"test\",  \nssl={\"ca\": \"ca.pem\"}",
//                    "connection_example": "import MySQLdb\n\nconnection = MySQLdb.connect(\n  host=\"${host}\",\n  port=${port},\n  user=\"${username}\",\n  password=\"${password}\",\n  database=\"test\",\n  ssl={\n    \"ca\": \"ca.pem\"\n  }\n)\n\nwith connection:\n  with connection.cursor() as cursor:\n    cursor.execute(\"SELECT DATABASE();\")\n    m = cursor.fetchone()\n    print(m[0])"
//                }
//            ]
//        },
//        {
//            "id": "general",
//            "display_name": "General",
//            "language": "",
//            "content": [
//                {
//                    "type": "serverless",
//                    "comment": "",
//                    "connection_string": "Host: ${host}\nUsername: ${username}\nPort: ${port}\nPassword: ${password}",
//                    "connection_example": ""
//                },
//                {
//                    "type": "dedicated",
//                    "comment": "",
//                    "connection_string": "Host: ${host}\nUsername: ${username}\nPort: ${port}\nPassword: ${password}",
//                    "connection_example": ""
//                }
//            ]
//		}
//    ]
//}
//`
//
//type ConnectInfoSuite struct {
//	suite.Suite
//	h          *internal.Helper
//	mockClient *mock.TiDBCloudClient
//}
//
//func (suite *ConnectInfoSuite) SetupTest() {
//	if err := os.Setenv("NO_COLOR", "true"); err != nil {
//		suite.T().Error(err)
//	}
//
//	suite.mockClient = new(mock.TiDBCloudClient)
//	suite.h = &internal.Helper{
//		Client: func() (cloud.TiDBCloudClient, error) {
//			return suite.mockClient, nil
//		},
//		IOStreams: iostream.Test(),
//	}
//}
//
//func (suite *ConnectInfoSuite) TestConnectInfoArgs() {
//	assert := require.New(suite.T())
//
//	branchID := "12345"
//	clusterID := "12345"
//	clientCLI := "mysql_cli"
//	clientDriver := "python_mysqlclient"
//	clientParameter := "general"
//	operatingSystem := "macos"
//
//	connectInfoBody := &connectInfoModel.ConnectInfo{}
//	err := json.Unmarshal([]byte(getConnectInfoResultStr), connectInfoBody)
//	assert.Nil(err)
//	getConnectInfoResult := &connectInfoService.GetInfoOK{
//		Payload: connectInfoBody,
//	}
//	suite.mockClient.On("GetConnectInfo", connectInfoService.NewGetInfoParams()).
//		Return(getConnectInfoResult, nil)
//
//	branchBody := &branchModel.OpenapiBranch{}
//	err = json.Unmarshal([]byte(getBranchResultStr), branchBody)
//	assert.Nil(err)
//	getClusterResult := &branchApi.GetBranchOK{
//		Payload: branchBody,
//	}
//	suite.mockClient.On("GetBranch", branchApi.NewGetBranchParams().
//		WithBranchID(branchID).WithClusterID(clusterID)).
//		Return(getClusterResult, nil)
//
//	tests := []struct {
//		name         string
//		args         []string
//		err          error
//		stdoutString string
//		stderrString string
//	}{
//		{
//			name:         "get MySQl CLI connecting string",
//			args:         []string{"-b", branchID, "-c", clusterID, "--client", clientCLI, "--operating-system", operatingSystem},
//			stdoutString: "\nmysql -u '49dDUPpoxGXdsY9.root' -h gateway01.us-east-1.dev.shared.aws.tidbcloud.com -P 4000 -D test --ssl-mode=VERIFY_IDENTITY --ssl-ca=/etc/ssl/cert.pem -p${password}\n",
//		},
//		{
//			name: "get mysqlclient connecting string",
//			args: []string{"-b", branchID, "-c", clusterID, "--client", clientDriver, "--operating-system", operatingSystem},
//			stdoutString: `
//host="gateway01.us-east-1.dev.shared.aws.tidbcloud.com",
//user="49dDUPpoxGXdsY9.root",
//password="${password}",
//port=4000,
//database="test",
//ssl={"ca": "/etc/ssl/cert.pem"}
//`,
//		},
//		{
//			name: "get standard connection parameter",
//			args: []string{"-b", branchID, "-c", clusterID, "--client", clientParameter, "--operating-system", operatingSystem},
//			stdoutString: `
//Host:    gateway01.us-east-1.dev.shared.aws.tidbcloud.com
//Port:    4000
//User:    49dDUPpoxGXdsY9.root
//`,
//		},
//		{
//			name: "with unsupported client name",
//			args: []string{"-b", branchID, "-c", clusterID, "--client", "JAVA", "--operating-system", operatingSystem},
//			err:  errors.New(fmt.Sprintf("Unsupported client. Run \"%[1]s cluster connect-info -h\" to check supported clients list", config.CliName)),
//		},
//		{
//			name: "with unsupported operating system",
//			args: []string{"-b", branchID, "-c", clusterID, "--client", "python_mysqlclient", "--operating-system", "Manjaro"},
//			err:  errors.New(fmt.Sprintf("Unsupported operating system. Run \"%[1]s cluster connect-info -h\" to check supported operating systems list", config.CliName)),
//		},
//	}
//
//	for _, tt := range tests {
//		suite.T().Run(tt.name, func(t *testing.T) {
//			cmd := ConnectInfoCmd(suite.h)
//			suite.h.IOStreams.Out.(*bytes.Buffer).Reset()
//			suite.h.IOStreams.Err.(*bytes.Buffer).Reset()
//			cmd.SetArgs(tt.args)
//			err := cmd.Execute()
//			if err != nil {
//				assert.Equal(tt.err.Error(), err.Error())
//			}
//
//			assert.Equal(tt.stdoutString, suite.h.IOStreams.Out.(*bytes.Buffer).String())
//			assert.Equal(tt.stderrString, suite.h.IOStreams.Err.(*bytes.Buffer).String())
//			if tt.err == nil {
//				suite.mockClient.AssertExpectations(suite.T())
//			}
//		})
//	}
//}
//
//func TestConnectInfoSuite(t *testing.T) {
//	suite.Run(t, new(ConnectInfoSuite))
//}
