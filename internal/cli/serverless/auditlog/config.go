// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auditlog

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
)

type ConfigOpts struct {
	interactive bool
}

func (c ConfigOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.AuditLogUnRedacted,
	}
}

type mutableField string

const (
	Unredacted mutableField = "unredacted"
	Enabled    mutableField = "enabled"
)

var mutableFields = []string{
	string(Unredacted),
	string(Enabled),
}

func (c *ConfigOpts) MarkInteractive(cmd *cobra.Command) error {
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
		err := cmd.MarkFlagRequired(flag.ClusterID)
		if err != nil {
			return err
		}
		cmd.MarkFlagsOneRequired(flag.AuditLogUnRedacted, flag.Enabled)
	}
	return nil
}

func ConfigCmd(h *internal.Helper) *cobra.Command {
	opts := ConfigOpts{
		interactive: true,
	}

	var configCmd = &cobra.Command{
		Use:         "config",
		Short:       "Configure the database audit logging",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Conigure the database audit logging in interactive mode:
  $ %[1]s serverless audit-log config

  Unredacted the database audit logging in non-interactive mode:
  $ %[1]s serverless audit-log config -c <cluster-id> --unredacted

  Enable the database audit logging in non-interactive mode:
  $ %[1]s serverless audit-log config -c <cluster-id> --enabled

  Disable the database audit logging in non-interactive mode:
  $ %[1]s serverless audit-log config -c <cluster-id> --enabled=false`, config.CliName),
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
			var unredacted, enabled *bool
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				selectedCluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = selectedCluster.ID

				fieldName, err := cloud.GetSelectedField(mutableFields, "Choose one field to config:")
				if err != nil {
					return err
				}

				switch fieldName {
				case string(Unredacted):
					prompt := &survey.Confirm{
						Message: "unredacted the database audit logging?",
						Default: false,
					}
					var confirm bool
					err = survey.AskOne(prompt, &confirm)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					if confirm {
						unredacted = aws.Bool(true)
					} else {
						unredacted = aws.Bool(false)
					}
				case string(Enabled):
					prompt := &survey.Confirm{
						Message: "enable the database audit logging?",
						Default: false,
					}
					var confirm bool
					err = survey.AskOne(prompt, &confirm)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					if confirm {
						enabled = aws.Bool(true)
					} else {
						enabled = aws.Bool(false)
					}
				}
			} else {
				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID = cID
				if cmd.Flags().Changed(flag.AuditLogUnRedacted) {
					u, err := cmd.Flags().GetBool(flag.AuditLogUnRedacted)
					if err != nil {
						return errors.Trace(err)
					}
					unredacted = &u
				}
				if cmd.Flags().Changed(flag.Enabled) {
					u, err := cmd.Flags().GetBool(flag.Enabled)
					if err != nil {
						return errors.Trace(err)
					}
					enabled = &u
				}
			}

			body := &cluster.V1beta1ServerlessServicePartialUpdateClusterBody{
				Cluster: &cluster.RequiredTheClusterToBeUpdated{
					AuditLogConfig: &cluster.V1beta1ClusterAuditLogConfig{},
				},
				UpdateMask: "auditLogConfig",
			}
			if unredacted != nil {
				body.Cluster.AuditLogConfig.Unredacted = *cluster.NewNullableBool(unredacted)
			}
			if enabled != nil {
				body.Cluster.AuditLogConfig.Enabled = *cluster.NewNullableBool(enabled)
			}
			_, err = d.PartialUpdateCluster(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("configure cluster %s database audit logging succeed", clusterID)))
			return nil
		},
	}

	configCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster to be updated.")
	configCmd.Flags().Bool(flag.AuditLogUnRedacted, false, "unredacted or redacted the database audit logging.")
	configCmd.Flags().Bool(flag.Enabled, false, "enable or disable the database audit logging.")
	return configCmd
}
