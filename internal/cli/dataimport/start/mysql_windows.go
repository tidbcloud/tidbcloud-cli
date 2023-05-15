//go:build windows
// +build windows

package start

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tidbcloud-cli/internal/config"

	"github.com/pingcap/errors"
	exec "golang.org/x/sys/execabs"
)

func (m *MySQLHelperImpl) DumpFromMySQL(args []string) error {
	fmt.Println(strings.Join(args, " "))
	c1 := exec.Command("powershell", "/C", strings.Join(args, " ")) //nolint:gosec
	var stderr bytes.Buffer
	c1.Stderr = &stderr

	err := c1.Run()
	if err != nil {
		return errors.Annotate(err, stderr.String())
	}

	return nil
}

func (m *MySQLHelperImpl) ImportToServerless(sqlCacheFile string, connectionString string) error {
	home, _ := os.UserHomeDir()
	caFile := filepath.Join(home, config.HomePath, "isrgrootx1.pem")
	_, err := os.Stat(caFile)
	if os.IsNotExist(err) {
		err := m.DownloadCaFile(caFile)
		if err != nil {
			return err
		}
	}
	connectionString = strings.Replace(connectionString, "<path_to_ca_cert>", caFile, -1)
	fmt.Println(connectionString)

	var stderr bytes.Buffer
	// PowerShell not support "<" operator, so we use "-e" to execute the command
	c1 := exec.Command("powershell", "/C", fmt.Sprintf("%s -e \"source %s\"", connectionString, sqlCacheFile)) //nolint:gosec
	stderr = bytes.Buffer{}
	c1.Stderr = &stderr

	err = c1.Run()
	if err != nil {
		return errors.Annotate(err, stderr.String())
	}

	return nil
}
