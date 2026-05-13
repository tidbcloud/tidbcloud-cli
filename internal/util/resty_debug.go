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

package util

import (
	"github.com/go-resty/resty/v2"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/redact"
)

// ConfigureRestyDebug enables Resty's debug mode with redaction callbacks.
// Resty does not expose a callback for the request URI in its debug log, so do
// not use this helper for pre-signed URL transfers.
func ConfigureRestyDebug(client *resty.Client, debug bool) {
	client.SetDebug(debug)
	if !debug {
		return
	}

	client.OnRequestLog(func(log *resty.RequestLog) error {
		log.Header = redact.RedactHeader(log.Header)
		log.Body = redact.RedactBodyString(log.Body)
		return nil
	})
	client.OnResponseLog(func(log *resty.ResponseLog) error {
		log.Header = redact.RedactHeader(log.Header)
		log.Body = redact.RedactBodyString(log.Body)
		return nil
	})
}
