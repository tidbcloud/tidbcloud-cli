// Copyright 2026 PingCAP, Inc.
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

package cli

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
)

func TestMissingCredentialsJSONDoesNotPrintHumanPreamble(t *testing.T) {
	oldOutput := color.Output
	defer func() { color.Output = oldOutput; viper.Reset() }()
	viper.Reset()
	config.SetActiveProfile("missing-credentials-test")
	var out bytes.Buffer
	color.Output = &out
	_, _, _, err := loadCredentials(true)
	require.ErrorIs(t, err, ErrMissingCredentials)
	require.Empty(t, out.String())
	_, _, _, err = loadCredentials(false)
	require.ErrorIs(t, err, ErrMissingCredentials)
	require.Contains(t, out.String(), "auth login", "legacy/human guidance is retained")
}
