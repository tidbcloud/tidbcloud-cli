// Copyright 2023 PingCAP, Inc.
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
	"fmt"
	"os"
	"sort"
	"time"

	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/internal/version"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml"
	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"go.uber.org/zap"
)

var (
	activeProfile = &Profile{}
)

const keyringService = "ticloud_access_token"

type Profile struct {
	name string
}

func ValidateProfile(profileName string) error {
	profiles, err := GetAllProfiles()
	if err != nil {
		return err
	}

	if !util.ElemInSlice(profiles, profileName) {
		return fmt.Errorf("profile %s not found", profileName)
	}

	return nil
}

func GetAllProfiles() ([]string, error) {
	s := viper.AllSettings()
	keys := make([]string, 0, len(s))
	// Profile names and global properties are at the same level in the config file, filter out global properties.
	for k := range s {
		if !util.ElemInSlice(prop.GlobalProperties(), k) {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)
	return keys, nil
}

func SetActiveProfile(name string) {
	if name == "" {
		name = "default"
		viper.Set(prop.CurProfile, name)
		err := viper.WriteConfig()
		if err != nil {
			log.Debug("failed to save current profile to config file", zap.Error(err))
		}
	}
	activeProfile.name = name
}

func ActiveProfileName() string {
	return activeProfile.name
}

func TelemetryEnabled() bool { return activeProfile.TelemetryEnabled() }
func (p *Profile) TelemetryEnabled() bool {
	// If we're running a dev build, we don't want to send telemetry.
	if version.Version == "dev" {
		return false
	}

	key := fmt.Sprintf("%s.%s", p.name, prop.TelemetryEnabled)
	// Telemetry is disabled by default.
	if !viper.IsSet(key) {
		return false
	}
	return viper.GetBool(key)
}

func GetPublicKey() (publicKey string) { return activeProfile.GetPublicKey() }
func (p *Profile) GetPublicKey() (publicKey string) {
	publicKey = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.PublicKey))
	return
}

func GetPrivateKey() (privateKey string) { return activeProfile.GetPrivateKey() }
func (p *Profile) GetPrivateKey() (privateKey string) {
	privateKey = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.PrivateKey))
	return
}

func GetApiUrl() (apiUrl string) { return activeProfile.GetApiUrl() }
func (p *Profile) GetApiUrl() (apiUrl string) {
	apiUrl = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.ApiUrl))
	return
}

func GetServerlessEndpoint() (apiUrl string) { return activeProfile.GetServerlessEndpoint() }
func (p *Profile) GetServerlessEndpoint() (newApiUrl string) {
	newApiUrl = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.ServerlessEndpoint))
	return
}

func GetOAuthEndpoint() (apiUrl string) { return activeProfile.GetOAuthEndpoint() }
func (p *Profile) GetOAuthEndpoint() (newApiUrl string) {
	newApiUrl = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.OAuthEndpoint))
	if newApiUrl == "" {
		return "https://oauth.tidbcloud.com"
	}
	return
}

func GetOAuthClientID() (apiUrl string) { return activeProfile.GetOAuthClientID() }
func (p *Profile) GetOAuthClientID() (clientID string) {
	clientID = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.OAuthClientID))
	if clientID == "" {
		return ""
	}
	return
}

func GetOAuthClientSecret() (apiUrl string) { return activeProfile.GetOAuthClientSecret() }
func (p *Profile) GetOAuthClientSecret() (clientID string) {
	clientID = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.OAuthClientSecret))
	if clientID == "" {
		return ""
	}
	return
}

func SaveAccessToken(expireAt time.Time, tokenType string, token string, insecureStorageUsed bool) error {
	return activeProfile.SaveAccessToken(expireAt, tokenType, token, insecureStorageUsed)
}
func (p *Profile) SaveAccessToken(expireAt time.Time, tokenType string, token string, insecureStorageUsed bool) error {
	if !insecureStorageUsed {
		err := keyring.Set(keyringService, p.name, token)
		if err != nil {
			log.Debug("failed to save access token to keyring", zap.Error(err))
			color.Yellow("failed to save access token to keyring, save to config file instead")
			// If failed to save to keyring, fallback to save to config file.
			insecureStorageUsed = true
		}
	}

	if insecureStorageUsed {
		viper.Set(fmt.Sprintf("%s.%s", p.name, prop.AccessToken), token)
	}
	viper.Set(fmt.Sprintf("%s.%s", p.name, prop.TokenExpiredAt), expireAt)
	viper.Set(fmt.Sprintf("%s.%s", p.name, prop.TokenType), tokenType)
	err := viper.WriteConfig()
	if err != nil {
		return errors.Annotate(err, "failed to save token info to config file")
	}

	return nil
}

func GetAccessToken() (string, error) {
	return activeProfile.GetAccessToken()
}
func (p *Profile) GetAccessToken() (string, error) {
	token := viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.AccessToken))
	if token != "" {
		return token, nil
	}

	return keyring.Get(keyringService, p.name)
}

func ValidateToken() error {
	return activeProfile.ValidateToken()
}
func (p *Profile) ValidateToken() error {
	tokenExpiredAt := viper.GetTime(fmt.Sprintf("%s.%s", p.name, prop.TokenExpiredAt))
	if tokenExpiredAt.Before(time.Now()) {
		err := DeleteAccessToken()
		log.Debug("failed to delete access token", zap.Error(err))

		return fmt.Errorf("access token expired, please login again")
	}

	return nil
}

func DeleteAccessToken() error {
	return activeProfile.DeleteAccessToken()
}
func (p *Profile) DeleteAccessToken() error {
	settings := viper.AllSettings()
	t, err := toml.TreeFromMap(settings)
	if err != nil {
		return errors.Trace(err)
	}

	if t.Has(fmt.Sprintf("%s.%s", p.name, prop.AccessToken)) {
		err := t.Delete(fmt.Sprintf("%s.%s", p.name, prop.AccessToken))
		if err != nil {
			return err
		}
	} else {
		err := keyring.Delete(keyringService, p.name)
		if err != nil {
			return err
		}
	}

	err = t.Delete(fmt.Sprintf("%s.%s", p.name, prop.TokenType))
	if err != nil {
		return err
	}
	err = t.Delete(fmt.Sprintf("%s.%s", p.name, prop.TokenExpiredAt))
	if err != nil {
		return err
	}

	fs := afero.NewOsFs()
	file, err := fs.OpenFile(viper.ConfigFileUsed(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Trace(err)
	}

	defer file.Close()

	s := t.String()
	_, err = file.WriteString(s)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}
