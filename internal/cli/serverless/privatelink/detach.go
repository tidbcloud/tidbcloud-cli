package privatelink

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	plapi "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"
)

type DetachDomainOpts struct {
	interactive bool
}

func (o DetachDomainOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.PrivateLinkConnectionID,
		flag.PLCAttachDomainID,
	}
}

func (o *DetachDomainOpts) MarkInteractive(cmd *cobra.Command) error {
	o.interactive = true
	for _, fn := range o.NonInteractiveFlags() {
		if f := cmd.Flags().Lookup(fn); f != nil && f.Changed {
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

func DetachDomainCmd(h *internal.Helper) *cobra.Command {
	opts := &DetachDomainOpts{interactive: true}

	cmd := &cobra.Command{
		Use:   "detach",
		Short: "Detach domains from a private link connection",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Detach domains (interactive):
  $ %[1]s serverless private-link-connection detach

  Detach domains (non-interactive):
  $ %[1]s serverless private-link-connection detach -c <cluster-id> --private-link-connection-id <plc-id> --plc-attach-domain-id <attach-id>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.MarkInteractive(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			var clusterID, plcID, attachID string
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
				attach, err := cloud.GetSelectedAttachDomain(ctx, cluster.ID, privatelink.ID, d)
				if err != nil {
					return err
				}
				attachID = attach.ID
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
				attachID, err = cmd.Flags().GetString(flag.PLCAttachDomainID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			body := &plapi.PrivateLinkConnectionServiceDetachDomainsBody{
				AttachDomainId: attachID,
			}
			resp, err := d.DetachPrivateLinkDomains(ctx, clusterID, plcID, body)
			if err != nil {
				return errors.Trace(err)
			}
			domains := make([]string, 0, len(resp.Domains))
			for _, domain := range resp.Domains {
				domains = append(domains, *domain.Name)
			}
			fmt.Fprintf(h.IOStreams.Out, "Successfully detached domains:\n%s\n", color.GreenString("%v", strings.Join(domains, "\n")))
			return nil
		},
	}

	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID.")
	cmd.Flags().String(flag.PrivateLinkConnectionID, "", "The private link connection ID.")
	cmd.Flags().String(flag.PLCAttachDomainID, "", "The private link connection attach domain ID.")
	return cmd
}
