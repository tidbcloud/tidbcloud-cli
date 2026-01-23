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

package util

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
)

func ConvertToAuthorizedNetwork(startIP string, endIP string, displayName string) (cluster.EndpointsPublicAuthorizedNetwork, error) {
	s, err := parseIPv4(startIP)
	if err != nil {
		return cluster.EndpointsPublicAuthorizedNetwork{}, err
	}
	e, err := parseIPv4(endIP)
	if err != nil {
		return cluster.EndpointsPublicAuthorizedNetwork{}, err
	}

	if bytes.Compare(s, e) == 1 {
		return cluster.EndpointsPublicAuthorizedNetwork{}, fmt.Errorf("start IP address %s must not be greater than end IP address %s", startIP, endIP)
	}

	if displayName == "" {
		displayName = GenerateIDAuthorizedNetworkDisplayName()
	}

	return cluster.EndpointsPublicAuthorizedNetwork{
		StartIpAddress: startIP,
		EndIpAddress:   endIP,
		DisplayName:    displayName,
	}, nil
}

func parseIPv4(ipStr string) (net.IP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To4() == nil {
		return nil, fmt.Errorf("invalid IPv4 address: %s", ipStr)
	}

	return ip.To4(), nil
}

func GenerateIDAuthorizedNetworkDisplayName() string {
	now := time.Now()
	timeStr := now.Format("2006_0102_1504")
	micros := now.Nanosecond() / 1e3
	return fmt.Sprintf("Allowlist_%s_%d", timeStr, micros)
}
