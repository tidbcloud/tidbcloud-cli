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
	filePathIdx localImportField = iota
	databaseIdx
	tableIdx
)

type LocalOpts struct {
	concurrency int
	h           *internal.Helper
	interactive bool
}

func (o LocalOpts) SupportedFileTypes() []string {
	return []string{
		string(importModel.V1beta1ImportOptionsFileTypeCSV),
	}
}

func (o LocalOpts) Run(cmd *cobra.Command) error {
	ctx := cmd.Context()
	var clusterID, fileType, targetDatabase, targetTable, separator, delimiter, filePath string
	var backslashEscape, trimLastSeparator bool
	d, err := o.h.Client()
	if err != nil {
		return err
	}
	uploader := o.h.Uploader(d)
	err = uploader.SetConcurrency(o.concurrency)
	if err != nil {
		return err
	}

	if o.interactive {
		cmd.Annotations[telemetry.InteractiveMode] = "true"
		if !o.h.IOStreams.CanPrompt {
			return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
		}

		// interactive mode
		project, err := cloud.GetSelectedProject(o.h.QueryPageSize, d)
		if err != nil {
			return err
		}

		cluster, err := cloud.GetSelectedCluster(project.ID, o.h.QueryPageSize, d)
		if err != nil {
			return err
		}
		clusterID = cluster.ID

		var fileTypes []interface{}
		for _, f := range o.SupportedFileTypes() {
			fileTypes = append(fileTypes, f)
		}
		model, err := ui.InitialSelectModel(fileTypes, "Choose the source file type:")
		if err != nil {
			return err
		}
		p := tea.NewProgram(model)
		formatModel, err := p.Run()
		if err != nil {
			return errors.Trace(err)
		}
		if m, _ := formatModel.(ui.SelectModel); m.Interrupted {
			return util.InterruptError
		}
		fileType = formatModel.(ui.SelectModel).Choices[formatModel.(ui.SelectModel).Selected].(string)

		// variables for input
		p = tea.NewProgram(initialLocalInputModel())
		inputModel, err := p.Run()
		if err != nil {
			return errors.Trace(err)
		}
		if inputModel.(ui.TextInputModel).Interrupted {
			return util.InterruptError
		}

		filePath = inputModel.(ui.TextInputModel).Inputs[filePathIdx].Value()
		if len(filePath) == 0 {
			return errors.New("Local file path is required")
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
		fileType = cmd.Flag(flag.FileType).Value.String()
		if !util.ElemInSlice(o.SupportedFileTypes(), fileType) {
			return fmt.Errorf("file type \"%s\" is not supported, please use one of %q", fileType, o.SupportedFileTypes())
		}
		targetDatabase = cmd.Flag(flag.LocalTargetDatabase).Value.String()
		targetTable = cmd.Flag(flag.LocalTargetTable).Value.String()
		f := cmd.Flags().Lookup(flag.LocalFilePath)
		if !f.Changed {
			return errors.New("required flag(s) \"local.file-path\" not set")
		}
		filePath = f.Value.String()

		// optional flags
		backslashEscape, err = cmd.Flags().GetBool(flag.CSVBackslashEscape)
		if err != nil {
			return errors.Trace(err)
		}
		separator, err = cmd.Flags().GetString(flag.CSVSeparator)
		if err != nil {
			return errors.Trace(err)
		}
		delimiter, err = cmd.Flags().GetString(flag.CSVDelimiter)
		if err != nil {
			return errors.Trace(err)
		}
		trimLastSeparator, err = cmd.Flags().GetBool(flag.CSVTrimLastSeparator)
		if err != nil {
			return errors.Trace(err)
		}
	}

	cmd.Annotations[telemetry.ClusterID] = clusterID

	uploadFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = uploadFile.Close()
	}()

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
	if o.h.IOStreams.CanPrompt {
		uploadID, err = spinnerWaitUploadOp(ctx, o.h, uploader, input)
		if err != nil {
			return err
		}
	} else {
		uploadID, err = waitUploadOp(ctx, o.h, uploader, input)
		if err != nil {
			return err
		}
	}

	body := importOp.ImportServiceCreateImportBody{}
	err = body.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"importOptions": {
				"fileType": "%s",
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
			"source": {
				"local": {
					"uploadId": "%s",
					"targetDatabase": "%s",
					"targetTable": "%s"
				},
				"type": "LOCAL"
			}
			}`, fileType, uploadID, targetDatabase, targetTable)))
	if err != nil {
		return errors.Trace(err)
	}

	body.ImportOptions.CsvFormat.Separator = separator
	body.ImportOptions.CsvFormat.Delimiter = delimiter
	body.ImportOptions.CsvFormat.BackslashEscape = aws.Bool(backslashEscape)
	body.ImportOptions.CsvFormat.TrimLastSeparator = aws.Bool(trimLastSeparator)

	params := importOp.NewImportServiceCreateImportParams().WithClusterID(clusterID).
		WithBody(body).WithContext(ctx)
	if o.h.IOStreams.CanPrompt {
		err := spinnerWaitStartOp(ctx, o.h, d, params)
		if err != nil {
			return err
		}
	} else {
		err := waitStartOp(o.h, d, params)
		if err != nil {
			return err
		}
	}

	return nil
}

func initialLocalInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.Cursor.Style = config.FocusedStyle
		t.CharLimit = 0
		f := localImportField(i)

		switch f {
		case filePathIdx:
			t.Placeholder = "Local file path"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case databaseIdx:
			t.Placeholder = "Target database"
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
