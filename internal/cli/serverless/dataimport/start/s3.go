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
	importOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pingcap/errors"
	"github.com/spf13/cobra"
)

var s3ImportField = map[string]int{
	flag.S3URI:            0,
	flag.S3TargetDatabase: 1,
}

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
		string(importModel.V1beta1ImportOptionsFileTypeCSV),
		string(importModel.V1beta1ImportOptionsFileTypeParquet),
		string(importModel.V1beta1ImportOptionsFileTypeSQL),
		string(importModel.V1beta1ImportOptionsFileTypeAuroraSnapshot),
	}
}

func (o S3Opts) Run(cmd *cobra.Command) error {
	ctx := cmd.Context()
	var clusterID, fileType, targetDatabase, separator, delimiter, s3Uri, s3Arn, accessKeyID, secretAccessKey string
	var backslashEscape, trimLastSeparator bool
	var authType importModel.V1beta1AuthType
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

		var fileTypes []interface{}
		for _, f := range o.SupportedFileTypes() {
			fileTypes = append(fileTypes, f)
		}
		model, err := ui.InitialSelectModel(fileTypes, "Choose the source file type:")
		if err != nil {
			return err
		}
		p := tea.NewProgram(model)
		fileTypeModel, err := p.Run()
		if err != nil {
			return errors.Trace(err)
		}
		if m, _ := fileTypeModel.(ui.SelectModel); m.Interrupted {
			return util.InterruptError
		}
		fileType = fileTypeModel.(ui.SelectModel).Choices[fileTypeModel.(ui.SelectModel).Selected].(string)

		// variables for input
		p = tea.NewProgram(o.initialInputModel())
		inputModel, err := p.Run()
		if err != nil {
			return errors.Trace(err)
		}
		if inputModel.(ui.TextInputModel).Interrupted {
			return util.InterruptError
		}
		s3Uri = inputModel.(ui.TextInputModel).Inputs[s3ImportField[flag.S3URI]].Value()
		if len(s3Uri) == 0 {
			return errors.New("S3 fold path is required")
		}
		targetDatabase = inputModel.(ui.TextInputModel).Inputs[s3ImportField[flag.S3TargetDatabase]].Value()
		if len(targetDatabase) == 0 {
			return errors.New("Target database is required")
		}

		authTypes := []interface{}{importModel.V1beta1AuthTypeROLEARN, importModel.V1beta1AuthTypeACCESSKEY}
		model, err = ui.InitialSelectModel(authTypes, "Choose the auth type:")
		if err != nil {
			return err
		}
		p = tea.NewProgram(model)
		authTypeModel, err := p.Run()
		if err != nil {
			return errors.Trace(err)
		}
		if m, _ := authTypeModel.(ui.SelectModel); m.Interrupted {
			return util.InterruptError
		}
		authType = authTypeModel.(ui.SelectModel).Choices[authTypeModel.(ui.SelectModel).Selected].(importModel.V1beta1AuthType)

		if authType == importModel.V1beta1AuthTypeROLEARN {
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
		} else if authType == importModel.V1beta1AuthTypeACCESSKEY {
			// variables for input
			p = tea.NewProgram(o.initialAccessKeyInputModel())
			inputModel, err := p.Run()
			if err != nil {
				return errors.Trace(err)
			}
			if inputModel.(ui.TextInputModel).Interrupted {
				return util.InterruptError
			}
			accessKeyID = inputModel.(ui.TextInputModel).Inputs[s3ImportField[flag.S3AccessKeyID]].Value()
			if len(accessKeyID) == 0 {
				return errors.New("S3 access key id is required")
			}
			secretAccessKey = inputModel.(ui.TextInputModel).Inputs[s3ImportField[flag.S3SecretAccessKey]].Value()
			if len(secretAccessKey) == 0 {
				return errors.New("S3 secret access key is required")
			}
		} else {
			return fmt.Errorf("invalid auth type :%s", authType)
		}

		separator, delimiter, backslashEscape, trimLastSeparator, err = getCSVFormat()
		if err != nil {
			return err
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
		targetDatabase, err = cmd.Flags().GetString(flag.S3TargetDatabase)
		if err != nil {
			return errors.Trace(err)
		}
		s3Uri, err = cmd.Flags().GetString(flag.S3URI)
		if err != nil {
			return errors.Trace(err)
		}

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
			authType = importModel.V1beta1AuthTypeROLEARN
		} else if accessKeyID != "" && secretAccessKey != "" {
			authType = importModel.V1beta1AuthTypeACCESSKEY
		} else {
			return fmt.Errorf("either role arn or access key id and secret access key must be provided")
		}
	}

	cmd.Annotations[telemetry.ClusterID] = clusterID

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
				"s3": {
					"s3Uri": "%s",
					"targetDatabase": "%s"
				},
				"type": "S3"
			}
			}`, fileType, s3Uri, targetDatabase)))
	if err != nil {
		return errors.Trace(err)
	}

	if authType == importModel.V1beta1AuthTypeROLEARN {
		body.Source.S3.Type = authType
		body.Source.S3.RoleArn = &importModel.V1beta1ImportSourceRoleArn{
			RoleArn: s3Arn,
		}
	} else if authType == importModel.V1beta1AuthTypeACCESSKEY {
		body.Source.S3.Type = authType
		body.Source.S3.AccessKey = &importModel.V1beta1ImportSourceAccessKey{
			ID:     accessKeyID,
			Secret: secretAccessKey,
		}
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

func (o S3Opts) initialInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(s3ImportField)),
	}

	var t textinput.Model
	for k, v := range s3ImportField {
		t = textinput.New()
		t.Cursor.Style = config.FocusedStyle
		t.CharLimit = 0

		switch k {
		case flag.S3URI:
			t.Placeholder = "S3 fold uri"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.S3TargetDatabase:
			t.Placeholder = "Target database"
		}

		m.Inputs[v] = t
	}

	return m
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
