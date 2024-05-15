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

package ui

import (
	"io"
	"strings"
	"tidbcloud-cli/internal/util"
)

type progressWriter struct {
	id             int
	downloadedSize int
	reader         io.Reader
	onResult       func(int, error, JobStatus)
	path           string
	fileName       string
}

func (pw *progressWriter) Read(p []byte) (n int, err error) {
	n, err = pw.reader.Read(p)
	if err == nil || err == io.EOF {
		pw.downloadedSize += n
	}
	return
}

func (pw *progressWriter) Start() {
	// create temp file
	tempFile, err := util.CreateTempFile(pw.path, pw.fileName)
	if err != nil {
		if strings.Contains(err.Error(), "file already exists") {
			pw.onResult(pw.id, err, Skipped)
		} else {
			pw.onResult(pw.id, err, Failed)
		}
		return
	}
	defer tempFile.Close()
	_, err = io.Copy(tempFile, pw)
	if err != nil {
		_ = util.DeleteFile(pw.path, tempFile.Name())
		pw.onResult(pw.id, err, Failed)
		return
	}

	err = util.RenameFile(pw.path, tempFile.Name(), pw.fileName)
	if err != nil {
		pw.onResult(pw.id, err, Failed)
		return
	}
	pw.onResult(pw.id, nil, Succeeded)
}
