// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	appconfig "github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/config/store"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestSecureStorageLifecycle(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	keyring.MockInit()
	file := filepath.Join(t.TempDir(), "config")
	viper.SetConfigFile(file)
	viper.SetConfigType("toml")
	viper.SetConfigPermissions(0600)
	require.NoError(t, viper.WriteConfig())
	h := &internal.Helper{IOStreams: iostream.Test()}
	create := CreateCmd(h)
	create.SetArgs([]string{"--profile-name", "secure", "--public-key", "public", "--private-key", "first-secret", "--secure-storage"})
	require.NoError(t, create.Execute())
	appconfig.SetActiveProfile("secure")
	data, err := os.ReadFile(file)
	require.NoError(t, err)
	require.NotContains(t, string(data), "first-secret")
	require.Equal(t, "public", appconfig.GetPublicKey())
	value, err := appconfig.GetPrivateKey()
	require.NoError(t, err)
	require.Equal(t, "first-secret", value)

	set := SetCmd(h)
	set.SetArgs([]string{"private-key", "second-secret"})
	require.NoError(t, set.Execute())
	value, err = appconfig.GetPrivateKey()
	require.NoError(t, err)
	require.Equal(t, "second-secret", value)
	data, err = os.ReadFile(file)
	require.NoError(t, err)
	require.NotContains(t, string(data), "second-secret")

	deleteCmd := DeleteCmd(h)
	deleteCmd.SetArgs([]string{"secure", "--force"})
	require.NoError(t, deleteCmd.Execute())
	_, err = store.GetAPIPrivateKey("secure")
	require.ErrorIs(t, err, keyring.ErrNotFound)
}

func TestSetSecureStorageRequiresPrivateKey(t *testing.T) {
	h := &internal.Helper{IOStreams: iostream.Test()}
	cmd := SetCmd(h)
	cmd.SetArgs([]string{"public-key", "public", "--secure-storage"})
	require.ErrorContains(t, cmd.Execute(), "only supported for private-key")
}
