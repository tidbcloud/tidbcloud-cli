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

package export

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	exportApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/client/export_service"
	exportModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/models"
)

type TargetType string

const (
	TargetTypeS3      TargetType = "S3"
	TargetTypeLOCAL   TargetType = "LOCAL"
	TargetTypeUnknown TargetType = "UNKNOWN"
)

type FileType string

const (
	FileTypeSQL     FileType = "SQL"
	FileTypeCSV     FileType = "CSV"
	FileTypeUnknown FileType = "UNKNOWN"
)

var S3InputFields = map[string]int{
	flag.BucketURI:       0,
	flag.AccessKeyID:     1,
	flag.SecretAccessKey: 2,
}

var FilterInputFields = map[string]int{
	flag.Database: 0,
	flag.Table:    1,
}

type CreateOpts struct {
	interactive bool
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func (c *CreateOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := c.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	// Mark required flags
	if !c.interactive {
		for _, fn := range flags {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a serverless cluster export",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create an export in interactive mode:
  $ %[1]s serverless export create

  Create an export with local type in non-interactive mode:
  $ %[1]s serverless export create -c <cluster-id> --database <database> --table <table>

  Create an export with s3 type in non-interactive mode:
  $ %[1]s serverless export create -c <cluster-id> --bucket-uri <bucket-uri> --access-key-id <access-key-id> --secret-access-key <secret-access-key>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterId string
			var bucketURI, accessKeyID, secretAccessKey, database, table, targetType, fileType, compression string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}

				cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterId = cluster.ID

				selectedTargetType, err := GetSelectedTargetType()
				if err != nil {
					return err
				}
				if selectedTargetType == TargetTypeUnknown {
					return errors.New("target type must be LOCAL or S3")
				}
				targetType = string(selectedTargetType)

				if selectedTargetType == TargetTypeS3 {
					s3InputModel, err := GetS3Input()
					if err != nil {
						return err
					}
					bucketURI = s3InputModel.(ui.TextInputModel).Inputs[S3InputFields[flag.BucketURI]].Value()
					if bucketURI == "" {
						return errors.New("bucket URI is required when target type is S3")
					}
					accessKeyID = s3InputModel.(ui.TextInputModel).Inputs[S3InputFields[flag.AccessKeyID]].Value()
					if accessKeyID == "" {
						return errors.New("accessKeyId is required when target type is S3")
					}
					secretAccessKey = s3InputModel.(ui.TextInputModel).Inputs[S3InputFields[flag.SecretAccessKey]].Value()
					if secretAccessKey == "" {
						return errors.New("secretAccessKey is required when target type is S3")
					}
				}

				selectedFileType, err := GetSelectedFileType()
				if err != nil {
					return err
				}
				if selectedFileType == FileTypeUnknown {
					return errors.New("file type must be LOCAL or S3")
				}
				fileType = string(selectedFileType)

				compression, err = GetSelectedCompression()
				if err != nil {
					return err
				}

				fmt.Fprintln(h.IOStreams.Out, color.BlueString("Please input the following options"))

				filterInputModel, err := GetFilterInput()
				if err != nil {
					return err
				}
				database = filterInputModel.(ui.TextInputModel).Inputs[FilterInputFields[flag.Database]].Value()
				if (database == "" || database == "*") && selectedTargetType == TargetTypeLOCAL {
					return errors.New("you must specify the database when target type is LOCAL")
				}
				table = filterInputModel.(ui.TextInputModel).Inputs[FilterInputFields[flag.Table]].Value()
			} else {
				// non-interactive mode, get values from flags
				var err error
				clusterId, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				targetType, err = cmd.Flags().GetString(flag.TargetType)
				if err != nil {
					return errors.Trace(err)
				}
				fileType, err = cmd.Flags().GetString(flag.FileType)
				if err != nil {
					return errors.Trace(err)
				}
				if targetType == string(TargetTypeS3) {
					bucketURI, err = cmd.Flags().GetString(flag.BucketURI)
					if err != nil {
						return errors.Trace(err)
					}
					if bucketURI == "" {
						return errors.New("bucket URI is required when target type is S3")
					}
					accessKeyID, err = cmd.Flags().GetString(flag.AccessKeyID)
					if err != nil {
						return errors.Trace(err)
					}
					if accessKeyID == "" {
						return errors.New("accessKeyId is required when target type is S3")
					}
					secretAccessKey, err = cmd.Flags().GetString(flag.SecretAccessKey)
					if err != nil {
						return errors.Trace(err)
					}
					if secretAccessKey == "" {
						return errors.New("secretAccessKey is required when target type is S3")
					}
				}
				database, err = cmd.Flags().GetString(flag.Database)
				if err != nil {
					return errors.Trace(err)
				}
				table, err = cmd.Flags().GetString(flag.Table)
				if err != nil {
					return errors.Trace(err)
				}
				if (database == "" || database == "*") && targetType == string(TargetTypeLOCAL) {
					return errors.New("you must specify the database when target type is LOCAL")
				}
				compression, err = cmd.Flags().GetString(flag.Compression)
				if err != nil {
					return errors.Trace(err)
				}
			}

			params := exportApi.NewExportServiceCreateExportParams().WithClusterID(clusterId).WithBody(
				exportApi.ExportServiceCreateExportBody{
					ExportOptions: &exportModel.V1beta1ExportOptions{
						Database: database,
						Table:    table,
						FileType: exportModel.ExportOptionsFileType(fileType),
					},
					Target: &exportModel.V1beta1Target{
						Type: exportModel.TargetTargetType(targetType),
						S3: &exportModel.TargetS3Target{
							BucketURI: bucketURI,
							AccessKey: &exportModel.S3TargetAccessKey{
								ID:     accessKeyID,
								Secret: secretAccessKey,
							},
						},
					},
				})
			if compression != "" {
				params.Body.ExportOptions.Compression = exportModel.ExportOptionsCompressionType(compression)
			}
			resp, err := d.CreateExport(params)
			if err != nil {
				return errors.Trace(err)
			}
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("export %s is running now", resp.Payload.ExportID))
			if err != nil {
				return err
			}
			return nil
		},
	}

	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster, in which the export will be created")
	createCmd.Flags().String(flag.Database, "*", "The database name you want to export")
	createCmd.Flags().String(flag.Table, "*", "The table name you want to export")
	createCmd.Flags().String(flag.FileType, "CSV", "The export file type. One of [\"CSV\" \"SQL\"]")
	createCmd.Flags().String(flag.TargetType, "LOCAL", "The export Target. One of [\"LOCAL\" \"S3\"]")
	createCmd.Flags().String(flag.BucketURI, "", "The bucket URI of the S3 bucket. Required when target type is S3")
	createCmd.Flags().String(flag.AccessKeyID, "", "The access key ID of the S3 bucket. Required when target type is S3")
	createCmd.Flags().String(flag.SecretAccessKey, "", "The secret access key of the S3 bucket. Required when target type is S3")
	createCmd.Flags().String(flag.Compression, "", "The compression algorithm of the export file. One of [\"gzip\" \"snappy\" \"zstd\" \"none\"]")
	return createCmd
}

func GetSelectedTargetType() (TargetType, error) {
	targetTypes := make([]interface{}, 0, 2)
	targetTypes = append(targetTypes, TargetTypeLOCAL, TargetTypeS3)
	model, err := ui.InitialSelectModel(targetTypes, "Choose export target")
	if err != nil {
		return TargetTypeUnknown, errors.Trace(err)
	}

	p := tea.NewProgram(model)
	targetTypeModel, err := p.Run()
	if err != nil {
		return TargetTypeUnknown, errors.Trace(err)
	}
	if m, _ := targetTypeModel.(ui.SelectModel); m.Interrupted {
		return TargetTypeUnknown, util.InterruptError
	}
	targetType := targetTypeModel.(ui.SelectModel).GetSelectedItem().(TargetType)
	return targetType, nil
}

func GetSelectedFileType() (FileType, error) {
	fileTypes := make([]interface{}, 0, 2)
	fileTypes = append(fileTypes, FileTypeSQL, FileTypeCSV)
	model, err := ui.InitialSelectModel(fileTypes, "Choose export file type")
	if err != nil {
		return FileTypeUnknown, errors.Trace(err)
	}

	p := tea.NewProgram(model)
	fileTypeModel, err := p.Run()
	if err != nil {
		return FileTypeUnknown, errors.Trace(err)
	}
	if m, _ := fileTypeModel.(ui.SelectModel); m.Interrupted {
		return FileTypeUnknown, util.InterruptError
	}
	fileType := fileTypeModel.(ui.SelectModel).GetSelectedItem().(FileType)
	return fileType, nil
}

func GetSelectedCompression() (string, error) {
	compressions := make([]interface{}, 0, 4)
	compressions = append(compressions, "GZIP", "SNAPPY", "ZSTD", "NONE")
	model, err := ui.InitialSelectModel(compressions, "Choose the compression algorithm")
	if err != nil {
		return "", errors.Trace(err)
	}

	p := tea.NewProgram(model)
	fileTypeModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := fileTypeModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	compression := fileTypeModel.(ui.SelectModel).GetSelectedItem().(string)
	return compression, nil
}

func initialS3InputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(S3InputFields)),
	}
	for k, v := range S3InputFields {
		t := textinput.New()
		switch k {
		case flag.BucketURI:
			t.Placeholder = "Bucket Uri"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.AccessKeyID:
			t.Placeholder = "Access Key ID"
		case flag.SecretAccessKey:
			t.Placeholder = "Secret Access Key"
		}
		m.Inputs[v] = t
	}
	return m
}

func GetS3Input() (tea.Model, error) {
	p := tea.NewProgram(initialS3InputModel())
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}

func initialFilterInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(FilterInputFields)),
	}
	for k, v := range FilterInputFields {
		t := textinput.New()
		switch k {
		case flag.Database:
			t.Placeholder = "Database Name. Press Enter to skip and export all databases (Can't skip in LOCAL target)."
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.Table:
			t.Placeholder = "Table Name. Press Enter to skip and export all tables."
		}
		m.Inputs[v] = t
	}
	return m
}

func GetFilterInput() (tea.Model, error) {
	p := tea.NewProgram(initialFilterInputModel())
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}
