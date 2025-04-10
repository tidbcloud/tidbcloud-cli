// Copyright 2025 PingCAP, Inc.
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
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal/ui"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"
)

const (
	CSVSeparatorDefaultValue       = ","
	CSVDelimiterDefaultValue       = "\""
	CSVNullValueDefaultValue       = "\\N"
	CSVSkipHeaderDefaultValue      = false
	CompressionDefaultValue        = export.EXPORTCOMPRESSIONTYPEENUM_GZIP
	ParquetCompressionDefaultValue = export.EXPORTPARQUETCOMPRESSIONTYPEENUM_ZSTD
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
		flag.DisplayName,
		flag.OSSURI,
		flag.OSSAccessKeyID,
		flag.OSSAccessKeySecret,
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
		Short: "Export data from a TiDB Cloud Serverless cluster",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create an export in interactive mode:
  $ %[1]s serverless export create

  Export all data with local type in non-interactive mode:
  $ %[1]s serverless export create -c <cluster-id>

  Export all data with S3 type in non-interactive mode:
  $ %[1]s serverless export create -c <cluster-id> --target-type S3 --s3.uri <s3-uri> --s3.access-key-id <access-key-id> --s3.secret-access-key <secret-access-key>

  Export all data and customize CSV format in non-interactive mode:
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
			var targetType export.ExportTargetTypeEnum
			var fileType export.ExportFileTypeEnum
			var compression export.ExportCompressionTypeEnum
			var clusterId, sql, where string
			var patterns []string
			// csv format
			var csvSeparator, csvDelimiter, csvNullValue string
			var csvSkipHeader bool
			// parquet options
			var parquetCompression export.ExportParquetCompressionTypeEnum
			// s3
			var s3URI, accessKeyID, secretAccessKey, s3RoleArn string
			// gcs
			var gcsURI, gcsServiceAccountKey string
			// azure
			var azBlobURI, azBlobSasToken string
			// oss
			var ossURI, ossAccessKeyID, ossAccessKeySecret string
			var displayName string

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

				// display name
				fmt.Fprintln(h.IOStreams.Out, color.HiGreenString("Input the display name (optional):"))
				inputs := []string{flag.DisplayName}
				textInput, err := ui.InitialInputModel(inputs, inputDescription)
				if err != nil {
					return err
				}
				displayName = textInput.Inputs[0].Value()

				// target
				targetType, err = GetSelectedTargetType()
				if err != nil {
					return err
				}
				selectedAuthType, err := GetSelectedAuthType(targetType, *cluster.CloudProvider)
				if err != nil {
					return err
				}
				switch selectedAuthType {
				// Both S3 and OSS supports ACCESS_KEY
				case string(export.EXPORTS3AUTHTYPEENUM_ACCESS_KEY):
					if targetType == export.EXPORTTARGETTYPEENUM_S3 {
						inputs := []string{flag.S3URI, flag.S3AccessKeyID, flag.S3SecretAccessKey}
						textInput, err := ui.InitialInputModel(inputs, inputDescription)
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
					}
					if targetType == export.EXPORTTARGETTYPEENUM_OSS {
						inputs := []string{flag.OSSURI, flag.OSSAccessKeyID, flag.OSSAccessKeySecret}
						textInput, err := ui.InitialInputModel(inputs, inputDescription)
						if err != nil {
							return err
						}
						ossURI = textInput.Inputs[0].Value()
						if ossURI == "" {
							return errors.New("empty OSS URI")
						}
						ossAccessKeyID = textInput.Inputs[1].Value()
						if ossAccessKeyID == "" {
							return errors.New("empty OSS access key Id")
						}
						ossAccessKeySecret = textInput.Inputs[2].Value()
						if ossAccessKeySecret == "" {
							return errors.New("empty OSS access key secret")
						}
					}
				case string(export.EXPORTS3AUTHTYPEENUM_ROLE_ARN):
					inputs := []string{flag.S3URI, flag.S3RoleArn}
					textInput, err := ui.InitialInputModel(inputs, inputDescription)
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
				case string(export.EXPORTGCSAUTHTYPEENUM_SERVICE_ACCOUNT_KEY):
					inputs := []string{flag.GCSURI}
					textInput, err := ui.InitialInputModel(inputs, inputDescription)
					if err != nil {
						return err
					}
					gcsURI = textInput.Inputs[0].Value()
					if gcsURI == "" {
						return errors.New("empty GCS URI")
					}
					areaInput, err := ui.InitialTextAreaModel(inputDescription[flag.GCSServiceAccountKey])
					if err != nil {
						return errors.Trace(err)
					}
					gcsServiceAccountKey = areaInput.Textarea.Value()
					if gcsServiceAccountKey == "" {
						return errors.New("empty GCS service account key")
					}
				case string(export.EXPORTAZUREBLOBAUTHTYPEENUM_SAS_TOKEN):
					inputs := []string{flag.AzureBlobURI, flag.AzureBlobSASToken}
					textInput, err := ui.InitialInputModel(inputs, inputDescription)
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
					textInput, err := ui.InitialInputModel(inputs, inputDescription)
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
					textInput, err := ui.InitialInputModel(inputs, inputDescription)
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

				fileType, err = GetSelectedFileType(filterType)
				if err != nil {
					return err
				}
				switch fileType {
				case export.EXPORTFILETYPEENUM_CSV:
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
						textInput, err := ui.InitialInputModel(inputs, inputDescription)
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
				case export.EXPORTFILETYPEENUM_PARQUET:
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

				if fileType != export.EXPORTFILETYPEENUM_PARQUET {
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
				displayName, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}
				clusterId, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				targetTypeStr, err := cmd.Flags().GetString(flag.TargetType)
				if err != nil {
					return errors.Trace(err)
				}
				targetType = export.ExportTargetTypeEnum(strings.ToUpper(targetTypeStr))
				if targetType != "" && !targetType.IsValid() {
					return errors.New("unsupported target type: " + targetTypeStr)
				}
				fileTypeStr, err := cmd.Flags().GetString(flag.FileType)
				if err != nil {
					return errors.Trace(err)
				}
				fileType = export.ExportFileTypeEnum(strings.ToUpper(fileTypeStr))
				if fileType != "" && !fileType.IsValid() {
					return errors.New("unsupported file type: " + fileTypeStr)
				}

				switch targetType {
				case export.EXPORTTARGETTYPEENUM_S3:
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
				case export.EXPORTTARGETTYPEENUM_GCS:
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
				case export.EXPORTTARGETTYPEENUM_AZURE_BLOB:
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
				case export.EXPORTTARGETTYPEENUM_OSS:
					ossURI, err = cmd.Flags().GetString(flag.OSSURI)
					if err != nil {
						return errors.Trace(err)
					}
					if ossURI == "" {
						return errors.New("OSS URI is required when target type is OSS")
					}
					ossAccessKeyID, err = cmd.Flags().GetString(flag.OSSAccessKeyID)
					if err != nil {
						return errors.Trace(err)
					}
					ossAccessKeySecret, err = cmd.Flags().GetString(flag.OSSAccessKeySecret)
					if err != nil {
						return errors.Trace(err)
					}
					if ossAccessKeyID == "" || ossAccessKeySecret == "" {
						return errors.New("OSS access key id and access key secret are required when target type is OSS")
					}
				}

				compressionStr, err := cmd.Flags().GetString(flag.Compression)
				if err != nil {
					return errors.Trace(err)
				}
				compression = export.ExportCompressionTypeEnum(strings.ToUpper(compressionStr))
				if compression != "" && !compression.IsValid() {
					return errors.New("unsupported compression: " + compressionStr)
				}

				switch fileType {
				case export.EXPORTFILETYPEENUM_CSV:
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
					if csvDelimiter == "" && !cmd.Flag(flag.CSVDelimiter).Changed {
						csvDelimiter = CSVDelimiterDefaultValue
					}
					csvNullValue, err = cmd.Flags().GetString(flag.CSVNullValue)
					if err != nil {
						return errors.Trace(err)
					}
					if csvNullValue == "" && !cmd.Flag(flag.CSVNullValue).Changed {
						csvNullValue = CSVNullValueDefaultValue
					}
					csvSkipHeader, err = cmd.Flags().GetBool(flag.CSVSkipHeader)
					if err != nil {
						return errors.Trace(err)
					}
				case export.EXPORTFILETYPEENUM_PARQUET:
					parquetCompressionStr, err := cmd.Flags().GetString(flag.ParquetCompression)
					if err != nil {
						return errors.Trace(err)
					}
					parquetCompression = export.ExportParquetCompressionTypeEnum(strings.ToUpper(parquetCompressionStr))
					if parquetCompression != "" && !parquetCompression.IsValid() {
						return errors.New("unsupported parquet compression: " + parquetCompressionStr)
					}
					if compression != "" {
						return errors.New("--compression is not supported when file type is parquet, please use --parquet.compression instead")
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

			// apply default values
			if fileType == export.EXPORTFILETYPEENUM_PARQUET {
				if parquetCompression == "" {
					parquetCompression = ParquetCompressionDefaultValue
				}
			} else if compression == "" {
				compression = CompressionDefaultValue
			}
			// build param to create export
			params := &export.ExportServiceCreateExportBody{
				ExportOptions: &export.ExportOptions{
					FileType: &fileType,
				},
				Target: &export.ExportTarget{
					Type: &targetType,
				},
			}
			if displayName != "" {
				params.DisplayName = &displayName
			}
			// add target
			switch targetType {
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
			case export.EXPORTTARGETTYPEENUM_OSS:
				params.Target.Oss = &export.OSSTarget{
					Uri:      ossURI,
					AuthType: export.EXPORTOSSAUTHTYPEENUM_ACCESS_KEY,
					AccessKey: &export.OSSTargetAccessKey{
						Id:     ossAccessKeyID,
						Secret: ossAccessKeySecret,
					},
				}
			}
			// add compression
			if compression != "" {
				params.ExportOptions.Compression = &compression
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
			switch fileType {
			case export.EXPORTFILETYPEENUM_CSV:
				params.ExportOptions.CsvFormat = &export.ExportOptionsCSVFormat{
					Separator:  &csvSeparator,
					Delimiter:  *export.NewNullableString(&csvDelimiter),
					NullValue:  *export.NewNullableString(&csvNullValue),
					SkipHeader: &csvSkipHeader,
				}
			case export.EXPORTFILETYPEENUM_PARQUET:
				params.ExportOptions.ParquetFormat = &export.ExportOptionsParquetFormat{
					Compression: &parquetCompression,
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
	createCmd.Flags().String(flag.FileType, "CSV", fmt.Sprintf("The export file type. One of %q.", export.AllowedExportFileTypeEnumEnumValues))
	createCmd.Flags().String(flag.TargetType, "LOCAL", fmt.Sprintf("The export target. One of %q.", export.AllowedExportTargetTypeEnumEnumValues))
	createCmd.Flags().String(flag.S3URI, "", "The S3 URI in s3://<bucket>/<path> format. Required when target type is S3.")
	createCmd.Flags().String(flag.S3AccessKeyID, "", "The access key ID of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	createCmd.Flags().String(flag.S3SecretAccessKey, "", "The secret access key of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	createCmd.Flags().String(flag.Compression, "", fmt.Sprintf("The compression algorithm of the export file. One of %q.", export.AllowedExportCompressionTypeEnumEnumValues))
	createCmd.Flags().StringSlice(flag.TableFilter, nil, "Specify the exported table(s) with table filter patterns. See https://docs.pingcap.com/tidb/stable/table-filter to learn table filter.")
	createCmd.Flags().String(flag.TableWhere, "", "Filter the exported table(s) with the where condition.")
	createCmd.Flags().String(flag.SQL, "", "Filter the exported data with SQL SELECT statement.")
	createCmd.Flags().BoolVar(&force, flag.Force, false, "Create without confirmation. You need to confirm when you want to export the whole cluster in non-interactive mode.")
	createCmd.Flags().String(flag.CSVDelimiter, "", "Delimiter of string type variables in CSV files. (default \"\"\")")
	createCmd.Flags().String(flag.CSVSeparator, CSVSeparatorDefaultValue, "Separator of each value in CSV files.")
	createCmd.Flags().String(flag.CSVNullValue, "", "Representation of null values in CSV files. (default \"\\N\")")
	createCmd.Flags().Bool(flag.CSVSkipHeader, CSVSkipHeaderDefaultValue, "Export CSV files of the tables without header.")
	createCmd.Flags().String(flag.S3RoleArn, "", "The role arn of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	createCmd.Flags().String(flag.GCSURI, "", "The GCS URI in gs://<bucket>/<path> format. Required when target type is GCS.")
	createCmd.Flags().String(flag.GCSServiceAccountKey, "", "The base64 encoded service account key of GCS.")
	createCmd.Flags().String(flag.AzureBlobURI, "", "The Azure Blob URI in azure://<account>.blob.core.windows.net/<container>/<path> format. Required when target type is AZURE_BLOB.")
	createCmd.Flags().String(flag.AzureBlobSASToken, "", "The SAS token of Azure Blob.")
	createCmd.Flags().String(flag.OSSURI, "", "The OSS URI in oss://<bucket>/<path> format. Required when target type is OSS.")
	createCmd.Flags().String(flag.OSSAccessKeyID, "", "The access key ID of the OSS.")
	createCmd.Flags().String(flag.OSSAccessKeySecret, "", "The access key secret of the OSS.")
	createCmd.Flags().String(flag.ParquetCompression, "ZSTD", fmt.Sprintf("The parquet compression algorithm. One of %q.", export.AllowedExportParquetCompressionTypeEnumEnumValues))
	createCmd.Flags().String(flag.DisplayName, "", "The display name of the export. (default \"SNAPSHOT_<snapshot_time>\")")

	createCmd.MarkFlagsMutuallyExclusive(flag.TableFilter, flag.SQL)
	createCmd.MarkFlagsMutuallyExclusive(flag.TableWhere, flag.SQL)
	createCmd.MarkFlagsMutuallyExclusive(flag.S3RoleArn, flag.S3AccessKeyID)
	createCmd.MarkFlagsMutuallyExclusive(flag.S3RoleArn, flag.S3SecretAccessKey)
	return createCmd
}
