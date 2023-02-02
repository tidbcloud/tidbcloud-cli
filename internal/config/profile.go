package config

import (
	"fmt"
	"sort"

	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/internal/version"

	"github.com/spf13/viper"
)

var (
	activeProfile = &Profile{}
)

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
	for k := range s {
		if !util.ElemInSlice(prop.GlobalProperties(), k) {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)
	return keys, nil
}

func SetActiveProfile(name string) {
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
	if !viper.IsSet(key) {
		return true
	}
	return viper.GetBool(key)
}

func CheckAuth() error { return activeProfile.CheckAuth() }
func (p *Profile) CheckAuth() error {
	if p.name == "" {
		return fmt.Errorf("no active profile for auth, please use `config create` to create one")
	}

	publicKey := viper.Get(fmt.Sprintf("%s.%s", p.name, prop.PublicKey))
	if publicKey == nil {
		return fmt.Errorf("no public key configured for auth, please use `set %s <publicKey>` to set one", prop.PublicKey)
	}

	privateKey := viper.Get(fmt.Sprintf("%s.%s", p.name, prop.PrivateKey))
	if privateKey == nil {
		return fmt.Errorf("no private key configured for auth, please use `set %s <privateKey>` to set one", prop.PrivateKey)
	}

	return nil
}

func GePublicKey() (publicKey string) { return activeProfile.GePublicKey() }
func (p *Profile) GePublicKey() (publicKey string) {
	publicKey = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.PublicKey))
	return
}

func GePrivateKey() (privateKey string) { return activeProfile.GePrivateKey() }
func (p *Profile) GePrivateKey() (privateKey string) {
	privateKey = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.PrivateKey))
	return
}

func GetApiUrl() (apiUrl string) { return activeProfile.GetApiUrl() }
func (p *Profile) GetApiUrl() (apiUrl string) {
	apiUrl = viper.GetString(fmt.Sprintf("%s.%s", p.name, prop.ApiUrl))
	return
}
