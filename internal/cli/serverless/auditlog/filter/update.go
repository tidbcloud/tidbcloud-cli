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
		flag.AuditLogFilterRuleID,
		flag.AuditLogFilterRule,
		flag.DisplayName,
	}
}

func (o UpdateFilterRuleOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.AuditLogFilterRuleID,
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
	Rule        mutableField = "rule"
	Enabled     mutableField = "enabled"
	DisplayName mutableField = "display-name"
)

var mutableFilterRuleFields = []string{
	string(Rule),
	string(Enabled),
	string(DisplayName),
}

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateFilterRuleOpts{interactive: true}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an audit log filter rule",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Update an audit log filter rule in interactive mode:
  $ %[1]s serverless audit-log filter update

  Enable audit log filter rule in non-interactive mode:
  $ %[1]s serverless audit-log filter update --cluster-id <cluster-id> --filter-rule-id <rule-id> --enabled

  Disable audit log filter rule in non-interactive mode:
  $ %[1]s serverless audit-log filter update --cluster-id <cluster-id> --filter-rule-id <rule-id> --enabled=false

  Update filters of an audit log filter rule in non-interactive mode:
  $ %[1]s serverless audit-log filter update --cluster-id <cluster-id> --filter-rule-id <rule-id> --rule '{"users":["%%@%%"],"filters":[{"classes":["QUERY"],"tables":["test.t"]}]}'
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

			var clusterID, ruleID, dispalyName string
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
				rule, err := cloud.GetSelectedFilterRule(ctx, clusterID, d)
				if err != nil {
					return err
				}
				ruleID = rule.FilterRuleId

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
					textInput, err := ui.InitialInputModel(inputs, InputDescription)
					if err != nil {
						return err
					}
					ruleStr := textInput.Inputs[0].Value()
					var f FilterRuleWithoutName
					if err := json.Unmarshal([]byte(ruleStr), &f); err != nil {
						return errors.New("Invalid rule, please use JSON format. Use \"ticloud serverless audit-log filter template\"" +
							" to see filter templates.")
					}
					filterRule = &f
				case string(DisplayName):
					inputs := []string{flag.DisplayName}
					textInput, err := ui.InitialInputModel(inputs, InputDescription)
					if err != nil {
						return err
					}
					dispalyName = textInput.Inputs[0].Value()
				default:
					return errors.Errorf("invalid field %s", fieldName)
				}
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				ruleID, err = cmd.Flags().GetString(flag.AuditLogFilterRuleID)
				if err != nil {
					return errors.Trace(err)
				}
				dispalyName, err = cmd.Flags().GetString(flag.DisplayName)
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
						return errors.New("Invalid rule, please use JSON format. Use \"ticloud serverless audit-log filter template\"" +
							" to see filter templates.")
					}
					filterRule = &f
				}
			}
			if enabled == nil && filterRule == nil && dispalyName == "" {
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
			body := &auditlog.DatabaseAuditLogServiceUpdateAuditLogFilterRuleBody{}
			if filterRule != nil {
				body.Users = filterRule.Users
				body.Filters = filterRule.Filters
			}
			if enabled != nil {
				disabled := !*enabled
				body.Disabled = *auditlog.NewNullableBool(&disabled)
			}
			if dispalyName != "" {
				body.DisplayName = &dispalyName
			}
			_, err = d.UpdateAuditLogFilterRule(ctx, clusterID, ruleID, body)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("filter rule %s updated", ruleID)))
			return nil
		},
	}

	updateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	updateCmd.Flags().String(flag.AuditLogFilterRuleID, "", "The ID of the filter rule.")
	updateCmd.Flags().String(flag.AuditLogFilterRule, "", "Complete filter rule expressions, use \"ticloud serverless audit-log filter template\" to see filter templates.")
	updateCmd.Flags().String(flag.DisplayName, "", "The display name of the filter rule.")
	updateCmd.Flags().Bool(flag.Enabled, false, "Enable or disable the filter rule.")

	return updateCmd
}
