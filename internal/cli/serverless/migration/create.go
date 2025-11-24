package migration

import (
	"context"
	"fmt"
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
	  $ %[1]s serverless migration create -c <cluster-id> --display-name <name> --sources '<sources-json>' --target '<target-json>' --mode <mode>

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
			sourcesStr, err := cmd.Flags().GetString(flag.MigrationSources)
			if err != nil {
				return errors.Trace(err)
			}
			targetStr, err := cmd.Flags().GetString(flag.MigrationTarget)
			if err != nil {
				return errors.Trace(err)
			}
			modeStr, err := cmd.Flags().GetString(flag.MigrationMode)
			if err != nil {
				return errors.Trace(err)
			}

			if strings.TrimSpace(name) == "" {
				return errors.New("display name is required")
			}
			sources, err := parseMigrationSources(sourcesStr)
			if err != nil {
				return err
			}
			target, err := parseMigrationTarget(targetStr)
			if err != nil {
				return err
			}
			mode, err := parseMigrationMode(modeStr)
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
	cmd.Flags().String(flag.MigrationSources, "", "Sources definition in JSON. Use \"ticloud serverless migration template --type sources\" for a template.")
	cmd.Flags().String(flag.MigrationTarget, "", "Target definition in JSON. Use \"ticloud serverless migration template --type target\" for a template.")
	cmd.Flags().String(flag.MigrationMode, "", fmt.Sprintf("Migration mode, one of %v.", taskModeValues()))
	cmd.Flags().Bool(flag.MigrationPrecheckOnly, false, "Run a migration precheck with the provided inputs and exit without creating a task.")

	return cmd
}

func markCreateMigrationRequiredFlags(cmd *cobra.Command) error {
	for _, fn := range []string{flag.ClusterID, flag.DisplayName, flag.MigrationSources, flag.MigrationTarget, flag.MigrationMode} {
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
	var lastStatus string
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			result, err := client.GetMigrationPrecheck(ctx, clusterID, precheckID)
			if err != nil {
				return errors.Trace(err)
			}
			status := strings.ToUpper(aws.ToString(result.Status))
			if status == "" {
				status = "PENDING"
			}
			if status != lastStatus {
				fmt.Fprintf(h.IOStreams.Out, "precheck %s status: %s\n", precheckID, status)
				lastStatus = status
			}
			if isPrecheckPending(status) {
				continue
			}
			if err := printPrecheckSummary(precheckID, status, result, h); err != nil {
				return err
			}
			if strings.EqualFold(status, "FAILED") || aws.ToInt32(result.FailedCnt) > 0 {
				return errors.New("migration precheck failed")
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("migration precheck %s passed", precheckID))
			return nil
		}
	}
}

func isPrecheckPending(status string) bool {
	switch status {
	case "PENDING", "RUNNING", "PROCESSING", "IN_PROGRESS", "":
		return true
	default:
		return false
	}
}

func printPrecheckSummary(id, status string, result *pkgmigration.MigrationPrecheck, h *internal.Helper) error {
	fmt.Fprintf(h.IOStreams.Out, "precheck %s finished with status %s\n", id, status)
	fmt.Fprintf(h.IOStreams.Out, "Total: %d, Success: %d, Warn: %d, Failed: %d\n",
		aws.ToInt32(result.Total), aws.ToInt32(result.SuccessCnt), aws.ToInt32(result.WarnCnt), aws.ToInt32(result.FailedCnt))
	if len(result.Items) == 0 {
		return nil
	}
	columns := []output.Column{"Type", "Status", "Description", "Reason", "Solution"}
	rows := make([]output.Row, 0, len(result.Items))
	for _, item := range result.Items {
		rows = append(rows, output.Row{
			precheckItemType(item.Type),
			aws.ToString(item.Status),
			aws.ToString(item.Description),
			aws.ToString(item.Reason),
			aws.ToString(item.Solution),
		})
	}
	return output.PrintHumanTable(h.IOStreams.Out, columns, rows)
}

func precheckItemType(value *pkgmigration.PrecheckItemType) string {
	if value == nil {
		return ""
	}
	return string(*value)
}
