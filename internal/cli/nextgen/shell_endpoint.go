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
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func parseEndpointAddress(address string) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil || host == "" || strings.ContainsAny(host, "/\\?#@ \t\r\n") {
		return "", fmt.Errorf("--endpoint must be an API-returned host:port, without a URL scheme or path")
	}
	number, err := strconv.ParseUint(port, 10, 16)
	if err != nil || number == 0 {
		return "", fmt.Errorf("--endpoint port must be between 1 and 65535")
	}
	return net.JoinHostPort(host, strconv.FormatUint(number, 10)), nil
}

// resolvePrivateEndpoint only selects addresses published for this instance.
// An API reachability value of true is not proof of client-side connectivity.
// A nil selector keeps the explicit-instance path free of additional prompts.
func resolvePrivateEndpoint(instance *api.Nextgenv1beta2Tidb, address string, selectEndpoint func([]string) (int, error)) (*api.TidbEndpoint, error) {
	candidates := make([]*api.TidbEndpoint, 0, len(instance.Endpoints))
	addresses := make([]string, 0, len(instance.Endpoints))
	found := false
	for i := range instance.Endpoints {
		endpoint := &instance.Endpoints[i]
		if endpoint.GetConnectionType() != api.ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT {
			continue
		}
		found = true
		// Some providers publish the endpoint before its address is available.
		if endpoint.GetHost() == "" || endpoint.GetPort() <= 0 || endpoint.GetPort() > 65535 {
			continue
		}
		candidates = append(candidates, endpoint)
		addresses = append(addresses, net.JoinHostPort(endpoint.GetHost(), strconv.Itoa(int(endpoint.GetPort()))))
	}
	if !found {
		return nil, fmt.Errorf("instance has no PRIVATE_ENDPOINT endpoint; configure a private endpoint connection first")
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("PRIVATE_ENDPOINT endpoint is not ready: the API has not returned a usable host and port")
	}
	if address != "" {
		for i, candidateAddress := range addresses {
			if candidateAddress == address {
				return validateEndpoint(candidates[i])
			}
		}
		return nil, fmt.Errorf("--endpoint %q is not a PRIVATE_ENDPOINT address returned for this instance; available addresses: %s", address, strings.Join(addresses, ", "))
	}
	selected := 0
	if len(candidates) > 1 {
		if selectEndpoint == nil {
			return nil, fmt.Errorf("multiple PRIVATE_ENDPOINT addresses are available; specify --endpoint with one of: %s", strings.Join(addresses, ", "))
		}
		var err error
		selected, err = selectEndpoint(addresses)
		if err != nil {
			return nil, err
		}
	}
	return validateEndpoint(candidates[selected])
}

func promptPrivateEndpoint(addresses []string) (int, error) {
	var selected int
	err := survey.AskOne(&survey.Select{Message: "Choose the private endpoint:", Options: addresses}, &selected)
	if err == terminal.InterruptErr {
		return 0, util.InterruptError
	}
	return selected, err
}
