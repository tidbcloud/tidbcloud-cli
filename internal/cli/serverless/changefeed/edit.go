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

type UpdateOpts struct {
	interactive bool
}

func (c UpdateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ChangefeedID,
		flag.ChangefeedName,
		flag.ChangefeedKafka,
		flag.ChangefeedFilter,
	}
}

func (c UpdateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ChangefeedID,
	}
}

func (c *UpdateOpts) MarkInteractive(cmd *cobra.Command) error {
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

func EditCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateOpts{interactive: true}

	var editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit a changefeed (kafka and filter must be fully specified)",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Update a changefeed in interactive mode:
  $ %[1]s serverless changefeed update

  Update the name, kafka, and filter of a changefeed in non-interactive mode:
  $ %[1]s serverless changefeed update -c <cluster-id> --changefeed-id <changefeed-id> --name newname --kafka '{"broker":{"address":"localhost:9092"}}' --filter '{"filterRule":["test.t1"]}'

  Note: --kafka and --filter must be fully specified, partial updates are not supported for these fields.
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

			var clusterID, changefeedID, kafkaStr, filterStr string
			var kafkaInfo cdc.KafkaInfo
			var filter cdc.CDCFilter
			var name *string

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

				cf, err := cloud.GetSelectedChangefeed(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				changefeedID = cf.ID

				inputs := []string{flag.ChangefeedName, flag.ChangefeedKafka, flag.ChangefeedFilter}
				textInput, err := ui.InitialInputModel(inputs, updateInputDescriptionInteractive)
				if err != nil {
					return err
				}
				nameStr := textInput.Inputs[0].Value()
				kafkaStr = textInput.Inputs[1].Value()
				filterStr = textInput.Inputs[2].Value()
				if nameStr != "" {
					name = &nameStr
				}
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				changefeedID, err = cmd.Flags().GetString(flag.ChangefeedID)
				if err != nil {
					return errors.Trace(err)
				}
				if cmd.Flags().Changed(flag.ChangefeedName) {
					nameStr, err := cmd.Flags().GetString(flag.ChangefeedName)
					if err != nil {
						return errors.Trace(err)
					}
					name = &nameStr
				}
				kafkaStr, err = cmd.Flags().GetString(flag.ChangefeedKafka)
				if err != nil {
					return errors.Trace(err)
				}
				filterStr, err = cmd.Flags().GetString(flag.ChangefeedFilter)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if filterStr == "" {
				return errors.New("filter (--filter) is required and must be fully specified")
			}
			if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
				return errors.New("invalid filter, please use JSON format")
			}

			body := &cdc.ConnectorServiceEditConnectorBody{
				ConnectorId: changefeedID,
				Name:        name,
				Filter:      filter,
			}

			changefeed, err := d.GetConnector(ctx, clusterID, changefeedID)
			if err != nil {
				return errors.Trace(err)
			}
			switch changefeed.Sink.Type {
			case cdc.CONNECTORTYPEENUM_KAFKA:
				if kafkaStr == "" {
					return errors.New("kafka info (--kafka) is required and must be fully specified")
				}
				if err := json.Unmarshal([]byte(kafkaStr), &kafkaInfo); err != nil {
					return errors.New("invalid kafka info, please use JSON format")
				}
				body.Sink = cdc.SinkInfo{
					Type:  changefeed.Sink.Type,
					Kafka: &kafkaInfo,
				}
			default:
				return errors.Errorf("unsupported changefeed sink type: %s", changefeed.Sink.Type)
			}

			_, err = d.EditConnector(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("changefeed %s updated", changefeedID))
			if err != nil {
				return err
			}
			return nil
		},
	}

	editCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	editCmd.Flags().String(flag.ChangefeedID, "", "The ID of the changefeed to be updated.")
	editCmd.Flags().String(flag.ChangefeedName, "", "The new name of the changefeed (optional).")
	editCmd.Flags().String(flag.ChangefeedKafka, "", "Kafka info in JSON format, see KafkaInfo struct. (required, must be fully specified)")
	editCmd.Flags().String(flag.ChangefeedFilter, "", "Filter in JSON format, see CDCFilter struct. (required, must be fully specified)")

	return editCmd
}

// updateInputDescriptionInteractive 用于交互式输入提示（kafka 和 filter 必填）
var updateInputDescriptionInteractive = map[string]string{
	flag.ChangefeedName:   "The new name of the changefeed (optional)",
	flag.ChangefeedKafka:  "Kafka info in JSON format, see KafkaInfo struct (required, must be fully specified)",
	flag.ChangefeedFilter: "Filter in JSON format, see CDCFilter struct (required, must be fully specified)",
}
