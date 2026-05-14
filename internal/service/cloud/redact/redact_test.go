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
	"io"
	"net/http"
	"strings"
	"testing"
)

const (
	bearerSecret      = "eyJhbGciOiJIUzI1NiJ9.tokenpayload.signature123"
	digestRaw         = `Digest username="pub", realm="r", nonce="n", uri="/x", response="abcdef0123456789", qop=auth, nc=00000001, cnonce="c"`
	awsSecretAccess   = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	ossSecretAccess   = "ossSecretValue12345"
	awsAccessKeyID    = "AKIAIOSFODNN7EXAMPLE"
	gcsServiceAccount = `{"type":"service_account","private_key":"-----BEGIN PRIVATE KEY-----\nMIIEvA...\n-----END PRIVATE KEY-----"}`
	azureSAS          = "sv=2024-01-01&sig=abc123XYZ&se=2030-01-01"
)

func newRequest(t *testing.T, method, url, contentType string, body []byte) *http.Request {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	r, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	return r
}

func TestDumpRequestOut_BearerToken(t *testing.T) {
	r := newRequest(t, http.MethodGet, "https://example.com/v1/clusters", "", nil)
	r.Header.Set("Authorization", "Bearer "+bearerSecret)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	s := string(out)
	if strings.Contains(s, bearerSecret) {
		t.Fatalf("raw bearer token leaked into dump:\n%s", s)
	}
	if !strings.Contains(s, "Bearer "+Placeholder) {
		t.Fatalf("expected 'Bearer %s' in dump, got:\n%s", Placeholder, s)
	}
	// Original header must be restored on the request so it can still be sent.
	if r.Header.Get("Authorization") != "Bearer "+bearerSecret {
		t.Fatalf("original Authorization header was not restored: %q", r.Header.Get("Authorization"))
	}
}

func TestDumpRequestOut_DigestAuth(t *testing.T) {
	r := newRequest(t, http.MethodGet, "https://example.com/v1/x", "", nil)
	r.Header.Set("Authorization", digestRaw)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "response=\"abcdef") {
		t.Fatalf("digest response leaked:\n%s", s)
	}
	if !strings.Contains(s, "Digest "+Placeholder) {
		t.Fatalf("expected 'Digest %s' in dump, got:\n%s", Placeholder, s)
	}
}

func TestDumpRequestOut_CookieAndProxyAuth(t *testing.T) {
	r := newRequest(t, http.MethodGet, "https://example.com/", "", nil)
	r.Header.Set("Cookie", "session=supersecretcookievalue")
	r.Header.Set("Proxy-Authorization", "Basic dXNlcjpwYXNz")

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "supersecretcookievalue") {
		t.Fatalf("cookie value leaked:\n%s", s)
	}
	if strings.Contains(s, "dXNlcjpwYXNz") {
		t.Fatalf("proxy basic creds leaked:\n%s", s)
	}
	if !strings.Contains(s, "Basic "+Placeholder) {
		t.Fatalf("expected 'Basic %s', got:\n%s", Placeholder, s)
	}
}

func TestDumpRequestOut_S3ImportBody(t *testing.T) {
	body := []byte(`{"source":{"type":"S3","s3":{"uri":"s3://bucket/k","accessKey":{"id":"` + awsAccessKeyID + `","secret":"` + awsSecretAccess + `"}}}}`)
	r := newRequest(t, http.MethodPost, "https://example.com/v1/imports", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	s := string(out)
	if strings.Contains(s, awsSecretAccess) {
		t.Fatalf("S3 secret leaked:\n%s", s)
	}
	// Identifier (access key id) is not a secret per the report; should remain visible.
	if !strings.Contains(s, awsAccessKeyID) {
		t.Fatalf("expected access key id %q in dump, got:\n%s", awsAccessKeyID, s)
	}
}

func TestDumpRequestOut_OSSImportBody(t *testing.T) {
	body := []byte(`{"source":{"oss":{"accessKey":{"id":"LTAI...","secret":"` + ossSecretAccess + `"}}}}`)
	r := newRequest(t, http.MethodPost, "https://example.com/v1/imports", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	if strings.Contains(string(out), ossSecretAccess) {
		t.Fatalf("OSS secret leaked:\n%s", string(out))
	}
}

func TestDumpRequestOut_ExportS3Body(t *testing.T) {
	body := []byte(`{"target":{"s3":{"uri":"s3://b/p","accessKey":{"id":"AKIA","secret":"` + awsSecretAccess + `"}}}}`)
	r := newRequest(t, http.MethodPost, "https://example.com/v1/exports", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	if strings.Contains(string(out), awsSecretAccess) {
		t.Fatalf("export S3 secret leaked:\n%s", string(out))
	}
}

func TestDumpRequestOut_GCSServiceAccountKey(t *testing.T) {
	saKey, _ := json.Marshal(gcsServiceAccount)
	body := []byte(`{"target":{"gcs":{"uri":"gs://b/p","serviceAccountKey":` + string(saKey) + `}}}`)
	r := newRequest(t, http.MethodPost, "https://example.com/v1/exports", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	if strings.Contains(string(out), "BEGIN PRIVATE KEY") {
		t.Fatalf("GCS private key leaked:\n%s", string(out))
	}
	if strings.Contains(string(out), "service_account") {
		t.Fatalf("GCS service account JSON leaked:\n%s", string(out))
	}
}

func TestDumpRequestOut_AzureSasToken(t *testing.T) {
	body := []byte(`{"target":{"azureBlob":{"uri":"https://a.blob.core.windows.net/c/p","sasToken":"` + azureSAS + `"}}}`)
	r := newRequest(t, http.MethodPost, "https://example.com/v1/exports", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	if strings.Contains(string(out), azureSAS) {
		t.Fatalf("Azure SAS token leaked:\n%s", string(out))
	}
}

func TestDumpRequestOut_AuditLogStorageBody(t *testing.T) {
	body := []byte(`{"auditLogConfig":{"cloudStorage":{"s3":{"accessKey":{"id":"AKIA","secret":"` + awsSecretAccess + `"}}}}}`)
	r := newRequest(t, http.MethodPatch, "https://example.com/v1/clusters/x/auditLogConfig", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	if strings.Contains(string(out), awsSecretAccess) {
		t.Fatalf("audit log S3 secret leaked:\n%s", string(out))
	}
}

func TestDumpRequestOut_NestedJSONSecret(t *testing.T) {
	body := []byte(`{"a":{"b":{"c":{"secret":"deeplyNestedSecret"}}}}`)
	r := newRequest(t, http.MethodPost, "https://example.com/", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	if strings.Contains(string(out), "deeplyNestedSecret") {
		t.Fatalf("nested secret leaked:\n%s", string(out))
	}
}

func TestDumpRequestOut_JSONArrayWithSecrets(t *testing.T) {
	body := []byte(`{"users":[{"name":"u1","password":"p1secret"},{"name":"u2","password":"p2secret"}]}`)
	r := newRequest(t, http.MethodPost, "https://example.com/", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "p1secret") || strings.Contains(s, "p2secret") {
		t.Fatalf("password leaked from array:\n%s", s)
	}
	if !strings.Contains(s, `"u1"`) || !strings.Contains(s, `"u2"`) {
		t.Fatalf("non-sensitive fields lost:\n%s", s)
	}
}

func TestDumpRequestOut_NonJSONBody(t *testing.T) {
	body := []byte("raw=upload&signature=abc123XYZdangerous")
	r := newRequest(t, http.MethodPost, "https://example.com/", "application/x-www-form-urlencoded", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "abc123XYZdangerous") {
		t.Fatalf("non-JSON body content leaked:\n%s", s)
	}
	if !strings.Contains(s, "REDACTED non-JSON body") {
		t.Fatalf("expected non-JSON placeholder, got:\n%s", s)
	}
}

func TestDumpRequestOut_EmptyBody(t *testing.T) {
	r := newRequest(t, http.MethodGet, "https://example.com/v1/clusters", "", nil)
	r.Header.Set("Authorization", "Bearer "+bearerSecret)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	if strings.Contains(string(out), bearerSecret) {
		t.Fatalf("token leaked on empty-body request:\n%s", string(out))
	}
}

func TestDumpRequestOut_BodyStillReadable(t *testing.T) {
	original := []byte(`{"name":"x","secret":"y"}`)
	r := newRequest(t, http.MethodPost, "https://example.com/", "application/json", original)

	if _, err := DumpRequestOut(r); err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	got, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read restored body: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Fatalf("restored body mismatch: got %q want %q", got, original)
	}
}

func TestDumpRequestOut_PreservesNonSecretFields(t *testing.T) {
	body := []byte(`{"clusterId":"cluster-123","displayName":"my-cluster","region":"us-east-1","secret":"hidden"}`)
	r := newRequest(t, http.MethodPost, "https://example.com/", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	s := string(out)
	for _, want := range []string{"cluster-123", "my-cluster", "us-east-1"} {
		if !strings.Contains(s, want) {
			t.Fatalf("expected %q in dump, got:\n%s", want, s)
		}
	}
	if strings.Contains(s, "hidden") {
		t.Fatalf("secret leaked alongside non-secret fields:\n%s", s)
	}
}

func TestDumpRequestOut_CaseInsensitiveJSONKey(t *testing.T) {
	body := []byte(`{"Secret":"upperCaseSecret","SecretAccessKey":"mixedCaseSecret","ACCESS_TOKEN":"shoutingSecret"}`)
	r := newRequest(t, http.MethodPost, "https://example.com/", "application/json", body)

	out, err := DumpRequestOut(r)
	if err != nil {
		t.Fatalf("DumpRequestOut: %v", err)
	}
	s := string(out)
	for _, leak := range []string{"upperCaseSecret", "mixedCaseSecret", "shoutingSecret"} {
		if strings.Contains(s, leak) {
			t.Fatalf("%s leaked despite case-insensitive match:\n%s", leak, s)
		}
	}
}

func TestDumpResponse_SetCookieRedacted(t *testing.T) {
	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"ok":true}`))),
	}
	resp.Header.Set("Set-Cookie", "sid=responseCookieValue; HttpOnly")
	resp.Header.Set("Content-Type", "application/json")

	out, err := DumpResponse(resp)
	if err != nil {
		t.Fatalf("DumpResponse: %v", err)
	}
	if strings.Contains(string(out), "responseCookieValue") {
		t.Fatalf("Set-Cookie value leaked:\n%s", string(out))
	}
}

func TestDumpResponse_JSONBodyRedacted(t *testing.T) {
	body := []byte(`{"access_token":"shouldNotLeak","refresh_token":"alsoSecret","user":"alice"}`)
	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(body)),
	}

	out, err := DumpResponse(resp)
	if err != nil {
		t.Fatalf("DumpResponse: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "shouldNotLeak") || strings.Contains(s, "alsoSecret") {
		t.Fatalf("token leaked in response body:\n%s", s)
	}
	if !strings.Contains(s, "alice") {
		t.Fatalf("non-sensitive 'user' lost:\n%s", s)
	}
}

func TestDumpResponse_BodyStillReadable(t *testing.T) {
	original := []byte(`{"access_token":"x","user":"alice"}`)
	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(original)),
	}

	if _, err := DumpResponse(resp); err != nil {
		t.Fatalf("DumpResponse: %v", err)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read restored body: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Fatalf("restored body mismatch: got %q want %q", got, original)
	}
}
