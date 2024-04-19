// Copyright 2024 PingCAP, Inc.
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

package util

import (
	"fmt"
	"net/http"
	"os"
)

// GetResponse returns the response of a given URL
func GetResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url) // nolint:gosec
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("receiving status of %d for url: %s", resp.StatusCode, url)
	}
	if resp.ContentLength <= 0 {
		resp.Body.Close()
		return nil, fmt.Errorf("file is empty")
	}
	return resp, nil
}

// CreateFile creates a file if it does not exist
func CreateFile(path, fileName string) (*os.File, error) {
	if path == "" {
		path = "."
	}
	if _, err := os.Stat(path + "/" + fileName); err == nil {
		return nil, fmt.Errorf("file already exists")
	}
	file, err := os.Create(path + "/" + fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// CreateFolder creates a folder if it does not exist
func CreateFolder(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0755)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
