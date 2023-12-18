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

package billing

import (
	"fmt"
	"regexp"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	biApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/billing/client/billing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type Opts struct {
	interactive bool
}

func (c Opts) NonInteractiveFlags() []string {
	return []string{
		flag.Month,
	}
}

func (c *Opts) MarkInteractive(cmd *cobra.Command) error {
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
		for _, fn := range flags {
			err := cmd.MarkFlagRequired(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Cmd(h *internal.Helper) *cobra.Command {
	opts := Opts{
		interactive: true,
	}
	var cmd = &cobra.Command{
		Use:   "bill",
		Short: "Get the bill for the given month",
		Example: fmt.Sprintf(`  Get the bill in interactive mode
  $ %[1]s bill

  Get the bill in non-interactive mode:
  $ %[1]s bill --month <YYYY-MM>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.MarkInteractive(cmd)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := h.Client()
			if err != nil {
				return err
			}
			var month string
			if opts.interactive {
				// input month
				inputModel, err := GetBillingInput()
				if err != nil {
					return err
				}
				month = inputModel.(ui.TextInputModel).Inputs[0].Value()
			} else {
				month, err = cmd.Flags().GetString(flag.Month)
				if err != nil {
					return errors.Trace(err)
				}
			}

			// check month format
			err = CheckMonthFormat(month)
			if err != nil {
				return errors.Trace(err)
			}
			param := biApi.NewGetBillsBilledMonthParams().WithBilledMonth(month)
			resp, err := d.GetBillsBilledMonth(param)
			if err != nil {
				return errors.Trace(err)
			}
			err = output.PrintJson(h.IOStreams.Out, resp.Payload)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
	}

	cmd.Flags().String(flag.Month, "", "The month of this bill, format is YYYY-MM, for example '2023-08'")
	return cmd
}

func initialBillingInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 1),
	}
	billedMonth := textinput.New()
	billedMonth.Placeholder = "Billed Month. e.g., 2023-08"
	billedMonth.Focus()
	billedMonth.PromptStyle = config.FocusedStyle
	billedMonth.TextStyle = config.FocusedStyle
	m.Inputs[0] = billedMonth
	return m
}

func GetBillingInput() (tea.Model, error) {
	p := tea.NewProgram(initialBillingInputModel())
	inputModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if inputModel.(ui.TextInputModel).Interrupted {
		return nil, util.InterruptError
	}
	return inputModel, nil
}

func CheckMonthFormat(month string) error {
	// month should be YYYY-MM
	isMatched, _ := regexp.MatchString(`^\d{4}-\d{2}$`, month)
	if !isMatched {
		return errors.Errorf("invalid month format: %s", month)
	}
	return nil
}
