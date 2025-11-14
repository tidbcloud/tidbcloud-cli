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

package privatelink

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	plapi "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"
)

type CreateOpts struct {
	interactive bool
}

func (o CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.Spec, // JSON body
	}
}

func (o *CreateOpts) MarkInteractive(cmd *cobra.Command) error {
	o.interactive = true
	for _, fn := range o.NonInteractiveFlags() {
		f := cmd.Flags().Lookup(fn)
		if f != nil && f.Changed {
			o.interactive = false
		}
	}
	if !o.interactive {
		for _, fn := range o.NonInteractiveFlags() {
			if err := cmd.MarkFlagRequired(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := &CreateOpts{interactive: true}

	var spec string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a private link connection",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a private link connection (interactive):
  $ %[1]s serverless private-link-connection create

  Create a private link connection (non-interactive):
  $ %[1]s serverless private-link-connection create -c <cluster-id> \
    --spec '<json-spec>'
  Hint: Run "%[1]s serverless private-link-connection get-zones -c <cluster-id>" to get supported account and availability zones.`, config.CliName),
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

				prompt := &survey.Editor{
					Message:       "Enter the JSON spec of private link connection",
					FileName:      "*.json",
					HideDefault:   true,
					AppendDefault: true,
				}
				if err := survey.AskOne(prompt, &spec); err != nil {
					if err == terminal.InterruptErr {
						return internal.InterruptError
					}
					return err
				}
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				spec, err = cmd.Flags().GetString(flag.Spec)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if spec == "" {
				return errors.New("spec is required")
			}

			var body plapi.PrivateLinkConnectionServiceCreatePrivateLinkConnectionBody
			if err := json.Unmarshal([]byte(spec), &body); err != nil {
				return errors.Errorf("invalid spec JSON: %v", err)
			}

			res, err := d.CreatePrivateLinkConnection(ctx, clusterID, &body)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Private link connection created"))
			_ = output.PrintJson(h.IOStreams.Out, res)
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	cmd.Flags().String(flag.Spec, "", "The JSON spec for the private link connection.")
	return cmd
}
