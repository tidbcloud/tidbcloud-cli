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

package config

import (
	"os"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	exec "golang.org/x/sys/execabs"
)

const defaultEditor = "vi"

func EditCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "edit",
		Short: "Opens the config file with the default text editor.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := exec.Command(defaultEditor, viper.ConfigFileUsed()) //nolint:gosec
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			return errors.Trace(c.Run())
		},
	}

	return listCmd
}
