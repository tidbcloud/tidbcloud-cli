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

//go:build darwin
// +build darwin

package start

import (
	"bytes"

	"github.com/pingcap/errors"
	exec "golang.org/x/sys/execabs"
)

func (m *MySQLHelperImpl) DumpFromMySQL(arg string) error {
	c1 := exec.Command("sh", "-c", arg) //nolint:gosec
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
