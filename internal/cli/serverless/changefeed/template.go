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

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
)

const (
	KafkaInfoTemplateWithExplain = `{
        "network_info": {
                "network_type": "PUBLIC"
        },
        "broker": {
                // "kafka_version": "VERSION_2XX", "VERSION_3XX"
                "kafka_version": "VERSION_2XX",
                "broker_endpoints": "broker:9092",
                "tls_enable": false,
                // "compression": "NONE", "GZIP", "LZ4", "ZSTD", "SNAPPY"
                "compression": "NONE"
        },
        "authentication": {
                // "auth_type": "DISABLE", "SASL_PLAIN", "SASL_SCRAM_SHA_256", "SASL_SCRAM_SHA_512"
                "auth_type": "DISABLE",
                // required when auth_type is SASL_PLAIN, SASL_SCRAM_SHA_256, or SASL_SCRAM_SHA_512
                "user_name": "",
                // required when auth_type is SASL_PLAIN, SASL_SCRAM_SHA_256, or SASL_SCRAM_SHA_512
                "password": ""
        },
        "data_format": {
                // "protocol": "CANAL_JSON", "AVRO", "OPEN_PROTOCOL"
                "protocol": "CANAL_JSON",
                // available when protocol is CANAL_JSON
                "enable_tidb_extension": false,
                // available when protocol is AVRO
                "avro_config": {
                        "decimal_handling_mode": "PRECISE",
                        "bigint_unsigned_handling_mode": "LONG",
                        "schema_registry": {
                                "schema_registry_endpoints": "",
                                "enable_http_auth": false,
                                "user_name": "",
                                "password": ""
                        }
                }
        },
        "topic_partition_config": {
                // "dispatch_type": "ONE_TOPIC", "BY_TABLE", "BY_DATABASE"
                "dispatch_type": "ONE_TOPIC",
                "default_topic": "test-topic",
                // required when dispatch_type is BY_TABLE or BY_DATABASE
                "topic_prefix": "_prefix",
                 // required when dispatch_type is BY_TABLE or BY_DATABASE
                "separator": "_",
                 // required when dispatch_type is BY_TABLE or BY_DATABASE
                "topic_suffix": "_suffix",
                "replication_factor": 1,
                "partition_num": 1,
                "partition_dispatchers": [
                  {
                        // "partition_type": "TABLE", "INDEX_VALUE", "TS", "COLUMN"
                        "partition_type": "TABLE",
                        // available when partition_type is TABLE
                        "matcher": ["*.*"],
                        // available when partition_type is INDEX_VALUE
                        "index_name": "index1",
                        // available when partition_type is COLUMN
                        "columns": ["col1", "col2"]
                  }
                ]
        },
        "column_selectors": [
          {
              "matcher": ["*.*"],
              "columns": ["col1", "col2"]
          }
        ]
}`

	KafkaInfoTemplate = `{
	"network_info": {
		"network_type": "PUBLIC"
	},
	"broker": {
		"kafka_version": "VERSION_2XX",
		"broker_endpoints": "broker1:9092,broker2:9092",
		"tls_enable": false,
		"compression": "NONE"
	},
	"authentication": {
		"auth_type": "DISABLE",
		"user_name": "",
		"password": ""
	},
	"data_format": {
		"protocol": "CANAL_JSON",
		"enable_tidb_extension": false,
		"avro_config": {
			"decimal_handling_mode": "PRECISE",
			"bigint_unsigned_handling_mode": "LONG",
			"schema_registry": {
				"schema_registry_endpoints": "",
				"enable_http_auth": false,
				"user_name": "",
				"password": ""
			}
		}
	},
	"topic_partition_config": {
		"dispatch_type": "ONE_TOPIC",
		"default_topic": "test-topic",
		"topic_prefix": "_prefix",
		"separator": "_",
		"topic_suffix": "_suffix",
		"replication_factor": 1,
		"partition_num": 1,
		"partition_dispatchers": [{
			"partition_type": "TABLE",
			"matcher": ["*.*"],
			"index_name": "index1",
			"columns": ["col1", "col2"]
		}]
	},
	"column_selectors": [{
		"matcher": ["*.*"],
		"columns": ["col1", "col2"]
	}]
}
  `

	CDCFilterTemplateWithExplain = `{
  "filterRule": ["test.t1", "test.t2"],
  // "mode": "IGNORE_NOT_SUPPORT_TABLE", "FORCE_SYNC"
  "mode": "IGNORE_NOT_SUPPORT_TABLE",
  "eventFilterRule": [
    {
      "matcher": ["test.t1", "test.t2"],
      "ignore_event": ["all dml", "all ddl"]
    }
  ]
}`

	CDCFilterTemplate = `{
  "filterRule": ["test.t1", "test.t2"],
  "mode": "IGNORE_NOT_SUPPORT_TABLE",
  "eventFilterRule": [
    {
      "matcher": ["test.t1", "test.t2"],
      "ignore_event": ["all dml", "all ddl"]
    }
  ]
}`
)

func TemplateCmd(h *internal.Helper) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "template",
		Short: "Show changefeed KafkaInfo and CDCFilter JSON templates",
		Example: `  Show changefeed templates:
  $ tidbcloud serverless changefeed template`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			explain, err := cmd.Flags().GetBool(flag.Explain)
			if err != nil {
				return err
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("KafkaInfo JSON template:"))
			if explain {
				fmt.Fprintln(h.IOStreams.Out, KafkaInfoTemplateWithExplain)
			} else {
				fmt.Fprintln(h.IOStreams.Out, KafkaInfoTemplate)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("CDCFilter JSON template:"))
			if explain {
				fmt.Fprintln(h.IOStreams.Out, CDCFilterTemplateWithExplain)
			} else {
				fmt.Fprintln(h.IOStreams.Out, CDCFilterTemplate)
			}
			return nil
		},
	}

	cmd.Flags().Bool(flag.Explain, false, "show template with explanations")

	return cmd
}
