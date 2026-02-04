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

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type DeleteOpts struct {
	interactive bool
}

func (o DeleteOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ExampleID,
		flag.Force,
	}
}

func (o DeleteOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
	}
}

func (o *DeleteOpts) MarkInteractive(cmd *cobra.Command) error {
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

func DeleteCmd(h *internal.Helper) *cobra.Command {
	opts := &DeleteOpts{interactive: true}
	var force bool

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "Delete an example resource",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf(`  Delete an example resource (interactive):
  $ %[1]s serverless example delete

  Delete an example resource (non-interactive):
  $ %[1]s serverless example delete -c <cluster-id> --example-id <example-id>`, config.CliName),
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
				force, err = cmd.Flags().GetBool(flag.Force)
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

			if !force {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support prompt, run with --force to skip confirmation")
				}
				var confirm string
				if err := survey.AskOne(&survey.Input{Message: DeleteConfirmPrompt}, &confirm); err != nil {
					return err
				}
				if confirm != "yes" {
					return errors.New("deletion cancelled")
				}
			}

			// TODO implement the create logic
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("example %s deleted", "example-id"))
			if err != nil {
				return err
			}
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	cmd.Flags().String(flag.ExampleID, "", "The example ID.")
	cmd.Flags().BoolVar(&force, flag.Force, false, "Delete without confirmation.")
	return cmd
}
