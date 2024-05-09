// Copyright 2024 PingCAP, Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqluser

import (
	"fmt"
	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/util"

	// "tidbcloud-cli/internal/util"
	iamApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/client/account"
	iamModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/models"

	// "github.com/AlecAivazis/survey/v2"
	// "github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type UpdateOpts struct {
	interactive bool
}

func (c UpdateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.User,
	}
}

func (c *UpdateOpts) MarkInteractive(cmd *cobra.Command) error {
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

func UpdateCmd(h *internal.Helper) *cobra.Command {
	opts := UpdateOpts{
		interactive: true,
	}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a SQL user",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  Update a SQL user in interactive mode:
  $ %[1]s serverless sql-user update

  Update a SQL user in non-interactive mode:
  $ %[1]s serverless sql-user update -c <cluster-id> --user <user-name>`, config.CliName),
		Aliases: []string{"rm"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.MarkInteractive(cmd)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID string
			var userName string
			var userPrefix string
			var password string
			var userRole string
			if opts.interactive {
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				user, err := cloud.GetSelectedSQLUser(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				userName = user

			} else {
				// non-interactive mode doesn't need projectID
				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID = cID

				userPrefix, err = getUserPrefix(ctx, d, clusterID)
				if err != nil {
					return errors.Trace(err)
				}

				uName, err := cmd.Flags().GetString(flag.User)
				if err != nil {
					return errors.Trace(err)
				}
				// if the user name doesn't have the prefix, add the prefix
				uName = util.AddPrefix(uName, userPrefix)
				userName = uName

				pw, err := cmd.Flags().GetString(flag.Password)
				if err != nil {
					return errors.Trace(err)
				}
				password = pw

				uRole, err := cmd.Flags().GetString(flag.UserRole)
				if err != nil {
					return errors.Trace(err)
				}
				userRole = uRole
			}

			builtinRole, customRoles, err := getBuiltinRoleAndCustomRoles(userRole)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Println("customRoles: ", customRoles)

			body := &iamModel.APIUpdateSQLUserReq{}

			if builtinRole != "" {
				// if the role is role_readonly or role_readwrite, add the prefix
				if builtinRole != util.ADMIN_ROLE {
					builtinRole = util.AddPrefix(builtinRole, userPrefix)
				}
				body.BuiltinRole = builtinRole
			}

			body.CustomRoles = customRoles

			if password != "" {
				body.Password = password
			}
			fmt.Println(body)
			params := iamApi.NewPatchV1beta1ClustersClusterIDSQLUsersUserNameParams().
				WithClusterID(clusterID).
				WithUserName(userName).
				WithSQLUser(body).
				WithContext(ctx)
			_, err = d.UpdateSQLUser(params)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("SQL user %s is updated", userName))
			return nil
		},
	}

	updateCmd.Flags().StringP(flag.User, flag.UserShort, "", "The name of the SQL user to be updated.")
	updateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The cluster ID of the SQL user to be updated.")
	updateCmd.Flags().StringP(flag.Password, "", "", "The new password of the SQL user.")
	updateCmd.Flags().StringP(flag.UserRole, "", "", "The new role of the SQL user.")

	return updateCmd
}
