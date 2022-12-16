package dataimport

import (
	"fmt"
	"os"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
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

func StartCmd(h *internal.Helper) *cobra.Command {
	opts := StartOpts{}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start a data import task",
		Example: fmt.Sprintf(`  Start an import task in interactive mode:
  $ %[1]s import start

  Start an import task in non-interactive mode:
  $ %[1]s import start --project-id <project-id> --cluster-name <cluster-name> --aws-role-arn <aws-role-arn> --data-format <data-format> --source-url <source-url>`,
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

				dataFormats := []interface{}{importModel.OpenapiDataFormatCSV,
					importModel.OpenapiDataFormatSQLFile, importModel.OpenapiDataFormatParquet,
					importModel.OpenapiDataFormatAuroraSnapshot}
				model, err := ui.InitialSelectModel(dataFormats, "Choose the cloud region:")
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
				sourceUrl = cmd.Flag(flag.SourceUrl).Value.String()
			}

			body := importOp.CreateImportBody{}
			err = body.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"aws_role_arn": "%s",
			"data_format": "%s",
			"source_url": "%s"
			}`, awsRoleArn, dataFormat, sourceUrl)))
			if err != nil {
				return errors.Trace(err)
			}

			params := importOp.NewCreateImportParams().WithProjectID(projectID).WithClusterID(clusterID).
				WithBody(body)
			res, err := d.CreateImport(params)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprint(h.IOStreams.Out, color.GreenString("Import task %s starts.", *res.Payload.ID))
			return nil
		},
	}

	return startCmd
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
