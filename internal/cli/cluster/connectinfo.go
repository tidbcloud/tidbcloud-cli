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

package cluster

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	connectInfoModel "tidbcloud-cli/pkg/tidbcloud/connect_info/models"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

// xxxDisplayName is used to display in interactive mode
// xxxInputName is used to input in non-interactive mode and display in help message
const (
	GeneralParameterID              string = "general"
	GeneralParameterDisplayName     string = "General"
	GeneralParameterInputName       string = "general"
	MysqlCliID                      string = "mysql_cli"
	MysqlCliDisplayName             string = "MySQL CLI"
	MysqlCliInputName               string = "mysql_cli"
	MyCliID                         string = "mycli"
	MyCliDisplayName                string = "MyCLI"
	MyCliInputName                  string = "mycli"
	LibMysqlClientID                string = "libmysqlclient"
	LibMysqlClientDisplayName       string = "libmysqlclient"
	LibMysqlClientInputName         string = "libmysqlclient"
	MysqlClientID                   string = "mysqlclient"
	MysqlClientDisplayName          string = "mysqlclient (Python)"
	MysqlClientInputName            string = "python_mysqlclient"
	PyMysqlID                       string = "pymysql"
	PyMysqlDisplayName              string = "PyMySQL"
	PyMysqlInputName                string = "pymysql"
	MysqlConnectorPythonID          string = "mysql_connector_python"
	MysqlConnectorPythonDisplayName string = "MySQL Connector/Python"
	MysqlConnectorPythonInputName   string = "mysql_connector_python"
	MysqlConnectorJavaID            string = "mysql_connector_java"
	MysqlConnectorJavaDisplayName   string = "MySQL Connector/Java"
	MysqlConnectorJavaInputName     string = "mysql_connector_java"
	GoMysqlDriverID                 string = "go_mysql_driver"
	GoMysqlDriverDisplayName        string = "Go MySQL Driver"
	GoMysqlDriverInputName          string = "go_mysql_driver"
	NodeMysql2ID                    string = "node_mysql2"
	NodeMysql2DisplayName           string = "Node MySQL 2"
	NodeMysql2InputName             string = "node_mysql2"
	Mysql2RubyID                    string = "mysql2_ruby"
	Mysql2RubyDisplayName           string = "Mysql2 (Ruby)"
	Mysql2RubyInputName             string = "ruby_mysql2"
	MysqliID                        string = "mysqli"
	MysqliDisplayName               string = "MySQLi (PHP)"
	MysqliInputName                 string = "php_mysqli"
	MysqlRustID                     string = "mysql_rust"
	MysqlRustDisplayName            string = "mysql (Rust)"
	MysqlRustInputName              string = "rust_mysql"
	MybatisID                       string = "mybatis"
	MybatisDisplayName              string = "MyBatis"
	MybatisInputName                string = "mybatis"
	HibernateID                     string = "hibernate"
	HibernateDisplayName            string = "Hibernate"
	HibernateInputName              string = "hibernate"
	SpringBootID                    string = "spring_boot"
	SpringBootDisplayName           string = "Spring Boot"
	SpringBootInputName             string = "spring_boot"
	GormID                          string = "gorm"
	GormDisplayName                 string = "GORM"
	GormInputName                   string = "gorm"
	PrismaID                        string = "prisma"
	PrismaDisplayName               string = "Prisma"
	PrismaInputName                 string = "prisma"
	SequelizeID                     string = "sequelize_mysql2"
	SequelizeDisplayName            string = "Sequelize (mysql2)"
	SequelizeInputName              string = "sequelize_mysql2"
	DjangoID                        string = "django_tidb"
	DjangoDisplayName               string = "Django (django_tidb)"
	DjangoInputName                 string = "django_tidb"
	SQLAlchemyID                    string = "sqlalchemy_mysqlclient"
	SqlAlchemyDisplayName           string = "SQLAlchemy (mysqlclient)"
	SqlAlchemyInputName             string = "sqlalchemy_mysqlclient"
	ActiveRecordID                  string = "active_record"
	ActiveRecordDisplayName         string = "Active Record"
	ActiveRecordInputName           string = "active_record"
)

type connectInfoOpts struct {
	interactive bool
}

func (c connectInfoOpts) NonInteractiveRequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
	}
}

func (c connectInfoOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
		flag.ClientName,
		flag.OperatingSystem,
	}
}

// Display clients name orderly in interactive mode
var connectClientsList = []string{
	// pure parameter
	GeneralParameterDisplayName,

	// CLI
	MysqlCliDisplayName,
	MyCliDisplayName,

	// driver
	LibMysqlClientDisplayName,
	MysqlClientDisplayName,
	PyMysqlDisplayName,
	MysqlConnectorPythonDisplayName,
	MysqlConnectorJavaDisplayName,
	GoMysqlDriverDisplayName,
	NodeMysql2DisplayName,
	Mysql2RubyDisplayName,
	MysqliDisplayName,
	MysqlRustDisplayName,

	// ORM
	MybatisDisplayName,
	HibernateDisplayName,
	SpringBootDisplayName,
	GormDisplayName,
	PrismaDisplayName,
	SequelizeDisplayName,
	DjangoDisplayName,
	SqlAlchemyDisplayName,
	ActiveRecordDisplayName,
}

// Display clients name orderly in help message
var connectClientsListForHelp = []string{
	// pure parameter
	GeneralParameterInputName,

	// CLI
	MysqlCliInputName,
	MyCliInputName,

	// driver
	LibMysqlClientInputName,
	MysqlClientInputName,
	PyMysqlInputName,
	MysqlConnectorPythonInputName,
	MysqlConnectorJavaInputName,
	GoMysqlDriverInputName,
	NodeMysql2InputName,
	Mysql2RubyInputName,
	MysqliInputName,
	MysqlRustInputName,

	// ORM
	MybatisInputName,
	HibernateInputName,
	SpringBootInputName,
	GormInputName,
	PrismaInputName,
	SequelizeInputName,
	DjangoInputName,
	SqlAlchemyInputName,
	ActiveRecordInputName,
}

var clientsForInteractiveMap = map[string]string{
	// pure parameter
	GeneralParameterDisplayName: GeneralParameterID,

	// CLI
	MysqlCliDisplayName: MysqlCliID,
	MyCliDisplayName:    MyCliID,

	// driver
	LibMysqlClientDisplayName:       LibMysqlClientID,
	MysqlClientDisplayName:          MysqlClientID,
	PyMysqlDisplayName:              PyMysqlID,
	MysqlConnectorPythonDisplayName: MysqlConnectorPythonID,
	MysqlConnectorJavaDisplayName:   MysqlConnectorJavaID,
	GoMysqlDriverDisplayName:        GoMysqlDriverID,
	NodeMysql2DisplayName:           NodeMysql2ID,
	Mysql2RubyDisplayName:           Mysql2RubyID,
	MysqliDisplayName:               MysqliID,
	MysqlRustDisplayName:            MysqlRustID,

	// ORM
	MybatisDisplayName:      MybatisID,
	HibernateDisplayName:    HibernateID,
	SpringBootDisplayName:   SpringBootID,
	GormDisplayName:         GormID,
	PrismaDisplayName:       PrismaID,
	SequelizeDisplayName:    SequelizeID,
	DjangoDisplayName:       DjangoID,
	SqlAlchemyDisplayName:   SQLAlchemyID,
	ActiveRecordDisplayName: ActiveRecordID,
}

var clientsForHelpMap = map[string]string{
	// pure parameter
	GeneralParameterInputName: GeneralParameterID,

	// CLI
	MysqlCliInputName: MysqlCliID,
	MyCliInputName:    MyCliID,

	// driver
	LibMysqlClientInputName:       LibMysqlClientID,
	MysqlClientInputName:          MysqlClientID,
	PyMysqlInputName:              PyMysqlID,
	MysqlConnectorPythonInputName: MysqlConnectorPythonID,
	MysqlConnectorJavaInputName:   MysqlConnectorJavaID,
	GoMysqlDriverInputName:        GoMysqlDriverID,
	NodeMysql2InputName:           NodeMysql2ID,
	Mysql2RubyInputName:           Mysql2RubyID,
	MysqliInputName:               MysqliID,
	MysqlRustInputName:            MysqlRustID,

	// ORM
	MybatisInputName:      MybatisID,
	HibernateInputName:    HibernateID,
	SpringBootInputName:   SpringBootID,
	GormInputName:         GormID,
	PrismaInputName:       PrismaID,
	SequelizeInputName:    SequelizeID,
	DjangoInputName:       DjangoID,
	SqlAlchemyInputName:   SQLAlchemyID,
	ActiveRecordInputName: ActiveRecordID,
}

// Display operating system orderly in interactive mode
var operatingSystemList = []string{
	"macOS/Alpine",
	"CentOS/RedHat/Fedora",
	"Debian/Ubuntu/Arch",
	"Windows",
	"OpenSUSE",
	"Others",
}

// Display operating system orderly in help message
var operatingSystemListForHelp = []string{
	"macOS",
	"Windows",
	"Ubuntu",
	"CentOS",
	"RedHat",
	"Fedora",
	"Debian",
	"Arch",
	"OpenSUSE",
	"Alpine",
	"Others",
}

var caPath = map[string]string{
	"macos":    "/etc/ssl/cert.pem",
	"alpine":   "/etc/ssl/cert.pem",
	"centos":   "/etc/pki/tls/certs/ca-bundle.crt",
	"redhat":   "/etc/pki/tls/certs/ca-bundle.crt",
	"fedora":   "/etc/pki/tls/certs/ca-bundle.crt",
	"debian":   "/etc/ssl/certs/ca-certificates.crt",
	"ubuntu":   "/etc/ssl/certs/ca-certificates.crt",
	"arch":     "/etc/ssl/certs/ca-certificates.crt",
	"windows":  "<path_to_ca_cert>",
	"opensuse": "/etc/ssl/ca-bundle.pem",
	"others":   "<path_to_ca_cert>",
}

// Cluster type
const (
	SERVERLESS = "SERVERLESS"
	DEVELOPER  = "DEVELOPER"
)

func ConnectInfoCmd(h *internal.Helper) *cobra.Command {
	opts := connectInfoOpts{
		interactive: true,
	}

	// Detect operating system
	// TODO: detect linux operating system name
	os := runtime.GOOS
	if os == "windows" {
		os = "Windows"
	} else {
		os = "macOS"
	}

	cmd := &cobra.Command{
		Use:   "connect-info",
		Short: "Get connection string for the specified cluster",
		Example: fmt.Sprintf(`  Get connection string in interactive mode:
  $ %[1]s cluster connect-info

  Get connection string in non-interactive mode:
  $ %[1]s cluster connect-info --project-id <project-id> --cluster-id <cluster-id> --client <client-name> --operating-system <operating-system>
`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// mark required flags in non-interactive mode
			flags = opts.NonInteractiveRequiredFlags()
			if !opts.interactive {
				for _, fn := range flags {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// flags
			var projectID, clusterID, client, operatingSystem string

			// Get TiDBCloudClient
			d, err := h.Client()
			if err != nil {
				return err
			}

			if opts.interactive { // interactive mode
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// Get project id
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				projectID = project.ID

				// Get cluster id
				cluster, err := cloud.GetSelectedCluster(projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				// Get client
				clientNameForInteractive, err := cloud.GetSelectedConnectClient(connectClientsList)
				if err != nil {
					return err
				}
				client = clientsForInteractiveMap[clientNameForInteractive]

				if os != "" && os != "linux" {
					for id, value := range operatingSystemList {
						if strings.Contains(value, os) {
							operatingSystemValueWithFlag := operatingSystemList[id] + " (Detected)"
							operatingSystemList = append([]string{operatingSystemValueWithFlag}, append(operatingSystemList[:id], operatingSystemList[id+1:]...)...)
							break
						}
					}
				}

				// Get operating system
				operatingSystemCombination, err := cloud.GetSelectedConnectOs(operatingSystemList)
				if err != nil {
					return err
				}
				operatingSystem = strings.Split(operatingSystemCombination, "/")[0]

			} else { // non-interactive mode
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return err
				}

				projectID, err = cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return err
				}

				clientNameForHelp, err := cmd.Flags().GetString(flag.ClientName)
				if err != nil {
					return err
				}
				if v, ok := clientsForHelpMap[strings.ToLower(clientNameForHelp)]; ok {
					client = v
				} else {
					return errors.New(fmt.Sprintf("Unsupported client. Run \"%[1]s cluster connect-info -h\" to check supported clients list", config.CliName))
				}

				operatingSystem, err = cmd.Flags().GetString(flag.OperatingSystem)
				if err != nil {
					return err
				}
				if !contains(operatingSystem, operatingSystemListForHelp) {
					return errors.New(fmt.Sprintf("Unsupported operating system. Run \"%[1]s cluster connect-info -h\" to check supported operating systems list", config.CliName))
				}
			}

			// Get cluster info
			params := clusterApi.NewGetClusterParams().WithProjectID(projectID).WithClusterID(clusterID)
			clusterInfo, err := d.GetCluster(params)
			if err != nil {
				return err
			}

			// Resolve cluster information
			// Get connect parameter
			defaultUser := clusterInfo.Payload.Status.ConnectionStrings.DefaultUser
			host := clusterInfo.Payload.Status.ConnectionStrings.Standard.Host
			port := strconv.Itoa(int(clusterInfo.Payload.Status.ConnectionStrings.Standard.Port))
			clusterType := clusterInfo.Payload.ClusterType
			if clusterType == DEVELOPER {
				clusterType = SERVERLESS
			}

			// Get connection string
			connectInfo, err := cloud.RetrieveConnectInfo(d)
			if err != nil {
				return err
			}
			connectionString, err := generateConnectionString(connectInfo, client, host, defaultUser, port, clusterType, operatingSystem)
			if err != nil {
				return err
			}
			fmt.Fprintln(h.IOStreams.Out)
			fmt.Fprintln(h.IOStreams.Out, connectionString)

			return nil
		},
	}

	cmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "Project ID")
	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	cmd.Flags().String(flag.ClientName, MysqlCliInputName, strings.ReplaceAll(fmt.Sprintf("Connected client. Supported clients: %+q", connectClientsListForHelp), "\" \"", "\", \""))
	cmd.Flags().String(flag.OperatingSystem, os, strings.ReplaceAll(fmt.Sprintf("Operating system name. Supported operating systems: %q", operatingSystemListForHelp), "\" \"", "\", \""))

	return cmd
}

func generateConnectionString(connectInfo *connectInfoModel.ConnectInfo, client string, host string, user string, port string, clusterType string, operatingSystem string) (string, error) {
	if client == GeneralParameterID {
		return fmt.Sprintf(`Host:    %s
Port:    %s
User:    %s`,
			host, port, user), nil
	}

	for _, clientData := range connectInfo.ClientData {
		// find user chose client
		if strings.EqualFold(clientData.ID, client) {
			for _, content := range clientData.Content {
				if strings.EqualFold(clusterType, content.Type) {
					connectionString := content.ConnectionString
					connectionString = strings.Replace(connectionString, "${host}", host, -1)
					connectionString = strings.Replace(connectionString, "${username}", user, -1)
					connectionString = strings.Replace(connectionString, "${port}", port, -1)
					connectionString = strings.Replace(connectionString, "${ca_root_path}", caPath[strings.ToLower(operatingSystem)], -1)
					return connectionString, nil
				}
			}
		}
	}
	return "", errors.New("failed to generate connection string")
}

func contains(str string, vec []string) bool {
	for _, v := range vec {
		if strings.EqualFold(str, v) {
			return true
		}
	}
	return false
}
