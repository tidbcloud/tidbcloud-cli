// Copyright 2026 PingCAP, Inc.
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

package changefeed

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cdc"
)

type startPositionMode string

const (
	startPositionFromNow  startPositionMode = "Current Time"
	startPositionFromTSO  startPositionMode = "Input a TSO(example: 443852055297916932)"
	startPositionFromTIME startPositionMode = "Input a Time(example: 2024-01-01T00:00:00Z)"
)

var startPositionModes = []startPositionMode{
	startPositionFromNow,
	startPositionFromTSO,
	startPositionFromTIME,
}

func GetSelectedStartPositionMode() (startPositionMode, error) {
	modes := make([]interface{}, 0, len(startPositionModes))
	for _, v := range startPositionModes {
		modes = append(modes, v)
	}
	model, err := ui.InitialSelectModel(modes, "Input the start time:")
	if err != nil {
		return "", errors.Trace(err)
	}

	p := tea.NewProgram(model)
	targetTypeModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := targetTypeModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	result := targetTypeModel.(ui.SelectModel).GetSelectedItem()
	if result == nil {
		return "", errors.New("no start position mode selected")
	}
	return result.(startPositionMode), nil
}

func GetSelectedChangefeedType() (cdc.ChangefeedTypeEnum, error) {
	changefeedTypes := make([]interface{}, 0, len(cdc.AllowedChangefeedTypeEnumEnumValues))
	for _, v := range cdc.AllowedChangefeedTypeEnumEnumValues {
		changefeedTypes = append(changefeedTypes, v)
	}
	model, err := ui.InitialSelectModel(changefeedTypes, "Choose the changefeed type:")
	if err != nil {
		return "", errors.Trace(err)
	}

	p := tea.NewProgram(model)
	targetTypeModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := targetTypeModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	changefeedType := targetTypeModel.(ui.SelectModel).GetSelectedItem()
	if changefeedType == nil {
		return "", errors.New("no changefeed type selected")
	}
	return changefeedType.(cdc.ChangefeedTypeEnum), nil
}

var createChangefeedInputDescription = map[string]string{
	flag.DisplayName:         "The name of the changefeed, skip to use the default name.",
	flag.ChangefeedMySQL:     "MySQL information in JSON format, use \"ticloud serverless changefeed template\" to see templates.",
	flag.ChangefeedKafka:     "Kafka information in JSON format, use \"ticloud serverless changefeed template --type kafka\" to see templates.",
	flag.ChangefeedFilter:    "Filter in JSON format, use \"ticloud serverless changefeed template --type filter\" to see templates.",
	flag.ChangefeedStartTSO:  "Start TSO for the changefeed (e.g., 443852055297916932). See https://docs.pingcap.com/tidb/stable/tso/ for more information about TSO.",
	flag.ChangefeedStartTime: "Start Time for the changefeed (RFC3339 format, e.g., 2024-01-01T00:00:00Z).",
}

var updateChangefeedInputDescriptionInteractive = map[string]string{
	flag.DisplayName:      "The new name of the changefeed, skip to keep the current name.",
	flag.ChangefeedKafka:  "Complete Kafka information in JSON format, use \"ticloud serverless changefeed template --type kafka\" to see templates.",
	flag.ChangefeedMySQL:  "Complete MySQL information in JSON format, use \"ticloud serverless changefeed template --type mysql\" to see templates.",
	flag.ChangefeedFilter: "Complete Filter in JSON format, use \"ticloud serverless changefeed template --type filter\" to see templates.",
}
