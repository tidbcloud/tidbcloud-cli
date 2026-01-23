// Copyright 2026 PingCAP, Inc.
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
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
)

var inputDescription = map[string]string{
	flag.S3URI:                "Input your S3 URI in s3://<bucket>/<path> format",
	flag.S3AccessKeyID:        "Input your S3 access key id",
	flag.S3SecretAccessKey:    "Input your S3 secret access key",
	flag.S3RoleArn:            "Input your S3 role arn",
	flag.AzureBlobURI:         "Input your Azure Blob URI in azure://<account>.blob.core.windows.net/<container>/<path> format",
	flag.AzureBlobSASToken:    "Input your Azure Blob SAS token",
	flag.GCSURI:               "Input your GCS URI in gs://<bucket>/<path> format",
	flag.GCSServiceAccountKey: "Input your base64 encoded GCS service account key",
	flag.OSSURI:               "Input your OSS URI in oss://<bucket>/<path> format",
	flag.OSSAccessKeyID:       "Input your OSS access key id",
	flag.OSSAccessKeySecret:   "Input your OSS access key secret",
}

func GetSelectedCloudStorage() (auditlog.CloudStorageTypeEnum, error) {
	selectTypes := make([]interface{}, 0, len(auditlog.AllowedCloudStorageTypeEnumEnumValues))
	for _, v := range auditlog.AllowedCloudStorageTypeEnumEnumValues {
		selectTypes = append(selectTypes, v)
	}
	selectModel, err := ui.InitialSelectModel(selectTypes, "Choose the cloud storage:")
	if err != nil {
		return "", errors.Trace(err)
	}

	p := tea.NewProgram(selectModel)
	model, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := model.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	selectType := model.(ui.SelectModel).GetSelectedItem()
	if selectType == nil {
		return "", errors.New("no cloud storage selected")
	}
	return selectType.(auditlog.CloudStorageTypeEnum), nil
}

func GetSelectedAuthType(target auditlog.CloudStorageTypeEnum, provider cluster.V1beta1RegionCloudProvider) (_ string, err error) {
	var model *ui.SelectModel
	switch target {
	case auditlog.CLOUDSTORAGETYPEENUM_S3:
		if provider != cluster.V1BETA1REGIONCLOUDPROVIDER_AWS {
			return string(auditlog.S3CLOUDSTORAGES3AUTHTYPE_ACCESS_KEY), nil
		}
		authTypes := make([]interface{}, 0, len(auditlog.AllowedS3CloudStorageS3AuthTypeEnumValues))
		for _, v := range auditlog.AllowedS3CloudStorageS3AuthTypeEnumValues {
			authTypes = append(authTypes, string(v))
		}
		model, err = ui.InitialSelectModel(authTypes, "Choose and input the S3 auth:")
		if err != nil {
			return "", errors.Trace(err)
		}
	case auditlog.CLOUDSTORAGETYPEENUM_GCS:
		return string(auditlog.GCSCLOUDSTORAGEGCSAUTHTYPE_SERVICE_ACCOUNT_KEY), nil
	case auditlog.CLOUDSTORAGETYPEENUM_AZURE_BLOB:
		return string(auditlog.AZUREBLOBCLOUDSTORAGEAZUREBLOBAUTHTYPE_SAS_TOKEN), nil
	case auditlog.CLOUDSTORAGETYPEENUM_TIDB_CLOUD:
		return "", nil
	case auditlog.CLOUDSTORAGETYPEENUM_OSS:
		return string(auditlog.OSSCLOUDSTORAGEOSSAUTHTYPE_ACCESS_KEY), nil
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
