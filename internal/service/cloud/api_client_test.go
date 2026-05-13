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
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/config"
)

func TestDebugTransportRedactsRequestAndResponse(t *testing.T) {
	t.Setenv(config.DebugEnv, "1")

	var stdout bytes.Buffer
	restore := captureStdout(t, &stdout)

	transport := NewTransportWithBearToken(NewDebugTransport(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(string(body), "raw-secret") {
			t.Fatalf("inner transport did not receive original body: %s", body)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Set-Cookie":   []string{"session=raw-cookie"},
				"Content-Type": []string{"application/json"},
			},
			Body:    io.NopCloser(strings.NewReader(`{"access_token":"raw-access","name":"visible"}`)),
			Request: req,
		}, nil
	})), "raw-bearer")

	req, err := http.NewRequest(http.MethodPost, "https://example.com/import?X-Amz-Signature=raw-signature&safe=visible", strings.NewReader(`{"secret":"raw-secret","name":"visible"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(respBody), "raw-access") {
		t.Fatalf("response body was not restored after debug dump: %s", respBody)
	}

	restore()
	got := stdout.String()
	for _, secret := range []string{"raw-bearer", "raw-signature", "raw-secret", "raw-cookie", "raw-access"} {
		if strings.Contains(got, secret) {
			t.Fatalf("debug output leaked %q: %s", secret, got)
		}
	}
	if !strings.Contains(got, "safe=visible") || !strings.Contains(got, `"name":"visible"`) || !strings.Contains(got, "******") {
		t.Fatalf("debug output lost expected context or masks: %s", got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func captureStdout(t *testing.T, dst *bytes.Buffer) func() {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(dst, reader)
		close(done)
	}()
	return func() {
		_ = writer.Close()
		<-done
		os.Stdout = original
		_ = reader.Close()
	}
}
