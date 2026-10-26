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

package cloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal/redact"
)

// NextGenAPIError exposes sanitized metadata and text to all output consumers.
type NextGenAPIError struct {
	Status    int                    `json:"status,omitempty"`
	RequestID string                 `json:"requestId,omitempty"`
	Method    string                 `json:"method,omitempty"`
	Path      string                 `json:"path,omitempty"`
	Response  map[string]interface{} `json:"response,omitempty"`
	text      string
	cause     error
}

func (e *NextGenAPIError) Error() string { return e.text }

func (e *NextGenAPIError) Unwrap() error { return e.cause }

func newNextGenAPIError(response *http.Response, detail string, cause error) error {
	e := &NextGenAPIError{
		Status: response.StatusCode, RequestID: response.Header.Get("X-Debug-Trace-Id"), cause: cause,
	}
	if response.Request != nil && response.Request.URL != nil {
		e.Method, e.Path = response.Request.Method, response.Request.URL.Path
	}
	decoder := json.NewDecoder(strings.NewReader(detail))
	decoder.UseNumber()
	var body map[string]interface{}
	// A failed 2xx response is not a valid API error payload. Omit its body.
	if response.StatusCode >= 300 && json.Valid([]byte(detail)) && decoder.Decode(&body) == nil && body != nil {
		e.Response = redact.Any(body).(map[string]interface{})
	}
	safeDetail := fmt.Sprintf("<response body omitted, %d bytes>", len(detail))
	if e.Response != nil {
		if data, err := json.Marshal(e.Response); err == nil {
			safeDetail = string(data)
		}
	}
	e.text = fmt.Sprintf("[%s %s][%d %s][%s] %s", e.Method, e.Path, e.Status, http.StatusText(e.Status), e.RequestID, safeDetail)
	if reason := e.LocalFailureReason(); reason != "" {
		e.text += "; " + reason
	}
	return e
}

// LocalFailureReason returns safe local failure descriptions, never raw SDK text
// or URL-bearing transport errors.
func (e *NextGenAPIError) LocalFailureReason() string {
	for _, known := range []error{context.DeadlineExceeded, context.Canceled, errNextGenRedirect} {
		if errors.Is(e.cause, known) {
			return known.Error()
		}
	}
	if e.cause != nil && e.Status >= 200 && e.Status < 300 {
		return "failed to read or decode the API response"
	}
	return ""
}
