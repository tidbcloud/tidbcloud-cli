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

package util

import (
	"fmt"
	"strings"

	"tidbcloud-cli/pkg/tidbcloud/connect_info/models"

	"github.com/juju/errors"
)

type ConnectStringUsage string

const (
	GolangCommand = "golang"
	Shell         = "shell"
)

// xxxDisplayName is used to display in interactive mode
// xxxInputName is used to input in non-interactive mode and display in help message
const (
	GeneralParameterID          string = "general"
	GeneralParameterDisplayName string = "General"
	GeneralParameterInputName   string = "general"

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

var CaPath = map[string]string{
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

func GenerateConnectionString(connectInfo *models.ConnectInfo, client string, host string, user string, port string, clusterType string, operatingSystem string, usage ConnectStringUsage) (string, error) {
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
					if usage == GolangCommand {
						connectionString = strings.Replace(connectionString, "'${username}'", user, -1)
					} else {
						connectionString = strings.Replace(connectionString, "${username}", user, -1)
					}

					connectionString = strings.Replace(connectionString, "${port}", port, -1)
					caPath, exist := CaPath[strings.ToLower(operatingSystem)]
					if exist {
						connectionString = strings.Replace(connectionString, "${ca_root_path}", caPath, -1)
					} else {
						connectionString = strings.Replace(connectionString, "${ca_root_path}", "<path_to_ca_cert>", -1)
					}
					return connectionString, nil
				}
			}
		}
	}
	return "", errors.New("failed to generate connection string")
}
