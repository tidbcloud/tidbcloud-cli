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

//nolint:bodyclose // parseError closes every non-nil response body
package cloud

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/prop"
	"github.com/tidbcloud/tidbcloud-cli/internal/version"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"
)

const DefaultNextGenEndpoint = "https://cloud.tidbapi.com"

type NextGenClient interface {
	CreateTiDB(ctx context.Context, body *nextgen.Nextgenv1beta2Tidb) (*nextgen.Nextgenv1beta2Tidb, error)
	ListTiDBs(ctx context.Context, servicePlan nextgen.TidbServiceListTidbsServicePlanParameter, pageSize *int32, pageToken *string) (*nextgen.V1beta2ListTidbsResponse, error)
	GetTiDB(ctx context.Context, tidbID string) (*nextgen.Nextgenv1beta2Tidb, error)
	UpdateTiDB(ctx context.Context, tidbID string, body *nextgen.TheTiDBCloudPremiumInstanceToUpdate) (*nextgen.Nextgenv1beta2Tidb, error)
	DeleteTiDB(ctx context.Context, tidbID string) (*nextgen.Nextgenv1beta2Tidb, error)
	ResetRootPassword(ctx context.Context, tidbID string, body *nextgen.TidbServiceResetRootPasswordBody) error
	ListRegions(ctx context.Context, servicePlan nextgen.TidbServiceListTidbsServicePlanParameter, pageSize *int32, pageToken *string) (*nextgen.V1beta2ListRegionsResponse, error)
	GetCACertificateDownloadURL(ctx context.Context, tidbID string) (*nextgen.V1beta2CaCertificateDownloadUrl, error)
	UpdatePublicConnectionSetting(ctx context.Context, tidbID string, body *nextgen.V1beta2PublicConnectionSetting) (*nextgen.V1beta2PublicConnectionSetting, error)
	GetCmekAccessIAMPrincipal(ctx context.Context, regionID string) (*nextgen.V1beta2CmekAccessIamPrincipal, error)
	VerifyCmekAccessIAMPrincipal(ctx context.Context, key *nextgen.V1beta2CustomerManagedEncryptionKey) (*nextgen.V1beta2VerifyCmekAccessIamPrincipalResponse, error)
}

type NextGenClientDelegate struct {
	client *nextgen.APIClient
}

func NewNextGenClientDelegateWithToken(token, endpoint, caCertPath string) (*NextGenClientDelegate, error) {
	transport, err := nextGenTransport(caCertPath)
	if err != nil {
		return nil, err
	}
	return newNextGenClientDelegate(NewTransportWithBearToken(NewDebugTransport(transport), token), endpoint)
}

func NewNextGenClientDelegateWithAPIKey(publicKey, privateKey, endpoint, caCertPath string) (*NextGenClientDelegate, error) {
	transport, err := nextGenTransport(caCertPath)
	if err != nil {
		return nil, err
	}
	return newNextGenClientDelegate(NewTransportWithDigest(NewDebugTransport(transport), publicKey, privateKey), endpoint)
}

func nextGenTransport(caCertPath string) (http.RoundTripper, error) {
	if caCertPath == "" {
		return http.DefaultTransport, nil
	}
	pem, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("read NextGen CA certificate: %w", err)
	}
	roots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("load system CA certificates: %w", err)
	}
	if !roots.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("NextGen CA certificate is not valid PEM")
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12, RootCAs: roots}
	if transport.TLSClientConfig != nil {
		tlsConfig = transport.TLSClientConfig.Clone()
		if tlsConfig.MinVersion < tls.VersionTLS12 {
			tlsConfig.MinVersion = tls.VersionTLS12
		}
		tlsConfig.RootCAs = roots
	}
	transport.TLSClientConfig = tlsConfig
	return transport, nil
}

func newNextGenClientDelegate(rt http.RoundTripper, endpoint string) (*NextGenClientDelegate, error) {
	u, err := prop.ValidateNextGenApiUrl(endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = strings.TrimRight(u.Path, "/")
	if !strings.HasSuffix(u.Path, "/v1beta2") {
		u.Path += "/v1beta2"
	}

	cfg := nextgen.NewConfiguration()
	cfg.HTTPClient = &http.Client{Transport: rt}
	cfg.UserAgent = fmt.Sprintf("%s/%s", config.CliName, version.Version)
	cfg.Servers[0].URL = u.String()
	return &NextGenClientDelegate{client: nextgen.NewAPIClient(cfg)}, nil
}

func (d *NextGenClientDelegate) CreateTiDB(ctx context.Context, body *nextgen.Nextgenv1beta2Tidb) (*nextgen.Nextgenv1beta2Tidb, error) {
	request := d.client.TiDBCloudPremiumInstanceAPI.TidbServiceCreateTidb(ctx)
	if body != nil {
		request = request.Tidb(*body)
	}
	result, response, err := request.Execute()
	return result, parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) ListTiDBs(ctx context.Context, servicePlan nextgen.TidbServiceListTidbsServicePlanParameter, pageSize *int32, pageToken *string) (*nextgen.V1beta2ListTidbsResponse, error) {
	request := d.client.TiDBCloudPremiumInstanceAPI.TidbServiceListTidbs(ctx).ServicePlan(servicePlan)
	if pageSize != nil {
		request = request.PageSize(*pageSize)
	}
	if pageToken != nil {
		request = request.PageToken(*pageToken)
	}
	result, response, err := request.Execute()
	return result, parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) GetTiDB(ctx context.Context, tidbID string) (*nextgen.Nextgenv1beta2Tidb, error) {
	result, response, err := d.client.TiDBCloudPremiumInstanceAPI.TidbServiceGetTidb(ctx, tidbID).Execute()
	return result, parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) UpdateTiDB(ctx context.Context, tidbID string, body *nextgen.TheTiDBCloudPremiumInstanceToUpdate) (*nextgen.Nextgenv1beta2Tidb, error) {
	request := d.client.TiDBCloudPremiumInstanceAPI.TidbServiceUpdateTidb(ctx, tidbID)
	if body != nil {
		request = request.Tidb(*body)
	}
	result, response, err := request.Execute()
	return result, parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) DeleteTiDB(ctx context.Context, tidbID string) (*nextgen.Nextgenv1beta2Tidb, error) {
	result, response, err := d.client.TiDBCloudPremiumInstanceAPI.TidbServiceDeleteTidb(ctx, tidbID).Execute()
	return result, parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) ResetRootPassword(ctx context.Context, tidbID string, body *nextgen.TidbServiceResetRootPasswordBody) error {
	request := d.client.TiDBCloudPremiumInstanceAPI.TidbServiceResetRootPassword(ctx, tidbID)
	if body != nil {
		request = request.Body(*body)
	}
	_, response, err := request.Execute()
	return parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) ListRegions(ctx context.Context, servicePlan nextgen.TidbServiceListTidbsServicePlanParameter, pageSize *int32, pageToken *string) (*nextgen.V1beta2ListRegionsResponse, error) {
	request := d.client.RegionAPI.RegionServiceListRegions(ctx).ServicePlan(servicePlan)
	if pageSize != nil {
		request = request.PageSize(*pageSize)
	}
	if pageToken != nil {
		request = request.PageToken(*pageToken)
	}
	result, response, err := request.Execute()
	return result, parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) GetCACertificateDownloadURL(ctx context.Context, tidbID string) (*nextgen.V1beta2CaCertificateDownloadUrl, error) {
	result, response, err := d.client.TiDBCloudPremiumInstanceAPI.TidbServiceGetCaCertificateDownloadUrl(ctx, tidbID).Execute()
	return result, parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) UpdatePublicConnectionSetting(ctx context.Context, tidbID string, body *nextgen.V1beta2PublicConnectionSetting) (*nextgen.V1beta2PublicConnectionSetting, error) {
	request := d.client.PublicConnectionSettingServiceAPI.PublicConnectionSettingServiceUpdatePublicConnectionSetting(ctx, tidbID)
	if body != nil {
		request = request.PublicConnectionSetting(*body)
	}
	result, response, err := request.Execute()
	return result, parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) GetCmekAccessIAMPrincipal(ctx context.Context, regionID string) (*nextgen.V1beta2CmekAccessIamPrincipal, error) {
	result, response, err := d.client.CustomerManagedEncryptionKeyServiceAPI.CustomerManagedEncryptionKeyServiceGetCmekAccessIamPrincipal(ctx).
		RegionId(regionID).
		ServicePlan(nextgen.TIDBSERVICELISTTIDBSSERVICEPLANPARAMETER_PREMIUM).
		Execute()
	return result, parseNextGenError(err, response)
}

func (d *NextGenClientDelegate) VerifyCmekAccessIAMPrincipal(ctx context.Context, key *nextgen.V1beta2CustomerManagedEncryptionKey) (*nextgen.V1beta2VerifyCmekAccessIamPrincipalResponse, error) {
	result, response, err := d.client.CustomerManagedEncryptionKeyServiceAPI.CustomerManagedEncryptionKeyServiceVerifyCmekAccessIamPrincipal(ctx).
		Key(*key).Execute()
	return result, parseNextGenError(err, response)
}
