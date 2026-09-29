// Copyright 2026 PingCAP, Inc.
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

package nextgen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/spf13/cobra"
)

type planSpec struct {
	commandName string
	aliases     []string
	displayName string
	servicePlan api.V1beta1ServicePlan
	queryPlan   api.TidbServiceListTidbsServicePlanParameter
}

var premiumPlan = planSpec{
	commandName: "premium",
	displayName: "TiDB Cloud Premium",
	servicePlan: api.V1BETA1SERVICEPLAN_PREMIUM,
	queryPlan:   api.TIDBSERVICELISTTIDBSSERVICEPLANPARAMETER_PREMIUM,
}

var essentialV2Plan = planSpec{
	commandName: "essential-v2",
	aliases:     []string{"essential"},
	displayName: "TiDB Cloud Essential V2",
	servicePlan: api.V1BETA1SERVICEPLAN_ESSENTIAL_V2,
	queryPlan:   api.TIDBSERVICELISTTIDBSSERVICEPLANPARAMETER_ESSENTIAL_V2,
}

func PremiumCmd(h *internal.Helper) *cobra.Command {
	return cmd(h, premiumPlan)
}

func EssentialCmd(h *internal.Helper) *cobra.Command {
	return cmd(h, essentialV2Plan)
}

func cmd(h *internal.Helper, plan planSpec) *cobra.Command {
	root := &cobra.Command{
		Use:     plan.commandName,
		Aliases: plan.aliases,
		Short:   fmt.Sprintf("Manage %s instances", plan.displayName),
		Example: fmt.Sprintf(`  $ %[1]s %[2]s --create --display-name <name> --region <region-id> --max-rcu <max-rcu>
  $ %[1]s %[2]s --create --project-id <project-id> --display-name <name> --region <region-id> --max-rcu <max-rcu>
  $ %[1]s %[2]s --list
  $ %[1]s %[2]s --describe -c <instance-id>
  $ %[1]s %[2]s --update -c <instance-id> --max-rcu <max-rcu>
  $ %[1]s %[2]s --delete -c <instance-id>
  $ %[1]s %[2]s public-endpoint enable -c <instance-id>
  $ %[1]s %[2]s public-endpoint disable -c <instance-id>
  $ %[1]s %[2]s password -c <instance-id>`, config.CliName, plan.commandName),
	}
	actionCommands := []*cobra.Command{
		createCmd(h, plan),
		listCmd(h, plan),
		describeCmd(h, plan),
		updateCmd(h, plan),
		deleteCmd(h, plan),
	}
	for _, command := range actionCommands {
		configureActionCommand(command, plan)
		root.AddCommand(command)
	}
	root.Long = actionHelp(plan, actionCommands)
	root.AddCommand(shellCmd(h, plan))
	root.AddCommand(passwordCmd(h, plan))
	root.AddCommand(publicEndpointCmd(h, plan))
	root.AddCommand(regionCmd(h, plan))
	if plan.servicePlan == api.V1BETA1SERVICEPLAN_PREMIUM {
		root.AddCommand(cmekCmd(h))
	}
	root.Flags().Bool("create", false, fmt.Sprintf("Create a %s instance.", plan.displayName))
	root.Flags().Bool("list", false, fmt.Sprintf("List %s instances.", plan.displayName))
	root.Flags().Bool("describe", false, fmt.Sprintf("Describe a %s instance.", plan.displayName))
	root.Flags().Bool("update", false, fmt.Sprintf("Update a %s instance.", plan.displayName))
	root.Flags().Bool("delete", false, fmt.Sprintf("Delete a %s instance.", plan.displayName))
	return root
}

func configureActionCommand(command *cobra.Command, plan planSpec) {
	command.Hidden = true
	usage := fmt.Sprintf("%s %s --%s [flags]", config.CliName, plan.commandName, command.Name())
	command.SetUsageTemplate(strings.Replace(command.UsageTemplate(), "{{.UseLine}}", usage, 1))
}

func actionHelp(plan planSpec, commands []*cobra.Command) string {
	var result strings.Builder
	fmt.Fprintf(&result, "Manage %s instances.\n\nAction-specific flags:", plan.displayName)
	for _, command := range commands {
		flagUsages := strings.TrimSuffix(command.Flags().FlagUsages(), "\n")
		flagUsages = "    " + strings.ReplaceAll(flagUsages, "\n", "\n    ")
		fmt.Fprintf(&result, "\n\n--%s\n\n%s", command.Name(), flagUsages)
	}
	return strings.TrimSpace(result.String())
}

// NormalizeActionArgs preserves the action-flag command form in the release
// contract while keeping the Cobra subcommands as the canonical implementation.
func NormalizeActionArgs(args []string) ([]string, error) {
	if len(args) > 0 && (args[0] == cobra.ShellCompRequestCmd || args[0] == cobra.ShellCompNoDescRequestCmd) {
		normalized, err := normalizeActionArgs(args[1:])
		if err != nil {
			return nil, err
		}
		return append([]string{args[0]}, normalized...), nil
	}
	return normalizeActionArgs(args)
}

func normalizeActionArgs(args []string) ([]string, error) {
	normalized := append([]string(nil), args...)
	planIndex := rootCommandIndex(normalized)
	if planIndex < 0 || !isPlanCommand(normalized[planIndex]) {
		return normalized, nil
	}

	if hasCanonicalSubcommand(normalized[planIndex+1:]) {
		return normalized, nil
	}

	actionIndex := -1
	action := ""
	for index := planIndex + 1; index < len(normalized); index++ {
		if normalized[index] == "--" {
			break
		}
		candidate, enabled, matched, err := actionCommand(normalized[index])
		if err != nil {
			return nil, err
		}
		if !matched {
			continue
		}
		if !enabled {
			return nil, fmt.Errorf("--%s=false does not select an action", candidate)
		}
		if action != "" {
			return nil, fmt.Errorf("only one action flag may be specified: --%s and --%s", action, candidate)
		}
		action = candidate
		actionIndex = index
	}
	if action == "" {
		return normalized, nil
	}

	result := make([]string, 0, len(normalized))
	result = append(result, normalized[:planIndex+1]...)
	result = append(result, action)
	result = append(result, normalized[planIndex+1:actionIndex]...)
	result = append(result, normalized[actionIndex+1:]...)
	return result, nil
}

func rootCommandIndex(args []string) int {
	for index := 0; index < len(args); index++ {
		switch {
		case args[index] == "--profile" || args[index] == "-P":
			index++
		case strings.HasPrefix(args[index], "--profile=") || (strings.HasPrefix(args[index], "-P") && len(args[index]) > 2):
			continue
		case isRootBooleanFlag(args[index]):
			continue
		case strings.HasPrefix(args[index], "-"):
			return -1
		default:
			return index
		}
	}
	return -1
}

func isPlanCommand(arg string) bool {
	return arg == premiumPlan.commandName || arg == essentialV2Plan.commandName || arg == essentialV2Plan.aliases[0]
}

func hasCanonicalSubcommand(args []string) bool {
	for index := 0; index < len(args); index++ {
		switch {
		case args[index] == "--profile" || args[index] == "-P":
			index++
		case strings.HasPrefix(args[index], "--profile=") || (strings.HasPrefix(args[index], "-P") && len(args[index]) > 2):
			continue
		case isRootBooleanFlag(args[index]):
			continue
		default:
			switch args[index] {
			case "create", "list", "describe", "update", "delete", "password", "public-endpoint", "shell", "region", "cmek":
				return true
			}
			return false
		}
	}
	return false
}

func isRootBooleanFlag(arg string) bool {
	return arg == "--debug" || arg == "-D" || arg == "--no-color" ||
		strings.HasPrefix(arg, "--debug=") || strings.HasPrefix(arg, "-D=") || strings.HasPrefix(arg, "--no-color=")
}

func actionCommand(arg string) (action string, enabled bool, matched bool, err error) {
	name, value, hasValue := strings.Cut(arg, "=")
	switch name {
	case "--create", "--list", "--describe", "--update", "--delete":
		action = strings.TrimPrefix(name, "--")
	default:
		return "", false, false, nil
	}
	if !hasValue {
		return action, true, true, nil
	}
	enabled, err = strconv.ParseBool(value)
	if err != nil {
		return "", false, true, fmt.Errorf("invalid value %q for --%s: expected a boolean", value, action)
	}
	return action, enabled, true, nil
}

func validatedOutputFormat(cmd *cobra.Command) (string, error) {
	format, err := cmd.Flags().GetString(flag.Output)
	if err != nil {
		return "", err
	}
	if format != output.JsonFormat && format != output.HumanFormat {
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
	return format, nil
}

func (p planSpec) validate(instance *api.Nextgenv1beta2Tidb) error {
	if instance == nil {
		return fmt.Errorf("the API returned an empty instance")
	}
	if instance.ServicePlan != p.servicePlan {
		return fmt.Errorf("instance %q uses service plan %q; the %s command only manages %q instances", instance.GetTidbId(), instance.ServicePlan, p.commandName, p.servicePlan)
	}
	return nil
}
