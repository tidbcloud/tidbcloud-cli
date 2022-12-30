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

package dataimport

import (
	"fmt"
	"os"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
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

type startImportField int

const (
	awsRoleArnIdx startImportField = iota
	sourceUrlIdx
)

type StartOpts struct {
	interactive bool
}

func (c StartOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
		flag.AwsRoleArn,
		flag.DataFormat,
		flag.SourceUrl,
	}
}

func (c StartOpts) SupportedDataFormats() []string {
	return []string{
		string(importModel.OpenapiDataFormatCSV),
		string(importModel.OpenapiDataFormatSQLFile),
		string(importModel.OpenapiDataFormatParquet),
		string(importModel.OpenapiDataFormatAuroraSnapshot),
	}
}

func StartCmd(h *internal.Helper) *cobra.Command {
	opts := StartOpts{
		interactive: true,
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start a data import task",
		Example: fmt.Sprintf(`  Start an import task in interactive mode:
  $ %[1]s import start

  Start an import task in non-interactive mode:
  $ %[1]s import start --project-id <project-id> --cluster-id <cluster-id> --aws-role-arn <aws-role-arn> --data-format <data-format> --source-url <source-url>`,
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
			var projectID, clusterID, awsRoleArn, dataFormat, sourceUrl string
			d, err := h.Client()
			if err != nil {
				return err
			}

			if opts.interactive {
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
				p = tea.NewProgram(initialStartInputModel())
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
					return errors.New("source url is required")
				}
			} else {
				// non-interactive mode
				projectID = cmd.Flag(flag.ProjectID).Value.String()
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				awsRoleArn = cmd.Flag(flag.AwsRoleArn).Value.String()
				dataFormat = cmd.Flag(flag.DataFormat).Value.String()
				if !util.ElemInSlice(opts.SupportedDataFormats(), dataFormat) {
					return fmt.Errorf("data format %s is not supported, please use one of CSV, SqlFile, Parquet, AuroraSnapshot", dataFormat)
				}
				sourceUrl = cmd.Flag(flag.SourceUrl).Value.String()
			}

			body := importOp.CreateImportBody{}
			err = body.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"aws_role_arn": "%s",
			"data_format": "%s",
			"source_url": "%s",
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

	startCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "Project ID")
	startCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	startCmd.Flags().String(flag.AwsRoleArn, "", "AWS S3 IAM Role ARN")
	startCmd.Flags().String(flag.DataFormat, "", "Data format, one of CSV, SqlFile, Parquet, AuroraSnapshot")
	startCmd.Flags().String(flag.SourceUrl, "", "The S3 path where the source data file is stored")
	return startCmd
}

func waitStartOp(h *internal.Helper, d cloud.TiDBCloudClient, params *importOp.CreateImportParams) error {
	fmt.Fprintf(h.IOStreams.Out, "... Starting the import task\n")
	res, err := d.CreateImport(params)
	if err != nil {
		return err
	}

	fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", *(res.Payload.ID)))
	return nil
}

func spinnerWaitStartOp(h *internal.Helper, d cloud.TiDBCloudClient, params *importOp.CreateImportParams) error {
	task := func() tea.Msg {
		errChan := make(chan error)

		go func() {
			res, err := d.CreateImport(params)
			if err != nil {
				errChan <- err
				return
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", *(res.Payload.ID)))
			errChan <- nil
		}()

		ticker := time.NewTicker(1 * time.Second)
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

func initialStartInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = config.FocusedStyle
		t.CharLimit = 0
		f := startImportField(i)

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
