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
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
)

var inputDescription = map[string]string{
	flag.S3URI:                "Input your S3 URI in s3://<bucket>/<path> format",
	flag.S3AccessKeyID:        "Input your S3 access key id",
	flag.S3SecretAccessKey:    "Input your S3 secret access key",
	flag.S3RoleArn:            "Input your S3 role arn",
	flag.AzureBlobURI:         "Input your Azure Blob URI in azure://<account>.blob.core.windows.net/<container>/<path> format",
	flag.AzureBlobSASToken:    "Input your Azure Blob SAS token",
	flag.GCSURI:               "Input your GCS URI in gcs://<bucket>/<path> format",
	flag.GCSServiceAccountKey: "Input your base64 encoded GCS service account key",
	flag.SQL:                  "Input the SELECT SQL statement",
	flag.TableFilter:          "Input the table filter patterns (comma separated). Example: database.table,database.*,`database-1`.`table-1`",
	flag.TableWhere:           "Input the where clause which will apply to all filtered tables. Example: id > 10",
	flag.CSVSeparator:         "Input the CSV separator: separator of each value in CSV files, skip to use default value (,)",
	flag.CSVDelimiter:         "Input the CSV delimiter: delimiter of string type variables in CSV files, skip to use default value (\"). If you want to set empty string, please use non-interactive mode",
	flag.CSVNullValue:         "Input the CSV null value: representation of null values in CSV files, skip to use default value (\\N). If you want to set empty string, please use non-interactive mode",
	flag.CSVSkipHeader:        "Input the CSV skip header: export CSV files of the tables without header. Type `true` to skip header, others will not skip header",
	flag.DisplayName:          "Input the name of export. You can skip and use the default name SNAPSHOT_${snapshot_time} by pressing Enter",
}

func GetSelectedParquetCompression() (export.ExportParquetCompressionTypeEnum, error) {
	compressions := make([]interface{}, 0, 3)
	compressions = append(compressions, export.EXPORTPARQUETCOMPRESSIONTYPEENUM_SNAPPY, export.EXPORTPARQUETCOMPRESSIONTYPEENUM_GZIP, export.EXPORTPARQUETCOMPRESSIONTYPEENUM_NONE)
	model, err := ui.InitialSelectModel(compressions, "Choose the parquet compression algorithm:")
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
	compression := fileTypeModel.(ui.SelectModel).GetSelectedItem()
	if compression == nil {
		return "", errors.New("no compression algorithm selected")
	}
	return compression.(export.ExportParquetCompressionTypeEnum), nil
}

func GetSelectedTargetType() (export.ExportTargetTypeEnum, error) {
	targetTypes := make([]interface{}, 0, 4)
	targetTypes = append(targetTypes, export.EXPORTTARGETTYPEENUM_LOCAL, export.EXPORTTARGETTYPEENUM_S3, export.EXPORTTARGETTYPEENUM_GCS, export.EXPORTTARGETTYPEENUM_AZURE_BLOB)
	model, err := ui.InitialSelectModel(targetTypes, "Choose the export target:")
	if err != nil {
		return "", errors.Trace(err)
	}

	p := tea.NewProgram(model)
	targetTypeModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := targetTypeModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	targetType := targetTypeModel.(ui.SelectModel).GetSelectedItem()
	if targetType == nil {
		return "", errors.New("no export target selected")
	}
	return targetType.(export.ExportTargetTypeEnum), nil
}

func GetSelectedAuthType(target export.ExportTargetTypeEnum) (_ string, err error) {
	var model *ui.SelectModel
	switch target {
	case export.EXPORTTARGETTYPEENUM_S3:
		authTypes := make([]interface{}, 0, 2)
		authTypes = append(authTypes, string(export.EXPORTS3AUTHTYPEENUM_ROLE_ARN), string(export.EXPORTS3AUTHTYPEENUM_ACCESS_KEY))
		model, err = ui.InitialSelectModel(authTypes, "Choose and input the S3 auth:")
		if err != nil {
			return "", errors.Trace(err)
		}
	case export.EXPORTTARGETTYPEENUM_GCS:
		return string(export.EXPORTGCSAUTHTYPEENUM_SERVICE_ACCOUNT_KEY), nil
	case export.EXPORTTARGETTYPEENUM_AZURE_BLOB:
		return string(export.EXPORTAZUREBLOBAUTHTYPEENUM_SAS_TOKEN), nil
	case export.EXPORTTARGETTYPEENUM_LOCAL:
		return "", nil
	}
	if model == nil {
		return "", errors.New("unknown auth type")
	}
	p := tea.NewProgram(model)
	authTypeModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := authTypeModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	authType := authTypeModel.(ui.SelectModel).GetSelectedItem()
	if authType == nil {
		return "", errors.New("no auth type selected")
	}
	return authType.(string), nil
}

func GetSelectedFileType(filterType FilterType) (_ export.ExportFileTypeEnum, err error) {
	var model *ui.SelectModel
	switch filterType {
	case FilterSQL:
		fileTypes := make([]interface{}, 0, 2)
		fileTypes = append(fileTypes, export.EXPORTFILETYPEENUM_CSV, export.EXPORTFILETYPEENUM_PARQUET)
		model, err = ui.InitialSelectModel(fileTypes, "Choose the exported file type:")
	default:
		fileTypes := make([]interface{}, 0, 3)
		fileTypes = append(fileTypes, export.EXPORTFILETYPEENUM_SQL, export.EXPORTFILETYPEENUM_CSV, export.EXPORTFILETYPEENUM_PARQUET)
		model, err = ui.InitialSelectModel(fileTypes, "Choose the exported file type:")
	}
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
	fileType := fileTypeModel.(ui.SelectModel).GetSelectedItem()
	if fileType == nil {
		return "", errors.New("no export file type selected")
	}
	return fileType.(export.ExportFileTypeEnum), nil
}

func GetSelectedCompression() (export.ExportCompressionTypeEnum, error) {
	compressions := make([]interface{}, 0, 3)
	compressions = append(compressions, export.EXPORTCOMPRESSIONTYPEENUM_SNAPPY, export.EXPORTCOMPRESSIONTYPEENUM_ZSTD, export.EXPORTCOMPRESSIONTYPEENUM_NONE)
	model, err := ui.InitialSelectModel(compressions, "Choose the compression algorithm:")
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
	compression := fileTypeModel.(ui.SelectModel).GetSelectedItem()
	if compression == nil {
		return "", errors.New("no compression algorithm selected")
	}
	return compression.(export.ExportCompressionTypeEnum), nil
}

type FilterType string

const (
	FilterNone  FilterType = "Export all data"
	FilterTable FilterType = "Table (export specific database(s) or table(s) using table filter)"
	FilterSQL   FilterType = "SQL (export data using SELECT SQL statement)"
)

func GetSelectedFilterType() (FilterType, error) {
	filterTypes := make([]interface{}, 0, 3)

	filterTypes = append(filterTypes, FilterTable, FilterSQL, FilterNone)
	model, err := ui.InitialSelectModel(filterTypes, "Choose how to filter your data:")
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
	filterType := fileTypeModel.(ui.SelectModel).GetSelectedItem()
	if filterType == nil {
		return "", errors.New("no filter type selected")
	}
	return filterType.(FilterType), nil
}

func initialDownloadPathInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(DownloadPathInputFields)),
	}
	for k, v := range DownloadPathInputFields {
		t := textinput.New()
		switch k {
		case flag.OutputPath:
			t.Placeholder = "Where you want to download the file. Press Enter to skip and download to the current directory"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		}
		m.Inputs[v] = t
	}
	return m
}

func GetDownloadPathInput() (tea.Model, error) {
	p := tea.NewProgram(initialDownloadPathInputModel())
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}
