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

package changefeed

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type PrivateLinkDeleteOpts struct {
	interactive bool
}

func (c PrivateLinkDeleteOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkServiceName,
	}
}

func (c *PrivateLinkDeleteOpts) MarkInteractive(cmd *cobra.Command) error {
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

func PrivateLinkDeleteCmd(h *internal.Helper) *cobra.Command {
	opts := PrivateLinkDeleteOpts{
		interactive: true,
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a changefeed private link endpoint",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Delete a changefeed private link endpoint in interactive mode:
  $ %[1]s serverless changefeed private-link delete

  Delete a changefeed private link endpoint in non-interactive mode:
  $ %[1]s serverless changefeed private-link delete -c <cluster-id> --private-link-service-name <service-name>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID, privateLinkServiceName string
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

				privateLinkServiceName, err = promptPrivateLinkServiceName()
				if err != nil {
					return err
				}
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				privateLinkServiceName, err = cmd.Flags().GetString(flag.PrivateLinkServiceName)
				if err != nil {
					return errors.Trace(err)
				}
			}

			privateLinkServiceName = strings.TrimSpace(privateLinkServiceName)
			if privateLinkServiceName == "" {
				return errors.New("private link service name is required")
			}

			_, err = d.DeletePrivateLinkEndpoint(ctx, clusterID, privateLinkServiceName)
			if err != nil {
				return errors.Trace(err)
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("private link endpoint %s is deleted", privateLinkServiceName))
			if err != nil {
				return err
			}
			return nil
		},
	}

	deleteCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	deleteCmd.Flags().StringP(flag.PrivateLinkServiceName, "", "", "The private link service name.")
	return deleteCmd
}
