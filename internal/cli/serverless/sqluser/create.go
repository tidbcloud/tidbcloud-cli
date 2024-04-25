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
	"strings"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/output"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	iamApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/client/account"
	iamModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/models"

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
)

var createSQLUserField = map[string]int{
	flag.UserName: 0,
	flag.Password: 1,
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.UserName,
		flag.Password,
		flag.UserRole,
	}
}

func (c CreateOpts) RequiredFlags() []string {
	return []string{
		flag.ClusterID,
		flag.UserName,
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
		Short:       "Create a SQL users",
		Aliases:     []string{"c"},
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Create a TiDB Serverless SQL user in interactive mode:
$ %[1]s serverless sql-user create

Create a TiDB Serverless SQL user in non-interactive mode:
$ %[1]s serverless sql-user create --name <user-name> --password <password> --role <role> --cluster-id <cluster-id>`,
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
			d, err := h.Client()
			if err != nil {
				return err
			}

			var clusterID string
			var userName string
			var password string
			var userRole string
			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

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

				userRole, err = cloud.GetSelectedBuildinRole()
				if err != nil {
					return err
				}

				// variables for input
				fmt.Fprintln(h.IOStreams.Out, color.BlueString("Please input the following options"))

				p := tea.NewProgram(initialCreateInputModel())
				inputModel, err := p.Run()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				userName = inputModel.(ui.TextInputModel).Inputs[createSQLUserField[flag.UserName]].Value()
				password = inputModel.(ui.TextInputModel).Inputs[createSQLUserField[flag.Password]].Value()

			} else {
				// non-interactive mode doesn't need projectID
				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID = cID

				uName, err := cmd.Flags().GetString(flag.UserName)
				if err != nil {
					return errors.Trace(err)
				}
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

			// generate the built-in role
			builtinRole := GetBuiltInRole(userRole)

			params := iamApi.NewPostV1beta1ClustersClusterIDSQLUsersParams().
				WithClusterID(clusterID).
				WithSQLUser(
					&iamModel.APICreateSQLUserReq{
						AuthMethod:  util.MYSQLNATIVEPASSWORD,
						UserName:    userName,
						Password:    password,
						BuiltinRole: builtinRole,
					},
				)

			_, err = d.CreateSQLUser(params)
			if err != nil {
				return errors.Trace(err)
			}

			_, err = fmt.Fprintln(h.IOStreams.Out, color.GreenString("SQL user %s is created", userName))
			if err != nil {
				return err
			}
			return nil

		},
	}

	CreateCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	CreateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster.")
	CreateCmd.Flags().StringP(flag.UserName, "", "", "The name of the SQL user.")
	CreateCmd.Flags().StringP(flag.Password, "", "", "The password of the SQL user.")
	CreateCmd.Flags().StringP(flag.UserRole, "", "", "The built-in role of the SQL user.")

	return CreateCmd
}

func GetBuiltInRole(userRole string) string {
	role := strings.ToLower(userRole)
	switch role {
	case util.ADMIN:
		role = util.ADMIN_ROLE
	case util.READWRITE:
		role = util.READWRITE_ROLE
	case util.READONLY:
		role = util.READONLY_ROLE
	default:
		role = userRole
	}

	return role
}

func initialCreateInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(createSQLUserField)),
	}

	for k, v := range createSQLUserField {
		t := textinput.New()
		t.Cursor.Style = config.CursorStyle
		t.CharLimit = 32

		switch k {
		case flag.UserName:
			t.Placeholder = "User Name"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case flag.Password:
			t.Placeholder = "Password"
		}
		m.Inputs[v] = t
	}
	return m
}
