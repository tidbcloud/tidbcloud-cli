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

package example

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type DescribeOpts struct {
	interactive bool
}

func (o DescribeOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ExampleID,
	}
}

func (o DescribeOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ExampleID,
	}
}

func (o *DescribeOpts) MarkInteractive(cmd *cobra.Command) error {
	o.interactive = true
	for _, fn := range o.NonInteractiveFlags() {
		if f := cmd.Flags().Lookup(fn); f != nil && f.Changed {
			o.interactive = false
		}
	}
	if !o.interactive {
		for _, fn := range o.RequiredFlags() {
			if err := cmd.MarkFlagRequired(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func DescribeCmd(h *internal.Helper) *cobra.Command {
	opts := &DescribeOpts{interactive: true}

	cmd := &cobra.Command{
		Use:     "describe",
		Aliases: []string{"get"},
		Short:   "Describe an example resource",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Describe an example resource (interactive):
  $ %[1]s serverless example describe

  Describe an example resource (non-interactive):
  $ %[1]s serverless example describe -c <cluster-id> --example-id <example-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}
			var clusterID, exampleID string
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

				example, err := GetSelectedExample(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				exampleID = example.ID
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				exampleID, err = cmd.Flags().GetString(flag.ExampleID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if clusterID == "" {
				return errors.New("cluster id is required")
			}
			if exampleID == "" {
				return errors.New("example id is required")
			}

			// TODO implement the get logic, now just mock the payload
			payload := map[string]string{
				"message":    UnwiredFeaturePrompt,
				"cluster_id": clusterID,
				"example_id": exampleID,
			}
			return output.PrintJson(h.IOStreams.Out, payload)
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	cmd.Flags().String(flag.ExampleID, "", "The example ID.")
	return cmd
}
