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

package config

import (
	"os"

	"github.com/tidbcloud/tidbcloud-cli/internal/util"

	"github.com/charmbracelet/lipgloss"
)

const (
	DebugEnv      = "TICLOUD_DEBUG"
	cliName       = "ticloud"
	cliNameInTiUP = "cloud"
	HomePath      = ".ticloud"

	Confirmed     = "yes"
	OAuthEndpoint = "https://oauth.tidbcloud.com"
	ClientID      = "wiSBy7f27zWBaBCxS16tDm7DDj2T3POgwFFbefTrgx8FAXKhzaPzv1Uta9NTck2r"
	ClientSecret  = "ieKdWDao0QFmHbfYXpQJHuYZ9nLdpRptfE9d3W30WettIFLZL66JKewznvVkY39IkEbBQiZK60pjcnm7BN7Lj9uRiCSpIC4n2aI3IUyHlLKtvxZLrrfXugHsC7qhb1Js"
)

var (
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	CursorStyle  = FocusedStyle.Copy()
	CliName      = cliName
	IsUnderTiUP  = false
)

func init() {
	var binpath string
	if exepath, err := os.Executable(); err == nil {
		binpath = exepath
	}

	IsUnderTiUP = util.IsUnderTiUP(binpath)

	if IsUnderTiUP {
		CliName = cliNameInTiUP
	} else {
		CliName = cliName
	}
}
