package serverless

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type ListAuditLogOpts struct {
	interactive bool
}

func (c ListAuditLogOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func ListAuditLogCmd(h *internal.Helper) *cobra.Command {
	opts := ListOpts{
		interactive: true,
	}

	var listCmd = &cobra.Command{
		Use:         "list-auditlog",
		Short:       "List TiDB Cloud Serverless cluster database audit logs",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  List TiDB Cloud Serverless cluster database audit logs in interactive mode):
  $ %[1]s serverless list_auditlog

  List all TiDB Cloud Serverless clusters in non-interactive mode:
  $ %[1]s serverless list_auditlog -c <cluster-id>

  List all TiDB Cloud Serverless clusters in non-interactive mode:
  $ %[1]s serverless list_auditlog -c <cluster-id> -o json`, config.CliName),
		Aliases: []string{"ls-al"},
		PreRun: func(cmd *cobra.Command, args []string) {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID string
			ctx := cmd.Context()
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
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			total, items, err := cloud.RetrieveClusters(context, pID, h.QueryPageSize, d)
			if err != nil {
				return err
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			// for terminal which can prompt, humanFormat is the default format.
			// for other terminals, json format is the default format.
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				res := &cluster.TidbCloudOpenApiserverlessv1beta1ListClustersResponse{
					Clusters:  items,
					TotalSize: &total,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"DisplayName",
					"State",
					"Version",
					"Cloud",
					"Region",
					"CreateTime",
				}

				var rows []output.Row
				for _, item := range items {
					rows = append(rows, output.Row{
						*item.ClusterId,
						item.DisplayName,
						string(*item.State),
						*item.Version,
						string(*item.Region.CloudProvider),
						*item.Region.DisplayName,
						item.CreateTime.Format(time.RFC3339),
					})
				}

				err := output.PrintHumanTable(h.IOStreams.Out, columns, rows)
				if err != nil {
					return errors.Trace(err)
				}
				return nil
			} else {
				return fmt.Errorf("unsupported output format: %s", format)
			}

			return nil
		},
	}

	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	listCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The ID of the project.")
	return listCmd
}
