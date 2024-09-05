// Copyright 2024 PingCAP, Inc.
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
	"slices"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/aws/s3"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	imp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/import"

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
	clusterId   string
}

func (o LocalOpts) SupportedFileTypes() []string {
	return []string{
		string(imp.IMPORTFILETYPEENUM_CSV),
	}
}

func (o LocalOpts) Run(cmd *cobra.Command) error {
	ctx := cmd.Context()
	var fileType, targetDatabase, targetTable, filePath string
	var format *imp.CSVFormat
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
		// interactive mode
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
		p = tea.NewProgram(o.initialInputModel())
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

		if fileType == string(imp.IMPORTFILETYPEENUM_CSV) {
			format, err = getCSVFormat()
			if err != nil {
				return err
			}
		}
	} else {
		// non-interactive mode
		fileType = cmd.Flag(flag.FileType).Value.String()
		if !slices.Contains(o.SupportedFileTypes(), fileType) {
			return fmt.Errorf("file type \"%s\" is not supported, please use one of %q", fileType, o.SupportedFileTypes())
		}
		targetDatabase = cmd.Flag(flag.LocalTargetDatabase).Value.String()
		targetTable = cmd.Flag(flag.LocalTargetTable).Value.String()
		filePath, err = cmd.Flags().GetString(flag.LocalFilePath)
		if err != nil {
			return errors.Trace(err)
		}
		if filePath == "" {
			return errors.New("required flag(s) \"local.file-path\" not set")
		}

		if fileType == string(imp.IMPORTFILETYPEENUM_CSV) {
			format, err = getCSVFlagValue(cmd)
			if err != nil {
				return errors.Trace(err)
			}
		}
	}

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
		ClusterID:     o.clusterId,
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

	source := imp.NewImportSource(imp.IMPORTSOURCETYPEENUM_LOCAL)
	source.Local = imp.NewLocalSource(uploadID, targetDatabase, targetTable)
	options := imp.NewImportOptions(imp.ImportFileTypeEnum(fileType))
	if fileType == string(imp.IMPORTFILETYPEENUM_CSV) {
		options.CsvFormat = format
	}
	body := imp.NewImportServiceCreateImportBody(*options, *source)

	if o.h.IOStreams.CanPrompt {
		err := spinnerWaitStartOp(ctx, o.h, d, o.clusterId, body)
		if err != nil {
			return err
		}
	} else {
		err := waitStartOp(ctx, o.h, d, o.clusterId, body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o LocalOpts) initialInputModel() ui.TextInputModel {
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

	fmt.Fprint(h.IOStreams.Out, color.GreenString("Start uploading...\n"))

	processModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if processModel.(ui.ProcessModel).Interrupted {
		return "", util.InterruptError
	}
	if processModel.(ui.ProcessModel).Err != nil {
		return "", processModel.(ui.ProcessModel).Err
	}

	fmt.Fprintln(h.IOStreams.Out, color.GreenString("File has been uploaded"))
	return uploadID, nil
}
