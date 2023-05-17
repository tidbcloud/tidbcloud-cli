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

package start

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"tidbcloud-cli/internal/config"

	"github.com/pingcap/errors"
	exec "golang.org/x/sys/execabs"
)

type MySQLHelperImpl struct {
}

func (m *MySQLHelperImpl) GenerateSqlCachePath() string {
	home, _ := os.UserHomeDir()
	sqlCacheFile := filepath.Join(home, config.HomePath, ".cache", "dump-"+time.Now().Format("2006-01-02T15-04-05")+".sql")
	return sqlCacheFile
}

func (m *MySQLHelperImpl) DownloadCaFile(caFile string) error {
	// 下载文件的 URL
	url := "https://letsencrypt.org/certs/isrgrootx1.pem"

	// 创建 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return errors.Annotate(err, "Failed to download ca file")
	}
	defer resp.Body.Close()

	// 创建文件并打开
	file, err := os.Create(caFile)
	if err != nil {
		return errors.Annotate(err, "Failed to create ca file")
	}
	defer file.Close()

	// 将 HTTP 响应体复制到文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Annotate(err, "Failed to copy ca file")
	}

	return nil
}

// CheckMySQLClient checks whether the 'mysql' client exists and is configured in $PATH
func (m *MySQLHelperImpl) CheckMySQLClient() error {
	_, err := exec.LookPath("mysql")
	if err == nil {
		return nil
	}

	msg := "couldn't find the 'mysql' command-line tool required to run this command."

	switch runtime.GOOS {
	case "darwin":
		if HasHomebrew() {
			return fmt.Errorf("%s\nTo install, run: brew install mysql-client", msg)
		}
	}

	return fmt.Errorf("%s\nPlease install it and add to $PATH", msg)
}

// HasHomebrew check whether the user has installed brew
func HasHomebrew() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func (m *MySQLHelperImpl) DumpFromMySQL(c1 *exec.Cmd, sqlCacheFile string) error {
	var stderr bytes.Buffer
	c1.Stderr = &stderr
	output, err := os.Create(sqlCacheFile)
	if err != nil {
		return err
	}
	defer output.Close()
	c1.Stdout = output

	err = c1.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return err
	}

	return nil
}

func (m *MySQLHelperImpl) ImportToServerless(c1 *exec.Cmd, sqlCacheFile string) error {
	var stderr bytes.Buffer
	c1.Stderr = &stderr
	input, err := os.Open(sqlCacheFile)
	if err != nil {
		return err
	}
	defer input.Close()
	c1.Stdin = input

	err = c1.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return err
	}

	return nil
}
