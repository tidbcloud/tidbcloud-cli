// Copyright 2026 PingCAP, Inc.
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

package nextgen

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"
)

func privateEndpointCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	cmd := &cobra.Command{Use: "private-endpoint", Short: "Manage private endpoint connections", Args: cobra.NoArgs,
		Long: "Bind an endpoint already created in your cloud account to an instance. Cloud endpoint creation, DNS and network configuration remain your responsibility. Availability depends on the service plan, cloud and region.",
		RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() }}
	for _, action := range []string{"service", "list", "bind", "describe", "unbind"} {
		cmd.AddCommand(privateEndpointActionCmd(h, plan, action, confirmUnbindPrivateEndpoint))
	}
	return cmd
}

func privateEndpointActionCmd(h *internal.Helper, plan planSpec, action string, confirm func(string, string) error) *cobra.Command {
	var id, endpointID, connectionID, domain, ip string
	var force bool
	descriptions := map[string]string{
		"service":  "Get the private link service used to create an endpoint in your cloud account",
		"list":     "List all private endpoint connections and their states",
		"bind":     "Bind an existing cloud endpoint; acceptance does not mean the connection is ACTIVE",
		"describe": "Describe a private endpoint connection",
		"unbind":   "Unbind a connection; this interrupts traffic and does not delete the endpoint in your cloud account",
	}
	cmd := &cobra.Command{Use: action, Short: descriptions[action], Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			format, err := validatedOutputFormat(cmd)
			if err != nil {
				return err
			}
			if strings.TrimSpace(id) == "" && cmd.Flags().Changed(flag.ClusterID) {
				return fmt.Errorf("--%s must not be empty", flag.ClusterID)
			}
			if id == "" && !h.IOStreams.CanPrompt {
				return fmt.Errorf("the terminal doesn't support prompts; use --%s to specify the instance", flag.ClusterID)
			}
			if action == "bind" && strings.TrimSpace(endpointID) == "" {
				return fmt.Errorf("--endpoint-id is required; create the endpoint in your cloud account first")
			}
			if action == "describe" || action == "unbind" {
				if strings.TrimSpace(connectionID) == "" || strings.ContainsAny(connectionID, "/\\?#") || connectionID == "." || connectionID == ".." {
					return fmt.Errorf("--connection-id must be a connection ID returned by private-endpoint list, not a cloud endpoint ID or resource path")
				}
			}
			if action == "unbind" && !force && !h.IOStreams.CanPrompt {
				return fmt.Errorf("unbinding interrupts private connections; use --force to confirm in a non-interactive terminal")
			}
			if id != "" {
				if _, err := strconv.ParseUint(id, 10, 64); err != nil {
					return fmt.Errorf("--%s must be a numeric instance ID", flag.ClusterID)
				}
			}
			if cmd.Flags().Changed("private-ip-address") && (net.ParseIP(ip) == nil || net.ParseIP(ip).To4() == nil) {
				return fmt.Errorf("--private-ip-address must be an IPv4 address")
			}
			if cmd.Flags().Changed("private-domain-name") && strings.TrimSpace(domain) == "" {
				return fmt.Errorf("--private-domain-name must not be empty")
			}
			client, err := h.NextGenClient()
			if err != nil {
				return err
			}
			if id == "" {
				instance, err := selectInstance(cmd.Context(), client, plan, int32(h.QueryPageSize))
				if err != nil {
					return err
				}
				id = instance.GetTidbId()
			}
			instance, err := client.GetTiDB(cmd.Context(), id)
			if err != nil {
				return err
			}
			if err := plan.validate(instance); err != nil {
				return err
			}
			var result interface{}
			switch action {
			case "service":
				result, err = client.GetPrivateLinkService(cmd.Context(), id)
			case "list":
				all := &api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse{}
				all.PrivateEndpointConnections = []api.Nextgenv1beta2PrivateEndpointConnection{}
				token := ""
				seen := map[string]bool{}
				for {
					page, e := client.ListPrivateEndpointConnections(cmd.Context(), id, int32(h.QueryPageSize), token)
					if e != nil {
						return e
					}
					if page == nil {
						return fmt.Errorf("the API returned an empty private endpoint connection response")
					}
					all.PrivateEndpointConnections = append(all.PrivateEndpointConnections, page.GetPrivateEndpointConnections()...)
					token = page.GetNextPageToken()
					if token == "" {
						break
					}
					if seen[token] {
						return fmt.Errorf("the API returned a repeated private endpoint pagination token")
					}
					seen[token] = true
				}
				all.SetTotalSize(int32(len(all.PrivateEndpointConnections)))
				result = all
			case "bind":
				body := api.NewNextgenv1beta2PrivateEndpointConnection(instance.GetRegionId(), endpointID)
				if cmd.Flags().Changed("private-domain-name") {
					body.SetPrivateDomainName(domain)
				}
				if cmd.Flags().Changed("private-ip-address") {
					body.SetPrivateIpAddress(ip)
				}
				result, err = client.CreatePrivateEndpointConnection(cmd.Context(), id, body)
			case "describe", "unbind":
				connection, e := client.GetPrivateEndpointConnection(cmd.Context(), id, connectionID)
				if e != nil {
					return e
				}
				// Do not mutate a connection belonging to another instance, even if the server returns it.
				if connection == nil || connection.GetTidbId() != id {
					return fmt.Errorf("the returned private endpoint connection does not belong to the selected instance")
				}
				result = connection
				if action == "unbind" {
					if !force {
						if e := confirm(id, connectionID); e != nil {
							if errors.Is(e, terminal.InterruptErr) {
								return util.InterruptError
							}
							return e
						}
					}
					result, err = client.DeletePrivateEndpointConnection(cmd.Context(), id, connectionID)
				}
			}
			if err != nil {
				return err
			}
			if api.IsNil(result) {
				return fmt.Errorf("the API returned an empty private endpoint response")
			}
			if format == output.HumanFormat {
				if action == "bind" || action == "unbind" {
					fmt.Fprintln(h.IOStreams.Out, "Request accepted; use private-endpoint list or describe to check completion before retrying.")
				}

				switch value := result.(type) {
				case *api.Nextgenv1beta2PrivateLinkService:
					return output.PrintHumanTable(h.IOStreams.Out, []output.Column{"Service", "State", "DNS", "Port", "Available zones"}, []output.Row{{value.GetServiceName(), string(value.GetState()), value.GetServiceDnsName(), fmt.Sprint(value.GetServicePort()), strings.Join(value.GetAvailableZones(), ", ")}})
				case *api.Nextgenv1beta2PrivateEndpointConnection:
					return output.PrintHumanTable(h.IOStreams.Out, privateEndpointColumns, []output.Row{privateEndpointRow(*value)})
				case *api.Nextgenv1beta2ListPrivateEndpointConnectionsResponse:
					rows := make([]output.Row, 0, len(value.PrivateEndpointConnections))
					for _, connection := range value.PrivateEndpointConnections {
						rows = append(rows, privateEndpointRow(connection))
					}
					return output.PrintHumanTable(h.IOStreams.Out, privateEndpointColumns, rows)
				}
			}
			return output.PrintJson(h.IOStreams.Out, result)
		}}
	cmd.Example = fmt.Sprintf("  ticloud %s private-endpoint %s -c <instance-id>", plan.commandName, action)
	if action == "bind" {
		cmd.Example += " --endpoint-id <cloud-endpoint-id>"
	}
	if action == "describe" || action == "unbind" {
		cmd.Example += " --connection-id <connection-id>"
	}
	cmd.Flags().StringVarP(&id, flag.ClusterID, flag.ClusterIDShort, "", "Instance ID. If omitted, select interactively.")
	cmd.Flags().StringP(flag.Output, flag.OutputShort, output.JsonFormat, flag.OutputHelp)
	if action == "bind" {
		cmd.Flags().StringVar(&endpointID, "endpoint-id", "", "Existing cloud endpoint ID (AWS: vpce-...).")
		cmd.Flags().StringVar(&domain, "private-domain-name", "", "Private endpoint domain name (Alibaba Cloud only).")
		cmd.Flags().StringVar(&ip, "private-ip-address", "", "Private endpoint IPv4 address (Azure only).")
		cmd.Long = descriptions[action] + ". Use list or describe to wait for ACTIVE before connecting. No automatic retry is performed."
	}
	if action == "describe" || action == "unbind" {
		cmd.Flags().StringVar(&connectionID, "connection-id", "", "TiDB Cloud connection ID returned by private-endpoint list.")
	}
	if action == "unbind" {
		cmd.Flags().BoolVar(&force, flag.Force, false, "Unbind without confirmation. Query list or describe to check asynchronous completion.")
	}
	return cmd
}

func confirmUnbindPrivateEndpoint(id, connectionID string) error {
	var answer string
	if err := survey.AskOne(&survey.Input{Message: fmt.Sprintf("Unbinding connection %s from instance %s interrupts private traffic. Type yes to confirm:", connectionID, id)}, &answer); err != nil {
		return err
	}
	if answer != "yes" {
		return fmt.Errorf("incorrect confirmation string; skipping private endpoint unbinding")
	}
	return nil
}

var privateEndpointColumns = []output.Column{"Connection ID", "Endpoint ID", "State", "Host", "Port"}

func privateEndpointRow(connection api.Nextgenv1beta2PrivateEndpointConnection) output.Row {
	return output.Row{connection.GetPrivateEndpointConnectionId(), connection.GetEndpointId(), string(connection.GetEndpointState()), connection.GetHost(), connection.GetPort()}
}
