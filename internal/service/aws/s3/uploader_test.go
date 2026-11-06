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

package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"
)

type uploadClientStub struct {
	cloud.TiDBCloudClient
	url       string
	partCount int32
	completed bool
	cancelled bool
	parts     []imp.CompletePart
}

func (c *uploadClientStub) StartUpload(_ context.Context, _ string, _, _, _ *string, count *int32) (*imp.StartUploadResponse, error) {
	c.partCount = *count
	id := "test-upload"
	urls := make([]string, *count)
	for i := range urls {
		urls[i] = c.url
	}
	return &imp.StartUploadResponse{UploadId: &id, UploadUrl: urls}, nil
}

func (c *uploadClientStub) CompleteUpload(_ context.Context, _ string, _ *string, parts *[]imp.CompletePart) error {
	c.completed = true
	c.parts = *parts
	return nil
}

func (c *uploadClientStub) CancelUpload(_ context.Context, _ string, _ *string) error {
	c.cancelled = true
	return nil
}

type uploadTestTransport func(*http.Request) (*http.Response, error)

func (f uploadTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestUploadDestinations(t *testing.T) {
	const origin = "https://storage.example/upload?X-Amz-Signature=test%2Bsignature"
	for _, multipart := range []bool{false, true} {
		for _, status := range []int{http.StatusTemporaryRedirect, http.StatusPermanentRedirect} {
			for _, tc := range []struct {
				name      string
				initial   string
				redirects map[string]string
				wantError bool
			}{
				{name: "HTTPS", initial: origin},
				{name: "HTTP", initial: "http://storage.example/upload", wantError: true},
				{name: "other scheme", initial: "ftp://storage.example/upload", wantError: true},
				{name: "missing host", initial: "https:///upload", wantError: true},
				{name: "same origin", initial: origin, redirects: map[string]string{origin: "https://storage.example/final"}},
				{name: "equivalent origin", initial: origin, redirects: map[string]string{origin: "https://STORAGE.example:443/final"}},
				{name: "cross host", initial: origin, redirects: map[string]string{origin: "https://attacker.example/final"}, wantError: true},
				{name: "cross port", initial: origin, redirects: map[string]string{origin: "https://storage.example:8443/final"}, wantError: true},
				{name: "downgrade", initial: origin, redirects: map[string]string{origin: "http://storage.example/final"}, wantError: true},
				{name: "loopback", initial: origin, redirects: map[string]string{origin: "http://127.0.0.1:8080/final"}, wantError: true},
				{name: "second hop", initial: origin, redirects: map[string]string{origin: "https://storage.example/next", "https://storage.example/next": "https://attacker.example/final"}, wantError: true},
				{name: "redirect loop", initial: origin, redirects: map[string]string{origin: origin}, wantError: true},
			} {
				t.Run(fmt.Sprintf("multipart=%t/%d/%s", multipart, status, tc.name), func(t *testing.T) {
					api := &uploadClientStub{url: tc.initial}
					cfg := NewUploader(api).(*UploaderImpl)
					cfg.PartSize = MinUploadPartSize
					cfg.Concurrency = 2
					var mu sync.Mutex
					var received int
					// Replace only the network transport, preserving the production
					// HTTPS guard and http.Client redirect policy. No network is used.
					guard := cfg.httpClient.GetClient().Transport.(uploadTransport)
					guard.inner = uploadTestTransport(func(req *http.Request) (*http.Response, error) {
						if req.URL.Scheme != "https" || !strings.EqualFold(req.URL.Hostname(), "storage.example") ||
							(req.URL.Port() != "" && req.URL.Port() != "443") {
							return nil, fmt.Errorf("unexpected destination reached network transport: %s", req.URL)
						}
						if req.Method != http.MethodPut {
							return nil, fmt.Errorf("unexpected upload method: %s", req.Method)
						}
						data, err := io.ReadAll(req.Body)
						if err != nil {
							return nil, err
						}
						code := http.StatusOK
						header := http.Header{"ETag": []string{"test-etag"}}
						if target, ok := tc.redirects[req.URL.String()]; ok {
							code = status
							header.Set("Location", target)
						} else {
							mu.Lock()
							received += len(data)
							mu.Unlock()
						}
						return &http.Response{StatusCode: code, Header: header, Body: io.NopCloser(strings.NewReader("")), Request: req}, nil
					})
					cfg.httpClient.SetTransport(guard)
					payload := "synthetic-import-data"
					if multipart {
						payload = strings.Repeat("x", int(MinUploadPartSize)) + payload
					}
					size := int64(len(payload))
					id, err := cfg.Upload(context.Background(), &PutObjectInput{
						ClusterID: "cluster", ContentLength: &size, Body: strings.NewReader(payload), OnProgress: func(float64) {},
					})
					if tc.wantError {
						require.Error(t, err)
						require.NotContains(t, err.Error(), "unexpected destination reached network transport")
						require.Zero(t, received)
						require.False(t, api.completed)
						require.True(t, api.cancelled)
						return
					}
					require.NoError(t, err)
					require.Equal(t, "test-upload", id)
					require.Equal(t, len(payload), received)
					require.True(t, api.completed)
					require.False(t, api.cancelled)
					if multipart {
						require.EqualValues(t, 2, api.partCount)
						require.Len(t, api.parts, 2)
						require.EqualValues(t, 1, api.parts[0].PartNumber)
						require.EqualValues(t, 2, api.parts[1].PartNumber)
					} else {
						require.EqualValues(t, 1, api.partCount)
					}
				})
			}
		}
	}
}
