// Copyright 2023 PingCAP, Inc.
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
	"fmt"
	"strconv"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
	serverlessModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type SpendingLimitOpts struct {
	interactive bool
}

func (c SpendingLimitOpts) NonInteractiveFlags() []string {
	return []string{
		flag.Monthly,
		flag.ClusterID,
	}
}

var spendingLimitFields = []string{
	flag.Monthly,
}

var SpendingLimitMonthlyMask = "spendingLimit.monthly"

func SpendingLimitCmd(h *internal.Helper) *cobra.Command {
	opts := SpendingLimitOpts{
		interactive: true,
	}

	var spendingLimitCmd = &cobra.Command{
		Use:         "spending-limit",
		Short:       "Set spending limit for a TiDB Serverless cluster",
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Set spending limit for a TiDB Serverless cluster in interactive mode:
  $ %[1]s serverless spending-limit

  Set spending limit for a TiDB Serverless cluster in non-interactive mode:
  $ %[1]s serverless spending-limit -c <cluster-id> --monthly <spending-limit-monthly>`, config.CliName),
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
				for _, fn := range flags {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID string
			var monthly int32
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}

				cluster, err := cloud.GetSelectedCluster(project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				field, err := cloud.GetSpendingLimitField(spendingLimitFields)
				if err != nil {
					return err
				}

				switch field {
				case flag.Monthly:
					inputModel, err := GetMonthlyInput()
					if err != nil {
						return err
					}
					monthlyValue := inputModel.(ui.TextInputModel).Inputs[0].Value()
					monthlyInt64, err := strconv.ParseInt(monthlyValue, 10, 32)
					if err != nil {
						return errors.Errorf("invalid monthly spending limit %s", monthlyValue)
					}
					monthly = int32(monthlyInt64)
				default:
					return errors.Errorf("unsupported spending limit field %s", field)
				}
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				monthly, err = cmd.Flags().GetInt32(flag.Monthly)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if monthly <= 0 {
				return errors.Errorf("invalid monthly spending limit %d", monthly)
			}

			body := &serverlessApi.ServerlessServicePartialUpdateClusterBody{
				Cluster: &serverlessApi.ServerlessServicePartialUpdateClusterParamsBodyCluster{
					SpendingLimit: &serverlessModel.ClusterSpendingLimit{},
				},
			}
			body.UpdateMask = &SpendingLimitMonthlyMask
			body.Cluster.SpendingLimit.Monthly = monthly
			params := serverlessApi.NewServerlessServicePartialUpdateClusterParams().WithClusterClusterID(clusterID).WithBody(*body)
			_, err = d.PartialUpdateCluster(params)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("set spending limit to %d cents success", monthly)))

			return nil
		},
	}

	spendingLimitCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	spendingLimitCmd.Flags().Int32(flag.Monthly, 0, "Maximum monthly spending limit in USD cents.")
	return spendingLimitCmd
}

func GetMonthlyInput() (tea.Model, error) {

	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 1),
	}
	t := textinput.New()
	t.Cursor.Style = config.CursorStyle
	t.CharLimit = 64
	t.Placeholder = "Input monthly spending limit in USD cents (int)"
	t.Focus()
	t.PromptStyle = config.FocusedStyle
	t.TextStyle = config.FocusedStyle
	m.Inputs[0] = t

	p := tea.NewProgram(m)
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}
