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

package store

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const namespace = "ticloud_access_token"

var ErrNotSupported = errors.New("keyring is not supported, you can add `--insecure-storage true` to save token to config file instead.\n" +
	"or see https://github.com/zalando/go-keyring#dependencies for more details of keyring dependencies")

func Get(profile string) (string, error) {
	if err := AssertKeyringSupported(); err != nil {
		return "", err
	}
	val, err := keyring.Get(namespace, profile)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", err
		}
		return "", fmt.Errorf("failed to load token: %s", err)
	}
	return val, nil
}

func Set(profile, token string) error {
	if err := AssertKeyringSupported(); err != nil {
		return err
	}
	if err := keyring.Set(namespace, profile, token); err != nil {
		return fmt.Errorf("failed to set token: %s", err)
	}
	return nil
}

func Delete(profile string) error {
	if err := AssertKeyringSupported(); err != nil {
		return err
	}
	if err := keyring.Delete(namespace, profile); err != nil {
		if !errors.Is(err, keyring.ErrNotFound) {
			return fmt.Errorf("failed to delete token: %s", err)
		}
	}
	return nil
}
