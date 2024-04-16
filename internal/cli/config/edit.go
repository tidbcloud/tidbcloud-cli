// Copyright 2024 PingCAP, Inc.
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

package config

import (
	"fmt"
	"os"
	"runtime"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	exec "golang.org/x/sys/execabs"
)

const defaultEditor = "vi"

func EditCmd(h *internal.Helper) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "edit",
		Short: "Open the config file with the default text editor",
		Example: fmt.Sprintf(`  To open the config
  $ %[1]s config edit`, config.CliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runtime.GOOS == "windows" {
				fmt.Fprintln(h.IOStreams.Out, color.YellowString("Currently, opening config file is not supported in Windows.\nThe config file path is %s", viper.ConfigFileUsed()))
				return nil
			}

			c := exec.Command(defaultEditor, viper.ConfigFileUsed()) //nolint:gosec
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			return errors.Trace(c.Run())
		},
	}

	return listCmd
}
