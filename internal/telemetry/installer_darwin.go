//go:build darwin
// +build darwin

package telemetry

import "tidbcloud-cli/internal/config"

func readInstaller() *string {
	if config.IsUnderTiUP {
		s := "TiUP"
		return &s
	}

	return nil
}
