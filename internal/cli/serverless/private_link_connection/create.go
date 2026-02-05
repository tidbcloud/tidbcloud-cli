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

package private_link_connection

import (
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type CreateOpts struct {
	interactive bool
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
		flag.PrivateLinkConnectionType,
		flag.AWSEndpointServiceName,
		flag.AWSEndpointServiceRegion,
		flag.AlicloudEndpointServiceName,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
		flag.PrivateLinkConnectionType,
	}
}

func (c *CreateOpts) MarkInteractive(cmd *cobra.Command) error {
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

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a private link connection for dataflow",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a private link connection (interactive):
  $ %[1]s serverless private-link-connection create

  Create a private link connection which connect to alicloud endpoint service (non-interactive):
  $ %[1]s serverless private-link-connection create -c <cluster-id> --display-name <name> --type ALICLOUD_ENDPOINT_SERVICE --alicloud.endpoint-service-name <name>

  Create a private link connection which connect to aws endpoint service (non-interactive):
  $ %[1]s serverless private-link-connection create -c <cluster-id> --display-name <name> --type AWS_ENDPOINT_SERVICE --aws.endpoint-service-name <name>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.MarkInteractive(cmd); err != nil {
				return errors.Trace(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID string
			var displayName string
			var endpointServiceName string
			var endpointServiceRegion string
			var alicloudEndpointServiceName string
			var plcType privatelink.PrivateLinkConnectionTypeEnum
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

				plcType, err = GetPrivateLinkConnectionType()
				if err != nil {
					return err
				}

				switch plcType {
				case privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE:
					inputModel, err := GetCreateAWSInput()
					if err != nil {
						return err
					}
					displayName = inputModel.(ui.TextInputModel).Inputs[createAWSField[flag.DisplayName]].Value()
					endpointServiceName = inputModel.(ui.TextInputModel).Inputs[createAWSField[flag.AWSEndpointServiceName]].Value()
					endpointServiceRegion = inputModel.(ui.TextInputModel).Inputs[createAWSField[flag.AWSEndpointServiceRegion]].Value()
				case privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE:
					inputModel, err := GetCreateAlicloudInput()
					if err != nil {
						return err
					}
					displayName = inputModel.(ui.TextInputModel).Inputs[createAlicloudField[flag.DisplayName]].Value()
					alicloudEndpointServiceName = inputModel.(ui.TextInputModel).Inputs[createAlicloudField[flag.AlicloudEndpointServiceName]].Value()
				default:
					return fmt.Errorf("unsupported private link connection type: %s", plcType)
				}
			} else {
				displayName, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				typeValue, err := cmd.Flags().GetString(flag.PrivateLinkConnectionType)
				if err != nil {
					return errors.Trace(err)
				}
				plcType, err = normalizePrivateLinkConnectionType(typeValue)
				if err != nil {
					return err
				}
				endpointServiceName, err = cmd.Flags().GetString(flag.AWSEndpointServiceName)
				if err != nil {
					return errors.Trace(err)
				}
				endpointServiceRegion, err = cmd.Flags().GetString(flag.AWSEndpointServiceRegion)
				if err != nil {
					return errors.Trace(err)
				}
				alicloudEndpointServiceName, err = cmd.Flags().GetString(flag.AlicloudEndpointServiceName)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if displayName == "" {
				return fmt.Errorf("display name is required")
			}
			switch plcType {
			case privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE:
				if endpointServiceName == "" {
					return fmt.Errorf("aws endpoint service name is required")
				}
			case privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE:
				if alicloudEndpointServiceName == "" {
					return fmt.Errorf("alicloud endpoint service name is required")
				}
			default:
				return fmt.Errorf("unsupported private link connection type: %s", plcType)
			}

			plc := privatelink.NewPrivateLinkConnection(clusterID, displayName, plcType)
			switch plcType {
			case privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE:
				awsService := privatelink.NewAwsEndpointService(endpointServiceName)
				if endpointServiceRegion != "" {
					awsService.SetRegion(endpointServiceRegion)
				}
				plc.SetAwsEndpointService(*awsService)
			case privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE:
				alicloudService := privatelink.NewAlicloudEndpointService(alicloudEndpointServiceName)
				plc.SetAlicloudEndpointService(*alicloudService)
			}

			body := privatelink.NewPrivateLinkConnectionServiceCreatePrivateLinkConnectionBody(*plc)
			resp, err := d.CreatePrivateLinkConnection(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}
			if resp != nil && resp.PrivateLinkConnectionId != nil {
				fmt.Fprintln(h.IOStreams.Out, color.GreenString("private link connection %s created.", *resp.PrivateLinkConnectionId))
				return nil
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("private link connection created."))
			return nil
		},
	}

	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	createCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "Display name for the private link connection.")
	createCmd.Flags().String(flag.PrivateLinkConnectionType, "", "Type of the private link connection, one of [\"AWS_ENDPOINT_SERVICE\" \"ALICLOUD_ENDPOINT_SERVICE\"]")
	createCmd.Flags().String(flag.AWSEndpointServiceName, "", "AWS endpoint service name.")
	createCmd.Flags().String(flag.AWSEndpointServiceRegion, "", "AWS endpoint service region.")
	createCmd.Flags().String(flag.AlicloudEndpointServiceName, "", "Alicloud endpoint service name.")

	return createCmd
}
