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
	"slices"
	"strconv"
	"time"

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
		flag.ChangefeedFilter,
		flag.ChangefeedStartTSO,
		flag.ChangefeedName,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ChangefeedType,
		flag.ChangefeedKafka,
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

			var clusterID, name, changefeedType, kafkaStr, filterStr string
			var startTSO uint64
			var kafkaInfo cdc.KafkaInfo
			var filter cdc.CDCFilter

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

				inputs := []string{flag.ChangefeedName, flag.ChangefeedType, flag.ChangefeedKafka, flag.ChangefeedFilter, flag.ChangefeedStartTSO}
				textInput, err := ui.InitialInputModel(inputs, inputDescription)
				if err != nil {
					return err
				}
				name = textInput.Inputs[0].Value()
				changefeedType = textInput.Inputs[1].Value()
				kafkaStr = textInput.Inputs[2].Value()
				filterStr = textInput.Inputs[3].Value()
				startTSOStr := textInput.Inputs[4].Value()
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
				changefeedType, err = cmd.Flags().GetString(flag.ChangefeedType)
				if err != nil {
					return errors.Trace(err)
				}
				if changefeedType == "" {
					return errors.New("type is required")
				}
				kafkaStr, err = cmd.Flags().GetString(flag.ChangefeedKafka)
				if err != nil {
					return errors.Trace(err)
				}
				if kafkaStr == "" {
					return errors.New("kafka info is required")
				}
				if err := json.Unmarshal([]byte(kafkaStr), &kafkaInfo); err != nil {
					return errors.New("invalid kafka info, please use JSON format")
				}
				filterStr, err = cmd.Flags().GetString(flag.ChangefeedFilter)
				if err != nil {
					return errors.Trace(err)
				}
				if filterStr == "" {
					return errors.New("filter is required")
				}
				if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
					return errors.New("invalid filter, please use JSON format")
				}
				startTSO, err = cmd.Flags().GetUint64(flag.ChangefeedStartTSO)
				if err != nil {
					return errors.Trace(err)
				}
			}

			// check all the parameters
			if changefeedType == "" {
				return errors.New("type is required")
			}
			if kafkaStr == "" {
				return errors.New("kafka info is required")
			}
			if err := json.Unmarshal([]byte(kafkaStr), &kafkaInfo); err != nil {
				return errors.New("invalid kafka info, please use JSON format")
			}
			if filterStr == "" {
				return errors.New("filter is required")
			}
			if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
				return errors.New("invalid filter, please use JSON format")
			}

			if !slices.Contains(cdc.AllowedConnectorTypeEnumEnumValues, cdc.ConnectorTypeEnum(changefeedType)) {
				return errors.New("currently only KAFKA type is supported")
			}

			// create the changefeed
			body := &cdc.ConnectorServiceCreateConnectorBody{
				Name: &name,
				Sink: cdc.SinkInfo{
					Type: cdc.ConnectorTypeEnum(changefeedType),
				},
				Filter: filter,
			}

			switch body.Sink.Type {
			case cdc.CONNECTORTYPEENUM_KAFKA:
				body.Sink.Kafka = &kafkaInfo
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

			now := time.Now().UTC().Format(time.RFC3339)
			mode := cdc.STARTMODEENUM_FROM_UTC
			body.StartPosition = &cdc.StartPosition{
				Mode: &mode,
				Utc:  &now,
			}

			resp, err := d.CreateConnector(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("changefeed %s created", resp.ConnectorId))
			if err != nil {
				return err
			}
			return nil
		},
	}

	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	createCmd.Flags().String(flag.ChangefeedName, "", "The name of the changefeed.")
	createCmd.Flags().String(flag.ChangefeedType, "", fmt.Sprintf("The type of the changefeed, one of %q", cdc.AllowedConnectorTypeEnumEnumValues))
	createCmd.Flags().String(flag.ChangefeedKafka, "", "Kafka information in JSON format, use \"ticloud serverless changefeed template\" to see templates.")
	createCmd.Flags().String(flag.ChangefeedFilter, "", "Filter in JSON format, use \"ticloud serverless changefeed template\" to see templates.")
	createCmd.Flags().Uint64(flag.ChangefeedStartTSO, 0, "Start TSO for the changefeed, default to current TSO. See https://docs.pingcap.com/tidb/stable/tso/ for more information about TSO.")

	return createCmd
}

// inputDescription 用于交互式输入提示
var inputDescription = map[string]string{
	flag.ChangefeedName:     "The name of the changefeed, skip to use the default name",
	flag.ChangefeedType:     fmt.Sprintf("The type of the changefeed, one of %q", cdc.AllowedConnectorTypeEnumEnumValues),
	flag.ChangefeedKafka:    "Kafka information in JSON format, use \"ticloud serverless changefeed template\" to see templates.",
	flag.ChangefeedFilter:   "Filter in JSON format, use \"ticloud serverless changefeed template\" to see templates.",
	flag.ChangefeedStartTSO: "Start TSO (uint64) for the changefeed, skip to use the current TSO. See https://docs.pingcap.com/tidb/stable/tso/ for more information about TSO.",
}
