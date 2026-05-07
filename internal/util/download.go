// Copyright 2026 PingCAP, Inc.
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
	"runtime"
	"strings"
	"unicode"

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
	filePath, err := safeDownloadPath(fileName)
	if err != nil {
		return nil, err
	}

	if path == "" {
		path = "."
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, err
	}
	defer root.Close()

	file, err := root.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("file already exists")
		}
		return nil, fmt.Errorf("download destination %q cannot be created: %w", fileName, err)
	}
	return file, nil
}

func safeDownloadPath(fileName string) (string, error) {
	if fileName == "" {
		return "", downloadDestinationError(fileName, "cannot be empty")
	}
	if fileName == "." || fileName == ".." {
		return "", downloadDestinationError(fileName, "must refer to a file")
	}
	if filepath.IsAbs(fileName) || isWindowsAbs(fileName) {
		return "", downloadDestinationError(fileName, "must be relative to the output path")
	}
	if containsControlCharacter(fileName) {
		return "", downloadDestinationError(fileName, "contains unsupported characters")
	}

	clean := filepath.Clean(fileName)
	if clean == "." || clean == ".." {
		return "", downloadDestinationError(fileName, "must refer to a file")
	}
	if strings.HasPrefix(clean, ".."+string(os.PathSeparator)) || filepath.IsAbs(clean) {
		return "", downloadDestinationError(fileName, "is outside the output path")
	}

	return clean, nil
}

func downloadDestinationError(fileName, reason string) error {
	return fmt.Errorf("download destination %q %s", fileName, reason)
}

func containsControlCharacter(s string) bool {
	return strings.ContainsRune(s, 0) || strings.IndexFunc(s, unicode.IsControl) >= 0
}

func isWindowsAbs(path string) bool {
	if runtime.GOOS == "windows" {
		return false
	}
	if len(path) >= 3 && path[1] == ':' && (path[2] == '\\' || path[2] == '/') {
		c := path[0]
		return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
	}
	return strings.HasPrefix(path, `\\`)
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
