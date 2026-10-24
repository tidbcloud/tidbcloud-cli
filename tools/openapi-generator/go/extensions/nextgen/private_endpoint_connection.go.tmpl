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

package nextgen

// Hand-written extension for the mgmt private endpoint write and detail APIs.
// Source: tools/openapi-generator/go/extensions/nextgen/private_endpoint_connection.go.tmpl
// The Makefile restores this extension after regenerating the SDK.
// Reuse the generated client's authentication transport and request machinery.
import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

func (a *PrivateEndpointConnectionServiceAPIService) CreateConnection(ctx context.Context, tidbID string, body *Nextgenv1beta2PrivateEndpointConnection) (*Nextgenv1beta2PrivateEndpointConnection, *http.Response, error) {
	return a.connectionRequest(ctx, http.MethodPost, tidbID, "", body)
}
func (a *PrivateEndpointConnectionServiceAPIService) GetConnection(ctx context.Context, tidbID, connectionID string) (*Nextgenv1beta2PrivateEndpointConnection, *http.Response, error) {
	return a.connectionRequest(ctx, http.MethodGet, tidbID, connectionID, nil)
}
func (a *PrivateEndpointConnectionServiceAPIService) DeleteConnection(ctx context.Context, tidbID, connectionID string) (*Nextgenv1beta2PrivateEndpointConnection, *http.Response, error) {
	return a.connectionRequest(ctx, http.MethodDelete, tidbID, connectionID, nil)
}
func (a *PrivateEndpointConnectionServiceAPIService) connectionRequest(ctx context.Context, method, tidbID, connectionID string, body *Nextgenv1beta2PrivateEndpointConnection) (*Nextgenv1beta2PrivateEndpointConnection, *http.Response, error) {
	base, err := a.client.cfg.ServerURLWithContext(ctx, "PrivateEndpointConnectionServiceAPIService")
	if err != nil {
		return nil, nil, err
	}
	path := base + "/tidbs/" + url.PathEscape(tidbID) + "/privateEndpointConnections"
	if connectionID != "" {
		path += "/" + url.PathEscape(connectionID)
	}
	var payload interface{}
	if body != nil {
		payload = body
	}
	req, err := a.client.prepareRequest(ctx, path, method, payload, map[string]string{"Accept": "application/json", "Content-Type": "application/json"}, url.Values{}, url.Values{}, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := a.client.callAPI(req)
	if err != nil || resp == nil {
		return nil, resp, err
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(data))
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode >= 300 {
		return nil, resp, &GenericOpenAPIError{body: data, error: resp.Status}
	}
	var result *Nextgenv1beta2PrivateEndpointConnection
	err = a.client.decode(&result, data, resp.Header.Get("Content-Type"))
	return result, resp, err
}
