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
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	//"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"
)

type CreateOpts struct {
	interactive bool
}

func (o CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.DisplayName,
		flag.PrivateLinkConnectionType,
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

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a private link connection",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Create a private link connection (interactive):
  $ %[1]s serverless private-link-connection create

  Create a private link connection which connect to alicloud endpoint service (non-interactive):
  $ %[1]s serverless private-link-connection create -c <cluster-id> --display-name <name> --type ALICLOUD_ENDPOINT_SERVICE --alicloud.endpoint-service.name <name>

  Create a private link connection which connect to aws endpoint service (non-interactive):
  $ %[1]s serverless private-link-connection create -c <cluster-id> --display-name <name> --type AWS_ENDPOINT_SERVICE --aws.endpoint-service.name <name>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, displayName string
			var connectionType privatelink.PrivateLinkConnectionTypeEnum
			var awsEndpointServiceName, AWSEndpointServiceRegion string
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

				if err := survey.AskOne(
					&survey.Input{Message: DisplayNamePrompt},
					&displayName,
					survey.WithValidator(survey.Required),
				); err != nil {
					return err
				}
				if displayName == "" {
					return errors.New("display name is required")
				}

				connectionType, err = GetSelectedPrivateLinkConnectionType()
				if err != nil {
					return err
				}
				switch connectionType {
				case privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE:
					if err = survey.AskOne(&survey.Input{Message: AWSEndpointServiceNamePrompt}, &awsEndpointServiceName, survey.WithValidator(survey.Required)); err != nil {
						return err
					}
					crossRegion := false
					if err = survey.AskOne(&survey.Confirm{Message: AWSEndpointServiceRegionConfirmPrompt, Default: false}, &crossRegion); err != nil {
						return err
					}
					if crossRegion {
						if err := survey.AskOne(&survey.Input{Message: AWSEndpointServiceRegionPrompt}, &AWSEndpointServiceRegion, survey.WithValidator(survey.Required)); err != nil {
							return err
						}
					}
				case privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE:
					if err := survey.AskOne(&survey.Input{Message: AlicloudEndpointServiceNamePrompt}, &alicloudEndpointServiceName, survey.WithValidator(survey.Required)); err != nil {
						return err
					}
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
				connectionTypeStr, err := cmd.Flags().GetString(flag.PrivateLinkConnectionType)
				if err != nil {
					return errors.Trace(err)
				}
				connectionType = privatelink.PrivateLinkConnectionTypeEnum(connectionTypeStr)

				awsEndpointServiceName, err = cmd.Flags().GetString(flag.AWSEndpointServiceName)
				if err != nil {
					return errors.Trace(err)
				}
				alicloudEndpointServiceName, err = cmd.Flags().GetString(flag.AlicloudEndpointServiceName)
				if err != nil {
					return errors.Trace(err)
				}
				AWSEndpointServiceRegion, err = cmd.Flags().GetString(flag.AWSEndpointServiceRegion)
				if err != nil {
					return errors.Trace(err)
				}
			}

			// build request body
			body := privatelink.PrivateLinkConnection{
				ClusterId: clusterID,
			}

			if displayName == "" {
				return errors.New("display name is required")
			}
			body.DisplayName = displayName

			switch connectionType {
			case privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE:
				if awsEndpointServiceName == "" {
					return errors.New("aws endpoint service name is required for aws endpoint service type")
				}
				body.Type = privatelink.PRIVATELINKCONNECTIONTYPEENUM_AWS_ENDPOINT_SERVICE
				var region *string
				if AWSEndpointServiceRegion != "" {
					region = &AWSEndpointServiceRegion
				}
				body.AwsEndpointService = &privatelink.AwsEndpointService{
					Name:   awsEndpointServiceName,
					Region: region,
				}
			case privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE:
				if alicloudEndpointServiceName == "" {
					return errors.New("alicloud endpoint service name is required for alicloud endpoint service type")
				}
				body.Type = privatelink.PRIVATELINKCONNECTIONTYPEENUM_ALICLOUD_ENDPOINT_SERVICE
				body.AlicloudEndpointService = &privatelink.AlicloudEndpointService{
					Name: alicloudEndpointServiceName,
				}
			default:
				return errors.Errorf("unsupported private link connection type: %s", connectionType)
			}

			res, err := d.CreatePrivateLinkConnection(ctx,
				clusterID,
				&privatelink.PrivateLinkConnectionServiceCreatePrivateLinkConnectionBody{
					PrivateLinkConnection: body,
				},
			)
			if err != nil {
				return errors.Trace(err)
			}
			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("Private link connectio %s created", res.PrivateLinkConnectionId))
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	cmd.Flags().String(flag.DisplayName, "", "Display name for the private link connection.")
	cmd.Flags().String(flag.PrivateLinkConnectionType, "", fmt.Sprintf("Type of the private link connection, one of %q", privatelink.AllowedPrivateLinkConnectionTypeEnumEnumValues))
	cmd.Flags().String(flag.AWSEndpointServiceName, "", "AWS endpoint service name")
	cmd.Flags().String(flag.AlicloudEndpointServiceName, "", "Alicloud endpoint service name")
	cmd.Flags().String(flag.AWSEndpointServiceRegion, "", "AWS endpoint service region")
	return cmd
}
