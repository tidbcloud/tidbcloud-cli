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

package serverless

import (
	"fmt"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/br"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type RestoreOpts struct {
	interactive bool
}

func (c RestoreOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.BackupTime,
		flag.BackupID,
	}
}

func (c *RestoreOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := c.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	// Mark flags
	if !c.interactive {
		cmd.MarkFlagsMutuallyExclusive(flag.BackupID, flag.ClusterID)
		cmd.MarkFlagsMutuallyExclusive(flag.BackupID, flag.BackupTime)
		if !cmd.Flags().Changed(flag.BackupID) {
			cmd.MarkFlagsRequiredTogether(flag.ClusterID, flag.BackupTime)
		}
	}
	return nil
}

func RestoreCmd(h *internal.Helper) *cobra.Command {
	opts := RestoreOpts{
		interactive: true,
	}

	var restoreCmd = &cobra.Command{
		Use:         "restore",
		Short:       "Restore a TiDB Cloud Serverless cluster",
		Annotations: make(map[string]string),
		Args:        cobra.NoArgs,
		Example: fmt.Sprintf(`  Restore a TiDB Cloud Serverless cluster in interactive mode:
 $ %[1]s serverless restore

 Restore a TiDB Cloud Serverless cluster with snaphot mode in non-interactive mode:
 $ %[1]s serverless restore --backup-id <backup-id>

 Restore a TiDB Cloud Serverless cluster with pointInTime mode in non-interactive mode:
 $ %[1]s serverless restore -c <cluster-id> --backup-time <backup-time>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.MarkInteractive(cmd)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID string
			var backupID string
			var backupTimeStr string
			var restoreMode string
			if opts.interactive {
				restoreMode, err = cloud.GetSelectedRestoreMode()
				if err != nil {
					return err
				}

				if restoreMode == cloud.RestoreModeSnapshot {
					project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
					if err != nil {
						return err
					}
					cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
					if err != nil {
						return err
					}
					backup, err := cloud.GetSelectedServerlessBackup(ctx, cluster.ID, int32(h.QueryPageSize), d)
					if err != nil {
						return err
					}
					backupID = backup.ID
				} else if restoreMode == cloud.RestoreModePointInTime {
					project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
					if err != nil {
						return err
					}
					cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
					if err != nil {
						return err
					}
					clusterID = cluster.ID
					// variables for input
					inputModel, err := GetRestoreInput()
					if err != nil {
						return err
					}
					backupTimeStr = inputModel.(ui.TextInputModel).Inputs[0].Value()
				} else {
					return errors.New("invalid restore mode")
				}
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				backupID, err = cmd.Flags().GetString(flag.BackupID)
				if err != nil {
					return errors.Trace(err)
				}
				backupTimeStr, err = cmd.Flags().GetString(flag.BackupTime)
				if err != nil {
					return errors.Trace(err)
				}
			}

			body := &br.V1beta1RestoreRequest{}
			if backupID != "" {
				body.Snapshot = &br.RestoreRequestSnapshot{
					BackupId: &backupID,
				}
			} else {
				if backupTimeStr == "" {
					return errors.New("backup time is required in point-in-time mode")
				}
				backupTime, err := time.Parse(time.RFC3339, backupTimeStr)
				if err != nil {
					return errors.New(fmt.Sprintf("invalid backup time %s. Please input the backup time with the 2006-01-02T15:04:05Z formate", backupTimeStr))
				}
				body.PointInTime = &br.RestoreRequestPointInTime{
					ClusterId:  &clusterID,
					BackupTime: &backupTime,
				}
			}
			resp, err := d.Restore(ctx, body)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("restore to clsuter %s, use \"ticloud serverless get -c %s\" to check the restore process", resp.ClusterId, resp.ClusterId)))
			return nil
		},
	}

	restoreCmd.Flags().String(flag.BackupID, "", "The ID of the backup. Used in snapshot restore mode.")
	restoreCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of cluster. Used in point-in-time restore mode. Please specify the --backup-time together.")
	restoreCmd.Flags().String(flag.BackupTime, "", "The time to restore to (e.g. 2023-12-13T07:00:00Z). Used with point-in-time restore mode. Please specify the --cluster-id together.")
	return restoreCmd
}

func initialRestoreInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 1),
	}
	backupTime := textinput.New()
	backupTime.Placeholder = "Backup Time. e.g., 2023-12-13T07:00:00Z"
	backupTime.Focus()
	backupTime.PromptStyle = config.FocusedStyle
	backupTime.TextStyle = config.FocusedStyle
	m.Inputs[0] = backupTime
	return m
}

func GetRestoreInput() (tea.Model, error) {
	p := tea.NewProgram(initialRestoreInputModel())
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}
