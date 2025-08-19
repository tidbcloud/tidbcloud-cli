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

package auditlog

import (
	"context"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/export"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
)

func DownloadFilesWithoutPrompt(h *internal.Helper, path string,
	concurrency int, clusterID string, fileNames []string, client cloud.TiDBCloudClient) error {

	generateFunc := func(ctx context.Context, fileNames []string) (map[string]*string, error) {
		resp, err := client.DownloadAuditLogs(ctx, clusterID, &auditlog.DatabaseAuditLogServiceDownloadAuditLogFilesBody{
			AuditLogNames: fileNames,
		})
		if err != nil {
			return nil, err
		}
		fileMap := make(map[string]*string)
		for _, file := range resp.AuditLogFiles {
			fileMap[*file.Name] = file.Url
		}
		return fileMap, nil
	}

	exportDownloadPool, err := export.NewDownloadPool(h, path, concurrency, fileNames, generateFunc)
	if err != nil {
		return err
	}
	err = exportDownloadPool.Start()
	if err != nil {
		return err
	}
	return nil
}
