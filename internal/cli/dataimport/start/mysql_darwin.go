//go:build darwin
// +build darwin

package start

import (
	"bytes"
	"strings"

	"github.com/pingcap/errors"
	exec "golang.org/x/sys/execabs"
)

func (m *MySQLHelperImpl) DumpFromMySQL(args []string) error {
	c1 := exec.Command("sh", "-c", strings.Join(args, " ")) //nolint:gosec
	var stderr bytes.Buffer
	c1.Stderr = &stderr

	err := c1.Run()
	if err != nil {
		return errors.Annotate(err, stderr.String())
	}

	return nil
}

func (m *MySQLHelperImpl) ImportToServerless(sqlCacheFile string, connectionString string) error {
	var stderr bytes.Buffer

	c1 := exec.Command("sh", "-c", connectionString+" < "+sqlCacheFile) //nolint:gosec
	stderr = bytes.Buffer{}
	c1.Stderr = &stderr

	err := c1.Run()
	if err != nil {
		return errors.Annotate(err, stderr.String())
	}

	return nil
}
