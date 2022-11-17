package config

import (
	"fmt"
	"sort"

	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/spf13/viper"
)

type Config struct {
	ActiveProfile string
}

func ValidateProfile(profileName string) error {
	profiles, err := GetAllProfiles()
	if err != nil {
		return err
	}

	if !util.StringInSlice(profiles, profileName) {
		return fmt.Errorf("profile %s not found", profileName)
	}

	return nil
}

func GetAllProfiles() ([]string, error) {
	s := viper.AllSettings()
	keys := make([]string, 0, len(s))
	for k := range s {
		if !util.StringInSlice(prop.GlobalProperties(), k) {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)
	return keys, nil
}
