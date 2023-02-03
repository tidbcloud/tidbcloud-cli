// Copyright 2023 PingCAP, Inc.
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
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	importOp "tidbcloud-cli/pkg/tidbcloud/import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/import/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type localImportField int

const (
	databaseIdx localImportField = iota
	tableIdx
)

type LocalOpts struct {
	interactive bool
}

func (c LocalOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
		flag.DataFormat,
		flag.TargetDatabase,
		flag.TargetTable,
	}
}

func (c LocalOpts) SupportedDataFormats() []string {
	return []string{
		string(importModel.OpenapiDataFormatCSV),
	}
}

func LocalCmd(h *internal.Helper) *cobra.Command {
	opts := LocalOpts{
		interactive: true,
	}

	var localCmd = &cobra.Command{
		Use:         "local <file-path>",
		Short:       "Import a local file to TiDB Cloud",
		Args:        util.RequiredArgs("file-path"),
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Start an import task in interactive mode:
  $ %[1]s import start local <file-path>

  Start an import task in non-interactive mode:
  $ %[1]s import start local <file-path> --project-id <project-id> --cluster-id <cluster-id> --data-format <data-format> --target-database <target-database> --target-table <target-table>
	
  Start an import task with custom CSV format:
  $ %[1]s import start local <file-path> --project-id <project-id> --cluster-id <cluster-id> --data-format CSV --target-database <target-database> --target-table <target-table> --separator \" --delimiter \' --backslash-escape=false --trim-last-separator=true
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
			var projectID, clusterID, dataFormat, targetDatabase, targetTable, separator, delimiter string
			var backslashEscape, trimLastSeparator bool
			d, err := h.Client()
			if err != nil {
				return err
			}

			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				projectID = project.ID

				cluster, err := cloud.GetSelectedCluster(projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				var dataFormats []interface{}
				for _, f := range opts.SupportedDataFormats() {
					dataFormats = append(dataFormats, f)
				}
				model, err := ui.InitialSelectModel(dataFormats, "Choose the data format:")
				if err != nil {
					return err
				}
				p := tea.NewProgram(model)
				formatModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if m, _ := formatModel.(ui.SelectModel); m.Interrupted {
					os.Exit(130)
				}
				dataFormat = formatModel.(ui.SelectModel).Choices[formatModel.(ui.SelectModel).Selected].(string)

				// variables for input
				p = tea.NewProgram(initialLocalInputModel())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return nil
				}

				targetDatabase = inputModel.(ui.TextInputModel).Inputs[databaseIdx].Value()
				if len(targetDatabase) == 0 {
					return errors.New("Target database is required")
				}
				targetTable = inputModel.(ui.TextInputModel).Inputs[tableIdx].Value()
				if len(targetTable) == 0 {
					return errors.New("Target table is required")
				}

				separator, delimiter, backslashEscape, trimLastSeparator, err = getCSVFormat()
				if err != nil {
					return err
				}
			} else {
				// non-interactive mode
				projectID = cmd.Flag(flag.ProjectID).Value.String()
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				dataFormat = cmd.Flag(flag.DataFormat).Value.String()
				if !util.ElemInSlice(opts.SupportedDataFormats(), dataFormat) {
					return fmt.Errorf("data format %s is not supported, please use one of %q", dataFormat, opts.SupportedDataFormats())
				}
				targetDatabase = cmd.Flag(flag.TargetDatabase).Value.String()
				targetTable = cmd.Flag(flag.TargetTable).Value.String()

				// optional flags
				backslashEscape, err = cmd.Flags().GetBool(flag.BackslashEscape)
				if err != nil {
					return errors.Trace(err)
				}
				separator, err = cmd.Flags().GetString(flag.Separator)
				if err != nil {
					return errors.Trace(err)
				}
				delimiter, err = cmd.Flags().GetString(flag.Delimiter)
				if err != nil {
					return errors.Trace(err)
				}
				trimLastSeparator, err = cmd.Flags().GetBool(flag.TrimLastSeparator)
				if err != nil {
					return errors.Trace(err)
				}
			}

			cmd.Annotations[telemetry.ProjectID] = projectID

			filePath := args[0]
			uploadFile, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer uploadFile.Close()

			stat, err := uploadFile.Stat()
			if err != nil {
				return err
			}
			size := strconv.FormatInt(stat.Size(), 10)
			name := stat.Name()
			urlRes, err := d.GenerateUploadURL(importOp.NewGenerateUploadURLParams().WithProjectID(projectID).WithClusterID(clusterID).WithBody(importOp.GenerateUploadURLBody{
				ContentLength: &size,
				FileName:      &name,
			}))
			if err != nil {
				return err
			}
			url := urlRes.Payload.UploadURL

			if h.IOStreams.CanPrompt {
				err := spinnerWaitUploadOp(h, d, url, uploadFile, stat.Size())
				if err != nil {
					return err
				}
			} else {
				err := waitUploadOp(h, d, url, uploadFile, stat.Size())
				if err != nil {
					return err
				}
			}

			body := importOp.CreateImportBody{}
			err = body.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"type": "LOCAL",
			"data_format": "%s",
			"file_name": "%s",
			"csv_format": {
                "separator": ",",
				"delimiter": "\"",
				"header": true,
				"backslash_escape": true,
				"null": "\\N",
				"trim_last_separator": false,
				"not_null": false
			},
			"target_table": {
				"schema": "%s",
				"table": "%s"
			}}`, dataFormat, *urlRes.Payload.NewFileName, targetDatabase, targetTable)))
			if err != nil {
				return errors.Trace(err)
			}

			body.CsvFormat.Separator = separator
			body.CsvFormat.Delimiter = delimiter
			body.CsvFormat.BackslashEscape = backslashEscape
			body.CsvFormat.TrimLastSeparator = trimLastSeparator

			params := importOp.NewCreateImportParams().WithProjectID(projectID).WithClusterID(clusterID).
				WithBody(body)
			if h.IOStreams.CanPrompt {
				err := spinnerWaitStartOp(h, d, params)
				if err != nil {
					return err
				}
			} else {
				err := waitStartOp(h, d, params)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	localCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "Project ID")
	localCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	localCmd.Flags().String(flag.DataFormat, "", fmt.Sprintf("Data format, one of %q", opts.SupportedDataFormats()))
	localCmd.Flags().String(flag.TargetDatabase, "", "Target database to which import data")
	localCmd.Flags().String(flag.TargetTable, "", "Target table to which import data")
	return localCmd
}

func initialLocalInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = config.FocusedStyle
		t.CharLimit = 0
		f := localImportField(i)

		switch f {
		case databaseIdx:
			t.Placeholder = "Target database"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case tableIdx:
			t.Placeholder = "Target table"
		}

		m.Inputs[i] = t
	}

	return m
}

func waitUploadOp(h *internal.Helper, d cloud.TiDBCloudClient, url *string, uploadFile *os.File, size int64) error {
	fmt.Fprintf(h.IOStreams.Out, "... Uploading file\n")
	err := d.PreSignedUrlUpload(url, uploadFile, size)
	if err != nil {
		return err
	}

	fmt.Fprintln(h.IOStreams.Out, "File has been uploaded")
	return nil
}

func spinnerWaitUploadOp(h *internal.Helper, d cloud.TiDBCloudClient, url *string, uploadFile *os.File, size int64) error {
	task := func() tea.Msg {
		errChan := make(chan error, 1)

		go func() {
			err := d.PreSignedUrlUpload(url, uploadFile, size)
			if err != nil {
				errChan <- err
				return
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("File has been uploaded"))
			errChan <- nil
		}()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		timer := time.After(2 * time.Minute)
		for {
			select {
			case <-timer:
				return fmt.Errorf("timeout waiting for uploading file")
			case <-ticker.C:
				// continue
			case err := <-errChan:
				if err != nil {
					return err
				} else {
					return ui.Result("File has been uploaded")
				}
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Uploading file"))
	createModel, err := p.StartReturningModel()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := createModel.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	} else {
		fmt.Fprintf(h.IOStreams.Out, color.GreenString(m.Output))
	}

	return nil
}
