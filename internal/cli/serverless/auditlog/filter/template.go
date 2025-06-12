/*
Copyright 2025 PingCAP, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package filter

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
)

type TemplateOpts struct {
	interactive bool
}

func (o TemplateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

const (
	FilterTemplate    = `{"users":["%s.root@%%"],"filters":[{"classes":["QUERY"],"tables":["test.t"],"statusCodes":[1]}]}`
	FilterAllTemplate = `{"users":["%@%"],"filters":[{}]}`
)

func TemplateCmd(h *internal.Helper) *cobra.Command {
	opts := TemplateOpts{interactive: true}

	var cmd = &cobra.Command{
		Use:   "template",
		Short: "Show audit log filter rule templates",
		Example: fmt.Sprintf(`  Show filter templates in interactive mode:
  $ %[1]s serverless audit-log filter template

  Show filter templates in non-interactive mode:
  $ %[1]s serverless audit-log filter template --cluster-id <cluster_id>`, config.CliName),
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
				cid, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID = cid
			}

			c, err := d.GetCluster(ctx, clusterID, cluster.SERVERLESSSERVICEGETCLUSTERVIEWPARAMETER_BASIC)
			if err != nil {
				return errors.Annotatef(err, "failed to get cluster %s", clusterID)
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Basic template:"))
			fmt.Fprintln(h.IOStreams.Out, fmt.Sprintf(FilterTemplate, *c.UserPrefix))
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Filter all template:"))
			fmt.Fprintln(h.IOStreams.Out, FilterAllTemplate)
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	return cmd
}
