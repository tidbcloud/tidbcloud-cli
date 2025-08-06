// Copyright 2025 PingCAP, Inc.
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

package changefeed

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
)

const (
	KafkaInfoTemplateWithExplain = `{
    "network": {
        "networkType": "PUBLIC",
        "publicEndpoints": "broker:9092"
    },
    "broker": {
        // "kafkaVersion": "VERSION_2XX", "VERSION_3XX"
        "kafkaVersion": "VERSION_2XX",
        // "compression": "NONE", "GZIP", "LZ4", "ZSTD", "SNAPPY"
        "compression": "NONE"
    },
    "authentication": {
        // "authType": "DISABLE", "SASL_PLAIN", "SASL_SCRAM_SHA_256", "SASL_SCRAM_SHA_512"
        "authType": "DISABLE",
        // required when authType is SASL_PLAIN, SASL_SCRAM_SHA_256, or SASL_SCRAM_SHA_512
        "userName": "",
        // required when authType is SASL_PLAIN, SASL_SCRAM_SHA_256, or SASL_SCRAM_SHA_512
        "password": "",
        "enableTls": false
    },
    "dataFormat": {
        // "protocol": "CANAL_JSON", "AVRO", "OPEN_PROTOCOL"
        "protocol": "CANAL_JSON",
        // available when protocol is CANAL_JSON
        "enableTidbExtension": false,
        // available when protocol is AVRO
        "avroConfig": {
            "decimalHandlingMode": "PRECISE",
            "bigintUnsignedHandlingMode": "LONG",
            // one of "confluentSchemaRegistry", "awsGlueSchemaRegistry"
            "confluentSchemaRegistry": {
                "endpoint": "",
                "enableHttpAuth": false,
                "userName": "",
                "password": ""
            },
            "awsGlueSchemaRegistry": {
                "region": "",
                "name": "",
                "accessKeyId": "",
                "secretAccessKey": ""
            }
        }
    },
    "topicPartitionConfig": {
        // "dispatchType": "ONE_TOPIC", "BY_TABLE", "BY_DATABASE"
        "dispatchType": "ONE_TOPIC",
        "defaultTopic": "test-topic",
        // required when dispatchType is BY_TABLE or BY_DATABASE
        "topicPrefix": "_prefix",
        // required when dispatchType is BY_TABLE or BY_DATABASE
        "separator": "_",
        // required when dispatchType is BY_TABLE or BY_DATABASE
        "topicSuffix": "_suffix",
        "replicationFactor": 1,
        "partitionNum": 1,
        "partitionDispatchers": [
            {
                // "partitionType": "TABLE", "INDEX_VALUE", "TS", "COLUMN"
                "partitionType": "TABLE",
                // available when partitionType is TABLE
                "matcher": ["*.*"],
                // available when partitionType is INDEX_VALUE
                "indexName": "index1",
                // available when partitionType is COLUMN
                "columns": ["col1", "col2"]
            }
        ]
    },
    "columnSelectors": [
        {
            "matcher": ["*.*"],
            "columns": ["col1", "col2"]
        }
    ]
}`

	KafkaInfoTemplate = `{
    "network": {
        "networkType": "PUBLIC",
        "publicEndpoints": "broker1:9092,broker2:9092"
    },
    "broker": {
        "kafkaVersion": "VERSION_2XX",
        "compression": "NONE"
    },
    "authentication": {
        "authType": "DISABLE",
        "userName": "",
        "password": "",
        "enableTls": false
    },
    "dataFormat": {
        "protocol": "CANAL_JSON",
        "enableTidbExtension": false,
        "avroConfig": {
            "decimalHandlingMode": "PRECISE",
            "bigintUnsignedHandlingMode": "LONG",
            "confluentSchemaRegistry": {
                "schemaRegistryEndpoints": "",
                "enableHttpAuth": false,
                "userName": "",
                "password": ""
            },
            "awsGlueSchemaRegistry": {
                "region": "",
                "name": "",
                "accessKeyId": "",
                "secretAccessKey": ""
            }
        }
    },
    "topicPartitionConfig": {
        "dispatchType": "ONE_TOPIC",
        "defaultTopic": "test-topic",
        "topicPrefix": "_prefix",
        "separator": "_",
        "topicSuffix": "_suffix",
        "replicationFactor": 1,
        "partitionNum": 1,
        "partitionDispatchers": [{
            "partitionType": "TABLE",
            "matcher": ["*.*"],
            "indexName": "index1",
            "columns": ["col1", "col2"]
        }]
    },
    "columnSelectors": [{
        "matcher": ["*.*"],
        "columns": ["col1", "col2"]
    }]
}`

	CDCFilterTemplateWithExplain = `{
    "filterRule": ["test.t1", "test.t2"],
    // "mode": "IGNORE_NOT_SUPPORT_TABLE", "FORCE_SYNC"
    "mode": "IGNORE_NOT_SUPPORT_TABLE",
    "eventFilterRule": [
        {
            "matcher": ["test.t1", "test.t2"],
            "ignoreEvent": ["all dml", "all ddl"]
        }
    ]
}`

	CDCFilterTemplate = `{
    "filterRule": ["test.t1", "test.t2"],
    "mode": "IGNORE_NOT_SUPPORT_TABLE",
    "eventFilterRule": [
        {
            "matcher": ["test.t1", "test.t2"],
            "ignoreEvent": ["all dml", "all ddl"]
        }
    ]
}`

	MySQLTemplateWithExplain = `{
    "network": {
        // required "PUBLIC", "PRIVATE"
        "networkType": "PUBLIC",
        "publicEndpoint": "127.0.0.1:3306"
    },
    "authentication": {
        // required the user name for MySQL
        "userName": "",
        // required the password for MySQL
        "password": "",
        // optional, enable TLS for MySQL connection
        "enableTls": false
    }
}`

	MySQLTemplate = `{
    "network": {
        "networkType": "PUBLIC",
        "publicEndpoint": "127.0.0.1:3306"
    },
    "authentication": {
        "userName": "",
        "password": "",
        "enableTls": false
    }
}`
)

func TemplateCmd(h *internal.Helper) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "template",
		Short: "Show changefeed KafkaInfo and CDCFilter JSON templates",
		Example: `  Show all changefeed templates:
  $ tidbcloud serverless changefeed template
  
  Show Kafka JSON template:
  $ tidbcloud serverless changefeed template --type kafka
  `,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			explain, err := cmd.Flags().GetBool(flag.Explain)
			if err != nil {
				return err
			}
			templateType, err := cmd.Flags().GetString(flag.ChangfeedTemplateType)
			if err != nil {
				return err
			}

			switch strings.ToLower(templateType) {
			case "kafka":
				if explain {
					fmt.Fprintln(h.IOStreams.Out, KafkaInfoTemplateWithExplain)
				} else {
					fmt.Fprintln(h.IOStreams.Out, KafkaInfoTemplate)
				}
			case "filter":
				if explain {
					fmt.Fprintln(h.IOStreams.Out, CDCFilterTemplateWithExplain)
				} else {
					fmt.Fprintln(h.IOStreams.Out, CDCFilterTemplate)
				}
			case "mysql":
				if explain {
					fmt.Fprintln(h.IOStreams.Out, MySQLTemplateWithExplain)
				} else {
					fmt.Fprintln(h.IOStreams.Out, MySQLTemplate)
				}
			default:
				fmt.Fprintln(h.IOStreams.Out, color.GreenString("KafkaInfo JSON template:"))
				if explain {
					fmt.Fprintln(h.IOStreams.Out, KafkaInfoTemplateWithExplain)
				} else {
					fmt.Fprintln(h.IOStreams.Out, KafkaInfoTemplate)
				}
				fmt.Fprintln(h.IOStreams.Out, color.GreenString("MySQLInfo JSON template:"))
				if explain {
					fmt.Fprintln(h.IOStreams.Out, MySQLTemplateWithExplain)
				} else {
					fmt.Fprintln(h.IOStreams.Out, MySQLTemplate)
				}
				fmt.Fprintln(h.IOStreams.Out, color.GreenString("CDCFilter JSON template:"))
				if explain {
					fmt.Fprintln(h.IOStreams.Out, CDCFilterTemplateWithExplain)
				} else {
					fmt.Fprintln(h.IOStreams.Out, CDCFilterTemplate)
				}

			}
			return nil
		},
	}

	cmd.Flags().Bool(flag.Explain, false, "show template with explanations")
	cmd.Flags().String(flag.ChangfeedTemplateType, "", "the type of changefeed template to show (kafka, mysql, filter)")

	return cmd
}
