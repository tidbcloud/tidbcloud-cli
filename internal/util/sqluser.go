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
	"slices"
	"strings"
)

func GetDisplayRole(builtinRole string, customRoles []string) string {
	if builtinRole != "" {
		builtinRole = trimRolePrefix(builtinRole)
	}

	// put built-in role in the first place
	customRoles = slices.Insert(customRoles, 0, builtinRole)
	joinedRoles := strings.Join(customRoles, ",")
	return joinedRoles
}

// trimRolePrefix trims the role prefix user, only read_write, read_only started with the role prefix.
func trimRolePrefix(role string) string {
	roleTrimmed := role
	dotIndex := strings.Index(role, ".")
	if dotIndex != -1 {
		trimmedStr := role[dotIndex+1:]
		roleTrimmed = trimmedStr
	}

	return roleTrimmed
}

func IsBuiltinRole(role string) bool {
	return role == ADMIN_ROLE || role == READWRITE_ROLE || role == READONLY_ROLE
}

func TrimUserNamePrefix(userName string, prefix string) string {
	prefix = prefix + "."
	if strings.HasPrefix(userName, prefix) {
		return userName[len(prefix):]
	}
	return userName
}

func AddPrefix(s string, prefix string) string {
	prefix = prefix + "."
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return fmt.Sprintf("%s%s", prefix, s)
}
