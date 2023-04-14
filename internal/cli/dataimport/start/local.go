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
	"bytes"
	"compress/gzip"
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

	"github.com/alecthomas/units"
	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/import_operations"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/go-openapi/strfmt"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type localImportField int

const (
	databaseIdx localImportField = iota
	tableIdx

	maxFileSize = 50 * units.MiB
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
		import_operations.CreateImportTaskParamsBodySpecSourceFormatTypeCSV,
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
  $ %[1]s import start local <file-path> --project-id <project-id> --cluster-id <cluster-id> --data-format CSV --target-database <target-database> --target-table <target-table> --backslash-escape=false --has-header-row=false
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
			var projectID, clusterID, dataFormat, targetDatabase, targetTable, quote, delimiter string
			var backslashEscape, hasHeaderRow bool
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

				quote, delimiter, backslashEscape, hasHeaderRow, err = getCSVFormat()
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
				quote, err = cmd.Flags().GetString(flag.Quote)
				if err != nil {
					return errors.Trace(err)
				}
				delimiter, err = cmd.Flags().GetString(flag.Delimiter)
				if err != nil {
					return errors.Trace(err)
				}
				hasHeaderRow, err = cmd.Flags().GetBool(flag.HasHeaderRow)
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
			if stat.Size() > int64(maxFileSize) {
				return fmt.Errorf("file size %s exceeds the maximum size %s", humanize.IBytes(uint64(stat.Size())), humanize.IBytes(uint64(maxFileSize)))
			}
			size := strconv.FormatInt(stat.Size(), 10)
			fileName := stat.Name()
			content := make([]byte, stat.Size())
			_, err = uploadFile.Read(content)
			if err != nil {
				return err
			}

			encodeContent, err := gzipCompress(content)
			if err != nil {
				return err
			}
			baseContent := strfmt.Base64(encodeContent)
			params := import_operations.NewUploadLocalFileParams().WithProjectID(projectID).WithClusterID(clusterID).WithBody(import_operations.UploadLocalFileBody{
				LocalFileName: &fileName,
				Payload: &import_operations.UploadLocalFileParamsBodyPayload{
					Content:        &baseContent,
					TotalSizeBytes: &size,
				},
			})

			var stubID *string
			if h.IOStreams.CanPrompt {
				id, err := spinnerWaitUploadOp(h, d, params)
				stubID = id
				if err != nil {
					return err
				}
			} else {
				id, err := waitUploadOp(h, d, params)
				stubID = id
				if err != nil {
					return err
				}
			}

			dataSourceType := "LOCAL_FILE"
			uri := fmt.Sprintf("file://%s/", *stubID)
			previewImportData, err := d.PreviewImportData(import_operations.NewPreviewImportDataParams().WithProjectID(projectID).WithClusterID(clusterID).WithBody(import_operations.PreviewImportDataBody{
				Spec: &import_operations.PreviewImportDataParamsBodySpec{
					Source: &import_operations.PreviewImportDataParamsBodySpecSource{
						Format: &import_operations.PreviewImportDataParamsBodySpecSourceFormat{
							CsvConfig: &import_operations.PreviewImportDataParamsBodySpecSourceFormatCsvConfig{
								BackslashEscape: &backslashEscape,
								Delimiter:       &delimiter,
								HasHeaderRow:    &hasHeaderRow,
								Quote:           &quote,
							},
							Type: &dataFormat,
						},
						Type: &dataSourceType,
						URI:  &uri,
					},
					Target: &import_operations.PreviewImportDataParamsBodySpecTarget{
						Tables: []*import_operations.PreviewImportDataParamsBodySpecTargetTablesItems0{
							{
								DatabaseName:    &targetDatabase,
								FileNamePattern: fileName,
								TableName:       &targetTable,
							},
						},
					},
				},
			}))
			if err != nil {
				return err
			}

			tablePreview := previewImportData.Payload.TablePreviews[0]
			createParams := import_operations.NewCreateImportTaskParams().WithProjectID(projectID).WithClusterID(clusterID).WithBody(import_operations.CreateImportTaskBody{
				Options: &import_operations.CreateImportTaskParamsBodyOptions{
					PreCreateTables: []*import_operations.CreateImportTaskParamsBodyOptionsPreCreateTablesItems0{
						{
							DatabaseName: tablePreview.DatabaseName,
							Schema:       parseSchema(tablePreview.SchemaPreview),
							TableName:    tablePreview.TableName,
						},
					},
				},
				Spec: &import_operations.CreateImportTaskParamsBodySpec{
					Source: &import_operations.CreateImportTaskParamsBodySpecSource{
						Format: &import_operations.CreateImportTaskParamsBodySpecSourceFormat{
							CsvConfig: &import_operations.CreateImportTaskParamsBodySpecSourceFormatCsvConfig{
								BackslashEscape: &backslashEscape,
								Delimiter:       &delimiter,
								HasHeaderRow:    &hasHeaderRow,
								Quote:           &quote,
							},
							Type: &dataFormat,
						},
						Type: &dataSourceType,
						URI:  &uri,
					},
					Target: &import_operations.CreateImportTaskParamsBodySpecTarget{
						Tables: []*import_operations.CreateImportTaskParamsBodySpecTargetTablesItems0{
							{
								DatabaseName:    tablePreview.DatabaseName,
								FileNamePattern: fileName,
								TableName:       tablePreview.TableName,
							},
						},
					},
				},
			})

			if h.IOStreams.CanPrompt {
				err := spinnerWaitStartOp(h, d, createParams)
				if err != nil {
					return err
				}
			} else {
				err := waitStartOp(h, d, createParams)
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

func gzipCompress(content []byte) ([]byte, error) {
	// Compress the contents with gzip
	var compressedContent bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedContent)
	if _, err := gzipWriter.Write(content); err != nil {
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	return compressedContent.Bytes(), nil
}

func parseSchema(schemaPreview *import_operations.PreviewImportDataOKBodyTablePreviewsItems0SchemaPreview) *import_operations.CreateImportTaskParamsBodyOptionsPreCreateTablesItems0Schema {
	columns := make([]*import_operations.CreateImportTaskParamsBodyOptionsPreCreateTablesItems0SchemaColumnDefinitionsItems0, 0, len(schemaPreview.ColumnDefinitions))

	for _, column := range schemaPreview.ColumnDefinitions {
		columns = append(columns, &import_operations.CreateImportTaskParamsBodyOptionsPreCreateTablesItems0SchemaColumnDefinitionsItems0{
			ColumnName: column.ColumnName,
			ColumnType: column.ColumnType,
		})
	}

	primaryKeys := make([]string, 0, len(schemaPreview.PrimaryKeyColumns))
	for _, pk := range schemaPreview.PrimaryKeyColumns {
		primaryKeys = append(primaryKeys, pk)
	}
	return &import_operations.CreateImportTaskParamsBodyOptionsPreCreateTablesItems0Schema{
		ColumnDefinitions: columns,
		PrimaryKeyColumns: primaryKeys,
	}
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

func waitUploadOp(h *internal.Helper, d cloud.TiDBCloudClient, params *import_operations.UploadLocalFileParams) (*string, error) {
	fmt.Fprintf(h.IOStreams.Out, "... Uploading file\n")
	res, err := d.UploadLocalFile(params)
	if err != nil {
		return nil, err
	}

	if !res.IsSuccess() {
		return nil, fmt.Errorf(res.Error())
	}

	fmt.Fprintln(h.IOStreams.Out, "File has been uploaded")
	return res.GetPayload().UploadStubID, nil
}

func spinnerWaitUploadOp(h *internal.Helper, d cloud.TiDBCloudClient, params *import_operations.UploadLocalFileParams) (*string, error) {
	var res *import_operations.UploadLocalFileOK

	task := func() tea.Msg {
		errChan := make(chan error, 1)

		go func() {
			var err error
			res, err = d.UploadLocalFile(params)
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
		return nil, errors.Trace(err)
	}
	if m, _ := createModel.(ui.SpinnerModel); m.Err != nil {
		return nil, m.Err
	} else {
		fmt.Fprintf(h.IOStreams.Out, color.GreenString(m.Output))
	}

	return res.GetPayload().UploadStubID, nil
}
