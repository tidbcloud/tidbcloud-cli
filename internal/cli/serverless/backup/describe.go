package backup

import (
	"fmt"
	"tidbcloud-cli/internal/output"

	"github.com/spf13/cobra"
	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	brAPI "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"

	"github.com/juju/errors"
)

type DescribeOpts struct {
	interactive bool
}

func (c DescribeOpts) NonInteractiveFlags() []string {
	return []string{
		flag.BackupID,
	}
}

func (c *DescribeOpts) MarkInteractive(cmd *cobra.Command) error {
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
		for _, fn := range flags {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DescribeCmd(h *internal.Helper) *cobra.Command {
	opts := DescribeOpts{
		interactive: true,
	}

	var describeCmd = &cobra.Command{
		Use:     "describe",
		Short:   "Describe a serverless cluster backup",
		Aliases: []string{"get"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Get the backup in interactive mode:
  $ %[1]s serverless backup describe

  Get the backup in non-interactive mode:
  $ %[1]s serverless backup describe --backup-id <backup-id>`, config.CliName),
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

			var backupID string
			var clusterID string
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

				backup, err := cloud.GetSelectedServerlessBackup(clusterID, int32(h.QueryPageSize), d)
				if err != nil {
					return err
				}
				backupID = backup.ID
			} else {
				// non-interactive mode, get values from flags
				backupID, err = cmd.Flags().GetString(flag.BackupID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			params := brAPI.NewBackupRestoreServiceGetBackupParams().WithBackupID(backupID)
			if err != nil {
				return errors.Trace(err)
			}
			backup, err := d.GetBackup(params)
			if err != nil {
				return errors.Trace(err)
			}
			err = output.PrintJson(h.IOStreams.Out, backup.Payload)
			if err != nil {
				return errors.Trace(err)
			}

			return nil
		},
	}

	describeCmd.Flags().String(flag.BackupID, "", "The ID of the backup to be described")
	return describeCmd
}
