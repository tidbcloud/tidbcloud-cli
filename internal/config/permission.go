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

package config

import (
	"io/fs"
	"os"
	"runtime"
)

// RestrictPermissions clears every permission bit of path that is not in
// limit, and keeps the other bits unchanged. For example, with limit 0600 a
// 0644 file becomes 0600 and a 0400 file stays 0400. It does nothing on
// Windows, where os.Chmod only controls the read-only attribute.
//
// The config directory and file can exist from an older version or be created
// by another tool with broader permissions; the create-time permissions set by
// os.MkdirAll and viper do not apply to them.
func RestrictPermissions(path string, limit fs.FileMode) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	perm := info.Mode().Perm()
	if perm&^limit == 0 {
		return nil
	}
	return os.Chmod(path, perm&limit)
}
