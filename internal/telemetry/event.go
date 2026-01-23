// Copyright 2026 PingCAP, Inc.
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
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/version"

	"github.com/mattn/go-isatty"
	"github.com/pingcap/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	InteractiveMode = "interactive"
	ProjectID       = "project-id"
	ClusterID       = "cluster-id"
)

type Event struct {
	Timestamp  time.Time              `json:"timestamp"`
	Source     string                 `json:"source"`
	Properties map[string]interface{} `json:"properties"`
}

type eventOpt func(Event)

func WithInteractive(cmd *cobra.Command) eventOpt {
	return func(event Event) {
		if cmd.Annotations[InteractiveMode] == "true" {
			event.Properties[InteractiveMode] = true
		} else {
			event.Properties[InteractiveMode] = false
		}
	}
}

func withCommandPath(cmd *cobra.Command) eventOpt {
	return func(event Event) {
		cmdPath := cmd.CommandPath()
		event.Properties["command"] = strings.ReplaceAll(cmdPath, " ", "-")
		if cmd.CalledAs() != "" {
			event.Properties["alias"] = cmd.CalledAs()
		}
	}
}

func withDuration(cmd *cobra.Command) eventOpt {
	return func(event Event) {
		if cmd.Context() == nil {
			log.Debug("telemetry: context not found")
			return
		}

		ctxValue, found := cmd.Context().Value(contextKey).(telemetryContextValue)
		if !found {
			log.Debug("telemetry: context not found")
			return
		}

		event.Properties["duration"] = event.Timestamp.Sub(ctxValue.startTime).Milliseconds()
	}
}

func withFlags(cmd *cobra.Command) eventOpt {
	return func(event Event) {
		setFlags := make([]string, 0, cmd.Flags().NFlag())
		cmd.Flags().Visit(func(f *pflag.Flag) {
			setFlags = append(setFlags, f.Name)
		})

		if len(setFlags) > 0 {
			event.Properties["flags"] = setFlags
		}
	}
}

func withVersion() eventOpt {
	return func(event Event) {
		event.Properties["version"] = version.Version
		event.Properties["git_commit"] = version.Commit
	}
}

func withOS() eventOpt {
	return func(event Event) {
		event.Properties["os"] = runtime.GOOS
		event.Properties["arch"] = runtime.GOARCH
	}
}

func withAuthMethod() eventOpt {
	return func(event Event) {
		event.Properties["auth_method"] = "api_key"
	}
}

func withProjectID(cmd *cobra.Command) eventOpt {
	return func(event Event) {
		if id, found := cmd.Annotations[ProjectID]; found {
			event.Properties["project_id"] = id
			return
		}
	}
}

func withTerminal() eventOpt {
	return func(event Event) {
		if isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			event.Properties["terminal"] = "cygwin"
		}

		if isatty.IsTerminal(os.Stdout.Fd()) {
			event.Properties["terminal"] = "tty"
			return
		}
	}
}

func withInstaller(installer *string) eventOpt {
	return func(event Event) {
		if installer != nil {
			event.Properties["installer"] = *installer
		}
	}
}

func withError(err error) eventOpt {
	return func(event Event) {
		event.Properties["result"] = "ERROR"

		errorMessage := strings.Split(err.Error(), "\n")[0] // only first line

		event.Properties["error"] = errorMessage
	}
}

func newEvent(opts ...eventOpt) Event {
	var event = Event{
		Timestamp: time.Now(),
		Source:    config.CliName,
		Properties: map[string]interface{}{
			"result": "SUCCESS",
		},
	}

	for _, fn := range opts {
		fn(event)
	}

	return event
}
