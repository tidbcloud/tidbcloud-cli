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

package private_link_connection

import (
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
)

var createAWSField = map[string]int{
	flag.DisplayName:              0,
	flag.AWSEndpointServiceName:   1,
	flag.AWSEndpointServiceRegion: 2,
}

var createAlicloudField = map[string]int{
	flag.DisplayName:                 0,
	flag.AlicloudEndpointServiceName: 1,
}

func GetPrivateLinkConnectionType() (privatelink.PrivateLinkConnectionTypeEnum, error) {
	items := []interface{}{
		string(privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE),
		string(privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE),
	}
	model, err := ui.InitialSelectModel(items, "Choose the private link connection type:")
	if err != nil {
		return "", errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)

	p := tea.NewProgram(model)
	typeModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := typeModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	selected := typeModel.(ui.SelectModel).GetSelectedItem()
	if selected == nil {
		return "", errors.New("no private link connection type selected")
	}
	return privatelink.PrivateLinkConnectionTypeEnum(selected.(string)), nil
}

func initialCreateAWSInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(createAWSField)),
	}

	for k, v := range createAWSField {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 128

		switch k {
		case flag.DisplayName:
			t.Placeholder = "Display Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.AWSEndpointServiceName:
			t.Placeholder = "AWS Endpoint Service Name"
		case flag.AWSEndpointServiceRegion:
			t.Placeholder = "AWS Endpoint Service Region (optional)"
		}
		m.Inputs[v] = t
	}
	return m
}

func GetCreateAWSInput() (tea.Model, error) {
	p := tea.NewProgram(initialCreateAWSInputModel())
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}

func initialCreateAlicloudInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(createAlicloudField)),
	}

	for k, v := range createAlicloudField {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 128

		switch k {
		case flag.DisplayName:
			t.Placeholder = "Display Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.AlicloudEndpointServiceName:
			t.Placeholder = "Alicloud Endpoint Service Name"
		}
		m.Inputs[v] = t
	}
	return m
}

func GetCreateAlicloudInput() (tea.Model, error) {
	p := tea.NewProgram(initialCreateAlicloudInputModel())
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}

func normalizePrivateLinkConnectionType(value string) (privatelink.PrivateLinkConnectionTypeEnum, error) {
	enumValue := privatelink.PrivateLinkConnectionTypeEnum(value)
	if !enumValue.IsValid() {
		return "", fmt.Errorf("invalid private link connection type: %s", value)
	}
	return enumValue, nil
}
