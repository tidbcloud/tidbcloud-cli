// Copyright 2024 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
