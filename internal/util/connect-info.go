// Copyright 2024 PingCAP, Inc.
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

// Code is generated. DO NOT EDIT.

package util

type ConnectInfo struct {
	Endpoint   []Endpoint
	Os         []Os
	Ca         []Ca
	Client     []Client
	Variable   []Variable
	Connection []Connection
}
type Endpoint struct {
	ID   string
	Name string
}

func (e Endpoint) String() string {
	return e.Name
}

type Os struct {
	ID   string
	Name string
}

func (o Os) String() string {
	return o.Name
}

type Ca struct {
	Os   string
	Type string
	Path string
}
type Options struct {
	ID   string
	Name string
}

func (o Options) String() string {
	return o.Name
}

type Client struct {
	ID      string
	Name    string
	Options []Options
}

func (c Client) String() string {
	return c.Name
}

type Variable struct {
	ID          string
	Placeholder string
}
type Connection struct {
	Endpoint   string
	Client     string
	Type       string
	Path       string
	DownloadCa []string
	Doc        string
	Content    string
}

var ConnectInfoEndpoint = []Endpoint{
	{ID: "public", Name: "Public"},
	{ID: "private", Name: "Private"},
}
var ConnectInfoOs = []Os{
	{ID: "macos", Name: "macOS"},
	{ID: "debian", Name: "Debian/Ubuntu/Arch"},
	{ID: "centos", Name: "CentOS/RedHat/Fedora"},
	{ID: "alpine", Name: "Alpine"},
	{ID: "suse", Name: "OpenSUSE"},
	{ID: "windows", Name: "Windows"},
	{ID: "other", Name: "Others"},
}
var ConnectInfoCa = []Ca{
	{Os: "macos", Type: "local", Path: "/etc/ssl/cert.pem"},
	{Os: "debian", Type: "local", Path: "/etc/ssl/certs/ca-certificates.crt"},
	{Os: "centos", Type: "local", Path: "/etc/pki/tls/certs/ca-bundle.crt"},
	{Os: "alpine", Type: "local", Path: "/etc/ssl/cert.pem"},
	{Os: "suse", Type: "local", Path: "/etc/ssl/ca-bundle.pem"},
	{Os: "windows", Type: "link", Path: "https://letsencrypt.org/certs/isrgrootx1.pem"},
	{Os: "other", Type: "link", Path: "https://letsencrypt.org/certs/isrgrootx1.pem"},
}
var ConnectInfoClient = []Client{
	{ID: "general", Name: "General", Options: []Options{
		{ID: "generalparams", Name: "Parameters"},
		{ID: "generalstring", Name: "Connection String"},
	}},
	{ID: "mysqlcli", Name: "MySQL CLI", Options: []Options{}},
	{ID: "mariadbcli", Name: "MariaDB CLI", Options: []Options{}},
	{ID: "dotenv", Name: ".env", Options: []Options{}},
	{ID: "mysqlconnectorj", Name: "Java(MySQL Connector)", Options: []Options{}},
	{ID: "mariadbconnectorj", Name: "Java(MariaDB Connector)", Options: []Options{}},
	{ID: "pythonmysqlclient", Name: "Python(mysqlclient)", Options: []Options{}},
	{ID: "pythonpymysql", Name: "Python(PyMySQL)", Options: []Options{}},
	{ID: "pythonmysqlconnector", Name: "Python(MySQL Connector)", Options: []Options{}},
	{ID: "sqlalchemy", Name: "SQLAlchemy", Options: []Options{
		{ID: "sqlalchemymysqlclient", Name: "mysqlclient"},
		{ID: "sqlalchemypymysql", Name: "PyMySQL"},
		{ID: "sqlalchemymysqlconnector", Name: "MySQL Connector"},
	}},
	{ID: "django", Name: "Django", Options: []Options{}},
	{ID: "nodemysql2", Name: "Node.js(mysql2)", Options: []Options{}},
	{ID: "prisma", Name: "Prisma", Options: []Options{}},
	{ID: "serverlessjs", Name: "Serverless Driver", Options: []Options{}},
	{ID: "go", Name: "Go", Options: []Options{}},
	{ID: "rails", Name: "Rails", Options: []Options{}},
	{ID: "wordpress", Name: "WordPress", Options: []Options{}},
	{ID: "datagrip", Name: "DataGrip", Options: []Options{}},
	{ID: "vscode", Name: "VS Code", Options: []Options{}},
	{ID: "dbeaver", Name: "DBeaver", Options: []Options{}},
	{ID: "mysqlworkbench", Name: "MySQL Workbench", Options: []Options{}},
	{ID: "navicat", Name: "Navicat", Options: []Options{}},
}
var ConnectInfoVariable = []Variable{
	{ID: "host", Placeholder: "<HOST>"},
	{ID: "port", Placeholder: "<PORT>"},
	{ID: "username", Placeholder: "<USERNAME>"},
	{ID: "password", Placeholder: "<PASSWORD>"},
	{ID: "database", Placeholder: "<DB>"},
	{ID: "ca_path", Placeholder: "<CA_PATH>"},
}
var ConnectInfoConnection = []Connection{
	{Endpoint: "public", Client: "generalparams", Type: "parameter", Path: "./public/general/params", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
CA ${ca_path}
`},
	{Endpoint: "private", Client: "generalparams", Type: "parameter", Path: "./private/general/params", DownloadCa: []string{}, Doc: "", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
`},
	{Endpoint: "public", Client: "generalstring", Type: "code", Path: "./public/general/string", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "", Content: `mysql://${username}:${password}@${host}:${port}/${database}
`},
	{Endpoint: "private", Client: "generalstring", Type: "code", Path: "./private/general/string", DownloadCa: []string{}, Doc: "", Content: `mysql://${username}:${password}@${host}:${port}/${database}
`},
	{Endpoint: "public", Client: "dotenv", Type: "code", Path: "./public/db.env", DownloadCa: []string{}, Doc: "", Content: `DB_HOST=${host}
DB_PORT=${port}
DB_USERNAME=${username}
DB_PASSWORD=${password}
DB_DATABASE=${database}
`},
	{Endpoint: "private", Client: "dotenv", Type: "code", Path: "./private/db.env", DownloadCa: []string{}, Doc: "", Content: `DB_HOST=${host}
DB_PORT=${port}
DB_USERNAME=${username}
DB_PASSWORD=${password}
DB_DATABASE=${database}
`},
	{Endpoint: "public", Client: "mysqlcli", Type: "code", Path: "./public/mysql.sh", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "", Content: `mysql --comments -u '${username}' -h ${host} -P ${port} -D '${database}' --ssl-mode=VERIFY_IDENTITY --ssl-ca=${ca_path} -p'${password}'
`},
	{Endpoint: "private", Client: "mysqlcli", Type: "code", Path: "./private/mysql.sh", DownloadCa: []string{}, Doc: "", Content: `mysql --comments -u '${username}' -h ${host} -P ${port} -D '${database}' -p'${password}'
`},
	{Endpoint: "public", Client: "mariadbcli", Type: "code", Path: "./public/mariadb.sh", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "", Content: `mysql --comments -u '${username}' -h ${host} -P ${port} -D '${database}' --ssl-verify-server-cert --ssl-ca=${ca_path} -p'${password}'
`},
	{Endpoint: "private", Client: "mariadbcli", Type: "code", Path: "./private/mariadb.sh", DownloadCa: []string{}, Doc: "", Content: `mysql --comments -u '${username}' -h ${host} -P ${port} -D '${database}' -p'${password}'
`},
	{Endpoint: "public", Client: "mysqlconnectorj", Type: "code", Path: "./public/java/mysql", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-java-jdbc", Content: `jdbc:mysql://${username}:${password}@${host}:${port}/${database}?sslMode=VERIFY_IDENTITY
`},
	{Endpoint: "private", Client: "mysqlconnectorj", Type: "code", Path: "./private/java/mysql", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-java-jdbc", Content: `jdbc:mysql://${username}:${password}@${host}:${port}/${database}
`},
	{Endpoint: "public", Client: "mariadbconnectorj", Type: "code", Path: "./public/java/mariadb", DownloadCa: []string{}, Doc: "", Content: `jdbc:mariadb://${host}:${port}/${database}?user=${username}&password=${password}&sslMode=verify-full
`},
	{Endpoint: "private", Client: "mariadbconnectorj", Type: "code", Path: "./private/java/mariadb", DownloadCa: []string{}, Doc: "", Content: `jdbc:mariadb://${host}:${port}/${database}?user=${username}&password=${password}
`},
	{Endpoint: "public", Client: "pythonmysqlclient", Type: "code", Path: "./public/python/mysqlclient.py", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-mysqlclient", Content: `import MySQLdb

connection = MySQLdb.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
  ssl_mode = "VERIFY_IDENTITY",
  ssl = {
    "ca": "${ca_path}"
  }
)
`},
	{Endpoint: "private", Client: "pythonmysqlclient", Type: "code", Path: "./private/python/mysqlclient.py", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-mysqlclient", Content: `import MySQLdb

connection = MySQLdb.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
)
`},
	{Endpoint: "public", Client: "pythonpymysql", Type: "code", Path: "./public/python/pymysql.py", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-pymysql", Content: `import pymysql

connection = pymysql.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
  ssl_verify_cert = True,
  ssl_verify_identity = True,
  ssl_ca = "${ca_path}"
)
`},
	{Endpoint: "private", Client: "pythonpymysql", Type: "code", Path: "./private/python/pymysql.py", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-pymysql", Content: `import pymysql

connection = pymysql.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
)
`},
	{Endpoint: "public", Client: "pythonmysqlconnector", Type: "code", Path: "./public/python/connector.py", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-mysql-connector", Content: `import mysql.connector

connection = mysql.connector.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
  ssl_ca = "${ca_path}",
  ssl_verify_cert = True,
  ssl_verify_identity = True
)
`},
	{Endpoint: "private", Client: "pythonmysqlconnector", Type: "code", Path: "./private/python/connector.py", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-mysql-connector", Content: `import mysql.connector

connection = mysql.connector.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
)
`},
	{Endpoint: "public", Client: "sqlalchemymysqlclient", Type: "code", Path: "./public/sqlalchemy/mysqlclient", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-sqlalchemy", Content: `mysql+mysqldb://${username}:${password}@${host}:${port}/${database}?ssl_mode=VERIFY_IDENTITY&ssl_ca=${ca_path}
`},
	{Endpoint: "private", Client: "sqlalchemymysqlclient", Type: "code", Path: "./private/sqlalchemy/mysqlclient", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-sqlalchemy", Content: `mysql+mysqldb://${username}:${password}@${host}:${port}/${database}
`},
	{Endpoint: "public", Client: "sqlalchemypymysql", Type: "code", Path: "./public/sqlalchemy/pymysql", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-sqlalchemy", Content: `mysql+pymysql://${username}:${password}@${host}:${port}/${database}?ssl_ca=${ca_path}&ssl_verify_cert=true&ssl_verify_identity=true
`},
	{Endpoint: "private", Client: "sqlalchemypymysql", Type: "code", Path: "./private/sqlalchemy/pymysql", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-sqlalchemy", Content: `mysql+pymysql://${username}:${password}@${host}:${port}/${database}
`},
	{Endpoint: "public", Client: "sqlalchemymysqlconnector", Type: "code", Path: "./public/sqlalchemy/connector", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-sqlalchemy", Content: `mysql+mysqlconnector://${username}:${password}@${host}:${port}/${database}?ssl_ca=${ca_path}&ssl_verify_cert=true&ssl_verify_identity=true
`},
	{Endpoint: "private", Client: "sqlalchemymysqlconnector", Type: "code", Path: "./private/sqlalchemy/connector", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-sqlalchemy", Content: `mysql+mysqlconnector://${username}:${password}@${host}:${port}/${database}
`},
	{Endpoint: "public", Client: "django", Type: "code", Path: "./public/django.py", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-django", Content: `DATABASES = {
    'default': {
        'ENGINE': 'django_tidb',
        'NAME': '${database}',
        'USER': '${username}',
        'PASSWORD': '${password}',
        'HOST': '${host}',
        'PORT': ${port},
        'OPTIONS': {
            'ssl_mode': 'VERIFY_IDENTITY',
            'ssl': {'ca': '${ca_path}'}
        }
    },
}
`},
	{Endpoint: "private", Client: "django", Type: "code", Path: "./private/django.py", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-python-django", Content: `DATABASES = {
    'default': {
        'ENGINE': 'django_tidb',
        'NAME': '${database}',
        'USER': '${username}',
        'PASSWORD': '${password}',
        'HOST': '${host}',
        'PORT': ${port},
    },
}
`},
	{Endpoint: "public", Client: "nodemysql2", Type: "code", Path: "./public/js/mysql2", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-nodejs-mysql2", Content: `mysql://${username}:${password}@${host}:${port}/${database}?ssl={"rejectUnauthorized":true}
`},
	{Endpoint: "private", Client: "nodemysql2", Type: "code", Path: "./private/js/mysql2", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-nodejs-mysql2", Content: `mysql://${username}:${password}@${host}:${port}/${database}
`},
	{Endpoint: "public", Client: "prisma", Type: "code", Path: "./public/js/prisma", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-nodejs-prisma", Content: `mysql://${username}:${password}@${host}:${port}/${database}?sslaccept=strict
`},
	{Endpoint: "private", Client: "prisma", Type: "code", Path: "./private/js/prisma", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-nodejs-prisma", Content: `mysql://${username}:${password}@${host}:${port}/${database}
`},
	{Endpoint: "public", Client: "serverlessjs", Type: "code", Path: "./public/js/serverless", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/serverless-driver-node-example", Content: `mysql://${username}:${password}@${host}/${database}
`},
	{Endpoint: "private", Client: "serverlessjs", Type: "code", Path: "./private/js/serverless", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/serverless-driver-node-example", Content: `mysql://${username}:${password}@${host}/${database}
`},
	{Endpoint: "public", Client: "go", Type: "code", Path: "./public/mysql.go", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-golang-sql-driver", Content: `mysql.RegisterTLSConfig("tidb", &tls.Config{
  MinVersion: tls.VersionTLS12,
  ServerName: "${host}",
})

db, err := sql.Open("mysql", "${username}:${password}@tcp(${host}:${port})/${database}?tls=tidb")
`},
	{Endpoint: "private", Client: "go", Type: "code", Path: "./private/mysql.go", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-golang-sql-driver", Content: `db, err := sql.Open("mysql", "${username}:${password}@tcp(${host}:${port})/${database}")
`},
	{Endpoint: "public", Client: "rails", Type: "code", Path: "./public/rails.yml", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-ruby-rails", Content: `development:
  <<: *default
  adapter: mysql2
  url: mysql2://${username}:${password}@${host}:${port}/${database}?ssl_mode=verify_identity  
`},
	{Endpoint: "private", Client: "rails", Type: "code", Path: "./private/rails.yml", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-sample-application-ruby-rails", Content: `development:
  <<: *default
  adapter: mysql2
  url: mysql2://${username}:${password}@${host}:${port}/${database}  
`},
	{Endpoint: "public", Client: "datagrip", Type: "code", Path: "./public/gui/datagrip", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-datagrip", Content: `jdbc:mysql://${host}:${port}/${database}?user=${username}&password=${password}&sslMode=VERIFY_IDENTITY
`},
	{Endpoint: "private", Client: "datagrip", Type: "code", Path: "./private/gui/datagrip", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-datagrip", Content: `jdbc:mysql://${host}:${port}/${database}?user=${username}&password=${password}
`},
	{Endpoint: "public", Client: "dbeaver", Type: "code", Path: "./public/java/mysql", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-dbeaver", Content: `jdbc:mysql://${username}:${password}@${host}:${port}/${database}?sslMode=VERIFY_IDENTITY
`},
	{Endpoint: "private", Client: "dbeaver", Type: "code", Path: "./private/java/mysql", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-dbeaver", Content: `jdbc:mysql://${username}:${password}@${host}:${port}/${database}
`},
	{Endpoint: "public", Client: "vscode", Type: "parameter", Path: "./public/general/params", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-vscode-sqltools", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
CA ${ca_path}
`},
	{Endpoint: "private", Client: "vscode", Type: "parameter", Path: "./private/general/params", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-vscode-sqltools", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
`},
	{Endpoint: "public", Client: "mysqlworkbench", Type: "parameter", Path: "./public/general/params", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-mysql-workbench", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
CA ${ca_path}
`},
	{Endpoint: "private", Client: "mysqlworkbench", Type: "parameter", Path: "./private/general/params", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-mysql-workbench", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
`},
	{Endpoint: "public", Client: "navicat", Type: "parameter", Path: "./public/general/params", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-navicat", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
CA ${ca_path}
`},
	{Endpoint: "private", Client: "navicat", Type: "parameter", Path: "./private/general/params", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-gui-navicat", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
`},
	{Endpoint: "public", Client: "wordpress", Type: "parameter", Path: "./public/general/params", DownloadCa: []string{
		"windows",
		"other",
	}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-wordpress", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
CA ${ca_path}
`},
	{Endpoint: "private", Client: "wordpress", Type: "parameter", Path: "./private/general/params", DownloadCa: []string{}, Doc: "https://docs.pingcap.com/tidbcloud/dev-guide-wordpress", Content: `HOST ${host}
PORT ${port}
USERNAME ${username}
PASSWORD ${password}
DATABASE ${database}
`},
}
