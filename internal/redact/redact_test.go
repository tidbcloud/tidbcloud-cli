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
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDumpMasksCredentials(t *testing.T) {
	assert := require.New(t)
	const (
		bearer = "BEARER_SECRET_TOKEN"
		aksk   = "AKSK_SECRET"
		sig    = "PRESIGNED_SIG"
		cookie = "SESSION_COOKIE"
	)

	var received []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.Header().Set("Set-Cookie", "s="+cookie)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"` + bearer + `","refresh_token":"` + bearer + `","device_code":"` + bearer + `","user_code":"ABCD-EFGH","name":"me"}`))
	}))
	defer srv.Close()

	var out bytes.Buffer
	body := `{"target":{"s3":{"uri":"s3://b","accessKey":{"id":"AKIA","secret":"` + aksk + `"}}},"password":"` + aksk + `","sasToken":"` + aksk + `","device_code":"` + aksk + `"}`
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/v1?X-Amz-Signature="+sig+"&X-Amz-Credential="+sig+"&sig="+sig+"&plain=1", strings.NewReader(body))
	assert.NoError(err)
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("Content-Type", "application/json")

	out.WriteString(DumpRequest(req))
	resp, err := http.DefaultTransport.RoundTrip(req)
	assert.NoError(err)
	defer resp.Body.Close()
	out.WriteString(DumpResponse(resp))

	dump := out.String()
	for _, secret := range []string{bearer, aksk, sig, cookie} {
		assert.NotContains(dump, secret)
	}
	assert.Contains(dump, `"id":"AKIA"`)
	assert.Contains(dump, "plain=1")
	assert.Contains(dump, `"name":"me"`)
	assert.Contains(dump, `"user_code":"ABCD-EFGH"`)

	// bodies must still be readable by the caller after dumping
	assert.Equal(body, string(received))
	got, _ := io.ReadAll(resp.Body)
	assert.Contains(string(got), bearer)
}

func TestNonJSONBodyIsNotEchoed(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, "https://s3/x?X-Amz-Signature=abc", strings.NewReader("FILECONTENT"))
	req.Header.Set("Content-Type", "application/octet-stream")
	dump := DumpRequest(req)
	require.NotContains(t, dump, "FILECONTENT")
	require.NotContains(t, dump, "abc")
	sent, _ := io.ReadAll(req.Body)
	require.Equal(t, "FILECONTENT", string(sent))
}

func TestValueAndAny(t *testing.T) {
	assert := require.New(t)
	assert.Equal(Mask, Value("private-key", "x"))
	assert.Equal(Mask, Value("oauth-client-secret", "x"))
	assert.Equal("x", Value("public-key", "x"))
	assert.Equal("x", Value("serverless-endpoint", "x"))
	got := Any(map[string]interface{}{"public-key": "pk", "private-key": "sk", "access-token": "tk", "token-type": "Bearer"})
	assert.Equal(map[string]interface{}{"public-key": "pk", "private-key": Mask, "access-token": Mask, "token-type": "Bearer"}, got)
	u, _ := url.Parse("https://h/p?a=1")
	assert.Equal("https://h/p?a=1", URL(u))
	u, _ = url.Parse("https://alice:USERINFO_PASS@h/p")
	assert.Equal("https://alice:xxxxx@h/p", URL(u))
	u, _ = url.Parse("https://alice@h/p?a=1")
	assert.Equal("https://alice@h/p?a=1", URL(u))
}
