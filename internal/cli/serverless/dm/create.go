// Copyright 2025 PingCAP, Inc.
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

package dm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/dm"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

const (
// We can remove this since it's now in flag package
)

type CreateOpts struct {
	interactive bool
}

func (c CreateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ConfigFile,
	}
}

func CreateCmd(h *internal.Helper) *cobra.Command {
	opts := CreateOpts{
		interactive: true,
	}

	var createCmd = &cobra.Command{
		Use:         "create",
		Aliases:     []string{"start"},
		Short:       "Create a DM task",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Create a DM task in interactive mode:
  $ %[1]s serverless dm create

  Create a DM task in non-interactive mode:
  $ %[1]s serverless dm create --cluster-id <cluster-id> --config-file <config-file>

  Example config file content:
  {
    "name": "my-dm-task",
    "mode": "incremental",
    "sourceConfig": {
      "host": "source.example.com",
      "port": 3306,
      "user": "root",
      "password": "password"
    },
    "targetConfig": {
      "database": "target_db"
    }
  }`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// Check if generate-json is requested
			generateJSON, _ := cmd.Flags().GetBool(flag.GenerateJSON)
			if generateJSON {
				opts.interactive = false
				return nil
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
			generateJSON, err := cmd.Flags().GetBool(flag.GenerateJSON)
			if err != nil {
				return err
			}

			if generateJSON {
				exampleJSON, err := GenerateExampleJSON("create")
				if err != nil {
					return fmt.Errorf("failed to generate example JSON: %w", err)
				}
				fmt.Println(exampleJSON)
				return nil
			}

			var clusterID, configFile string
			d, err := h.Client()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

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

				cluster, err := cloud.GetSelectedCluster(ctx, project.ID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				clusterID = cluster.ID

				// TODO: Add interactive prompts for task configuration
				return errors.New("Interactive mode for DM task creation is not yet implemented. Please use non-interactive mode with --config-file flag")
			} else {
				// non-interactive mode
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				configFile = cmd.Flag(flag.ConfigFile).Value.String()
			}

			cmd.Annotations[telemetry.ClusterID] = clusterID

			// Read and parse the config file
			var file io.Reader
			if configFile == "-" {
				file = os.Stdin
			} else {
				f, err := os.Open(configFile)
				if err != nil {
					return errors.Trace(err)
				}
				defer f.Close()
				file = f
			}

			data, err := io.ReadAll(file)
			if err != nil {
				return errors.Trace(err)
			}

			var taskBody dm.DMServiceCreateTaskBody
			err = json.Unmarshal(data, &taskBody)
			if err != nil {
				return errors.Trace(err)
			}

			task, err := d.CreateTask(ctx, clusterID, &taskBody)
			if err != nil {
				return errors.Trace(err)
			}

			format, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return errors.Trace(err)
			}

			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				err := output.PrintJson(h.IOStreams.Out, task)
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				fmt.Fprintf(h.IOStreams.Out, "DM task %s created successfully.\n", *task.Id)
			}

			return nil
		},
	}

	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	createCmd.Flags().String(flag.ConfigFile, "", "Config file path (JSON format). Use '-' to read from stdin.")
	createCmd.Flags().Bool(flag.GenerateJSON, false, "Generate example JSON configuration file")
	createCmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	return createCmd
}
