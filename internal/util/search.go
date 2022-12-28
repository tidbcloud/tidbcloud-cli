// Copyright 2022 PingCAP, Inc.
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
	"strings"
)

const tiupBinPrefix = "/.tiup/components/"

func ElemInSlice[T comparable](slice []T, o T) bool {
	for _, b := range slice {
		if b == o {
			return true
		}
	}
	return false
}

// IsUnderTiUP checks whether the given binary is under the TiUP path.
func IsUnderTiUP(binpath string) bool {
	if binpath == "" {
		return false
	}

	if strings.Contains(binpath, tiupBinPrefix) {
		return true
	}

	return false
}
