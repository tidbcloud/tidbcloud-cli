// Copyright 2026 PingCAP, Inc.
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

// Keep API private keys separate from OAuth access tokens.
const apiPrivateKeyNamespace = "ticloud_api_private_key"

func GetAPIPrivateKey(profile string) (string, error) {
	value, err := keyring.Get(apiPrivateKeyNamespace, profile)
	if err != nil {
		return "", fmt.Errorf("failed to load API private key from credential store: %w", err)
	}
	return value, nil
}

func SetAPIPrivateKey(profile, value string) error {
	if err := keyring.Set(apiPrivateKeyNamespace, profile, value); err != nil {
		return fmt.Errorf("failed to save API private key to credential store: %w", err)
	}
	return nil
}

func DeleteAPIPrivateKey(profile string) error {
	if err := keyring.Delete(apiPrivateKeyNamespace, profile); err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("failed to delete API private key from credential store: %w", err)
	}
	return nil
}
