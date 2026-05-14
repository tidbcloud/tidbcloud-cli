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

// Package redact provides redacted HTTP request/response dumps for debug
// logging. It is a drop-in replacement for httputil.DumpRequestOut and
// httputil.DumpResponse with body=true, but strips credential-bearing
// headers and JSON body fields before serializing.
package redact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

// Placeholder is the string substituted in place of a redacted value.
const Placeholder = "[REDACTED]"

// sensitiveHeaders are HTTP headers whose values must never be printed.
var sensitiveHeaders = []string{
	"Authorization",
	"Proxy-Authorization",
	"Cookie",
	"Set-Cookie",
}

// sensitiveJSONKeys are JSON object key names (case-insensitive) whose
// values must be replaced with the placeholder.
var sensitiveJSONKeys = map[string]struct{}{
	"secret":              {},
	"secretaccesskey":     {},
	"accesskeysecret":     {},
	"serviceaccountkey":   {},
	"sastoken":            {},
	"sas_token":           {},
	"private-key":         {},
	"privatekey":          {},
	"private_key":         {},
	"oauth-client-secret": {},
	"oauthclientsecret":   {},
	"oauth_client_secret": {},
	"clientsecret":        {},
	"client_secret":       {},
	"access-token":        {},
	"access_token":        {},
	"accesstoken":         {},
	"refresh-token":       {},
	"refresh_token":       {},
	"refreshtoken":        {},
	"password":            {},
	"token":               {},
}

// DumpRequestOut behaves like httputil.DumpRequestOut(r, true) but redacts
// sensitive headers and JSON body fields. The original request's headers
// and body are restored before returning so the caller can still send it.
func DumpRequestOut(r *http.Request) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("redact: nil request")
	}

	origHeader := cloneHeader(r.Header)
	origBody, err := drainBody(&r.Body)
	if err != nil {
		return nil, err
	}

	// Apply redactions in place.
	r.Header = cloneHeader(origHeader)
	redactHeaders(r.Header)
	redactedBody := redactBody(origBody, r.Header.Get("Content-Type"))
	if origBody != nil {
		r.Body = io.NopCloser(bytes.NewReader(redactedBody))
		r.ContentLength = int64(len(redactedBody))
	}

	dump, dumpErr := httputil.DumpRequestOut(r, true)

	// Restore original headers and body so the request can still be sent.
	r.Header = origHeader
	if origBody != nil {
		r.Body = io.NopCloser(bytes.NewReader(origBody))
		r.ContentLength = int64(len(origBody))
	}

	return dump, dumpErr
}

// DumpResponse behaves like httputil.DumpResponse(resp, true) but redacts
// sensitive headers and JSON body fields. The response body is restored
// before returning so the caller can still read it.
func DumpResponse(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("redact: nil response")
	}

	origHeader := cloneHeader(resp.Header)
	origBody, err := drainBody(&resp.Body)
	if err != nil {
		return nil, err
	}

	resp.Header = cloneHeader(origHeader)
	redactHeaders(resp.Header)
	redactedBody := redactBody(origBody, resp.Header.Get("Content-Type"))
	if origBody != nil {
		resp.Body = io.NopCloser(bytes.NewReader(redactedBody))
		resp.ContentLength = int64(len(redactedBody))
	}

	dump, dumpErr := httputil.DumpResponse(resp, true)

	resp.Header = origHeader
	if origBody != nil {
		resp.Body = io.NopCloser(bytes.NewReader(origBody))
		resp.ContentLength = int64(len(origBody))
	}

	return dump, dumpErr
}

// redactHeaders mutates h, replacing values of sensitive headers with the
// placeholder. For Authorization-style headers the auth scheme prefix
// (Bearer / Digest / Basic) is preserved so the auth method stays visible
// in logs.
func redactHeaders(h http.Header) {
	for _, name := range sensitiveHeaders {
		values := h.Values(name)
		if len(values) == 0 {
			continue
		}
		h.Del(name)
		for _, v := range values {
			h.Add(name, redactHeaderValue(name, v))
		}
	}
}

func redactHeaderValue(name, value string) string {
	if strings.EqualFold(name, "Authorization") || strings.EqualFold(name, "Proxy-Authorization") {
		parts := strings.SplitN(value, " ", 2)
		if len(parts) == 2 {
			scheme := parts[0]
			if strings.EqualFold(scheme, "Bearer") ||
				strings.EqualFold(scheme, "Digest") ||
				strings.EqualFold(scheme, "Basic") {
				return scheme + " " + Placeholder
			}
		}
	}
	return Placeholder
}

// redactBody returns a redacted copy of body. JSON bodies have sensitive
// keys replaced; non-JSON bodies are replaced with a placeholder noting
// their length, since they may contain pre-signed URLs or form-encoded
// credentials.
func redactBody(body []byte, contentType string) []byte {
	if len(body) == 0 {
		return body
	}

	if isJSONContentType(contentType) || looksLikeJSON(body) {
		var parsed interface{}
		if err := json.Unmarshal(body, &parsed); err == nil {
			out, marshalErr := json.Marshal(redactJSONValue(parsed))
			if marshalErr == nil {
				return out
			}
		}
	}

	return []byte(fmt.Sprintf("[REDACTED non-JSON body, %d bytes]", len(body)))
}

func redactJSONValue(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, val := range t {
			if isSensitiveJSONKey(k) {
				out[k] = Placeholder
			} else {
				out[k] = redactJSONValue(val)
			}
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(t))
		for i, item := range t {
			out[i] = redactJSONValue(item)
		}
		return out
	default:
		return v
	}
}

func isSensitiveJSONKey(k string) bool {
	_, ok := sensitiveJSONKeys[strings.ToLower(k)]
	return ok
}

func isJSONContentType(ct string) bool {
	ct = strings.ToLower(ct)
	return strings.HasPrefix(ct, "application/json") || strings.Contains(ct, "+json")
}

func looksLikeJSON(body []byte) bool {
	for _, b := range body {
		switch b {
		case ' ', '\t', '\n', '\r':
			continue
		case '{', '[':
			return true
		default:
			return false
		}
	}
	return false
}

func cloneHeader(h http.Header) http.Header {
	if h == nil {
		return http.Header{}
	}
	out := make(http.Header, len(h))
	for k, v := range h {
		dup := make([]string, len(v))
		copy(dup, v)
		out[k] = dup
	}
	return out
}

// drainBody fully reads *bodyPtr (if non-nil), closes the original, and
// returns the bytes. *bodyPtr is left as nil; the caller is responsible
// for setting a fresh reader.
func drainBody(bodyPtr *io.ReadCloser) ([]byte, error) {
	if bodyPtr == nil || *bodyPtr == nil || *bodyPtr == http.NoBody {
		return nil, nil
	}
	b, err := io.ReadAll(*bodyPtr)
	if err != nil {
		return nil, err
	}
	_ = (*bodyPtr).Close()
	*bodyPtr = nil
	return b, nil
}
