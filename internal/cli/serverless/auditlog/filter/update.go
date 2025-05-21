package filter

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	alutil "github.com/tidbcloud/tidbcloud-cli/internal/cli/serverless/auditlog/util"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
)

type UpdateFilterRuleOpts struct {
	interactive bool
}

func (o UpdateFilterRuleOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.AuditLogFilterRuleName,
		flag.AuditLogFilterRuleUsers,
		flag.AuditLogFilterRuleFilters,
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
		cmd.MarkFlagsOneRequired(flag.AuditLogFilterRuleUsers, flag.AuditLogFilterRuleFilters)
	}
	return nil
}

var mutableFilterRuleFields = []string{
	"users",
	"filters",
}

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateFilterRuleOpts{interactive: true}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an audit log filter rule",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Update an audit log filter rule in interactive mode:
  $ %[1]s serverless auditlog filter-rule update

  Update users of an audit log filter rule in non-interactive mode:
  $ %[1]s serverless auditlog filter-rule update --cluster-id <cluster-id> --rule-name <rule-name> --users user1,user2

  Update filters of an audit log filter rule in non-interactive mode:
  $ %[1]s serverless auditlog filter-rule update --cluster-id <cluster-id> --rule-name <rule-name> --filters '{"classes": ["QUERY", "EXECUTE"], "tables": ["test.t1"]}' --filters '{"classes": ["QUERY"]}'
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
			var users []string
			var filters []auditlog.AuditLogFilter

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
				case "users":
					inputs := []string{flag.AuditLogFilterRuleUsers}
					textInput, err := ui.InitialInputModel(inputs, alutil.InputDescription)
					if err != nil {
						return err
					}
					usersStr := textInput.Inputs[0].Value()
					if err := json.Unmarshal([]byte(usersStr), &users); err != nil {
						return errors.New(fmt.Sprintf("invalid users format: %s, please use JSON format", usersStr))
					}
					if len(users) == 0 {
						return errors.New("empty users")
					}
				case "filters":
					inputs := []string{flag.AuditLogFilterRuleFilters}
					textInput, err := ui.InitialInputModel(inputs, alutil.InputDescription)
					if err != nil {
						return err
					}
					filtersStr := textInput.Inputs[0].Value()
					if err := json.Unmarshal([]byte(filtersStr), &filters); err != nil {
						return errors.New(fmt.Sprintf("invalid filters format: %s, please use JSON format", filtersStr))
					}
					if len(filters) == 0 {
						return errors.New("empty filters")
					}
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
				if cmd.Flags().Changed(flag.AuditLogFilterRuleUsers) {
					users, err = cmd.Flags().GetStringSlice(flag.AuditLogFilterRuleUsers)
					if err != nil {
						return errors.Trace(err)
					}
				}
				if cmd.Flags().Changed(flag.AuditLogFilterRuleFilters) {
					filtersStr, err := cmd.Flags().GetStringArray(flag.AuditLogFilterRuleFilters)
					if err != nil {
						return errors.Trace(err)
					}
					for _, f := range filtersStr {
						var filter auditlog.AuditLogFilter
						if err := json.Unmarshal([]byte(f), &filter); err != nil {
							return errors.New(fmt.Sprintf("invalid filter format: %s, please use JSON format", f))
						}
						filters = append(filters, filter)
					}
				}
			}
			if len(users) == 0 && len(filters) == 0 {
				return errors.New("nothing to update")
			}
			body := &auditlog.AuditLogServiceUpdateAuditLogFilterRuleBody{}
			if len(users) > 0 {
				body.Users = users
			}
			if len(filters) > 0 {
				body.Filters = filters
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
	updateCmd.Flags().StringSlice(flag.AuditLogFilterRuleUsers, nil, "Users to apply the rule to. e.g. %@%.")
	updateCmd.Flags().StringArray(flag.AuditLogFilterRuleFilters, nil, "Filter expressions. e.g. '{\"classes\": [\"QUERY\"]' or '{}' to filter all audit logs.")

	return updateCmd
}
