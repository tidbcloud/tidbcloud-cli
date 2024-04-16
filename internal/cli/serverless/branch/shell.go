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

package branch

import (
	"fmt"
	"strconv"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/util"
	branchApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client/branch_service"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	_ "github.com/xo/usql/drivers/mysql"
	"github.com/xo/usql/env"
)

type ShellOpts struct {
	interactive bool
}

func (c ShellOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.BranchID,
	}
}

func (c *ShellOpts) MarkInteractive(cmd *cobra.Command) error {
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

func ShellCmd(h *internal.Helper) *cobra.Command {
	opts := ShellOpts{
		interactive: true,
	}

	var shellCmd = &cobra.Command{
		Use:   "shell",
		Short: "Connect to a branch",
		Long: `Connect to a branch. 
The connection forces the [ANSI SQL mode](https://dev.mysql.com/doc/refman/8.0/en/sql-mode.html#sqlmode_ansi) for the session.`,
		Example: fmt.Sprintf(`  Connect to a branch in interactive mode:
  $ %[1]s serverless branch shell

  Connect to a branch with default user in non-interactive mode:
  $ %[1]s serverless branch shell -c <cluster-id> -b <branch-id>

  Connect to a branch with default user and password in non-interactive mode:
  $ %[1]s serverless branch shell -c <cluster-id> -b <branch-id> --password <password>

  Connect to a branch with specific user and password in non-interactive mode:
  $ %[1]s serverless branch shell -c <cluster-id> -b <branch-id> -u <user-name> --password <password>`, config.CliName),

		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.MarkInteractive(cmd)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !h.IOStreams.CanPrompt {
				return fmt.Errorf("the stdout is not a terminal")
			}
			ctx := cmd.Context()

			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID, branchID, userName string
			var pass *string
			if opts.interactive {
				// interactive mode
				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				projectID := project.ID

				cluster, err := cloud.GetSelectedCluster(projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				branch, err := cloud.GetSelectedBranch(clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				branchID = branch.ID

				useDefaultUser := false
				prompt := &survey.Confirm{
					Message: "Use the default user?",
					Default: true,
				}
				err = survey.AskOne(prompt, &useDefaultUser)
				if err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					} else {
						return err
					}
				}

				var userInput string
				if !useDefaultUser {
					input := &survey.Input{
						Message: "Please input the user name:",
					}
					err = survey.AskOne(input, &userInput, survey.WithValidator(survey.Required))
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
					userName = userInput
				}
			} else {
				// non-interactive mode, get values from flags
				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}

				clusterID = cID

				// options flags
				bID, err := cmd.Flags().GetString(flag.BranchID)
				if err != nil {
					return errors.Trace(err)
				}
				branchID = bID

				uName, err := cmd.Flags().GetString(flag.User)
				if err != nil {
					return errors.Trace(err)
				}
				userName = uName

				if cmd.Flags().Changed(flag.Password) {
					password, err := cmd.Flags().GetString(flag.Password)
					if err != nil {
						return errors.Trace(err)
					}
					pass = &password
				}
			}

			var host, name, port string
			params := branchApi.NewBranchServiceGetBranchParams().WithClusterID(clusterID).WithBranchID(branchID)
			branchInfo, err := d.GetBranch(params)
			if err != nil {
				return errors.Trace(err)
			}
			host = branchInfo.Payload.Endpoints.Public.Host
			port = strconv.Itoa(int(branchInfo.Payload.Endpoints.Public.Port))
			name = *branchInfo.Payload.DisplayName
			if userName == "" {
				userName = fmt.Sprintf("%s.root", branchInfo.Payload.UserPrefix)
				fmt.Fprintln(h.IOStreams.Out, color.GreenString("Current user: ")+color.HiGreenString(userName))
			}

			// Set prompt style, see https://github.com/xo/usql/commit/d5db12eaa6fe48cd0a697831ad03d61611290576
			err = env.Set("PROMPT1", "%n"+"@"+name+"%/%R%#")
			if err != nil {
				return err
			}

			err = util.ExecuteSqlDialog(ctx, util.SERVERLESS, userName, host, port, pass, h.IOStreams.Out)
			if err != nil {
				return err
			}

			return nil
		},
	}

	shellCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	shellCmd.Flags().StringP(flag.BranchID, flag.BranchIDShort, "", "The ID of the branch.")
	shellCmd.Flags().String(flag.Password, "", "The password of the user.")
	shellCmd.Flags().StringP(flag.User, flag.UserShort, "", "A specific user for login if not using the default user.")
	return shellCmd
}
