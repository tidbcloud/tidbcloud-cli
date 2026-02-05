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

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
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
		flag.AwsEndpointServiceName,
		flag.AwsEndpointServiceRegion,
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
		Short: "Create a private link connection",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a private link connection in interactive mode:
  $ %[1]s serverless private-link-connection create

  Create a private link connection in non-interactive mode:
  $ %[1]s serverless private-link-connection create -c <cluster-id> --display-name <display-name> --type AWS_ENDPOINT_SERVICE --aws-endpoint-service-name <service-name> --aws-endpoint-service-region <region>`, config.CliName),
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
			var connectionType privatelink.PrivateLinkConnectionTypeEnum
			var awsEndpointServiceName string
			var awsEndpointServiceRegion string
			var alicloudEndpointServiceName string

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

				types := []string{
					string(privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE),
					string(privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE),
				}
				selectedType, err := cloud.GetSelectedField(types, "Choose the private link connection type:")
				if err != nil {
					return err
				}
				connectionType = privatelink.PrivateLinkConnectionTypeEnum(selectedType)

				var inputs []string
				inputDescription := map[string]string{
					flag.DisplayName: "Display name",
				}
				if connectionType == privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE {
					inputs = []string{flag.DisplayName, flag.AwsEndpointServiceName, flag.AwsEndpointServiceRegion}
					inputDescription[flag.AwsEndpointServiceName] = "AWS endpoint service name"
					inputDescription[flag.AwsEndpointServiceRegion] = "AWS endpoint service region (optional)"
				} else {
					inputs = []string{flag.DisplayName, flag.AlicloudEndpointServiceName}
					inputDescription[flag.AlicloudEndpointServiceName] = "Alicloud endpoint service name"
				}

				inputModel, err := ui.InitialInputModel(inputs, inputDescription)
				if err != nil {
					return err
				}

				displayName = inputModel.Inputs[0].Value()
				if connectionType == privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE {
					awsEndpointServiceName = inputModel.Inputs[1].Value()
					if len(inputModel.Inputs) > 2 {
						awsEndpointServiceRegion = inputModel.Inputs[2].Value()
					}
				} else {
					alicloudEndpointServiceName = inputModel.Inputs[1].Value()
				}
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				displayName, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}
				typeValue, err := cmd.Flags().GetString(flag.PrivateLinkConnectionType)
				if err != nil {
					return errors.Trace(err)
				}
				connectionType = privatelink.PrivateLinkConnectionTypeEnum(typeValue)
				if !connectionType.IsValid() {
					return fmt.Errorf("invalid private link connection type: %s", typeValue)
				}
				awsEndpointServiceName, err = cmd.Flags().GetString(flag.AwsEndpointServiceName)
				if err != nil {
					return errors.Trace(err)
				}
				awsEndpointServiceRegion, err = cmd.Flags().GetString(flag.AwsEndpointServiceRegion)
				if err != nil {
					return errors.Trace(err)
				}
				alicloudEndpointServiceName, err = cmd.Flags().GetString(flag.AlicloudEndpointServiceName)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if displayName == "" {
				return errors.New("display name is required")
			}

			connection := privatelink.NewPrivateLinkConnection(clusterID, displayName, connectionType)
			switch connectionType {
			case privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE:
				if awsEndpointServiceName == "" {
					return errors.New("AWS endpoint service name is required")
				}
				awsEndpointService := privatelink.NewAwsEndpointService(awsEndpointServiceName)
				if awsEndpointServiceRegion != "" {
					awsEndpointService.SetRegion(awsEndpointServiceRegion)
				}
				connection.AwsEndpointService = awsEndpointService
			case privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE:
				if alicloudEndpointServiceName == "" {
					return errors.New("Alicloud endpoint service name is required")
				}
				connection.AlicloudEndpointService = privatelink.NewAlicloudEndpointService(alicloudEndpointServiceName)
			default:
				return fmt.Errorf("unsupported private link connection type: %s", connectionType)
			}

			body := &privatelink.PrivateLinkConnectionServiceCreatePrivateLinkConnectionBody{
				PrivateLinkConnection: *connection,
			}
			created, err := d.CreatePrivateLinkConnection(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}
			if format == output.JsonFormat {
				return errors.Trace(output.PrintJson(h.IOStreams.Out, created))
			}

			connectionID := ""
			if created.PrivateLinkConnectionId != nil {
				connectionID = *created.PrivateLinkConnectionId
			}
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("Private link connection %s is created", connectionID))
			return err
		},
	}

	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	createCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The display name of the private link connection.")
	createCmd.Flags().String(flag.PrivateLinkConnectionType, "", "The private link connection type. One of [\"AWS_ENDPOINT_SERVICE\", \"ALICLOUD_ENDPOINT_SERVICE\"].")
	createCmd.Flags().String(flag.AwsEndpointServiceName, "", "The AWS endpoint service name.")
	createCmd.Flags().String(flag.AwsEndpointServiceRegion, "", "The AWS endpoint service region.")
	createCmd.Flags().String(flag.AlicloudEndpointServiceName, "", "The Alicloud endpoint service name.")
	createCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)

	return createCmd
}
