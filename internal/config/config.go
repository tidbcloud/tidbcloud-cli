// Copyright 2022 PingCAP, Inc.
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

package config

import (
	"fmt"
	"sort"

	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/util"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

const (
	cliName       = "ticloud"
	cliNameInTiUP = "cloud"
	HomePath      = ".ticloud"
	DevVersion    = "dev"
	Repo          = "tidbcloud/tidbcloud-cli"

	Confirmed = "yes"
)

var (
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	CursorStyle  = FocusedStyle.Copy()
	CliName      = cliName
)

type Config struct {
	ActiveProfile string
}

func SetCliName(isUnderTiUP bool) {
	if isUnderTiUP {
		CliName = cliNameInTiUP
	} else {
		CliName = cliName
	}
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
