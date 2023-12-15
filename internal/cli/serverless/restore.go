package serverless

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-openapi/strfmt"
	"github.com/juju/errors"
	"tidbcloud-cli/internal/service/cloud"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	brApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"
	brModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/models"

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
	err := cmd.MarkFlagRequired(flag.RestoreMode)
	if err != nil {
		return err
	}
	cmd.MarkFlagsMutuallyExclusive(flag.ClusterID, flag.BranchID)
	cmd.MarkFlagsMutuallyExclusive(flag.ClusterID, flag.BackupTime)
	cmd.MarkFlagsRequiredTogether(flag.ClusterID, flag.BackupTime)
	return nil
}

func RestoreCmd(h *internal.Helper) *cobra.Command {
	opts := RestoreOpts{
		interactive: true,
	}

	var restoreCmd = &cobra.Command{
		Use:         "restore",
		Short:       "Restore a serverless cluster",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Restore a serverless cluster in interactive mode):
 $ %[1]s serverless restore

 Restore a serverless cluster with snaphot mode in non-interactive mode:
 $ %[1]s serverless restore --backup-id <back-id>

 Restore a serverless cluster with pointInTime mode in non-interactive mode:
 $ %[1]s serverless restore -p <project-id> -o json`, config.CliName),
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
					project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
					if err != nil {
						return err
					}
					cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
					if err != nil {
						return err
					}
					backup, err := cloud.GetSelectedServerlessBackup(cluster.ID, int32(h.QueryPageSize), d)
					if err != nil {
						return err
					}
					backupID = backup.ID
				} else if restoreMode == cloud.RestoreModePointInTime {
					project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
					if err != nil {
						return err
					}
					cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
					if err != nil {
						return err
					}
					clusterID = cluster.ID

				} else {
					return errors.New("invalid restore mode")
				}
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				backupID, err = cmd.Flags().GetString(flag.BackupID)
				backupTimeStr, err = cmd.Flags().GetString(flag.BackupTime)
			}

			var backupTime strfmt.DateTime
			if backupTimeStr != "" {
				_, err := strfmt.ParseDateTime(backupTimeStr)
				if err != nil {
					return errors.Trace(err)
				}
			}

			params := brApi.NewBackupRestoreServiceRestoreParams().WithBody(&brModel.V1beta1RestoreRequest{
				Snapshot: &brModel.RestoreRequestSnapshot{
					BackupID: backupID,
				},
				PointInTime: &brModel.RestoreRequestPointInTime{
					ClusterID:  clusterID,
					BackupTime: backupTime,
				},
			})
			resp, err := d.Restore(params)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("restore to clsuter %s, use \"ticloud serverless get -c %s\" to check the cluster status", *resp.Payload.ClusterID, *resp.Payload.ClusterID)))
			return nil
		},
	}

	restoreCmd.Flags().String(flag.BackupID, "", "The ID of the backup, used with Snapshot restore mode.")
	restoreCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of cluster, used with PointInTime restore mode. Please specify the --backup-time together.")
	restoreCmd.Flags().String(flag.BackupTime, "", "The time to restore to (e.g. 2023-12-13T07:00:00Z), used with PointInTime restore mode. Please specify the --cluster-id together.")
	return restoreCmd
}
