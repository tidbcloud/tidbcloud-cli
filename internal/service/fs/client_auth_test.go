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

package fs

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientInjectsZeroInstanceIDHeader(t *testing.T) {
	var capturedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"entries": []}`))
	}))
	defer server.Close()

	instanceID := "test-instance-123"
	client := NewClient(server.URL, &http.Client{}, "", instanceID)

	_, err := client.List("/")
	if err != nil {
		t.Logf("Expected error: %v", err)
	}

	if capturedHeaders == nil {
		t.Fatal("No request was made to server")
	}

	instanceHeader := capturedHeaders.Get("X-TIDBCLOUD-ZERO-INSTANCE-ID")
	if instanceHeader != instanceID {
		t.Errorf("X-TIDBCLOUD-ZERO-INSTANCE-ID header = %q, want %q", instanceHeader, instanceID)
	}

	if capturedHeaders.Get("X-TIDBCLOUD-CLUSTER-ID") != "" {
		t.Error("X-TIDBCLOUD-CLUSTER-ID should not be set when zero instance ID is used")
	}
}

func TestClientInjectsClusterIDHeader(t *testing.T) {
	var capturedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"entries": []}`))
	}))
	defer server.Close()

	clusterID := "test-cluster-456"
	client := NewClient(server.URL, &http.Client{}, clusterID, "")

	_, err := client.List("/")
	if err != nil {
		t.Logf("Expected error: %v", err)
	}

	if capturedHeaders == nil {
		t.Fatal("No request was made to server")
	}

	clusterHeader := capturedHeaders.Get("X-TIDBCLOUD-CLUSTER-ID")
	if clusterHeader != clusterID {
		t.Errorf("X-TIDBCLOUD-CLUSTER-ID header = %q, want %q", clusterHeader, clusterID)
	}

	if capturedHeaders.Get("X-TIDBCLOUD-ZERO-INSTANCE-ID") != "" {
		t.Error("X-TIDBCLOUD-ZERO-INSTANCE-ID should not be set when cluster ID is used")
	}
}

func TestClientClusterIDWinsOverZeroInstanceID(t *testing.T) {
	var capturedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"entries": []}`))
	}))
	defer server.Close()

	clusterID := "test-cluster-456"
	instanceID := "test-instance-123"
	client := NewClient(server.URL, &http.Client{}, clusterID, instanceID)

	_, err := client.List("/")
	if err != nil {
		t.Logf("Expected error: %v", err)
	}

	if capturedHeaders.Get("X-TIDBCLOUD-CLUSTER-ID") != clusterID {
		t.Errorf("X-TIDBCLOUD-CLUSTER-ID header = %q, want %q", capturedHeaders.Get("X-TIDBCLOUD-CLUSTER-ID"), clusterID)
	}
	if capturedHeaders.Get("X-TIDBCLOUD-ZERO-INSTANCE-ID") != "" {
		t.Error("X-TIDBCLOUD-ZERO-INSTANCE-ID should not be set when both IDs are provided (cluster wins)")
	}
}

func TestClientNoDBHeadersWhenNotConfigured(t *testing.T) {
	var capturedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"entries": []}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, &http.Client{}, "", "")

	_, err := client.List("/")
	if err != nil {
		t.Logf("Expected error: %v", err)
	}

	if capturedHeaders.Get("X-TIDBCLOUD-CLUSTER-ID") != "" {
		t.Error("X-TIDBCLOUD-CLUSTER-ID should not be set when not configured")
	}
	if capturedHeaders.Get("X-TIDBCLOUD-ZERO-INSTANCE-ID") != "" {
		t.Error("X-TIDBCLOUD-ZERO-INSTANCE-ID should not be set when not configured")
	}
}
