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
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	_ "github.com/xo/usql/drivers/mysql"
	"github.com/xo/usql/env"
)

const maxCACertificateSize = 1 << 20

const privateEndpointConnectTimeout = 30 * time.Second

const publicDNSConnectTimeout = 30 * time.Second

func shellCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	return shellCmdWithSelector(h, plan, selectInstance)
}

type instanceSelector func(context.Context, cloud.NextGenClient, planSpec, int32) (*api.Nextgenv1beta2Tidb, error)

func shellCmdWithSelector(h *internal.Helper, plan planSpec, selectInstance instanceSelector) *cobra.Command {
	var connection, address string
	var connectionType api.EndpointConnectionType
	command := &cobra.Command{
		Use:   "shell",
		Short: fmt.Sprintf("Connect to a %s instance through a public or private endpoint", plan.displayName),
		Long: "Connect through a public endpoint by default, or an existing private endpoint with --connection-type private-endpoint. " +
			"With --endpoint host:port, connect directly without discovering or checking the address against API results. " +
			"Otherwise, discovering private addresses requires private endpoint read permissions and has a 30-second total timeout. " +
			"Private networking and DNS must already be configured. This command does not create network resources. " +
			"Private connections have a 30-second initial connection timeout and retain TLS certificate verification; there is no fallback to public connections. " +
			"For a known public DNS probe failure, the CLI checks that public access is enabled, then warns and attempts a TLS-verified connection with a 30-second initial timeout after password entry.",
		Args: cobra.NoArgs,
		Example: fmt.Sprintf(`  $ %[1]s %[2]s shell
  $ %[1]s %[2]s shell -c <instance-id>
  $ %[1]s %[2]s shell -c <instance-id> --password <password>
  $ %[1]s %[2]s shell -c <instance-id> -u <user-name> --password <password>
  $ %[1]s %[2]s shell --connection-type private-endpoint
  $ %[1]s %[2]s shell -c <instance-id> --connection-type private-endpoint
  $ %[1]s %[2]s shell -c <instance-id> --connection-type private-endpoint --endpoint <host>:4000`, config.CliName, plan.commandName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			switch connection {
			case "public":
				connectionType = api.ENDPOINTCONNECTIONTYPE_PUBLIC
			case "private-endpoint":
				connectionType = api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT
			default:
				return fmt.Errorf("unsupported connection type %q; supported values are public and private-endpoint", connection)
			}
			if cmd.Flags().Changed("endpoint") {
				if connectionType != api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT {
					return fmt.Errorf("--endpoint requires --connection-type private-endpoint")
				}
				var err error
				address, err = parseEndpointAddress(address)
				if err != nil {
					return err
				}
			}
			if shellIsInteractive(cmd) {
				return nil
			}
			id, err := cmd.Flags().GetString(flag.ClusterID)
			if err != nil {
				return err
			}
			if id == "" {
				return fmt.Errorf("required flag(s) %q not set", flag.ClusterID)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !h.IOStreams.CanPrompt {
				return fmt.Errorf("the stdout is not a terminal")
			}
			client, err := h.NextGenClient()
			if err != nil {
				return err
			}

			interactive := shellIsInteractive(cmd)
			var id string
			if interactive {
				selected, err := selectInstance(cmd.Context(), client, plan, int32(h.QueryPageSize))
				if err != nil {
					return err
				}
				id = selected.GetTidbId()
			} else {
				id, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return err
				}
			}

			instance, err := client.GetTiDB(cmd.Context(), id)
			if err != nil {
				return err
			}
			if err := plan.validate(instance); err != nil {
				return err
			}
			if instance.GetState() != api.V1BETA1CLUSTERSTATE_ACTIVE {
				return fmt.Errorf("instance %q is in state %q; wait until it is ACTIVE before connecting", id, instance.GetState())
			}

			var endpoint *api.TidbEndpoint
			var connectTimeout time.Duration
			var dnsWarning bool
			if connectionType == api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT {
				var endpoints []api.TidbEndpoint
				if address == "" {
					endpoints, err = loadPrivateEndpoints(cmd.Context(), client, id, instance.GetCloudProvider(), int32(h.QueryPageSize))
					if err != nil {
						return err
					}
				}
				var selectEndpoint func([]string) (int, error)
				if interactive {
					selectEndpoint = promptPrivateEndpoint
				}
				endpoint, err = resolvePrivateEndpoint(endpoints, address, selectEndpoint)
				connectTimeout = privateEndpointConnectTimeout
			} else {
				endpoint, dnsWarning, err = resolvePublicEndpoint(cmd.Context(), client, id, instance)
			}
			if err != nil {
				return err
			}
			if dnsWarning {
				fmt.Fprintf(h.IOStreams.Err, "Warning: the API reports %s for PUBLIC endpoint %q; public access is enabled. Attempting a TLS-verified connection with a 30-second initial timeout after password entry.\n", endpoint.ConnectionReachability.Detail.GetMessage(), endpoint.GetHost())
				connectTimeout = publicDNSConnectTimeout
			}
			caURL, err := client.GetCACertificateDownloadURL(cmd.Context(), id)
			if err != nil {
				return err
			}
			if caURL == nil || caURL.GetUri() == "" {
				return fmt.Errorf("the API returned an empty CA certificate download URL")
			}
			roots, err := downloadCACertificate(cmd.Context(), caURL.GetUri())
			if err != nil {
				return err
			}

			var userName string
			if interactive {
				userName, err = promptShellUser()
				if err != nil {
					return err
				}
			} else {
				userName, err = cmd.Flags().GetString(flag.User)
				if err != nil {
					return err
				}
			}
			if userName == "" {
				userName = "root"
				fmt.Fprintln(h.IOStreams.Out, color.GreenString("Current user: ")+color.HiGreenString(userName))
			}
			var password *string
			if !interactive && cmd.Flags().Changed(flag.Password) {
				value, err := cmd.Flags().GetString(flag.Password)
				if err != nil {
					return err
				}
				password = &value
			}

			if err := env.Set("PROMPT1", "%n@"+instance.DisplayName+"%R%#"); err != nil {
				return err
			}
			tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12, ServerName: endpoint.GetHost(), RootCAs: roots}
			err = util.ExecuteSqlDialogWithTLS(cmd.Context(), userName, endpoint.GetHost(), strconv.Itoa(int(endpoint.GetPort())), password, tlsConfig, connectTimeout)
			return annotateShellTLSError(connectionType, endpoint.GetHost(), err)
		},
	}
	command.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the instance.")
	command.Flags().StringP(flag.User, flag.UserShort, "", "The SQL user name. The default is root.")
	command.Flags().String(flag.Password, "", "The password of the SQL user.")
	command.Flags().StringVar(&connection, "connection-type", "public", "The connection type: public or private-endpoint. Does not configure networking.")
	command.Flags().StringVar(&address, "endpoint", "", "Connect directly to host:port without private endpoint discovery or existence checks. Requires --connection-type private-endpoint; TLS certificate verification remains enabled.")
	return command
}

func annotateShellTLSError(connectionType api.EndpointConnectionType, host string, err error) error {
	if err == nil {
		return nil
	}
	var unknownAuthority x509.UnknownAuthorityError
	if errors.As(err, &unknownAuthority) {
		return fmt.Errorf("the CA certificate returned by the API does not verify %s endpoint %q: %w", connectionType, host, err)
	}
	return err
}

func shellIsInteractive(command *cobra.Command) bool {
	for _, name := range []string{flag.ClusterID, flag.User, flag.Password} {
		if command.Flags().Changed(name) {
			return false
		}
	}
	return true
}

type instanceChoice struct {
	instance *api.Nextgenv1beta2Tidb
}

func (c instanceChoice) String() string {
	instance := c.instance
	result := fmt.Sprintf("%s(%s)[%s][%s]", instance.GetDisplayName(), instance.GetTidbId(), instance.GetState(), instance.GetRegionId())
	if projectID := instance.GetLabels()[projectIDLabel]; projectID != "" {
		result += fmt.Sprintf("[project:%s]", projectID)
	}
	return result
}

func selectInstance(ctx context.Context, client cloud.NextGenClient, plan planSpec, pageSize int32) (*api.Nextgenv1beta2Tidb, error) {
	instances, err := retrieveTiDBs(ctx, client, plan, pageSize)
	if err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		return nil, fmt.Errorf("no available %s instances found", plan.displayName)
	}

	items := make([]interface{}, 0, len(instances))
	for i := range instances {
		items = append(items, &instanceChoice{instance: &instances[i]})
	}
	model, err := ui.InitialSelectModel(items, "Choose the instance:")
	if err != nil {
		return nil, err
	}
	model.EnablePagination(6)
	model.EnableFilter()

	result, err := tea.NewProgram(model).Run()
	if err != nil {
		return nil, err
	}
	selection := result.(ui.SelectModel)
	if selection.Interrupted {
		return nil, util.InterruptError
	}
	selected := selection.GetSelectedItem()
	if selected == nil {
		return nil, fmt.Errorf("no instance selected")
	}
	return selected.(*instanceChoice).instance, nil
}

func promptShellUser() (string, error) {
	useDefaultUser := false
	if err := survey.AskOne(&survey.Confirm{Message: "Use the default user?", Default: true}, &useDefaultUser); err != nil {
		if err == terminal.InterruptErr {
			return "", util.InterruptError
		}
		return "", err
	}
	if useDefaultUser {
		return "", nil
	}

	var userName string
	if err := survey.AskOne(&survey.Input{Message: "Please input the user name:"}, &userName, survey.WithValidator(survey.Required)); err != nil {
		if err == terminal.InterruptErr {
			return "", util.InterruptError
		}
		return "", err
	}
	return userName, nil
}

// resolvePublicEndpoint permits only known DNS probe failures to reach the SQL
// driver. The setting is a configuration snapshot, not proof of connectivity.
func resolvePublicEndpoint(ctx context.Context, client cloud.NextGenClient, id string, instance *api.Nextgenv1beta2Tidb) (*api.TidbEndpoint, bool, error) {
	for i := range instance.Endpoints {
		endpoint := &instance.Endpoints[i]
		if endpoint.GetConnectionType() != api.ENDPOINTCONNECTIONTYPE_PUBLIC {
			continue
		}
		if !isKnownPublicDNSProbeFailure(endpoint) {
			selected, err := validateEndpoint(endpoint)
			return selected, false, err
		}
		address := net.JoinHostPort(endpoint.GetHost(), strconv.Itoa(int(endpoint.GetPort())))
		if _, err := parseEndpointAddress(address); err != nil {
			return nil, false, fmt.Errorf("PUBLIC endpoint is not ready: the API has not returned a usable host and port")
		}
		settingCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		setting, err := client.GetPublicConnectionSetting(settingCtx, id)
		if err != nil {
			return nil, false, fmt.Errorf("cannot confirm PUBLIC access is enabled: %w", err)
		}
		enabled, _ := setting.GetEnabledOk()
		if enabled == nil {
			return nil, false, fmt.Errorf("cannot confirm PUBLIC access is enabled: the API returned no enabled value")
		}
		if !*enabled {
			return nil, false, fmt.Errorf("PUBLIC access is disabled")
		}
		return endpoint, true, nil
	}
	return nil, false, fmt.Errorf("instance has no PUBLIC endpoint")
}

func isKnownPublicDNSProbeFailure(endpoint *api.TidbEndpoint) bool {
	if endpoint.GetConnectionType() != api.ENDPOINTCONNECTIONTYPE_PUBLIC {
		return false
	}
	r := endpoint.ConnectionReachability
	if r == nil || r.Reachable == nil || r.GetReachable() || r.Detail == nil {
		return false
	}
	d := r.Detail
	if !d.GetServiceActive() || !d.GetEndpointActive() || d.DnsReachable == nil || d.GetDnsReachable() {
		return false
	}
	// Message is human-readable in the API schema. Match only these established
	// server values; new or absent reasons must retain the existing rejection.
	return d.GetMessage() == "DNS_NOT_REACHABLE" || d.GetMessage() == "DNS_PROBE_FAILED"
}

func validateEndpoint(endpoint *api.TidbEndpoint) (*api.TidbEndpoint, error) {
	connectionType := endpoint.GetConnectionType()
	if endpoint.ConnectionReachability != nil && endpoint.ConnectionReachability.Reachable != nil && !endpoint.ConnectionReachability.GetReachable() {
		detail := endpoint.ConnectionReachability.GetDetail()
		if message := detail.GetMessage(); message != "" {
			return nil, fmt.Errorf("%s endpoint is not reachable: %s", connectionType, message)
		}
		return nil, fmt.Errorf("%s endpoint is not reachable", connectionType)
	}
	if endpoint.GetHost() == "" || endpoint.GetPort() <= 0 {
		return nil, fmt.Errorf("%s endpoint is not ready", connectionType)
	}
	return endpoint, nil
}

func downloadCACertificate(ctx context.Context, uri string) (*x509.CertPool, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	if request.URL.Scheme != "https" {
		return nil, fmt.Errorf("the CA certificate download URL must use HTTPS")
	}
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(next *http.Request, via []*http.Request) error {
			if next.URL.Scheme != "https" {
				return fmt.Errorf("the CA certificate download redirect must use HTTPS")
			}
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects while downloading the CA certificate")
			}
			return nil
		},
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("download CA certificate: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("download CA certificate: unexpected HTTP status %s", response.Status)
	}
	pem, err := io.ReadAll(io.LimitReader(response.Body, maxCACertificateSize+1))
	if err != nil {
		return nil, fmt.Errorf("read CA certificate: %w", err)
	}
	if len(pem) > maxCACertificateSize {
		return nil, fmt.Errorf("CA certificate exceeds %d bytes", maxCACertificateSize)
	}
	roots, err := x509.SystemCertPool()
	if err != nil || roots == nil {
		roots = x509.NewCertPool()
	}
	if !roots.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("downloaded CA certificate is not valid PEM")
	}
	return roots, nil
}
