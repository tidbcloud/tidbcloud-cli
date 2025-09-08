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

type PrecheckOpts struct {
	interactive bool
}

func (p PrecheckOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.ConfigFile,
	}
}

func PrecheckCmd(h *internal.Helper) *cobra.Command {
	opts := PrecheckOpts{
		interactive: true,
	}

	cmd := &cobra.Command{
		Use:         "precheck",
		Short:       "Run precheck for a DM task",
		Long:        "Run precheck for a data migration (DM) task using a configuration file.",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Run precheck in interactive mode:
  $ %[1]s serverless dm precheck

  Run precheck in non-interactive mode:
  $ %[1]s serverless dm precheck --cluster-id <cluster-id> --config-file <config-file>

  Generate example JSON configuration:
  $ %[1]s serverless dm precheck --generate-json`,
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
				exampleJSON, err := GenerateExampleJSON("precheck")
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
				return errors.New("Interactive mode for DM precheck is not yet implemented. Please use non-interactive mode with --config-file flag")
			} else {
				// non-interactive mode
				clusterID, err = cmd.Flags().GetString(flag.ClusterID)
				if err != nil {
					return errors.Trace(err)
				}

				configFile, err = cmd.Flags().GetString(flag.ConfigFile)
				if err != nil {
					return errors.Trace(err)
				}
			}

			// Read configuration from file or stdin
			var reader io.Reader
			if configFile == "" || configFile == "-" {
				reader = os.Stdin
			} else {
				file, err := os.Open(configFile)
				if err != nil {
					return fmt.Errorf("failed to open config file: %w", err)
				}
				defer file.Close()
				reader = file
			}

			configData, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("failed to read config: %w", err)
			}

			var body dm.DMServicePrecheckBody
			if err := json.Unmarshal(configData, &body); err != nil {
				return fmt.Errorf("failed to parse config JSON: %w", err)
			}

			// Start precheck
			result, err := d.Precheck(ctx, clusterID, &body)
			if err != nil {
				return err
			}

			// Output result
			outputFormat, err := cmd.Flags().GetString(flag.Output)
			if err != nil {
				return err
			}

			if outputFormat == output.JsonFormat {
				return output.PrintJson(h.IOStreams.Out, result)
			}

			fmt.Printf("Precheck started successfully.\n")
			fmt.Printf("Precheck ID: %s\n", *result.Id)

			return nil
		},
	}

	cmd.Flags().StringP(flag.ConfigFile, "", "", "Config file path (use '-' for stdin)")
	cmd.Flags().Bool(flag.GenerateJSON, false, "Generate example JSON configuration file")
	cmd.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, "Output format (human|json)")
	cmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID")

	return cmd
}
