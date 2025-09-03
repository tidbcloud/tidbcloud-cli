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

package serverless

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

var createClusterField = map[string]int{
	flag.DisplayName: 0,
}

const (
	WaitInterval = 5 * time.Second
	WaitTimeout  = 2 * time.Minute
)

type CreateOpts struct {
	serverlessProviders []cluster.Commonv1beta1Region
	interactive         bool
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.DisplayName,
		flag.Region,
		flag.ProjectID,
		flag.SpendingLimitMonthly,
		flag.Encryption,
		flag.PublicEndpointDisabled,
		flag.MinRCU,
		flag.MaxRCU,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.DisplayName,
		flag.Region,
	}
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:         "create",
		Short:       "Create a TiDB Cloud Serverless cluster",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Create a TiDB Cloud Serverless cluster in interactive mode:
  $ %[1]s serverless create

  Create a TiDB Cloud Serverless cluster of the default project in non-interactive mode:
  $ %[1]s serverless create --display-name <cluster-name> --region <region>

  Create a TiDB Cloud Serverless cluster in non-interactive mode:
  $ %[1]s serverless create --project-id <project-id> --display-name <cluster-name> --region <region>`,
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
			var minRcu, maxRcu int32
			var encryption bool
			var publicEndpointDisabled bool
			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				regions, err := d.ListProviderRegions(ctx)
				if err != nil {
					return errors.Trace(err)
				}
				opts.serverlessProviders = regions.Regions

				// distinct cloud providers
				providers := hashset.New()
				for _, provider := range opts.serverlessProviders {
					providers.Add(string(*provider.CloudProvider))
				}
				cloudProvider, err = GetProvider(providers)
				if err != nil {
					return errors.Trace(err)
				}

				// filter out regions for the selected cloud provider
				regionSet := hashset.New()
				for _, provider := range opts.serverlessProviders {
					if string(*provider.CloudProvider) == cloudProvider {
						regionSet.Add(cloud.Region{
							Name:        *provider.Name,
							DisplayName: *provider.DisplayName,
							Provider:    string(*provider.CloudProvider),
						})
					}
				}
				region, err = GetRegion(regionSet)
				if err != nil {
					return errors.Trace(err)
				}

				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				projectID = project.ID

				// variables for input
				p := tea.NewProgram(initialCreateInputModel())
				inputModel, err := p.Run()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				clusterName = inputModel.(ui.TextInputModel).Inputs[createClusterField[flag.DisplayName]].Value()
				// check clusterName
				err = checkClusterName(clusterName)
				if err != nil {
					return errors.Trace(err)
				}

				// advanced options: spending limit/capacity and enhanced encryption
				if cloudProvider == string(cluster.V1BETA1REGIONCLOUDPROVIDER_AWS) {
					var spendingLimitString string
					spendingLimitPrompt := &survey.Input{
						Message: "Set spending limit monthly in USD cents (Example: 10, default is 0)?",
						Default: "0",
					}
					err = survey.AskOne(spendingLimitPrompt, &spendingLimitString)
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					spendingLimitMonthly, err = getAndCheckNumber(spendingLimitString, "monthly spending limit")
					if err != nil {
						return errors.Trace(err)
					}
					if spendingLimitMonthly > 0 {
						encryptionPrompt := &survey.Confirm{
							Message: "Enable Enhanced Encryption at Rest?",
							Default: false,
						}
						err = survey.AskOne(encryptionPrompt, &encryption)
						if err != nil {
							if err == terminal.InterruptErr {
								return util.InterruptError
							} else {
								return err
							}
						}
					}
				}

				if cloudProvider == string(cluster.V1BETA1REGIONCLOUDPROVIDER_ALICLOUD) {
					clusterPlan, err := GetClusterPlan()
					if err != nil {
						return err
					}
					if clusterPlan == cluster.CLUSTERCLUSTERPLAN_ESSENTIAL {
						var minRcuString, maxRcuString string
						minRcuPrompt := &survey.Input{
							Message: "Set minimum RCU (default is 2000)?",
							Default: "2000",
						}
						err = survey.AskOne(minRcuPrompt, &minRcuString)
						if err != nil {
							if err == terminal.InterruptErr {
								return util.InterruptError
							} else {
								return err
							}
						}
						minRcu, err = getAndCheckNumber(minRcuString, "minimum RCU")
						if err != nil {
							return errors.Trace(err)
						}
						maxRcuPrompt := &survey.Input{
							Message: "Set maximum RCU (default is 4000)?",
							Default: "4000",
						}
						err = survey.AskOne(maxRcuPrompt, &maxRcuString)
						if err != nil {
							if err == terminal.InterruptErr {
								return util.InterruptError
							} else {
								return err
							}
						}
						maxRcu, err = getAndCheckNumber(maxRcuString, "maximum RCU")
						if err != nil {
							return errors.Trace(err)
						}
						err = checkCapacity(minRcu, maxRcu)
						if err != nil {
							return errors.Trace(err)
						}
						encryptionPrompt := &survey.Confirm{
							Message: "Enable Enhanced Encryption at Rest?",
							Default: false,
						}
						err = survey.AskOne(encryptionPrompt, &encryption)
						if err != nil {
							if err == terminal.InterruptErr {
								return util.InterruptError
							} else {
								return err
							}
						}
					}
					if clusterPlan == cluster.CLUSTERCLUSTERPLAN_STARTER {
						var spendingLimitString string
						spendingLimitPrompt := &survey.Input{
							Message: "Set spending limit monthly in USD cents (Example: 10, default is 0)?",
							Default: "0",
						}
						err = survey.AskOne(spendingLimitPrompt, &spendingLimitString)
						if err != nil {
							if err == terminal.InterruptErr {
								return util.InterruptError
							} else {
								return err
							}
						}
						spendingLimitMonthly, err = getAndCheckNumber(spendingLimitString, "monthly spending limit")
						if err != nil {
							return errors.Trace(err)
						}
					}
				}
			} else {
				// non-interactive mode, get values from flags
				var err error
				clusterName, err = cmd.Flags().GetString(flag.DisplayName)
				if err != nil {
					return errors.Trace(err)
				}
				region, err = cmd.Flags().GetString(flag.Region)
				if err != nil {
					return errors.Trace(err)
				}
				projectID, err = cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return errors.Trace(err)
				}
				encryption, err = cmd.Flags().GetBool(flag.Encryption)
				if err != nil {
					return errors.Trace(err)
				}
				publicEndpointDisabled, err = cmd.Flags().GetBool(flag.PublicEndpointDisabled)
				if err != nil {
					return errors.Trace(err)
				}
				spendingLimitMonthly, err = cmd.Flags().GetInt32(flag.SpendingLimitMonthly)
				if err != nil {
					return errors.Trace(err)
				}
				minRcu, err = cmd.Flags().GetInt32(flag.MinRCU)
				if err != nil {
					return errors.Trace(err)
				}
				maxRcu, err = cmd.Flags().GetInt32(flag.MaxRCU)
				if err != nil {
					return errors.Trace(err)
				}
				// check clusterName
				err = checkClusterName(clusterName)
				if err != nil {
					return errors.Trace(err)
				}
				if spendingLimitMonthly < 0 {
					return errors.New("Spending limit monthly should be non-negative")
				}
				if minRcu != 0 || maxRcu != 0 {
					err = checkCapacity(minRcu, maxRcu)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			cmd.Annotations[telemetry.ProjectID] = projectID

			v1Cluster := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{
				DisplayName: clusterName,
				Region: cluster.Commonv1beta1Region{
					Name: &region,
				},
			}
			// optional fields
			if projectID != "" {
				v1Cluster.Labels = &map[string]string{"tidb.cloud/project": projectID}
			}
			if spendingLimitMonthly != 0 {
				v1Cluster.SpendingLimit = &cluster.ClusterSpendingLimit{
					Monthly: &spendingLimitMonthly,
				}
			}
			if maxRcu != 0 || minRcu != 0 {
				v1Cluster.AutoScaling = &cluster.V1beta1ClusterAutoScaling{
					MinRcu: toInt64Ptr(minRcu),
					MaxRcu: toInt64Ptr(maxRcu),
				}
			}
			if encryption {
				v1Cluster.EncryptionConfig = &cluster.V1beta1ClusterEncryptionConfig{
					EnhancedEncryptionEnabled: &encryption,
				}
			}

			if publicEndpointDisabled {
				v1Cluster.Endpoints = &cluster.V1beta1ClusterEndpoints{
					Public: &cluster.EndpointsPublic{
						Disabled: &publicEndpointDisabled,
					},
				}
			}

			if h.IOStreams.CanPrompt {
				err := CreateAndSpinnerWait(ctx, h, d, v1Cluster)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				err := CreateAndWaitReady(ctx, h, d, v1Cluster)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	createCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "Display name of the cluster to de created.")
	createCmd.Flags().StringP(flag.Region, flag.RegionShort, "", "The name of cloud region. You can use \"ticloud serverless region\" to see all regions.")
	createCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The ID of the project, in which the cluster will be created. (default \"default project\")")
	createCmd.Flags().Int32(flag.SpendingLimitMonthly, 0, "Maximum monthly spending limit in USD cents.")
	createCmd.Flags().Bool(flag.Encryption, false, "Whether Enhanced Encryption at Rest is enabled.")
	createCmd.Flags().Bool(flag.PublicEndpointDisabled, false, "Whether the public endpoint is disabled.")
	createCmd.Flags().Int32(flag.MinRCU, 0, "Minimum RCU for the cluster, at least 2000.")
	createCmd.Flags().Int32(flag.MaxRCU, 0, "Maximum RCU for the cluster, at most 100000.")
	createCmd.MarkFlagsMutuallyExclusive(flag.SpendingLimitMonthly, flag.MinRCU)
	createCmd.MarkFlagsMutuallyExclusive(flag.SpendingLimitMonthly, flag.MaxRCU)
	createCmd.MarkFlagsRequiredTogether(flag.MinRCU, flag.MaxRCU)
	return createCmd
}

func CreateAndWaitReady(ctx context.Context, h *internal.Helper, d cloud.TiDBCloudClient, v1Cluster *cluster.TidbCloudOpenApiserverlessv1beta1Cluster) error {
	createClusterResult, err := d.CreateCluster(ctx, v1Cluster)
	if err != nil {
		return errors.Trace(err)
	}
	newClusterID := *createClusterResult.ClusterId

	fmt.Fprintln(h.IOStreams.Out, "... Waiting for cluster to be ready")
	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()
	timer := time.After(WaitTimeout)
	for {
		select {
		case <-timer:
			return errors.New(fmt.Sprintf("Timeout waiting for cluster %s to be ready, please check status on dashboard.", newClusterID))
		case <-ticker.C:
			clusterResult, err := d.GetCluster(ctx, newClusterID, cluster.CLUSTERSERVICEGETCLUSTERVIEWPARAMETER_BASIC)
			if err != nil {
				return errors.Trace(err)
			}
			if *clusterResult.State == cluster.COMMONV1BETA1CLUSTERSTATE_ACTIVE {
				fmt.Fprint(h.IOStreams.Out, color.GreenString("Cluster %s is ready.", newClusterID))
				return nil
			}
		}
	}
}

func CreateAndSpinnerWait(ctx context.Context, h *internal.Helper, d cloud.TiDBCloudClient, v1Cluster *cluster.TidbCloudOpenApiserverlessv1beta1Cluster) error {
	// use spinner to indicate that the cluster is being created
	task := func() tea.Msg {
		createClusterResult, err := d.CreateCluster(ctx, v1Cluster)
		if err != nil {
			return errors.Trace(err)
		}
		newClusterID := *createClusterResult.ClusterId

		ticker := time.NewTicker(WaitInterval)
		defer ticker.Stop()
		timer := time.After(WaitTimeout)
		for {
			select {
			case <-timer:
				return ui.Result(fmt.Sprintf("Timeout waiting for cluster %s to be ready, please check status on dashboard.", newClusterID))
			case <-ticker.C:
				clusterResult, err := d.GetCluster(ctx, newClusterID, cluster.CLUSTERSERVICEGETCLUSTERVIEWPARAMETER_BASIC)
				if err != nil {
					return errors.Trace(err)
				}
				if *clusterResult.State == cluster.COMMONV1BETA1CLUSTERSTATE_ACTIVE {
					return ui.Result(fmt.Sprintf("Cluster %s is ready.", newClusterID))
				}
			case <-ctx.Done():
				return util.InterruptError
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Waiting for cluster to be ready"))
	createModel, err := p.Run()
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
		Inputs: make([]textinput.Model, len(createClusterField)),
	}

	for k, v := range createClusterField {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 64

		switch k {
		case flag.DisplayName:
			t.Placeholder = "Display Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		}
		m.Inputs[v] = t
	}
	return m
}

func checkClusterName(name string) error {
	if len(name) == 0 {
		return errors.New("cluster name is required")
	}
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

func checkCapacity(minRcu, maxRcu int32) error {
	if minRcu < 2000 {
		return errors.New("Minimum RCU should be at least 2000")
	}
	if maxRcu > 100000 {
		return errors.New("Maximum RCU should be at most 100000")
	}
	if minRcu > maxRcu {
		return errors.New("Minimum RCU should be not greater than maximum RCU")
	}
	return nil
}

func getAndCheckNumber(n string, hint string) (int32, error) {
	if len(n) == 0 {
		return 0, nil
	} else {
		s, err := strconv.Atoi(n)
		if err != nil {
			return 0, fmt.Errorf("%s should be int type", hint)
		}
		if s > math.MaxInt32 || s < 0 {
			return 0, fmt.Errorf("%s out of range", hint)
		}
		return int32(s), nil //nolint:gosec // will not overflow
	}
}

func toInt64Ptr(v int32) *int64 {
	val := int64(v)
	return &val
}

func GetProvider(providers *hashset.Set) (string, error) {
	var cloudProvider string
	if providers.Size() == 0 {
		return cloudProvider, errors.New("No cloud provider available")
	}
	if providers.Size() == 1 {
		return providers.Values()[0].(string), nil
	}

	values := providers.Values()
	sort.Slice(values, func(i, j int) bool {
		s1 := values[i].(string)
		s2 := values[j].(string)
		return s1 < s2
	})
	model, err := ui.InitialSelectModel(values, "Choose the cloud provider:")
	if err != nil {
		return cloudProvider, err
	}
	p := tea.NewProgram(model)
	providerModel, err := p.Run()
	if err != nil {
		return cloudProvider, errors.Trace(err)
	}
	if m, _ := providerModel.(ui.SelectModel); m.Interrupted {
		return cloudProvider, util.InterruptError
	}
	cloudProvider = providerModel.(ui.SelectModel).Choices[providerModel.(ui.SelectModel).Selected].(string)
	return cloudProvider, nil
}

func GetClusterPlan() (cluster.ClusterClusterPlan, error) {
	choices := make([]interface{}, len(cluster.AllowedClusterClusterPlanEnumValues))
	for i, v := range cluster.AllowedClusterClusterPlanEnumValues {
		choices[i] = v
	}
	model, err := ui.InitialSelectModel(choices, "Choose the cluster plan:")
	if err != nil {
		return "", err
	}
	p := tea.NewProgram(model)
	planModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := planModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	clusterPlan := planModel.(ui.SelectModel).Choices[planModel.(ui.SelectModel).Selected].(cluster.ClusterClusterPlan)
	return clusterPlan, nil
}

func GetRegion(regionSet *hashset.Set) (string, error) {
	var region string
	if regionSet.Size() == 0 {
		return region, errors.New("No region available")
	}

	values := regionSet.Values()
	sort.Slice(values, func(i, j int) bool {
		s1 := values[i].(cloud.Region).Name
		s2 := values[j].(cloud.Region).Name
		return s1 < s2
	})
	model, err := ui.InitialSelectModel(values, "Choose the cloud region:")
	if err != nil {
		return region, err
	}
	p := tea.NewProgram(model)
	regionModel, err := p.Run()
	if err != nil {
		return region, errors.Trace(err)
	}
	if m, _ := regionModel.(ui.SelectModel); m.Interrupted {
		return region, util.InterruptError
	}
	region = regionModel.(ui.SelectModel).Choices[regionModel.(ui.SelectModel).Selected].(cloud.Region).Name
	return region, nil
}
