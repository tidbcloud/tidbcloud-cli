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
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/tidbcloud/tidbcloud-cli/internal/config"
)

// captureStdout swaps os.Stdout for a pipe while fn runs and returns the
// captured bytes. Not safe for concurrent use; these tests must not run in
// parallel.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	done := make(chan struct{})
	var buf bytes.Buffer
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	fn()

	_ = w.Close()
	<-done
	os.Stdout = old
	_ = r.Close()
	return buf.String()
}

func TestDebugTransport_RedactsBearerToken(t *testing.T) {
	t.Setenv(config.DebugEnv, "1")

	const token = "leakyBearerTokenABCDEF.payload.signature"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Set-Cookie", "sid=responseLeakCookie; HttpOnly")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"responseAccessTokenSecret","user":"alice"}`))
	}))
	defer srv.Close()

	transport := NewBearTokenTransport(token)
	client := &http.Client{Transport: transport}

	out := captureStdout(t, func() {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL+"/v1/x", nil)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("client.Do: %v", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	})

	if strings.Contains(out, token) {
		t.Fatalf("raw bearer token leaked to stdout:\n%s", out)
	}
	if !strings.Contains(out, "Bearer [REDACTED]") {
		t.Fatalf("expected 'Bearer [REDACTED]' in stdout, got:\n%s", out)
	}
	if strings.Contains(out, "responseAccessTokenSecret") {
		t.Fatalf("response access_token leaked to stdout:\n%s", out)
	}
	if strings.Contains(out, "responseLeakCookie") {
		t.Fatalf("Set-Cookie value leaked to stdout:\n%s", out)
	}
}

func TestDebugTransport_NoDebug_NoOutput(t *testing.T) {
	t.Setenv(config.DebugEnv, "")

	const token = "shouldNotBePrintedToken"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &http.Client{Transport: NewBearTokenTransport(token)}

	out := captureStdout(t, func() {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("client.Do: %v", err)
		}
		_ = resp.Body.Close()
	})

	if out != "" {
		t.Fatalf("expected no stdout output when debug off, got %q", out)
	}
}

func TestDebugTransport_RedactsS3ImportSecret(t *testing.T) {
	t.Setenv(config.DebugEnv, "1")

	const (
		token             = "bearerForS3Test"
		secretAccessValue = "wJalrXUtnFEMI_S3_SECRET_LEAK_CANARY"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// drain the body so DumpRequestOut sees something to dump
		_, _ = io.Copy(io.Discard, r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"imp-1"}`))
	}))
	defer srv.Close()

	body := bytes.NewBufferString(
		`{"source":{"type":"S3","s3":{"uri":"s3://b/k","accessKey":{"id":"AKIATEST","secret":"` + secretAccessValue + `"}}}}`,
	)

	client := &http.Client{Transport: NewBearTokenTransport(token)}

	out := captureStdout(t, func() {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL+"/v1/imports", body)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("client.Do: %v", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	})

	if strings.Contains(out, secretAccessValue) {
		t.Fatalf("S3 secret leaked to stdout:\n%s", out)
	}
	if !strings.Contains(out, "AKIATEST") {
		t.Fatalf("access key id should remain visible (it is an identifier, not a secret):\n%s", out)
	}
	if strings.Contains(out, token) {
		t.Fatalf("bearer token leaked to stdout:\n%s", out)
	}
}

func TestDebugTransport_BodyForwardedToServer(t *testing.T) {
	t.Setenv(config.DebugEnv, "1")

	const payload = `{"clusterId":"c-1","secret":"x"}`

	var received []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		received = b
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &http.Client{Transport: NewBearTokenTransport("t")}

	captureStdout(t, func() {
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			srv.URL+"/",
			bytes.NewBufferString(payload),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("client.Do: %v", err)
		}
		_ = resp.Body.Close()
	})

	if string(received) != payload {
		t.Fatalf("server received %q, want %q (debug dump must not alter wire body)", string(received), payload)
	}
}
