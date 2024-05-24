// Copyright 2024 PingCAP, Inc.
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

//go:build linux

package store

import (
	"bytes"
	"os"

	"github.com/pingcap/log"
	ss "github.com/zalando/go-keyring/secret_service"
	"go.uber.org/zap"
)

func AssertKeyringSupported() error {
	// Suggested check: https://github.com/microsoft/WSL/issues/423
	if f, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		if bytes.Contains(f, []byte("WSL")) || bytes.Contains(f, []byte("Microsoft")) {
			return ErrNotSupported
		}
	}

	// Check gnome-keyring, see https://github.com/zalando/go-keyring/blob/v0.2.4/keyring_unix.go#L16
	svc, err := ss.NewSecretService()
	if err != nil {
		log.Debug("failed to create secret service", zap.Error(err))
		return ErrNotSupported
	}

	session, err := svc.OpenSession()
	if err != nil {
		log.Debug("failed to open dbus session", zap.Error(err))
		return ErrNotSupported
	}
	defer svc.Close(session)

	return nil
}
