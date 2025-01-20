// Copyright 2025 PingCAP, Inc.
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
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/mock"

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
		sender: mockSender,
		cmd:    cmd,
	}
	mockSender.On("SendEvent", mockUtils.Anything).Return(nil)

	err := tracker.trackCommand(TrackOptions{})
	mockSender.AssertExpectations(t)
	a.NoError(err)
}

func TestTrackCommandWithError(t *testing.T) {
	mockSender := new(mock.EventsSender)

	a := require.New(t)

	cmd := &cobra.Command{
		Use: "test-command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("test command error")
		},
	}
	cmd.SetArgs([]string{})
	err := cmd.ExecuteContext(NewTelemetryContext(context.Background()))
	a.Error(err)

	tracker := &tracker{
		sender: mockSender,
		cmd:    cmd,
	}
	mockSender.On("SendEvent", mockUtils.Anything).Return(nil)

	err = tracker.trackCommand(TrackOptions{
		Err: err,
	})
	mockSender.AssertExpectations(t)
	a.NoError(err)
}

func TestTrackCommandWithSendError(t *testing.T) {
	mockSender := new(mock.EventsSender)

	a := require.New(t)
	mockSender.On("SendEvent", mockUtils.Anything).Return(errors.New("test send error"))

	cmd := &cobra.Command{
		Use: "test-command",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.SetArgs([]string{})
	err := cmd.ExecuteContext(NewTelemetryContext(context.Background()))
	a.NoError(err)

	tracker := &tracker{
		sender: mockSender,
		cmd:    cmd,
	}

	err = tracker.trackCommand(TrackOptions{
		Err: err,
	})
	mockSender.AssertExpectations(t)
	a.Error(err)
}
