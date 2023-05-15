package start

import (
	"fmt"
	"io"
	"net/http"
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
	"github.com/spf13/cobra"
	exec "golang.org/x/sys/execabs"
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

type MysqlOpts struct {
	interactive bool
}

func (c MysqlOpts) NonInteractiveFlags() []string {
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

func MysqlCmd(h *internal.Helper) *cobra.Command {
	opts := MysqlOpts{
		interactive: true,
	}

	var mysqlCmd = &cobra.Command{
		Use:         "mysql",
		Short:       "Import from mysql into TiDB Cloud",
		Annotations: make(map[string]string),
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
			var projectID, clusterID, sourceHost, sourcePort, sourceUser, sourcePassword, sourceTable, sourceDatabase, userName, password, databaseName string
			var skipCreateTable bool
			helper := &mysqlHelperImpl{}
			err := helper.CheckMySQLClient()
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
				p := tea.NewProgram(initialMysqlInputModel())
				inputModel, err := p.StartReturningModel()
				if err != nil {
					return errors.Trace(err)
				}
				if inputModel.(ui.TextInputModel).Interrupted {
					return nil
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
				if len(sourcePassword) == 0 {
					return errors.New("Source password is required")
				}
				sourceDatabase = inputModel.(ui.TextInputModel).Inputs[sourceDatabaseIdx].Value()
				if len(sourceDatabase) == 0 {
					return errors.New("Source database is required")
				}
				sourceTable = inputModel.(ui.TextInputModel).Inputs[sourceTableIdx].Value()
				if len(sourceTable) == 0 {
					return errors.New("Source table is required")
				}

				fmt.Fprintln(h.IOStreams.Out, "Please select the cluster to import:")
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
						os.Exit(130)
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
							os.Exit(130)
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
						os.Exit(130)
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
						os.Exit(130)
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
						os.Exit(130)
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

			argsMysqldump := []string{
				"mysqldump",
				"-h",
				sourceHost,
				"-P",
				sourcePort,
				"-u",
				sourceUser,
				"-p" + sourcePassword,
				"--skip-add-drop-table",
				"--skip-add-locks",
				"--skip-triggers",
			}

			if skipCreateTable {
				argsMysqldump = append(argsMysqldump, "--no-create-info")
			}

			home, _ := os.UserHomeDir()
			sqlCacheFile := filepath.Join(home, config.HomePath, ".cache", "dump-"+time.Now().Format("2006-01-02T15-04-05")+".sql")
			defer deleteSqlCacheFile(h, sqlCacheFile)
			argsMysqldump = append(argsMysqldump, "-r")
			argsMysqldump = append(argsMysqldump, sqlCacheFile)
			argsMysqldump = append(argsMysqldump, sourceDatabase)
			argsMysqldump = append(argsMysqldump, sourceTable)

			goOS := runtime.GOOS
			if goOS == "darwin" {
				goOS = "macOS"
			}
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

			// Get connection string
			connectInfo, err := cloud.RetrieveConnectInfo(d)
			if err != nil {
				return err
			}
			connectionString, err := util.GenerateConnectionString(connectInfo, util.MysqlCliID, host, userName, port, clusterType, goOS)
			if err != nil {
				return err
			}
			connectionString = strings.Replace(connectionString, "${password}", password, -1)
			connectionString = strings.Replace(connectionString, "-D test", fmt.Sprintf("-D %s", databaseName), -1)
			fmt.Println(connectionString)

			if h.IOStreams.CanPrompt {
				err := updateAndSpinnerWait(h, helper, argsMysqldump, sqlCacheFile, connectionString)
				if err != nil {
					return err
				}
			} else {
				err := updateAndWaitReady(h, helper, argsMysqldump, sqlCacheFile, connectionString)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	mysqlCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "Project ID")
	mysqlCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")
	mysqlCmd.Flags().String(flag.SourceHost, "", "The host of the source Mysql")
	mysqlCmd.Flags().String(flag.SourcePort, "", "The port of the source Mysql")
	mysqlCmd.Flags().String(flag.SourceUser, "", "The user to login source Mysql")
	mysqlCmd.Flags().String(flag.SourcePassword, "", "The password to login source Mysql")
	mysqlCmd.Flags().String(flag.SourceDatabase, "", "The database of the source Mysql")
	mysqlCmd.Flags().String(flag.SourceTable, "", "The table to dump")
	mysqlCmd.Flags().String(flag.Database, "", "The target database")
	mysqlCmd.Flags().String(flag.User, "", "The user to login serverless cluster, default is <token>.root")
	mysqlCmd.Flags().Bool(flag.SkipCreateTable, false, "Skip create table step")
	mysqlCmd.Flags().String(flag.Password, "", "The password to login serverless cluster")

	return mysqlCmd
}

func initialMysqlInputModel() ui.TextInputModel {
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
			t.Placeholder = "The host of the source Mysql"
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		case sourcePortIdx:
			t.Placeholder = "The port of the source Mysql"
		case sourceUserIdx:
			t.Placeholder = "The user to login source Mysql"
		case sourcePasswordIdx:
			t.Placeholder = "The password to login source Mysql"
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

func updateAndWaitReady(h *internal.Helper, helper *mysqlHelperImpl, args []string, sqlCacheFile string, connectionString string) error {
	fmt.Fprintf(h.IOStreams.Out, "... Dumping data from source Mysql\n")

	err := helper.DumpFromMysql(args)
	if err != nil {
		return err
	}

	fmt.Fprintf(h.IOStreams.Out, "... Importing data to serverless\n")
	err = helper.ImportToServerless(sqlCacheFile, connectionString)
	if err != nil {
		return err
	}

	return nil
}

func updateAndSpinnerWait(h *internal.Helper, helper *mysqlHelperImpl, args []string, sqlCacheFile string, connectionString string) error {
	task := func() tea.Msg {
		res := make(chan error, 1)

		go func() {
			err := helper.DumpFromMysql(args)
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

	p := tea.NewProgram(ui.InitialSpinnerModel(task, "Dumping data from source Mysql"))
	model, err := p.StartReturningModel()
	if err != nil {
		return errors.Trace(err)
	}
	if m, _ := model.(ui.SpinnerModel); m.Err != nil {
		return m.Err
	} else {
		fmt.Fprintln(h.IOStreams.Out, color.GreenString(m.Output))
	}

	task = func() tea.Msg {
		res := make(chan error, 1)

		go func() {
			err := helper.ImportToServerless(sqlCacheFile, connectionString)
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

type MysqlHelper interface {
	DownloadCaFile(caFile string) error
	CheckMySQLClient() error
	DumpFromMysql(args []string) error
	ImportToServerless(sqlCacheFile string, connectionString string) error
}

type mysqlHelperImpl struct {
}

func (m *mysqlHelperImpl) DownloadCaFile(caFile string) error {
	// 下载文件的 URL
	url := "https://letsencrypt.org/certs/isrgrootx1.pem"

	// 创建 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return errors.Annotate(err, "Failed to download ca file")
	}
	defer resp.Body.Close()

	// 创建文件并打开
	file, err := os.Create(caFile)
	if err != nil {
		return errors.Annotate(err, "Failed to create ca file")
	}
	defer file.Close()

	// 将 HTTP 响应体复制到文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Annotate(err, "Failed to copy ca file")
	}

	return nil
}

// CheckMySQLClient checks whether the 'mysql' client exists and is configured in $PATH
func (m *mysqlHelperImpl) CheckMySQLClient() error {
	_, err := exec.LookPath("mysql")
	if err == nil {
		return nil
	}

	msg := "couldn't find the 'mysql' command-line tool required to run this command."

	switch runtime.GOOS {
	case "darwin":
		if HasHomebrew() {
			return fmt.Errorf("%s\nTo install, run: brew install mysql-client", msg)
		}
	}

	return fmt.Errorf("%s\nPlease install it and add to $PATH", msg)
}

// HasHomebrew check whether the user has installed brew
func HasHomebrew() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}
