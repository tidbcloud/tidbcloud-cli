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
	"context"
	"fmt"
	"strings"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	serverlessApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
	iamClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/iam"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type CreateOpts struct {
	interactive bool
}

const (
	WaitInterval = 5 * time.Second
	WaitTimeout  = 2 * time.Minute

	DefaultAutoPrefix = true
)

var createSQLUserField = map[string]int{
	flag.User:     0,
	flag.Password: 1,
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.User,
		flag.Password,
		flag.UserRole,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.User,
		flag.Password,
		flag.UserRole,
	}
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var CreateCmd = &cobra.Command{
		Use:         "create",
		Short:       "Create a SQL user",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Create a SQL user in interactive mode:
  $ %[1]s serverless sql-user create

  Create a SQL user in non-interactive mode:
  $ %[1]s serverless sql-user create --user <user-name> --password <password> --role <role> --cluster-id <cluster-id>`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// mark required flags in non-interactive mode
			if !opts.interactive {
				for _, fn := range opts.RequiredFlags() {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
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
			var password string
			var userRole []string
			var userPrefix string
			var customRoles []string
			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				project, err := cloud.GetSelectedProject(ctx, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				projectID := project.ID

				cluster, err := cloud.GetSelectedCluster(ctx, projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID
				userPrefix = cluster.UserPrefix

				uRole, err := cloud.GetSelectedBuiltinRole()
				if err != nil {
					return err
				}
				userRole = append(userRole, uRole)

				// variables for input
				fmt.Fprintln(h.IOStreams.Out, color.BlueString("Please input the following options"))

				p := tea.NewProgram(initialCreateInputModel(userPrefix))
				inputModel, err := p.Run()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				userName = inputModel.(ui.TextInputModel).Inputs[createSQLUserField[flag.User]].Value()
				password = inputModel.(ui.TextInputModel).Inputs[createSQLUserField[flag.Password]].Value()

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
				// if the user name has the prefix, remove it
				uName = util.TrimUserNamePrefix(uName, userPrefix)
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
			}

			builtinRole, customRoles, err := getBuiltinRoleAndCustomRoles(userRole)
			if err != nil {
				return errors.Trace(err)
			}

			authMethod := util.MYSQLNATIVEPASSWORD
			autoPrefix := DefaultAutoPrefix
			params := &iamClient.ApiCreateSqlUserReq{
				AuthMethod:  &authMethod,
				UserName:    &userName,
				BuiltinRole: builtinRole,
				CustomRoles: customRoles,
				Password:    &password,
				AutoPrefix:  &autoPrefix,
			}

			_, err = d.CreateSQLUser(ctx, clusterID, params)
			if err != nil {
				return errors.Trace(err)
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("SQL user %s.%s is created", userPrefix, userName))
			if err != nil {
				return err
			}
			return nil

		},
	}

	CreateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	CreateCmd.Flags().StringP(flag.User, flag.UserShort, "", "The name of the SQL user.")
	CreateCmd.Flags().StringP(flag.Password, "", "", "The password of the SQL user.")
	CreateCmd.Flags().StringSliceP(flag.UserRole, "", nil, "The role(s) of the SQL user.")

	return CreateCmd
}

func initialCreateInputModel(userPrefix string) ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(createSQLUserField)),
	}

	for k, v := range createSQLUserField {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 32

		switch k {
		case flag.User:
			// add a prefix showing the user prefix
			t.Placeholder = "User Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
			t.Prompt = "> " + userPrefix + "."
		case flag.Password:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		m.Inputs[v] = t
	}
	return m
}

func getUserPrefix(ctx context.Context, d cloud.TiDBCloudClient, clusterID string) (string, error) {
	params := serverlessApi.NewServerlessServiceGetClusterParams().WithClusterID(clusterID).WithContext(ctx)

	cluster, err := d.GetCluster(params)
	if err != nil {
		return "", errors.Trace(err)
	}

	return cluster.Payload.UserPrefix, nil
}

func getBuiltinRoleAndCustomRoles(roles []string) (*string, []string, error) {
	builtinRole := ""
	customRoles := make([]string, 0, len(roles))
	for _, role := range roles {
		role = strings.TrimSpace(role)
		if util.IsBuiltinRole(role) {
			if builtinRole == "" {
				switch role {
				case util.ADMIN_ROLE:
					builtinRole = util.ADMIN_ROLE
				case util.READWRITE_ROLE:
					builtinRole = util.READWRITE_ROLE
				case util.READONLY_ROLE:
					builtinRole = util.READONLY_ROLE
				}
			} else {
				return nil, []string{}, errors.New("only one built-in role is allowed")
			}
		} else {
			customRoles = append(customRoles, role)
		}
	}
	return &builtinRole, customRoles, nil
}
