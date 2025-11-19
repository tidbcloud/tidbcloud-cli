package migration

import (
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type ListOpts struct {
	interactive bool
}

func (c ListOpts) NonInteractiveFlags() []string {
	return []string{flag.ClusterID}
}

func (c *ListOpts) MarkInteractive(cmd *cobra.Command) error {
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

func ListCmd(h *internal.Helper) *cobra.Command {
	opts := ListOpts{interactive: true}

	var cmd = &cobra.Command{
		Use:     "list",
		Short:   "List migration tasks",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  List migrations in interactive mode:
  $ %[1]s serverless migration list

  List migrations in non-interactive mode with JSON output:
  $ %[1]s serverless migration list -c <cluster-id> -o json`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID string
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
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			pageSize := int32(h.QueryPageSize)
			resp, err := d.ListMigrationTasks(ctx, clusterID, &pageSize, nil, nil)
			if err != nil {
				return errors.Trace(err)
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				return errors.Trace(output.PrintJson(h.IOStreams.Out, resp))
			}

			if format != output.HumanFormat {
				return fmt.Errorf("unsupported output format: %s", format)
			}

			columns := []output.Column{"ID", "Name", "Mode", "State", "CreatedAt"}
			var rows []output.Row
			for _, task := range resp.Tasks {
				id := safeString(task.Id)
				name := safeString(task.Name)
				if name == "" {
					name = id
				}
				mode := ""
				if task.Mode != nil {
					mode = string(*task.Mode)
				}
				state := ""
				if task.State != nil {
					state = string(*task.State)
				}
				created := formatTime(task.CreateTime)
				rows = append(rows, output.Row{id, name, mode, state, created})
			}
			return errors.Trace(output.PrintHumanTable(h.IOStreams.Out, columns, rows))
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID of the migration tasks to list.")
	cmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return cmd
}

func safeString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func formatTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(time.RFC3339)
}
