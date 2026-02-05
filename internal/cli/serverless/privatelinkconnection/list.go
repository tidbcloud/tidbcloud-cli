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

package privatelinkconnection

import (
	"fmt"
	"strings"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type ListOpts struct {
	interactive bool
}

func (c ListOpts) NonInteractiveFlags() []string {
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
		for _, fn := range flags {
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
		Use:   "list",
		Short: "List private link connections",
		Args:  cobra.NoArgs,
		Aliases: []string{
			"ls",
		},
		Example: fmt.Sprintf(`  List private link connections in interactive mode:
  $ %[1]s serverless private-link-connection list

  List private link connections in non-interactive mode:
  $ %[1]s serverless private-link-connection list -c <cluster-id>

  List private link connections in non-interactive mode with state filter:
  $ %[1]s serverless private-link-connection list -c <cluster-id> --state ACTIVE`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.MarkInteractive(cmd); err != nil {
				return errors.Trace(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID string
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
			}

			var statePtr *privatelink.PrivateLinkConnectionServiceListPrivateLinkConnectionsStateParameter
			stateValue, err := cmd.Flags().GetString(flag.PrivateLinkConnectionState)
			if err != nil {
				return errors.Trace(err)
			}
			if stateValue != "" {
				state := privatelink.PrivateLinkConnectionServiceListPrivateLinkConnectionsStateParameter(stateValue)
				if !state.IsValid() {
					return fmt.Errorf("invalid private link connection state: %s", stateValue)
				}
				statePtr = &state
			}

			total, items, err := cloud.RetrievePrivateLinkConnections(ctx, clusterID, h.QueryPageSize, statePtr, d)
			if err != nil {
				return err
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				res := &privatelink.ListPrivateLinkConnectionsResponse{
					PrivateLinkConnections: items,
					TotalSize:              &total,
				}
				return errors.Trace(output.PrintJson(h.IOStreams.Out, res))
			} else if format == output.HumanFormat {
				columns := []output.Column{
					"ID",
					"DisplayName",
					"Type",
					"State",
					"CreateTime",
				}
				var rows []output.Row
				for _, item := range items {
					id := ""
					if item.PrivateLinkConnectionId != nil {
						id = *item.PrivateLinkConnectionId
					}
					state := ""
					if item.State != nil {
						state = string(*item.State)
					}
					createTime := ""
					if item.CreateTime != nil {
						createTime = item.CreateTime.Format(time.RFC3339)
					}
					rows = append(rows, output.Row{
						id,
						item.DisplayName,
						string(item.Type),
						state,
						createTime,
					})
				}
				return errors.Trace(output.PrintHumanTable(h.IOStreams.Out, columns, rows))
			}
			return fmt.Errorf("unsupported output format: %s", strings.ToLower(format))
		},
	}

	listCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	listCmd.Flags().String(flag.PrivateLinkConnectionState, "", "Filter by private link connection state.")
	listCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return listCmd
}
