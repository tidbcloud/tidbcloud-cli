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
	"io"
	"os"
	"strings"
	"sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"
)

const DefaultConcurrency = 3

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

			resp, err := d.DownloadExport(ctx, clusterID, exportID)
			if err != nil {
				return errors.Trace(err)
			}

			var totalSize int64
			for _, download := range resp.Downloads {
				totalSize += *download.Size
			}
			fileMessage := fmt.Sprintf("There are %d files to download, total size is %s.", len(resp.Downloads), humanize.IBytes(uint64(totalSize)))

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
				err = DownloadFilesPrompt(h, resp.Downloads, path, concurrency)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				err = DownloadFilesWithoutPrompt(h, resp.Downloads, path, concurrency)
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

func DownloadFilesPrompt(h *internal.Helper, urls []export.DownloadUrl, path string, concurrency int) error {
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}

	// create the path if not exist
	err := util.CreateFolder(path)
	if err != nil {
		return err
	}

	// init the concurrency progress model
	var p *tea.Program
	urlMsgs := make([]ui.URLMsg, 0)
	for _, u := range urls {
		url := ui.URLMsg{
			Name: *u.Name,
			Url:  *u.Url,
			Size: *u.Size,
		}
		urlMsgs = append(urlMsgs, url)
	}
	m := ui.NewProcessDownloadModel(
		urlMsgs,
		concurrency,
		path,
	)

	// run the program
	p = tea.NewProgram(m)
	m.SetProgram(p)
	model, err := p.Run()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := model.(*ui.ProcessDownloadModel); m.Interrupted {
		return util.InterruptError
	}

	succeededCount := 0
	failedCount := 0
	skippedCount := 0
	for _, f := range m.GetFinishedJobs() {
		switch f.GetStatus() {
		case ui.Succeeded:
			succeededCount++
		case ui.Failed:
			failedCount++
		case ui.Skipped:
			skippedCount++
		}
	}
	fmt.Fprint(h.IOStreams.Out, generateDownloadSummary(succeededCount, skippedCount, failedCount))
	index := 0
	for _, f := range m.GetFinishedJobs() {
		if f.GetStatus() != ui.Succeeded {
			index++
			fmt.Fprintf(h.IOStreams.Out, "%d.%s\n", index, f.GetErrorString())
		}
	}

	if failedCount > 0 {
		return errors.New(fmt.Sprintf("%d file(s) failed to download", failedCount))
	}
	return nil
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

var wg sync.WaitGroup

type downloadJob struct {
	url  export.DownloadUrl
	path string
}

type downloadResult struct {
	name   string
	err    error
	status ui.JobStatus
}

func (r *downloadResult) GetErrorString() string {
	if r.status == ui.Succeeded {
		return ""
	}
	if r.err == nil {
		return fmt.Sprintf("%s %s", r.name, r.status)
	}
	return fmt.Sprintf("%s %s: %s", r.name, r.status, r.err.Error())
}

func DownloadFilesWithoutPrompt(h *internal.Helper, urls []export.DownloadUrl, path string, concurrency int) error {
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}
	// create the path if not exist
	err := util.CreateFolder(path)
	if err != nil {
		return err
	}

	jobs := make(chan *downloadJob, len(urls))
	results := make(chan *downloadResult, len(urls))
	// Start consumers:
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go consume(h, jobs, results)
	}
	// Start producing
	for _, u := range urls {
		jobs <- &downloadJob{url: u, path: path}
	}
	close(jobs)
	wg.Wait()
	close(results)

	succeededCount := 0
	failedCount := 0
	skippedCount := 0
	downloadResults := make([]*downloadResult, 0)
	for result := range results {
		switch result.status {
		case ui.Succeeded:
			succeededCount++
		case ui.Failed:
			failedCount++
		case ui.Skipped:
			skippedCount++
		}
		downloadResults = append(downloadResults, result)
	}
	fmt.Fprint(h.IOStreams.Out, generateDownloadSummary(succeededCount, skippedCount, failedCount))
	index := 0
	for _, f := range downloadResults {
		if f.status != ui.Succeeded {
			index++
			fmt.Fprintf(h.IOStreams.Out, "%d.%s\n", index, f.GetErrorString())
		}
	}
	if failedCount > 0 {
		return errors.New(fmt.Sprintf("%d file(s) failed to download", failedCount))
	}
	return nil
}

func consume(h *internal.Helper, jobs <-chan *downloadJob, results chan *downloadResult) {
	defer wg.Done()
	for job := range jobs {
		func() {
			var err error
			defer func() {
				if err != nil {
					if strings.Contains(err.Error(), "file already exists") {
						fmt.Fprintf(h.IOStreams.Out, "download %s skipped: %s\n", *job.url.Name, err.Error())
						results <- &downloadResult{name: *job.url.Name, err: err, status: ui.Skipped}
					} else {
						fmt.Fprintf(h.IOStreams.Out, "download %s failed: %s\n", *job.url.Name, err.Error())
						results <- &downloadResult{name: *job.url.Name, err: err, status: ui.Failed}
					}
				} else {
					fmt.Fprintf(h.IOStreams.Out, "download %s succeeded\n", *job.url.Name)
					results <- &downloadResult{name: *job.url.Name, err: nil, status: ui.Succeeded}
				}
			}()

			fmt.Fprintf(h.IOStreams.Out, "downloading %s | %s\n", *job.url.Name, humanize.IBytes(uint64(*job.url.Size)))

			// request the url
			resp, err := util.GetResponse(*job.url.Url, os.Getenv(config.DebugEnv) != "")
			if err != nil {
				return
			}
			defer resp.Body.Close()

			file, err := util.CreateFile(job.path, *job.url.Name)
			if err != nil {
				return
			}
			defer file.Close()
			_, err = io.Copy(file, resp.Body)
		}()
	}
}

func generateDownloadSummary(succeededCount, skippedCount, failedCount int) string {
	summaryMessage := fmt.Sprintf("%s %s %s", color.BlueString("download summary:"), color.GreenString("succeeded: %d", succeededCount), color.GreenString("skipped: %d", skippedCount))
	if failedCount > 0 {
		summaryMessage += color.RedString(" failed: %d", failedCount)
	} else {
		summaryMessage += fmt.Sprintf(" failed: %d", failedCount)
	}
	return summaryMessage + "\n"
}
