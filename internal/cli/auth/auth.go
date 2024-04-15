// Copyright 2024 PingCAP, Inc.
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

package auth

import (
	"tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

const (
	tokenTypeHint = "access_token"
	grantType     = "urn:ietf:params:oauth:grant-type:device_code"

	authPath   = "/v1/device_authorization"
	accessPath = "/v1/token"
	revokePath = "/v1/revoke"

	errSlowDown = "slow_down"
	errPending  = "authorization_pending"
)

func AuthCmd(h *internal.Helper) *cobra.Command {
	var authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Login and logout via TiDB Cloud API",
	}

	authCmd.AddCommand(LoginCmd(h))
	authCmd.AddCommand(LogoutCmd(h))
	return authCmd
}
