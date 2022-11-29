// Copyright 2022 PingCAP, Inc.
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

package version

import (
	"fmt"
	"regexp"
	"strings"

	"tidbcloud-cli/internal"

	"github.com/spf13/cobra"
)

func VersionCmd(h *internal.Helper, ver, commit, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(h.IOStreams.Out, Format(ver, commit, buildDate))
		},
	}

	return cmd
}

// Format formats a version string with the given information.
func Format(ver, commit, buildDate string) string {
	if ver == "" && buildDate == "" && commit == "" {
		return "pscale version (built from source)"
	}

	ver = strings.TrimPrefix(ver, "v")

	return fmt.Sprintf("pscale version %s (build date: %s commit: %s)\n%s\n", ver, buildDate, commit, changelogURL(ver))
}

func changelogURL(version string) string {
	path := "https://github.com/tidbcloud/tidbcloud-cli"
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", path)
	}

	url := fmt.Sprintf("%s/releases/tag/v%s", path, strings.TrimPrefix(version, "v"))
	return url
}
