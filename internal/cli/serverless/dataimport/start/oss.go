// Copyright 2025 PingCAP, Inc.
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
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pingcap/errors"
	"github.com/spf13/cobra"
)

type OSSOpts struct {
	h           *internal.Helper
	interactive bool
	clusterId   string
}

func (o OSSOpts) Run(cmd *cobra.Command) error {
	ctx := cmd.Context()
	var ossUri, accessKeyID, secretAccessKey string
	var fileType imp.ImportFileTypeEnum
	var authType imp.ImportOSSAuthTypeEnum
	var format *imp.CSVFormat
	d, err := o.h.Client()
	if err != nil {
		return err
	}

	if o.interactive {
		// interactive mode
		authTypes := []interface{}{imp.IMPORTOSSAUTHTYPEENUM_ACCESS_KEY}
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
		authType = authTypeModel.(ui.SelectModel).Choices[authTypeModel.(ui.SelectModel).Selected].(imp.ImportOSSAuthTypeEnum)

		if authType == imp.IMPORTOSSAUTHTYPEENUM_ACCESS_KEY {
			inputs := []string{flag.OSSURI, flag.OSSAccessKeyID, flag.OSSSecretAccessKey}
			textInput, err := ui.InitialInputModel(inputs, inputDescription)
			if err != nil {
				return err
			}
			ossUri = textInput.Inputs[0].Value()
			if ossUri == "" {
				return errors.New("empty OSS URI")
			}
			accessKeyID = textInput.Inputs[1].Value()
			if accessKeyID == "" {
				return errors.New("empty OSS access key Id")
			}
			secretAccessKey = textInput.Inputs[2].Value()
			if secretAccessKey == "" {
				return errors.New("empty OSS secret access key")
			}
		} else {
			return fmt.Errorf("invalid auth type :%s", authType)
		}

		var fileTypes []interface{}
		for _, f := range imp.AllowedImportFileTypeEnumEnumValues {
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
		fileType = fileTypeModel.(ui.SelectModel).Choices[fileTypeModel.(ui.SelectModel).Selected].(imp.ImportFileTypeEnum)

		if fileType == imp.IMPORTFILETYPEENUM_CSV {
			format, err = getCSVFormat()
			if err != nil {
				return err
			}
		}
	} else {
		// non-interactive mode
		fileTypeStr, err := cmd.Flags().GetString(flag.FileType)
		if err != nil {
			return errors.Trace(err)
		}
		fileType = imp.ImportFileTypeEnum(fileTypeStr)
		if !fileType.IsValid() {
			return fmt.Errorf("file type \"%s\" is not supported, please use one of %q", fileTypeStr, imp.AllowedImportFileTypeEnumEnumValues)
		}
		ossUri, err = cmd.Flags().GetString(flag.OSSURI)
		if err != nil {
			return errors.Trace(err)
		}
		if ossUri == "" {
			return errors.New("empty OSS URI")
		}

		// optional flags
		if fileType == imp.IMPORTFILETYPEENUM_CSV {
			format, err = getCSVFlagValue(cmd)
			if err != nil {
				return errors.Trace(err)
			}
		}
		accessKeyID, err = cmd.Flags().GetString(flag.OSSAccessKeyID)
		if err != nil {
			return errors.Trace(err)
		}
		secretAccessKey, err = cmd.Flags().GetString(flag.OSSSecretAccessKey)
		if err != nil {
			return errors.Trace(err)
		}
		if accessKeyID != "" && secretAccessKey != "" {
			authType = imp.IMPORTOSSAUTHTYPEENUM_ACCESS_KEY
		} else {
			return fmt.Errorf("access key id and secret access key must be provided")
		}
	}

	source := imp.NewImportSource(imp.IMPORTSOURCETYPEENUM_OSS)
	source.Oss = imp.NewOSSSource(ossUri, authType)
	source.Oss.AuthType = authType
	source.Oss.AccessKey = &imp.OSSSourceAccessKey{
		Id:     accessKeyID,
		Secret: secretAccessKey,
	}

	options := imp.NewImportOptions(fileType)
	if fileType == imp.IMPORTFILETYPEENUM_CSV {
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
