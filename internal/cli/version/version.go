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
