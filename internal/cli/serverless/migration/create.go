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

package migration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	pkgmigration "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/migration"
)

func CreateCmd(h *internal.Helper) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a migration",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a migration:
	  $ %[1]s serverless migration create -c <cluster-id> --display-name <name> --config-file /path/to/config.json

	  Run migration precheck only with shared inputs:
  $ %[1]s serverless migration create --precheck-only`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return markCreateMigrationRequiredFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			precheckOnly, err := cmd.Flags().GetBool(flag.MigrationPrecheckOnly)
			if err != nil {
				return errors.Trace(err)
			}
			clusterID, err := cmd.Flags().GetString(flag.ClusterID)
			if err != nil {
				return errors.Trace(err)
			}
			name, err := cmd.Flags().GetString(flag.DisplayName)
			if err != nil {
				return errors.Trace(err)
			}
			configPath, err := cmd.Flags().GetString(flag.MigrationConfigFile)
			if err != nil {
				return errors.Trace(err)
			}
			configPath = strings.TrimSpace(configPath)
			if configPath == "" {
				return errors.New("config file path is required")
			}
			definitionBytes, err := os.ReadFile(configPath)
			if err != nil {
				return errors.Annotatef(err, "failed to read config file %q", configPath)
			}
			definitionStr := string(definitionBytes)

			if strings.TrimSpace(name) == "" {
				return errors.New("display name is required")
			}
			sources, target, mode, err := parseMigrationDefinition(definitionStr)
			if err != nil {
				return err
			}

			createBody := &pkgmigration.MigrationServiceCreateMigrationBody{
				DisplayName: name,
				Sources:     sources,
				Target:      target,
				Mode:        mode,
			}
			precheckBody := &pkgmigration.MigrationServicePrecheckBody{
				DisplayName: name,
				Sources:     sources,
				Target:      target,
				Mode:        mode,
			}

			if precheckOnly {
				return runMigrationPrecheck(ctx, d, clusterID, precheckBody, h)
			}

			resp, err := d.CreateMigrationTask(ctx, clusterID, createBody)
			if err != nil {
				return errors.Trace(err)
			}

			taskID := aws.ToString(resp.MigrationId)
			if taskID == "" {
				taskID = aws.ToString(resp.DisplayName)
			}
			if taskID == "" {
				taskID = "<unknown>"
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("migration %s created", taskID))
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the target cluster.")
	cmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "Display name for the migration.")
	cmd.Flags().String(flag.MigrationConfigFile, "", "Path to a migration config JSON file. Use \"ticloud serverless migration template --modetype <mode>\" to print templates.")
	cmd.Flags().Bool(flag.MigrationPrecheckOnly, false, "Run a migration precheck with the provided inputs and exit without creating a task.")

	return cmd
}

func markCreateMigrationRequiredFlags(cmd *cobra.Command) error {
	for _, fn := range []string{flag.ClusterID, flag.DisplayName, flag.MigrationConfigFile} {
		if err := cmd.MarkFlagRequired(fn); err != nil {
			return err
		}
	}
	return nil
}

const precheckPollInterval = 5 * time.Second

func runMigrationPrecheck(ctx context.Context, client cloud.TiDBCloudClient, clusterID string, body *pkgmigration.MigrationServicePrecheckBody, h *internal.Helper) error {
	resp, err := client.CreateMigrationPrecheck(ctx, clusterID, body)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.PrecheckId == nil || *resp.PrecheckId == "" {
		return errors.New("precheck created but ID is empty")
	}
	precheckID := *resp.PrecheckId
	fmt.Fprintf(h.IOStreams.Out, "migration precheck %s created, polling results...\n", precheckID)

	ticker := time.NewTicker(precheckPollInterval)
	defer ticker.Stop()
	var lastStatus pkgmigration.MigrationPrecheckStatus
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			result, err := client.GetMigrationPrecheck(ctx, clusterID, precheckID)
			if err != nil {
				return errors.Trace(err)
			}
			status := precheckStatusOrDefault(result.Status)
			if status != lastStatus {
				fmt.Fprintf(h.IOStreams.Out, "precheck %s status: %s\n", precheckID, status)
				lastStatus = status
			}
			if isPrecheckUnfinished(status) {
				continue
			}
			if err := printPrecheckSummary(precheckID, status, result, h); err != nil {
				return err
			}
			if status == pkgmigration.MIGRATIONPRECHECKSTATUS_FAILED {
				fmt.Fprintln(h.IOStreams.Out, color.RedString("migration precheck %s failed", precheckID))
				return errors.New("migration precheck failed")
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("migration precheck %s passed", precheckID))
			return nil
		}
	}
}

func precheckStatusOrDefault(value *pkgmigration.MigrationPrecheckStatus) pkgmigration.MigrationPrecheckStatus {
	if value == nil || *value == "" {
		return pkgmigration.MIGRATIONPRECHECKSTATUS_PENDING
	}
	return *value
}

func isPrecheckUnfinished(status pkgmigration.MigrationPrecheckStatus) bool {
	switch status {
	case pkgmigration.MIGRATIONPRECHECKSTATUS_PENDING,
		pkgmigration.MIGRATIONPRECHECKSTATUS_RUNNING:
		return true
	default:
		return false
	}
}

func printPrecheckSummary(id string, status pkgmigration.MigrationPrecheckStatus, result *pkgmigration.MigrationPrecheck, h *internal.Helper) error {
	fmt.Fprintf(h.IOStreams.Out, "precheck %s finished with status %s\n", id, status)
	fmt.Fprintf(h.IOStreams.Out, "Total: %d, Success: %d, Warn: %d, Failed: %d\n",
		aws.ToInt32(result.Total), aws.ToInt32(result.SuccessCnt), aws.ToInt32(result.WarnCnt), aws.ToInt32(result.FailedCnt))
	if len(result.Items) == 0 {
		return nil
	}
	columns := []output.Column{"Type", "Status", "Description", "Reason", "Solution"}
	rows := make([]output.Row, 0, len(result.Items))
	for _, item := range result.Items {
		if !shouldPrintPrecheckItem(item.Status) {
			continue
		}
		var status string
		if item.Status != nil {
			status = string(*item.Status)
		}
		rows = append(rows, output.Row{
			precheckItemType(item.Type),
			status,
			aws.ToString(item.Description),
			aws.ToString(item.Reason),
			aws.ToString(item.Solution),
		})
	}
	if len(rows) == 0 {
		fmt.Fprintln(h.IOStreams.Out, "No warning or failure details returned.")
		return nil
	}
	return output.PrintHumanTable(h.IOStreams.Out, columns, rows)
}

func precheckItemType(value *pkgmigration.PrecheckItemType) string {
	if value == nil {
		return ""
	}
	return string(*value)
}

// shouldPrintPrecheckItem reports whether a precheck item should be shown to users.
// Currently only WARNING and FAILED statuses surface because SUCCESS does not
// provide actionable information.
func shouldPrintPrecheckItem(status *pkgmigration.PrecheckItemStatus) bool {
	if status == nil {
		return false
	}
	switch *status {
	case pkgmigration.PRECHECKITEMSTATUS_WARNING,
		pkgmigration.PRECHECKITEMSTATUS_FAILED:
		return true
	default:
		return false
	}
}
