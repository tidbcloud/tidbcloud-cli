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
	"context"
	"fmt"
	"os"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/aws/s3"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	importOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/bubbles/progress"
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
		flag.DataFormat,
		flag.TargetDatabase,
		flag.TargetTable,
	}
}

func (c LocalOpts) SupportedDataFormats() []string {
	return []string{
		string(importModel.V1beta1DataFormatCSV),
	}
}

func LocalCmd(h *internal.Helper) *cobra.Command {
	opts := LocalOpts{
		interactive: true,
	}
	var partSize int64
	var concurrency int

	var localCmd = &cobra.Command{
		Use:         "local <file-path>",
		Short:       "Import a local file to TiDB Cloud",
		Args:        util.RequiredArgs("file-path"),
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Start an import task in interactive mode:
  $ %[1]s serverless import start local <file-path>

  Start an import task in non-interactive mode:
  $ %[1]s serverless import start local <file-path> --cluster-id <cluster-id> --data-format <data-format> --target-database <target-database> --target-table <target-table>
	
  Start an import task with custom CSV format:
  $ %[1]s serverless import start local <file-path> --cluster-id <cluster-id> --data-format CSV --target-database <target-database> --target-table <target-table> --separator \" --delimiter \' --backslash-escape=false --trim-last-separator=true
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
			ctx := cmd.Context()
			var clusterID, dataFormat, targetDatabase, targetTable, separator, delimiter string
			var backslashEscape, trimLastSeparator bool
			d, err := h.Client()
			if err != nil {
				return err
			}
			uploader := h.Uploader(d)
			err = uploader.SetConcurrency(concurrency)
			if err != nil {
				return err
			}
			err = uploader.SetPartSize(partSize * 1024 * 1024)
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

				cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
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
					return util.InterruptError
				}
				dataFormat = formatModel.(ui.SelectModel).Choices[formatModel.(ui.SelectModel).Selected].(string)

				// variables for input
				p = tea.NewProgram(initialLocalInputModel())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
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

			cmd.Annotations[telemetry.ClusterID] = clusterID

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

			var uploadID string
			input := &s3.PutObjectInput{
				FileName:      aws.String(stat.Name()),
				DatabaseName:  aws.String(targetDatabase),
				TableName:     aws.String(targetTable),
				ContentLength: aws.Int64(stat.Size()),
				ClusterID:     aws.String(clusterID),
				Body:          uploadFile,
			}
			if h.IOStreams.CanPrompt {
				uploadID, err = spinnerWaitUploadOp(ctx, h, uploader, input)
				if err != nil {
					return err
				}
			} else {
				uploadID, err = waitUploadOp(ctx, h, uploader, input)
				if err != nil {
					return err
				}
			}

			body := importOp.ImportServiceCreateImportBody{}
			err = body.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"type": "LOCAL",
			"dataFormat": "%s",
			"importOptions": {
				"csvFormat": {
                	"separator": ",",
					"delimiter": "\"",
					"header": true,
					"backslashEscape": true,
					"null": "\\N",
					"trimLastSeparator": false,
					"notNull": false
				}
			},
			"target": {
				"local": {
					"uploadId": "%s",
					"targetTable": {
						"schema": "%s",
						"table": "%s"
					},
					"fileName": "%s"
				},
				"type": "LOCAL"
			}
			}`, dataFormat, uploadID, targetDatabase, targetTable, stat.Name())))
			if err != nil {
				return errors.Trace(err)
			}

			body.ImportOptions.CsvFormat.Separator = separator
			body.ImportOptions.CsvFormat.Delimiter = delimiter
			body.ImportOptions.CsvFormat.BackslashEscape = backslashEscape
			body.ImportOptions.CsvFormat.TrimLastSeparator = trimLastSeparator

			params := importOp.NewImportServiceCreateImportParams().WithClusterID(clusterID).
				WithBody(body).WithContext(ctx)
			if h.IOStreams.CanPrompt {
				err := spinnerWaitStartOp(ctx, h, d, params)
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

	localCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	localCmd.Flags().String(flag.DataFormat, "", fmt.Sprintf("Data format, one of %q", opts.SupportedDataFormats()))
	localCmd.Flags().String(flag.TargetDatabase, "", "Target database to which import data")
	localCmd.Flags().String(flag.TargetTable, "", "Target table to which import data")
	localCmd.Flags().String(flag.Delimiter, "\"", "The delimiter used for quoting of CSV file")
	localCmd.Flags().String(flag.Separator, ",", "The field separator of CSV file")
	localCmd.Flags().Bool(flag.TrimLastSeparator, false, "In CSV file whether to treat Separator as the line terminator and trim all trailing separators")
	localCmd.Flags().Bool(flag.BackslashEscape, true, "In CSV file whether to parse backslash inside fields as escape characters")
	localCmd.Flags().Int64Var(&partSize, flag.PartSize, 5, "The part size for uploading file(MiB), default is 5")
	localCmd.Flags().IntVar(&concurrency, flag.Concurrency, 5, "The concurrency for uploading file, default is 5")
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

func waitUploadOp(ctx context.Context, h *internal.Helper, u s3.Uploader, input *s3.PutObjectInput) (string, error) {
	fmt.Fprintf(h.IOStreams.Out, "... Uploading file\n")

	p := make(chan float64)
	e := make(chan error)
	input.OnProgress = func(ratio float64) {
		p <- ratio
	}

	var id string
	var err error
	go func() {
		id, err = u.Upload(ctx, input)
		e <- err
	}()
	timer := time.After(2 * time.Hour)
	for {
		select {
		case progress := <-p:
			fmt.Fprintf(h.IOStreams.Out, "upload progress: %.2f%%\n", progress*100)
		case <-timer:
			return "", fmt.Errorf("time out when uploading file")
		case err := <-e:
			if err != nil {
				return "", err
			}
			fmt.Fprintln(h.IOStreams.Out, "File has been uploaded")
			return id, nil
		}
	}
}

func spinnerWaitUploadOp(ctx context.Context, h *internal.Helper, u s3.Uploader, input *s3.PutObjectInput) (string, error) {
	var uploadID string
	m := ui.ProcessModel{
		Progress: progress.New(progress.WithDefaultGradient()),
	}

	p := tea.NewProgram(m)
	input.OnProgress = func(ratio float64) {
		p.Send(ui.ProgressMsg(ratio))
	}

	go func() {
		var err error
		uploadID, err = u.Upload(ctx, input)
		if err != nil {
			p.Send(ui.ProgressErrMsg{
				Err: err,
			})
		}
		input.OnProgress(1.0)
	}()

	fmt.Fprintf(h.IOStreams.Out, color.GreenString("Start uploading...\n"))

	processModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if processModel.(ui.ProcessModel).Interrupted {
		return "", util.InterruptError
	}

	fmt.Fprintln(h.IOStreams.Out, color.GreenString("File has been uploaded"))
	return uploadID, nil
}
