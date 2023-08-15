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
	"tidbcloud-cli/internal/util"

	serverlessApi "tidbcloud-cli/pkg/tidbcloud/serverless/client/serverless_service"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
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
var ConnectClientsList = []string{
	// pure parameter
	util.GeneralParameterDisplayName,

	// CLI
	util.MysqlCliDisplayName,
	util.MyCliDisplayName,

	// driver
	util.LibMysqlClientDisplayName,
	util.MysqlClientDisplayName,
	util.PyMysqlDisplayName,
	util.MysqlConnectorPythonDisplayName,
	util.MysqlConnectorJavaDisplayName,
	util.GoMysqlDriverDisplayName,
	util.NodeMysql2DisplayName,
	util.Mysql2RubyDisplayName,
	util.MysqliDisplayName,
	util.MysqlRustDisplayName,

	// ORM
	util.MybatisDisplayName,
	util.HibernateDisplayName,
	util.SpringBootDisplayName,
	util.GormDisplayName,
	util.PrismaDisplayName,
	util.SequelizeDisplayName,
	util.DjangoDisplayName,
	util.SqlAlchemyDisplayName,
	util.ActiveRecordDisplayName,
}

// Display clients name orderly in help message
var ConnectClientsListForHelp = []string{
	// pure parameter
	util.GeneralParameterInputName,

	// CLI
	util.MysqlCliInputName,
	util.MyCliInputName,

	// driver
	util.LibMysqlClientInputName,
	util.MysqlClientInputName,
	util.PyMysqlInputName,
	util.MysqlConnectorPythonInputName,
	util.MysqlConnectorJavaInputName,
	util.GoMysqlDriverInputName,
	util.NodeMysql2InputName,
	util.Mysql2RubyInputName,
	util.MysqliInputName,
	util.MysqlRustInputName,

	// ORM
	util.MybatisInputName,
	util.HibernateInputName,
	util.SpringBootInputName,
	util.GormInputName,
	util.PrismaInputName,
	util.SequelizeInputName,
	util.DjangoInputName,
	util.SqlAlchemyInputName,
	util.ActiveRecordInputName,
}

var ClientsForInteractiveMap = map[string]string{
	// pure parameter
	util.GeneralParameterDisplayName: util.GeneralParameterID,

	// CLI
	util.MysqlCliDisplayName: util.MysqlCliID,
	util.MyCliDisplayName:    util.MyCliID,

	// driver
	util.LibMysqlClientDisplayName:       util.LibMysqlClientID,
	util.MysqlClientDisplayName:          util.MysqlClientID,
	util.PyMysqlDisplayName:              util.PyMysqlID,
	util.MysqlConnectorPythonDisplayName: util.MysqlConnectorPythonID,
	util.MysqlConnectorJavaDisplayName:   util.MysqlConnectorJavaID,
	util.GoMysqlDriverDisplayName:        util.GoMysqlDriverID,
	util.NodeMysql2DisplayName:           util.NodeMysql2ID,
	util.Mysql2RubyDisplayName:           util.Mysql2RubyID,
	util.MysqliDisplayName:               util.MysqliID,
	util.MysqlRustDisplayName:            util.MysqlRustID,

	// ORM
	util.MybatisDisplayName:      util.MybatisID,
	util.HibernateDisplayName:    util.HibernateID,
	util.SpringBootDisplayName:   util.SpringBootID,
	util.GormDisplayName:         util.GormID,
	util.PrismaDisplayName:       util.PrismaID,
	util.SequelizeDisplayName:    util.SequelizeID,
	util.DjangoDisplayName:       util.DjangoID,
	util.SqlAlchemyDisplayName:   util.SQLAlchemyID,
	util.ActiveRecordDisplayName: util.ActiveRecordID,
}

var ClientsForHelpMap = map[string]string{
	// pure parameter
	util.GeneralParameterInputName: util.GeneralParameterID,

	// CLI
	util.MysqlCliInputName: util.MysqlCliID,
	util.MyCliInputName:    util.MyCliID,

	// driver
	util.LibMysqlClientInputName:       util.LibMysqlClientID,
	util.MysqlClientInputName:          util.MysqlClientID,
	util.PyMysqlInputName:              util.PyMysqlID,
	util.MysqlConnectorPythonInputName: util.MysqlConnectorPythonID,
	util.MysqlConnectorJavaInputName:   util.MysqlConnectorJavaID,
	util.GoMysqlDriverInputName:        util.GoMysqlDriverID,
	util.NodeMysql2InputName:           util.NodeMysql2ID,
	util.Mysql2RubyInputName:           util.Mysql2RubyID,
	util.MysqliInputName:               util.MysqliID,
	util.MysqlRustInputName:            util.MysqlRustID,

	// ORM
	util.MybatisInputName:      util.MybatisID,
	util.HibernateInputName:    util.HibernateID,
	util.SpringBootInputName:   util.SpringBootID,
	util.GormInputName:         util.GormID,
	util.PrismaInputName:       util.PrismaID,
	util.SequelizeInputName:    util.SequelizeID,
	util.DjangoInputName:       util.DjangoID,
	util.SqlAlchemyInputName:   util.SQLAlchemyID,
	util.ActiveRecordInputName: util.ActiveRecordID,
}

// Display operating system orderly in interactive mode
var OperatingSystemList = []string{
	"macOS/Alpine",
	"CentOS/RedHat/Fedora",
	"Debian/Ubuntu/Arch",
	"Windows",
	"OpenSUSE",
	"Others",
}

// Display operating system orderly in help message
var OperatingSystemListForHelp = []string{
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
				clientNameForInteractive, err := cloud.GetSelectedConnectClient(ConnectClientsList)
				if err != nil {
					return err
				}
				client = ClientsForInteractiveMap[clientNameForInteractive]

				// Detect operating system
				// TODO: detect linux operating system name
				goOS := runtime.GOOS
				if goOS == "darwin" {
					goOS = "macOS"
				} else if goOS == "windows" {
					goOS = "Windows"
				}
				if goOS != "" && goOS != "linux" {
					for id, value := range OperatingSystemList {
						if strings.Contains(value, goOS) {
							operatingSystemValueWithFlag := OperatingSystemList[id] + " (Detected)"
							OperatingSystemList = append([]string{operatingSystemValueWithFlag},
								append(OperatingSystemList[:id], OperatingSystemList[id+1:]...)...)
							break
						}
					}
				}

				// Get operating system
				operatingSystemCombination, err := cloud.GetSelectedConnectOs(OperatingSystemList)
				if err != nil {
					return err
				}
				operatingSystem = strings.Split(operatingSystemCombination, "/")[0]

			} else { // non-interactive mode
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return err
				}

				_, err = cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return err
				}

				clientNameForHelp, err := cmd.Flags().GetString(flag.ClientName)
				if err != nil {
					return err
				}
				if v, ok := ClientsForHelpMap[strings.ToLower(clientNameForHelp)]; ok {
					client = v
				} else {
					return errors.New(fmt.Sprintf("Unsupported client. Run \"%[1]s cluster connect-info -h\" to check supported clients list", config.CliName))
				}

				operatingSystem, err = cmd.Flags().GetString(flag.OperatingSystem)
				if err != nil {
					return err
				}
				if !Contains(operatingSystem, OperatingSystemListForHelp) {
					return errors.New(fmt.Sprintf("Unsupported operating system. Run \"%[1]s cluster connect-info -h\" to check supported operating systems list", config.CliName))
				}
			}

			// Get cluster info
			params := serverlessApi.NewServerlessServiceGetClusterParams().WithClusterID(clusterID)
			clusterInfo, err := d.GetCluster(params)
			if err != nil {
				return err
			}

			// Resolve cluster information
			// Get connect parameter
			defaultUser := fmt.Sprintf("%s.root", clusterInfo.Payload.UserPrefix)
			host := clusterInfo.Payload.Endpoints.PublicEndpoint.Host
			port := strconv.Itoa(int(clusterInfo.Payload.Endpoints.PublicEndpoint.Port))
			clusterType := SERVERLESS

			// Get connection string
			connectInfo, err := cloud.RetrieveConnectInfo(d)
			if err != nil {
				return err
			}
			connectionString, err := util.GenerateConnectionString(connectInfo, client, host, defaultUser, port, clusterType, operatingSystem, util.Shell)
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
	cmd.Flags().String(flag.ClientName, "", fmt.Sprintf("Connected client. Supported clients: %q", ConnectClientsListForHelp))
	cmd.Flags().String(flag.OperatingSystem, "", fmt.Sprintf("Operating system name. "+
		"Supported operating systems: %q", OperatingSystemListForHelp))

	return cmd
}

func Contains(str string, vec []string) bool {
	for _, v := range vec {
		if strings.EqualFold(str, v) {
			return true
		}
	}
	return false
}
