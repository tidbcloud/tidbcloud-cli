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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal/config/store"
	"github.com/tidbcloud/tidbcloud-cli/internal/prop"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func credentialConfig(t *testing.T) string {
	t.Helper()
	viper.Reset()
	t.Cleanup(viper.Reset)
	file := filepath.Join(t.TempDir(), "config")
	viper.SetConfigFile(file)
	viper.SetConfigType("toml")
	viper.SetConfigPermissions(0600)
	require.NoError(t, viper.WriteConfig())
	return file
}

func TestPrivateKeyStorage(t *testing.T) {
	for _, profile := range []string{"default", "other", "'special'", `"special"`} {
		t.Run(profile, func(t *testing.T) {
			file := credentialConfig(t)
			keyring.MockInitWithError(errors.New("keyring unavailable"))
			p := &Profile{name: profile}
			require.NoError(t, SetPrivateKey(profile, "legacy-secret", false))
			require.NoError(t, viper.WriteConfig())
			value, err := p.GetPrivateKey()
			require.NoError(t, err)
			require.Equal(t, "legacy-secret", value)

			// An explicit secure write must not fall back to plaintext on failure.
			require.ErrorContains(t, SetPrivateKey(profile, "new-secret", true), "failed to save")
			require.Equal(t, "legacy-secret", viper.GetString(profile+"."+prop.PrivateKey))
			require.Empty(t, viper.GetString(profile+"."+prop.PrivateKeyStorage))

			keyring.MockInit()
			require.NoError(t, keyring.Set("ticloud_access_token", profile, "oauth-secret"))
			require.NoError(t, SetPrivateKey(profile, "keyring-secret", true))
			require.NoError(t, viper.WriteConfig())
			data, err := os.ReadFile(file)
			require.NoError(t, err)
			require.NotContains(t, string(data), "legacy-secret")
			require.NotContains(t, string(data), "keyring-secret")

			// Reload from disk to verify a later invocation can find the private key.
			viper.Reset()
			viper.SetConfigFile(file)
			viper.SetConfigType("toml")
			require.NoError(t, viper.ReadInConfig())
			value, err = p.GetPrivateKey()
			require.NoError(t, err)
			require.Equal(t, "keyring-secret", value)
			viper.Set(profile+"."+prop.PrivateKey, "stale-plaintext")
			value, err = p.GetPrivateKey()
			require.NoError(t, err)
			require.Equal(t, "keyring-secret", value)
			oauth, err := keyring.Get("ticloud_access_token", profile)
			require.NoError(t, err)
			require.Equal(t, "oauth-secret", oauth)
			unrelated, err := (&Profile{name: "unrelated"}).GetPrivateKey()
			require.NoError(t, err)
			require.Empty(t, unrelated)

			// Once opted in, an ordinary update must remain in the keyring.
			require.NoError(t, SetPrivateKey(profile, "rotated-secret", false))
			require.NoError(t, viper.WriteConfig())
			value, err = p.GetPrivateKey()
			require.NoError(t, err)
			require.Equal(t, "rotated-secret", value)
			data, err = os.ReadFile(file)
			require.NoError(t, err)
			require.NotContains(t, string(data), "rotated-secret")

			require.NoError(t, store.DeleteAPIPrivateKey(profile))
			_, err = p.GetPrivateKey()
			require.ErrorIs(t, err, keyring.ErrNotFound)
			// A missing keyring entry must not silently disable API-key authentication.
			require.ErrorContains(t, err, "failed to load API private key")
			viper.Set(profile+"."+prop.PrivateKey, "stale-plaintext")
			_, err = p.GetPrivateKey()
			require.ErrorIs(t, err, keyring.ErrNotFound)
		})
	}
}

func TestAccessTokenStoragePriority(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		secure    string
		loadErr   error
		want      string
	}{
		{name: "prefer keyring", plaintext: "stale", secure: "secure", want: "secure"},
		{name: "insecure storage", plaintext: "plaintext", loadErr: keyring.ErrNotFound, want: "plaintext"},
		{name: "unsupported keyring", plaintext: "plaintext", loadErr: store.ErrNotSupported, want: "plaintext"},
		{name: "locked keyring preserves legacy access", plaintext: "plaintext", loadErr: errors.New("locked"), want: "plaintext"},
		{name: "empty keyring", plaintext: "plaintext", want: "plaintext"},
		{name: "missing token", loadErr: keyring.ErrNotFound},
		{name: "keyring error without plaintext", loadErr: errors.New("locked")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			credentialConfig(t)
			viper.Set("test."+prop.AccessToken, tc.plaintext)
			p := &Profile{name: "test"}
			got, err := p.getAccessToken(func(profile string) (string, error) {
				require.Equal(t, "test", profile)
				return tc.secure, tc.loadErr
			})
			require.Equal(t, tc.want, got)
			if tc.want == "" {
				require.ErrorIs(t, err, tc.loadErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteAccessTokenClearsViperOverrides(t *testing.T) {
	file := credentialConfig(t)
	keyring.MockInit()
	p := &Profile{name: "target"}
	// Values written during this process live in Viper's override map.
	viper.Set("target."+prop.AccessToken, "old-secret")
	viper.Set("target."+prop.TokenType, "Bearer")
	viper.Set("target."+prop.TokenExpiredAt, time.Now().Add(time.Hour))
	viper.Set("unrelated."+prop.AccessToken, "keep-secret")
	require.NoError(t, viper.WriteConfig())
	require.NoError(t, p.DeleteAccessToken())
	require.Empty(t, viper.GetString("target."+prop.AccessToken))
	require.Equal(t, "keep-secret", viper.GetString("unrelated."+prop.AccessToken))
	require.NoError(t, viper.WriteConfig())
	data, err := os.ReadFile(file)
	require.NoError(t, err)
	require.NotContains(t, string(data), "old-secret")
}

func TestSaveAccessTokenUsesReceiverProfile(t *testing.T) {
	credentialConfig(t)
	keyring.MockInit()
	previous := activeProfile
	activeProfile = &Profile{name: "active"}
	t.Cleanup(func() { activeProfile = previous })
	viper.Set("active."+prop.AccessToken, "keep-secret")
	viper.Set("target."+prop.AccessToken, "old-secret")
	require.NoError(t, viper.WriteConfig())
	p := &Profile{name: "target"}
	require.NoError(t, p.SaveAccessToken(time.Now().Add(time.Hour), "Bearer", "new-secret", true))
	require.Equal(t, "keep-secret", viper.GetString("active."+prop.AccessToken))
	require.Equal(t, "new-secret", viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.AccessToken)))
}
