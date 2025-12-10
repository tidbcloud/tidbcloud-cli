// Copyright 2025 PingCAP, Inc.
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

package privatelink

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"
)

const (
	DisplayNamePrompt                     = "Input the display name:"
	AWSEndpointServiceNamePrompt          = "Input the AWS endpoint service name:"
	AlicloudEndpointServiceNamePrompt     = "Input the Alicloud endpoint service name:"
	AWSEndpointServiceRegionPrompt        = "Input the AWS endpoint service region:"
	AWSEndpointServiceRegionConfirmPrompt = "Is the endpoint service cross region?"
)

func GetSelectedPrivateLinkConnectionType() (privatelink.PrivateLinkConnectionTypeEnum, error) {
	types := make([]interface{}, 0, len(privatelink.AllowedPrivateLinkConnectionTypeEnumEnumValues))
	for _, v := range privatelink.AllowedPrivateLinkConnectionTypeEnumEnumValues {
		types = append(types, v)
	}
	selectModel, err := ui.InitialSelectModel(types, "Choose the private link connection type:")
	if err != nil {
		return "", errors.Trace(err)
	}
	p := tea.NewProgram(selectModel)
	model, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := model.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}

	resp := model.(ui.SelectModel).GetSelectedItem()
	if resp == nil {
		return "", errors.New("no private link connection type selected")
	}
	return resp.(privatelink.PrivateLinkConnectionTypeEnum), nil
}

func GetSelectedPLCAttachDomainType() (privatelink.PrivateLinkConnectionDomainTypeEnum, error) {
	types := make([]interface{}, 0, len(privatelink.AllowedPrivateLinkConnectionDomainTypeEnumEnumValues))
	for _, v := range privatelink.AllowedPrivateLinkConnectionDomainTypeEnumEnumValues {
		types = append(types, v)
	}
	selectModel, err := ui.InitialSelectModel(types, "Choose the private link connection domain type:")
	if err != nil {
		return "", errors.Trace(err)
	}
	p := tea.NewProgram(selectModel)
	model, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := model.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}

	resp := model.(ui.SelectModel).GetSelectedItem()
	if resp == nil {
		return "", errors.New("no private link connection domain type selected")
	}
	return resp.(privatelink.PrivateLinkConnectionDomainTypeEnum), nil
}
