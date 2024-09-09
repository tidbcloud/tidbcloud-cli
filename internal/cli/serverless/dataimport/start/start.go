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
	"context"
	"fmt"
	"strconv"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	imp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/import"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/aws/aws-sdk-go-v2/aws"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

const (
	defaultCsvSeparator         = ","
	defaultCsvDelimiter         = "\""
	defaultCsvNullValue         = "\\N"
	defaultCsvNotNull           = false
	defaultCsvSkipHeader        = false
	defaultCsvBackslashEscape   = true
	defaultCsvTrimLastSeparator = false
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
	flag.CSVSeparator:         "Input the CSV separator: separator of each value in CSV files, skip to use default value (,)",
	flag.CSVDelimiter:         "Input the CSV delimiter: delimiter of string type variables in CSV files, skip to use default value (\"). If you want to set empty string, please use non-interactive mode",
	flag.CSVNullValue:         "Input the CSV null value: representation of null values in CSV files, skip to use default value (\\N). If you want to set empty string, please use non-interactive mode",
	flag.CSVSkipHeader:        "Input the CSV skip header: export CSV files of the tables without header. Type `true` to skip header, others will not skip header, default value (false)",
	flag.CSVBackslashEscape:   "Input the CSV backslash-escape: whether to interpret backslash escapes inside fields, skip to use default value (true)",
	flag.CSVTrimLastSeparator: "Input the CSV trim-last-separator: remove the last separator when a line ends with a separator, skip to use default value (false)",
	flag.CSVNotNull:           "Input the CSV not-null: whether the CSV can contains any NULL value, skip to use default value (false)",
}

type StartOpts struct {
	interactive bool
}

func (o StartOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.FileType,
		flag.SourceType,
		flag.LocalFilePath,
		flag.LocalTargetDatabase,
		flag.LocalTargetTable,
		flag.LocalConcurrency,
		flag.S3URI,
		flag.S3AccessKeyID,
		flag.S3SecretAccessKey,
		flag.S3RoleArn,
		flag.GCSURI,
		flag.GCSServiceAccountKey,
		flag.AzureBlobURI,
		flag.AzureBlobSASToken,
		flag.CSVSeparator,
		flag.CSVDelimiter,
		flag.CSVNotNull,
		flag.CSVNullValue,
		flag.CSVBackslashEscape,
		flag.CSVTrimLastSeparator,
		flag.CSVSkipHeader,
	}
}

func (o StartOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func StartCmd(h *internal.Helper) *cobra.Command {
	var concurrency int
	opts := StartOpts{
		interactive: true,
	}
	startCmd := &cobra.Command{
		Use:         "start",
		Short:       "Start a data import task",
		Aliases:     []string{"create"},
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Start an import task in interactive mode:
  $ %[1]s serverless import start

  Start a local import task in non-interactive mode:
  $ %[1]s serverless import start --local.file-path <file-path> --cluster-id <cluster-id> --file-type <file-type> --local.target-database <target-database> --local.target-table <target-table>

  Start a local import task with custom upload concurrency:
  $ %[1]s serverless import start --local.file-path <file-path> --cluster-id <cluster-id> --file-type <file-type> --local.target-database <target-database> --local.target-table <target-table> --local.concurrency 10
	
  Start a local import task with custom CSV format:
  $ %[1]s serverless import start --local.file-path <file-path> --cluster-id <cluster-id> --file-type CSV --local.target-database <target-database> --local.target-table <target-table> --csv.separator \" --csv.delimiter \' --csv.backslash-escape=false --csv.trim-last-separator=true

  Start an S3 import task in non-interactive mode:
  $ %[1]s serverless import start --source-type S3 --s3.uri <s3-uri> --cluster-id <cluster-id> --file-type <file-type> --s3.role-arn <role-arn>

  Start a GCS import task in non-interactive mode:
  $ %[1]s serverless import start --source-type GCS --gcs.uri <gcs-uri> --cluster-id <cluster-id> --file-type <file-type> --gcs.service-account-key <service-account-key>

  Start an Azure Blob import task in non-interactive mode:
  $ %[1]s serverless import start --source-type AZURE_BLOB --azblob.uri <azure-blob-uri> --cluster-id <cluster-id> --file-type <file-type> --azblob.sas-token <sas-token>
`,
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
				for _, fn := range opts.RequiredFlags() {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var sourceType imp.ImportSourceTypeEnum
			var clusterId string
			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}
				var err error
				ctx := cmd.Context()
				d, err := h.Client()
				if err != nil {
					return err
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

				sourceType, err = getSelectedSourceType()
				if err != nil {
					return err
				}
			} else {
				sourceType = imp.ImportSourceTypeEnum(cmd.Flag(flag.SourceType).Value.String())
				var err error
				clusterId, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			cmd.Annotations[telemetry.ClusterID] = clusterId

			if sourceType == imp.IMPORTSOURCETYPEENUM_LOCAL {
				localOpts := LocalOpts{
					concurrency: concurrency,
					h:           h,
					interactive: opts.interactive,
					clusterId:   clusterId,
				}
				return localOpts.Run(cmd)
			} else if sourceType == imp.IMPORTSOURCETYPEENUM_S3 {
				s3Opts := S3Opts{
					h:           h,
					interactive: opts.interactive,
					clusterId:   clusterId,
				}
				return s3Opts.Run(cmd)
			} else if sourceType == imp.IMPORTSOURCETYPEENUM_GCS {
				gcsOpts := GCSOpts{
					h:           h,
					interactive: opts.interactive,
					clusterId:   clusterId,
				}
				return gcsOpts.Run(cmd)
			} else if sourceType == imp.IMPORTSOURCETYPEENUM_AZURE_BLOB {
				azBlobOpts := AzBlobOpts{
					h:           h,
					interactive: opts.interactive,
					clusterId:   clusterId,
				}
				return azBlobOpts.Run(cmd)
			} else {
				return errors.New("unsupported import source type")
			}
		},
	}

	startCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	startCmd.Flags().String(flag.SourceType, "LOCAL", fmt.Sprintf("The import source type, one of %q.", imp.AllowedImportSourceTypeEnumEnumValues))
	startCmd.Flags().String(flag.FileType, "", fmt.Sprintf("The import file type, one of %q.", imp.AllowedImportFileTypeEnumEnumValues))

	startCmd.Flags().String(flag.LocalFilePath, "", "The local file path to import.")
	startCmd.Flags().String(flag.LocalTargetDatabase, "", "Target database to which import data.")
	startCmd.Flags().String(flag.LocalTargetTable, "", "Target table to which import data.")
	startCmd.Flags().IntVar(&concurrency, flag.LocalConcurrency, 5, "The concurrency for uploading file.")

	startCmd.Flags().String(flag.S3AccessKeyID, "", "The access key ID of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	startCmd.Flags().String(flag.S3SecretAccessKey, "", "The secret access key of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	startCmd.Flags().String(flag.S3RoleArn, "", "The role arn of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	startCmd.Flags().String(flag.S3URI, "", "The S3 URI in s3://<bucket>/<path> format. Required when source type is S3.")
	startCmd.MarkFlagsMutuallyExclusive(flag.S3RoleArn, flag.S3AccessKeyID)
	startCmd.MarkFlagsMutuallyExclusive(flag.S3RoleArn, flag.S3SecretAccessKey)
	startCmd.MarkFlagsRequiredTogether(flag.S3AccessKeyID, flag.S3SecretAccessKey)

	startCmd.Flags().String(flag.GCSURI, "", "The GCS URI in gcs://<bucket>/<path> format. Required when source type is GCS.")
	startCmd.Flags().String(flag.GCSServiceAccountKey, "", "The base64 encoded service account key of GCS.")

	startCmd.Flags().String(flag.AzureBlobURI, "", "The Azure Blob URI in azure://<account>.blob.core.windows.net/<container>/<path> format.")
	startCmd.Flags().String(flag.AzureBlobSASToken, "", "The SAS token of Azure Blob.")

	startCmd.Flags().String(flag.CSVDelimiter, defaultCsvDelimiter, "The delimiter used for quoting of CSV file.")
	startCmd.Flags().String(flag.CSVSeparator, defaultCsvSeparator, "The field separator of CSV file.")
	startCmd.Flags().Bool(flag.CSVTrimLastSeparator, defaultCsvTrimLastSeparator, "Specifies whether to treat separator as the line terminator and trim all trailing separators in the CSV file.")
	startCmd.Flags().Bool(flag.CSVBackslashEscape, defaultCsvBackslashEscape, "Specifies whether to interpret backslash escapes inside fields in the CSV file.")
	startCmd.Flags().Bool(flag.CSVNotNull, defaultCsvNotNull, "Specifies whether a CSV file can contain any NULL values.")
	startCmd.Flags().String(flag.CSVNullValue, defaultCsvNullValue, "The representation of NULL values in the CSV file.")
	startCmd.Flags().Bool(flag.CSVSkipHeader, defaultCsvSkipHeader, "Specifies whether the CSV file contains a header line.")
	return startCmd
}

func getSelectedSourceType() (imp.ImportSourceTypeEnum, error) {
	SourceTypes := make([]interface{}, 0, len(imp.AllowedImportSourceTypeEnumEnumValues))
	for _, sourceType := range imp.AllowedImportSourceTypeEnumEnumValues {
		SourceTypes = append(SourceTypes, sourceType)
	}
	model, err := ui.InitialSelectModel(SourceTypes, "Choose import source type:")
	if err != nil {
		return "", errors.Trace(err)
	}

	p := tea.NewProgram(model)
	SourceTypeModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := SourceTypeModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	sourceType := SourceTypeModel.(ui.SelectModel).GetSelectedItem()
	if sourceType == nil {
		return "", errors.New("no source type selected")
	}
	return sourceType.(imp.ImportSourceTypeEnum), nil
}

func waitStartOp(ctx context.Context, h *internal.Helper, d cloud.TiDBCloudClient, clusterId string, body *imp.ImportServiceCreateImportBody) error {
	fmt.Fprintf(h.IOStreams.Out, "... Starting the import task\n")
	res, err := d.CreateImport(ctx, clusterId, body)
	if err != nil {
		return err
	}

	fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", *res.ImportId))
	return nil
}

func spinnerWaitStartOp(ctx context.Context, h *internal.Helper, d cloud.TiDBCloudClient, clusterId string, body *imp.ImportServiceCreateImportBody) error {
	task := func() tea.Msg {
		errChan := make(chan error, 1)

		go func() {
			res, err := d.CreateImport(ctx, clusterId, body)
			if err != nil {
				errChan <- err
				return
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", *res.ImportId))
			errChan <- nil
		}()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		timer := time.After(2 * time.Minute)
		for {
			select {
			case <-timer:
				return fmt.Errorf("timeout waiting for import task to start")
			case <-ticker.C:
				// continue
			case err := <-errChan:
				if err != nil {
					return err
				} else {
					return ui.Result("")
				}
			case <-ctx.Done():
				return util.InterruptError
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Starting import task"))
	model, err := p.Run()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := model.(ui.SpinnerModel); m.Interrupted {
		return util.InterruptError
	}
	if m, _ := model.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	}

	return nil
}

func getCSVFormat() (format *imp.CSVFormat, errToReturn error) {
	separator, delimiter, nullValue := defaultCsvSeparator, defaultCsvDelimiter, defaultCsvNullValue
	backslashEscape, trimLastSeparator, skipHeader, notNull := defaultCsvBackslashEscape, defaultCsvTrimLastSeparator, defaultCsvSkipHeader, defaultCsvNotNull

	needCustomCSV := false
	prompt := &survey.Confirm{
		Message: "Do you need to customize CSV format?",
	}
	err := survey.AskOne(prompt, &needCustomCSV)
	if err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			errToReturn = util.InterruptError
			return
		} else {
			errToReturn = err
			return
		}
	}

	if needCustomCSV {
		// variables for input
		inputs := []string{flag.CSVSeparator, flag.CSVDelimiter, flag.CSVBackslashEscape, flag.CSVTrimLastSeparator,
			flag.CSVNotNull, flag.CSVNullValue, flag.CSVSkipHeader}
		inputModel, err := ui.InitialInputModel(inputs, inputDescription)
		if err != nil {
			return nil, err
		}

		// If user input is blank, use the default value.
		v := inputModel.Inputs[0].Value()
		if len(v) > 0 {
			separator = v
		}
		v = inputModel.Inputs[1].Value()
		if len(v) > 0 {
			delimiter = v
		}
		v = inputModel.Inputs[2].Value()
		if len(v) > 0 {
			backslashEscape, err = strconv.ParseBool(v)
			if err != nil {
				errToReturn = errors.Annotate(err, "backslashEscape must be true or false")
				return
			}
		}
		v = inputModel.Inputs[3].Value()
		if len(v) > 0 {
			trimLastSeparator, err = strconv.ParseBool(v)
			if err != nil {
				errToReturn = errors.Annotate(err, "trimLastSeparator must be true or false")
				return
			}
		}
		v = inputModel.Inputs[4].Value()
		if len(v) > 0 {
			notNull, err = strconv.ParseBool(v)
			if err != nil {
				errToReturn = errors.Annotate(err, "notNull must be true or false")
				return
			}
		}
		v = inputModel.Inputs[5].Value()
		if len(v) > 0 {
			nullValue = v
		}
		v = inputModel.Inputs[6].Value()
		if len(v) > 0 {
			skipHeader, err = strconv.ParseBool(v)
			if err != nil {
				errToReturn = errors.Annotate(err, "skipHeader must be true or false")
				return
			}
		}
	}

	format = &imp.CSVFormat{
		Separator:         aws.String(separator),
		Delimiter:         *imp.NewNullableString(aws.String(delimiter)),
		BackslashEscape:   *imp.NewNullableBool(aws.Bool(backslashEscape)),
		TrimLastSeparator: *imp.NewNullableBool(aws.Bool(trimLastSeparator)),
		Null:              *imp.NewNullableString(aws.String(nullValue)),
		Header:            *imp.NewNullableBool(aws.Bool(!skipHeader)),
		NotNull:           *imp.NewNullableBool(aws.Bool(notNull)),
	}
	return
}

func getCSVFlagValue(cmd *cobra.Command) (*imp.CSVFormat, error) {
	// optional flags
	backslashEscape, err := cmd.Flags().GetBool(flag.CSVBackslashEscape)
	if err != nil {
		return nil, errors.Trace(err)
	}
	separator, err := cmd.Flags().GetString(flag.CSVSeparator)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(separator) == 0 {
		return nil, fmt.Errorf("CSV separator must not be empty")
	}
	delimiter, err := cmd.Flags().GetString(flag.CSVDelimiter)
	if err != nil {
		return nil, errors.Trace(err)
	}
	trimLastSeparator, err := cmd.Flags().GetBool(flag.CSVTrimLastSeparator)
	if err != nil {
		return nil, errors.Trace(err)
	}
	nullValue, err := cmd.Flags().GetString(flag.CSVNullValue)
	if err != nil {
		return nil, errors.Trace(err)
	}
	skipHeader, err := cmd.Flags().GetBool(flag.CSVSkipHeader)
	if err != nil {
		return nil, errors.Trace(err)
	}
	notNull, err := cmd.Flags().GetBool(flag.CSVNotNull)
	if err != nil {
		return nil, errors.Trace(err)
	}

	format := &imp.CSVFormat{
		Separator:         aws.String(separator),
		Delimiter:         *imp.NewNullableString(aws.String(delimiter)),
		BackslashEscape:   *imp.NewNullableBool(aws.Bool(backslashEscape)),
		TrimLastSeparator: *imp.NewNullableBool(aws.Bool(trimLastSeparator)),
		Null:              *imp.NewNullableString(aws.String(nullValue)),
		Header:            *imp.NewNullableBool(aws.Bool(!skipHeader)),
		NotNull:           *imp.NewNullableBool(aws.Bool(notNull)),
	}
	return format, nil
}
