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
	"net/http"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/charmbracelet/bubbles/progress"
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
	exportApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/client/export_service"
	exportModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/models"
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
		Short: "Download the local type export",
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

			var exportID, clusterID, path string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				export, err := cloud.GetSelectedLocalExport(clusterID, h.QueryPageSize, d)
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

			params := exportApi.NewExportServiceDownloadExportParams().
				WithClusterID(clusterID).WithExportID(exportID)
			resp, err := d.DownloadExport(params)
			if err != nil {
				return errors.Trace(err)
			}

			var totalSize int64
			for _, download := range resp.Payload.Downloads {
				totalSize += download.Size
			}
			fileMessage := fmt.Sprintf("There are %d files to download, total size is %s.", len(resp.Payload.Downloads), humanize.IBytes(uint64(totalSize)))

			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to download")
				}

				confirmationMessage := fmt.Sprintf("\n%s %s %s %s", color.BlueString(fileMessage), color.BlueString("Please type"), color.HiBlueString(confirmed), color.BlueString("to download:"))
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
				fmt.Fprintf(h.IOStreams.Out, "\n%s\n", color.BlueString(fileMessage))
			}

			err = DownloadFiles(h, resp.Payload.Downloads, path)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
	}

	downloadCmd.Flags().StringP(flag.ExportID, flag.ExportIDShort, "", "The ID of the export.")
	downloadCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID of the export.")
	downloadCmd.Flags().String(flag.OutputPath, "", "Where you want to download to. If not specified, download to the current directory.")
	downloadCmd.Flags().BoolVar(&force, flag.Force, false, "Download without confirmation.")
	downloadCmd.MarkFlagsRequiredTogether(flag.ExportID, flag.ClusterID)
	return downloadCmd
}

func DownloadFiles(h *internal.Helper, urls []*exportModel.V1beta1DownloadURL, path string) error {
	if path == "" {
		path = "."
	}

	// create the path if not exist
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0755)
		}
		if err != nil {
			return err
		}
	}
	interrupt := false
	for _, downloadUrl := range urls {
		if interrupt {
			return util.InterruptError
		}
		func() {
			fileName := downloadUrl.Name
			url := downloadUrl.URL
			size := humanize.IBytes(uint64(downloadUrl.Size))
			fmt.Fprintf(h.IOStreams.Out, "\ndownload %s(%s) to %s\n", fileName, size, path+"/"+fileName)

			// send the request
			resp, err := http.Get(url) // nolint:gosec
			if err != nil {
				fmt.Fprintf(h.IOStreams.Out, "download file error: %v\n", err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(h.IOStreams.Out, "download file error with status code: %d\n", resp.StatusCode)
				return
			}
			defer resp.Body.Close()

			// check the response
			if resp.ContentLength <= 0 {
				fmt.Fprintf(h.IOStreams.Out, "content length less than 0, aborting download")
				return
			}

			// skip if the file exists
			if _, err := os.Stat(path + "/" + fileName); err == nil {
				fmt.Fprintf(h.IOStreams.Out, "file %s already exists, skipping download\n", path+"/"+fileName)
				return
			}

			// create the file and download
			file, err := os.Create(path + "/" + fileName)
			if err != nil {
				fmt.Fprintf(h.IOStreams.Out, "create file error: %v\n", err)
				return
			}
			defer file.Close()

			if h.IOStreams.CanPrompt {
				err = processDownload(int(resp.ContentLength), file, resp.Body)
				if err != nil {
					if err == util.InterruptError {
						interrupt = true
						return
					}
					fmt.Fprintf(h.IOStreams.Out, "download file error: %v\n", err)
					return
				}
			} else {
				_, err = io.Copy(file, resp.Body)
				if err != nil {
					fmt.Fprintf(h.IOStreams.Out, "download file error: %v\n", err)
					return
				}
			}
		}()
	}
	if interrupt {
		return util.InterruptError
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

func processDownload(contentType int, file *os.File, reader io.Reader) error {
	var p *tea.Program
	pw := &ui.ProgressWriter{
		Total:  contentType,
		File:   file,
		Reader: reader,
		OnProgress: func(ratio float64) {
			p.Send(ui.ProgressMsg(ratio))
		},
	}

	m := ui.ProcessModel{
		Progress: progress.New(progress.WithDefaultGradient()),
	}
	// Start Bubble Tea
	p = tea.NewProgram(m)
	// Start the download
	go pw.Start()
	processModel, err := p.Run()
	if err != nil {
		return errors.Trace(err)
	}
	if processModel.(ui.ProcessModel).Interrupted {
		return util.InterruptError
	}
	if processModel.(ui.ProcessModel).Err != nil {
		return processModel.(ui.ProcessModel).Err
	}
	return nil
}
