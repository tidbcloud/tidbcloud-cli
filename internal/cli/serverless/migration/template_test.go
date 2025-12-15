// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package migration

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	pkgmigration "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/migration"
)

func TestTemplateCmd(t *testing.T) {
	assert := require.New(t)
	if err := os.Setenv("NO_COLOR", "true"); err != nil {
		t.Error(err)
	}

	h := testHelper()

	tests := []struct {
		name         string
		args         []string
		err          error
		stdoutString string
		stderrString string
	}{
		{
			name:         "render ALL template",
			args:         []string{"--mode", "all"},
			stdoutString: fmt.Sprintf("%s\n%s\n", definitionTemplates[pkgmigration.TASKMODE_ALL].heading, migrationDefinitionAllTemplate),
		},
		{
			name:         "render INCREMENTAL template",
			args:         []string{"--mode", "incremental"},
			stdoutString: fmt.Sprintf("%s\n%s\n", definitionTemplates[pkgmigration.TASKMODE_INCREMENTAL].heading, migrationDefinitionIncrementalTemplate),
		},
		{
			name: "invalid mode",
			args: []string{"--mode", "invalid"},
			err:  fmt.Errorf("unknown mode %q, allowed values: %s", "invalid", strings.Join(allowedTemplateModeStrings(), ", ")),
		},
		{
			name: "missing mode flag",
			args: []string{},
			err:  fmt.Errorf("required flag(s) \"%s\" not set", flag.MigrationMode),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := TemplateCmd(h)
			h.IOStreams.Out.(*bytes.Buffer).Reset()
			h.IOStreams.Err.(*bytes.Buffer).Reset()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if tt.err != nil {
				assert.EqualError(err, tt.err.Error())
			} else {
				assert.NoError(err)
			}

			assert.Equal(tt.stdoutString, h.IOStreams.Out.(*bytes.Buffer).String())
			assert.Equal(tt.stderrString, h.IOStreams.Err.(*bytes.Buffer).String())
		})
	}
}

// testHelper creates a helper with no client dependency for template command tests.
func testHelper() *internal.Helper {
	return &internal.Helper{
		Client:        func() (cloud.TiDBCloudClient, error) { return nil, nil },
		QueryPageSize: 1,
		IOStreams:     iostream.Test(),
	}
}
