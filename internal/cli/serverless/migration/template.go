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
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	pkgmigration "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/migration"
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
            // Required: source database type. Supported values: MYSQL, ALICLOUD_RDS_MYSQL, AWS_RDS_MYSQL
            "sourceType": "MYSQL",
            "connProfile": {
                // Required connection type. Supported values: PUBLIC, PRIVATE_LINK
                // PUBLIC connections require host
                "connType": "PUBLIC",
                "host": "10.0.0.8",
                // PRIVATE_LINK connections use endpointId. Get endpointId by 'ticloud plc' commands.
                "connType": "PRIVATE_LINK",
                "endpointId": "pl-xxxxxxxx",
                "host": "10.0.0.2",
                "port": 3306,
                "user": "dm_sync_user",
                "password": "Passw0rd!",
                // optional TLS settings
                "security": {
                    // TLS materials must be Base64 encoded
                    "sslCaContent": "<base64-of-ca.pem>",
                    "sslCertContent": "<base64-of-client-cert.pem>",
                    "sslKeyContent": "<base64-of-client-key.pem>",
                    "certAllowedCn": ["client-cn"]
                }
            },
            // Optional block/allow rules to control synced schemas/tables (mutually exclusive with routeRules)
            "baRules": {
                "doDbs": ["app_db"],
                "doTables": [
                    {"schema": "app_db", "table": "orders"},
                    {"schema": "app_db", "table": "customers"}
                ]
            },
            // Optional route rules to rename objects during migration (mutually exclusive with baRules)
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
            // Required: source database type. Supported values: MYSQL, ALICLOUD_RDS_MYSQL, AWS_RDS_MYSQL
            "sourceType": "MYSQL",
            "connProfile": {
                // Required connection type. Supported values: PUBLIC, PRIVATE_LINK
                // PUBLIC connections require host
                "connType": "PUBLIC",
                "host": "10.0.0.8",
                // PRIVATE_LINK connections use endpointId. Get endpointId by 'ticloud plc' commands.
                "connType": "PRIVATE_LINK",
                "endpointId": "pl-xxxxxxxx",
                "port": 3306,
                "user": "dm_sync_user",
                "password": "Passw0rd!",
                // optional TLS settings
                "security": {
                    // TLS materials must be Base64 encoded
                    "sslCaContent": "<base64-of-ca.pem>",
                    "sslCertContent": "<base64-of-client-cert.pem>",
                    "sslKeyContent": "<base64-of-client-key.pem>",
                    "certAllowedCn": ["client-cn"]
                }
            },
            // Optional block/allow rules when only part of the data should be replicated (mutually exclusive with routeRules)
            "baRules": {
                "doDbs": ["app_db"],
                "doTables": [
                    {"schema": "app_db", "table": "orders"}
                ]
            },
            // Optional route rule sample for remapping objects during incremental sync (mutually exclusive with baRules)
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
)

type templateVariant struct {
	heading string
	body    string
}

var allowedTemplateModes = []pkgmigration.TaskMode{pkgmigration.TASKMODE_ALL, pkgmigration.TASKMODE_INCREMENTAL}

var definitionTemplates = map[pkgmigration.TaskMode]templateVariant{
	pkgmigration.TASKMODE_ALL: {
		heading: "Definition template (mode = ALL)",
		body:    migrationDefinitionAllTemplate,
	},
	pkgmigration.TASKMODE_INCREMENTAL: {
		heading: "Definition template (mode = INCREMENTAL)",
		body:    migrationDefinitionIncrementalTemplate,
	},
}

func TemplateCmd(h *internal.Helper) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Short:   "Show migration JSON templates",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("  Show the ALL mode migration template:\n  $ %[1]s serverless migration template --mode all\n\n  Show the INCREMENTAL migration template:\n  $ %[1]s serverless migration template --mode incremental\n", config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cmd.MarkFlagRequired(flag.MigrationMode)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			modeValue, err := cmd.Flags().GetString(flag.MigrationMode)
			if err != nil {
				return err
			}
			mode, err := parseTemplateMode(modeValue)
			if err != nil {
				return err
			}
			return renderMigrationTemplate(h, mode)
		},
	}

	cmd.Flags().String(
		flag.MigrationMode,
		"",
		fmt.Sprintf(
			"Migration mode template to show, one of [%s].",
			strings.Join(allowedTemplateModeStrings(), ", "),
		),
	)
	return cmd
}

func renderMigrationTemplate(h *internal.Helper, mode pkgmigration.TaskMode) error {
	variant, ok := definitionTemplates[mode]
	if !ok {
		return fmt.Errorf("unknown mode %q, allowed values: %s", mode, strings.Join(allowedTemplateModeStrings(), ", "))
	}

	fmt.Fprintln(h.IOStreams.Out, color.GreenString(variant.heading))
	fmt.Fprintln(h.IOStreams.Out, variant.body)
	return nil
}

func parseTemplateMode(raw string) (pkgmigration.TaskMode, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("mode is required; use --%s", flag.MigrationMode)
	}
	normalized := strings.ToUpper(trimmed)
	mode := pkgmigration.TaskMode(normalized)
	if _, ok := definitionTemplates[mode]; ok {
		return mode, nil
	}
	return "", fmt.Errorf("unknown mode %q, allowed values: %s", trimmed, strings.Join(allowedTemplateModeStrings(), ", "))
}

func allowedTemplateModeStrings() []string {
	values := make([]string, 0, len(allowedTemplateModes))
	for _, mode := range allowedTemplateModes {
		values = append(values, strings.ToLower(string(mode)))
	}
	return values
}
