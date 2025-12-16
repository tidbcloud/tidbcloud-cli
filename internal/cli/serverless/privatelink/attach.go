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
	"strings"

	"github.com/AlecAivazis/survey/v2"
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

type AttachDomainOpts struct {
	interactive bool
}

func (o AttachDomainOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkConnectionID,
		flag.PLCAttachDomainType,
		flag.PLCAttachDomainUniqueName,
	}
}

func (o AttachDomainOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkConnectionID,
		flag.PLCAttachDomainType,
	}
}

func (o *AttachDomainOpts) MarkInteractive(cmd *cobra.Command) error {
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

func AttachDomainCmd(h *internal.Helper) *cobra.Command {
	opts := &AttachDomainOpts{interactive: true}
	cmd := &cobra.Command{
		Use:   "attach-domains",
		Short: "Attach domains to a private link connection",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Attach domain (interactive):
  $ %[1]s serverless private-link-connection attach-domains

  Attach domain (non-interactive):
  $ %[1]s serverless private-link-connection attach-domains -c <cluster-id> --private-link-connection-id <plc-id> --type <type> --unique-name <unique-name>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, plcID string
			var domainType plapi.PrivateLinkConnectionDomainTypeEnum
			var uniqueName string
			var dryRun bool
			dryRun, err = cmd.Flags().GetBool(flag.DryRun)
			if err != nil {
				return errors.Trace(err)
			}
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

				privatelink, err := cloud.GetSelectedPrivateLinkConnection(ctx, cluster.ID, int32(h.QueryPageSize), d)
				if err != nil {
					return err
				}
				plcID = privatelink.ID

				domainType, err = GetSelectedPLCAttachDomainType()
				if err != nil {
					return err
				}

				switch domainType {
				case plapi.PRIVATELINKCONNECTIONDOMAINTYPEENUM_CONFLUENT:
					if err := survey.AskOne(&survey.Input{Message: "Domain unique name:"}, &uniqueName); err != nil {
						return err
					}
				case plapi.PRIVATELINKCONNECTIONDOMAINTYPEENUM_TIDBCLOUD_MANAGED:
					if !dryRun {
						if err := survey.AskOne(&survey.Input{Message: "Domain unique name:"}, &uniqueName); err != nil {
							return err
						}
					}
				default:
					return errors.New("invalid domain type")
				}
			} else {
				var err error
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				plcID, err = cmd.Flags().GetString(flag.PrivateLinkConnectionID)
				if err != nil {
					return errors.Trace(err)
				}
				domainTypeStr, err := cmd.Flags().GetString(flag.PLCAttachDomainType)
				if err != nil {
					return errors.Trace(err)
				}
				domainType = plapi.PrivateLinkConnectionDomainTypeEnum(domainTypeStr)
				uniqueName, err = cmd.Flags().GetString(flag.PLCAttachDomainUniqueName)
				if err != nil {
					return errors.Trace(err)
				}

				if domainType == "" {
					return errors.New("domain type is required")
				}
			}

			body := &plapi.PrivateLinkConnectionServiceAttachDomainsBody{
				AttachDomain: plapi.AttachDomain{
					Type:       domainType,
					UniqueName: &uniqueName,
				},
				ValidateOnly: &dryRun,
			}

			resp, err := d.AttachPrivateLinkDomains(ctx, clusterID, plcID, body)
			if err != nil {
				return errors.Trace(err)
			}
			if !dryRun {
				return output.PrintJson(h.IOStreams.Out, resp)
			}
			domains := make([]string, 0, len(resp.Domains))
			for _, domain := range resp.Domains {
				domains = append(domains, *domain.Name)
			}
			fmt.Fprintf(h.IOStreams.Out, "unique name %s:\n%s\n",
				color.BlueString("%v", *resp.UniqueName),
				color.GreenString("%v", strings.Join(domains, "\n")))
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	cmd.Flags().String(flag.PrivateLinkConnectionID, "", "The private link connection ID.")
	cmd.Flags().String(flag.PLCAttachDomainType, "", fmt.Sprintf("The type of domain to attach, one of: %v", plapi.AllowedPrivateLinkConnectionDomainTypeEnumEnumValues))
	cmd.Flags().String(flag.PLCAttachDomainUniqueName, "", "The unique name of the domain to attach, you can use --dry-run to generate the unique name when attaching a TiDB Cloud managed domain.")
	cmd.Flags().Bool(flag.DryRun, false, "Set dry run mode to only show generated domains without attaching them.")

	return cmd
}
