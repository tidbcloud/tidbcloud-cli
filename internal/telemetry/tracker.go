// Copyright 2025 PingCAP, Inc.
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
	"github.com/spf13/cobra"
)

type tracker struct {
	cmd       *cobra.Command
	sender    EventsSender
	args      []string
	installer *string
}

func newTracker(cmd *cobra.Command, args []string) (*tracker, error) {
	return &tracker{
		cmd:       cmd,
		sender:    NewSender(),
		args:      args,
		installer: readInstaller(),
	}, nil
}

func (t *tracker) defaultCommandOptions() []eventOpt {
	return []eventOpt{withCommandPath(t.cmd), WithInteractive(t.cmd), withFlags(t.cmd), withVersion(), withOS(), withAuthMethod(), withProjectID(t.cmd), withTerminal(), withInstaller(t.installer)}
}

func (t *tracker) trackCommand(data TrackOptions) error {
	options := append(t.defaultCommandOptions(), withDuration(t.cmd))

	if data.Err != nil {
		options = append(options, withError(data.Err))
	}
	event := newEvent(options...)
	err := t.sender.SendEvent(event)
	if err != nil {
		return err
	}
	return nil
}
