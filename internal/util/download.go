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
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"path/filepath"
)

// GetResponse returns the response of a given URL
func GetResponse(url string, debug bool) (*http.Response, error) {
	httpClient := resty.New()
	httpClient.SetDebug(debug)
	resp, err := httpClient.GetClient().Get(url) // nolint:gosec
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
	filePath := filepath.Join(path, fileName)
	if _, err := os.Stat(filePath); err == nil {
		return nil, fmt.Errorf("file already exists")
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// CreateFolder creates a folder if it does not exist
func CreateFolder(path string) error {
	if path == "" {
		return nil
	}
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
