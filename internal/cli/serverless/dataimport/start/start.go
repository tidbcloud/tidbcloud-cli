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
	"context"
	"fmt"
	"strconv"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	importOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type csvFormatField int

const (
	separatorIdx csvFormatField = iota
	delimiterIdx
	backslashEscapeIdx
	trimLastSeparatorIdx
)

type SourceType string

const (
	SourceTypeLOCAL   SourceType = "LOCAL"
	SourceTypeUnknown SourceType = "UNKNOWN"
)

type StartOpts struct {
	interactive bool
}

func (o StartOpts) SupportedFileTypes() []string {
	return []string{
		string(importModel.V1beta1ImportOptionsFileTypeCSV),
	}
}

func (c StartOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.FileType,
	}
}

func StartCmd(h *internal.Helper) *cobra.Command {
	var concurrency int
	opts := StartOpts{
		interactive: true,
	}
	startCmd := &cobra.Command{
		Use:         "start",
		Short:       "Start a data import task",
		Aliases:     []string{"create"},
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Start an import task in interactive mode:
  $ %[1]s serverless import start

  Start a local import task in non-interactive mode:
  $ %[1]s serverless import start --local.file-path <file-path> --cluster-id <cluster-id> --file-type <file-type> --local.target-database <target-database> --local.target-table <target-table>

  Start a local import task with custom upload concurrency:
  $ %[1]s serverless import start --local.file-path <file-path> --cluster-id <cluster-id> --file-type <file-type> --local.target-database <target-database> --local.target-table <target-table> --local.concurrency 10
	
  Start a local import task with custom CSV format:
  $ %[1]s serverless import start --local.file-path <file-path> --cluster-id <cluster-id> --file-type CSV --local.target-database <target-database> --local.target-table <target-table> --csv.separator \" --csv.delimiter \' --csv.backslash-escape=false --csv.trim-last-separator=true
`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// mark required flags in non-interactive mode
			if !opts.interactive {
				for _, fn := range flags {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var sourceType SourceType
			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}
				var err error
				sourceType, err = getSelectedSourceType()
				if err != nil {
					return err
				}
			} else {
				sourceType = SourceType(cmd.Flag(flag.SourceType).Value.String())
			}

			if sourceType == SourceTypeLOCAL {
				localOpts := LocalOpts{
					concurrency: concurrency,
					h:           h,
					interactive: opts.interactive,
				}
				return localOpts.Run(cmd)
			} else {
				return errors.New("unsupported import source type")
			}
		},
	}

	startCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	startCmd.Flags().String(flag.SourceType, "LOCAL", fmt.Sprintf("The import source type, one of %q.", []string{string(SourceTypeLOCAL)}))
	startCmd.Flags().String(flag.FileType, "", fmt.Sprintf("The import file type, one of %q.", opts.SupportedFileTypes()))

	startCmd.Flags().String(flag.LocalFilePath, "", "The local file path to import.")
	startCmd.Flags().String(flag.LocalTargetDatabase, "", "Target database to which import data.")
	startCmd.Flags().String(flag.LocalTargetTable, "", "Target table to which import data.")
	startCmd.Flags().IntVar(&concurrency, flag.LocalConcurrency, 5, "The concurrency for uploading file.")

	startCmd.Flags().String(flag.CSVDelimiter, "\"", "The delimiter used for quoting of CSV file.")
	startCmd.Flags().String(flag.CSVSeparator, ",", "The field separator of CSV file.")
	startCmd.Flags().Bool(flag.CSVTrimLastSeparator, false, "In CSV file whether to treat separator as the line terminator and trim all trailing separators.")
	startCmd.Flags().Bool(flag.CSVBackslashEscape, true, "In CSV file whether to parse backslash inside fields as escape characters.")

	return startCmd
}

func getSelectedSourceType() (SourceType, error) {
	SourceTypes := make([]interface{}, 0, 1)
	SourceTypes = append(SourceTypes, SourceTypeLOCAL)
	model, err := ui.InitialSelectModel(SourceTypes, "Choose import source type:")
	if err != nil {
		return SourceTypeUnknown, errors.Trace(err)
	}

	p := tea.NewProgram(model)
	SourceTypeModel, err := p.Run()
	if err != nil {
		return SourceTypeUnknown, errors.Trace(err)
	}
	if m, _ := SourceTypeModel.(ui.SelectModel); m.Interrupted {
		return SourceTypeUnknown, util.InterruptError
	}
	fileType := SourceTypeModel.(ui.SelectModel).GetSelectedItem()
	if fileType == nil {
		return SourceTypeUnknown, errors.New("no source type selected")
	}
	return fileType.(SourceType), nil
}

func waitStartOp(h *internal.Helper, d cloud.TiDBCloudClient, params *importOp.ImportServiceCreateImportParams) error {
	fmt.Fprintf(h.IOStreams.Out, "... Starting the import task\n")
	res, err := d.CreateImport(params)
	if err != nil {
		return err
	}

	fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", res.Payload.ID))
	return nil
}

func spinnerWaitStartOp(ctx context.Context, h *internal.Helper, d cloud.TiDBCloudClient, params *importOp.ImportServiceCreateImportParams) error {
	task := func() tea.Msg {
		errChan := make(chan error, 1)

		go func() {
			res, err := d.CreateImport(params)
			if err != nil {
				errChan <- err
				return
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", res.Payload.ID))
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
			case <-ctx.Done():
				return util.InterruptError
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Starting import task"))
	model, err := p.Run()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := model.(ui.SpinnerModel); m.Interrupted {
		return util.InterruptError
	}
	if m, _ := model.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	}

	return nil
}

func getCSVFormat() (separator string, delimiter string, backslashEscape bool, trimLastSeparator bool, errToReturn error) {
	separator, delimiter = ",", "\""
	backslashEscape, trimLastSeparator = true, false

	needCustomCSV := false
	prompt := &survey.Confirm{
		Message: "Do you need to custom CSV format?",
	}
	err := survey.AskOne(prompt, &needCustomCSV)
	if err != nil {
		if err == terminal.InterruptErr {
			errToReturn = util.InterruptError
			return
		} else {
			errToReturn = err
			return
		}
	}

	if needCustomCSV {
		// variables for input
		p := tea.NewProgram(initialCSVFormatInputModel())
		inputModel, err := p.Run()
		if err != nil {
			errToReturn = errors.Trace(err)
			return
		}
		if inputModel.(ui.TextInputModel).Interrupted {
			errToReturn = util.InterruptError
			return
		}

		// If user input is blank, use the default value.
		v := inputModel.(ui.TextInputModel).Inputs[separatorIdx].Value()
		if len(v) > 0 {
			separator = v
		}
		v = inputModel.(ui.TextInputModel).Inputs[delimiterIdx].Value()
		if len(v) > 0 {
			delimiter = v
		}
		v = inputModel.(ui.TextInputModel).Inputs[backslashEscapeIdx].Value()
		if len(v) > 0 {
			backslashEscape, err = strconv.ParseBool(v)
			if err != nil {
				errToReturn = errors.Annotate(err, "backslash escape must be true or false")
				return
			}
		}
		v = inputModel.(ui.TextInputModel).Inputs[trimLastSeparatorIdx].Value()
		if len(v) > 0 {
			trimLastSeparator, err = strconv.ParseBool(v)
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
		t.Cursor.Style = config.FocusedStyle
		t.CharLimit = 0
		f := csvFormatField(i)

		switch f {
		case separatorIdx:
			t.Placeholder = "separator, default is ',', empty to use default"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case delimiterIdx:
			t.Placeholder = "delimiter, default is '\"', empty to use default"
		case backslashEscapeIdx:
			t.Placeholder = "backslashEscape, default is true, empty to use default"
		case trimLastSeparatorIdx:
			t.Placeholder = "trimLastSeparator, default is false, empty to use default"
		}

		m.Inputs[i] = t
	}

	return m
}
