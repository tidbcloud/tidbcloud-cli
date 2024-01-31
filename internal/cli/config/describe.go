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
	"encoding/json"
	"fmt"
	"strings"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/util"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DescribeCmd(h *internal.Helper) *cobra.Command {
	describeCmd := &cobra.Command{
		Use:     "describe <profile-name>",
		Aliases: []string{"get"},
		Short:   "Describe a specific profile",
		Example: fmt.Sprintf(`  Describe the profile configuration:
  $ %[1]s config describe <profile-name>`, config.CliName),
		Args: util.RequiredArgs("profile-name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.ToLower(args[0])
			err := config.ValidateProfile(name)
			if err != nil {
				return err
			}

			value := viper.Get(name)
			v, err := json.MarshalIndent(value, "", "  ")
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, string(v))
			return nil
		},
	}

	return describeCmd
}
