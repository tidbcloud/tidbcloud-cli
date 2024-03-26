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

package serverless

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/util"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
)

type connectInfoOpts struct {
	interactive bool
}

func (c connectInfoOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ClientName,
		flag.OperatingSystem,
	}
}

func ConnectInfoCmd(h *internal.Helper) *cobra.Command {
	opts := connectInfoOpts{
		interactive: true,
	}

	cmd := &cobra.Command{
		Use:   "connect-info",
		Short: "Get connection string for a serverless cluster",
		Example: fmt.Sprintf(`  Get connection string in interactive mode:
  $ %[1]s serverless connect-info

  Get connection string in non-interactive mode:
  $ %[1]s serverless connect-info --cluster-id <cluster-id> --client <client-name> --operating-system <operating-system>`, config.CliName),
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
			var clusterID, clientID, operatingSystemID string

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
				projectID := project.ID

				// Get cluster id
				cluster, err := cloud.GetSelectedCluster(projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				// Detect operating system
				// TODO: detect linux operating system name
				goOS := runtime.GOOS
				if goOS == "darwin" {
					goOS = "macos"
				} else if goOS == "windows" {
					goOS = "windows"
				}
				println("Detected operating system:", goOS)
				operatingSystems := make([]util.Os, len(util.ConnectInfoOs))
				for i, os := range util.ConnectInfoOs {
					operatingSystems[i] = util.Os{
						ID:   os.ID,
						Name: os.Name,
					}
					if strings.Contains(os.ID, goOS) {
						operatingSystems[i].Name = operatingSystems[i].Name + " (Detected)"
						operatingSystems[0], operatingSystems[i] = operatingSystems[i], operatingSystems[0]
					}
				}

				// Get operating system
				operatingSystem, err := cloud.GetSelectedConnectOs(operatingSystems)
				if err != nil {
					return err
				}
				operatingSystemID = operatingSystem.ID

				// Get client
				client, err := cloud.GetSelectedConnectClient(util.ConnectInfoClient)
				if err != nil {
					return err
				}
				clientID = client.ID
				if client.Options != nil && len(client.Options) > 0 {
					// get options
					option, err := cloud.GetSelectedConnectClientOptions(client.Options)
					if err != nil {
						return err
					}
					clientID = option.ID
				}
			} else { // non-interactive mode
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return err
				}

				clientName, err := cmd.Flags().GetString(flag.ClientName)
				if err != nil {
					return err
				}
				for _, v := range util.ConnectInfoClient {
					name := v.Name
					if v.Options != nil && len(v.Options) > 0 {
						for _, option := range v.Options {
							name = fmt.Sprintf("%s_%s", v.Name, option.Name)
						}
					}
					if clientName == name {
						clientID = v.ID
						break
					}
				}
				if clientID == "" {
					return errors.New(fmt.Sprintf("Unsupported client. Run \"%[1]s serverless connect-info -h\" to check supported clients list", config.CliName))
				}

				operatingSystemName, err := cmd.Flags().GetString(flag.OperatingSystem)
				if err != nil {
					return err
				}
				for _, v := range util.ConnectInfoOs {
					osInfos := strings.Split(v.Name, "/")
					for _, osInfo := range osInfos {
						if operatingSystemName == osInfo {
							operatingSystemID = v.ID
							break
						}
					}
				}
				if operatingSystemID == "" {
					return errors.New(fmt.Sprintf("Unsupported operating system. Run \"%[1]s serverless connect-info -h\" to check supported operating systems list", config.CliName))
				}
			}

			// get cluster info
			params := serverlessApi.NewServerlessServiceGetClusterParams().WithClusterID(clusterID)
			clusterInfo, err := d.GetCluster(params)
			if err != nil {
				return err
			}
			defaultUser := fmt.Sprintf("%s.root", clusterInfo.Payload.UserPrefix)
			host := clusterInfo.Payload.Endpoints.PublicEndpoint.Host
			port := strconv.Itoa(int(clusterInfo.Payload.Endpoints.PublicEndpoint.Port))

			// get connect string
			connectString, err := cloud.GetConnectString(clientID, operatingSystemID, defaultUser, host, port, "test")
			if err != nil {
				return err
			}

			fmt.Fprintln(h.IOStreams.Out)
			fmt.Fprintln(h.IOStreams.Out, connectString)
			return nil
		},
	}

	var ConnectClientName []string
	for _, v := range util.ConnectInfoClient {
		if v.Options != nil && len(v.Options) > 0 {
			for _, option := range v.Options {
				ConnectClientName = append(ConnectClientName, fmt.Sprintf("%s(%s)", v.Name, option.Name))
			}
		} else {
			ConnectClientName = append(ConnectClientName, v.Name)
		}
	}

	var ConnectOsName []string
	for _, v := range util.ConnectInfoOs {
		osInfos := strings.Split(v.Name, "/")
		ConnectOsName = append(ConnectOsName, osInfos...)
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the Cluster")
	cmd.Flags().String(flag.ClientName, "", fmt.Sprintf("Connected client. Supported clients: %q", ConnectClientName))
	cmd.Flags().String(flag.OperatingSystem, "", fmt.Sprintf("Operating system name. "+
		"Supported operating systems: %q", ConnectOsName))

	return cmd
}
