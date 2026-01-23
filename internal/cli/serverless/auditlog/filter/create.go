// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
)

var InputDescription = map[string]string{
	flag.DisplayName:        "Input the filter rule name",
	flag.AuditLogFilterRule: "Input the filter rule expression, use \"ticloud serverless audit-log filter template\" to get the template",
}

type FilterRuleWithoutName struct {
	Users   []string                  `json:"users"`
	Filters []auditlog.AuditLogFilter `json:"filters"`
}

type FilterRuleOpts struct {
	interactive bool
}

func (o FilterRuleOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
		flag.AuditLogFilterRule,
	}
}

func (o FilterRuleOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
		flag.AuditLogFilterRule,
	}
}

func (o *FilterRuleOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := o.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			o.interactive = false
			break
		}
	}
	if !o.interactive {
		for _, fn := range o.RequiredFlags() {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := FilterRuleOpts{interactive: true}

	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create an audit log filter rule",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a filter rule in interactive mode:
  $ %[1]s serverless audit-log filter create

  Create a filter rule which filters all audit logs in non-interactive mode:
  $ %[1]s serverless audit-log filter create --cluster-id <cluster-id> --display-name <rule-name> --rule '{"users":["%%@%%"],"filters":[{}]}'

  Create a filter rule which filters QUERY and EXECUTE for test.t and filter QUERY for all tables in non-interactive mode:
  $ %[1]s serverless audit-log filter create --cluster-id <cluster-id> --display-name <rule-name> --rule '{"users":["%%@%%"],"filters":[{"classes":["QUERY","EXECUTE"],"tables":["test.t"]},{"classes":["QUERY"]}]}'
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

			var clusterID, name string
			var filterRule FilterRuleWithoutName
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

				inputs := []string{flag.DisplayName, flag.AuditLogFilterRule}
				textInput, err := ui.InitialInputModel(inputs, InputDescription)
				if err != nil {
					return err
				}
				name = textInput.Inputs[0].Value()
				ruleStr := textInput.Inputs[1].Value()
				if err := json.Unmarshal([]byte(ruleStr), &filterRule); err != nil {
					return errors.New("Invalid rule, please use JSON format. Use \"ticloud serverless audit-log filter template\"" +
						" to see filter templates.")
				}
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				name, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}
				ruleStr, err := cmd.Flags().GetString(flag.AuditLogFilterRule)
				if err != nil {
					return errors.Trace(err)
				}
				if err := json.Unmarshal([]byte(ruleStr), &filterRule); err != nil {
					return errors.New("Invalid rule, please use JSON format. Use \"ticloud serverless audit-log filter template\"" +
						" to see filter templates.")
				}
			}
			if name == "" {
				return errors.New("empty filter rule name, please specify a name")
			}
			if len(filterRule.Users) == 0 {
				return errors.New("empty users, please specify at least one user")
			}
			if len(filterRule.Filters) == 0 {
				return errors.New("empty filters, please specify at least one filter")
			}
			params := auditlog.AuditLogFilterRule{
				DisplayName: name,
				Users:       filterRule.Users,
				Filters:     filterRule.Filters,
			}

			resp, err := d.CreateAuditLogFilterRule(
				ctx,
				clusterID,
				&auditlog.DatabaseAuditLogServiceCreateAuditLogFilterRuleBody{
					FilterRule: params,
				})
			if err != nil {
				return errors.Trace(err)
			}
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("audit log filter rule %s created", *resp.FilterRuleId))
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	cmd.Flags().String(flag.DisplayName, "", "The display name of the filter rule.")
	cmd.Flags().String(flag.AuditLogFilterRule, "", "Filter rule expressions, use \"ticloud serverless audit-log filter template\" to see filter templates.")

	return cmd
}
