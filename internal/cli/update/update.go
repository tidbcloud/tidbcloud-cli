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

package update

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"tidbcloud-cli/internal"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func UpdateCmd(h *internal.Helper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the CLI to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 创建一个上下文并且设置超时时间
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			c1 := exec.CommandContext(ctx, "curl", "https://raw.githubusercontent.com/tidbcloud/tidbcloud-cli/main/install.sh")
			if ctx.Err() == context.DeadlineExceeded {
				return errors.New("timeout when download the install.sh script")
			}
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("start to download the install.sh script"))
			out, err := c1.Output()
			if err != nil {
				return errors.Annotate(err, "failed to download the install.sh script")
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("done!"))
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("start to execute the install.sh script"))
			result, err := exec.CommandContext(ctx, "/bin/sh", "-c", string(out)).Output() //nolint:gosec
			if ctx.Err() == context.DeadlineExceeded {
				return errors.New("timeout when execute the install.sh script")
			}
			if err != nil {
				return errors.Annotate(err, "execute the install.sh script")
			}

			fmt.Fprintln(h.IOStreams.Out, string(result))
			return nil
		},
	}

	return cmd
}
