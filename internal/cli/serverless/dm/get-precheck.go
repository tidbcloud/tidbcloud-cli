// Copyright 2025 PingCAP, Inc.
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

package dm

import (
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type GetPrecheckOpts struct {
	interactive bool
}

func (c GetPrecheckOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrecheckID,
	}
}

func GetPrecheckCmd(h *internal.Helper) *cobra.Command {
	opts := GetPrecheckOpts{
		interactive: true,
	}

	var getPrecheckCmd = &cobra.Command{
		Use:         "get-precheck",
		Short:       "Get precheck result for a DM task",
		Aliases:     []string{"describe-precheck"},
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Get precheck result in interactive mode:
  $ %[1]s serverless dm get-precheck

  Get precheck result in non-interactive mode:
  $ %[1]s serverless dm get-precheck --cluster-id <cluster-id> --precheck-id <precheck-id>`,
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
			var clusterID, precheckID string
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
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

				// TODO: Add interactive selection for precheck
				return errors.New("Interactive mode for DM get-precheck is not yet implemented. Please use non-interactive mode with --precheck-id flag")
			} else {
				// non-interactive mode
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				precheckID = cmd.Flag(flag.PrecheckID).Value.String()
			}

			cmd.Annotations[telemetry.ClusterID] = clusterID

			precheck, err := d.GetPrecheck(ctx, clusterID, precheckID)
			if err != nil {
				return errors.Trace(err)
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				err := output.PrintJson(h.IOStreams.Out, precheck)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				return fmt.Errorf("unsupported output format: %s", format)
			}

			return nil
		},
	}

	getPrecheckCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	getPrecheckCmd.Flags().String(flag.PrecheckID, "", "Precheck ID.")
	getPrecheckCmd.Flags().StringP(flag.Output, flag.OutputShort, output.JsonFormat, flag.OutputHelp)
	return getPrecheckCmd
}
