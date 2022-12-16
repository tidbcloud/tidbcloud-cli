package dataimport

import (
	"encoding/json"
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	importOp "tidbcloud-cli/pkg/tidbcloud/import/client/import_service"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type DescribeOpts struct {
	interactive bool
}

func (c DescribeOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
		flag.ImportID,
	}
}

func DescribeCmd(h *internal.Helper) *cobra.Command {
	opts := DescribeOpts{}

	var describeCmd = &cobra.Command{
		Use:   "describe",
		Short: "Describe a data import task",
		Example: fmt.Sprintf(`  Describe an import task in interactive mode:
  $ %[1]s import describe

  Describe an import task in non-interactive mode:
  $ %[1]s import describe --project-id <project-id> --cluster-name <cluster-name> --aws-role-arn <aws-role-arn> --data-format <data-format> --source-url <source-url>`,
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
			var projectID, clusterID, importID string
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

				selectedImport, err := cloud.GetSelectedImport(projectID, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				importID = selectedImport.ID
			} else {
				// non-interactive mode
				projectID = cmd.Flag(flag.ProjectID).Value.String()
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				importID = cmd.Flag(flag.ImportID).Value.String()
			}

			params := importOp.NewGetImportParams().WithProjectID(projectID).WithClusterID(clusterID).WithID(importID)
			importTask, err := d.GetImport(params)
			if err != nil {
				return errors.Trace(err)
			}

			v, err := json.MarshalIndent(importTask.Payload, "", "  ")
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, string(v))
			return nil
		},
	}

	return describeCmd
}
