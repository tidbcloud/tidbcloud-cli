// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
)

type ConfigOpts struct {
	interactive bool
}

func (c ConfigOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.AuditLogUnRedacted,
		flag.Enabled,
		flag.CloudStorageType,
		flag.S3URI,
		flag.S3AccessKeyID,
		flag.S3SecretAccessKey,
		flag.S3RoleArn,
		flag.GCSURI,
		flag.GCSServiceAccountKey,
		flag.AzureBlobURI,
		flag.AzureBlobSASToken,
		flag.OSSURI,
		flag.OSSAccessKeyID,
		flag.OSSAccessKeySecret,
	}
}

type mutableField string

const (
	Unredacted              mutableField = "unredacted"
	Enabled                 mutableField = "enabled"
	CloudStorage            mutableField = "cloud-storage"
	RotationIntervalMinutes mutableField = "rotation-interval-minutes"
	RotationSizeMib         mutableField = "rotation-size-mib"
)

var mutableFields = []string{
	string(Unredacted),
	string(Enabled),
	string(CloudStorage),
	string(RotationIntervalMinutes),
	string(RotationSizeMib),
}

func (c *ConfigOpts) MarkInteractive(cmd *cobra.Command) error {
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
		err := cmd.MarkFlagRequired(flag.ClusterID)
		if err != nil {
			return err
		}
		cmd.MarkFlagsOneRequired(flag.AuditLogUnRedacted, flag.Enabled, flag.CloudStorageType, flag.RotationIntervalMinutes, flag.RotationSizeMib)
	}
	return nil
}

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := ConfigOpts{
		interactive: true,
	}

	var updateCmd = &cobra.Command{
		Use:         "update",
		Short:       "Update the database audit logging configuration",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Conigure database audit logging in interactive mode:
  $ %[1]s serverless audit-log config update

  Unredact the database audit log in non-interactive mode:
  $ %[1]s serverless audit-log config update -c <cluster-id> --unredacted

  Enable database audit logging with S3 Cloud Storage in non-interactive mode:
  $ %[1]s serverless audit-log config update -c <cluster-id> --enabled --cloud-storage S3 --s3.uri <s3-uri> --s3.access-key-id <s3-access-key-id> --s3.secret-access-key <s3-secret-access-key>

  Configure database audit logging rotation strategy in non-interactive mode:
  $ %[1]s serverless audit-log config update -c <cluster-id> --rotation-interval-minutes <rotation-interval-minutes> --rotation-size-mib <rotation-size-mib>

  Disable database audit logging in non-interactive mode:
  $ %[1]s serverless audit-log config update -c <cluster-id> --enabled=false`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID string
			var unredacted, enabled *bool
			var rotationIntervalMinutes, rotationSizeMib int32
			var cloudStorage auditlog.CloudStorageTypeEnum
			// s3
			var s3URI, accessKeyID, secretAccessKey, s3RoleArn string
			// gcs
			var gcsURI, gcsServiceAccountKey string
			// azure
			var azBlobURI, azBlobSasToken string
			// oss
			var ossURI, ossAccessKeyID, ossAccessKeySecret string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				selectedCluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = selectedCluster.ID

				fieldName, err := cloud.GetSelectedField(mutableFields, "Choose one field to config:")
				if err != nil {
					return err
				}

				switch fieldName {
				case string(Unredacted):
					prompt := &survey.Confirm{
						Message: "unredact the database audit log?",
						Default: false,
					}
					var confirm bool
					err = survey.AskOne(prompt, &confirm)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					if confirm {
						unredacted = aws.Bool(true)
					} else {
						unredacted = aws.Bool(false)
					}
				case string(Enabled):
					prompt := &survey.Confirm{
						Message: "enable database audit logging?",
						Default: false,
					}
					var confirm bool
					err = survey.AskOne(prompt, &confirm)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					if confirm {
						enabled = aws.Bool(true)
					} else {
						enabled = aws.Bool(false)
					}
				case string(CloudStorage):
					cloudStorage, err = GetSelectedCloudStorage()
					if err != nil {
						return err
					}
					selectedAuthType, err := GetSelectedAuthType(cloudStorage, *selectedCluster.CloudProvider)
					if err != nil {
						return err
					}
					switch selectedAuthType {
					// Both S3 and OSS supports ACCESS_KEY
					case string(auditlog.S3CLOUDSTORAGES3AUTHTYPE_ACCESS_KEY):
						if cloudStorage == auditlog.CLOUDSTORAGETYPEENUM_S3 {
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
						if cloudStorage == auditlog.CLOUDSTORAGETYPEENUM_OSS {
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
					case string(auditlog.S3CLOUDSTORAGES3AUTHTYPE_ROLE_ARN):
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
					case string(auditlog.GCSCLOUDSTORAGEGCSAUTHTYPE_SERVICE_ACCOUNT_KEY):
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
					case string(auditlog.AZUREBLOBCLOUDSTORAGEAZUREBLOBAUTHTYPE_SAS_TOKEN):
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
				case string(RotationIntervalMinutes):
					var rotationIntervalStr string
					prompt := &survey.Input{
						Message: "input rotation interval in minutes (range [10, 1440]):",
					}
					err = survey.AskOne(prompt, &rotationIntervalStr)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					rotationIntervalMinutesInt64, err := strconv.ParseInt(rotationIntervalStr, 10, 32)
					if err != nil {
						return errors.Trace(err)
					}
					rotationIntervalMinutes = int32(rotationIntervalMinutesInt64)
				case string(RotationSizeMib):
					var rotationSizeStr string
					prompt := &survey.Input{
						Message: "input rotation size in MiB (range [1, 1024]):",
					}
					err = survey.AskOne(prompt, &rotationSizeStr)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					rotationSizeInt64, err := strconv.ParseInt(rotationSizeStr, 10, 32)
					if err != nil {
						return errors.Trace(err)
					}
					rotationSizeMib = int32(rotationSizeInt64)
				}
			} else {
				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID = cID
				if cmd.Flags().Changed(flag.AuditLogUnRedacted) {
					u, err := cmd.Flags().GetBool(flag.AuditLogUnRedacted)
					if err != nil {
						return errors.Trace(err)
					}
					unredacted = &u
				}
				if cmd.Flags().Changed(flag.Enabled) {
					u, err := cmd.Flags().GetBool(flag.Enabled)
					if err != nil {
						return errors.Trace(err)
					}
					enabled = &u
				}
				rotationIntervalMinutes, err = cmd.Flags().GetInt32(flag.RotationIntervalMinutes)
				if err != nil {
					return errors.Trace(err)
				}
				rotationSizeMib, err = cmd.Flags().GetInt32(flag.RotationSizeMib)
				if err != nil {
					return errors.Trace(err)
				}
				cloudStorageStr, err := cmd.Flags().GetString(flag.CloudStorageType)
				if err != nil {
					return errors.Trace(err)
				}
				cloudStorage = auditlog.CloudStorageTypeEnum(strings.ToUpper(cloudStorageStr))
				if cloudStorage != "" && !cloudStorage.IsValid() {
					return errors.New("unsupported target type: " + cloudStorageStr)
				}
				switch cloudStorage {
				case auditlog.CLOUDSTORAGETYPEENUM_S3:
					s3URI, err = cmd.Flags().GetString(flag.S3URI)
					if err != nil {
						return errors.Trace(err)
					}
					if s3URI == "" {
						return errors.New("S3 URI is required when cloud storage is S3")
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
				case auditlog.CLOUDSTORAGETYPEENUM_GCS:
					gcsURI, err = cmd.Flags().GetString(flag.GCSURI)
					if err != nil {
						return errors.Trace(err)
					}
					if gcsURI == "" {
						return errors.New("GCS URI is required when cloud storage is GCS")
					}
					gcsServiceAccountKey, err = cmd.Flags().GetString(flag.GCSServiceAccountKey)
					if err != nil {
						return errors.Trace(err)
					}
					if gcsServiceAccountKey == "" {
						return errors.New("GCS service account key is required when cloud storage is GCS")
					}
				case auditlog.CLOUDSTORAGETYPEENUM_AZURE_BLOB:
					azBlobURI, err = cmd.Flags().GetString(flag.AzureBlobURI)
					if err != nil {
						return errors.Trace(err)
					}
					if azBlobURI == "" {
						return errors.New("Azure Blob URI is required when cloud storage is AZURE_BLOB")
					}
					azBlobSasToken, err = cmd.Flags().GetString(flag.AzureBlobSASToken)
					if err != nil {
						return errors.Trace(err)
					}
					if azBlobSasToken == "" {
						return errors.New("Azure Blob SAS token is required when cloud storage is AZURE_BLOB")
					}
				case auditlog.CLOUDSTORAGETYPEENUM_OSS:
					ossURI, err = cmd.Flags().GetString(flag.OSSURI)
					if err != nil {
						return errors.Trace(err)
					}
					if ossURI == "" {
						return errors.New("OSS URI is required when cloud storage is OSS")
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
						return errors.New("OSS access key id and access key secret are required when cloud storage is OSS")
					}
				}
			}

			body := &auditlog.DatabaseAuditLogServiceUpdateAuditLogConfigBody{
				AuditLogConfig: auditlog.AuditLogConfig{},
			}
			if unredacted != nil {
				body.AuditLogConfig.Unredacted = *auditlog.NewNullableBool(unredacted)
			}
			if enabled != nil {
				body.AuditLogConfig.Enabled = *auditlog.NewNullableBool(enabled)
			}
			if rotationIntervalMinutes > 0 {
				body.AuditLogConfig.RotationIntervalMinutes = &rotationIntervalMinutes
			}
			if rotationSizeMib > 0 {
				body.AuditLogConfig.RotationSizeMib = &rotationSizeMib
			}

			if cloudStorage != "" {
				body.AuditLogConfig.CloudStorage = &auditlog.CloudStorage{
					Type: cloudStorage,
				}
				switch cloudStorage {
				case auditlog.CLOUDSTORAGETYPEENUM_S3:
					if s3RoleArn != "" {
						body.AuditLogConfig.CloudStorage.S3 = &auditlog.S3CloudStorage{
							Uri:      s3URI,
							AuthType: auditlog.S3CLOUDSTORAGES3AUTHTYPE_ROLE_ARN,
							RoleArn:  aws.String(s3RoleArn),
						}
					} else {
						body.AuditLogConfig.CloudStorage.S3 = &auditlog.S3CloudStorage{
							Uri:      s3URI,
							AuthType: auditlog.S3CLOUDSTORAGES3AUTHTYPE_ACCESS_KEY,
							AccessKey: &auditlog.S3CloudStorageAccessKey{
								Id:     accessKeyID,
								Secret: secretAccessKey,
							},
						}
					}
				case auditlog.CLOUDSTORAGETYPEENUM_GCS:
					body.AuditLogConfig.CloudStorage.Gcs = &auditlog.GCSCloudStorage{
						Uri:               gcsURI,
						AuthType:          auditlog.GCSCLOUDSTORAGEGCSAUTHTYPE_SERVICE_ACCOUNT_KEY,
						ServiceAccountKey: aws.String(gcsServiceAccountKey),
					}
				case auditlog.CLOUDSTORAGETYPEENUM_AZURE_BLOB:
					body.AuditLogConfig.CloudStorage.AzureBlob = &auditlog.AzureBlobCloudStorage{
						Uri:      azBlobURI,
						AuthType: auditlog.AZUREBLOBCLOUDSTORAGEAZUREBLOBAUTHTYPE_SAS_TOKEN,
						SasToken: aws.String(azBlobSasToken),
					}
				case auditlog.CLOUDSTORAGETYPEENUM_OSS:
					body.AuditLogConfig.CloudStorage.Oss = &auditlog.OSSCloudStorage{
						Uri:      ossURI,
						AuthType: auditlog.OSSCLOUDSTORAGEOSSAUTHTYPE_ACCESS_KEY,
						AccessKey: &auditlog.OSSCloudStorageAccessKey{
							Id:     ossAccessKeyID,
							Secret: ossAccessKeySecret,
						},
					}
				}
			}
			_, err = d.UpdateAuditLogConfig(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("update cluster %s database audit logging configuration succeed", clusterID)))
			return nil
		},
	}

	updateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster to be updated.")
	updateCmd.Flags().Bool(flag.AuditLogUnRedacted, false, "unredact or redact the database audit log.")
	updateCmd.Flags().Bool(flag.Enabled, false, "enable or disable database audit logging.")
	updateCmd.Flags().String(flag.S3URI, "", "The S3 URI in s3://<bucket>/<path> format. Required when cloud storage is S3.")
	updateCmd.Flags().String(flag.S3RoleArn, "", "The role arn of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	updateCmd.Flags().String(flag.S3AccessKeyID, "", "The access key ID of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	updateCmd.Flags().String(flag.S3SecretAccessKey, "", "The secret access key of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].")
	updateCmd.Flags().String(flag.GCSURI, "", "The GCS URI in gs://<bucket>/<path> format. Required when cloud storage is GCS.")
	updateCmd.Flags().String(flag.GCSServiceAccountKey, "", "The base64 encoded service account key of GCS.")
	updateCmd.Flags().String(flag.AzureBlobURI, "", "The Azure Blob URI in azure://<account>.blob.core.windows.net/<container>/<path> format. Required when cloud storage is AZURE_BLOB.")
	updateCmd.Flags().String(flag.AzureBlobSASToken, "", "The SAS token of Azure Blob.")
	updateCmd.Flags().String(flag.OSSURI, "", "The OSS URI in oss://<bucket>/<path> format. Required when cloud storage is OSS.")
	updateCmd.Flags().String(flag.OSSAccessKeyID, "", "The access key ID of the OSS.")
	updateCmd.Flags().String(flag.OSSAccessKeySecret, "", "The access key secret of the OSS.")
	updateCmd.Flags().String(flag.CloudStorageType, "", fmt.Sprintf("The cloud storage. One of %q.", auditlog.AllowedCloudStorageTypeEnumEnumValues))
	updateCmd.Flags().Int32(flag.RotationIntervalMinutes, 0, "The rotation interval in minutes, range [10, 1440].")
	updateCmd.Flags().Int32(flag.RotationSizeMib, 0, "The rotation size in MiB, range [1, 1024].")

	updateCmd.MarkFlagsMutuallyExclusive(flag.S3RoleArn, flag.S3AccessKeyID)
	updateCmd.MarkFlagsMutuallyExclusive(flag.S3RoleArn, flag.S3SecretAccessKey)
	return updateCmd
}
