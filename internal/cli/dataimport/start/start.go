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

package start

import (
	"fmt"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
	importOp "tidbcloud-cli/pkg/tidbcloud/import/client/import_service"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func StartCmd(h *internal.Helper) *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start an import task",
	}

	startCmd.AddCommand(LocalCmd(h))
	startCmd.AddCommand(S3Cmd(h))
	return startCmd
}

func waitStartOp(h *internal.Helper, d cloud.TiDBCloudClient, params *importOp.CreateImportParams) error {
	fmt.Fprintf(h.IOStreams.Out, "... Starting the import task\n")
	res, err := d.CreateImport(params)
	if err != nil {
		return err
	}

	fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", *(res.Payload.ID)))
	return nil
}

func spinnerWaitStartOp(h *internal.Helper, d cloud.TiDBCloudClient, params *importOp.CreateImportParams) error {
	task := func() tea.Msg {
		errChan := make(chan error)

		go func() {
			res, err := d.CreateImport(params)
			if err != nil {
				errChan <- err
				return
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("Import task %s started.", *(res.Payload.ID)))
			errChan <- nil
		}()

		ticker := time.NewTicker(1 * time.Second)
		timer := time.After(2 * time.Minute)
		for {
			select {
			case <-timer:
				return fmt.Errorf("timeout waiting for import task to start")
			case <-ticker.C:
				// continue
			case err := <-errChan:
				if err != nil {
					return err
				} else {
					return ui.Result("")
				}
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Starting import task"))
	createModel, err := p.StartReturningModel()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := createModel.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	}

	return nil
}
