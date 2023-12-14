// Copyright 2023 PingCAP, Inc.
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

package branch

import (
	"encoding/json"
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	branchApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client/branch_service"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type DescribeOpts struct {
	interactive bool
}

func (c DescribeOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.BranchID,
	}
}

func (c *DescribeOpts) MarkInteractive(cmd *cobra.Command) error {
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
		for _, fn := range flags {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DescribeCmd(h *internal.Helper) *cobra.Command {
	opts := DescribeOpts{
		interactive: true,
	}

	var describeCmd = &cobra.Command{
		Use:     "describe",
		Short:   "Describe a branch",
		Aliases: []string{"get"},
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Get the branch info in interactive mode:
  $ %[1]s serverless branch describe

  Get the branch info with Basic view in non-interactive mode:
  $ %[1]s serverless branch describe -c <cluster-id> -b <branch-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.MarkInteractive(cmd)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}

			var branchID string
			var clusterID string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				branch, err := cloud.GetSelectedBranch(clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				branchID = branch.ID
			} else {
				// non-interactive mode, get values from flags
				bID, err := cmd.Flags().GetString(flag.BranchID)
				if err != nil {
					return errors.Trace(err)
				}

				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				branchID = bID
				clusterID = cID
			}

			params := branchApi.NewBranchServiceGetBranchParams().
				WithClusterID(clusterID).WithBranchID(branchID)
			if err != nil {
				return errors.Trace(err)
			}
			branch, err := d.GetBranch(params)
			if err != nil {
				return errors.Trace(err)
			}

			v, err := json.MarshalIndent(branch.Payload, "", "  ")
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, string(v))
			return nil
		},
	}

	describeCmd.Flags().StringP(flag.BranchID, flag.BranchIDShort, "", "The ID of the branch to be described")
	describeCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID of the branch to be described")
	describeCmd.MarkFlagsRequiredTogether(flag.BranchID, flag.ClusterID)
	return describeCmd
}
