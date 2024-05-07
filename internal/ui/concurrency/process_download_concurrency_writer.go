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

package ui_concurrency

import (
	"io"
	"os"
	"time"
)

type progressConcurrencyWriter struct {
	id         int
	total      int
	downloaded int
	file       *os.File
	reader     io.Reader
	onProgress func(int, float64)
	onError    func(int, error)
	percent    float64
}

func (pw *progressConcurrencyWriter) Read(p []byte) (n int, err error) {
	n, err = pw.reader.Read(p)
	pw.downloaded += len(p)
	return
}

func (pw *progressConcurrencyWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	go pw.Watch()
	_, err := io.Copy(pw.file, pw)
	if err != nil {
		pw.onError(pw.id, err)
	}
}

func (pw *progressConcurrencyWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	if pw.total > 0 && pw.onProgress != nil {
		percentNow := float64(pw.downloaded) / float64(pw.total)
		if percentNow > 0.9 || percentNow-pw.percent > 0.05 {
			pw.percent = percentNow
			pw.onProgress(pw.id, pw.percent)
		}
	}
	return len(p), nil
}

func (pw *progressConcurrencyWriter) Watch() {
	for {
		time.Sleep(500 * time.Millisecond)
		percentNow := float64(pw.downloaded) / float64(pw.total)
		if percentNow == pw.percent {
			continue
		}
		if percentNow > 0.9 || percentNow-pw.percent > 0.05 {
			pw.percent = percentNow
			pw.onProgress(pw.id, pw.percent)
		}
		if percentNow >= 1 {
			break
		}
	}
}
