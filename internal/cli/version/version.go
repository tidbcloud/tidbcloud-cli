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

package version

import (
	"fmt"
	"regexp"
	"strings"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/version"

	"github.com/spf13/cobra"
)

func VersionCmd(h *internal.Helper) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Args:   cobra.NoArgs,
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(h.IOStreams.Out, Format(version.Version, version.Commit, version.Date))
		},
	}

	return cmd
}

// Format formats a version string with the given information.
func Format(ver, commit, buildDate string) string {
	fmt.Println("ver:", ver)
	if ver == version.DevVersion && buildDate == "" && commit == "" {
		return fmt.Sprintf("%s version (built from source)", config.CliName)
	}

	if strings.Contains(ver, version.NightlyVersion) {
		return fmt.Sprintf("%s version %s-%s (commit: %s)\n", config.CliName, ver, buildDate, commit)
	}

	ver = strings.TrimPrefix(ver, "v")

	return fmt.Sprintf("%s version %s (build date: %s commit: %s)\n%s\n", config.CliName, ver, buildDate, commit, changelogURL(ver))
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
