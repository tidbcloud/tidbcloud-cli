// Copyright 2026 PingCAP, Inc.
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

package upgrade

import (
	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func Cmd(_ *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Show how to upgrade the CLI",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if config.IsUnderTiUP {
				return errors.New("the CLI is managed by TiUP, please upgrade it by `tiup update cloud`")
			}
			return errors.New("automatic upgrade is unavailable because releases are not authenticated; install a new version through a trusted channel: https://github.com/tidbcloud/tidbcloud-cli/releases")
		},
	}
}
