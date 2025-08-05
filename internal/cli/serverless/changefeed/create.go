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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cdc"
)

type CreateOpts struct {
	interactive bool
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ChangefeedType,
		flag.ChangefeedKafka,
		flag.ChangefeedMySQL,
		flag.ChangefeedFilter,
		flag.ChangefeedStartTSO,
		flag.ChangefeedName,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ChangefeedType,
		flag.ChangefeedFilter,
	}
}

func (c *CreateOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := c.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	// Mark required flags
	if !c.interactive {
		for _, fn := range c.RequiredFlags() {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{interactive: true}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a changefeed",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a changefeed in interactive mode:
  $ %[1]s serverless changefeed create

  Create a changefeed in non-interactive mode:
  $ %[1]s serverless changefeed create -c <cluster-id> --type KAFKA --kafka '{"network_info":{"network_type":"PUBLIC"},"broker":{"kafka_version":"VERSION_2XX","broker_endpoints":"52.34.156.155:9092","compression":"NONE"},"authentication":{"auth_type":"DISABLE"},"topic_partition_config":{"dispatch_type":"ONE_TOPIC","default_topic":"default-topic","replication_factor":1,"partition_num":1,"partition_dispatchers":[{"partition_type":"TABLE","matcher":["*.*"]}]},"data_format":{"protocol":"CANAL_JSON"}}' --filter '{"filterRule":["test.*"], "mode": "IGNORE_NOT_SUPPORT_TABLE"}'

  Create a changefeed named myfeed with specified start tso a in non-interactive mode:
  $ %[1]s serverless changefeed create -c <cluster-id> --name myfeed --type KAFKA --kafka <kafka-json> --filter <filter-json> --start-tso 458996254096228352
`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, name, kafkaStr, mysqlStr, filterStr string
			var startTSO uint64
			var filter cdc.CDCFilter
			var kafkaInfo cdc.Kafka
			var mysqlInfo cdc.MySQL
			var changefeedType cdc.ChangefeedTypeEnum

			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID
				changefeedType, err = GetSelectedChangefeedType()
				if err != nil {
					return err
				}

				var startTSOStr string
				switch changefeedType {
				case cdc.CHANGEFEEDTYPEENUM_KAFKA:
					inputs := []string{flag.ChangefeedName, flag.ChangefeedKafka, flag.ChangefeedFilter, flag.ChangefeedStartTSO}
					textInput, err := ui.InitialInputModel(inputs, createKafkaInputDescription)
					if err != nil {
						return err
					}
					name = textInput.Inputs[0].Value()
					kafkaStr = textInput.Inputs[1].Value()
					filterStr = textInput.Inputs[2].Value()
					startTSOStr = textInput.Inputs[3].Value()
				case cdc.CHANGEFEEDTYPEENUM_MYSQL:
					inputs := []string{flag.ChangefeedName, flag.ChangefeedMySQL, flag.ChangefeedFilter, flag.ChangefeedStartTSO}
					textInput, err := ui.InitialInputModel(inputs, createKafkaInputDescription)
					if err != nil {
						return err
					}
					name = textInput.Inputs[0].Value()
					mysqlStr = textInput.Inputs[1].Value()
					filterStr = textInput.Inputs[2].Value()
					startTSOStr = textInput.Inputs[3].Value()
				default:
					return errors.Errorf("currently only %s and %s type is supported", cdc.CHANGEFEEDTYPEENUM_KAFKA, cdc.CHANGEFEEDTYPEENUM_MYSQL)
				}
				if startTSOStr == "" {
					startTSO = 0
				} else {
					_, err = fmt.Sscanf(startTSOStr, "%d", &startTSO)
					if err != nil {
						return errors.New("invalid start-tso, must be uint64")
					}
				}
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				name, err = cmd.Flags().GetString(flag.ChangefeedName)
				if err != nil {
					return errors.Trace(err)
				}
				changefeedTypeStr, err := cmd.Flags().GetString(flag.ChangefeedType)
				if err != nil {
					return errors.Trace(err)
				}
				changefeedType = cdc.ChangefeedTypeEnum(changefeedTypeStr)
				kafkaStr, err = cmd.Flags().GetString(flag.ChangefeedKafka)
				if err != nil {
					return errors.Trace(err)
				}
				mysqlStr, err = cmd.Flags().GetString(flag.ChangefeedMySQL)
				if err != nil {
					return errors.Trace(err)
				}
				filterStr, err = cmd.Flags().GetString(flag.ChangefeedFilter)
				if err != nil {
					return errors.Trace(err)
				}
				startTSO, err = cmd.Flags().GetUint64(flag.ChangefeedStartTSO)
				if err != nil {
					return errors.Trace(err)
				}
			}

			// check all the parameters
			switch changefeedType {
			case cdc.CHANGEFEEDTYPEENUM_KAFKA:
				if kafkaStr == "" {
					return errors.New("kafka info is required")
				}
				if err := json.Unmarshal([]byte(kafkaStr), &kafkaInfo); err != nil {
					return errors.New("invalid kafka info, please use JSON format")
				}
			case cdc.CHANGEFEEDTYPEENUM_MYSQL:
				if mysqlStr == "" {
					return errors.New("mysql info is required")
				}
				if err := json.Unmarshal([]byte(mysqlStr), &mysqlInfo); err != nil {
					return errors.New("invalid mysql info, please use JSON format")
				}
			default:
				return errors.New("currently only KAFKA and MYSQL type is supported")
			}

			if filterStr == "" {
				return errors.New("filter is required")
			}
			if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
				return errors.New("invalid filter, please use JSON format")
			}

			// create the changefeed
			body := &cdc.ChangefeedServiceCreateChangefeedBody{
				DisplayName: &name,
				Sink: cdc.SinkInfo{
					Type: cdc.ChangefeedTypeEnum(changefeedType),
				},
				Filter: filter,
			}

			switch body.Sink.Type {
			case cdc.CHANGEFEEDTYPEENUM_KAFKA:
				body.Sink.Kafka = &kafkaInfo
			case cdc.CHANGEFEEDTYPEENUM_MYSQL:
				body.Sink.Mysql = &mysqlInfo
			}
			if startTSO == 0 {
				mode := cdc.STARTMODEENUM_FROM_NOW
				body.StartPosition = &cdc.StartPosition{
					Mode: &mode,
				}
			} else {
				mode := cdc.STARTMODEENUM_FROM_TSO
				body.StartPosition = &cdc.StartPosition{
					Mode: &mode,
					Tso:  aws.String(strconv.FormatUint(startTSO, 10)),
				}
			}

			resp, err := d.CreateChangefeed(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("changefeed %s created", resp.ChangefeedId))
			if err != nil {
				return err
			}
			return nil
		},
	}

	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	createCmd.Flags().String(flag.ChangefeedName, "", "The name of the changefeed.")
	createCmd.Flags().String(flag.ChangefeedType, "", fmt.Sprintf("The type of the changefeed, one of %q", cdc.AllowedChangefeedTypeEnumEnumValues))
	createCmd.Flags().String(flag.ChangefeedKafka, "", "Kafka information in JSON format, use \"ticloud serverless changefeed template\" to see templates.")
	createCmd.Flags().String(flag.ChangefeedMySQL, "", "MySQL information in JSON format, use \"ticloud serverless changefeed template\" to see templates.")
	createCmd.Flags().String(flag.ChangefeedFilter, "", "Filter in JSON format, use \"ticloud serverless changefeed template\" to see templates.")
	createCmd.Flags().Uint64(flag.ChangefeedStartTSO, 0, "Start TSO for the changefeed, default to current TSO. See https://docs.pingcap.com/tidb/stable/tso/ for more information about TSO.")

	return createCmd
}
