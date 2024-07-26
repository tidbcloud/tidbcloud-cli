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
	stdErr "errors"
	"fmt"
	"slices"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	importOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type AzBlobOpts struct {
	h           *internal.Helper
	interactive bool
}

func (o AzBlobOpts) SupportedFileTypes() []string {
	return []string{
		string(importModel.V1beta1ImportOptionsFileTypeCSV),
		string(importModel.V1beta1ImportOptionsFileTypePARQUET),
		string(importModel.V1beta1ImportOptionsFileTypeSQL),
		string(importModel.V1beta1ImportOptionsFileTypeAURORASNAPSHOT),
	}
}

func (o AzBlobOpts) Run(cmd *cobra.Command) error {
	ctx := cmd.Context()
	var clusterID, fileType, sasUrl string
	var authType importModel.V1beta1AzureBlobSourceAuthType
	var format *importModel.V1beta1CSVFormat
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

		authTypes := []interface{}{importModel.V1beta1AzureBlobSourceAuthTypeSASTOKEN}
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
		authType = authTypeModel.(ui.SelectModel).Choices[authTypeModel.(ui.SelectModel).Selected].(importModel.V1beta1AzureBlobSourceAuthType)

		if authType == importModel.V1beta1AzureBlobSourceAuthTypeSASTOKEN {
			input := &survey.Input{
				Message: "Please input the blob sas url:",
			}
			err = survey.AskOne(input, &sasUrl, survey.WithValidator(survey.Required))
			if err != nil {
				if stdErr.Is(err, terminal.InterruptErr) {
					return util.InterruptError
				} else {
					return err
				}
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

		if fileType == string(importModel.V1beta1ImportOptionsFileTypeCSV) {
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

		sasUrl, err = cmd.Flags().GetString(flag.AzureBlobSASUrl)
		if err != nil {
			return errors.Trace(err)
		}
		if sasUrl == "" {
			return fmt.Errorf("azure blob sas url is required")
		}
		authType = importModel.V1beta1AzureBlobSourceAuthTypeSASTOKEN

		// optional flags
		format, err = getCSVFlagValue(cmd)
		if err != nil {
			return errors.Trace(err)
		}
	}

	cmd.Annotations[telemetry.ClusterID] = clusterID

	body := &importModel.ImportServiceCreateImportBody{}
	err = body.UnmarshalBinary([]byte(fmt.Sprintf(`{
			"importOptions": {
				"fileType": "%s"
			},
			"source": {
				"type": "AzBlob"
			}
			}`, fileType)))
	if err != nil {
		return errors.Trace(err)
	}
	body.ImportOptions.CsvFormat = format

	body.Source.AzureBlob = &importModel.V1beta1AzureBlobSource{
		Type:   &authType,
		SasURL: sasUrl,
	}

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
