// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package example

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"

	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
)

const (
	ExampleNamePrompt    = "Input the example name:"
	DeleteConfirmPrompt  = "Type 'yes' to confirm deletion:"
	UnwiredFeaturePrompt = "This command is not wired to an API yet."
)

var s3InputDescriptions = map[string]string{
	flag.S3URI:             "S3 URI (e.g., s3://bucket/path)",
	flag.S3AccessKeyID:     "S3 access key ID",
	flag.S3SecretAccessKey: "S3 secret access key",
}

type Example struct {
	ID          string
	DisplayName string
}

func (e Example) String() string {
	return fmt.Sprintf("%s(%s)", e.DisplayName, e.ID)
}

func GetSelectedExample(ctx context.Context, clusterID string, pageSize int64, client cloud.TiDBCloudClient) (*Example, error) {
	// examples, err := client.ListExamples(ctx, clusterID, pageSize, nil)
	// if err != nil {
	// 	return nil, errors.Trace(err)
	// }

	// var items = make([]interface{}, 0, len(examples))
	// for _, example := range examples {
	// 	items = append(items, &Example{ID: *example.ExampleId, DisplayName: *example.DisplayName})
	// }

	items := []interface{}{
		&Example{ID: "example-1", DisplayName: "Example One"},
		&Example{ID: "example-2", DisplayName: "Example Two"},
		&Example{ID: "example-3", DisplayName: "Example Three"},
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no available examples found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the example:")
	if err != nil {
		return nil, errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	exampleModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if m, _ := exampleModel.(ui.SelectModel); m.Interrupted {
		return nil, util.InterruptError
	}
	selected := exampleModel.(ui.SelectModel).GetSelectedItem()
	if selected == nil {
		return nil, errors.New("no example selected")
	}
	return selected.(*Example), nil
}

func GetDisplayNameInput() (string, error) {
	model, err := ui.InitialOneInputModel(ExampleNamePrompt, "example-name")
	if err != nil {
		return "", errors.Trace(err)
	}
	if model.Interrupted {
		return "", util.InterruptError
	}
	if model.Err != nil {
		return "", model.Err
	}
	return model.Input.Value(), nil
}

func GetS3Inputs() (string, string, string, error) {
	fmt.Fprintln(h.IOStreams.Out, "Pl")
	inputs := []string{flag.S3URI, flag.S3AccessKeyID, flag.S3SecretAccessKey}
	textInput, err := ui.InitialInputModel(inputs, s3InputDescriptions)
	if err != nil {
		return "", "", "", err
	}

	s3URI := textInput.Inputs[0].Value()
	if s3URI == "" {
		return "", "", "", errors.New("empty S3 URI")
	}
	s3AccessKeyID := textInput.Inputs[1].Value()
	if s3AccessKeyID == "" {
		return "", "", "", errors.New("empty S3 access key ID")
	}
	s3SecretAccessKey := textInput.Inputs[2].Value()
	if s3SecretAccessKey == "" {
		return "", "", "", errors.New("empty S3 secret access key")
	}
	return s3URI, s3AccessKeyID, s3SecretAccessKey, nil
}
