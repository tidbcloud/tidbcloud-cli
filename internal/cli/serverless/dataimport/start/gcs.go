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
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	imp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/import"

	"github.com/aws/aws-sdk-go-v2/aws"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type GCSOpts struct {
	h           *internal.Helper
	interactive bool
	clusterId   string
}

func (o GCSOpts) SupportedFileTypes() []string {
	return []string{
		string(imp.IMPORTFILETYPEENUM_CSV),
		string(imp.IMPORTFILETYPEENUM_PARQUET),
		string(imp.IMPORTFILETYPEENUM_SQL),
		string(imp.IMPORTFILETYPEENUM_AURORA_SNAPSHOT),
	}
}

func (o GCSOpts) Run(cmd *cobra.Command) error {
	ctx := cmd.Context()
	var fileType, gcsUri, accountKey string
	var authType imp.ImportGcsAuthTypeEnum
	var format *imp.CSVFormat
	d, err := o.h.Client()
	if err != nil {
		return err
	}

	if o.interactive {
		// interactive mode
		authTypes := []interface{}{imp.IMPORTGCSAUTHTYPEENUM_SERVICE_ACCOUNT_KEY}
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
		authType = authTypeModel.(ui.SelectModel).Choices[authTypeModel.(ui.SelectModel).Selected].(imp.ImportGcsAuthTypeEnum)

		if authType == imp.IMPORTGCSAUTHTYPEENUM_SERVICE_ACCOUNT_KEY {
			inputs := []string{flag.GCSURI, flag.GCSServiceAccountKey}
			textInput, err := ui.InitialInputModel(inputs, inputDescription)
			if err != nil {
				return err
			}
			gcsUri = textInput.Inputs[0].Value()
			if gcsUri == "" {
				return errors.New("empty GCS URI")
			}
			accountKey = textInput.Inputs[1].Value()
			if accountKey == "" {
				return errors.New("empty GCS service account key")
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
		fileType, err = cmd.Flags().GetString(flag.FileType)
		if err != nil {
			return errors.Trace(err)
		}
		if !slices.Contains(o.SupportedFileTypes(), fileType) {
			return fmt.Errorf("file type \"%s\" is not supported, please use one of %q", fileType, o.SupportedFileTypes())
		}
		gcsUri, err = cmd.Flags().GetString(flag.GCSURI)
		if err != nil {
			return errors.Trace(err)
		}
		if gcsUri == "" {
			return fmt.Errorf("empty GCS URI")
		}

		accountKey, err = cmd.Flags().GetString(flag.GCSServiceAccountKey)
		if err != nil {
			return errors.Trace(err)
		}
		if accountKey == "" {
			return fmt.Errorf("empty GCS service account key")
		}
		authType = imp.IMPORTGCSAUTHTYPEENUM_SERVICE_ACCOUNT_KEY

		// optional flags
		if fileType == string(imp.IMPORTFILETYPEENUM_CSV) {
			format, err = getCSVFlagValue(cmd)
			if err != nil {
				return errors.Trace(err)
			}
		}
	}

	source := imp.NewImportSource(imp.IMPORTSOURCETYPEENUM_GCS)
	source.Gcs = imp.NewGCSSource(gcsUri, authType)
	source.Gcs.ServiceAccountKey = aws.String(accountKey)
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
