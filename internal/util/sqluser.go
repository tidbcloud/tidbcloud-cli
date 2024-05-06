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
	"strings"
)

func GetDisplayRole(builtinRole string, customRoles []string) string {
	displayRole := ""
	if builtinRole != "" {
		builtinRole = trimRolePrefix(builtinRole)
		switch builtinRole {
		case ADMIN_ROLE:
			displayRole = ADMIN_DISPLAY
		case READWRITE_ROLE:
			displayRole = READWRITE_DISPLAY
		case READONLY_ROLE:
			displayRole = READONLY_DISPLAY
		}
	}

	// put built-in role in the first place
	allRoles := make([]string, 0, len(customRoles)+1)
	if displayRole != "" {
		allRoles = append(allRoles, displayRole)
	}
	allRoles = append(allRoles, customRoles...)

	joinedRoles := strings.Join(allRoles, ", ")
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
	return role == ADMIN_ROLE || role == READWRITE_ROLE || role == READONLY_ROLE ||
		role == ADMIN_DISPLAY || role == READWRITE_DISPLAY || role == READONLY_DISPLAY
}
