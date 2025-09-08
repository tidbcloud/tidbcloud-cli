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
	"fmt"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/dm"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type OperateOpts struct {
	interactive bool
}

func (c OperateOpts) NonInteractiveFlags() []string {
	return []string{
		flag.ClusterID,
		flag.TaskID,
		flag.Operation,
	}
}

func OperateCmd(h *internal.Helper) *cobra.Command {
	opts := OperateOpts{
		interactive: true,
	}

	var operateCmd = &cobra.Command{
		Use:         "operate",
		Short:       "Operate a DM task (start, stop, pause, resume)",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Operate a DM task in interactive mode:
  $ %[1]s serverless dm operate

  Operate a DM task in non-interactive mode:
  $ %[1]s serverless dm operate --cluster-id <cluster-id> --task-id <task-id> --operation <operation>

  Available operations: pause, resume`,
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
			var clusterID, taskID, operation string
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

				task, err := cloud.GetSelectedDMTask(ctx, clusterID, h.QueryPageSize, d)
				if err != nil {
					return err
				}
				taskID = *task.Id

				// TODO: Add interactive selection for operation
				return errors.New("Interactive mode for DM task operation is not yet implemented. Please use non-interactive mode with --operation flag")
			} else {
				// non-interactive mode
				clusterID = cmd.Flag(flag.ClusterID).Value.String()
				taskID = cmd.Flag(flag.TaskID).Value.String()
				operation = cmd.Flag(flag.Operation).Value.String()
			}

			cmd.Annotations[telemetry.ClusterID] = clusterID

			// Validate operation and convert to enum
			var operationEnum *dm.OperateDMTaskReqOperation
			switch operation {
			case "pause":
				op := dm.OPERATEDMTASKREQOPERATION_OPERATION_PAUSE
				operationEnum = &op
			case "resume":
				op := dm.OPERATEDMTASKREQOPERATION_OPERATION_RESUME
				operationEnum = &op
			default:
				return fmt.Errorf("invalid operation: %s. Valid operations are: pause, resume", operation)
			}

			// Create operation request body
			operateBody := &dm.DMServiceOperateTaskBody{
				Op: operationEnum,
			}

			err = d.OperateTask(ctx, clusterID, taskID, operateBody)
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Fprintf(h.IOStreams.Out, "DM task %s operation '%s' has been executed successfully.\n", taskID, operation)
			return nil
		},
	}

	operateCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "Cluster ID.")
	operateCmd.Flags().String(flag.TaskID, "", "Task ID.")
	operateCmd.Flags().String(flag.Operation, "", "Operation to perform (pause, resume).")
	return operateCmd
}
