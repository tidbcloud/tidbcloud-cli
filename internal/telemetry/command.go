package telemetry

import (
	"tidbcloud-cli/internal/config"

	"github.com/pingcap/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type TrackOptions struct {
	Err error
}

var currentTracker *tracker

func StartTrackingCommand(cmd *cobra.Command, args []string) {
	if !config.TelemetryEnabled() {
		return
	}

	var err error
	currentTracker, err = newTracker(cmd, args)
	if err != nil {
		log.Debug("telemetry: failed to create tracker", zap.Error(err))
		return
	}
}

func FinishTrackingCommand(opt TrackOptions) {
	if !config.TelemetryEnabled() || currentTracker == nil {
		return
	}

	if err := currentTracker.trackCommand(opt); err != nil {
		log.Debug("telemetry: failed to track command", zap.Error(err))
	}
}
