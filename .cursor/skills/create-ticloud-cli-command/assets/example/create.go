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

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type CreateOpts struct {
	interactive bool
}

func (o CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
	}
}

func (o CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
	}
}

func (o *CreateOpts) MarkInteractive(cmd *cobra.Command) error {
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

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := &CreateOpts{interactive: true}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an example resource",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create an example resource (interactive):
  $ %[1]s serverless example create

  Create an example resource (non-interactive):
  $ %[1]s serverless example create -c <cluster-id> --display-name <name>

  Create with multi-line input flags (non-interactive):
  $ %[1]s serverless example create -c <cluster-id> --display-name <name> \
    --s3.uri <s3-uri> \
    --s3.access-key-id <access-key-id> \
    --s3.secret-access-key <secret-access-key>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}
			var clusterID, displayName string
			var s3URI, s3AccessKeyID, s3SecretAccessKey string
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

				displayName, err = GetDisplayNameInput()
				if err != nil {
					return err
				}

				s3URI, s3AccessKeyID, s3SecretAccessKey, err = GetS3Inputs()
				if err != nil {
					return err
				}
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				displayName, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}
				s3URI, err = cmd.Flags().GetString(flag.S3URI)
				if err != nil {
					return errors.Trace(err)
				}
				s3AccessKeyID, err = cmd.Flags().GetString(flag.S3AccessKeyID)
				if err != nil {
					return errors.Trace(err)
				}
				s3SecretAccessKey, err = cmd.Flags().GetString(flag.S3SecretAccessKey)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if clusterID == "" {
				return errors.New("cluster id is required")
			}
			if displayName == "" {
				return errors.New("display name is required")
			}

			// TODO implement the create logic
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("example %s created", "example-id"))
			if err != nil {
				return err
			}
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	cmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "Display name for the example resource.")
	cmd.Flags().String(flag.S3URI, "", "The S3 URI.")
	cmd.Flags().String(flag.S3AccessKeyID, "", "The S3 access key ID.")
	cmd.Flags().String(flag.S3SecretAccessKey, "", "The S3 secret access key.")
	return cmd
}
