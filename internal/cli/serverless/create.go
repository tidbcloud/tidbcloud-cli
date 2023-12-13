// Copyright 2022 PingCAP, Inc.
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

package serverless

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	serverlessApi "tidbcloud-cli/pkg/tidbcloud/serverless/client/serverless_service"
	serverlessModel "tidbcloud-cli/pkg/tidbcloud/serverless/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type createClusterField int

const (
	clusterNameIdx createClusterField = iota
	spendingLimitIdx
)

const (
	serverlessType = "SERVERLESS"
	WaitInterval   = 5 * time.Second
	WaitTimeout    = 2 * time.Minute
)

type CreateOpts struct {
	serverlessProviders []*serverlessModel.TidbCloudOpenApiserverlessv1beta1Region
	interactive         bool
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterName,
		flag.ClusterType,
		flag.CloudProvider,
		flag.Region,
		flag.ProjectID,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterName,
		flag.CloudProvider,
		flag.Region,
	}
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:         "create",
		Short:       "Create a serverless cluster",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Create a serverless cluster in interactive mode:
  $ %[1]s serverless create

  Create a serverless cluster of the default ptoject in non-interactive mode:
  $ %[1]s serverless create --cluster-name <cluster-name> --cloud-provider <cloud-provider> -r <region>

  Create a serverless cluster in non-interactive mode:
  $ %[1]s serverless create --project-id <project-id> --cluster-name <cluster-name> --cloud-provider <cloud-provider> -r <region>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// mark required flags in non-interactive mode
			if !opts.interactive {
				for _, fn := range opts.RequiredFlags() {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterName string
			var cloudProvider string
			var region string
			var projectID string
			var spendingLimitMonthly int32
			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				regions, err := d.ListProviderRegions(serverlessApi.NewServerlessServiceListRegionsParams())
				if err != nil {
					return errors.Trace(err)
				}
				opts.serverlessProviders = regions.Payload.Regions

				// distinct cloud providers
				providers := hashset.New()
				for _, provider := range opts.serverlessProviders {
					providers.Add(string(*provider.Provider))
				}
				model, err := ui.InitialSelectModel(providers.Values(), "Choose the cloud provider:")
				if err != nil {
					return err
				}
				p := tea.NewProgram(model)
				providerModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if m, _ := providerModel.(ui.SelectModel); m.Interrupted {
					return util.InterruptError
				}
				cloudProvider = providerModel.(ui.SelectModel).Choices[providerModel.(ui.SelectModel).Selected].(string)

				// filter out regions for the selected cloud provider
				regionSet := hashset.New()
				for _, provider := range opts.serverlessProviders {
					if string(*provider.Provider) == cloudProvider {
						regionSet.Add(cloud.Region{
							Name:        *provider.Name,
							DisplayName: provider.DisplayName,
							Provider:    string(*provider.Provider),
						})
					}
				}
				model, err = ui.InitialSelectModel(regionSet.Values(), "Choose the cloud region:")
				if err != nil {
					return err
				}
				p = tea.NewProgram(model)
				regionModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if m, _ := regionModel.(ui.SelectModel); m.Interrupted {
					return util.InterruptError
				}
				region = regionModel.(ui.SelectModel).Choices[regionModel.(ui.SelectModel).Selected].(cloud.Region).Name

				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				projectID = project.ID

				// variables for input
				p = tea.NewProgram(initialCreateInputModel())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				clusterName = inputModel.(ui.TextInputModel).Inputs[clusterNameIdx].Value()
				spendingLimitString := inputModel.(ui.TextInputModel).Inputs[spendingLimitIdx].Value()
				if len(clusterName) == 0 {
					return errors.New("cluster name is required")
				}
				if len(spendingLimitString) == 0 {
					spendingLimitMonthly = 0
				} else {
					s, err := strconv.Atoi(spendingLimitString)
					if err != nil {
						return errors.New("monthly spending limit should be int type")
					}
					if s > math.MaxInt32 || s < math.MinInt32 {
						return errors.New("monthly spending limit out of range")
					}
					spendingLimitMonthly = int32(s) //nolint:gosec // will not overflow
				}
			} else {
				// non-interactive mode, get values from flags
				var err error
				clusterName, err = cmd.Flags().GetString(flag.ClusterName)
				if err != nil {
					return errors.Trace(err)
				}
				cloudProvider, err = cmd.Flags().GetString(flag.CloudProvider)
				if err != nil {
					return errors.Trace(err)
				}
				region, err = cmd.Flags().GetString(flag.Region)
				if err != nil {
					return errors.Trace(err)
				}
				spendingLimitMonthly, err = cmd.Flags().GetInt32(flag.SpendingLimitMonthly)
				if err != nil {
					return errors.Trace(err)
				}

				// generate region name
				region = fmt.Sprintf("regions/%s-%s", strings.ToLower(cloudProvider), strings.ToLower(region))
				projectID, err = cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return errors.Trace(err)
				}
			}

			cmd.Annotations[telemetry.ProjectID] = projectID

			// check clusterName
			err = checkClusterName(clusterName)
			if err != nil {
				return errors.Trace(err)
			}

			v1Cluster := &serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster{
				DisplayName: &clusterName,
				Region: &serverlessModel.TidbCloudOpenApiserverlessv1beta1Region{
					Name: &region,
				},
			}
			// optional fields
			if projectID != "" {
				v1Cluster.Labels = map[string]string{"tidb.cloud/project": projectID}
			}
			if spendingLimitMonthly != 0 {
				v1Cluster.SpendingLimit = &serverlessModel.ClusterSpendingLimit{
					Monthly: spendingLimitMonthly,
				}
			}

			if h.IOStreams.CanPrompt {
				err := CreateAndSpinnerWait(ctx, d, v1Cluster, h)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				err := CreateAndWaitReady(h, d, v1Cluster)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	createCmd.Flags().String(flag.ClusterName, "", "Name of the cluster to de created")
	createCmd.Flags().String(flag.CloudProvider, "", "Cloud provider, one of [\"AWS\"]")
	createCmd.Flags().StringP(flag.Region, flag.RegionShort, "", "Cloud region")
	createCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The ID of the project, in which the cluster will be created (optional: default \"default project\")")
	createCmd.Flags().Int32(flag.SpendingLimitMonthly, 0, "The monthly spending limit of the cluster (optional)")
	return createCmd
}

func CreateAndWaitReady(h *internal.Helper, d cloud.TiDBCloudClient, v1Cluster *serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster) error {
	createClusterResult, err := d.CreateCluster(serverlessApi.NewServerlessServiceCreateClusterParams().WithCluster(v1Cluster))
	if err != nil {
		return errors.Trace(err)
	}
	newClusterID := createClusterResult.GetPayload().ClusterID

	fmt.Fprintln(h.IOStreams.Out, "... Waiting for cluster to be ready")
	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()
	timer := time.After(WaitTimeout)
	for {
		select {
		case <-timer:
			return errors.New(fmt.Sprintf("Timeout waiting for cluster %s to be ready, please check status on dashboard.", newClusterID))
		case <-ticker.C:
			clusterResult, err := d.GetCluster(serverlessApi.NewServerlessServiceGetClusterParams().
				WithClusterID(newClusterID))
			if err != nil {
				return errors.Trace(err)
			}
			s := *clusterResult.GetPayload().State
			if s == "ACTIVE" {
				fmt.Fprint(h.IOStreams.Out, color.GreenString("Cluster %s is ready.", newClusterID))
				return nil
			}
		}
	}
}

func CreateAndSpinnerWait(ctx context.Context, d cloud.TiDBCloudClient, v1Cluster *serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster, h *internal.Helper) error {
	// use spinner to indicate that the cluster is being created
	task := func() tea.Msg {
		createClusterResult, err := d.CreateCluster(serverlessApi.NewServerlessServiceCreateClusterParams().WithCluster(v1Cluster))
		if err != nil {
			return errors.Trace(err)
		}
		newClusterID := createClusterResult.GetPayload().ClusterID

		ticker := time.NewTicker(WaitInterval)
		defer ticker.Stop()
		timer := time.After(WaitTimeout)
		for {
			select {
			case <-timer:
				return ui.Result(fmt.Sprintf("Timeout waiting for cluster %s to be ready, please check status on dashboard.", newClusterID))
			case <-ticker.C:
				clusterResult, err := d.GetCluster(serverlessApi.NewServerlessServiceGetClusterParams().
					WithClusterID(newClusterID))
				if err != nil {
					return errors.Trace(err)
				}
				s := *clusterResult.GetPayload().State
				if s == "ACTIVE" {
					return ui.Result(fmt.Sprintf("Cluster %s is ready.", newClusterID))
				}
			case <-ctx.Done():
				return util.InterruptError
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Waiting for cluster to be ready"))
	createModel, err := p.StartReturningModel()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := createModel.(ui.SpinnerModel); m.Interrupted {
		return util.InterruptError
	}
	if m, _ := createModel.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	} else {
		fmt.Fprintln(h.IOStreams.Out, color.GreenString(m.Output))
	}
	return nil
}

func initialCreateInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = config.CursorStyle
		t.CharLimit = 64
		f := createClusterField(i)

		switch f {
		case clusterNameIdx:
			t.Placeholder = "Cluster Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case spendingLimitIdx:
			t.Placeholder = "Spending Limit Monthly($), e.g., 10. Skip it by press 0 or enter"
		}
		m.Inputs[i] = t
	}

	return m
}

func checkClusterName(name string) error {
	if len(name) < 4 || len(name) > 64 || !isNumber(name[0]) && !isLetter(name[0]) || !isNumber(name[len(name)-1]) && !isLetter(name[len(name)-1]) {
		return errors.New("Cluster name must be 4~64 characters that can only include numbers, lowercase or uppercase letters, and hyphens. The first and last character must be a letter or number.")
	}
	return nil
}
func isNumber(s byte) bool {
	return s >= '0' && s <= '9'
}

func isLetter(s byte) bool {
	return s >= 'a' && s <= 'z' || s >= 'A' && s <= 'Z'
}
