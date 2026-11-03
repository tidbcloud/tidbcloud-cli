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
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIErrorBodyRedacted(t *testing.T) {
	for _, body := range []string{
		`{"password":"body-secret","message":"invalid request","nested":{"access_token":"token-secret"},"url":"https://example.com?X-Amz-Signature=url-secret"}`,
		"unstructured body-secret",
	} {
		resp := &http.Response{
			Status: "400 Bad Request",
			Header: http.Header{"X-Debug-Trace-Id": []string{"trace-123"}},
			Body:   io.NopCloser(strings.NewReader(body)),
		}
		err := parseError(errors.New("SDK may duplicate body-secret"), resp)
		for _, secret := range []string{"body-secret", "token-secret", "url-secret"} {
			require.NotContains(t, err.Error(), secret)
		}
		require.Contains(t, err.Error(), "400 Bad Request")
		require.Contains(t, err.Error(), "trace-123")
	}
}
