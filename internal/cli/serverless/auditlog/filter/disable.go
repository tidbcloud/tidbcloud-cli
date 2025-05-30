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
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
)

type DisableFilterRuleOpts struct {
	interactive bool
}

func (o DisableFilterRuleOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.AuditLogFilterRuleName,
	}
}

func (o *DisableFilterRuleOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := o.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			o.interactive = false
			break
		}
	}
	if !o.interactive {
		for _, fn := range flags {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DisableCmd(h *internal.Helper) *cobra.Command {
	opts := DisableFilterRuleOpts{
		interactive: true,
	}

	var disableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disable an audit log filter rule",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Disable an audit log filter rule in interactive mode:
  $ %[1]s serverless auditlog filter-rule disable

  Disable an audit log filter rule in non-interactive mode:
  $ %[1]s serverless auditlog filter-rule disable --cluster-id <cluster-id> --rule-name <rule-name>
`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, ruleName string
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
				ruleName, err = cloud.GetSelectedRuleName(ctx, cluster.ID, d)
				if err != nil {
					return err
				}
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				ruleName, err = cmd.Flags().GetString(flag.AuditLogFilterRuleName)
				if err != nil {
					return errors.Trace(err)
				}
			}

			disabled := true
			_, err = d.UpdateAuditLogFilterRule(ctx, clusterID, ruleName, &auditlog.AuditLogServiceUpdateAuditLogFilterRuleBody{
				Disabled: *auditlog.NewNullableBool(&disabled),
			})
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("filter rule %s is disabled", ruleName)))
			return nil
		},
	}

	disableCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	disableCmd.Flags().String(flag.AuditLogFilterRuleName, "", "The name of the filter rule.")

	return disableCmd
}
