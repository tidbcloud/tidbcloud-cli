// Copyright 2025 PingCAP, Inc.
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
	"encoding/base32"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
)

func ConvertToAuthorizedNetworks(authorizedNetworksStrList []string) ([]cluster.EndpointsPublicAuthorizedNetwork, error) {
	authorizedNetworks := make([]cluster.EndpointsPublicAuthorizedNetwork, 0, len(authorizedNetworksStrList))
	for _, str := range authorizedNetworksStrList {
		authorizedNetwork, err := ConvertToAuthorizedNetwork(str, "")
		if err != nil {
			return nil, fmt.Errorf("failed to convert authorized network %s: %w", str, err)
		}
		authorizedNetworks = append(authorizedNetworks, authorizedNetwork)
	}
	return authorizedNetworks, nil
}

func ConvertToAuthorizedNetwork(authorizedNetworksStr string, displayName string) (cluster.EndpointsPublicAuthorizedNetwork, error) {
	startIP, endIP, err := parseIPv4Range(authorizedNetworksStr)
	if err != nil {
		return cluster.EndpointsPublicAuthorizedNetwork{}, err
	}

	if displayName == "" {
		displayName = GenerateIDAuthorizedNetworkDisplayName()
	}

	return cluster.EndpointsPublicAuthorizedNetwork{
		StartIpAddress: startIP.String(),
		EndIpAddress:   endIP.String(),
		DisplayName:    displayName,
	}, nil
}

func parseIPv4Range(ipRange string) (net.IP, net.IP, error) {
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid IP range format, expected '{start-ip}-{end-ip}'")
	}

	startIPStr := strings.TrimSpace(parts[0])
	endIPStr := strings.TrimSpace(parts[1])

	startIP := net.ParseIP(startIPStr)
	if startIP == nil || startIP.To4() == nil {
		return nil, nil, fmt.Errorf("invalid IPv4 start address: %s", startIPStr)
	}

	endIP := net.ParseIP(endIPStr)
	if endIP == nil || endIP.To4() == nil {
		return nil, nil, fmt.Errorf("invalid IPv4 end address: %s", endIPStr)
	}

	return startIP.To4(), endIP.To4(), nil
}

var base32Encoder = base32.StdEncoding.WithPadding(base32.NoPadding)

func GenerateIDAuthorizedNetworkDisplayName() string {
	uuidV4 := uuid.New()
	now := time.Now()
	dateStr := now.Format("20060102")
	return "Allowlist_" + dateStr + "_" + strings.ToLower(base32Encoder.EncodeToString(uuidV4[:]))
}
