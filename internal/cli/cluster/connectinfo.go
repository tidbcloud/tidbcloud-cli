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

	"github.com/fatih/color"
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

// Display clients name orderly
var connectClientsList = []string{
	// pure parameter
	"Standard Connection Parameter",

	// CLI
	"MySQL CLI",
	"MyCLI",

	// driver
	"libmysqlclient",
	"mysqlclient",
	"PyMySQL",
	"MySQL Connector/Python",
	"MySQL Connector/J",
	"Go MySQL Driver",
	"Node MySQL 2",
	"Mysql2 (Ruby)",
	"MySQLi",
	"mysql (Rust)",

	// ORM
	"MyBatis",
	"Hibernate",
	"Spring Boot",
	"GORM",
	"Prisma",
	"Sequelize",
	"Django",
	"SQLAlchemy",
	"Active Record",
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
		Short: "Get connection string for your specified cluster",
		Example: fmt.Sprintf(`  Get connection string in interactive mode:
  $ %[1]s connect-info

  Get connection string in non-interactive mode:
  $ %[1]s connect-info --project-id <project-id> --cluster-id <cluster-id> --client '<client-name>' --operating-system <operating-system>
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
			var projectID, clusterID, clientName, operatingSystem string

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
				clientName, err = cloud.GetSelectedConnectClient(connectClientsList)
				if err != nil {
					return err
				}

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

				clientName, err = cmd.Flags().GetString(flag.ClientName)
				if err != nil {
					return err
				}
				if !contains(clientName, connectClientsList) {
					return errors.New("Unsupported clients. Run \"ticloud cluster connect-info -h\" to check supported clients list")
				}

				operatingSystem, err = cmd.Flags().GetString(flag.OperatingSystem)
				if err != nil {
					return err
				}
				if !contains(operatingSystem, operatingSystemListForHelp) {
					return errors.New("Unsupported operating system. Run \"ticloud cluster connect-info -h\" to check supported operating systems list")
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
			connectionString, err := generateConnectionString(connectInfo, clientName, host, defaultUser, port, clusterType, operatingSystem)
			if err != nil {
				return err
			}
			fmt.Fprintln(h.IOStreams.Out, color.BlueString("Connection string"))
			fmt.Fprintln(h.IOStreams.Out, connectionString)

			return nil
		},
	}

	cmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "Project ID")
	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	cmd.Flags().String(flag.ClientName, "", "Connected client. Supported clients: ["+strings.Join(connectClientsList, ", ")+"]")
	cmd.Flags().String(flag.OperatingSystem, "", "Operating system name. Supported operating systems: ["+strings.Join(operatingSystemListForHelp, ", ")+"]")

	return cmd
}

func generateConnectionString(connectInfo *connectInfoModel.ConnectInfo, clientName string, host string, user string, port string, clusterType string, operatingSystem string) (string, error) {
	for _, clientData := range connectInfo.ClientData {
		if strings.EqualFold(clientData.DisplayName, clientName) {
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
