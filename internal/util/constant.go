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

// Cluster type
const (
	SERVERLESS = "SERVERLESS"
	DEDICATED  = "DEDICATED"
)

// Built-in SQL User Role types
const (
	ADMIN     = "admin"
	READWRITE = "readwrite"
	READONLY  = "readonly"
)

// Server accepted built-in SQL User Role types
const (
	ADMIN_ROLE     = "role_admin"
	READWRITE_ROLE = "role_readwrite"
	READONLY_ROLE  = "role_readonly"
)

const (
	ADMIN_DISPLAY     = "Database Admin"
	READWRITE_DISPLAY = "Database Read-Write"
	READONLY_DISPLAY  = "Database Read-Only"
)

const (
	MYSQLNATIVEPASSWORD = "mysql_native_password"
)
