// Copyright 2026 PingCAP, Inc.
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

package privatelink

import (
	"fmt"
	"strings"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	pl "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type ListOpts struct {
	interactive bool
}

func (c ListOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkConnectionState,
	}
}

func (c ListOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
	}
}

func (c *ListOpts) MarkInteractive(cmd *cobra.Command) error {
	flags := c.NonInteractiveFlags()
	for _, fn := range flags {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			c.interactive = false
			break
		}
	}
	if !c.interactive {
		for _, fn := range c.RequiredFlags() {
			if err := cmd.MarkFlagRequired(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func ListCmd(h *internal.Helper) *cobra.Command {
	opts := ListOpts{
		interactive: true,
	}

	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List private link connections",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Example: fmt.Sprintf(`  List private link connections in interactive mode:
  $ %[1]s serverless private-link-connection list

  List private link connections in non-interactive mode:
  $ %[1]s serverless private-link-connection list -c <cluster-id>

  List private link connections with json format in non-interactive mode:
  $ %[1]s serverless private-link-connection list -c <cluster-id> -o json`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID string
			var stateFilter *pl.PrivateLinkConnectionServiceListPrivateLinkConnectionsStateParameter
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
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				stateRaw, err := cmd.Flags().GetString(flag.PrivateLinkConnectionState)
				if err != nil {
					return errors.Trace(err)
				}
				if strings.TrimSpace(stateRaw) != "" {
					parsedState, err := normalizePrivateLinkConnectionState(stateRaw)
					if err != nil {
						return err
					}
					stateFilter = &parsedState
				}
			}

			total, items, err := cloud.RetrievePrivateLinkConnections(ctx, clusterID, h.QueryPageSize, stateFilter, d)
			if err != nil {
				return err
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				res := &pl.ListPrivateLinkConnectionsResponse{
					PrivateLinkConnections: items,
					TotalSize:              &total,
				}
				err := output.PrintJson(h.IOStreams.Out, res)
				if err != nil {
					return errors.Trace(err)
				}
				return nil
			}
			if format != output.HumanFormat {
				return fmt.Errorf("unsupported output format: %s", format)
			}

			columns := []output.Column{
				"ID",
				"DisplayName",
				"Type",
				"State",
				"EndpointServiceName",
				"EndpointServiceRegion",
				"CreateTime",
			}

			var rows []output.Row
			for _, item := range items {
				rows = append(rows, output.Row{
					valueOrEmpty(item.PrivateLinkConnectionId),
					item.DisplayName,
					string(item.Type),
					privateLinkConnectionState(item.State),
					endpointServiceName(item),
					endpointServiceRegion(item),
					privateLinkConnectionTime(item.CreateTime),
				})
			}

			err = output.PrintHumanTable(h.IOStreams.Out, columns, rows)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
	}

	listCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	listCmd.Flags().StringP(flag.PrivateLinkConnectionState, "", "", "Filter by private link connection state.")
	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return listCmd
}

func normalizePrivateLinkConnectionState(input string) (pl.PrivateLinkConnectionServiceListPrivateLinkConnectionsStateParameter, error) {
	value := strings.TrimSpace(strings.ToUpper(input))
	state := pl.PrivateLinkConnectionServiceListPrivateLinkConnectionsStateParameter(value)
	if state.IsValid() {
		return state, nil
	}
	return "", fmt.Errorf("invalid private link connection state: %s", input)
}

func endpointServiceName(item pl.PrivateLinkConnection) string {
	if item.AwsEndpointService != nil {
		return item.AwsEndpointService.Name
	}
	if item.AlicloudEndpointService != nil {
		return item.AlicloudEndpointService.Name
	}
	return ""
}

func endpointServiceRegion(item pl.PrivateLinkConnection) string {
	if item.AwsEndpointService != nil && item.AwsEndpointService.Region != nil {
		return *item.AwsEndpointService.Region
	}
	return ""
}

func privateLinkConnectionState(state *pl.PrivateLinkConnectionStateEnum) string {
	if state == nil {
		return ""
	}
	return string(*state)
}

func privateLinkConnectionTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
