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
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type CapacityOpts struct {
	interactive bool
}

func (c CapacityOpts) NonInteractiveFlags() []string {
	return []string{
		flag.MaxRCU,
		flag.MinRCU,
		flag.ClusterID,
	}
}

var capacityFields = []string{
	flag.MinRCU,
	flag.MaxRCU,
}

var CapacityMask = "auto_scaling"

func CapacityCmd(h *internal.Helper) *cobra.Command {
	opts := CapacityOpts{
		interactive: true,
	}

	var capacityCmd = &cobra.Command{
		Use:         "capacity",
		Short:       "Set capacity for a TiDB Cloud Serverless cluster",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Set capacity for a TiDB Cloud Serverless cluster in interactive mode:
  $ %[1]s serverless capacity

  Set capacity for a TiDB Cloud Serverless cluster in non-interactive mode:
  $ %[1]s serverless capacity -c <cluster-id> --max-rcu <maximum-rcu> --min-rcu <minimum-rcu>`, config.CliName),
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
			ctx := cmd.Context()

			var clusterID string
			var maxRcu, minRcu int32
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}

				c, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = c.ID

				fmt.Fprintln(h.IOStreams.Out, color.BlueString("Please input the capacity for the cluster:"))
				p := tea.NewProgram(initialCapacityInputModel())
				inputModel, err := p.Run()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				minRcuString := inputModel.(ui.TextInputModel).Inputs[0].Value()
				maxRcuString := inputModel.(ui.TextInputModel).Inputs[1].Value()
				minRcu, err = getAndCheckNumber(minRcuString, "minimum RCU")
				if err != nil {
					return errors.Trace(err)
				}
				maxRcu, err = getAndCheckNumber(maxRcuString, "maximum RCU")
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				maxRcu, err = cmd.Flags().GetInt32(flag.MaxRCU)
				if err != nil {
					return errors.Trace(err)
				}
				minRcu, err = cmd.Flags().GetInt32(flag.MinRCU)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if err := checkCapacity(minRcu, maxRcu); err != nil {
				return errors.Trace(err)
			}

			body := &cluster.V1beta1ServerlessServicePartialUpdateClusterBody{
				Cluster: &cluster.RequiredTheClusterToBeUpdated{
					AutoScaling: &cluster.V1beta1ClusterAutoScaling{},
				},
			}
			body.UpdateMask = CapacityMask
			body.Cluster.AutoScaling.MinRcu = toInt64Ptr(minRcu)
			body.Cluster.AutoScaling.MaxRcu = toInt64Ptr(maxRcu)
			_, err = d.PartialUpdateCluster(ctx, clusterID, body)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString(fmt.Sprintf("set capacity to [%d, %d] cents success", minRcu, maxRcu)))

			return nil
		},
	}

	capacityCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	capacityCmd.Flags().Int32(flag.MinRCU, 0, "Minimum RCU for the cluster, at least 2000.")
	capacityCmd.Flags().Int32(flag.MaxRCU, 0, "Maximum RCU for the cluster, at most 100000.")
	return capacityCmd
}

func initialCapacityInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(capacityFields)),
	}

	for k, v := range capacityFields {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 64

		switch v {
		case flag.MinRCU:
			t.Placeholder = "Minimum RCU, at least 2000"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.MaxRCU:
			t.Placeholder = "Maximum RCU, at most 100000"
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		}
		m.Inputs[k] = t
	}
	return m
}
