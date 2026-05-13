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
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDumpRequestOutRedactsSensitiveData(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "https://example.com/upload?X-Amz-Signature=raw-signature&token=raw-token&safe=visible", strings.NewReader(`{"accessKey":{"id":"raw-access-key-id","secret":"raw-secret"},"nested":{"serviceAccountKey":"raw-gcs-key"},"name":"visible"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer raw-bearer")
	req.Header.Set("Cookie", "session=raw-cookie")
	req.Header.Set("Content-Type", "application/json")

	dump, err := DumpRequestOut(req, true)
	if err != nil {
		t.Fatal(err)
	}
	got := string(dump)

	assertNotContains(t, got, "raw-signature", "raw-token", "raw-access-key-id", "raw-secret", "raw-gcs-key", "raw-bearer", "raw-cookie")
	assertContains(t, got, "safe=visible", `"name":"visible"`, Mask)
}

func TestDumpResponseRedactsSensitiveData(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Set-Cookie":    []string{"session=raw-cookie"},
			"Content-Type":  []string{"application/json"},
			"X-Request-Id":  []string{"visible-request-id"},
			"Authorization": []string{"Bearer raw-bearer"},
		},
		Body: io.NopCloser(strings.NewReader(`{"access_token":"raw-access-token","display":"visible"}`)),
	}

	dump, err := DumpResponse(resp, true)
	if err != nil {
		t.Fatal(err)
	}
	got := string(dump)

	assertNotContains(t, got, "raw-cookie", "raw-bearer", "raw-access-token")
	assertContains(t, got, "visible-request-id", `"display":"visible"`, Mask)
}

func TestMaskAnyPreservesNonSensitiveValues(t *testing.T) {
	input := map[string]interface{}{
		"public-key":          "public",
		"private-key":         "private",
		"oauth-client-secret": "client-secret",
		"nested": map[string]string{
			"sasToken": "sas-token",
			"region":   "us-west-2",
		},
	}

	got := MaskAny(input).(map[string]interface{})
	nested := got["nested"].(map[string]interface{})

	if got["public-key"] != "public" {
		t.Fatalf("public key should not be masked: %#v", got["public-key"])
	}
	if got["private-key"] != Mask || got["oauth-client-secret"] != Mask || nested["sasToken"] != Mask {
		t.Fatalf("sensitive fields were not masked: %#v", got)
	}
	if nested["region"] != "us-west-2" {
		t.Fatalf("non-sensitive nested field changed: %#v", nested["region"])
	}
}

func TestRedactURLMasksAzureSASSignature(t *testing.T) {
	got := RedactURL("https://account.blob.core.windows.net/container/file?sp=r&sig=raw-sas-signature&name=visible")

	assertNotContains(t, got, "raw-sas-signature")
	assertContains(t, got, "name=visible", "sig=%2A%2A%2A%2A%2A%2A")
}

func assertContains(t *testing.T, got string, needles ...string) {
	t.Helper()
	for _, needle := range needles {
		if !strings.Contains(got, needle) {
			t.Fatalf("expected %q to contain %q", got, needle)
		}
	}
}

func assertNotContains(t *testing.T, got string, needles ...string) {
	t.Helper()
	for _, needle := range needles {
		if strings.Contains(got, needle) {
			t.Fatalf("expected %q not to contain %q", got, needle)
		}
	}
}
