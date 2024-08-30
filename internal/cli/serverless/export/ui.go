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
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

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
	flag.CSVSeparator:         "Input the csv separator: separator of each value in CSV files, skip to use default value (,)",
	flag.CSVDelimiter:         "Input the csv delimiter: delimiter of string type variables in CSV files, skip to use default value (\"). If you want to set empty string, please use non-interactive mode",
	flag.CSVNullValue:         "Input the csv null value: representation of null values in CSV files, skip to use default value (\\N). If you want to set empty string, please use non-interactive mode",
	flag.CSVSkipHeader:        "Input the csv skip header: export CSV files of the tables without header. Type `true` to skip header, others will not skip header",
}

func GetSelectedParquetCompression() (string, error) {
	compressions := make([]interface{}, 0, 4)
	compressions = append(compressions, "SNAPPY", "GZIP", "NONE")
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
	return compression.(string), nil
}

func GetSelectedTargetType() (TargetType, error) {
	targetTypes := make([]interface{}, 0, 2)
	targetTypes = append(targetTypes, TargetTypeLOCAL, TargetTypeS3, TargetTypeGCS, TargetTypeAZBLOB)
	model, err := ui.InitialSelectModel(targetTypes, "Choose the export target:")
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
	targetType := targetTypeModel.(ui.SelectModel).GetSelectedItem()
	if targetType == nil {
		return TargetTypeUnknown, errors.New("no export target selected")
	}
	return targetType.(TargetType), nil
}

func GetSelectedAuthType(target TargetType) (_ AuthType, err error) {
	var model *ui.SelectModel
	switch target {
	case TargetTypeS3:
		authTypes := make([]interface{}, 0, 2)
		authTypes = append(authTypes, AuthTypeS3RoleArn, AuthTypeS3AccessKey)
		model, err = ui.InitialSelectModel(authTypes, "Choose and input the S3 auth:")
		if err != nil {
			return "", errors.Trace(err)
		}
	case TargetTypeGCS:
		return AuthTypeGCSServiceAccountKey, nil
	case TargetTypeAZBLOB:
		return AuthTypeAzBlobSasToken, nil
	case TargetTypeLOCAL:
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
	return authType.(AuthType), nil
}

func GetSelectedFileType(filterType FilterType) (_ FileType, err error) {
	var model *ui.SelectModel
	switch filterType {
	case FilterSQL:
		fileTypes := make([]interface{}, 0, 2)
		fileTypes = append(fileTypes, FileTypeCSV, FileTypePARQUET)
		model, err = ui.InitialSelectModel(fileTypes, "Choose the exported file type:")
	default:
		fileTypes := make([]interface{}, 0, 3)
		fileTypes = append(fileTypes, FileTypeSQL, FileTypeCSV, FileTypePARQUET)
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
		return FileTypeUnknown, util.InterruptError
	}
	fileType := fileTypeModel.(ui.SelectModel).GetSelectedItem()
	if fileType == nil {
		return "", errors.New("no export file type selected")
	}
	return fileType.(FileType), nil
}

func GetSelectedCompression() (string, error) {
	compressions := make([]interface{}, 0, 4)
	compressions = append(compressions, "SNAPPY", "ZSTD", "NONE")
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
	return compression.(string), nil
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
