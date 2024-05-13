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
	"strings"
	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	// "tidbcloud-cli/internal/util"
	iamApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/client/account"
	iamModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/models"

	// "github.com/AlecAivazis/survey/v2"
	// "github.com/AlecAivazis/survey/v2/terminal"
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
  $ %[1]s serverless sql-user update -c <cluster-id> --user <user-name>`, config.CliName),
		Aliases: []string{"rm"},
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

				// check if at least one of the SQL user artributes is set
				flag1 := cmd.Flags().Lookup(flag.Password)
				flag2 := cmd.Flags().Lookup(flag.UserRole)
				flag3 := cmd.Flags().Lookup(flag.AddRole)
				flag4 := cmd.Flags().Lookup(flag.DeleteRole)
				if !flag1.Changed && !flag2.Changed && !flag3.Changed && !flag4.Changed {
					return errors.New(fmt.Sprintf("at least one of %s, %s, %s, %s must be set", flag.Password, flag.UserRole, flag.AddRole, flag.DeleteRole))
				}
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
				fmt.Fprintln(h.IOStreams.Out, color.BlueString("Please input the following update options"))

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

				userRole = strings.Split(inputUserRole, ",")
				// if inputUserRole is "", set userRole to nil
				if len(userRole) == 1 && userRole[0] == "" {
					userRole = nil
				}

				if password == "" && len(userRole) == 0 {
					return errors.New("At least one of password and user role must be set")
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

			getUserParams := iamApi.NewGetV1beta1ClustersClusterIDSQLUsersUserNameParams().
				WithClusterID(clusterID).
				WithUserName(userName).
				WithContext(ctx)
			getSQLUserResult, err := d.GetSQLUser(getUserParams)
			if err != nil {
				return errors.Trace(err)
			}
			u := getSQLUserResult.GetPayload()

			originBuiltinRole := u.BuiltinRole
			originCustomRoles := u.CustomRoles

			var builtinRole string
			var customRoles []string

			if len(userRole) != 0 {
				builtinRole, customRoles, err = getBuiltinRoleAndCustomRoles(userRole)
				if err != nil {
					return errors.Trace(err)
				}
			}

			if len(addRole) != 0 {
				addBuiltinRole, addCustomRoles, err := getBuiltinRoleAndCustomRoles(addRole)
				if err != nil {
					return errors.Trace(err)
				}
				if addBuiltinRole != "" {
					return errors.New("Built-in role can't be added when updating SQL user")
				}
				builtinRole = originBuiltinRole
				customRoles = append(originCustomRoles, addCustomRoles...)
			}

			if len(deleteRole) != 0 {
				deleteBuiltinRole, deleteCustomRoles, err := getBuiltinRoleAndCustomRoles(deleteRole)
				if err != nil {
					return errors.Trace(err)
				}
				if deleteBuiltinRole != "" {
					return errors.New("Built-in role can't be deleted when updating SQL user")
				}
				for _, role := range deleteCustomRoles {
					index := slices.Index(originCustomRoles, role)
					if index == -1 {
						return errors.New(fmt.Sprintf("Role %s doesn't exist in the SQL user", role))
					}
					originCustomRoles = slices.Delete(originCustomRoles, index, index+1)
				}
				builtinRole = originBuiltinRole
				customRoles = originCustomRoles
			}

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
	updateCmd.Flags().StringSliceP(flag.UserRole, "", nil, "The new role of the SQL user.")
	updateCmd.Flags().StringSliceP(flag.AddRole, "", nil, "The role to be added to the SQL user.")
	updateCmd.Flags().StringSliceP(flag.DeleteRole, "", nil, "The role to be deleted from the SQL user.")

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
			// add a prefix showing the user prefix
			t.Placeholder = "New Password"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case flag.UserRole:
			t.Placeholder = "New SQL User Roles, separated by comma"
		}
		m.Inputs[v] = t
	}
	return m
}
