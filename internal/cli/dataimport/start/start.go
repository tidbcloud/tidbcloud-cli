// Copyright 2022 PingCAP, Inc.
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

package start

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/import_operations"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type csvFormatField int

const (
	delimiterIdx csvFormatField = iota
	quoteIdx
	backslashEscapeIdx
	hasHeaderRowIdx
)

func StartCmd(h *internal.Helper) *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start an import task",
	}

	startCmd.PersistentFlags().String(flag.Delimiter, ",", "The delimiter character used to separate fields in the CSV data")
	startCmd.PersistentFlags().String(flag.Quote, "\"", "The character used to quote the fields in the CSV data")
	startCmd.PersistentFlags().Bool(flag.HasHeaderRow, true, "In CSV file whether the CSV data has a header row, which is not part of the data")
	startCmd.PersistentFlags().Bool(flag.BackslashEscape, true, "In CSV file whether to parse backslash(\\) inside fields as escape characters")

	startCmd.AddCommand(LocalCmd(h))
	startCmd.AddCommand(S3Cmd(h))
	return startCmd
}

func waitStartOp(h *internal.Helper, d cloud.TiDBCloudClient, params *import_operations.CreateImportTaskParams) error {
	fmt.Fprintf(h.IOStreams.Out, "... Starting the import task\n")
	res, err := d.CreateImportTask(params)
	if err != nil {
		return err
	}

	fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", *(res.Payload.ID)))
	return nil
}

func spinnerWaitStartOp(h *internal.Helper, d cloud.TiDBCloudClient, params *import_operations.CreateImportTaskParams) error {
	task := func() tea.Msg {
		errChan := make(chan error, 1)

		go func() {
			res, err := d.CreateImportTask(params)
			if err != nil {
				errChan <- err
				return
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", *(res.Payload.ID)))
			errChan <- nil
		}()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		timer := time.After(2 * time.Minute)
		for {
			select {
			case <-timer:
				return fmt.Errorf("timeout waiting for import task to start")
			case <-ticker.C:
				// continue
			case err := <-errChan:
				if err != nil {
					return err
				} else {
					return ui.Result("")
				}
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Starting import task"))
	createModel, err := p.StartReturningModel()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := createModel.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	}

	return nil
}

func getCSVFormat() (delimiter string, quote string, backslashEscape bool, hasHeaderRow bool, errToReturn error) {
	delimiter, quote = ",", "\""
	backslashEscape, hasHeaderRow = true, true

	needCustomCSV := false
	prompt := &survey.Confirm{
		Message: "Do you need to custom CSV format?",
	}
	err := survey.AskOne(prompt, &needCustomCSV)
	if err != nil {
		if err == terminal.InterruptErr {
			os.Exit(130)
		} else {
			errToReturn = err
			return
		}
	}

	if needCustomCSV {
		// variables for input
		p := tea.NewProgram(initialCSVFormatInputModel())
		inputModel, err := p.StartReturningModel()
		if err != nil {
			errToReturn = errors.Trace(err)
			return
		}
		if inputModel.(ui.TextInputModel).Interrupted {
			os.Exit(130)
		}

		// If user input is blank, use the default value.
		v := inputModel.(ui.TextInputModel).Inputs[delimiterIdx].Value()
		if len(v) > 0 {
			delimiter = v
		}
		v = inputModel.(ui.TextInputModel).Inputs[quoteIdx].Value()
		if len(v) > 0 {
			quote = v
		}
		v = inputModel.(ui.TextInputModel).Inputs[backslashEscapeIdx].Value()
		if len(v) > 0 {
			backslashEscape, err = strconv.ParseBool(v)
			if err != nil {
				errToReturn = errors.Annotate(err, "backslash escape must be true or false")
				return
			}
		}
		v = inputModel.(ui.TextInputModel).Inputs[hasHeaderRowIdx].Value()
		if len(v) > 0 {
			hasHeaderRow, err = strconv.ParseBool(v)
			if err != nil {
				errToReturn = errors.Annotate(err, "backslash escape must be true or false")
				return
			}
		}
	}
	return
}

func initialCSVFormatInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 4),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = config.FocusedStyle
		t.CharLimit = 0
		f := csvFormatField(i)

		switch f {
		case delimiterIdx:
			t.Placeholder = "delimiter, default is ',', empty to use default"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case quoteIdx:
			t.Placeholder = "quote, default is '\"', empty to use default"
		case backslashEscapeIdx:
			t.Placeholder = "backslashEscape, default is true, empty to use default"
		case hasHeaderRowIdx:
			t.Placeholder = "hasHeaderRow, default is false, empty to use default"
		}

		m.Inputs[i] = t
	}

	return m
}
