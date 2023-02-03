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
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type s3ImportField int

const (
	awsRoleArnIdx s3ImportField = iota
	sourceUrlIdx
)

type S3Opts struct {
	interactive bool
}

func (c S3Opts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
		flag.AwsRoleArn,
		flag.DataFormat,
		flag.SourceUrl,
	}
}

func (c S3Opts) SupportedDataFormats() []string {
	return []string{
		string(importModel.OpenapiDataFormatCSV),
		string(importModel.OpenapiDataFormatSQLFile),
		string(importModel.OpenapiDataFormatParquet),
		string(importModel.OpenapiDataFormatAuroraSnapshot),
	}
}

func S3Cmd(h *internal.Helper) *cobra.Command {
	opts := S3Opts{
		interactive: true,
	}

	var s3Cmd = &cobra.Command{
		Use:   "s3",
		Short: "Import files from Amazon S3 into TiDB Cloud",
		Example: fmt.Sprintf(`  Start an import task in interactive mode:
  $ %[1]s import start s3

  Start an import task in non-interactive mode:
  $ %[1]s import start s3 --project-id <project-id> --cluster-id <cluster-id> --aws-role-arn <aws-role-arn> --data-format <data-format> --source-url <source-url>

  Start an import task with custom CSV format:
  $ %[1]s import start s3 --project-id <project-id> --cluster-id <cluster-id> --aws-role-arn <aws-role-arn> --data-format CSV --source-url <source-url> --separator \" --delimiter \' --backslash-escape=false --trim-last-separator=true
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
			var projectID, clusterID, awsRoleArn, dataFormat, sourceUrl, separator, delimiter string
			var backslashEscape, trimLastSeparator bool

			d, err := h.Client()
			if err != nil {
				return err
			}

			if opts.interactive {
				cmd.Annotations = map[string]string{telemetry.InteractiveMode: "true"}
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
				p = tea.NewProgram(initialS3InputModel())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return nil
				}

				awsRoleArn = inputModel.(ui.TextInputModel).Inputs[awsRoleArnIdx].Value()
				if len(awsRoleArn) == 0 {
					return errors.New("AWS role ARN is required")
				}
				sourceUrl = inputModel.(ui.TextInputModel).Inputs[sourceUrlIdx].Value()
				if len(sourceUrl) == 0 {
					return errors.New("Source url is required")
				}

				if dataFormat == string(importModel.OpenapiDataFormatCSV) {
					separator, delimiter, backslashEscape, trimLastSeparator, err = getCSVFormat()
					if err != nil {
						return err
					}
				}
			} else {
				// non-interactive mode
				projectID = cmd.Flag(flag.ProjectID).Value.String()
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				awsRoleArn = cmd.Flag(flag.AwsRoleArn).Value.String()
				dataFormat = cmd.Flag(flag.DataFormat).Value.String()
				if !util.ElemInSlice(opts.SupportedDataFormats(), dataFormat) {
					return fmt.Errorf("data format %s is not supported, please use one of %q", dataFormat, opts.SupportedDataFormats())
				}
				sourceUrl = cmd.Flag(flag.SourceUrl).Value.String()

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

			body := importOp.CreateImportBody{}
			err = body.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"aws_role_arn": "%s",
			"data_format": "%s",
			"source_url": "%s",
			"type": "S3",
			"csv_format": {
                "separator": ",",
				"delimiter": "\"",
				"header": true,
				"backslash_escape": true,
				"null": "\\N",
				"trim_last_separator": false,
				"not_null": false
			}
			}`, awsRoleArn, dataFormat, sourceUrl)))
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

	s3Cmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "Project ID")
	s3Cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	s3Cmd.Flags().String(flag.AwsRoleArn, "", "AWS S3 IAM Role ARN")
	s3Cmd.Flags().String(flag.DataFormat, "", fmt.Sprintf("Data format, one of %q", opts.SupportedDataFormats()))
	s3Cmd.Flags().String(flag.SourceUrl, "", "The S3 path where the source data file is stored")
	return s3Cmd
}

func initialS3InputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = config.FocusedStyle
		t.CharLimit = 0
		f := s3ImportField(i)

		switch f {
		case awsRoleArnIdx:
			t.Placeholder = "AWS role ARN"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case sourceUrlIdx:
			t.Placeholder = "Source url"
		}

		m.Inputs[i] = t
	}

	return m
}
