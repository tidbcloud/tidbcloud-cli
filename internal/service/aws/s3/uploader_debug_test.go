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

package s3

import (
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/config"
)

func TestNewUploaderDoesNotEnableRawRestyDebugForPresignedURLs(t *testing.T) {
	t.Setenv(config.DebugEnv, "1")

	uploader, ok := NewUploader(nil).(*UploaderImpl)
	if !ok {
		t.Fatalf("unexpected uploader type %T", uploader)
	}
	if uploader.httpClient.Debug {
		t.Fatal("raw Resty debug must stay disabled for pre-signed upload URLs")
	}
}
