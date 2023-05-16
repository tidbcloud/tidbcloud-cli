// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func (m *MySQLHelperImpl) DumpFromMySQL(arg string) error {
	c1 := exec.Command("powershell", "/C", arg) //nolint:gosec
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
