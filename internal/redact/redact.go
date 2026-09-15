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

// Package redact masks credentials in values that are printed for debugging
// (HTTP dumps, profile output) so that debug mode does not leak secrets.
package redact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const Mask = "***"

// normalize lower-cases a key and strips "-" and "_" so that
// "secret_access_key", "secretAccessKey" and "secret-access-key" all compare equal.
func normalize(key string) string {
	return strings.NewReplacer("-", "", "_", "").Replace(strings.ToLower(key))
}

var sensitiveKeys = map[string]struct{}{
	"secret":            {},
	"secretaccesskey":   {},
	"accesskeysecret":   {},
	"serviceaccountkey": {},
	"sastoken":          {},
	"azuresastoken":     {},
	"privatekey":        {},
	"oauthclientsecret": {},
	"clientsecret":      {},
	"accesstoken":       {},
	"devicecode":        {},
	"refreshtoken":      {},
	"token":             {},
	"password":          {},
	"rootpassword":      {},
	"sslkeycontent":     {},
}

var sensitiveHeaders = map[string]struct{}{
	"authorization":        {},
	"proxy-authorization":  {},
	"cookie":               {},
	"set-cookie":           {},
	"x-amz-security-token": {},
}

// IsSensitiveKey reports whether a JSON field, profile property or query
// parameter name holds a credential.
func IsSensitiveKey(key string) bool {
	_, ok := sensitiveKeys[normalize(key)]
	return ok
}

// isSensitiveQueryParam covers pre-signed URL parameters of S3, GCS, Azure and OSS.
func isSensitiveQueryParam(key string) bool {
	k := normalize(key)
	if IsSensitiveKey(k) || k == "sig" || k == "ossaccesskeyid" {
		return true
	}
	for _, suffix := range []string{"signature", "credential", "securitytoken"} {
		if strings.HasSuffix(k, suffix) {
			return true
		}
	}
	return false
}

// Value returns Mask when key is sensitive, otherwise value unchanged.
func Value(key, value string) string {
	if IsSensitiveKey(key) {
		return Mask
	}
	return value
}

// Any recursively masks sensitive keys in decoded JSON / viper values.
func Any(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, val := range t {
			if IsSensitiveKey(k) {
				out[k] = Mask
			} else {
				out[k] = Any(val)
			}
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(t))
		for i, val := range t {
			out[i] = Any(val)
		}
		return out
	default:
		return v
	}
}

// JSON masks sensitive fields of a JSON body. Non-JSON bodies are not echoed.
func JSON(body []byte) string {
	var v interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return fmt.Sprintf("<non-JSON body omitted, %d bytes>", len(body))
	}
	out, err := json.Marshal(Any(v))
	if err != nil {
		return fmt.Sprintf("<body omitted, %d bytes>", len(body))
	}
	return string(out)
}

// Headers returns a copy of h with credential-bearing headers masked.
func Headers(h http.Header) http.Header {
	out := h.Clone()
	for k := range out {
		if _, ok := sensitiveHeaders[strings.ToLower(k)]; ok {
			out[k] = []string{Mask}
		}
	}
	return out
}

// URL returns u as a string with credential-bearing query parameters masked.
func URL(u *url.URL) string {
	if u == nil {
		return ""
	}
	c := *u
	if c.RawQuery != "" {
		q := c.Query()
		for k := range q {
			if isSensitiveQueryParam(k) {
				q[k] = []string{Mask}
			}
		}
		c.RawQuery = q.Encode()
	}
	// Redacted replaces a userinfo password with "xxxxx".
	return c.Redacted()
}

// body echoes a JSON body (masked) and replaces *rc so it can still be read.
// Non-JSON bodies (file uploads, downloads) are never read.
func body(rc *io.ReadCloser, contentType string, contentLength int64) string {
	if rc == nil || *rc == nil || *rc == http.NoBody {
		return ""
	}
	if !strings.Contains(strings.ToLower(contentType), "json") {
		return fmt.Sprintf("<%s body omitted, %d bytes>", contentType, contentLength)
	}
	data, err := io.ReadAll(*rc)
	(*rc).Close()
	*rc = io.NopCloser(bytes.NewReader(data))
	if err != nil {
		return fmt.Sprintf("<body unreadable: %v>", err)
	}
	return JSON(data)
}

// DumpRequest renders r with credentials masked.
func DumpRequest(r *http.Request) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s %s\n", r.Method, URL(r.URL), r.Proto)
	host := r.Host
	if host == "" && r.URL != nil {
		host = r.URL.Host
	}
	fmt.Fprintf(&b, "Host: %s\n", host)
	_ = Headers(r.Header).Write(&b)
	b.WriteString("\n")
	b.WriteString(body(&r.Body, r.Header.Get("Content-Type"), r.ContentLength))
	b.WriteString("\n")
	return b.String()
}

// DumpResponse renders resp with credentials masked.
func DumpResponse(resp *http.Response) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s\n", resp.Proto, resp.Status)
	_ = Headers(resp.Header).Write(&b)
	b.WriteString("\n")
	b.WriteString(body(&resp.Body, resp.Header.Get("Content-Type"), resp.ContentLength))
	b.WriteString("\n")
	return b.String()
}

// DebugTransport prints masked request/response dumps to stdout when Debug is set.
type DebugTransport struct {
	Inner http.RoundTripper
	Debug bool
}

func NewDebugTransport(inner http.RoundTripper, debug bool) http.RoundTripper {
	if inner == nil {
		inner = http.DefaultTransport
	}
	return &DebugTransport{Inner: inner, Debug: debug}
}

func (dt *DebugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if dt.Debug {
		fmt.Printf("\n%s", DumpRequest(r))
	}
	resp, err := dt.Inner.RoundTrip(r)
	if err != nil {
		return resp, err
	}
	if dt.Debug {
		fmt.Printf("%s\n", DumpResponse(resp))
	}
	return resp, err
}
