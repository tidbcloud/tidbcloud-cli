// Copyright 2025 PingCAP, Inc.
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

package util

import "github.com/tidbcloud/tidbcloud-cli/internal/flag"

var InputDescription = map[string]string{
	flag.OutputPath:                "Input the download path, press Enter to skip and download to the current directory",
	flag.StartDate:                 "Input the start date of the download in the format of 'YYYY-MM-DD'",
	flag.EndDate:                   "Input the end date of the download in the format of 'YYYY-MM-DD'",
	flag.AuditLogFilterRuleName:    "Input the filter rule name",
	flag.AuditLogFilterRuleFilters: "Input the filter rule expression, use `ticloud serverless audit-log filter template` to get the template",
}
