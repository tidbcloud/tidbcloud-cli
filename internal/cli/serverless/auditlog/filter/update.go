// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	alutil "github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/auditlog/util"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
)

type UpdateFilterRuleOpts struct {
	interactive bool
}

func (o UpdateFilterRuleOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.AuditLogFilterRuleName,
		flag.AuditLogFilterRule,
	}
}

func (o UpdateFilterRuleOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.AuditLogFilterRuleName,
	}
}

func (o *UpdateFilterRuleOpts) MarkInteractive(cmd *cobra.Command) error {
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
		cmd.MarkFlagsOneRequired(flag.AuditLogFilterRule, flag.Enabled)
	}
	return nil
}

type mutableField string

const (
	Rule    mutableField = "rule"
	Enabled mutableField = "enabled"
)

var mutableFilterRuleFields = []string{
	string(Rule),
	string(Enabled),
}

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateFilterRuleOpts{interactive: true}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an audit log filter rule",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Update an audit log filter rule in interactive mode:
  $ %[1]s serverless auditlog filter-rule update

  Enable audit log filter rule in non-interactive mode:
  $ %[1]s serverless auditlog filter-rule update --cluster-id <cluster-id> --name <rule-name> --enabled

  Disable audit log filter rule in non-interactive mode:
  $ %[1]s serverless auditlog filter-rule update --cluster-id <cluster-id> --name <rule-name> --enabled=false

  Update filters of an audit log filter rule in non-interactive mode:
  $ %[1]s serverless auditlog filter-rule update --cluster-id <cluster-id> --name <rule-name> --rule '{"users":["%%@%%"],"filters":[{"classes":["QUERY"],"tables":["test.t"]}]}'
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
			var filterRule *FilterRuleWithoutName
			var enabled *bool
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

				// Select rule name
				ruleName, err = cloud.GetSelectedRuleName(ctx, clusterID, d)
				if err != nil {
					return err
				}

				// Select field to update
				fieldName, err := cloud.GetSelectedField(mutableFilterRuleFields, "Choose the field to update:")
				if err != nil {
					return err
				}

				switch fieldName {
				case string(Enabled):
					prompt := &survey.Confirm{
						Message: "enable the database audit logging?",
						Default: false,
					}
					var confirm bool
					err = survey.AskOne(prompt, &confirm)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					if confirm {
						enabled = aws.Bool(true)
					} else {
						enabled = aws.Bool(false)
					}
				case string(Rule):
					inputs := []string{flag.AuditLogFilterRule}
					textInput, err := ui.InitialInputModel(inputs, alutil.InputDescription)
					if err != nil {
						return err
					}
					ruleStr := textInput.Inputs[0].Value()
					var f FilterRuleWithoutName
					if err := json.Unmarshal([]byte(ruleStr), &f); err != nil {
						return errors.New("invalid filter, please use JSON format")
					}
					filterRule = &f
				default:
					return errors.Errorf("invalid field %s", fieldName)
				}
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				ruleName, err = cmd.Flags().GetString(flag.AuditLogFilterRuleName)
				if err != nil {
					return errors.Trace(err)
				}
				if cmd.Flags().Changed(flag.Enabled) {
					u, err := cmd.Flags().GetBool(flag.Enabled)
					if err != nil {
						return errors.Trace(err)
					}
					enabled = &u
				}
				if cmd.Flags().Changed(flag.AuditLogFilterRule) {
					ruleStr, err := cmd.Flags().GetString(flag.AuditLogFilterRule)
					if err != nil {
						return errors.Trace(err)
					}
					var f FilterRuleWithoutName
					if err := json.Unmarshal([]byte(ruleStr), &f); err != nil {
						return errors.New("invalid filter, please use JSON format")
					}
					filterRule = &f
				}
			}
			if enabled == nil && filterRule == nil {
				return errors.New("nothing to update")
			}
			if filterRule != nil {
				if len(filterRule.Users) == 0 {
					return errors.New("empty users, please specify at least one user")
				}
				if len(filterRule.Filters) == 0 {
					return errors.New("empty filters, please specify at least one filter")
				}
			}
			body := &auditlog.AuditLogServiceUpdateAuditLogFilterRuleBody{}
			if filterRule != nil {
				body.Users = filterRule.Users
				body.Filters = filterRule.Filters
			}
			if enabled != nil {
				disabled := !*enabled
				body.Disabled = *auditlog.NewNullableBool(&disabled)
			}
			_, err = d.UpdateAuditLogFilterRule(ctx, clusterID, ruleName, body)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("filter rule %s updated", ruleName)))
			return nil
		},
	}

	updateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	updateCmd.Flags().String(flag.AuditLogFilterRuleName, "", "The name of the filter rule to update.")
	updateCmd.Flags().String(flag.AuditLogFilterRule, "", "Complete filter rule expressions, use `ticloud serverless audit-log filter template` to see filter templates")
	updateCmd.Flags().Bool(flag.Enabled, false, "Enable or disable the filter rule.")

	return updateCmd
}
