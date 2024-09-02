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
	"fmt"
	"slices"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	imp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/import"

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
		string(imp.IMPORTFILETYPEENUM_CSV),
		string(imp.IMPORTFILETYPEENUM_PARQUET),
		string(imp.IMPORTFILETYPEENUM_SQL),
		string(imp.IMPORTFILETYPEENUM_AURORA_SNAPSHOT),
	}
}

func (o AzBlobOpts) Run(cmd *cobra.Command) error {
	ctx := cmd.Context()
	var clusterID, fileType, uri, sasToken string
	var authType imp.ImportAzureBlobAuthTypeEnum
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

		authTypes := []interface{}{imp.IMPORTAZUREBLOBAUTHTYPEENUM_SAS_TOKEN}
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
		authType = authTypeModel.(ui.SelectModel).Choices[authTypeModel.(ui.SelectModel).Selected].(imp.ImportAzureBlobAuthTypeEnum)

		if authType == imp.IMPORTAZUREBLOBAUTHTYPEENUM_SAS_TOKEN {
			inputs := []string{flag.AzureBlobURI, flag.AzureBlobSASToken}
			textInput, err := ui.InitialInputModel(inputs, inputDescription)
			if err != nil {
				return err
			}
			uri = textInput.Inputs[0].Value()
			if uri == "" {
				return errors.New("empty Azure Blob URI")
			}
			sasToken = textInput.Inputs[1].Value()
			if sasToken == "" {
				return errors.New("empty Azure Blob SAS token")
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

		uri, err = cmd.Flags().GetString(flag.AzureBlobURI)
		if err != nil {
			return errors.Trace(err)
		}
		if uri == "" {
			return fmt.Errorf("empty Azure Blob URI")
		}

		sasToken, err = cmd.Flags().GetString(flag.AzureBlobSASToken)
		if err != nil {
			return errors.Trace(err)
		}
		if sasToken == "" {
			return fmt.Errorf("empty Azure Blob SAS token")
		}
		authType = imp.IMPORTAZUREBLOBAUTHTYPEENUM_SAS_TOKEN

		// optional flags
		format, err = getCSVFlagValue(cmd)
		if err != nil {
			return errors.Trace(err)
		}
	}

	cmd.Annotations[telemetry.ClusterID] = clusterID

	source := imp.NewImportSource(imp.IMPORTSOURCETYPEENUM_AZURE_BLOB)
	source.AzureBlob = imp.NewAzureBlobSource(authType, uri)
	source.AzureBlob.AuthType = authType
	source.AzureBlob.SasToken = &sasToken
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
