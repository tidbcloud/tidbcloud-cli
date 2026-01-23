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

package ui

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
)

// ConfirmPrompt prompts the user to type a confirmation string.
// promptText: the string the user must type to confirm.
// canPrompt: whether the terminal supports prompting.
func ConfirmPrompt(promptText string, canPrompt bool) error {
	if !canPrompt {
		return errors.New("the terminal doesn't support prompt, please run with --force to proceed")
	}

	confrim := "yes"

	confirmationMessage := fmt.Sprintf("%s %s %s %s",
		color.BlueString(promptText),
		color.BlueString("Please type"),
		color.HiBlueString(confrim),
		color.BlueString("to confirm:"),
	)

	prompt := &survey.Input{
		Message: confirmationMessage,
	}

	var userInput string
	err := survey.AskOne(prompt, &userInput)
	if err != nil {
		if err == terminal.InterruptErr {
			return util.InterruptError
		} else {
			return err
		}
	}

	if userInput != confrim {
		return errors.New("incorrect confirm string entered, skipping operation")
	}

	return nil
}
