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
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
)

// GetResponse returns the response of a given AWS per-signed URL
func GetResponse(url string, debug bool) (*http.Response, error) {
	httpClient := resty.New()
	httpClient.SetDebug(debug)
	resp, err := httpClient.GetClient().Get(url) // nolint:gosec
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		// read the body to get the error message
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("receiving status of %d", resp.StatusCode)
		}
		type AwsError struct {
			Code    string `xml:"Code"`
			Message string `xml:"Message"`
		}
		v := AwsError{}
		err = xml.Unmarshal(body, &v)
		if err != nil {
			return nil, fmt.Errorf("receiving status of %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("receiving status of %d. code: %s, message: %s", resp.StatusCode, v.Code, v.Message)
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

// CreateTempFile create a temp file if it does not exist
func CreateTempFile(path, fileName string) (*os.File, error) {
	// check the original file
	filePath := filepath.Join(path, fileName)
	if _, err := os.Stat(filePath); err == nil {
		return nil, fmt.Errorf("file already exists")
	}
	// try to delete the temp file if it exists
	tempFile := filepath.Join(path, fileName+".tmp")
	if _, err := os.Stat(filePath); err == nil {
		err = DeleteFile(path, fileName+".tmp")
		if err != nil {
			return nil, err
		}
	}
	// create the temp file
	file, err := os.Create(tempFile)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func RenameFile(path, oldFileName, newFileName string) error {
	oldFilePath := filepath.Join(path, oldFileName)
	newFilePath := filepath.Join(path, newFileName)

	return os.Rename(oldFilePath, newFilePath)
}

func DeleteFile(path, fileName string) error {
	filePath := filepath.Join(path, fileName)
	return os.RemoveAll(filePath)
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
