package migration

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
)

const (
	migrationSourcesTemplateWithExplain = `[
    {
        // Required: source database type. Supported: SOURCE_TYPE_MYSQL, SOURCE_TYPE_ALICLOUD_RDS_MYSQL
        "sourceType": "SOURCE_TYPE_MYSQL",
        "connProfile": {
            // Optional connection type, PUBLIC or PRIVATE_LINK
            "connType": "PUBLIC",
            "host": "10.0.0.2",
            "port": 3306,
            "user": "dm_sync_user",
            "password": "Passw0rd!",
            // Optional when using private link
            "endpointId": "pl-xxxxxxxx",
            "security": {
                // Optional TLS materials encoded in Base64
                "sslCaContent": "<base64-of-ca.pem>",
                "sslCertContent": "<base64-of-client-cert.pem>",
                "sslKeyContent": "<base64-of-client-key.pem>",
                "certAllowedCn": ["client-cn"]
            }
        },
        // Optional block/allow rules to whitelist schemas/tables
        "baRules": {
            "doDbs": ["app_db"],
            "doTables": [
                {"schema": "app_db", "table": "orders"},
                {"schema": "app_db", "table": "customers"}
            ]
        },
        // Optional route rules for renaming schemas/tables
        "routeRules": [
            {
                "sourceTable": {
                    "schemaPattern": "app_db",
                    "tablePattern": "orders"
                },
                "targetTable": {
                    "schema": "app_db",
                    "table": "orders_copy"
                }
            }
        ],
        // Optional start position for incremental sync. Provide binlogName+binlogPos or binlogGtid
        "binlogName": "mysql-bin.000001",
        "binlogPos": 4,
        "binlogGtid": "3E11FA47-71CA-11E1-9E33-C80AA9429562:1-12345"
    }
 ]`

	migrationSourcesTemplate = `[
    {
        "sourceType": "SOURCE_TYPE_MYSQL",
        "connProfile": {
            "connType": "PUBLIC",
            "host": "10.0.0.2",
            "port": 3306,
            "user": "dm_sync_user",
            "password": "Passw0rd!"
        },
        "baRules": {
            "doDbs": ["app_db"],
            "doTables": [{"schema": "app_db", "table": "orders"}]
        },
        "routeRules": [
            {
                "sourceTable": {"schemaPattern": "app_db", "tablePattern": "orders"},
                "targetTable": {"schema": "app_db", "table": "orders_copy"}
            }
        ]
    }
 ]`

	migrationTargetTemplateWithExplain = `{
    // Target TiDB Cloud user used by the migration task
    "user": "migration_user",
    // Password corresponding to the target user
    "password": "Passw0rd!"
}`

	migrationTargetTemplate = `{
    "user": "migration_user",
    "password": "Passw0rd!"
}`
)

func TemplateCmd(h *internal.Helper) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "template",
		Short: "Show migration JSON templates",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Show all migration templates:
  $ %[1]s serverless migration template

  Show the sources template with explanations:
  $ %[1]s serverless migration template --type sources --explain`, config.CliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			explain, err := cmd.Flags().GetBool(flag.Explain)
			if err != nil {
				return err
			}
			templateType, err := cmd.Flags().GetString(flag.MigrationTemplateType)
			if err != nil {
				return err
			}

			return renderMigrationTemplate(h, strings.ToLower(templateType), explain)
		},
	}

	cmd.Flags().Bool(flag.Explain, false, "Show template with inline explanations.")
	cmd.Flags().String(flag.MigrationTemplateType, "", "Template type to show, one of [\"sources\", \"target\"]. Default prints both.")
	return cmd
}

func renderMigrationTemplate(h *internal.Helper, templateType string, explain bool) error {
	switch templateType {
	case "sources":
		if explain {
			fmt.Fprintln(h.IOStreams.Out, migrationSourcesTemplateWithExplain)
		} else {
			fmt.Fprintln(h.IOStreams.Out, migrationSourcesTemplate)
		}
	case "target":
		if explain {
			fmt.Fprintln(h.IOStreams.Out, migrationTargetTemplateWithExplain)
		} else {
			fmt.Fprintln(h.IOStreams.Out, migrationTargetTemplate)
		}
	case "":
		fmt.Fprintln(h.IOStreams.Out, color.GreenString("Sources template:"))
		if explain {
			fmt.Fprintln(h.IOStreams.Out, migrationSourcesTemplateWithExplain)
		} else {
			fmt.Fprintln(h.IOStreams.Out, migrationSourcesTemplate)
		}
		fmt.Fprintln(h.IOStreams.Out, color.GreenString("Target template:"))
		if explain {
			fmt.Fprintln(h.IOStreams.Out, migrationTargetTemplateWithExplain)
		} else {
			fmt.Fprintln(h.IOStreams.Out, migrationTargetTemplate)
		}
	default:
		return fmt.Errorf("unknown template type %q", templateType)
	}
	return nil
}
