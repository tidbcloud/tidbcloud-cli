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

// xxxForInteracive is used to display in interactive mode
// xxxForHelp is used to display in help message and input in non-interactive mode
const (
	GeneralParameterID                 string = "standard_connection_parameter"
	GeneralParameterForInteractive     string = "General Parameter"
	GeneralParameterForHelp            string = "general_parameter"
	MysqlCliID                         string = "mysql_cli"
	MysqlCliForInteractive             string = "MySQL CLI"
	MysqlCliForHelp                    string = "mysql_cli"
	MyCliID                            string = "mycli"
	MyCliForInteractive                string = "MyCLI"
	MyCliForHelp                       string = "mycli"
	LibMysqlClientID                   string = "libmysqlclient"
	LibMysqlClientForInteractive       string = "libmysqlclient"
	LibMysqlClientForHelp              string = "libmysqlclient"
	MysqlClientID                      string = "mysqlclient"
	MysqlClientForInteractive          string = "mysqlclient (Python)"
	MysqlClientForHelp                 string = "python_mysqlclient"
	PyMysqlID                          string = "pymysql"
	PyMysqlForInteractive              string = "PyMySQL"
	PyMysqlForHelp                     string = "pymysql"
	MysqlConnectorPythonID             string = "mysql_connector_python"
	MysqlConnectorPythonForInteractive string = "MySQL Connector/Python"
	MysqlConnectorPythonForHelp        string = "mysql_connector_python"
	MysqlConnectorJavaID               string = "mysql_connector_java"
	MysqlConnectorJavaForInteractive   string = "MySQL Connector/Java"
	MysqlConnectorJavaForHelp          string = "mysql_connector_java"
	GoMysqlDriverID                    string = "go_mysql_driver"
	GoMysqlDriverForInteractive        string = "Go MySQL Driver"
	GoMysqlDriverForHelp               string = "go_mysql_driver"
	NodeMysql2ID                       string = "node_mysql_2"
	NodeMysql2ForInteractive           string = "Node MySQL 2"
	NodeMysql2ForHelp                  string = "node.js_mysql2"
	Mysql2RubyID                       string = "mysql2_ruby"
	Mysql2RubyForInteractive           string = "Mysql2 (Ruby)"
	Mysql2RubyForHelp                  string = "ruby_mysql2"
	MysqliID                           string = "mysqli"
	MysqliForInteractive               string = "MySQLi (PHP)"
	MysqliForHelp                      string = "php_mysqli"
	MysqlRustID                        string = "mysql_rust"
	MysqlRustForInteractive            string = "mysql (Rust)"
	MysqlRustForHelp                   string = "rust_mysql"
	MybatisID                          string = "mybatis"
	MybatisForInteractive              string = "MyBatis"
	MybatisForHelp                     string = "mybatis"
	HibernateID                        string = "hibernate"
	HibernateForInteractive            string = "Hibernate"
	HibernateForHelp                   string = "hibernate"
	SpringBootID                       string = "spring_boot"
	SpringBootForInteractive           string = "Spring Boot"
	SpringBootForHelp                  string = "spring_boot"
	GormID                             string = "gorm"
	GormForInteractive                 string = "GORM"
	GormForHelp                        string = "gorm"
	PrismaID                           string = "prisma"
	PrismaForInteractive               string = "Prisma"
	PrismaForHelp                      string = "prisma"
	SequelizeID                        string = "sequelize"
	SequelizeForInteractive            string = "Sequelize (mysql2)"
	SequelizeForHelp                   string = "sequelize"
	DjangoID                           string = "django"
	DjangoForInteractive               string = "Django (django_tidb)"
	DjangoForHelp                      string = "django"
	SQLAlchemyID                       string = "sqlalchemy"
	SqlAlchemyForInteractive           string = "SQLAlchemy (mysqlclient)"
	SqlAlchemyForHelp                  string = "sqlalchemy"
	ActiveRecordID                     string = "active_record"
	ActiveRecordForInteractive         string = "Active Record"
	ActiveRecordForHelp                string = "active_record"
)

type connectInfoOpts struct {
	interactive bool
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
	GeneralParameterForInteractive,

	// CLI
	MysqlCliForInteractive,
	MyCliForInteractive,

	// driver
	LibMysqlClientForInteractive,
	MysqlClientForInteractive,
	PyMysqlForInteractive,
	MysqlConnectorPythonForInteractive,
	MysqlConnectorJavaForInteractive,
	GoMysqlDriverForInteractive,
	NodeMysql2ForInteractive,
	Mysql2RubyForInteractive,
	MysqliForInteractive,
	MysqlRustForInteractive,

	// ORM
	MybatisForInteractive,
	HibernateForInteractive,
	SpringBootForInteractive,
	GormForInteractive,
	PrismaForInteractive,
	SequelizeForInteractive,
	DjangoForInteractive,
	SqlAlchemyForInteractive,
	ActiveRecordForInteractive,
}

// Display clients name orderly in help message
var connectClientsListForHelp = []string{
	// pure parameter
	GeneralParameterForHelp,

	// CLI
	MysqlCliForHelp,
	MyCliForHelp,

	// driver
	LibMysqlClientForHelp,
	MysqlClientForHelp,
	PyMysqlForHelp,
	MysqlConnectorPythonForHelp,
	MysqlConnectorJavaForHelp,
	GoMysqlDriverForHelp,
	NodeMysql2ForHelp,
	Mysql2RubyForHelp,
	MysqliForHelp,
	MysqlRustForHelp,

	// ORM
	MybatisForHelp,
	HibernateForHelp,
	SpringBootForHelp,
	GormForHelp,
	PrismaForHelp,
	SequelizeForHelp,
	DjangoForHelp,
	SqlAlchemyForHelp,
	ActiveRecordForHelp,
}

var clientsForInteractiveMap = map[string]string{
	// pure parameter
	GeneralParameterForInteractive: GeneralParameterID,

	// CLI
	MysqlCliForInteractive: MysqlCliID,
	MyCliForInteractive:    MyCliID,

	// driver
	LibMysqlClientForInteractive:       LibMysqlClientID,
	MysqlClientForInteractive:          MysqlClientID,
	PyMysqlForInteractive:              PyMysqlID,
	MysqlConnectorPythonForInteractive: MysqlConnectorPythonID,
	MysqlConnectorJavaForInteractive:   MysqlConnectorJavaID,
	GoMysqlDriverForInteractive:        GoMysqlDriverID,
	NodeMysql2ForInteractive:           NodeMysql2ID,
	Mysql2RubyForInteractive:           Mysql2RubyID,
	MysqliForInteractive:               MysqliID,
	MysqlRustForInteractive:            MysqlRustID,

	// ORM
	MybatisForInteractive:      MybatisID,
	HibernateForInteractive:    HibernateID,
	SpringBootForInteractive:   SpringBootID,
	GormForInteractive:         GormID,
	PrismaForInteractive:       PrismaID,
	SequelizeForInteractive:    SequelizeID,
	DjangoForInteractive:       DjangoID,
	SqlAlchemyForInteractive:   SQLAlchemyID,
	ActiveRecordForInteractive: ActiveRecordID,
}

var clientsForHelpMap = map[string]string{
	// pure parameter
	GeneralParameterForHelp: GeneralParameterID,

	// CLI
	MysqlCliForHelp: MysqlCliID,
	MyCliForHelp:    MyCliID,

	// driver
	LibMysqlClientForHelp:       LibMysqlClientID,
	MysqlClientForHelp:          MysqlClientID,
	PyMysqlForHelp:              PyMysqlID,
	MysqlConnectorPythonForHelp: MysqlConnectorPythonID,
	MysqlConnectorJavaForHelp:   MysqlConnectorJavaID,
	GoMysqlDriverForHelp:        GoMysqlDriverID,
	NodeMysql2ForHelp:           NodeMysql2ID,
	Mysql2RubyForHelp:           Mysql2RubyID,
	MysqliForHelp:               MysqliID,
	MysqlRustForHelp:            MysqlRustID,

	// ORM
	MybatisForHelp:      MybatisID,
	HibernateForHelp:    HibernateID,
	SpringBootForHelp:   SpringBootID,
	GormForHelp:         GormID,
	PrismaForHelp:       PrismaID,
	SequelizeForHelp:    SequelizeID,
	DjangoForHelp:       DjangoID,
	SqlAlchemyForHelp:   SQLAlchemyID,
	ActiveRecordForHelp: ActiveRecordID,
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

				// Detect operating system
				// TODO: detect linux operating system name
				goOS := runtime.GOOS
				if goOS == "darwin" {
					goOS = "macOS"
				} else if goOS == "windows" {
					goOS = "Windows"
				}
				if goOS != "" && goOS != "linux" {
					for id, value := range operatingSystemList {
						if strings.Contains(value, goOS) {
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
	cmd.Flags().String(flag.ClientName, "", fmt.Sprintf("Connected client. Supported clients: %q", connectClientsListForHelp))
	cmd.Flags().String(flag.OperatingSystem, "", fmt.Sprintf("Operating system name. Supported operating systems: %q", operatingSystemListForHelp))

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
