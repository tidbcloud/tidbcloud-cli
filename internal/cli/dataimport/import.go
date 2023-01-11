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

package dataimport

import (
	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/cli/dataimport/start"

	"github.com/spf13/cobra"
)

func ImportCmd(h *internal.Helper) *cobra.Command {
	var importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import data into TiDB Cloud cluster",
	}

	importCmd.AddCommand(ListCmd(h))
	importCmd.AddCommand(CancelCmd(h))
	importCmd.AddCommand(start.StartCmd(h))
	importCmd.AddCommand(DescribeCmd(h))
	return importCmd
}
