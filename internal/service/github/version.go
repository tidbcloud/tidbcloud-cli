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

package github

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tidbcloud-cli/internal/config"
	ver "tidbcloud-cli/internal/version"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-version"
	"github.com/juju/errors"
	"gopkg.in/yaml.v3"
)

// ReleaseInfo stores information about a release
type ReleaseInfo struct {
	Version string `json:"tag_name"`
}

type StateEntry struct {
	CheckedForUpdateAt time.Time   `yaml:"checked_for_update_at"`
	LatestRelease      ReleaseInfo `yaml:"latest_release"`
}

// CheckForUpdate checks for updates and returns the latest release info if there is a newer version.
// For checking after every command, we should use stateEntry to avoid checking too frequently.
// For checking when using `update`, we should use forceCheck to ignore the stateEntry.
func CheckForUpdate(repo string, withRateLimit bool) (*ReleaseInfo, error) {
	home, _ := os.UserHomeDir()
	stateFilePath := filepath.Join(home, config.HomePath, "state.yml")
	if withRateLimit {
		stateEntry, _ := getStateEntry(stateFilePath)
		if stateEntry != nil && time.Since(stateEntry.CheckedForUpdateAt).Hours() < 24 {
			return nil, nil
		}
	}

	releaseInfo, err := getLatestReleaseInfo(repo)
	if err != nil {
		return nil, err
	}

	if withRateLimit {
		err = setStateEntry(stateFilePath, time.Now(), *releaseInfo)
		if err != nil {
			return nil, err
		}
	}

	if versionGreaterThan(releaseInfo.Version, ver.Version) {
		return releaseInfo, nil
	}

	return nil, nil
}

func getStateEntry(stateFilePath string) (*StateEntry, error) {
	content, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, err
	}

	var stateEntry StateEntry
	err = yaml.Unmarshal(content, &stateEntry)
	if err != nil {
		return nil, err
	}

	return &stateEntry, nil
}

func getLatestReleaseInfo(repo string) (*ReleaseInfo, error) {
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	response, err := client.
		R().
		Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/main/latest-version", repo))
	if err != nil {
		return nil, errors.Trace(err)
	}

	body := string(response.Body())
	_, err = version.NewVersion(body)
	if err != nil {
		return nil, err
	}
	latestRelease := ReleaseInfo{Version: body}
	if response.IsSuccess() {
		return &latestRelease, nil
	} else {
		return nil, errors.Errorf("failed to get latest release info: %s", response.Status()+" "+body)
	}
}

func setStateEntry(stateFilePath string, t time.Time, r ReleaseInfo) error {
	data := StateEntry{CheckedForUpdateAt: t, LatestRelease: r}
	content, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(stateFilePath), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(stateFilePath, content, 0600)
	return err
}

func versionGreaterThan(v, w string) bool {
	vv, ve := version.NewVersion(v)
	vw, we := version.NewVersion(w)

	return ve == nil && we == nil && vv.GreaterThan(vw)
}
