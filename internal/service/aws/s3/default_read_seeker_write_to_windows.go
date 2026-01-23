// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s3

import "github.com/aws/aws-sdk-go-v2/feature/s3/manager"

// BufferedReadSeekerWriteToPool is only used by default when executing in Windows environments.
// It exists to work around constraints with Go standard library's os file implementations for Windows.
func defaultUploadBufferProvider() manager.ReadSeekerWriteToProvider {
	return manager.NewBufferedReadSeekerWriteToPool(1024 * 1024)
}
