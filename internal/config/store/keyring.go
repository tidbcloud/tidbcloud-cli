package store

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
)

const namespace = "ticloud_access_token"

var ErrNotSupported = errors.New("Keyring is not supported on WSL")

func Get(profile string) (string, error) {
	if err := assertKeyringSupported(); err != nil {
		return "", err
	}
	val, err := keyring.Get(namespace, profile)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to load token: %s", err))
	}
	return val, nil
}

func Set(profile, token string) error {
	if err := assertKeyringSupported(); err != nil {
		return err
	}
	if err := keyring.Set(namespace, profile, token); err != nil {
		return errors.New(fmt.Sprintf("failed to set token: %s", err))
	}
	return nil
}

func Delete(profile string) error {
	if err := assertKeyringSupported(); err != nil {
		return err
	}
	if err := keyring.Delete(namespace, profile); err != nil {
		if !errors.Is(err, keyring.ErrNotFound) {
			return errors.New(fmt.Sprintf("failed to delete token: %s", err))
		}
	}
	return nil
}

func assertKeyringSupported() error {
	// Suggested check: https://github.com/microsoft/WSL/issues/423
	if f, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil && bytes.Contains(f, []byte("WSL")) {
		return ErrNotSupported
	}
	return nil
}
