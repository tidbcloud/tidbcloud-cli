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

package migration

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/juju/errors"

	pkgmigration "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/migration"
)

func parseMigrationDefinition(value string) ([]pkgmigration.Source, pkgmigration.Target, pkgmigration.TaskMode, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, pkgmigration.Target{}, "", errors.New("migration definition is required; use --definition")
	}
	var payload struct {
		Sources []pkgmigration.Source `json:"sources"`
		Target  *pkgmigration.Target  `json:"target"`
		Mode    string                `json:"mode"`
	}
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return nil, pkgmigration.Target{}, "", errors.Annotate(err, "invalid migration definition JSON")
	}
	if len(payload.Sources) == 0 {
		return nil, pkgmigration.Target{}, "", errors.New("migration definition must include at least one source")
	}
	if payload.Target == nil {
		return nil, pkgmigration.Target{}, "", errors.New("migration definition must include the target block")
	}
	mode, err := parseMigrationMode(payload.Mode)
	if err != nil {
		return nil, pkgmigration.Target{}, "", err
	}
	return payload.Sources, *payload.Target, mode, nil
}

func parseMigrationMode(value string) (pkgmigration.TaskMode, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", errors.New("mode is required in the migration definition")
	}
	normalized := strings.ToUpper(trimmed)
	if !strings.HasPrefix(normalized, "MODE_") {
		normalized = fmt.Sprintf("MODE_%s", normalized)
	}
	mode := pkgmigration.TaskMode(normalized)
	for _, allowed := range pkgmigration.AllowedTaskModeEnumValues {
		if mode == allowed {
			return mode, nil
		}
	}
	return "", errors.Errorf("invalid mode %q, allowed values: %s", value, strings.Join(taskModeValues(), ", "))
}

func taskModeValues() []string {
	values := make([]string, 0, len(pkgmigration.AllowedTaskModeEnumValues))
	for _, mode := range pkgmigration.AllowedTaskModeEnumValues {
		values = append(values, string(mode))
	}
	return values
}
