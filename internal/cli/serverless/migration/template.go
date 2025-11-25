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
	migrationDefinitionAllTemplate = `{
    // Required migration mode. Use "ALL" for full + incremental.
    "mode": "ALL",
    // Target TiDB Cloud user credentials used by the migration task
    "target": {
        "user": "migration_user",
        "password": "Passw0rd!"
    },
    // List at least one migration source
    "sources": [
        {
            // Required: source database type
            "sourceType": "SOURCE_TYPE_MYSQL",
            "connProfile": {
                // Optional connection type, PUBLIC or PRIVATE_LINK
                "connType": "PUBLIC",
                "host": "10.0.0.2",
                "port": 3306,
                "user": "dm_sync_user",
                "password": "Passw0rd!",
                // Optional fields below are needed only for private link or TLS
                "endpointId": "pl-xxxxxxxx",
                "security": {
                    // TLS materials must be Base64 encoded
                    "sslCaContent": "<base64-of-ca.pem>",
                    "sslCertContent": "<base64-of-client-cert.pem>",
                    "sslKeyContent": "<base64-of-client-key.pem>",
                    "certAllowedCn": ["client-cn"]
                }
            },
            // Optional block/allow rules to control synced schemas/tables
            "baRules": {
                "doDbs": ["app_db"],
                "doTables": [
                    {"schema": "app_db", "table": "orders"},
                    {"schema": "app_db", "table": "customers"}
                ]
            },
            // Optional route rules to rename objects during migration
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
            // Optional start position for incremental sync (binlog position or GTID)
            "binlogName": "mysql-bin.000001",
            "binlogPos": 4,
            "binlogGtid": "3E11FA47-71CA-11E1-9E33-C80AA9429562:1-12345"
        }
    ]
}`

	migrationDefinitionIncrementalTemplate = `{
    // Incremental-only mode keeps the source and target in sync
    "mode": "INCREMENTAL",
    "target": {
        "user": "migration_user",
        "password": "Passw0rd!"
    },
    "sources": [
        {
            "sourceType": "SOURCE_TYPE_MYSQL",
            "connProfile": {
                "connType": "PUBLIC",
                "host": "10.0.0.2",
                "port": 3306,
                "user": "dm_sync_user",
                "password": "Passw0rd!"
            },
            // Binlog coordinates are usually required when starting from existing data
            "binlogName": "mysql-bin.000777",
            "binlogPos": 12345,
            "binlogGtid": "3E11FA47-71CA-11E1-9E33-C80AA9429562:1-12345"
        }
    ]
}`
)

type templateVariant struct {
	heading string
	body    string
}

var allowedTemplateModes = []string{"all", "incremental"}

var definitionTemplates = map[string]templateVariant{
	"all": {
		heading: "Definition template (mode = ALL)",
		body:    migrationDefinitionAllTemplate,
	},
	"incremental": {
		heading: "Definition template (mode = INCREMENTAL)",
		body:    migrationDefinitionIncrementalTemplate,
	},
}

func TemplateCmd(h *internal.Helper) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Short:   "Show migration JSON templates",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("  Show the ALL mode migration template:\n  $ %[1]s serverless migration template --modetype all\n\n  Show the INCREMENTAL migration template:\n  $ %[1]s serverless migration template --modetype incremental\n", config.CliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := cmd.Flags().GetString(flag.MigrationModeType)
			if err != nil {
				return err
			}
			return renderMigrationTemplate(h, strings.ToLower(mode))
		},
	}

	cmd.Flags().String(flag.MigrationModeType, "", "Migration mode template to show, one of [\"all\", \"incremental\"].")
	if err := cmd.MarkFlagRequired(flag.MigrationModeType); err != nil {
		panic(err)
	}
	return cmd
}

func renderMigrationTemplate(h *internal.Helper, mode string) error {
	variant, ok := definitionTemplates[mode]
	if !ok {
		return fmt.Errorf("unknown mode %q, allowed values: %s", mode, strings.Join(allowedTemplateModes, ", "))
	}

	fmt.Fprintln(h.IOStreams.Out, color.GreenString(variant.heading))
	fmt.Fprintln(h.IOStreams.Out, variant.body)
	return nil
}
