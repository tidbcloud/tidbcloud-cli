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

package telemetry

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"tidbcloud-cli/internal/config"

	"github.com/pingcap/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	cacheFilename           = "telemetry"
	dirPermissions          = 0700
	filePermissions         = 0600
	defaultMaxCacheFileSize = 500_000 // 500KB
)

type tracker struct {
	fs               afero.Fs
	maxCacheFileSize int64
	cacheDir         string
	cmd              *cobra.Command
	sender           EventsSender
	args             []string
	installer        *string
}

func newTracker(cmd *cobra.Command, args []string) (*tracker, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	cacheDir = filepath.Join(cacheDir, config.CliName)

	return &tracker{
		fs:               afero.NewOsFs(),
		maxCacheFileSize: defaultMaxCacheFileSize,
		cacheDir:         cacheDir,
		cmd:              cmd,
		sender:           NewSender(),
		args:             args,
		installer:        readInstaller(),
	}, nil
}

func (t *tracker) defaultCommandOptions() []eventOpt {
	return []eventOpt{withCommandPath(t.cmd), WithInteractive(t.cmd), withHelpCommand(t.cmd, t.args), withFlags(t.cmd), withVersion(), withOS(), withAuthMethod(), withProjectID(t.cmd), withTerminal(), withInstaller(t.installer)}
}

func (t *tracker) trackCommand(data TrackOptions) error {
	options := append(t.defaultCommandOptions(), withDuration(t.cmd))

	if data.Err != nil {
		options = append(options, withError(data.Err))
	}
	event := newEvent(options...)
	events, err := t.read()
	if err != nil {
		log.Debug("telemetry: failed to read cache", zap.Error(err))
	}
	events = append(events, event)
	err = t.sender.SendEvents(events)
	if err != nil {
		log.Debug("telemetry: failed to send events", zap.Error(err))
		return t.save(event)
	}
	return t.remove()
}

func (t *tracker) openCacheFile() (afero.File, error) {
	exists, err := afero.DirExists(t.fs, t.cacheDir)
	if err != nil {
		return nil, err
	}
	if !exists {
		if mkdirError := t.fs.MkdirAll(t.cacheDir, dirPermissions); mkdirError != nil {
			return nil, mkdirError
		}
	}
	filename := filepath.Join(t.cacheDir, cacheFilename)
	exists, err = afero.Exists(t.fs, filename)
	if err != nil {
		return nil, err
	}
	if exists {
		info, statError := t.fs.Stat(filename)
		if statError != nil {
			return nil, statError
		}
		size := info.Size()
		if size > t.maxCacheFileSize {
			return nil, errors.New("telemetry cache file too large")
		}
	}
	file, err := t.fs.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, filePermissions)
	return file, err
}

// Append a single event to the cache file.
func (t *tracker) save(event Event) error {
	file, err := t.openCacheFile()
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = file.Write(data)
	return err
}

// Read all events in the cache file.
func (t *tracker) read() ([]Event, error) {
	initialSize := 100
	events := make([]Event, 0, initialSize)
	filename := filepath.Join(t.cacheDir, cacheFilename)
	exists, err := afero.Exists(t.fs, filename)
	if err != nil {
		return events, err
	}
	if exists {
		file, err := t.fs.Open(filename)
		if err != nil {
			return events, err
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		for decoder.More() {
			var event Event
			err = decoder.Decode(&event)
			if err != nil {
				return events, err
			}
			events = append(events, event)
		}
	}
	return events, nil
}

// Removes the cache file.
func (t *tracker) remove() error {
	filename := filepath.Join(t.cacheDir, cacheFilename)
	exists, err := afero.Exists(t.fs, filename)
	if exists && err == nil {
		return t.fs.Remove(filename)
	}
	return err
}
