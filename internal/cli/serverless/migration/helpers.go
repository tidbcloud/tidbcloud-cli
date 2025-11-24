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

	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	pkgmigration "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/migration"
)

var migrationInputDescription = map[string]string{
	flag.DisplayName:      "Display name for the migration.",
	flag.MigrationSources: "Sources definition in JSON. Use \"ticloud serverless migration template --type sources\" as a reference.",
	flag.MigrationTarget:  "Target definition in JSON. Use \"ticloud serverless migration template --type target\" as a reference.",
	flag.MigrationMode:    "Migration mode, one of MODE_ALL or MODE_INCREMENTAL.",
}

func parseMigrationSources(value string) ([]pkgmigration.Source, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, errors.New("sources is required, use --sources or provide it in interactive mode")
	}
	var sources []pkgmigration.Source
	if err := json.Unmarshal([]byte(trimmed), &sources); err != nil {
		return nil, errors.Annotate(err, "invalid sources JSON")
	}
	if len(sources) == 0 {
		return nil, errors.New("sources must contain at least one entry")
	}
	return sources, nil
}

func parseMigrationTarget(value string) (pkgmigration.Target, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return pkgmigration.Target{}, errors.New("target is required, use --target or provide it in interactive mode")
	}
	var target pkgmigration.Target
	if err := json.Unmarshal([]byte(trimmed), &target); err != nil {
		return pkgmigration.Target{}, errors.Annotate(err, "invalid target JSON")
	}
	return target, nil
}

func parseMigrationMode(value string) (pkgmigration.TaskMode, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", errors.New("mode is required, use --mode or provide it in interactive mode")
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
