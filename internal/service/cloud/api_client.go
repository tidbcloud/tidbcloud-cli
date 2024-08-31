// Copyright 2024 PingCAP, Inc.
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

//nolint:bodyclose // still warn even the body is closed
package cloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"os"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/version"
	pingchatClient "tidbcloud-cli/pkg/tidbcloud/pingchat/client"
	pingchatOp "tidbcloud-cli/pkg/tidbcloud/pingchat/client/operations"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/branch"
	serverlessClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client"
	serverlessOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"
	brClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client"
	brOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"
	serverlessImportClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client"
	serverlessImportOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"

	apiClient "github.com/c4pt0r/go-tidbcloud-sdk-v1/client"
	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/icholy/digest"
)

const (
	DefaultApiUrl             = "https://" + apiClient.DefaultHost
	DefaultServerlessEndpoint = "https://" + serverlessClient.DefaultHost
	DefaultIAMEndpoint        = "https://iam.tidbapi.com"
	userAgent                 = "User-Agent"
)

type TiDBCloudClient interface {
	CreateCluster(params *serverlessOp.ServerlessServiceCreateClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceCreateClusterOK, error)

	DeleteCluster(params *serverlessOp.ServerlessServiceDeleteClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceDeleteClusterOK, error)

	GetCluster(params *serverlessOp.ServerlessServiceGetClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceGetClusterOK, error)

	ListClustersOfProject(params *serverlessOp.ServerlessServiceListClustersParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceListClustersOK, error)

	PartialUpdateCluster(params *serverlessOp.ServerlessServicePartialUpdateClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServicePartialUpdateClusterOK, error)

	ListProviderRegions(params *serverlessOp.ServerlessServiceListRegionsParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceListRegionsOK, error)

	ListProjects(ctx context.Context, pageSize *int32, pageToken *string) (*iam.ApiListProjectsRsp, error)

	CancelImport(params *serverlessImportOp.ImportServiceCancelImportParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCancelImportOK, error)

	CreateImport(params *serverlessImportOp.ImportServiceCreateImportParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCreateImportOK, error)

	GetImport(params *serverlessImportOp.ImportServiceGetImportParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceGetImportOK, error)

	ListImports(params *serverlessImportOp.ImportServiceListImportsParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceListImportsOK, error)

	GetBranch(ctx context.Context, clusterId, branchId string) (*branch.Branch, error)

	ListBranches(ctx context.Context, clusterId string, pageSize *int32, pageToken *string) (*branch.ListBranchesResponse, error)

	CreateBranch(ctx context.Context, clusterId string, body *branch.Branch) (*branch.Branch, error)

	DeleteBranch(ctx context.Context, clusterId string, branchId string) (*branch.Branch, error)

	Chat(params *pingchatOp.ChatParams, opts ...pingchatOp.ClientOption) (*pingchatOp.ChatOK, error)

	DeleteBackup(params *brOp.BackupRestoreServiceDeleteBackupParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceDeleteBackupOK, error)

	GetBackup(params *brOp.BackupRestoreServiceGetBackupParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceGetBackupOK, error)

	ListBackups(params *brOp.BackupRestoreServiceListBackupsParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceListBackupsOK, error)

	Restore(params *brOp.BackupRestoreServiceRestoreParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceRestoreOK, error)

	StartUpload(params *serverlessImportOp.ImportServiceStartUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceStartUploadOK, error)

	CompleteUpload(params *serverlessImportOp.ImportServiceCompleteUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCompleteUploadOK, error)

	CancelUpload(params *serverlessImportOp.ImportServiceCancelUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCancelUploadOK, error)

	GetExport(ctx context.Context, clusterId string, exportId string) (*export.Export, error)

	CancelExport(ctx context.Context, clusterId string, exportId string) (*export.Export, error)

	CreateExport(ctx context.Context, clusterId string, body *export.ExportServiceCreateExportBody) (*export.Export, error)

	DeleteExport(ctx context.Context, clusterId string, exportId string) (*export.Export, error)

	ListExports(ctx context.Context, clusterId string, pageSize *int32, pageToken *string, orderBy *string) (*export.ListExportsResponse, error)

	DownloadExport(ctx context.Context, clusterId string, exportId string) (*export.DownloadExportsResponse, error)

	ListSQLUsers(ctx context.Context, clusterID string, pageSize *int32, pageToken *string) (*iam.ApiListSqlUsersRsp, error)

	CreateSQLUser(ctx context.Context, clusterID string, body *iam.ApiCreateSqlUserReq) (*iam.ApiSqlUser, error)

	GetSQLUser(ctx context.Context, clusterID string, userName string) (*iam.ApiSqlUser, error)

	DeleteSQLUser(ctx context.Context, clusterID string, userName string) (*iam.ApiBasicResp, error)

	UpdateSQLUser(ctx context.Context, clusterID string, userName string, body *iam.ApiUpdateSqlUserReq) (*iam.ApiSqlUser, error)
}

type ClientDelegate struct {
	ic  *iam.APIClient
	bc  *branch.APIClient
	pc  *pingchatClient.TidbcloudPingchat
	sc  *serverlessClient.TidbcloudServerless
	brc *brClient.TidbcloudServerless
	sic *serverlessImportClient.TidbcloudServerless
	ec  *export.APIClient
}

func NewClientDelegateWithToken(token string, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewBearTokenTransport(token)
	bc, sc, pc, brc, sic, ec, ic, err := NewApiClient(transport, apiUrl, serverlessEndpoint, iamEndpoint)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		bc:  bc,
		sc:  sc,
		pc:  pc,
		brc: brc,
		ec:  ec,
		ic:  ic,
		sic: sic,
	}, nil
}

func NewClientDelegateWithApiKey(publicKey string, privateKey string, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewDigestTransport(publicKey, privateKey)
	bc, sc, pc, brc, sic, ec, ic, err := NewApiClient(transport, apiUrl, serverlessEndpoint, iamEndpoint)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		bc:  bc,
		sc:  sc,
		pc:  pc,
		brc: brc,
		ec:  ec,
		ic:  ic,
		sic: sic,
	}, nil
}

func (d *ClientDelegate) CreateCluster(params *serverlessOp.ServerlessServiceCreateClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceCreateClusterOK, error) {
	return d.sc.ServerlessService.ServerlessServiceCreateCluster(params, opts...)
}

func (d *ClientDelegate) DeleteCluster(params *serverlessOp.ServerlessServiceDeleteClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceDeleteClusterOK, error) {
	return d.sc.ServerlessService.ServerlessServiceDeleteCluster(params, opts...)
}

func (d *ClientDelegate) GetCluster(params *serverlessOp.ServerlessServiceGetClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceGetClusterOK, error) {
	return d.sc.ServerlessService.ServerlessServiceGetCluster(params, opts...)
}

func (d *ClientDelegate) ListProviderRegions(params *serverlessOp.ServerlessServiceListRegionsParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceListRegionsOK, error) {
	return d.sc.ServerlessService.ServerlessServiceListRegions(params, opts...)
}

func (d *ClientDelegate) ListClustersOfProject(params *serverlessOp.ServerlessServiceListClustersParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceListClustersOK, error) {
	return d.sc.ServerlessService.ServerlessServiceListClusters(params, opts...)
}

func (d *ClientDelegate) PartialUpdateCluster(params *serverlessOp.ServerlessServicePartialUpdateClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServicePartialUpdateClusterOK, error) {
	return d.sc.ServerlessService.ServerlessServicePartialUpdateCluster(params, opts...)
}

func (d *ClientDelegate) ListProjects(ctx context.Context, pageSize *int32, pageToken *string) (*iam.ApiListProjectsRsp, error) {
	r := d.ic.AccountAPI.V1beta1ProjectsGet(ctx)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) CancelImport(params *serverlessImportOp.ImportServiceCancelImportParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCancelImportOK, error) {
	return d.sic.ImportService.ImportServiceCancelImport(params, opts...)
}

func (d *ClientDelegate) CreateImport(params *serverlessImportOp.ImportServiceCreateImportParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCreateImportOK, error) {
	return d.sic.ImportService.ImportServiceCreateImport(params, opts...)
}

func (d *ClientDelegate) GetImport(params *serverlessImportOp.ImportServiceGetImportParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceGetImportOK, error) {
	return d.sic.ImportService.ImportServiceGetImport(params, opts...)
}

func (d *ClientDelegate) ListImports(params *serverlessImportOp.ImportServiceListImportsParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceListImportsOK, error) {
	return d.sic.ImportService.ImportServiceListImports(params, opts...)
}

func (d *ClientDelegate) GetBranch(ctx context.Context, clusterId, branchId string) (*branch.Branch, error) {
	b, h, err := d.bc.BranchServiceAPI.BranchServiceGetBranch(ctx, clusterId, branchId).Execute()
	return b, parseError(err, h)
}

func (d *ClientDelegate) ListBranches(ctx context.Context, clusterId string, pageSize *int32, pageToken *string) (*branch.ListBranchesResponse, error) {
	r := d.bc.BranchServiceAPI.BranchServiceListBranches(ctx, clusterId)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	bs, h, err := r.Execute()
	return bs, parseError(err, h)
}

func (d *ClientDelegate) CreateBranch(ctx context.Context, clusterId string, body *branch.Branch) (*branch.Branch, error) {
	r := d.bc.BranchServiceAPI.BranchServiceCreateBranch(ctx, clusterId)
	if body != nil {
		r = r.Branch(*body)
	}
	b, h, err := r.Execute()
	return b, parseError(err, h)
}

func (d *ClientDelegate) DeleteBranch(ctx context.Context, clusterId string, branchId string) (*branch.Branch, error) {
	b, h, err := d.bc.BranchServiceAPI.BranchServiceDeleteBranch(ctx, clusterId, branchId).Execute()
	return b, parseError(err, h)
}

func (d *ClientDelegate) Chat(params *pingchatOp.ChatParams, opts ...pingchatOp.ClientOption) (*pingchatOp.ChatOK, error) {
	return d.pc.Operations.Chat(params, opts...)
}

func (d *ClientDelegate) DeleteBackup(params *brOp.BackupRestoreServiceDeleteBackupParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceDeleteBackupOK, error) {
	return d.brc.BackupRestoreService.BackupRestoreServiceDeleteBackup(params, opts...)
}

func (d *ClientDelegate) GetBackup(params *brOp.BackupRestoreServiceGetBackupParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceGetBackupOK, error) {
	return d.brc.BackupRestoreService.BackupRestoreServiceGetBackup(params, opts...)
}

func (d *ClientDelegate) ListBackups(params *brOp.BackupRestoreServiceListBackupsParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceListBackupsOK, error) {
	return d.brc.BackupRestoreService.BackupRestoreServiceListBackups(params, opts...)
}

func (d *ClientDelegate) Restore(params *brOp.BackupRestoreServiceRestoreParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceRestoreOK, error) {
	return d.brc.BackupRestoreService.BackupRestoreServiceRestore(params, opts...)
}

func (d *ClientDelegate) StartUpload(params *serverlessImportOp.ImportServiceStartUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceStartUploadOK, error) {
	return d.sic.ImportService.ImportServiceStartUpload(params, opts...)
}

func (d *ClientDelegate) CompleteUpload(params *serverlessImportOp.ImportServiceCompleteUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCompleteUploadOK, error) {
	return d.sic.ImportService.ImportServiceCompleteUpload(params, opts...)
}

func (d *ClientDelegate) CancelUpload(params *serverlessImportOp.ImportServiceCancelUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCancelUploadOK, error) {
	return d.sic.ImportService.ImportServiceCancelUpload(params, opts...)
}

func (d *ClientDelegate) GetExport(ctx context.Context, clusterId string, exportId string) (*export.Export, error) {
	res, h, err := d.ec.ExportServiceAPI.ExportServiceGetExport(ctx, clusterId, exportId).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) CancelExport(ctx context.Context, clusterId string, exportId string) (*export.Export, error) {
	res, h, err := d.ec.ExportServiceAPI.ExportServiceCancelExport(ctx, clusterId, exportId).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) CreateExport(ctx context.Context, clusterId string, body *export.ExportServiceCreateExportBody) (*export.Export, error) {
	r := d.ec.ExportServiceAPI.ExportServiceCreateExport(ctx, clusterId)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) DeleteExport(ctx context.Context, clusterId string, exportId string) (*export.Export, error) {
	res, h, err := d.ec.ExportServiceAPI.ExportServiceDeleteExport(ctx, clusterId, exportId).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) ListExports(ctx context.Context, clusterId string, pageSize *int32, pageToken *string, orderBy *string) (*export.ListExportsResponse, error) {
	r := d.ec.ExportServiceAPI.ExportServiceListExports(ctx, clusterId)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	if orderBy != nil {
		r = r.OrderBy(*orderBy)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) DownloadExport(ctx context.Context, clusterId string, exportId string) (*export.DownloadExportsResponse, error) {
	r := d.ec.ExportServiceAPI.ExportServiceDownloadExport(ctx, clusterId, exportId)
	r = r.Body(make(map[string]interface{}))
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) ListSQLUsers(ctx context.Context, clusterID string, pageSize *int32, pageToken *string) (*iam.ApiListSqlUsersRsp, error) {
	r := d.ic.AccountAPI.V1beta1ClustersClusterIdSqlUsersGet(ctx, clusterID)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) CreateSQLUser(ctx context.Context, clusterId string, body *iam.ApiCreateSqlUserReq) (*iam.ApiSqlUser, error) {
	r := d.ic.AccountAPI.V1beta1ClustersClusterIdSqlUsersPost(ctx, clusterId)
	if body != nil {
		r = r.SqlUser(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) GetSQLUser(ctx context.Context, clusterID string, userName string) (*iam.ApiSqlUser, error) {
	r := d.ic.AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNameGet(ctx, clusterID, userName)
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) DeleteSQLUser(ctx context.Context, clusterID string, userName string) (*iam.ApiBasicResp, error) {
	r := d.ic.AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNameDelete(ctx, clusterID, userName)
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) UpdateSQLUser(ctx context.Context, clusterID string, userName string, body *iam.ApiUpdateSqlUserReq) (*iam.ApiSqlUser, error) {
	r := d.ic.AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNamePatch(ctx, clusterID, userName)
	if body != nil {
		r = r.SqlUser(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func NewApiClient(rt http.RoundTripper, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*branch.APIClient, *serverlessClient.TidbcloudServerless, *pingchatClient.TidbcloudPingchat, *brClient.TidbcloudServerless, *serverlessImportClient.TidbcloudServerless, *export.APIClient, *iam.APIClient, error) {
	httpclient := &http.Client{
		Transport: rt,
	}

	// v1beta api
	u, err := prop.ValidateApiUrl(apiUrl)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	transport := httpTransport.NewWithClient(u.Host, u.Path, []string{u.Scheme}, httpclient)

	// v1beta1 api (serverless)
	serverlessURL, err := prop.ValidateApiUrl(serverlessEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	serverlessTransport := httpTransport.NewWithClient(serverlessURL.Host, serverlessClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	backRestoreTransport := httpTransport.NewWithClient(serverlessURL.Host, brClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	importTransport := httpTransport.NewWithClient(serverlessURL.Host, serverlessImportClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)

	iamCfg := iam.NewConfiguration()
	iamCfg.HTTPClient = httpclient
	iamURL, err := prop.ValidateApiUrl(iamEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	iamCfg.Host = iamURL.Host

	branchCfg := branch.NewConfiguration()
	branchCfg.HTTPClient = httpclient
	branchCfg.Host = serverlessURL.Host

	exportCfg := export.NewConfiguration()
	exportCfg.HTTPClient = httpclient
	exportCfg.Host = serverlessURL.Host

	return branch.NewAPIClient(branchCfg), serverlessClient.New(serverlessTransport, strfmt.Default),
		pingchatClient.New(transport, strfmt.Default), brClient.New(backRestoreTransport, strfmt.Default),
		serverlessImportClient.New(importTransport, strfmt.Default), export.NewAPIClient(exportCfg),
		iam.NewAPIClient(iamCfg), nil
}

func NewDigestTransport(publicKey, privateKey string) http.RoundTripper {
	return NewTransportWithAgent(&digest.Transport{
		Username:  publicKey,
		Password:  privateKey,
		Transport: NewDebugTransport(http.DefaultTransport),
	}, fmt.Sprintf("%s/%s", config.CliName, version.Version))
}

func NewBearTokenTransport(token string) http.RoundTripper {
	return NewTransportWithAgent(NewTransportWithBearToken(NewDebugTransport(http.DefaultTransport), token),
		fmt.Sprintf("%s/%s", config.CliName, version.Version))
}

// NewTransportWithAgent returns a new http.RoundTripper that add the User-Agent header,
// according to https://github.com/go-swagger/go-swagger/issues/1563.
func NewTransportWithAgent(inner http.RoundTripper, userAgent string) http.RoundTripper {
	return &UserAgentTransport{
		inner: inner,
		Agent: userAgent,
	}
}

type UserAgentTransport struct {
	inner http.RoundTripper
	Agent string
}

func (ug *UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set(userAgent, ug.Agent)
	return ug.inner.RoundTrip(r)
}

func NewTransportWithBearToken(inner http.RoundTripper, token string) http.RoundTripper {
	return &BearTokenTransport{
		inner: inner,
		Token: token,
	}
}

type BearTokenTransport struct {
	inner http.RoundTripper
	Token string
}

func (bt *BearTokenTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bt.Token))
	return bt.inner.RoundTrip(r)
}

func NewDebugTransport(inner http.RoundTripper) http.RoundTripper {
	return &DebugTransport{inner: inner}
}

type DebugTransport struct {
	inner http.RoundTripper
}

func (dt *DebugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	debug := os.Getenv(config.DebugEnv) == "true" || os.Getenv(config.DebugEnv) == "1"

	if debug {
		dump, err := httputil.DumpRequestOut(r, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("\n%s", string(dump))
	}

	resp, err := dt.inner.RoundTrip(r)
	if err != nil {
		return resp, err
	}

	if debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return resp, err
		}
		fmt.Printf("%s\n", string(dump))
	}

	return resp, err
}

func parseError(err error, resp *http.Response) error {
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err == nil {
		return nil
	}
	if resp == nil {
		return err
	}
	body, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		return err
	}
	path := "[path]"
	if resp.Request != nil {
		path = fmt.Sprintf("[%s %s]", resp.Request.Method, resp.Request.URL.Path)
	}
	return fmt.Errorf("%s[%s] %s", path, err.Error(), body)
}
