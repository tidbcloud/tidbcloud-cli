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
	"slices"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type UpdateOpts struct {
	interactive bool
}

var updateSQLUserField = map[string]int{
	flag.Password: 0,
	flag.UserRole: 1,
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
  $ %[1]s serverless sql-user update -c <cluster-id> --user <user-name> --password <password> --role <role>`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.MarkInteractive(cmd)
			if err != nil {
				return errors.Trace(err)
			}

			if !opts.interactive {
				err := cmd.MarkFlagRequired(flag.ClusterID)
				if err != nil {
					return err
				}
				cmd.MarkFlagsMutuallyExclusive(flag.UserRole, flag.AddRole, flag.DeleteRole)
				cmd.MarkFlagsOneRequired(flag.Password, flag.UserRole, flag.AddRole, flag.DeleteRole)
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
			var userRole []string
			var addRole []string
			var deleteRole []string
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
				userPrefix = cluster.UserPrefix

				user, err := cloud.GetSelectedSQLUser(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				userName = user

				// variables for input
				fmt.Fprintln(h.IOStreams.Out, color.BlueString("Please input the following update options, require at least one field"))

				p := tea.NewProgram(initialUpdateInputModel())
				inputModel, err := p.Run()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				password = inputModel.(ui.TextInputModel).Inputs[updateSQLUserField[flag.Password]].Value()
				inputUserRole := inputModel.(ui.TextInputModel).Inputs[updateSQLUserField[flag.UserRole]].Value()

				userRole, err = util.StringSliceConv(inputUserRole)
				if err != nil {
					return errors.Trace(err)
				}

				if password == "" && len(userRole) == 0 {
					return errors.New("at least one of password and user role must be set")
				}

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

				uRole, err := cmd.Flags().GetStringSlice(flag.UserRole)
				if err != nil {
					return errors.Trace(err)
				}
				userRole = uRole

				aRole, err := cmd.Flags().GetStringSlice(flag.AddRole)
				if err != nil {
					return errors.Trace(err)
				}
				addRole = aRole

				dRole, err := cmd.Flags().GetStringSlice(flag.DeleteRole)
				if err != nil {
					return errors.Trace(err)
				}
				deleteRole = dRole
			}

			u, err := d.GetSQLUser(ctx, clusterID, userName)
			if err != nil {
				return errors.Trace(err)
			}

			if len(userRole) != 0 {
				u.BuiltinRole, u.CustomRoles, err = getBuiltinRoleAndCustomRoles(userRole)
				if err != nil {
					return errors.Trace(err)
				}
				if util.IsNilOrEmpty(u.BuiltinRole) {
					return errors.New("built-in role must be set")
				}
			}

			if len(addRole) != 0 {
				addBuiltinRole, addCustomRoles, err := getBuiltinRoleAndCustomRoles(addRole)
				if err != nil {
					return errors.Trace(err)
				}
				if !util.IsNilOrEmpty(u.BuiltinRole) && !util.IsNilOrEmpty(addBuiltinRole) {
					return errors.New("built-in role already exists in the SQL user")
				} else if util.IsNilOrEmpty(u.BuiltinRole) && !util.IsNilOrEmpty(addBuiltinRole) {
					u.BuiltinRole = addBuiltinRole
				}
				u.CustomRoles = append(u.CustomRoles, addCustomRoles...)
			}

			if len(deleteRole) != 0 {
				deleteBuiltinRole, deleteCustomRoles, err := getBuiltinRoleAndCustomRoles(deleteRole)
				if err != nil {
					return errors.Trace(err)
				}
				if !util.IsNilOrEmpty(deleteBuiltinRole) {
					builtinRoleWithPrefix := util.AddPrefix(*u.BuiltinRole, userPrefix)
					deleteBuiltinRole = &builtinRoleWithPrefix
					if deleteBuiltinRole == u.BuiltinRole {
						// it doesn't work yet, because the API doesn't support to delete the builtin role
						u.BuiltinRole = nil
					} else {
						return errors.New("can not delte built-in role")
					}
				}
				for _, role := range deleteCustomRoles {
					index := slices.Index(u.CustomRoles, role)
					if index == -1 {
						return errors.New(fmt.Sprintf("role %s doesn't exist in the SQL user", role))
					}
					u.CustomRoles = slices.Delete(u.CustomRoles, index, index+1)
				}
			}

			body := &iam.ApiUpdateSqlUserReq{}
			if *u.BuiltinRole != util.ADMIN_ROLE {
				builtinRoleWithPrefix := util.AddPrefix(*u.BuiltinRole, userPrefix)
				body.BuiltinRole = &builtinRoleWithPrefix
			} else {
				body.BuiltinRole = u.BuiltinRole
			}
			body.CustomRoles = u.CustomRoles

			if password != "" {
				body.Password = &password
			}

			_, err = d.UpdateSQLUser(ctx, clusterID, userName, body)
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
	updateCmd.Flags().StringSliceP(flag.UserRole, "", nil, "The new role(s) of the SQL user. Passing this flag replaces preexisting data.")
	updateCmd.Flags().StringSliceP(flag.AddRole, "", nil, "The role(s) to be added to the SQL user.")
	updateCmd.Flags().StringSliceP(flag.DeleteRole, "", nil, "The role(s) to be deleted from the SQL user.")

	return updateCmd
}

func initialUpdateInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(updateSQLUserField)),
	}

	for k, v := range updateSQLUserField {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 32

		switch k {
		case flag.Password:
			t.Placeholder = "New password"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case flag.UserRole:
			t.Placeholder = "New SQL user roles which replaces preexisting data, separated by comma"
		}
		m.Inputs[v] = t
	}
	return m
}
