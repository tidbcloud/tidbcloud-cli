// Copyright 2023 PingCAP, Inc.
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

package start

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"
	"tidbcloud-cli/internal/telemetry"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	DEVELOPER  = "DEVELOPER"
	SERVERLESS = "SERVERLESS"
)

type mysqlImportField int

const (
	sourceHostIdx mysqlImportField = iota
	sourcePortIdx
	sourceUserIdx
	sourcePasswordIdx
	sourceDatabaseIdx
	sourceTableIdx
)

type MySQLOpts struct {
	interactive bool
}

func (c MySQLOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
		flag.Database,
		flag.SourceHost,
		flag.SourcePort,
		flag.SourceDatabase,
		flag.SourceTable,
		flag.SourceUser,
		flag.SourcePassword,
		flag.Password,
	}
}

func MySQLCmd(h *internal.Helper) *cobra.Command {
	opts := MySQLOpts{
		interactive: true,
	}

	var mysqlCmd = &cobra.Command{
		Use:   "mysql",
		Short: "Import from MySQL into TiDB Cloud serverless cluster",
		Long: `This command dumps data from MySQL and imports it into TiDB Cloud serverless cluster. 
It depends on 'mysql' command-line tool, please make sure you have installed it and add to path.`,
		Annotations: make(map[string]string),
		Args:        cobra.NoArgs,
		Example: fmt.Sprintf(`  Start an import task in interactive mode:
  $ %[1]s import start mysql

  Start an import task in non-interactive mode:
  $ %[1]s import start mysql --project-id <project-id> --cluster-id <cluster-id> --source-host <source-host> --source-port <source-port> --source-user <source-user> --source-password <source-password> --source-database <source-database> --source-table <source-table> --database <database> --password <password>

  Start an import task with a specific user:
  $ %[1]s import start mysql --project-id <project-id> --cluster-id <cluster-id> --source-host <source-host> --source-port <source-port> --source-user <source-user> --source-password <source-password> --source-database <source-database> --source-table <source-table> --database <database> --password <password> --user <user>

  Start an import task skipping create table:
  $ %[1]s import start mysql --project-id <project-id> --cluster-id <cluster-id> --source-host <source-host> --source-port <source-port> --source-user <source-user> --source-password <source-password> --source-database <source-database> --source-table <source-table> --database <database> --password <password> --skip-create-table	
`,
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
				for _, fn := range flags {
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
			var projectID, clusterID, sourceHost, sourcePort, sourceUser, sourcePassword, sourceTable, sourceDatabase, userName, password, databaseName string
			var skipCreateTable bool
			err := h.MySQLHelper.CheckMySQLClient()
			if err != nil {
				return err
			}

			d, err := h.Client()
			if err != nil {
				return err
			}

			if opts.interactive {
				cmd.Annotations[telemetry.InteractiveMode] = "true"
				if !h.IOStreams.CanPrompt {
					return errors.New("The terminal doesn't support interactive mode, please use non-interactive mode")
				}

				// interactive mode
				// variables for input
				p := tea.NewProgram(initialMySQLInputModel())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return util.InterruptError
				}

				sourceHost = inputModel.(ui.TextInputModel).Inputs[sourceHostIdx].Value()
				if len(sourceHost) == 0 {
					return errors.New("Source host is required")
				}
				sourcePort = inputModel.(ui.TextInputModel).Inputs[sourcePortIdx].Value()
				if len(sourcePort) == 0 {
					return errors.New("Source port is required")
				}
				sourceUser = inputModel.(ui.TextInputModel).Inputs[sourceUserIdx].Value()
				if len(sourceUser) == 0 {
					return errors.New("Source user is required")
				}
				sourcePassword = inputModel.(ui.TextInputModel).Inputs[sourcePasswordIdx].Value()
				sourceDatabase = inputModel.(ui.TextInputModel).Inputs[sourceDatabaseIdx].Value()
				if len(sourceDatabase) == 0 {
					return errors.New("Source database is required")
				}
				sourceTable = inputModel.(ui.TextInputModel).Inputs[sourceTableIdx].Value()
				if len(sourceTable) == 0 {
					return errors.New("Source table is required")
				}

				project, err := cloud.GetSelectedProject(h.QueryPageSize, d)
				if err != nil {
					return err
				}
				projectID = project.ID

				cluster, err := cloud.GetSelectedCluster(projectID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

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

				if !useDefaultUser {
					input := &survey.Input{
						Message: "Please input the user name:",
					}
					err = survey.AskOne(input, &userName, survey.WithValidator(survey.Required))
					if err != nil {
						if err == terminal.InterruptErr {
							return util.InterruptError
						} else {
							return err
						}
					}
				}

				passInput := &survey.Password{
					Message: "Please input the user password:",
				}
				err = survey.AskOne(passInput, &password, survey.WithValidator(survey.Required))
				if err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					} else {
						return err
					}
				}

				input := &survey.Input{
					Message: "Please input the target database:",
				}
				err = survey.AskOne(input, &databaseName, survey.WithValidator(survey.Required))
				if err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					} else {
						return err
					}
				}

				prompt = &survey.Confirm{
					Message: "Skip create table?",
					Default: false,
				}
				err = survey.AskOne(prompt, &skipCreateTable)
				if err != nil {
					if err == terminal.InterruptErr {
						return util.InterruptError
					} else {
						return err
					}
				}
			} else {
				// non-interactive mode, get values from flags
				var err error
				projectID, err = cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return errors.Trace(err)
				}
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				sourceHost, err = cmd.Flags().GetString(flag.SourceHost)
				if err != nil {
					return errors.Trace(err)
				}
				sourcePort, err = cmd.Flags().GetString(flag.SourcePort)
				if err != nil {
					return errors.Trace(err)
				}
				sourceUser, err = cmd.Flags().GetString(flag.SourceUser)
				if err != nil {
					return errors.Trace(err)
				}
				sourcePassword, err = cmd.Flags().GetString(flag.SourcePassword)
				if err != nil {
					return errors.Trace(err)
				}
				sourceDatabase, err = cmd.Flags().GetString(flag.SourceDatabase)
				if err != nil {
					return errors.Trace(err)
				}
				sourceTable, err = cmd.Flags().GetString(flag.SourceTable)
				if err != nil {
					return errors.Trace(err)
				}
				password, err = cmd.Flags().GetString(flag.Password)
				if err != nil {
					return errors.Trace(err)
				}
				skipCreateTable, err = cmd.Flags().GetBool(flag.SkipCreateTable)
				if err != nil {
					return errors.Trace(err)
				}
				databaseName, err = cmd.Flags().GetString(flag.Database)
				if err != nil {
					return errors.Trace(err)
				}

				if cmd.Flags().Changed(flag.User) {
					userName, err = cmd.Flags().GetString(flag.User)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			cmd.Annotations[telemetry.ProjectID] = projectID

			mysqldumpCommand := []string{
				"mysqldump",
				"-h",
				sourceHost,
				"-P",
				sourcePort,
				"-u",
				sourceUser,
				"--password=" + sourcePassword,
				"--skip-add-drop-table",
				"--skip-add-locks",
				"--skip-triggers",
			}

			if skipCreateTable {
				mysqldumpCommand = append(mysqldumpCommand, "--no-create-info")
			}

			sqlCacheFile := h.MySQLHelper.GenerateSqlCachePath()
			defer deleteSqlCacheFile(h, sqlCacheFile)
			mysqldumpCommand = append(mysqldumpCommand, "-r")
			mysqldumpCommand = append(mysqldumpCommand, sqlCacheFile)
			mysqldumpCommand = append(mysqldumpCommand, sourceDatabase)
			mysqldumpCommand = append(mysqldumpCommand, sourceTable)

			// Get cluster info
			params := clusterApi.NewGetClusterParams().WithProjectID(projectID).WithClusterID(clusterID)
			clusterInfo, err := d.GetCluster(params)
			if err != nil {
				return err
			}

			// Resolve cluster information
			// Get connect parameter
			if userName == "" {
				userName = clusterInfo.Payload.Status.ConnectionStrings.DefaultUser
			}
			host := clusterInfo.Payload.Status.ConnectionStrings.Standard.Host
			port := strconv.Itoa(int(clusterInfo.Payload.Status.ConnectionStrings.Standard.Port))
			clusterType := clusterInfo.Payload.ClusterType
			if clusterType != DEVELOPER {
				return errors.New("Only serverless cluster is supported")
			} else {
				clusterType = SERVERLESS
			}

			goOS := runtime.GOOS
			if goOS == "darwin" {
				goOS = "macOS"
			}
			// Get connection string
			connectInfo, err := cloud.RetrieveConnectInfo(d)
			if err != nil {
				return err
			}
			connectionString, err := util.GenerateConnectionString(connectInfo, util.MysqlCliID, host, userName, port, clusterType, goOS, util.GolangCommand)
			if err != nil {
				return err
			}
			connectionString = strings.Replace(connectionString, "${password}", password, -1)
			connectionString = strings.Replace(connectionString, "-D test", fmt.Sprintf("-D %s", databaseName), -1)
			if runtime.GOOS != "darwin" {
				home, _ := os.UserHomeDir()
				caFile := filepath.Join(home, config.HomePath, "isrgrootx1.pem")
				_, err := os.Stat(caFile)
				if os.IsNotExist(err) {
					err := h.MySQLHelper.DownloadCaFile(caFile)
					if err != nil {
						return err
					}
				}
				connectionString = strings.Replace(connectionString, "<path_to_ca_cert>", caFile, -1)
			}
			importCommand := strings.Split(connectionString, " ")
			log.Debug("Print dump command", zap.Any("command", mysqldumpCommand))
			log.Debug("Print import command", zap.Any("command", importCommand))

			if h.IOStreams.CanPrompt {
				err := updateAndSpinnerWait(ctx, h, mysqldumpCommand, importCommand, sqlCacheFile)
				if err != nil {
					return err
				}
			} else {
				err := updateAndWaitReady(ctx, h, mysqldumpCommand, importCommand, sqlCacheFile)
				if err != nil {
					return err
				}
			}

			fmt.Fprintln(h.IOStreams.Out, "Data has been imported successfully")
			return nil
		},
	}

	mysqlCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "Project ID")
	mysqlCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	mysqlCmd.Flags().String(flag.SourceHost, "", "The host of the source MySQL")
	mysqlCmd.Flags().String(flag.SourcePort, "", "The port of the source MySQL")
	mysqlCmd.Flags().String(flag.SourceUser, "", "The user to login source MySQL")
	mysqlCmd.Flags().String(flag.SourcePassword, "", "The password to login source MySQL")
	mysqlCmd.Flags().String(flag.SourceDatabase, "", "The database of the source MySQL")
	mysqlCmd.Flags().String(flag.SourceTable, "", "The table to dump")
	mysqlCmd.Flags().String(flag.Database, "", "The target database")
	mysqlCmd.Flags().String(flag.User, "", "The user to login serverless cluster, default is '<token>.root'")
	mysqlCmd.Flags().Bool(flag.SkipCreateTable, false, "Skip create table step, default create table")
	mysqlCmd.Flags().String(flag.Password, "", "The password to login serverless cluster")

	return mysqlCmd
}

func initialMySQLInputModel() ui.TextInputModel {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, 6),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = config.FocusedStyle
		t.CharLimit = 0
		f := mysqlImportField(i)

		switch f {
		case sourceHostIdx:
			t.Placeholder = "The host of the source MySQL"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case sourcePortIdx:
			t.Placeholder = "The port of the source MySQL"
		case sourceUserIdx:
			t.Placeholder = "The user to login source MySQL"
		case sourcePasswordIdx:
			t.Placeholder = "The password to login source MySQL"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case sourceDatabaseIdx:
			t.Placeholder = "The database to dump"
		case sourceTableIdx:
			t.Placeholder = "The table to dump"
		}

		m.Inputs[i] = t
	}

	return m
}

func updateAndWaitReady(ctx context.Context, h *internal.Helper, mysqlDumpCommand []string, importCommand []string, sqlCacheFile string) error {
	fmt.Fprintf(h.IOStreams.Out, "... Dumping data from source MySQL\n")

	err := h.MySQLHelper.DumpFromMySQL(ctx, mysqlDumpCommand, sqlCacheFile)
	if err != nil {
		return err
	}

	fmt.Fprintf(h.IOStreams.Out, "... Importing data to serverless\n")
	err = h.MySQLHelper.ImportToServerless(ctx, importCommand, sqlCacheFile)
	if err != nil {
		return err
	}

	return nil
}

func updateAndSpinnerWait(ctx context.Context, h *internal.Helper, mysqlDumpCommand []string, importCommand []string, sqlCacheFile string) error {
	task := func() tea.Msg {
		res := make(chan error, 1)

		go func() {
			err := h.MySQLHelper.DumpFromMySQL(ctx, mysqlDumpCommand, sqlCacheFile)
			if err != nil {
				res <- err
			}

			res <- nil
		}()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case err := <-res:
				if err != nil {
					return err
				} else {
					return ui.Result("Dump from MySQL finished!")
				}
			case <-ticker.C:
				// continue
			}
		}
	}

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Dumping data from source MySQL"))
	model, err := p.StartReturningModel()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := model.(ui.SpinnerModel); m.Interrupted {
		return util.InterruptError
	}
	if m, _ := model.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	} else {
		fmt.Fprintln(h.IOStreams.Out, color.GreenString(m.Output))
	}

	task = func() tea.Msg {
		res := make(chan error, 1)

		go func() {
			err := h.MySQLHelper.ImportToServerless(ctx, importCommand, sqlCacheFile)
			if err != nil {
				res <- err
			}

			res <- nil
		}()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case err := <-res:
				if err != nil {
					return err
				} else {
					return ui.Result("Import to serverless cluster finished!")
				}
			case <-ticker.C:
				// continue
			}
		}
	}

	p = tea.NewProgram(ui.InitialSpinnerModel(task, "Importing data to serverless cluster"))
	model, err = p.StartReturningModel()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := model.(ui.SpinnerModel); m.Interrupted {
		return util.InterruptError
	}
	if m, _ := model.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	} else {
		fmt.Fprintln(h.IOStreams.Out, color.GreenString(m.Output))
	}
	return nil
}

func deleteSqlCacheFile(h *internal.Helper, sqlCacheFile string) {
	_, err := os.Stat(sqlCacheFile)

	if os.IsNotExist(err) {
		return
	}

	// Remove the file
	err = os.Remove(sqlCacheFile)
	if err != nil {
		fmt.Fprintf(h.IOStreams.Err, "Failed to remove cache file %s, error: %s", sqlCacheFile, err)
	}
}
