package serverless

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

var inputDescription = map[string]string{
	flag.OutputPath: "Input the path where you want to download to. If not specified, download to the current directory.",
	flag.StartDate:  "Input the start date of the audit log you want to download in the format of 'YYYY-MM-DD', e.g. '2025-01-01'.",
	flag.EndDate:    "Input the end date of the audit log you want to download in the format of 'YYYY-MM-DD', e.g. '2025-01-01'.",
}

const (
	DefaultConcurrency = 3
	MaxBatchSize       = 100
	MaxDateRange       = 7 * 24 * time.Hour
)

var DownloadPathInputFields = map[string]int{
	flag.OutputPath: 0,
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

func DownloadAuditLogCmd(h *internal.Helper) *cobra.Command {
	var force bool
	opts := DownloadAuditLogOpts{
		interactive: true,
	}

	var downloadAuditLogCmd = &cobra.Command{
		Use:   "download-auditlog",
		Short: "Download the database audit logs",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Download the database audit logs in interactive mode:
  $ %[1]s serverless download-auditlog

  Download the database audit logs in non-interactive mode:
  $ %[1]s serverless download-auditlog -c <cluster-id> --date <date>`, config.CliName),
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

				pathInput, err := ui.InitialInputModel([]string{flag.OutputPath}, inputDescription)
				if err != nil {
					return err
				}
				path = pathInput.Inputs[0].Value()

				dateInput, err := ui.InitialInputModel([]string{flag.StartDate, flag.EndDate}, inputDescription)
				if err != nil {
					return err
				}
				startDate = dateInput.Inputs[0].Value()
				endDate = dateInput.Inputs[1].Value()
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
				totalSize += *log.Size
				auditLogNames = append(auditLogNames, *log.Name)
			}
			// ask for confirmation
			fileMessage := fmt.Sprintf("There are %d files to download, total size is %s.", len(auditLogNames), humanize.IBytes(uint64(totalSize)))
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

			// download the audit logs
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

	downloadAuditLogCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	downloadAuditLogCmd.Flags().String(flag.OutputPath, "", "The path where you want to download to. If not specified, download to the current directory.")
	downloadAuditLogCmd.Flags().BoolVar(&force, flag.Force, false, "Download without confirmation.")
	downloadAuditLogCmd.Flags().Int(flag.Concurrency, 3, "Download concurrency.")
	downloadAuditLogCmd.Flags().String(flag.StartDate, "", "The start date of the audit log you want to download in the format of 'YYYY-MM-DD', e.g. '2025-01-01'.")
	downloadAuditLogCmd.Flags().String(flag.EndDate, "", "The end date of the audit log you want to download in the format of 'YYYY-MM-DD', e.g. '2025-01-01'.")

	return downloadAuditLogCmd
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

	// init the concurrency progress model
	var p *tea.Program
	m := NewProcessDownloadModel(
		concurrency,
		path,
		exportID,
		clusterID,
		client,
		int(totalSize),
		fileNames,
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
	fmt.Fprint(h.IOStreams.Out, generateDownloadSummary(succeededCount, skippedCount, failedCount))
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
	exportDownloadPool, err := NewDownloadPool(h, path, concurrency, exportID, clusterID, fileNames, client)
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
