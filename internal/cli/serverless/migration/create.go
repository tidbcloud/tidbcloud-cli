// Copyright 2026 PingCAP, Inc.
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
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tailscale/hujson"

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
  $ %[1]s serverless migration create -c <cluster-id> --display-name <name> --config-file <file-path> --dry-run
  $ %[1]s serverless migration create -c <cluster-id> --display-name <name> --config-file <file-path>
`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return markCreateMigrationRequiredFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			dryRun, err := cmd.Flags().GetBool(flag.DryRun)
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
			if strings.TrimSpace(name) == "" {
				return errors.New("display name is required")
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

			sources, target, mode, importMode, err := parseMigrationDefinition(definitionStr)
			if err != nil {
				return err
			}

			if dryRun {
				precheckBody := &pkgmigration.MigrationServicePrecheckBody{
					DisplayName: name,
					Sources:     sources,
					Target:      target,
					Mode:        mode,
					ImportMode:  importMode,
				}
				return runMigrationPrecheck(ctx, d, clusterID, precheckBody, h)
			}

			createBody := &pkgmigration.MigrationServiceCreateMigrationBody{
				DisplayName: name,
				Sources:     sources,
				Target:      target,
				Mode:        mode,
				ImportMode:  importMode,
			}

			resp, err := d.CreateMigration(ctx, clusterID, createBody)
			if err != nil {
				return errors.Trace(err)
			}

			migrationID := aws.ToString(resp.MigrationId)
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("migration %s(%s) created", name, migrationID))
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the target cluster.")
	cmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "Display name for the migration.")
	cmd.Flags().String(flag.MigrationConfigFile, "", "Path to a migration config JSON file. Use \"ticloud serverless migration template --mode <mode>\" to print templates.")
	cmd.Flags().Bool(flag.DryRun, false, "Run a migration precheck (dry run) with the provided inputs without creating a migration.")

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

const (
	precheckPollInterval = 5 * time.Second
	precheckPollTimeout  = 2 * time.Minute
)

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
	pollCtx, cancel := context.WithTimeout(ctx, precheckPollTimeout)
	defer cancel()

	// Poll precheck status until it finishes or the overall timeout is hit.
	for {
		select {
		case <-pollCtx.Done():
			if pollCtx.Err() == context.DeadlineExceeded {
				return errors.Errorf("migration precheck polling timed out after %s", precheckPollTimeout)
			}
			return pollCtx.Err()
		case <-ticker.C:
			result, err := client.GetMigrationPrecheck(pollCtx, clusterID, precheckID)
			if err != nil {
				return errors.Trace(err)
			}
			finished, err := printPrecheckSummary(result, h)
			if err != nil {
				return err
			}
			if !finished {
				continue
			}
			if result.GetStatus() == pkgmigration.MIGRATIONPRECHECKSTATUS_FAILED {
				fmt.Fprintln(h.IOStreams.Out, color.RedString("migration precheck %s failed", precheckID))
				return errors.New("migration precheck failed")
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("migration precheck %s passed", precheckID))
			return nil
		}
	}
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

func printPrecheckSummary(result *pkgmigration.MigrationPrecheck, h *internal.Helper) (bool, error) {
	if isPrecheckUnfinished(result.GetStatus()) {
		fmt.Fprintf(h.IOStreams.Out, "precheck %s summary (status %s)\n", result.GetPrecheckId(), result.GetStatus())
		fmt.Fprintf(h.IOStreams.Out, "Total: %d, Success: %d, Warn: %d, Failed: %d\n",
			aws.ToInt32(result.Total), aws.ToInt32(result.SuccessCnt), aws.ToInt32(result.WarnCnt), aws.ToInt32(result.FailedCnt))
		return false, nil
	}

	fmt.Fprintf(h.IOStreams.Out, "precheck %s finished with status %s\n", result.GetPrecheckId(), result.GetStatus())
	fmt.Fprintf(h.IOStreams.Out, "Total: %d, Success: %d, Warn: %d, Failed: %d\n",
		aws.ToInt32(result.Total), aws.ToInt32(result.SuccessCnt), aws.ToInt32(result.WarnCnt), aws.ToInt32(result.FailedCnt))
	if len(result.Items) == 0 {
		return true, nil
	}
	columns := []output.Column{"Type", "Status", "Description", "Reason", "Solution"}
	rows := make([]output.Row, 0, len(result.Items))
	for _, item := range result.Items {
		if !shouldPrintPrecheckItem(item.Status) {
			continue
		}
		rows = append(rows, output.Row{
			string(pointer.Get(item.Type)),
			string(pointer.Get(item.Status)),
			pointer.Get(item.Description),
			pointer.Get(item.Reason),
			pointer.Get(item.Solution),
		})
	}
	if len(rows) == 0 {
		return true, nil
	}
	return true, output.PrintHumanTable(h.IOStreams.Out, columns, rows)
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

func parseMigrationDefinition(value string) ([]pkgmigration.Source, pkgmigration.Target, pkgmigration.TaskMode, *pkgmigration.ImportMode, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, pkgmigration.Target{}, "", nil, errors.New("migration config is required; use --config-file")
	}
	var payload struct {
		Sources    []pkgmigration.Source `json:"sources"`
		Target     *pkgmigration.Target  `json:"target"`
		Mode       string                `json:"mode"`
		ImportMode *string               `json:"importMode"`
	}
	stdJson, err := standardizeJSON([]byte(trimmed))
	if err != nil {
		return nil, pkgmigration.Target{}, "", nil, errors.Annotate(err, "invalid migration definition JSON")
	}
	if err := json.Unmarshal(stdJson, &payload); err != nil {
		return nil, pkgmigration.Target{}, "", nil, errors.Annotate(err, "invalid migration definition JSON")
	}
	if len(payload.Sources) == 0 {
		return nil, pkgmigration.Target{}, "", nil, errors.New("migration definition must include at least one source")
	}
	if payload.Target == nil {
		return nil, pkgmigration.Target{}, "", nil, errors.New("migration definition must include the target block")
	}
	mode, err := parseMigrationMode(payload.Mode)
	if err != nil {
		return nil, pkgmigration.Target{}, "", nil, err
	}

	importMode, err := parseImportMode(payload.ImportMode)
	if err != nil {
		return nil, pkgmigration.Target{}, "", nil, err
	}
	if mode == pkgmigration.TASKMODE_INCREMENTAL && importMode != nil {
		return nil, pkgmigration.Target{}, "", nil, errors.New("importMode is only applicable for mode=ALL; remove importMode or switch to mode=ALL")
	}

	return payload.Sources, *payload.Target, mode, importMode, nil
}

func parseMigrationMode(value string) (pkgmigration.TaskMode, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", errors.New("empty config file")
	}
	normalized := strings.ToUpper(trimmed)
	mode := pkgmigration.TaskMode(normalized)
	if slices.Contains(pkgmigration.AllowedTaskModeEnumValues, mode) {
		return mode, nil
	}
	return "", errors.Errorf("invalid mode %q, allowed values: %s", value, pkgmigration.AllowedTaskModeEnumValues)
}

func parseImportMode(raw *string) (*pkgmigration.ImportMode, error) {
	if raw == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil, nil
	}

	normalized := strings.ToUpper(trimmed)
	switch normalized {
	case "LOGICAL":
		normalized = "IMPORT_MODE_LOGICAL"
	case "PHYSICAL":
		normalized = "IMPORT_MODE_PHYSICAL"
	}

	mode := pkgmigration.ImportMode(normalized)
	if slices.Contains(pkgmigration.AllowedImportModeEnumValues, mode) {
		return &mode, nil
	}
	return nil, errors.Errorf("invalid importMode %q, allowed values: %s", trimmed, pkgmigration.AllowedImportModeEnumValues)
}

// standardizeJSON accepts JSON With Commas and Comments(JWCC) see
// https://nigeltao.github.io/blog/2021/json-with-commas-comments.html) and
// returns a standard JSON byte slice ready for json.Unmarshal.
func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}
