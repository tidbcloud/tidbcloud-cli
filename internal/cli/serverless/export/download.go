// Copyright 2026 PingCAP, Inc.
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
	"context"
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

const (
	DefaultConcurrency = 3
	MaxBatchSize       = 100
)

var DownloadPathInputFields = map[string]int{
	flag.OutputPath: 0,
}

type DownloadOpts struct {
	interactive bool
}

func (c DownloadOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ExportID,
		flag.OutputPath,
	}
}

func (c DownloadOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ExportID,
	}
}

func (c *DownloadOpts) MarkInteractive(cmd *cobra.Command) error {
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

func DownloadCmd(h *internal.Helper) *cobra.Command {
	var force bool
	opts := DownloadOpts{
		interactive: true,
	}

	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download the exported data",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Download the local type export in interactive mode:
  $ %[1]s serverless export download

  Download the local type export in non-interactive mode:
  $ %[1]s serverless export download -c <cluster-id> -e <export-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var exportID, clusterID, path string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				export, err := cloud.GetSelectedLocalExport(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				exportID = export.ID

				downloadPathInputModel, err := GetDownloadPathInput()
				if err != nil {
					return err
				}
				path = downloadPathInputModel.(ui.TextInputModel).
					Inputs[DownloadPathInputFields[flag.OutputPath]].Value()
			} else {
				// non-interactive mode, get values from flags
				exportID, err = cmd.Flags().GetString(flag.ExportID)
				if err != nil {
					return errors.Trace(err)
				}

				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				path, err = cmd.Flags().GetString(flag.OutputPath)
				if err != nil {
					return errors.Trace(err)
				}
			}

			concurrency, err := cmd.Flags().GetInt(flag.Concurrency)
			if err != nil {
				return errors.Trace(err)
			}

			exportFiles, err := cloud.GetAllExportFiles(ctx, clusterID, exportID, d)
			if err != nil {
				return errors.Trace(err)
			}

			var totalSize int64
			fileNames := make([]string, 0)
			for _, file := range exportFiles {
				totalSize += *file.Size
				if *file.Name == "metadata" {
					continue
				}
				fileNames = append(fileNames, *file.Name)
			}
			fileMessage := fmt.Sprintf("There are %d files to download, total size is %s.", len(fileNames), humanize.IBytes(uint64(totalSize)))
			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to download")
				}

				confirmationMessage := fmt.Sprintf("%s %s %s %s", color.BlueString(fileMessage), color.BlueString("Please type"), color.HiBlueString(confirmed), color.BlueString("to download:"))
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
					return errors.New("incorrect confirm string entered, skipping download")
				}
			} else {
				fmt.Fprintf(h.IOStreams.Out, "%s\n", color.BlueString(fileMessage))
			}
			if h.IOStreams.CanPrompt {
				err = DownloadFilesPrompt(h, path, concurrency, exportID, clusterID, totalSize, fileNames, d)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				err = DownloadFilesWithoutPrompt(h, path, concurrency, exportID, clusterID, fileNames, d)
				if err != nil {
					return errors.Trace(err)
				}
			}
			return nil
		},
	}

	downloadCmd.Flags().StringP(flag.ExportID, flag.ExportIDShort, "", "The ID of the export.")
	downloadCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID of the export.")
	downloadCmd.Flags().String(flag.OutputPath, "", "Where you want to download to. If not specified, download to the current directory.")
	downloadCmd.Flags().BoolVar(&force, flag.Force, false, "Download without confirmation.")
	downloadCmd.Flags().Int(flag.Concurrency, 3, "Download concurrency.")
	downloadCmd.MarkFlagsRequiredTogether(flag.ExportID, flag.ClusterID)
	return downloadCmd
}

func DownloadFilesPrompt(h *internal.Helper, path string,
	concurrency int, exportID, clusterID string, totalSize int64, fileNames []string, client cloud.TiDBCloudClient) error {
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}

	// create the path if not exist
	err := util.CreateFolder(path)
	if err != nil {
		return err
	}

	generateFunc := func(ctx context.Context, fileNames []string) (map[string]*string, error) {
		resp, err := client.DownloadExportFiles(ctx, clusterID, exportID, &export.ExportServiceDownloadExportFilesBody{
			FileNames: fileNames,
		})
		if err != nil {
			return nil, err
		}
		fileMap := make(map[string]*string)
		for _, file := range resp.Files {
			fileMap[*file.Name] = file.Url
		}
		return fileMap, nil
	}

	// init the concurrency progress model
	var p *tea.Program
	m := NewProcessDownloadModel(
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
		return errors.Trace(err)
	}
	if m, _ := model.(*ProcessDownloadModel); m.Interrupted {
		return util.InterruptError
	}

	succeededCount := 0
	failedCount := 0
	skippedCount := 0
	for _, f := range m.GetFinishedJobs() {
		switch f.GetStatus() {
		case Succeeded:
			succeededCount++
		case Failed:
			failedCount++
		case Skipped:
			skippedCount++
		}
	}
	fmt.Fprint(h.IOStreams.Out, GenerateDownloadSummary(succeededCount, skippedCount, failedCount))
	index := 0
	for _, f := range m.GetFinishedJobs() {
		if f.GetStatus() != Succeeded {
			index++
			fmt.Fprintf(h.IOStreams.Out, "%d.%s\n", index, f.GetResult())
		}
	}

	if failedCount > 0 {
		return errors.New(fmt.Sprintf("%d file(s) failed to download", failedCount))
	}
	return nil
}

func DownloadFilesWithoutPrompt(h *internal.Helper, path string,
	concurrency int, exportID, clusterID string, fileNames []string, client cloud.TiDBCloudClient) error {

	generateFunc := func(ctx context.Context, fileNames []string) (map[string]*string, error) {
		resp, err := client.DownloadExportFiles(ctx, clusterID, exportID, &export.ExportServiceDownloadExportFilesBody{
			FileNames: fileNames,
		})
		if err != nil {
			return nil, err
		}
		fileMap := make(map[string]*string)
		for _, file := range resp.Files {
			fileMap[*file.Name] = file.Url
		}
		return fileMap, nil
	}

	exportDownloadPool, err := NewDownloadPool(h, path, concurrency, fileNames, generateFunc)
	if err != nil {
		return errors.Trace(err)
	}
	err = exportDownloadPool.Start()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func getBatchSize(concurrency int) int {
	batchSize := 2 * concurrency
	if batchSize > MaxBatchSize {
		batchSize = MaxBatchSize
	}
	return batchSize
}
