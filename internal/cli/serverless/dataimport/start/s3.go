// Copyright 2024 PingCAP, Inc.
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

package start

import (
	stdErr "errors"
	"fmt"
	"slices"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	imp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/import"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pingcap/errors"
	"github.com/spf13/cobra"
)

var accessKeyImportField = map[string]int{
	flag.S3AccessKeyID:     0,
	flag.S3SecretAccessKey: 1,
}

type S3Opts struct {
	h           *internal.Helper
	interactive bool
}

func (o S3Opts) SupportedFileTypes() []string {
	return []string{
		string(imp.IMPORTFILETYPEENUM_CSV),
		string(imp.IMPORTFILETYPEENUM_PARQUET),
		string(imp.IMPORTFILETYPEENUM_SQL),
		string(imp.IMPORTFILETYPEENUM_AURORA_SNAPSHOT),
	}
}

func (o S3Opts) Run(cmd *cobra.Command) error {
	ctx := cmd.Context()
	var clusterID, fileType, s3Uri, s3Arn, accessKeyID, secretAccessKey string
	var authType imp.ImportS3AuthTypeEnum
	var format *imp.CSVFormat
	d, err := o.h.Client()
	if err != nil {
		return err
	}

	if o.interactive {
		cmd.Annotations[telemetry.InteractiveMode] = "true"
		if !o.h.IOStreams.CanPrompt {
			return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
		}

		// interactive mode
		project, err := cloud.GetSelectedProject(ctx, o.h.QueryPageSize, d)
		if err != nil {
			return err
		}

		cluster, err := cloud.GetSelectedCluster(ctx, project.ID, o.h.QueryPageSize, d)
		if err != nil {
			return err
		}
		clusterID = cluster.ID

		input := &survey.Input{
			Message: "Please input the s3 fold uri:",
		}
		err = survey.AskOne(input, &s3Uri, survey.WithValidator(survey.Required))
		if err != nil {
			if stdErr.Is(err, terminal.InterruptErr) {
				return util.InterruptError
			} else {
				return err
			}
		}

		authTypes := []interface{}{imp.IMPORTS3AUTHTYPEENUM_ROLE_ARN, imp.IMPORTS3AUTHTYPEENUM_ACCESS_KEY}
		model, err := ui.InitialSelectModel(authTypes, "Choose the auth type:")
		if err != nil {
			return err
		}
		p := tea.NewProgram(model)
		authTypeModel, err := p.Run()
		if err != nil {
			return errors.Trace(err)
		}
		if m, _ := authTypeModel.(ui.SelectModel); m.Interrupted {
			return util.InterruptError
		}
		authType = authTypeModel.(ui.SelectModel).Choices[authTypeModel.(ui.SelectModel).Selected].(imp.ImportS3AuthTypeEnum)

		if authType == imp.IMPORTS3AUTHTYPEENUM_ROLE_ARN {
			input := &survey.Input{
				Message: "Please input the arn:",
			}
			err = survey.AskOne(input, &s3Arn, survey.WithValidator(survey.Required))
			if err != nil {
				if stdErr.Is(err, terminal.InterruptErr) {
					return util.InterruptError
				} else {
					return err
				}
			}
		} else if authType == imp.IMPORTS3AUTHTYPEENUM_ACCESS_KEY {
			// variables for input
			p = tea.NewProgram(o.initialAccessKeyInputModel())
			inputModel, err := p.Run()
			if err != nil {
				return errors.Trace(err)
			}
			if inputModel.(ui.TextInputModel).Interrupted {
				return util.InterruptError
			}
			accessKeyID = inputModel.(ui.TextInputModel).Inputs[accessKeyImportField[flag.S3AccessKeyID]].Value()
			if len(accessKeyID) == 0 {
				return errors.New("S3 access key id is required")
			}
			secretAccessKey = inputModel.(ui.TextInputModel).Inputs[accessKeyImportField[flag.S3SecretAccessKey]].Value()
			if len(secretAccessKey) == 0 {
				return errors.New("S3 secret access key is required")
			}
		} else {
			return fmt.Errorf("invalid auth type :%s", authType)
		}

		var fileTypes []interface{}
		for _, f := range o.SupportedFileTypes() {
			fileTypes = append(fileTypes, f)
		}
		model, err = ui.InitialSelectModel(fileTypes, "Choose the source file type:")
		if err != nil {
			return err
		}
		p = tea.NewProgram(model)
		fileTypeModel, err := p.Run()
		if err != nil {
			return errors.Trace(err)
		}
		if m, _ := fileTypeModel.(ui.SelectModel); m.Interrupted {
			return util.InterruptError
		}
		fileType = fileTypeModel.(ui.SelectModel).Choices[fileTypeModel.(ui.SelectModel).Selected].(string)

		if fileType == string(imp.IMPORTFILETYPEENUM_CSV) {
			format, err = getCSVFormat()
			if err != nil {
				return err
			}
		}
	} else {
		// non-interactive mode
		clusterID, err = cmd.Flags().GetString(flag.ClusterID)
		if err != nil {
			return errors.Trace(err)
		}
		fileType, err = cmd.Flags().GetString(flag.FileType)
		if err != nil {
			return errors.Trace(err)
		}
		if !slices.Contains(o.SupportedFileTypes(), fileType) {
			return fmt.Errorf("file type \"%s\" is not supported, please use one of %q", fileType, o.SupportedFileTypes())
		}
		s3Uri, err = cmd.Flags().GetString(flag.S3URI)
		if err != nil {
			return errors.Trace(err)
		}

		// optional flags
		format, err = getCSVFlagValue(cmd)
		if err != nil {
			return errors.Trace(err)
		}
		s3Arn, err = cmd.Flags().GetString(flag.S3RoleArn)
		if err != nil {
			return errors.Trace(err)
		}
		accessKeyID, err = cmd.Flags().GetString(flag.S3AccessKeyID)
		if err != nil {
			return errors.Trace(err)
		}
		secretAccessKey, err = cmd.Flags().GetString(flag.S3SecretAccessKey)
		if err != nil {
			return errors.Trace(err)
		}
		if s3Arn != "" {
			authType = imp.IMPORTS3AUTHTYPEENUM_ROLE_ARN
		} else if accessKeyID != "" && secretAccessKey != "" {
			authType = imp.IMPORTS3AUTHTYPEENUM_ACCESS_KEY
		} else {
			return fmt.Errorf("either role arn or access key id and secret access key must be provided")
		}
	}

	cmd.Annotations[telemetry.ClusterID] = clusterID

	source := imp.NewImportSource(imp.IMPORTSOURCETYPEENUM_S3)
	source.S3 = imp.NewS3Source(s3Uri, authType)
	if authType == imp.IMPORTS3AUTHTYPEENUM_ROLE_ARN {
		source.S3.AuthType = authType
		source.S3.RoleArn = &s3Arn
	} else {
		source.S3.AuthType = authType
		source.S3.AccessKey = &imp.S3SourceAccessKey{
			Id:     accessKeyID,
			Secret: secretAccessKey,
		}
	}
	options := imp.NewImportOptions(imp.ImportFileTypeEnum(fileType))
	options.CsvFormat = format
	body := imp.NewImportServiceCreateImportBody(*options, *source)

	if o.h.IOStreams.CanPrompt {
		err := spinnerWaitStartOp(ctx, o.h, d, clusterID, body)
		if err != nil {
			return err
		}
	} else {
		err := waitStartOp(ctx, o.h, d, clusterID, body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o S3Opts) initialAccessKeyInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(accessKeyImportField)),
	}

	var t textinput.Model
	for k, v := range accessKeyImportField {
		t = textinput.New()
		t.Cursor.Style = config.FocusedStyle
		t.CharLimit = 0

		switch k {
		case flag.S3AccessKeyID:
			t.Placeholder = "S3 access key id"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.S3SecretAccessKey:
			t.Placeholder = "S3 secret access key"
		}

		m.Inputs[v] = t
	}

	return m
}
