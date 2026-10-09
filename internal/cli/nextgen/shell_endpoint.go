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
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

const privateEndpointLookupTimeout = 30 * time.Second

func parseEndpointAddress(address string) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil || host == "" || strings.ContainsAny(host, "/\\?#@ \t\r\n") {
		return "", fmt.Errorf("--endpoint must be host:port, without a URL scheme or path")
	}
	number, err := strconv.ParseUint(port, 10, 16)
	if err != nil || number == 0 {
		return "", fmt.Errorf("--endpoint port must be between 1 and 65535")
	}
	return net.JoinHostPort(host, strconv.FormatUint(number, 10)), nil
}

// loadPrivateEndpoints reads private networking APIs, independently of the
// instance detail's endpoint reachability. Service readiness is not proof of
// connectivity from the machine running the CLI.
func loadPrivateEndpoints(ctx context.Context, client cloud.NextGenClient, tidbID string, provider api.RegionCloudProvider, pageSize int32) ([]api.TidbEndpoint, error) {
	// One budget covers all pages and the optional AWS service lookup. Cancel it
	// before prompting or starting the separate SQL connection timeout.
	ctx, cancel := context.WithTimeout(ctx, privateEndpointLookupTimeout)
	defer cancel()

	var endpoints []api.TidbEndpoint
	seenAddresses := make(map[string]bool)
	add := func(host, port string) {
		address, err := parseEndpointAddress(net.JoinHostPort(host, port))
		if err != nil || seenAddresses[address] {
			return
		}
		_, normalizedPort, _ := net.SplitHostPort(address)
		number, _ := strconv.ParseInt(normalizedPort, 10, 32)
		endpoints = append(endpoints, api.TidbEndpoint{
			Host: api.PtrString(host), Port: api.PtrInt32(int32(number)),
			ConnectionType: api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT.Ptr(),
		})
		seenAddresses[address] = true
	}

	var token string
	seenTokens := make(map[string]bool)
	hasConnections := false
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		page, err := client.ListPrivateEndpointConnections(ctx, tidbID, pageSize, token)
		if err != nil {
			return nil, fmt.Errorf("list private endpoint connections: %w", err)
		}
		if page == nil {
			return nil, fmt.Errorf("the API returned an empty private endpoint connection response")
		}
		for _, connection := range page.GetPrivateEndpointConnections() {
			hasConnections = true
			if connection.GetEndpointState() == api.NEXTGENV1BETA2PRIVATEENDPOINTCONNECTIONENDPOINTSTATE_ACTIVE &&
				connection.GetPrivateLinkServiceState() == api.NEXTGENV1BETA2PRIVATELINKSERVICESTATE_ACTIVE {
				add(connection.GetHost(), connection.GetPort())
			}
		}
		token = page.GetNextPageToken()
		if token == "" {
			break
		}
		if seenTokens[token] {
			return nil, fmt.Errorf("private endpoint connections returned a repeated page token")
		}
		seenTokens[token] = true
	}

	// AWS publishes a service-wide hostname before a consumer connection exists.
	// Only an empty successful list can use this path; failures and unready
	// existing connections must not be bypassed via a different privilege.
	if !hasConnections && provider == api.REGIONCLOUDPROVIDER_AWS {
		service, err := client.GetPrivateLinkService(ctx, tidbID)
		if err != nil {
			return nil, fmt.Errorf("get private link service: %w", err)
		}
		if service != nil && service.GetCloudProvider() == api.REGIONCLOUDPROVIDER_AWS && service.GetState() == api.NEXTGENV1BETA2PRIVATELINKSERVICESTATE_ACTIVE {
			add(service.GetServiceDnsName(), strconv.FormatInt(int64(service.GetServicePort()), 10))
		}
	}
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no ready PRIVATE_ENDPOINT address is available; configure a private endpoint connection and wait until it is ACTIVE")
	}
	return endpoints, nil
}

// resolvePrivateEndpoint uses an explicit address without checking discovery
// results, or selects an address validated by loadPrivateEndpoints.
// A nil selector keeps the explicit-instance path free of additional prompts.
func resolvePrivateEndpoint(endpoints []api.TidbEndpoint, address string, selectEndpoint func([]string) (int, error)) (*api.TidbEndpoint, error) {
	if address != "" {
		normalized, err := parseEndpointAddress(address)
		if err != nil {
			return nil, err
		}
		host, port, _ := net.SplitHostPort(normalized)
		number, _ := strconv.ParseInt(port, 10, 32)
		return &api.TidbEndpoint{
			Host: api.PtrString(host), Port: api.PtrInt32(int32(number)),
			ConnectionType: api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT.Ptr(),
		}, nil
	}

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("instance has no PRIVATE_ENDPOINT endpoint; configure a private endpoint connection first")
	}
	selected := 0
	if len(endpoints) > 1 {
		addresses := make([]string, 0, len(endpoints))
		for _, endpoint := range endpoints {
			addresses = append(addresses, net.JoinHostPort(endpoint.GetHost(), strconv.Itoa(int(endpoint.GetPort()))))
		}
		if selectEndpoint == nil {
			return nil, fmt.Errorf("multiple PRIVATE_ENDPOINT addresses are available; specify --endpoint with one of: %s", strings.Join(addresses, ", "))
		}
		var err error
		selected, err = selectEndpoint(addresses)
		if err != nil {
			return nil, err
		}
	}
	return validateEndpoint(&endpoints[selected])
}

func promptPrivateEndpoint(addresses []string) (int, error) {
	var selected int
	err := survey.AskOne(&survey.Select{Message: "Choose the private endpoint:", Options: addresses}, &selected)
	if err == terminal.InterruptErr {
		return 0, util.InterruptError
	}
	return selected, err
}
