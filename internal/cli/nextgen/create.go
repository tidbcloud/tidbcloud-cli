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
	"regexp"
	"strconv"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	encryptionNone       = "none"
	encryptionDefaultKey = "default-key"
	encryptionCMEK       = "cmek"
	cmekKeyARNFlag       = "cmek-key-arn"
	projectNameFlag      = "project-name"
	projectIDLabel       = "tidb.cloud/project"
)

var displayNamePattern = regexp.MustCompile(`^[A-Za-z0-9][-A-Za-z0-9]{2,62}[A-Za-z0-9]$`)

func createCmd(h *internal.Helper, plan planSpec) *cobra.Command {
	command := &cobra.Command{
		Use:         "create",
		Short:       fmt.Sprintf("Create a %s instance", plan.displayName),
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  $ %[1]s %[2]s --create --display-name <name> --region <region-id> --max-rcu <max-rcu>
  $ %[1]s %[2]s --create --project-id <project-id> --display-name <name> --region <region-id> --max-rcu <max-rcu>
  $ %[1]s %[2]s --create --display-name <name> --region <region-id> --max-rcu <max-rcu> --encryption default-key`, config.CliName, plan.commandName),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := validatedOutputFormat(cmd)
			if err != nil {
				return err
			}

			displayName, err := cmd.Flags().GetString(flag.DisplayName)
			if err != nil {
				return err
			}
			if err := validateDisplayName(displayName); err != nil {
				return err
			}
			regionID, err := cmd.Flags().GetString(flag.Region)
			if err != nil {
				return err
			}
			maxRCU, err := cmd.Flags().GetInt64(flag.MaxRCU)
			if err != nil {
				return err
			}
			if maxRCU <= 0 {
				return fmt.Errorf("--%s must be greater than zero", flag.MaxRCU)
			}
			projectID, err := cmd.Flags().GetString(flag.ProjectID)
			if err != nil {
				return err
			}
			if projectID != "" {
				id, err := strconv.ParseUint(projectID, 10, 64)
				if err != nil || id == 0 {
					return fmt.Errorf("--%s must be a positive decimal project ID (uint64), not a project name", flag.ProjectID)
				}
			}
			projectName, err := cmd.Flags().GetString(projectNameFlag)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed(projectNameFlag) && projectName == "" {
				return fmt.Errorf("--%s must not be empty", projectNameFlag)
			}
			encryption, err := cmd.Flags().GetString(flag.Encryption)
			if err != nil {
				return err
			}
			var keyARN string
			if plan.servicePlan == api.V1BETA1SERVICEPLAN_PREMIUM {
				keyARN, err = cmd.Flags().GetString(cmekKeyARNFlag)
				if err != nil {
					return err
				}
				if cmd.Flags().Changed(cmekKeyARNFlag) && encryption != encryptionCMEK {
					return fmt.Errorf("--%s requires --%s cmek", cmekKeyARNFlag, flag.Encryption)
				}
			}

			body := api.NewNextgenv1beta2Tidb(displayName, regionID, plan.servicePlan)
			body.SetMaxRcu(strconv.FormatInt(maxRCU, 10))
			switch encryption {
			case encryptionNone:
				body.DualLayerDataEncryption = &api.TidbDualLayerDataEncryption{NoEncryption: map[string]interface{}{}}
			case encryptionDefaultKey:
				body.DualLayerDataEncryption = &api.TidbDualLayerDataEncryption{DefaultKey: map[string]interface{}{}}
			case encryptionCMEK:
				if plan.servicePlan != api.V1BETA1SERVICEPLAN_PREMIUM {
					return fmt.Errorf("CMEK is only supported for Premium instances")
				}
				if keyARN == "" {
					return fmt.Errorf("--%s is required with --%s cmek", cmekKeyARNFlag, flag.Encryption)
				}
				if _, err := cmekProvider(regionID); err != nil {
					return err
				}
			default:
				if plan.servicePlan == api.V1BETA1SERVICEPLAN_PREMIUM {
					return fmt.Errorf("unsupported encryption %q; supported values are %q, %q and %q", encryption, encryptionNone, encryptionDefaultKey, encryptionCMEK)
				}
				return fmt.Errorf("unsupported encryption %q; supported values are %q and %q", encryption, encryptionNone, encryptionDefaultKey)
			}

			if projectName != "" {
				projectClient, err := h.Client()
				if err != nil {
					return err
				}
				_, projects, err := cloud.RetrieveProjects(cmd.Context(), h.QueryPageSize, projectClient)
				if err != nil {
					return fmt.Errorf("resolve project name %q: %w", projectName, err)
				}
				for _, project := range projects {
					if project.GetName() != projectName {
						continue
					}
					if projectID != "" {
						return fmt.Errorf("multiple accessible projects named %q; specify --%s instead", projectName, flag.ProjectID)
					}
					projectID = project.GetId()
				}
				if projectID == "" {
					return fmt.Errorf("no accessible project named %q", projectName)
				}
			}
			if projectID != "" {
				body.Labels = &map[string]string{projectIDLabel: projectID}
				cmd.Annotations[telemetry.ProjectID] = projectID
			}
			client, err := h.NextGenClient()
			if err != nil {
				return err
			}
			if encryption == encryptionCMEK {
				key, err := verifiedCMEK(cmd.Context(), client, regionID, keyARN)
				if err != nil {
					return err
				}
				body.DualLayerDataEncryption = &api.TidbDualLayerDataEncryption{Cmek: key}
			}
			instance, err := client.CreateTiDB(cmd.Context(), body)
			if err != nil {
				return err
			}
			if err := plan.validate(instance); err != nil {
				return err
			}

			if format == output.JsonFormat {
				return output.PrintJson(h.IOStreams.Out, instance)
			}
			fmt.Fprintf(h.IOStreams.Out, "%s\n", color.GreenString("Create request accepted for instance %s (state: %s)", instance.GetTidbId(), instance.GetState()))
			return nil
		},
	}
	command.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The display name of the instance.")
	command.Flags().StringP(flag.Region, flag.RegionShort, "", "The region ID, for example aws-us-west-2.")
	command.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The positive decimal project ID (not the project name) in which the instance will be created. If neither --project-id nor --project-name is set, the default TiDB X project is used.")
	command.Flags().String(projectNameFlag, "", "The exact, case-sensitive name of an accessible project. Must be unique; mutually exclusive with --project-id. Uses the same IAM endpoint as project list.")
	command.MarkFlagsMutuallyExclusive(flag.ProjectID, projectNameFlag)
	command.Flags().Int64(flag.MaxRCU, 0, "The maximum number of Request Capacity Units (RCUs). Must be at least the service-calculated minimum; the default calculation is max(5000, max-rcu / 4), subject to organization and plan configuration.")
	command.Flags().String(flag.Encryption, encryptionNone, "Dual-layer data encryption, one of [\"none\" \"default-key\"]. default-key uses the TiDB Cloud default key; no key value is required.")
	if plan.servicePlan == api.V1BETA1SERVICEPLAN_PREMIUM {
		command.Flags().Lookup(flag.Encryption).Usage = "Dual-layer data encryption, one of [\"none\" \"default-key\" \"cmek\"]. default-key uses the TiDB Cloud default key; no key value is required."
		command.Flags().String(cmekKeyARNFlag, "", "The AWS or Alibaba Cloud KMS key ARN. Required with --encryption cmek.")
		command.Example += fmt.Sprintf("\n  $ %s premium --create --display-name <name> --region <region-id> --max-rcu <max-rcu> --encryption cmek --cmek-key-arn <key-arn>", config.CliName)
	}
	command.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	_ = command.MarkFlagRequired(flag.DisplayName)
	_ = command.MarkFlagRequired(flag.Region)
	_ = command.MarkFlagRequired(flag.MaxRCU)
	return command
}

func validateDisplayName(displayName string) error {
	if !displayNamePattern.MatchString(displayName) {
		return fmt.Errorf("display name must be 4-64 characters, start and end with an alphanumeric character, and contain only alphanumeric characters or hyphens")
	}
	return nil
}
