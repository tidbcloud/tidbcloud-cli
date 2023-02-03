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
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/mock"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	mockUtils "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTrackCommand(t *testing.T) {
	mockSender := new(mock.EventsSender)
	a := require.New(t)

	cmd := &cobra.Command{
		Use: "test-command",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.SetArgs([]string{})
	_ = cmd.ExecuteContext(NewTelemetryContext(context.Background()))

	tracker := &tracker{
		fs:               afero.NewMemMapFs(),
		maxCacheFileSize: defaultMaxCacheFileSize,
		sender:           mockSender,
		cmd:              cmd,
	}
	mockSender.On("SendEvents", mockUtils.Anything).Return(nil)

	err := tracker.trackCommand(TrackOptions{})
	mockSender.AssertExpectations(t)
	a.NoError(err)
}

func TestTrackCommandWithError(t *testing.T) {
	mockSender := new(mock.EventsSender)

	a := require.New(t)
	cacheDir, err := os.MkdirTemp(os.TempDir(), config.CliName+"*")
	a.NoError(err)

	cmd := &cobra.Command{
		Use: "test-command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("test command error")
		},
	}
	cmd.SetArgs([]string{})
	err = cmd.ExecuteContext(NewTelemetryContext(context.Background()))
	a.Error(err)

	tracker := &tracker{
		fs:               afero.NewMemMapFs(),
		maxCacheFileSize: defaultMaxCacheFileSize,
		cacheDir:         cacheDir,
		sender:           mockSender,
		cmd:              cmd,
	}
	mockSender.On("SendEvents", mockUtils.Anything).Return(nil)

	err = tracker.trackCommand(TrackOptions{
		Err: err,
	})
	mockSender.AssertExpectations(t)
	a.NoError(err)
}

func TestTrackCommandWithSendError(t *testing.T) {
	mockSender := new(mock.EventsSender)

	a := require.New(t)
	cacheDir, err := os.MkdirTemp(os.TempDir(), config.CliName+"*")
	a.NoError(err)
	mockSender.On("SendEvents", mockUtils.Anything).Return(errors.New("test send error"))

	cmd := &cobra.Command{
		Use: "test-command",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.SetArgs([]string{})
	err = cmd.ExecuteContext(NewTelemetryContext(context.Background()))
	a.NoError(err)

	tracker := &tracker{
		fs:               afero.NewMemMapFs(),
		maxCacheFileSize: defaultMaxCacheFileSize,
		cacheDir:         cacheDir,
		sender:           mockSender,
		cmd:              cmd,
	}

	err = tracker.trackCommand(TrackOptions{
		Err: err,
	})
	mockSender.AssertExpectations(t)
	a.NoError(err)

	// Verify that the file exists
	filename := filepath.Join(cacheDir, cacheFilename)
	info, statError := tracker.fs.Stat(filename)
	a.NoError(statError)
	// Verify the file name
	a.Equal(info.Name(), cacheFilename)
	// Verify that the file contains some data
	var minExpectedSize int64 = 10
	a.True(info.Size() > minExpectedSize)
}

func TestSave(t *testing.T) {
	a := require.New(t)
	cacheDir, err := os.MkdirTemp(os.TempDir(), config.CliName+"*")
	a.NoError(err)

	tracker := &tracker{
		fs:               afero.NewMemMapFs(),
		maxCacheFileSize: defaultMaxCacheFileSize,
		cacheDir:         cacheDir,
	}

	var properties = map[string]interface{}{
		"command": "mock-command",
	}
	var event = Event{
		Timestamp:  time.Now(),
		Source:     config.CliName,
		Properties: properties,
	}
	a.NoError(tracker.save(event))
	// Verify that the file exists
	filename := path.Join(cacheDir, cacheFilename)
	info, statError := tracker.fs.Stat(filename)
	a.NoError(statError)
	// Verify the file name
	a.Equal(info.Name(), cacheFilename)
	// Verify that the file contains some data
	var minExpectedSize int64 = 10
	a.True(info.Size() > minExpectedSize)
}

func TestSaveOverMaxCacheFileSize(t *testing.T) {
	a := require.New(t)
	cacheDir, err := os.MkdirTemp(os.TempDir(), config.CliName+"*")
	a.NoError(err)

	tracker := &tracker{
		fs:               afero.NewMemMapFs(),
		maxCacheFileSize: 10, // 10 bytes
		cacheDir:         cacheDir,
	}

	var properties = map[string]interface{}{
		"command": "mock-command",
	}
	var event = Event{
		Timestamp:  time.Now(),
		Source:     config.CliName,
		Properties: properties,
	}

	// First save will work as the cache file will be new
	a.NoError(tracker.save(event))
	// Second save should fail as the file will be larger than 10 bytes
	a.Error(tracker.save(event))
}

func TestOpenCacheFile(t *testing.T) {
	a := require.New(t)
	cacheDir, err := os.MkdirTemp(os.TempDir(), config.CliName+"*")
	a.NoError(err)

	tracker := &tracker{
		fs:               afero.NewMemMapFs(),
		maxCacheFileSize: 10, // 10 bytes
		cacheDir:         cacheDir,
	}

	_, err = tracker.openCacheFile()
	a.NoError(err)
	// Verify that the file exists
	filename := path.Join(cacheDir, cacheFilename)
	info, statError := tracker.fs.Stat(filename)
	a.NoError(statError)
	// Verify the file name
	a.Equal(info.Name(), cacheFilename)
	// Verify that the file is empty
	var expectedSize int64 // The nil value is zero
	a.Equal(info.Size(), expectedSize)
}
