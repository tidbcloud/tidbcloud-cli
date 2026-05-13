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

package redact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

const Mask = "******"

var (
	sensitiveKeys = map[string]struct{}{
		"authorization":      {},
		"proxyauthorization": {},
		"cookie":             {},
		"setcookie":          {},

		"accesskey":         {},
		"accesskeyid":       {},
		"accesskeysecret":   {},
		"accesstoken":       {},
		"awsaccesskeyid":    {},
		"clientsecret":      {},
		"credential":        {},
		"devicecode":        {},
		"googleaccessid":    {},
		"oauthclientsecret": {},
		"password":          {},
		"privatekey":        {},
		"refreshtoken":      {},
		"sastoken":          {},
		"secret":            {},
		"secretaccesskey":   {},
		"securitytoken":     {},
		"serviceaccountkey": {},
		"sig":               {},
		"signature":         {},
		"token":             {},

		"xamzcredential":     {},
		"xamzsecuritytoken":  {},
		"xamzsignature":      {},
		"xgoogcredential":    {},
		"xgoogsecuritytoken": {},
		"xgoogsignature":     {},
		"xosssecuritytoken":  {},
	}

	assignmentPattern = regexp.MustCompile(`(?i)(access[-_]?token|refresh[-_]?token|client[-_]?secret|oauth[-_]?client[-_]?secret|private[-_]?key|secret[-_]?access[-_]?key|access[-_]?key[-_]?secret|service[-_]?account[-_]?key|sas[-_]?token|password|token|secret)\s*([:=])\s*("[^"]*"|'[^']*'|[^\s,&}]+)`)
	bearerPattern     = regexp.MustCompile(`(?i)(Bearer\s+)[A-Za-z0-9._~+/=-]+`)
)

// IsSensitiveKey reports whether a header, query parameter, or body field name
// commonly carries credentials or tokens.
func IsSensitiveKey(key string) bool {
	_, ok := sensitiveKeys[normalizeKey(key)]
	return ok
}

// MaskValue returns Mask for sensitive key names and the original value otherwise.
func MaskValue(key, value string) string {
	if IsSensitiveKey(key) {
		return Mask
	}
	return value
}

// MaskAny returns a redacted deep copy of maps and slices. Scalar values are
// masked only when their parent key is sensitive.
func MaskAny(value interface{}) interface{} {
	return maskAny(reflect.ValueOf(value), "")
}

// RedactHeader returns a redacted copy of h.
func RedactHeader(h http.Header) http.Header {
	if h == nil {
		return nil
	}
	redacted := h.Clone()
	for key := range redacted {
		if IsSensitiveKey(key) {
			redacted[key] = []string{Mask}
		}
	}
	return redacted
}

// RedactURL redacts credential-bearing query parameters in a URL or request URI.
func RedactURL(raw string) string {
	if raw == "" {
		return raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	q := u.Query()
	changed := false
	for key, values := range q {
		if IsSensitiveKey(key) {
			q[key] = maskedValues(values)
			changed = true
		}
	}
	if !changed {
		return raw
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// RedactBodyString redacts sensitive fields in JSON or form-like bodies.
func RedactBodyString(body string) string {
	if strings.TrimSpace(body) == "" {
		return body
	}

	var v interface{}
	decoder := json.NewDecoder(strings.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err == nil {
		redacted, err := json.Marshal(MaskAny(v))
		if err == nil {
			return string(redacted)
		}
	}

	if values, err := url.ParseQuery(body); err == nil && len(values) > 0 {
		changed := false
		for key, value := range values {
			if IsSensitiveKey(key) {
				values[key] = maskedValues(value)
				changed = true
			}
		}
		if changed {
			return values.Encode()
		}
	}

	return RedactText(body)
}

// RedactBody redacts body bytes. JSON content is detected even if contentType is
// missing because HTTP dumps do not always retain enough context.
func RedactBody(body []byte, contentType string) []byte {
	if len(bytes.TrimSpace(body)) == 0 {
		return body
	}
	if isJSONContent(contentType) || json.Valid(body) || isFormContent(contentType) {
		return []byte(RedactBodyString(string(body)))
	}
	return []byte(RedactText(string(body)))
}

// RedactText applies conservative fallback redaction for non-JSON text.
func RedactText(text string) string {
	text = bearerPattern.ReplaceAllString(text, "${1}"+Mask)
	return assignmentPattern.ReplaceAllString(text, "$1$2"+Mask)
}

// DumpRequestOut is httputil.DumpRequestOut with sensitive headers, query
// parameters, and body fields masked before returning bytes to callers.
func DumpRequestOut(req *http.Request, body bool) ([]byte, error) {
	dump, err := httputil.DumpRequestOut(req, body)
	if err != nil {
		return nil, err
	}
	return RedactHTTPDump(dump), nil
}

// DumpResponse is httputil.DumpResponse with sensitive headers and body fields
// masked before returning bytes to callers.
func DumpResponse(resp *http.Response, body bool) ([]byte, error) {
	dump, err := httputil.DumpResponse(resp, body)
	if err != nil {
		return nil, err
	}
	return RedactHTTPDump(dump), nil
}

// RedactHTTPDump redacts a raw HTTP request or response dump.
func RedactHTTPDump(dump []byte) []byte {
	head, body, sep := splitHTTPDump(dump)
	redactedHead := redactHTTPHead(string(head))
	if sep == "" {
		return []byte(redactedHead)
	}
	return []byte(redactedHead + sep + string(RedactBody(body, headerContentType(redactedHead))))
}

func maskAny(value reflect.Value, parentKey string) interface{} {
	if !value.IsValid() {
		return nil
	}
	if parentKey != "" && IsSensitiveKey(parentKey) {
		return Mask
	}
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Map:
		out := make(map[string]interface{}, value.Len())
		for _, key := range value.MapKeys() {
			keyString := fmt.Sprint(key.Interface())
			out[keyString] = maskAny(value.MapIndex(key), keyString)
		}
		return out
	case reflect.Slice, reflect.Array:
		out := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			out[i] = maskAny(value.Index(i), parentKey)
		}
		return out
	default:
		return value.Interface()
	}
}

func normalizeKey(key string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, key)
}

func maskedValues(values []string) []string {
	if len(values) == 0 {
		return []string{Mask}
	}
	out := make([]string, len(values))
	for i := range out {
		out[i] = Mask
	}
	return out
}

func isJSONContent(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mediaType = contentType
	}
	mediaType = strings.ToLower(mediaType)
	return mediaType == "application/json" || strings.HasSuffix(mediaType, "+json")
}

func isFormContent(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mediaType = contentType
	}
	return strings.EqualFold(mediaType, "application/x-www-form-urlencoded")
}

func splitHTTPDump(dump []byte) ([]byte, []byte, string) {
	if idx := bytes.Index(dump, []byte("\r\n\r\n")); idx >= 0 {
		return dump[:idx], dump[idx+4:], "\r\n\r\n"
	}
	if idx := bytes.Index(dump, []byte("\n\n")); idx >= 0 {
		return dump[:idx], dump[idx+2:], "\n\n"
	}
	return dump, nil, ""
}

func redactHTTPHead(head string) string {
	lines := strings.Split(head, "\n")
	for i, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if i == 0 {
			lines[i] = redactStartLine(line)
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			lines[i] = line
			continue
		}
		name := line[:idx]
		value := strings.TrimSpace(line[idx+1:])
		if IsSensitiveKey(name) {
			lines[i] = name + ": " + Mask
			continue
		}
		lines[i] = name + ": " + RedactURL(value)
	}
	return strings.Join(lines, "\n")
}

func redactStartLine(line string) string {
	parts := strings.Split(line, " ")
	if len(parts) == 3 && strings.HasPrefix(parts[2], "HTTP/") {
		parts[1] = RedactURL(parts[1])
		return strings.Join(parts, " ")
	}
	return RedactURL(line)
}

func headerContentType(head string) string {
	for _, line := range strings.Split(head, "\n") {
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(line[:idx]), "Content-Type") {
			return strings.TrimSpace(line[idx+1:])
		}
	}
	return ""
}
