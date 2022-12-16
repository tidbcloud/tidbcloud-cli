package dataimport

import (
	"fmt"
	"strconv"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/service/cloud"
	importModel "tidbcloud-cli/pkg/tidbcloud/import/models"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type ListOpts struct {
	interactive bool
}

func (c ListOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
	}
}

func ListCmd(h *internal.Helper) *cobra.Command {
	opts := ListOpts{}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List a data import task",
		Example: fmt.Sprintf(`  List an import task in interactive mode:
  $ %[1]s import list

  List an import task in non-interactive mode:
  $ %[1]s import list --project-id <project-id> --cluster-name <cluster-name> --aws-role-arn <aws-role-arn> --data-format <data-format> --source-url <source-url>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// mark required flags in non-interactive mode
			if !opts.interactive {
				for _, fn := range flags {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectID, clusterID string
			d, err := h.Client()
			if err != nil {
				return err
			}

			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				projectID = project.ID

				cluster, err := cloud.GetSelectedCluster(projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID
			} else {
				// non-interactive mode
				projectID = cmd.Flag(flag.ProjectID).Value.String()
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
			}

			total, importTasks, err := cloud.RetrieveImports(projectID, clusterID, h.QueryPageSize, d)
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
				res := &importModel.OpenapiListImportsResp{
					Imports: importTasks,
					Total:   &total,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"Status",
					"CreatedAt",
					"SourceURL",
					"DataFormat",
					"ClusterID",
				}

				var rows []output.Row
				for _, item := range importTasks {

					rows = append(rows, output.Row{
						item.ID,
						string(*item.Status),
						item.CreatedAt.String(),
						*item.SourceURL,
						strconv.FormatUint(uint64(*item.DataFormat), 10),
						*item.ClusterID,
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

	return listCmd
}
