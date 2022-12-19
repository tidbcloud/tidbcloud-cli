package dataimport

import (
	"fmt"
	"strconv"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	importOp "tidbcloud-cli/pkg/tidbcloud/import/client/import_service"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type CancelOpts struct {
	interactive bool
}

func (c CancelOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
		flag.ImportID,
	}
}

func CancelCmd(h *internal.Helper) *cobra.Command {
	opts := CancelOpts{
		interactive: true,
	}
	var cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a data import task",
		Example: fmt.Sprintf(`  Cancel an import task in interactive mode:
  $ %[1]s import cancel

  Cancel an import task in non-interactive mode:
  $ %[1]s import cancel --project-id <project-id> --cluster-name <cluster-name> --import-id <import-id>`,
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
				importID = strconv.FormatUint(selectedImport.ID, 10)
			} else {
				// non-interactive mode
				projectID = cmd.Flag(flag.ProjectID).Value.String()
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				importID = cmd.Flag(flag.ImportID).Value.String()
			}

			params := importOp.NewCancelImportParams().WithProjectID(projectID).WithClusterID(clusterID).WithID(importID)
			_, err = d.CancelImport(params)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s is canceled.", importID))
			return nil
		},
	}

	cancelCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "Project ID")
	cancelCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	cancelCmd.Flags().String(flag.ImportID, "", "The ID of import task")
	return cancelCmd
}
