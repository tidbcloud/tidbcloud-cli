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
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestConfigureRestyDebugRedactsLogs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "session=raw-cookie")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"access_token":"raw-access-token","name":"visible"}`)
	}))
	defer server.Close()

	logger := &bufferLogger{}
	client := resty.New()
	client.SetLogger(logger)
	ConfigureRestyDebug(client, true)

	resp, err := client.R().
		SetHeader("Authorization", "Bearer raw-bearer").
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"client_secret": "raw-client-secret",
			"name":          "visible",
		}).
		Post(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsSuccess() {
		t.Fatalf("unexpected status: %s", resp.Status())
	}

	got := logger.String()
	for _, secret := range []string{"raw-bearer", "raw-client-secret", "raw-cookie", "raw-access-token"} {
		if strings.Contains(got, secret) {
			t.Fatalf("debug log leaked %q: %s", secret, got)
		}
	}
	if !strings.Contains(got, "visible") || !strings.Contains(got, "******") {
		t.Fatalf("debug log did not keep expected context and masks: %s", got)
	}
}

type bufferLogger struct {
	bytes.Buffer
}

func (l *bufferLogger) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(&l.Buffer, format, args...)
}

func (l *bufferLogger) Warnf(format string, args ...interface{}) {
	fmt.Fprintf(&l.Buffer, format, args...)
}

func (l *bufferLogger) Debugf(format string, args ...interface{}) {
	fmt.Fprintf(&l.Buffer, format, args...)
}
