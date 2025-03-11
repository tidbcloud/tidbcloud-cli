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
package auditlog

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/export"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
)

func DownloadFilesPrompt(h *internal.Helper, path string,
	concurrency int, clusterID string, totalSize int64, fileNames []string, client cloud.TiDBCloudClient) error {
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}

	// create the path if not exist
	err := util.CreateFolder(path)
	if err != nil {
		return err
	}

	generateFunc := func(ctx context.Context, fileNames []string) (map[string]*string, error) {
		resp, err := client.DownloadAuditLogs(ctx, clusterID, &auditlog.AuditLogServiceDownloadAuditLogsBody{
			AuditLogNames: fileNames,
		})
		if err != nil {
			return nil, err
		}
		fileMap := make(map[string]*string)
		for _, file := range resp.AuditLogs {
			fileMap[*file.Name] = file.Url
		}
		return fileMap, nil
	}

	// init the concurrency progress model
	var p *tea.Program
	m := export.NewProcessDownloadModel(
		concurrency,
		path,
		int(totalSize),
		fileNames,
		generateFunc,
	)

	// run the program
	p = tea.NewProgram(m)
	m.SetProgram(p)
	model, err := p.Run()
	if err != nil {
		return err
	}
	if m, _ := model.(*export.ProcessDownloadModel); m.Interrupted {
		return util.InterruptError
	}

	succeededCount := 0
	failedCount := 0
	skippedCount := 0
	for _, f := range m.GetFinishedJobs() {
		switch f.GetStatus() {
		case export.Succeeded:
			succeededCount++
		case export.Failed:
			failedCount++
		case export.Skipped:
			skippedCount++
		}
	}
	fmt.Fprint(h.IOStreams.Out, export.GenerateDownloadSummary(succeededCount, skippedCount, failedCount))
	index := 0
	for _, f := range m.GetFinishedJobs() {
		if f.GetStatus() != export.Succeeded {
			index++
			fmt.Fprintf(h.IOStreams.Out, "%d.%s\n", index, f.GetResult())
		}
	}

	if failedCount > 0 {
		return errors.New(fmt.Sprintf("%d file(s) failed to download", failedCount))
	}
	return nil
}
