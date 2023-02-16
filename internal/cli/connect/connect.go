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

package connect

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"os/user"
	"strconv"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/service/cloud"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/fatih/color"
	"github.com/go-sql-driver/mysql"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	_ "github.com/xo/usql/drivers/mysql"
	"github.com/xo/usql/env"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/rline"
)

const (
	SERVERLESS = "SERVERLESS"
	DEVELOPER  = "DEVELOPER"
	DEDICATED  = "DEDICATED"
)

type ConnectOpts struct {
	interactive bool
}

func (c ConnectOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ProjectID,
	}
}

func ConnectCmd(h *internal.Helper) *cobra.Command {
	opts := ConnectOpts{
		interactive: true,
	}

	var connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to a TiDB Cloud cluster",
		Long: `Connect to a TiDB Cloud cluster; 
the connection forces the [ANSI SQL mode](https://dev.mysql.com/doc/refman/8.0/en/sql-mode.html#sqlmode_ansi) for the session.`,
		Example: fmt.Sprintf(`  Connect to the TiDB Cloud cluster in interactive mode:
  $ %[1]s connect

  Use the default user to connect to the TiDB Cloud cluster in non-interactive mode:
  $ %[1]s connect -p <project-id> -c <cluster-id>

  Use the default user to connect to the TiDB Cloud cluster with password in non-interactive mode:
  $ %[1]s connect -p <project-id> -c <cluster-id> --password <password> 

  Use a specific user to connect to the TiDB Cloud cluster in non-interactive mode:
  $ %[1]s connect -p <project-id> -c <cluster-id> -u <user-name>`, config.CliName),
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
			if !h.IOStreams.CanPrompt {
				return fmt.Errorf("the stdout is not a terminal")
			}

			d, err := h.Client()
			if err != nil {
				return err
			}

			var projectID, clusterID, userName string
			var pass *string
			if opts.interactive {
				// interactive mode
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

				var userInput string
				if !useDefaultUser {
					input := &survey.Input{
						Message: "Please input the user name:",
					}
					err = survey.AskOne(input, &userInput, survey.WithValidator(survey.Required))
					if err != nil {
						if err == terminal.InterruptErr {
							os.Exit(130)
						} else {
							return err
						}
					}
					userName = userInput
				}
			} else {
				// non-interactive mode, get values from flags
				pID, err := cmd.Flags().GetString(flag.ProjectID)
				if err != nil {
					return errors.Trace(err)
				}

				cID, err := cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}
				projectID = pID
				clusterID = cID

				// options flags
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

			params := clusterApi.NewGetClusterParams().
				WithProjectID(projectID).
				WithClusterID(clusterID)
			cluster, err := d.GetCluster(params)
			if err != nil {
				return errors.Trace(err)
			}
			defaultUser := cluster.Payload.Status.ConnectionStrings.DefaultUser
			host := cluster.Payload.Status.ConnectionStrings.Standard.Host
			clusterName := cluster.Payload.Name
			port := strconv.Itoa(int(cluster.Payload.Status.ConnectionStrings.Standard.Port))
			clusterType := cluster.Payload.ClusterType
			if userName == "" {
				userName = defaultUser
				fmt.Fprintln(h.IOStreams.Out, color.GreenString("Current user: ")+color.HiGreenString(userName))
			}
			if clusterType == DEVELOPER {
				clusterType = SERVERLESS
			}

			// Set prompt style, see https://github.com/xo/usql/commit/d5db12eaa6fe48cd0a697831ad03d61611290576
			err = env.Set("PROMPT1", "%n"+"@"+clusterName+"%/%R%#")
			if err != nil {
				return err
			}

			err = ExecuteSqlDialog(clusterType, userName, host, port, pass, h.IOStreams.Out)
			if err != nil {
				return err
			}
			return nil
		},
	}

	connectCmd.Flags().StringP(flag.ProjectID, flag.ProjectIDShort, "", "The project ID of the cluster")
	connectCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster")
	connectCmd.Flags().String(flag.Password, "", "The password of the user")
	connectCmd.Flags().StringP(flag.User, flag.UserShort, "", "A specific user for login if not using the default user")
	return connectCmd
}

func ExecuteSqlDialog(clusterType, userName, host, port string, pass *string, out io.Writer) error {
	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("can't get current user: %s", err.Error())
	}
	l, err := rline.New(false, "", env.HistoryFile(u))
	if err != nil {
		return fmt.Errorf("can't open history file: %s", err.Error())
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	h := handler.New(l, u, wd, true)

	var dsn string
	if pass == nil {
		dsn, err = generateDsnWithoutPassword(clusterType, userName, host, port, h)
		if err != nil {
			return err
		}
	} else {
		dsn, err = generateDsnWithPassword(clusterType, userName, host, port, *pass)
		if err != nil {
			return err
		}
	}

	if err = h.Open(context.TODO(), dsn); err != nil {
		return fmt.Errorf("can't open connection to %s: %s", dsn, err.Error())
	}

	if err = h.Run(); err != io.EOF {
		return err
	}
	return nil
}

func generateDsnWithPassword(clusterType string, userName string, host string, port string, pass string) (string, error) {
	var dsn string
	if clusterType == SERVERLESS {
		err := mysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		})
		if err != nil {
			return "", err
		}
		dsn = fmt.Sprintf("tidb://%s:%s@%s:%s/test?tls=tidb", userName, pass, host, port)
	} else if clusterType == DEDICATED {
		dsn = fmt.Sprintf("tidb://%s:%s@%s:%s/test?tls=skip-verify", userName, pass, host, port)
	} else {
		return "", fmt.Errorf("unsupproted cluster type: %s", clusterType)
	}
	return dsn, nil
}

func generateDsnWithoutPassword(clusterType string, userName string, host string, port string, h *handler.Handler) (string, error) {
	var dsn string
	if clusterType == SERVERLESS {
		err := mysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		})
		if err != nil {
			return "", err
		}
		dsn = fmt.Sprintf("tidb://%s@%s:%s/test?tls=tidb", userName, host, port)
	} else if clusterType == DEDICATED {
		dsn = fmt.Sprintf("tidb://%s@%s:%s/test?tls=skip-verify", userName, host, port)
	} else {
		return "", fmt.Errorf("unsupproted cluster type: %s", clusterType)
	}

	// Prompt for password
	dsn, err := h.Password(dsn)
	if err != nil {
		return "", err
	}
	return dsn, nil
}
