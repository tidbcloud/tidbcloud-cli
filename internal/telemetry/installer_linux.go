//go:build !windows && !darwin
// +build !windows,!darwin

package telemetry

import (
	"os"

	"tidbcloud-cli/internal/config"
)

func readInstaller() *string {
	if config.IsUnderTiUP {
		s := "TiUP"
		return &s
	}

	if b, err := os.ReadFile("/etc/ticloud/installer"); err == nil {
		s := string(b)
		return &s
	}

	return nil
}
