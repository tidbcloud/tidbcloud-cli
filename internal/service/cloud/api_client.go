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
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/br"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/branch"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"
	imp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/import"

	apiClient "github.com/c4pt0r/go-tidbcloud-sdk-v1/client"
	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/icholy/digest"
)

const (
	DefaultApiUrl             = "https://" + apiClient.DefaultHost
	DefaultServerlessEndpoint = "https://serverless.tidbapi.com"
	DefaultIAMEndpoint        = "https://iam.tidbapi.com"
	userAgent                 = "User-Agent"
)

type TiDBCloudClient interface {
	CreateCluster(ctx context.Context, body *cluster.TidbCloudOpenApiserverlessv1beta1Cluster) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	DeleteCluster(ctx context.Context, clusterId string) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	GetCluster(ctx context.Context, clusterId string) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	ListClusters(ctx context.Context, filter *string, pageSize *int32, pageToken *string, orderBy *string, skip *int32) (*cluster.TidbCloudOpenApiserverlessv1beta1ListClustersResponse, error)

	PartialUpdateCluster(ctx context.Context, clusterId string, body *cluster.V1beta1ServerlessServicePartialUpdateClusterBody) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	ListProviderRegions(ctx context.Context) (*cluster.TidbCloudOpenApiserverlessv1beta1ListRegionsResponse, error)

	ListProjects(ctx context.Context, pageSize *int32, pageToken *string) (*iam.ApiListProjectsRsp, error)

	CancelImport(ctx context.Context, clusterId string, id string) error

	CreateImport(ctx context.Context, clusterId string, body *imp.ImportServiceCreateImportBody) (*imp.Import, error)

	GetImport(ctx context.Context, clusterId string, id string) (*imp.Import, error)

	ListImports(ctx context.Context, clusterId string, pageSize *int32, pageToken, orderBy *string) (*imp.ListImportsResp, error)

	GetBranch(ctx context.Context, clusterId, branchId string) (*branch.Branch, error)

	ListBranches(ctx context.Context, clusterId string, pageSize *int32, pageToken *string) (*branch.ListBranchesResponse, error)

	CreateBranch(ctx context.Context, clusterId string, body *branch.Branch) (*branch.Branch, error)

	DeleteBranch(ctx context.Context, clusterId string, branchId string) (*branch.Branch, error)

	Chat(params *pingchatOp.ChatParams, opts ...pingchatOp.ClientOption) (*pingchatOp.ChatOK, error)

	DeleteBackup(ctx context.Context, backupId string) (*br.V1beta1Backup, error)

	GetBackup(ctx context.Context, backupId string) (*br.V1beta1Backup, error)

	ListBackups(ctx context.Context, clusterId *string, pageSize *int32, pageToken *string) (*br.V1beta1ListBackupsResponse, error)

	Restore(ctx context.Context, body *br.V1beta1RestoreRequest) (*br.V1beta1RestoreResponse, error)

	StartUpload(ctx context.Context, clusterId string, fileName, targetDatabase, targetTable *string, partNumber *int32) (*imp.StartUploadResponse, error)

	CompleteUpload(ctx context.Context, clusterId string, uploadId *string, parts *[]imp.CompletePart) error

	CancelUpload(ctx context.Context, clusterId string, uploadId *string) error

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
	brc *br.APIClient
	sc  *cluster.APIClient
	sic *imp.APIClient
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

func (d *ClientDelegate) CreateCluster(ctx context.Context, body *cluster.TidbCloudOpenApiserverlessv1beta1Cluster) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	r := d.sc.ServerlessServiceAPI.ServerlessServiceCreateCluster(ctx)
	if body != nil {
		r = r.Cluster(*body)
	}
	c, h, err := r.Execute()
	return c, parseError(err, h)
}

func (d *ClientDelegate) DeleteCluster(ctx context.Context, clusterId string) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	c, h, err := d.sc.ServerlessServiceAPI.ServerlessServiceDeleteCluster(ctx, clusterId).Execute()
	return c, parseError(err, h)
}

func (d *ClientDelegate) GetCluster(ctx context.Context, clusterId string) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	c, h, err := d.sc.ServerlessServiceAPI.ServerlessServiceGetCluster(ctx, clusterId).Execute()
	return c, parseError(err, h)
}

func (d *ClientDelegate) ListProviderRegions(ctx context.Context) (*cluster.TidbCloudOpenApiserverlessv1beta1ListRegionsResponse, error) {
	resp, h, err := d.sc.ServerlessServiceAPI.ServerlessServiceListRegions(ctx).Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) ListClusters(ctx context.Context, filter *string, pageSize *int32, pageToken *string, orderBy *string, skip *int32) (*cluster.TidbCloudOpenApiserverlessv1beta1ListClustersResponse, error) {
	r := d.sc.ServerlessServiceAPI.ServerlessServiceListClusters(ctx)
	if filter != nil {
		r = r.Filter(*filter)
	}
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	if orderBy != nil {
		r = r.OrderBy(*orderBy)
	}
	if skip != nil {
		r = r.Skip(*skip)
	}
	resp, h, err := r.Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) PartialUpdateCluster(ctx context.Context, clusterId string, body *cluster.V1beta1ServerlessServicePartialUpdateClusterBody) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	r := d.sc.ServerlessServiceAPI.ServerlessServicePartialUpdateCluster(ctx, clusterId)
	if body != nil {
		r = r.Body(*body)
	}
	c, h, err := r.Execute()
	return c, parseError(err, h)
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

func (d *ClientDelegate) CancelImport(ctx context.Context, clusterId string, id string) error {
	_, h, err := d.sic.ImportServiceAPI.ImportServiceCancelImport(ctx, clusterId, id).Execute()
	return parseError(err, h)
}

func (d *ClientDelegate) CreateImport(ctx context.Context, clusterId string, body *imp.ImportServiceCreateImportBody) (*imp.Import, error) {
	r := d.sic.ImportServiceAPI.ImportServiceCreateImport(ctx, clusterId)
	if body != nil {
		r = r.Body(*body)
	}
	i, h, err := r.Execute()
	return i, parseError(err, h)
}

func (d *ClientDelegate) GetImport(ctx context.Context, clusterId string, id string) (*imp.Import, error) {
	i, h, err := d.sic.ImportServiceAPI.ImportServiceGetImport(ctx, clusterId, id).Execute()
	return i, parseError(err, h)
}

func (d *ClientDelegate) ListImports(ctx context.Context, clusterId string, pageSize *int32, pageToken, orderBy *string) (*imp.ListImportsResp, error) {
	r := d.sic.ImportServiceAPI.ImportServiceListImports(ctx, clusterId)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	if orderBy != nil {
		r = r.OrderBy(*orderBy)
	}
	is, h, err := r.Execute()
	return is, parseError(err, h)
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

func (d *ClientDelegate) DeleteBackup(ctx context.Context, backupId string) (*br.V1beta1Backup, error) {
	b, h, err := d.brc.BackupRestoreServiceAPI.BackupRestoreServiceDeleteBackup(ctx, backupId).Execute()
	return b, parseError(err, h)
}

func (d *ClientDelegate) GetBackup(ctx context.Context, backupId string) (*br.V1beta1Backup, error) {
	b, h, err := d.brc.BackupRestoreServiceAPI.BackupRestoreServiceGetBackup(ctx, backupId).Execute()
	return b, parseError(err, h)
}

func (d *ClientDelegate) ListBackups(ctx context.Context, clusterId *string, pageSize *int32, pageToken *string) (*br.V1beta1ListBackupsResponse, error) {
	r := d.brc.BackupRestoreServiceAPI.BackupRestoreServiceListBackups(ctx)
	if clusterId != nil {
		r = r.ClusterId(*clusterId)
	}
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	bs, h, err := r.Execute()
	return bs, parseError(err, h)
}

func (d *ClientDelegate) Restore(ctx context.Context, body *br.V1beta1RestoreRequest) (*br.V1beta1RestoreResponse, error) {
	r := d.brc.BackupRestoreServiceAPI.BackupRestoreServiceRestore(ctx)
	if body != nil {
		r = r.Body(*body)
	}
	bs, h, err := r.Execute()
	return bs, parseError(err, h)
}

func (d *ClientDelegate) StartUpload(ctx context.Context, clusterId string, fileName, targetDatabase, targetTable *string, partNumber *int32) (*imp.StartUploadResponse, error) {
	r := d.sic.ImportServiceAPI.ImportServiceStartUpload(ctx, clusterId)
	if fileName != nil {
		r = r.FileName(*fileName)
	}
	if targetDatabase != nil {
		r = r.TargetDatabase(*targetDatabase)
	}
	if targetTable != nil {
		r = r.TargetTable(*targetTable)
	}
	if partNumber != nil {
		r = r.PartNumber(*partNumber)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) CompleteUpload(ctx context.Context, clusterId string, uploadId *string, parts *[]imp.CompletePart) error {
	r := d.sic.ImportServiceAPI.ImportServiceCompleteUpload(ctx, clusterId)
	if uploadId != nil {
		r = r.UploadId(*uploadId)
	}
	if parts != nil {
		r = r.Parts(*parts)
	}
	_, h, err := r.Execute()
	return parseError(err, h)
}

func (d *ClientDelegate) CancelUpload(ctx context.Context, clusterId string, uploadId *string) error {
	r := d.sic.ImportServiceAPI.ImportServiceCancelUpload(ctx, clusterId)
	if uploadId != nil {
		r = r.UploadId(*uploadId)
	}
	_, h, err := r.Execute()
	return parseError(err, h)
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

func NewApiClient(rt http.RoundTripper, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*branch.APIClient, *cluster.APIClient, *pingchatClient.TidbcloudPingchat, *br.APIClient, *imp.APIClient, *export.APIClient, *iam.APIClient, error) {
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

	iamCfg := iam.NewConfiguration()
	iamCfg.HTTPClient = httpclient
	iamURL, err := prop.ValidateApiUrl(iamEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	iamCfg.Host = iamURL.Host

	clusterCfg := cluster.NewConfiguration()
	clusterCfg.HTTPClient = httpclient
	clusterCfg.Host = serverlessURL.Host

	branchCfg := branch.NewConfiguration()
	branchCfg.HTTPClient = httpclient
	branchCfg.Host = serverlessURL.Host

	exportCfg := export.NewConfiguration()
	exportCfg.HTTPClient = httpclient
	exportCfg.Host = serverlessURL.Host

	importCfg := imp.NewConfiguration()
	importCfg.HTTPClient = httpclient
	importCfg.Host = serverlessURL.Host

	backupRestoreCfg := br.NewConfiguration()
	backupRestoreCfg.HTTPClient = httpclient
	backupRestoreCfg.Host = serverlessURL.Host

	return branch.NewAPIClient(branchCfg), cluster.NewAPIClient(clusterCfg),
		pingchatClient.New(transport, strfmt.Default), br.NewAPIClient(backupRestoreCfg),
		imp.NewAPIClient(importCfg), export.NewAPIClient(exportCfg),
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
