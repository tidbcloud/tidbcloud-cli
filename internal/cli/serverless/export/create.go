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
	"slices"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"
)

type TargetType string

const (
	TargetTypeS3      TargetType = "S3"
	TargetTypeLOCAL   TargetType = "LOCAL"
	TargetTypeGCS     TargetType = "GCS"
	TargetTypeAZBLOB  TargetType = "AZURE_BLOB"
	TargetTypeUnknown TargetType = "UNKNOWN"
)

type FileType string

const (
	FileTypeSQL     FileType = "SQL"
	FileTypeCSV     FileType = "CSV"
	FileTypePARQUET FileType = "PARQUET"
	FileTypeUnknown FileType = "UNKNOWN"
)

type AuthType string

const (
	AuthTypeS3AccessKey          AuthType = "S3AccessKey"
	AuthTypeS3RoleArn            AuthType = "S3RoleArn"
	AuthTypeGCSServiceAccountKey AuthType = "GCSServiceAccountKey"
	AuthTypeAzBlobSasToken       AuthType = "AzBlobSasToken"
)

var (
	supportedFileType           = []string{string(FileTypeSQL), string(FileTypeCSV), string(FileTypePARQUET)}
	supportedTargetType         = []string{string(TargetTypeS3), string(TargetTypeLOCAL), string(TargetTypeGCS), string(TargetTypeAZBLOB)}
	supportedCompression        = []string{"GZIP", "SNAPPY", "ZSTD", "NONE"}
	supportedParquetCompression = []string{"GZIP", "SNAPPY", "ZSTD", "NONE"}
)

const (
	CSVSeparatorDefaultValue = ","
	CSVDelimiterDefaultValue = "\""
	CSVNullValueDefaultValue = "\\N"
)

type CreateOpts struct {
	interactive bool
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.FileType,
		flag.TargetType,
		flag.S3URI,
		flag.S3AccessKeyID,
		flag.S3SecretAccessKey,
		flag.Compression,
		flag.SQL,
		flag.TableFilter,
		flag.TableWhere,
		flag.CSVDelimiter,
		flag.CSVNullValue,
		flag.CSVSkipHeader,
		flag.CSVSeparator,
		flag.S3RoleArn,
		flag.GCSURI,
		flag.GCSServiceAccountKey,
		flag.AzureBlobURI,
		flag.AzureBlobSASToken,
		flag.ParquetCompression,
	}
}

func (c CreateOpts) RequiredFlags() []string {
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
		for _, fn := range c.RequiredFlags() {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	var force bool
	opts := CreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Export data from a TiDB Serverless cluster",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create an export in interactive mode:
  $ %[1]s serverless export create

  Export all data with local type in non-interactive mode:
  $ %[1]s serverless export create -c <cluster-id>

  Export all data with S3 type in non-interactive mode:
  $ %[1]s serverless export create -c <cluster-id> --target-type S3 --s3.uri <s3-uri> --s3.access-key-id <access-key-id> --s3.secret-access-key <secret-access-key>

  Export all data and customize csv format in non-interactive mode:
  $ %[1]s serverless export create -c <cluster-id> --file-type CSV --csv.separator ";" --csv.delimiter "\"" --csv.null-value 'NaN' --csv.skip-header

  Export test.t1 and test.t2 in non-interactive mode:
  $ %[1]s serverless export create -c <cluster-id> --filter 'test.t1,test.t2'

  Export tables with special characters, for example, if you want to export %[2]stest,%[2]s.%[2]st1%[2]s and %[2]s"test%[2]s.%[2]st1%[2]s:
  $ %[1]s serverless export create -c <cluster-id> --filter '"%[2]stest1,%[2]s.t1","%[2]s""test%[2]s.t1"'`,
			config.CliName, "`"),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			// options
			var targetType, fileType, compression, clusterId, sql, where string
			var patterns []string
			// csv format
			var csvSeparator, csvDelimiter, csvNullValue string
			var csvSkipHeader bool
			// parquet options
			var parquetCompression string
			// s3
			var s3URI, accessKeyID, secretAccessKey, s3RoleArn string
			// gcs
			var gcsURI, gcsServiceAccountKey string
			// azure
			var azBlobURI, azBlobSasToken string

			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}

				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterId = cluster.ID

				// target
				selectedTargetType, err := GetSelectedTargetType()
				if err != nil {
					return err
				}
				targetType = string(selectedTargetType)
				selectedAuthType, err := GetSelectedAuthType(selectedTargetType)
				if err != nil {
					return err
				}
				switch selectedAuthType {
				case AuthTypeS3AccessKey:
					inputs := []string{flag.S3URI, flag.S3AccessKeyID, flag.S3SecretAccessKey}
					textInput, err := InitialInputModel(inputs)
					if err != nil {
						return err
					}
					s3URI = textInput.Inputs[0].Value()
					if s3URI == "" {
						return errors.New("empty S3 URI")
					}
					accessKeyID = textInput.Inputs[1].Value()
					if accessKeyID == "" {
						return errors.New("empty S3 access key Id")
					}
					secretAccessKey = textInput.Inputs[2].Value()
					if secretAccessKey == "" {
						return errors.New("empty S3 secret access key")
					}
				case AuthTypeS3RoleArn:
					inputs := []string{flag.S3URI, flag.S3RoleArn}
					textInput, err := InitialInputModel(inputs)
					if err != nil {
						return err
					}
					s3URI = textInput.Inputs[0].Value()
					if s3URI == "" {
						return errors.New("empty S3 URI")
					}
					s3RoleArn = textInput.Inputs[1].Value()
					if s3RoleArn == "" {
						return errors.New("empty S3 role arn")
					}
				case AuthTypeGCSServiceAccountKey:
					inputs := []string{flag.GCSURI, flag.GCSServiceAccountKey}
					textInput, err := InitialInputModel(inputs)
					if err != nil {
						return err
					}
					gcsURI = textInput.Inputs[0].Value()
					if gcsURI == "" {
						return errors.New("empty GCS URI")
					}
					gcsServiceAccountKey = textInput.Inputs[1].Value()
					if gcsServiceAccountKey == "" {
						return errors.New("empty GCS service account key")
					}
				case AuthTypeAzBlobSasToken:
					inputs := []string{flag.AzureBlobURI, flag.AzureBlobSASToken}
					textInput, err := InitialInputModel(inputs)
					if err != nil {
						return err
					}
					azBlobURI = textInput.Inputs[0].Value()
					if azBlobURI == "" {
						return errors.New("empty Azure Blob URI")
					}
					azBlobSasToken = textInput.Inputs[1].Value()
					if azBlobSasToken == "" {
						return errors.New("empty Azure Blob SAS token")
					}
				}

				// Export options, including: filter, file type, compression
				filterType, err := GetSelectedFilterType()
				if err != nil {
					return err
				}
				switch filterType {
				case FilterNone:
				case FilterSQL:
					inputs := []string{flag.SQL}
					textInput, err := InitialInputModel(inputs)
					if err != nil {
						return err
					}
					sql = textInput.Inputs[0].Value()
					if sql == "" {
						return errors.New("sql is empty")
					}
				case FilterTable:
					fmt.Fprintln(h.IOStreams.Out, color.BlueString("Please input the following options, require at least one field"))
					inputs := []string{flag.TableFilter, flag.TableWhere}
					textInput, err := InitialInputModel(inputs)
					if err != nil {
						return err
					}
					patternString := textInput.Inputs[0].Value()
					patterns, err = util.StringSliceConv(patternString)
					if err != nil {
						return err
					}
					where = textInput.Inputs[1].Value()
					if len(patterns) == 0 && where == "" {
						return errors.New("both patterns and where are empty, require at least one field")
					}
				}

				selectedFileType, err := GetSelectedFileType(filterType)
				if err != nil {
					return err
				}
				fileType = string(selectedFileType)
				switch fileType {
				case string(FileTypeCSV):
					customCSVFormat := false
					prompt := &survey.Confirm{
						Message: "Do you want to customize the CSV format",
						Default: false,
					}
					err = survey.AskOne(prompt, &customCSVFormat)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					if customCSVFormat {
						inputs := []string{flag.CSVSeparator, flag.CSVDelimiter, flag.CSVNullValue, flag.CSVSkipHeader}
						textInput, err := InitialInputModel(inputs)
						if err != nil {
							return err
						}
						csvSeparator = textInput.Inputs[0].Value()
						csvDelimiter = textInput.Inputs[1].Value()
						csvNullValue = textInput.Inputs[2].Value()
						skipHeader := textInput.Inputs[3].Value()
						if skipHeader == "true" {
							csvSkipHeader = true
						} else {
							csvSkipHeader = false
						}
					}
					if csvSeparator == "" {
						csvSeparator = CSVSeparatorDefaultValue
					}
					if csvDelimiter == "" {
						csvDelimiter = CSVDelimiterDefaultValue
					}
					if csvNullValue == "" {
						csvNullValue = CSVNullValueDefaultValue
					}
				case string(FileTypePARQUET):
					customParquetCompression := false
					prompt := &survey.Confirm{
						Message: "Do you want change the default parquet compression algorithm ZSTD",
						Default: false,
					}
					err = survey.AskOne(prompt, &customParquetCompression)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}

					if customParquetCompression {
						parquetCompression, err = GetSelectedParquetCompression()
						if err != nil {
							return err
						}
					}
				}

				if fileType != string(FileTypePARQUET) {
					changeCompression := false
					prompt := &survey.Confirm{
						Message: "Do you want to change the default compression algorithm GZIP",
						Default: false,
					}
					err = survey.AskOne(prompt, &changeCompression)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					if changeCompression {
						compression, err = GetSelectedCompression()
						if err != nil {
							return err
						}
					}
				}
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
				if targetType != "" && !slices.Contains(supportedTargetType, strings.ToUpper(targetType)) {
					return errors.New("unsupported target type: " + targetType)
				}
				fileType, err = cmd.Flags().GetString(flag.FileType)
				if err != nil {
					return errors.Trace(err)
				}
				if fileType != "" && !slices.Contains(supportedFileType, strings.ToUpper(fileType)) {
					return errors.New("unsupported file type: " + fileType)
				}
				switch strings.ToUpper(targetType) {
				case string(TargetTypeS3):
					s3URI, err = cmd.Flags().GetString(flag.S3URI)
					if err != nil {
						return errors.Trace(err)
					}
					if s3URI == "" {
						return errors.New("S3 URI is required when target type is S3")
					}
					accessKeyID, err = cmd.Flags().GetString(flag.S3AccessKeyID)
					if err != nil {
						return errors.Trace(err)
					}
					secretAccessKey, err = cmd.Flags().GetString(flag.S3SecretAccessKey)
					if err != nil {
						return errors.Trace(err)
					}
					s3RoleArn, err = cmd.Flags().GetString(flag.S3RoleArn)
					if err != nil {
						return errors.Trace(err)
					}
					if s3RoleArn == "" && (accessKeyID == "" || secretAccessKey == "") {
						return errors.New("missing S3 auth information, require either role arn or access key id and secret access key")
					}
				case string(TargetTypeGCS):
					gcsURI, err = cmd.Flags().GetString(flag.GCSURI)
					if err != nil {
						return errors.Trace(err)
					}
					if gcsURI == "" {
						return errors.New("GCS URI is required when target type is GCS")
					}
					gcsServiceAccountKey, err = cmd.Flags().GetString(flag.GCSServiceAccountKey)
					if err != nil {
						return errors.Trace(err)
					}
					if gcsServiceAccountKey == "" {
						return errors.New("GCS service account key is required when target type is GCS")
					}
				case string(TargetTypeAZBLOB):
					azBlobURI, err = cmd.Flags().GetString(flag.AzureBlobURI)
					if err != nil {
						return errors.Trace(err)
					}
					if azBlobURI == "" {
						return errors.New("Azure Blob URI is required when target type is AZURE_BLOB")
					}
					azBlobSasToken, err = cmd.Flags().GetString(flag.AzureBlobSASToken)
					if err != nil {
						return errors.Trace(err)
					}
					if azBlobSasToken == "" {
						return errors.New("Azure Blob SAS token is required when target type is AZURE_BLOB")
					}
				}

				switch strings.ToUpper(fileType) {
				case string(FileTypeCSV):
					csvSeparator, err = cmd.Flags().GetString(flag.CSVSeparator)
					if err != nil {
						return errors.Trace(err)
					}
					if len(csvSeparator) == 0 {
						return errors.New("csv separator can not be empty")
					}
					csvDelimiter, err = cmd.Flags().GetString(flag.CSVDelimiter)
					if err != nil {
						return errors.Trace(err)
					}
					csvNullValue, err = cmd.Flags().GetString(flag.CSVNullValue)
					if err != nil {
						return errors.Trace(err)
					}
					csvSkipHeader, err = cmd.Flags().GetBool(flag.CSVSkipHeader)
					if err != nil {
						return errors.Trace(err)
					}
				case string(FileTypePARQUET):
					parquetCompression, err = cmd.Flags().GetString(flag.ParquetCompression)
					if err != nil {
						return errors.Trace(err)
					}
					if parquetCompression != "" && !slices.Contains(supportedParquetCompression, strings.ToUpper(parquetCompression)) {
						return errors.New("unsupported parquet compression: " + parquetCompression)
					}
				}
				if strings.ToUpper(fileType) != string(FileTypePARQUET) {
					compression, err = cmd.Flags().GetString(flag.Compression)
					if err != nil {
						return errors.Trace(err)
					}
					if compression != "" && !slices.Contains(supportedCompression, strings.ToUpper(compression)) {
						return errors.New("unsupported compression: " + compression)
					}
				}
				sql, err = cmd.Flags().GetString(flag.SQL)
				if err != nil {
					return errors.Trace(err)
				}
				where, err = cmd.Flags().GetString(flag.TableWhere)
				if err != nil {
					return errors.Trace(err)
				}
				patterns, err = cmd.Flags().GetStringSlice(flag.TableFilter)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if !opts.interactive && sql == "" && len(patterns) == 0 && !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to create export")
				}

				confirmationMessage := fmt.Sprintf("%s %s %s %s", color.BlueString("You will export the whole cluster."), color.BlueString("Please type"), color.HiBlueString(confirmed), color.BlueString("to continue:"))
				prompt := &survey.Input{
					Message: confirmationMessage,
				}
				var userInput string
				err := survey.AskOne(prompt, &userInput)
				if err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					} else {
						return err
					}
				}
				if userInput != confirmed {
					return errors.New("incorrect confirm string entered, skipping create")
				}
			}

			// build param to create export
			fileTypeEnum := export.ExportFileTypeEnum(strings.ToUpper(fileType))
			targetTypeEnum := export.ExportTargetTypeEnum(strings.ToUpper(targetType))
			params := &export.ExportServiceCreateExportBody{
				ExportOptions: &export.ExportOptions{
					FileType: &fileTypeEnum,
				},
				Target: &export.ExportTarget{
					Type: &targetTypeEnum,
				},
			}
			// add target
			switch targetTypeEnum {
			case export.EXPORTTARGETTYPEENUM_S3:
				if s3RoleArn != "" {
					params.Target.S3 = &export.S3Target{
						Uri:      &s3URI,
						RoleArn:  &s3RoleArn,
						AuthType: export.EXPORTS3AUTHTYPEENUM_ROLE_ARN,
					}
				} else {
					params.Target.S3 = &export.S3Target{
						Uri:      &s3URI,
						AuthType: export.EXPORTS3AUTHTYPEENUM_ACCESS_KEY,
						AccessKey: &export.S3TargetAccessKey{
							Id:     accessKeyID,
							Secret: secretAccessKey,
						},
					}
				}
			case export.EXPORTTARGETTYPEENUM_GCS:
				params.Target.Gcs = &export.GCSTarget{
					Uri:               gcsURI,
					AuthType:          export.EXPORTGCSAUTHTYPEENUM_SERVICE_ACCOUNT_KEY,
					ServiceAccountKey: &gcsServiceAccountKey,
				}
			case export.EXPORTTARGETTYPEENUM_AZURE_BLOB:
				params.Target.AzureBlob = &export.AzureBlobTarget{
					Uri:      azBlobURI,
					AuthType: export.EXPORTAZUREBLOBAUTHTYPEENUM_SAS_TOKEN,
					SasToken: &azBlobSasToken,
				}
			}
			// add compression
			if compression != "" {
				compressionEnum := export.ExportCompressionTypeEnum(strings.ToUpper(compression))
				params.ExportOptions.Compression = &compressionEnum
			}
			// add filter
			if sql != "" {
				params.ExportOptions.Filter = &export.ExportOptionsFilter{
					Sql: &sql,
				}
			}
			if len(patterns) > 0 || where != "" {
				params.ExportOptions.Filter = &export.ExportOptionsFilter{
					Table: &export.ExportOptionsFilterTable{
						Where:    &where,
						Patterns: patterns,
					},
				}
			}
			// add file type
			switch strings.ToUpper(fileType) {
			case string(FileTypeCSV):
				params.ExportOptions.CsvFormat = &export.ExportOptionsCSVFormat{
					Separator:  &csvSeparator,
					Delimiter:  *export.NewNullableString(&csvDelimiter),
					NullValue:  *export.NewNullableString(&csvNullValue),
					SkipHeader: &csvSkipHeader,
				}
			case string(FileTypePARQUET):
				if parquetCompression != "" {
					c := export.ExportParquetCompressionTypeEnum(strings.ToUpper(parquetCompression))
					params.ExportOptions.ParquetFormat = &export.ExportOptionsParquetFormat{
						Compression: &c,
					}
				}
			}

			resp, err := d.CreateExport(ctx, clusterId, params)
			if err != nil {
				return errors.Trace(err)
			}
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("export %s is running now", *resp.ExportId))
			if err != nil {
				return err
			}
			return nil
		},
	}

	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster, in which the export will be created.")
	createCmd.Flags().String(flag.FileType, "CSV", "The export file type. One of [\"CSV\" \"SQL\" \"PARQUET\"].")
	createCmd.Flags().String(flag.TargetType, "LOCAL", "The export target. One of [\"LOCAL\" \"S3\" \"GCS\" \"AZURE_BLOB\"].")
	createCmd.Flags().String(flag.S3URI, "", "The S3 URI in s3://<bucket>/<path> format. Required when target type is S3.")
	createCmd.Flags().String(flag.S3AccessKeyID, "", "The access key ID of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	createCmd.Flags().String(flag.S3SecretAccessKey, "", "The secret access key of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	createCmd.Flags().String(flag.Compression, "GZIP", "The compression algorithm of the export file. One of [\"GZIP\" \"SNAPPY\" \"ZSTD\" \"NONE\"]. Only take effect when file type is not PARQUET.")
	createCmd.Flags().StringSlice(flag.TableFilter, nil, "Specify the exported table(s) with table filter patterns. See https://docs.pingcap.com/tidb/stable/table-filter to learn table filter.")
	createCmd.Flags().String(flag.TableWhere, "", "Filter the exported table(s) with the where condition.")
	createCmd.Flags().String(flag.SQL, "", "Filter the exported data with SQL SELECT statement.")
	createCmd.Flags().BoolVar(&force, flag.Force, false, "Create without confirmation. You need to confirm when you want to export the whole cluster in non-interactive mode.")
	createCmd.Flags().String(flag.CSVDelimiter, CSVDelimiterDefaultValue, "Delimiter of string type variables in CSV files.")
	createCmd.Flags().String(flag.CSVSeparator, CSVSeparatorDefaultValue, "Separator of each value in CSV files.")
	createCmd.Flags().String(flag.CSVNullValue, CSVNullValueDefaultValue, "Representation of null values in CSV files.")
	createCmd.Flags().Bool(flag.CSVSkipHeader, false, "Export CSV files of the tables without header.")
	createCmd.Flags().String(flag.S3RoleArn, "", "The role arn of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	createCmd.Flags().String(flag.GCSURI, "", "The GCS URI in gcs://<bucket>/<path> format. Required when target type is GCS.")
	createCmd.Flags().String(flag.GCSServiceAccountKey, "", "The base64 encoded service account key of GCS.")
	createCmd.Flags().String(flag.AzureBlobURI, "", "The Azure Blob URI in azure://<account>.blob.core.windows.net/<container>/<path> format. Required when target type is AZURE_BLOB.")
	createCmd.Flags().String(flag.AzureBlobSASToken, "", "The SAS token of Azure Blob.")
	createCmd.Flags().String(flag.ParquetCompression, "ZSTD", "The parquet compression algorithm. One of [\"GZIP\" \"SNAPPY\" \"ZSTD\" \"NONE\"].")

	createCmd.MarkFlagsMutuallyExclusive(flag.TableFilter, flag.SQL)
	createCmd.MarkFlagsMutuallyExclusive(flag.TableWhere, flag.SQL)
	createCmd.MarkFlagsMutuallyExclusive(flag.S3RoleArn, flag.S3AccessKeyID)
	createCmd.MarkFlagsMutuallyExclusive(flag.S3RoleArn, flag.S3SecretAccessKey)
	return createCmd
}
