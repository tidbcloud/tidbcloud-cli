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

package auditlog

import (
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
)

const (
	DefaultConcurrency = 3
	MaxBatchSize       = 100
	MaxDateRange       = 7 * 24 * time.Hour
)

var DownloadPathInputFields = map[string]int{
	flag.OutputPath: 0,
}

var InputDescription = map[string]string{
	flag.OutputPath: "Input the download path, press Enter to skip and download to the current directory",
	flag.StartDate:  "Input the start date of the download in the format of 'YYYY-MM-DD'",
	flag.EndDate:    "Input the end date of the download in the format of 'YYYY-MM-DD'",
}

type DownloadAuditLogOpts struct {
	interactive bool
}

func (c DownloadAuditLogOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.StartDate,
		flag.EndDate,
		flag.OutputPath,
	}
}

func (c DownloadAuditLogOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func (c *DownloadAuditLogOpts) MarkInteractive(cmd *cobra.Command) error {
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
	opts := DownloadAuditLogOpts{
		interactive: true,
	}

	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download the database audit log files",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Download the database audit logs in interactive mode:
  $ %[1]s serverless audit-log download

  Download the database audit logs in non-interactive mode:
  $ %[1]s serverless audit-log download -c <cluster-id> --start-date <start-date> --end-date <end-date>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var startDate, endDate, clusterID, path string
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

				inputs := []string{flag.OutputPath, flag.StartDate, flag.EndDate}
				textInput, err := ui.InitialInputModel(inputs, InputDescription)
				if err != nil {
					return err
				}
				path = textInput.Inputs[0].Value()
				startDate = textInput.Inputs[1].Value()
				endDate = textInput.Inputs[2].Value()
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				path, err = cmd.Flags().GetString(flag.OutputPath)
				if err != nil {
					return errors.Trace(err)
				}
				startDate, err = cmd.Flags().GetString(flag.StartDate)
				if err != nil {
					return errors.Trace(err)
				}
				endDate, err = cmd.Flags().GetString(flag.EndDate)
				if err != nil {
					return errors.Trace(err)
				}
			}

			concurrency, err := cmd.Flags().GetInt(flag.Concurrency)
			if err != nil {
				return errors.Trace(err)
			}
			// check the date
			if err := checkDate(startDate, endDate); err != nil {
				return errors.Trace(err)
			}

			// list the audit logs
			auditLogs, err := cloud.GetAllAuditLogs(ctx, clusterID, startDate, endDate, d)
			if err != nil {
				return errors.Trace(err)
			}
			var totalSize int64
			auditLogNames := make([]string, 0)
			for _, log := range auditLogs {
				totalSize += *log.SizeBytes
				auditLogNames = append(auditLogNames, *log.Name)
			}

			// verify the size
			if len(auditLogNames) == 0 {
				fmt.Fprintln(h.IOStreams.Out, color.RedString("No audit logs found in the specified date range, skip the download."))
				return nil
			}

			// ask for confirmation
			fileMessage := fmt.Sprintf("There are %d files to download, total size is %s.", len(auditLogNames), humanize.IBytes(uint64(totalSize)))
			if !force {
				if !h.IOStreams.CanPrompt {
					return fmt.Errorf("the terminal doesn't support prompt, please run with --force to download")
				}

				confirmed := "yes"

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

			//download the audit logs
			if h.IOStreams.CanPrompt {
				err = DownloadFilesPrompt(h, path, concurrency, clusterID, totalSize, auditLogNames, d)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				err = DownloadFilesWithoutPrompt(h, path, concurrency, clusterID, auditLogNames, d)
				if err != nil {
					return errors.Trace(err)
				}
			}
			return nil
		},
	}

	downloadCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	downloadCmd.Flags().String(flag.OutputPath, "", "The path where you want to download to. If not specified, download to the current directory.")
	downloadCmd.Flags().BoolVar(&force, flag.Force, false, "Download without confirmation.")
	downloadCmd.Flags().Int(flag.Concurrency, 3, "Download concurrency.")
	downloadCmd.Flags().String(flag.StartDate, "", "The start date of the audit log you want to download in the format of 'YYYY-MM-DD', e.g. '2025-01-01'.")
	downloadCmd.Flags().String(flag.EndDate, "", "The end date of the audit log you want to download in the format of 'YYYY-MM-DD', e.g. '2025-01-01'.")

	return downloadCmd
}

func checkDate(startDate, endDate string) error {
	if startDate == "" {
		return errors.New("start date is required")
	}
	if endDate == "" {
		return errors.New("end date is required")
	}
	st, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return errors.New("invalid start date, please input the date in the format of 'YYYY-MM-DD'")
	}
	et, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return errors.New("invalid end date, please input the date in the format of 'YYYY-MM-DD'")
	}
	if st.Add(MaxDateRange).Before(et) {
		return errors.New("the date range should be within 7 days")
	}
	return nil
}
