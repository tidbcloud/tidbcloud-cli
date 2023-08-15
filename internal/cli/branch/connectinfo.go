// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package branch

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"tidbcloud-cli/internal"
	clu "tidbcloud-cli/internal/cli/cluster"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/util"
	branchApi "tidbcloud-cli/pkg/tidbcloud/branch/client/branch_service"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type connectInfoOpts struct {
	interactive bool
}

func (c connectInfoOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.BranchID,
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
		Short: "Get connection string for the specified branch",
		Example: fmt.Sprintf(`  Get connection string in interactive mode:
 $ %[1]s branch connect-info

 Get connection string in non-interactive mode:
 $ %[1]s branch connect-info --cluster-id <cluster-id> --branch-id <branch-id> --client <client-name> --operating-system <operating-system>
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
			var branchID, clusterID, client, operatingSystem string

			// Get TiDBCloudClient
			d, err := h.Client()
			if err != nil {
				return err
			}

			if opts.interactive { // interactive mode
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// Get cluster id
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				// Get branch id
				branch, err := cloud.GetSelectedBranch(clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				branchID = branch.ID

				// Get client
				clientNameForInteractive, err := cloud.GetSelectedConnectClient(clu.ConnectClientsList)
				if err != nil {
					return err
				}
				client = clu.ClientsForInteractiveMap[clientNameForInteractive]

				// Detect operating system
				// TODO: detect linux operating system name
				goOS := runtime.GOOS
				if goOS == "darwin" {
					goOS = "macOS"
				} else if goOS == "windows" {
					goOS = "Windows"
				}
				if goOS != "" && goOS != "linux" {
					for id, value := range clu.OperatingSystemList {
						if strings.Contains(value, goOS) {
							operatingSystemValueWithFlag := clu.OperatingSystemList[id] + " (Detected)"
							clu.OperatingSystemList = append([]string{operatingSystemValueWithFlag},
								append(clu.OperatingSystemList[:id], clu.OperatingSystemList[id+1:]...)...)
							break
						}
					}
				}

				// Get operating system
				operatingSystemCombination, err := cloud.GetSelectedConnectOs(clu.OperatingSystemList)
				if err != nil {
					return err
				}
				operatingSystem = strings.Split(operatingSystemCombination, "/")[0]

			} else { // non-interactive mode
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return err
				}

				branchID, err = cmd.Flags().GetString(flag.BranchID)
				if err != nil {
					return err
				}

				clientNameForHelp, err := cmd.Flags().GetString(flag.ClientName)
				if err != nil {
					return err
				}
				if v, ok := clu.ClientsForHelpMap[strings.ToLower(clientNameForHelp)]; ok {
					client = v
				} else {
					return errors.New(fmt.Sprintf("Unsupported client. Run \"%[1]s cluster connect-info -h\" to check supported clients list", config.CliName))
				}

				operatingSystem, err = cmd.Flags().GetString(flag.OperatingSystem)
				if err != nil {
					return err
				}
				if !clu.Contains(operatingSystem, clu.OperatingSystemListForHelp) {
					return errors.New(fmt.Sprintf("Unsupported operating system. Run \"%[1]s cluster connect-info -h\" to check supported operating systems list", config.CliName))
				}
			}

			// Get branch info
			params := branchApi.NewGetBranchParams().WithBranchID(branchID).WithClusterID(clusterID)
			branchInfo, err := d.GetBranch(params)
			if err != nil {
				return err
			}

			// Get connect parameter
			defaultUser := fmt.Sprintf("%s.root", *branchInfo.Payload.UserPrefix)
			host := branchInfo.Payload.Endpoints.PublicEndpoint.Host
			port := strconv.Itoa(int(branchInfo.Payload.Endpoints.PublicEndpoint.Port))

			// Get connection string
			connectInfo, err := cloud.RetrieveConnectInfo(d)
			if err != nil {
				return err
			}
			connectionString, err := util.GenerateConnectionString(connectInfo, client, host, defaultUser, port, clu.SERVERLESS, operatingSystem, util.Shell)
			if err != nil {
				return err
			}
			fmt.Fprintln(h.IOStreams.Out)
			fmt.Fprintln(h.IOStreams.Out, connectionString)

			return nil
		},
	}

	cmd.Flags().StringP(flag.BranchID, flag.BranchIDShort, "", "Branch ID")
	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	cmd.Flags().String(flag.ClientName, "", fmt.Sprintf("Connected client. Supported clients: %q", clu.ConnectClientsListForHelp))
	cmd.Flags().String(flag.OperatingSystem, "", fmt.Sprintf("Operating system name. "+
		"Supported operating systems: %q", clu.OperatingSystemListForHelp))

	return cmd
}
