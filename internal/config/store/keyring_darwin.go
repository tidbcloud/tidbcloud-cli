//go:build darwin

package store

import (
	"errors"
	"os/exec"
)

const execPathKeychain = "/usr/bin/security"

func assertKeyringSupported() error {
	if errors.Is(exec.Command(execPathKeychain).Run(), exec.ErrNotFound) {
		return ErrNotSupported
	}
	return nil
}
