package migration

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type PauseOpts struct {
	interactive bool
}

func (c PauseOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.MigrationTaskID,
	}
}

func (c *PauseOpts) MarkInteractive(cmd *cobra.Command) error {
	for _, fn := range c.NonInteractiveFlags() {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	if !c.interactive {
		for _, fn := range c.NonInteractiveFlags() {
			if err := cmd.MarkFlagRequired(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func PauseCmd(h *internal.Helper) *cobra.Command {
	opts := PauseOpts{interactive: true}

	var cmd = &cobra.Command{
		Use:   "pause",
		Short: "Pause a migration task",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Pause a migration task in interactive mode:
  $ %[1]s serverless migration pause

  Pause a migration task in non-interactive mode:
  $ %[1]s serverless migration pause -c <cluster-id> --migration-id <task-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, taskID string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID
				task, err := cloud.GetSelectedMigrationTask(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				taskID = task.ID
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				taskID, err = cmd.Flags().GetString(flag.MigrationTaskID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			emptyBody := map[string]interface{}{}
			if _, err := d.PauseMigrationTask(ctx, clusterID, taskID, &emptyBody); err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintf(h.IOStreams.Out, "migration task %s paused\n", taskID)
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID that owns the migration task.")
	cmd.Flags().StringP(flag.MigrationTaskID, flag.MigrationTaskIDShort, "", "ID of the migration task to pause.")
	return cmd
}
