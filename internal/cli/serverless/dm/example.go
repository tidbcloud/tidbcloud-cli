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

	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/dm"
)

// generatePrecheckExample creates a complete example of DMServicePrecheckBody with default values
func generatePrecheckExample() *dm.DMServicePrecheckBody {
	taskMode := dm.TASKMODE_MODE_INCREMENTAL

	return &dm.DMServicePrecheckBody{
		Name:                  dm.PtrString("example-dm-task"),
		Mode:                  &taskMode,
		FullInstanceMigration: dm.PtrBool(false),
		Sources: []dm.Source{
			{
				ConnProfile: &dm.ConnProfile{
					Host:     dm.PtrString("source.example.com"),
					Port:     dm.PtrInt32(3306),
					User:     dm.PtrString("root"),
					Password: dm.PtrString("password"),
				},
				BaRules: &dm.BlockAllowRules{
					DoDbs: []string{"source_db"},
					DoTables: []dm.Table{
						{
							Schema: dm.PtrString("source_db"),
							Table:  dm.PtrString("users"),
						},
					},
				},
			},
		},
		TargetDb: &dm.TargetDatabase{
			User:     dm.PtrString("target_user"),
			Password: dm.PtrString("target_password"),
		},
	}
}

// generateCreateTaskExample creates a complete example of DMServiceCreateTaskBody with default values
func generateCreateTaskExample() *dm.DMServiceCreateTaskBody {
	taskMode := dm.TASKMODE_MODE_INCREMENTAL

	return &dm.DMServiceCreateTaskBody{
		Name:                  dm.PtrString("example-dm-task"),
		Mode:                  &taskMode,
		FullInstanceMigration: dm.PtrBool(false),
		Sources: []dm.Source{
			{
				ConnProfile: &dm.ConnProfile{
					Host:     dm.PtrString("source.example.com"),
					Port:     dm.PtrInt32(3306),
					User:     dm.PtrString("root"),
					Password: dm.PtrString("password"),
				},
				BaRules: &dm.BlockAllowRules{
					DoDbs: []string{"source_db"},
					DoTables: []dm.Table{
						{
							Schema: dm.PtrString("source_db"),
							Table:  dm.PtrString("users"),
						},
					},
				},
			},
		},
		TargetDb: &dm.TargetDatabase{
			User:     dm.PtrString("target_user"),
			Password: dm.PtrString("target_password"),
		},
	}
}

// GenerateExampleJSON generates example JSON for the given type
func GenerateExampleJSON(exampleType string) (string, error) {
	var data interface{}

	switch exampleType {
	case "precheck":
		data = generatePrecheckExample()
	case "create":
		data = generateCreateTaskExample()
	default:
		data = generatePrecheckExample() // default to precheck
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
