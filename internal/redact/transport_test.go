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

package redact_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal/redact"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
)

func TestDebugTransportPreservesResponseReadError(t *testing.T) {
	const body = `{"name":"partial"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "100")
		_, _ = io.WriteString(w, body)
	}))
	defer srv.Close()

	for _, debug := range []bool{false, true} {
		t.Run(fmt.Sprintf("debug=%t", debug), func(t *testing.T) {
			resp, err := util.GetResponse(srv.URL, debug)
			require.NoError(t, err)
			defer resp.Body.Close()
			got, err := io.ReadAll(resp.Body)
			require.ErrorIs(t, err, io.ErrUnexpectedEOF)
			require.Equal(t, body, string(got))

			client := resty.New().SetTransport(redact.NewDebugTransport(srv.Client().Transport, debug))
			var result map[string]interface{}
			_, err = client.R().SetResult(&result).Get(srv.URL)
			require.ErrorIs(t, err, io.ErrUnexpectedEOF)
		})
	}
}

type failingReader struct {
	io.Reader
}

func (r failingReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return n, err
}

func TestDumpRequestPreservesReadError(t *testing.T) {
	const body = `{"name":"partial"}`
	req, err := http.NewRequest(http.MethodPost, "https://example.com", failingReader{strings.NewReader(body)})
	require.NoError(t, err)
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/json")
	redact.DumpRequest(req)
	got, err := io.ReadAll(req.Body)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.Equal(t, body, string(got))
}

func TestHeadersMaskURLValuesWithoutMutation(t *testing.T) {
	headers := http.Header{
		"Location":         {"https://alice:URL_PASSWORD@h/file?sig=ABSOLUTE_SIGNATURE&plain=1", "/file?sig=RELATIVE_SIGNATURE"},
		"Referer":          {"https://h/file?X-Amz-Signature=REFERER_SIGNATURE"},
		"Content-Location": {"/file?sig=CONTENT_SIGNATURE"},
	}
	original := headers.Clone()
	masked := redact.Headers(headers)
	require.Equal(t, original, headers)
	for _, secret := range []string{"URL_PASSWORD", "ABSOLUTE_SIGNATURE", "RELATIVE_SIGNATURE", "REFERER_SIGNATURE", "CONTENT_SIGNATURE"} {
		require.NotContains(t, fmt.Sprint(masked), secret)
	}
	require.Contains(t, masked.Get("Location"), "plain=1")
	require.Len(t, masked.Values("Location"), 2)
}

func TestDownloadRedirectMasksURLHeaders(t *testing.T) {
	for _, relative := range []bool{false, true} {
		t.Run(fmt.Sprintf("relative=%t", relative), func(t *testing.T) {
			const body = "column\nvalue\n"
			received := make(chan *http.Request, 1)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/start" {
					location := "/file?X-Amz-Signature=DESTINATION_SIGNATURE&plain=1"
					if !relative {
						location = "http://" + r.Host + location
					}
					w.Header().Set("Location", location)
					w.WriteHeader(http.StatusFound)
					return
				}
				received <- r
				w.Header().Set("Content-Type", "text/csv")
				_, _ = io.WriteString(w, body)
			}))
			defer srv.Close()

			dump := captureStdout(t, func() {
				resp, err := util.GetResponse(srv.URL+"/start?X-Amz-Signature=ORIGINAL_SIGNATURE", true)
				require.NoError(t, err)
				defer resp.Body.Close()
				got, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				require.Equal(t, body, string(got))
			})
			require.NotContains(t, dump, "ORIGINAL_SIGNATURE")
			require.NotContains(t, dump, "DESTINATION_SIGNATURE")
			require.Contains(t, dump, "Location:")
			require.Contains(t, dump, "Referer:")
			require.Contains(t, dump, "plain=1")
			req := <-received
			require.Equal(t, "DESTINATION_SIGNATURE", req.URL.Query().Get("X-Amz-Signature"))
			require.Equal(t, srv.URL+"/start?X-Amz-Signature=ORIGINAL_SIGNATURE", req.Referer())
		})
	}
}

func TestGeneratedClientDebugDoesNotLeakHeaders(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(fmt.Sprintf("transportDebug=%t", debug), func(t *testing.T) {
			received := make(chan http.Header, 1)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				received <- r.Header.Clone()
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Set-Cookie", "session=RESPONSE_COOKIE")
				_, _ = io.WriteString(w, `{"clusters":[]}`)
			}))
			defer srv.Close()

			cfg := cluster.NewConfiguration()
			cfg.Debug = true
			cfg.Servers = cluster.ServerConfigurations{{URL: srv.URL}}
			cfg.HTTPClient = &http.Client{Transport: redact.NewDebugTransport(srv.Client().Transport, debug)}
			cfg.AddDefaultHeader("Authorization", "Bearer SDK_TOKEN")
			cfg.AddDefaultHeader("Cookie", "session=REQUEST_COOKIE")
			var logs bytes.Buffer
			original := log.Writer()
			log.SetOutput(&logs)
			defer log.SetOutput(original)
			dump := captureStdout(t, func() {
				_, resp, err := cluster.NewAPIClient(cfg).ClusterServiceAPI.ClusterServiceListClusters(context.Background()).Execute()
				require.NoError(t, err)
				defer resp.Body.Close()
				require.Equal(t, "session=RESPONSE_COOKIE", resp.Header.Get("Set-Cookie"))
			})
			for _, secret := range []string{"SDK_TOKEN", "REQUEST_COOKIE", "RESPONSE_COOKIE"} {
				require.NotContains(t, dump+logs.String(), secret)
			}
			headers := <-received
			require.Equal(t, "Bearer SDK_TOKEN", headers.Get("Authorization"))
			require.Equal(t, "session=REQUEST_COOKIE", headers.Get("Cookie"))
		})
	}
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	output, err := os.CreateTemp(t.TempDir(), "stdout")
	require.NoError(t, err)
	original := os.Stdout
	os.Stdout = output
	defer func() {
		os.Stdout = original
		output.Close()
	}()
	f()
	_, err = output.Seek(0, io.SeekStart)
	require.NoError(t, err)
	data, err := io.ReadAll(output)
	require.NoError(t, err)
	return string(data)
}
