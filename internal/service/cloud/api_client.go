// Copyright 2025 PingCAP, Inc.
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

	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/prop"
	"github.com/tidbcloud/tidbcloud-cli/internal/version"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/br"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/branch"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cdc"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/privatelink"

	"github.com/icholy/digest"
)

const (
	DefaultServerlessEndpoint = "https://serverless.tidbapi.com"
	DefaultIAMEndpoint        = "https://iam.tidbapi.com"
)

type TiDBCloudClient interface {
	CreateCluster(ctx context.Context, body *cluster.TidbCloudOpenApiserverlessv1beta1Cluster) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	DeleteCluster(ctx context.Context, clusterId string) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	GetCluster(ctx context.Context, clusterId string, view cluster.ClusterServiceGetClusterViewParameter) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	ListClusters(ctx context.Context, filter *string, pageSize *int32, pageToken *string, orderBy *string, skip *int32) (*cluster.TidbCloudOpenApiserverlessv1beta1ListClustersResponse, error)

	PartialUpdateCluster(ctx context.Context, clusterId string, body *cluster.V1beta1ClusterServicePartialUpdateClusterBody) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	ListProviderRegions(ctx context.Context) (*cluster.TidbCloudOpenApiserverlessv1beta1ListRegionsResponse, error)

	ListProjects(ctx context.Context, pageSize *int32, pageToken *string) (*iam.ApiListProjectsRsp, error)

	CancelImport(ctx context.Context, clusterId string, id string) error

	CreateImport(ctx context.Context, clusterId string, body *imp.ImportServiceCreateImportBody) (*imp.Import, error)

	GetImport(ctx context.Context, clusterId string, id string) (*imp.Import, error)

	ListImports(ctx context.Context, clusterId string, pageSize *int32, pageToken, orderBy *string) (*imp.ListImportsResp, error)

	GetBranch(ctx context.Context, clusterId, branchId string, view branch.BranchServiceGetBranchViewParameter) (*branch.Branch, error)

	ListBranches(ctx context.Context, clusterId string, pageSize *int32, pageToken *string) (*branch.ListBranchesResponse, error)

	CreateBranch(ctx context.Context, clusterId string, body *branch.Branch) (*branch.Branch, error)

	DeleteBranch(ctx context.Context, clusterId string, branchId string) (*branch.Branch, error)

	ResetBranch(ctx context.Context, clusterId string, branchId string) (*branch.Branch, error)

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

	ListExportFiles(ctx context.Context, clusterId string, exportId string, pageSize *int32, pageToken *string, isGenerateUrl bool) (*export.ListExportFilesResponse, error)

	DownloadExportFiles(ctx context.Context, clusterId string, exportId string, body *export.ExportServiceDownloadExportFilesBody) (*export.DownloadExportFilesResponse, error)

	ListSQLUsers(ctx context.Context, clusterID string, pageSize *int32, pageToken *string) (*iam.ApiListSqlUsersRsp, error)

	CreateSQLUser(ctx context.Context, clusterID string, body *iam.ApiCreateSqlUserReq) (*iam.ApiSqlUser, error)

	GetSQLUser(ctx context.Context, clusterID string, userName string) (*iam.ApiSqlUser, error)

	DeleteSQLUser(ctx context.Context, clusterID string, userName string) (*iam.ApiBasicResp, error)

	UpdateSQLUser(ctx context.Context, clusterID string, userName string, body *iam.ApiUpdateSqlUserReq) (*iam.ApiSqlUser, error)

	DownloadAuditLogs(ctx context.Context, clusterID string, body *auditlog.DatabaseAuditLogServiceDownloadAuditLogFilesBody) (*auditlog.DownloadAuditLogFilesResponse, error)

	ListAuditLogs(ctx context.Context, clusterID string, pageSize *int32, pageToken *string, date *string) (*auditlog.ListAuditLogFilesResponse, error)

	// CDC Changefeed methods
	CreateChangefeed(ctx context.Context, clusterId string, body *cdc.ChangefeedServiceCreateChangefeedBody) (*cdc.Changefeed, error)
	DeleteChangefeed(ctx context.Context, clusterId, changefeedId string) (*cdc.Changefeed, error)
	EditChangefeed(ctx context.Context, clusterId, changefeedId string, body *cdc.ChangefeedServiceEditChangefeedBody) (*cdc.Changefeed, error)
	GetChangefeed(ctx context.Context, clusterId, changefeedId string) (*cdc.Changefeed, error)
	ListChangefeeds(ctx context.Context, clusterId string, pageSize *int32, pageToken *string, changefeedType *cdc.ChangefeedServiceListChangefeedsChangefeedTypeParameter) (*cdc.Changefeeds, error)
	StartChangefeed(ctx context.Context, clusterId, changefeedId string) (*cdc.Changefeed, error)
	StopChangefeed(ctx context.Context, clusterId, changefeedId string) (*cdc.Changefeed, error)
	TestChangefeed(ctx context.Context, clusterId string, body *cdc.ChangefeedServiceTestChangefeedBody) (map[string]interface{}, error)
	DescribeSchemaTable(ctx context.Context, clusterId string, body *cdc.ChangefeedServiceDescribeSchemaTableBody) (*cdc.DescribeSchemaTableResp, error)

	CreateAuditLogFilterRule(ctx context.Context, clusterID string, body *auditlog.DatabaseAuditLogServiceCreateAuditLogFilterRuleBody) (*auditlog.AuditLogFilterRule, error)

	DeleteAuditLogFilterRule(ctx context.Context, clusterID, ruleID string) (*auditlog.AuditLogFilterRule, error)

	GetAuditLogFilterRule(ctx context.Context, clusterID, ruleID string) (*auditlog.AuditLogFilterRule, error)

	ListAuditLogFilterRules(ctx context.Context, clusterID string) (*auditlog.ListAuditLogFilterRulesResponse, error)

	UpdateAuditLogFilterRule(ctx context.Context, clusterID, ruleID string, body *auditlog.DatabaseAuditLogServiceUpdateAuditLogFilterRuleBody) (*auditlog.AuditLogFilterRule, error)

	UpdateAuditLogConfig(ctx context.Context, clusterID string, body *auditlog.DatabaseAuditLogServiceUpdateAuditLogConfigBody) (*auditlog.AuditLogConfig, error)

	GetAuditLogConfig(ctx context.Context, clusterID string) (*auditlog.AuditLogConfig, error)
	// ===== Private Link Connection =====
	CreatePrivateLinkConnection(ctx context.Context, clusterId string, body *privatelink.PrivateLinkConnectionServiceCreatePrivateLinkConnectionBody) (*privatelink.PrivateLinkConnection, error)
	DeletePrivateLinkConnection(ctx context.Context, clusterId string, privateLinkConnectionId string) (*privatelink.PrivateLinkConnection, error)
	GetPrivateLinkConnection(ctx context.Context, clusterId string, privateLinkConnectionId string) (*privatelink.PrivateLinkConnection, error)
	ListPrivateLinkConnections(ctx context.Context, clusterId string, pageSize *int32, pageToken *string) (*privatelink.ListPrivateLinkConnectionsResponse, error)
	GetPrivateLinkAvailabilityZones(ctx context.Context, clusterId string) (*privatelink.GetAvailabilityZonesResponse, error)
	// ===== Private Link Connection =====
}

type ClientDelegate struct {
	ic  *iam.APIClient
	bc  *branch.APIClient
	brc *br.APIClient
	sc  *cluster.APIClient
	sic *imp.APIClient
	ec  *export.APIClient
	alc *auditlog.APIClient
	cdc *cdc.APIClient
	plc *privatelink.APIClient
}

func NewClientDelegateWithToken(token string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewBearTokenTransport(token)
	bc, sc, brc, sic, ec, ic, alc, cdc, plc, err := NewApiClient(transport, serverlessEndpoint, iamEndpoint)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		bc:  bc,
		sc:  sc,
		brc: brc,
		ec:  ec,
		ic:  ic,
		sic: sic,
		alc: alc,
		cdc: cdc,
		plc: plc,
	}, nil
}

func NewClientDelegateWithApiKey(publicKey string, privateKey string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewDigestTransport(publicKey, privateKey)
	bc, sc, brc, sic, ec, ic, alc, cdc, plc, err := NewApiClient(transport, serverlessEndpoint, iamEndpoint)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		bc:  bc,
		sc:  sc,
		brc: brc,
		ec:  ec,
		ic:  ic,
		sic: sic,
		alc: alc,
		cdc: cdc,
		plc: plc,
	}, nil
}

func (d *ClientDelegate) CreateCluster(ctx context.Context, body *cluster.TidbCloudOpenApiserverlessv1beta1Cluster) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	r := d.sc.ClusterServiceAPI.ClusterServiceCreateCluster(ctx)
	if body != nil {
		r = r.Cluster(*body)
	}
	c, h, err := r.Execute()
	return c, parseError(err, h)
}

func (d *ClientDelegate) DeleteCluster(ctx context.Context, clusterId string) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	c, h, err := d.sc.ClusterServiceAPI.ClusterServiceDeleteCluster(ctx, clusterId).Execute()
	return c, parseError(err, h)
}

func (d *ClientDelegate) GetCluster(ctx context.Context, clusterId string, view cluster.ClusterServiceGetClusterViewParameter) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	r := d.sc.ClusterServiceAPI.ClusterServiceGetCluster(ctx, clusterId)
	r = r.View(view)
	c, h, err := r.Execute()
	return c, parseError(err, h)
}

func (d *ClientDelegate) ListProviderRegions(ctx context.Context) (*cluster.TidbCloudOpenApiserverlessv1beta1ListRegionsResponse, error) {
	resp, h, err := d.sc.ClusterServiceAPI.ClusterServiceListRegions(ctx).Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) ListClusters(ctx context.Context, filter *string, pageSize *int32, pageToken *string, orderBy *string, skip *int32) (*cluster.TidbCloudOpenApiserverlessv1beta1ListClustersResponse, error) {
	r := d.sc.ClusterServiceAPI.ClusterServiceListClusters(ctx)
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

func (d *ClientDelegate) PartialUpdateCluster(ctx context.Context, clusterId string, body *cluster.V1beta1ClusterServicePartialUpdateClusterBody) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	r := d.sc.ClusterServiceAPI.ClusterServicePartialUpdateCluster(ctx, clusterId)
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

func (d *ClientDelegate) GetBranch(ctx context.Context, clusterId, branchId string, view branch.BranchServiceGetBranchViewParameter) (*branch.Branch, error) {
	r := d.bc.BranchServiceAPI.BranchServiceGetBranch(ctx, clusterId, branchId)
	r = r.View(view)
	b, h, err := r.Execute()
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

func (d *ClientDelegate) ResetBranch(ctx context.Context, clusterId string, branchId string) (*branch.Branch, error) {
	b, h, err := d.bc.BranchServiceAPI.BranchServiceResetBranch(ctx, clusterId, branchId).Execute()
	return b, parseError(err, h)
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
	res, h, err := d.ec.ExportAPI.ExportServiceGetExport(ctx, clusterId, exportId).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) CancelExport(ctx context.Context, clusterId string, exportId string) (*export.Export, error) {
	res, h, err := d.ec.ExportAPI.ExportServiceCancelExport(ctx, clusterId, exportId).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) CreateExport(ctx context.Context, clusterId string, body *export.ExportServiceCreateExportBody) (*export.Export, error) {
	r := d.ec.ExportAPI.ExportServiceCreateExport(ctx, clusterId)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) DeleteExport(ctx context.Context, clusterId string, exportId string) (*export.Export, error) {
	res, h, err := d.ec.ExportAPI.ExportServiceDeleteExport(ctx, clusterId, exportId).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) ListExports(ctx context.Context, clusterId string, pageSize *int32, pageToken *string, orderBy *string) (*export.ListExportsResponse, error) {
	r := d.ec.ExportAPI.ExportServiceListExports(ctx, clusterId)
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

func (d *ClientDelegate) ListExportFiles(ctx context.Context, clusterId string, exportId string, pageSize *int32,
	pageToken *string, isGenerateUrl bool) (*export.ListExportFilesResponse, error) {
	r := d.ec.ExportServiceAPI.ExportServiceListExportFiles(ctx, clusterId, exportId)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	if isGenerateUrl {
		r = r.GenerateUrl(isGenerateUrl)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) DownloadExportFiles(ctx context.Context, clusterId string, exportId string, body *export.ExportServiceDownloadExportFilesBody) (*export.DownloadExportFilesResponse, error) {
	r := d.ec.ExportServiceAPI.ExportServiceDownloadExportFiles(ctx, clusterId, exportId)
	if body != nil {
		r = r.Body(*body)
	}
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

func (d *ClientDelegate) ListAuditLogs(ctx context.Context, clusterID string, pageSize *int32, pageToken *string, date *string) (*auditlog.ListAuditLogFilesResponse, error) {
	r := d.alc.DatabaseAuditLogServiceAPI.DatabaseAuditLogServiceListAuditLogFiles(ctx, clusterID)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	if date != nil {
		r = r.Date(*date)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) DownloadAuditLogs(ctx context.Context, clusterID string, body *auditlog.DatabaseAuditLogServiceDownloadAuditLogFilesBody) (*auditlog.DownloadAuditLogFilesResponse, error) {
	r := d.alc.DatabaseAuditLogServiceAPI.DatabaseAuditLogServiceDownloadAuditLogFiles(ctx, clusterID)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) CreateAuditLogFilterRule(ctx context.Context, clusterID string, body *auditlog.DatabaseAuditLogServiceCreateAuditLogFilterRuleBody) (*auditlog.AuditLogFilterRule, error) {
	r := d.alc.DatabaseAuditLogServiceAPI.DatabaseAuditLogServiceCreateAuditLogFilterRule(ctx, clusterID)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) DeleteAuditLogFilterRule(ctx context.Context, clusterID, ruleID string) (*auditlog.AuditLogFilterRule, error) {
	res, h, err := d.alc.DatabaseAuditLogServiceAPI.DatabaseAuditLogServiceDeleteAuditLogFilterRule(ctx, clusterID, ruleID).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) GetAuditLogFilterRule(ctx context.Context, clusterID, ruleID string) (*auditlog.AuditLogFilterRule, error) {
	res, h, err := d.alc.DatabaseAuditLogServiceAPI.DatabaseAuditLogServiceGetAuditLogFilterRule(ctx, clusterID, ruleID).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) ListAuditLogFilterRules(ctx context.Context, clusterID string) (*auditlog.ListAuditLogFilterRulesResponse, error) {
	res, h, err := d.alc.DatabaseAuditLogServiceAPI.DatabaseAuditLogServiceListAuditLogFilterRules(ctx, clusterID).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) UpdateAuditLogFilterRule(ctx context.Context, clusterID, ruleID string, body *auditlog.DatabaseAuditLogServiceUpdateAuditLogFilterRuleBody) (*auditlog.AuditLogFilterRule, error) {
	r := d.alc.DatabaseAuditLogServiceAPI.DatabaseAuditLogServiceUpdateAuditLogFilterRule(ctx, clusterID, ruleID)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) UpdateAuditLogConfig(ctx context.Context, clusterID string, body *auditlog.DatabaseAuditLogServiceUpdateAuditLogConfigBody) (*auditlog.AuditLogConfig, error) {
	r := d.alc.DatabaseAuditLogServiceAPI.DatabaseAuditLogServiceUpdateAuditLogConfig(ctx, clusterID)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) GetAuditLogConfig(ctx context.Context, clusterID string) (*auditlog.AuditLogConfig, error) {
	res, h, err := d.alc.DatabaseAuditLogServiceAPI.DatabaseAuditLogServiceGetAuditLogConfig(ctx, clusterID).Execute()
	return res, parseError(err, h)
}

// ===== Private Link Connection =====
func (d *ClientDelegate) CreatePrivateLinkConnection(ctx context.Context, clusterId string, body *privatelink.PrivateLinkConnectionServiceCreatePrivateLinkConnectionBody) (*privatelink.PrivateLinkConnection, error) {
	r := d.plc.PrivateLinkConnectionServiceAPI.PrivateLinkConnectionServiceCreatePrivateLinkConnection(ctx, clusterId)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) DeletePrivateLinkConnection(ctx context.Context, clusterId string, privateLinkConnectionId string) (*privatelink.PrivateLinkConnection, error) {
	res, h, err := d.plc.PrivateLinkConnectionServiceAPI.PrivateLinkConnectionServiceDeletePrivateLinkConnection(ctx, clusterId, privateLinkConnectionId).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) GetPrivateLinkConnection(ctx context.Context, clusterId string, privateLinkConnectionId string) (*privatelink.PrivateLinkConnection, error) {
	res, h, err := d.plc.PrivateLinkConnectionServiceAPI.PrivateLinkConnectionServiceGetPrivateLinkConnection(ctx, clusterId, privateLinkConnectionId).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) ListPrivateLinkConnections(ctx context.Context, clusterId string, pageSize *int32, pageToken *string) (*privatelink.ListPrivateLinkConnectionsResponse, error) {
	r := d.plc.PrivateLinkConnectionServiceAPI.PrivateLinkConnectionServiceListPrivateLinkConnections(ctx, clusterId)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) GetPrivateLinkAvailabilityZones(ctx context.Context, clusterId string) (*privatelink.GetAvailabilityZonesResponse, error) {
	res, h, err := d.plc.PrivateLinkConnectionServiceAPI.PrivateLinkConnectionServiceGetAvailabilityZones(ctx, clusterId).Execute()
	return res, parseError(err, h)
}

// ===== Private Link Connection =====

func NewApiClient(rt http.RoundTripper, serverlessEndpoint string, iamEndpoint string) (*branch.APIClient, *cluster.APIClient, *br.APIClient, *imp.APIClient, *export.APIClient, *iam.APIClient, *auditlog.APIClient, *cdc.APIClient, *privatelink.APIClient, error) {
	httpclient := &http.Client{
		Transport: rt,
	}

	// v1beta1 api (serverless)
	serverlessURL, err := prop.ValidateApiUrl(serverlessEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	iamURL, err := prop.ValidateApiUrl(iamEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	userAgent := fmt.Sprintf("%s/%s", config.CliName, version.Version)

	iamCfg := iam.NewConfiguration()
	iamCfg.HTTPClient = httpclient
	iamCfg.Host = iamURL.Host
	iamCfg.UserAgent = userAgent

	clusterCfg := cluster.NewConfiguration()
	clusterCfg.HTTPClient = httpclient
	clusterCfg.Host = serverlessURL.Host
	clusterCfg.UserAgent = userAgent

	branchCfg := branch.NewConfiguration()
	branchCfg.HTTPClient = httpclient
	branchCfg.Host = serverlessURL.Host
	branchCfg.UserAgent = userAgent

	exportCfg := export.NewConfiguration()
	exportCfg.HTTPClient = httpclient
	exportCfg.Host = serverlessURL.Host
	exportCfg.UserAgent = userAgent

	importCfg := imp.NewConfiguration()
	importCfg.HTTPClient = httpclient
	importCfg.Host = serverlessURL.Host
	importCfg.UserAgent = userAgent

	backupRestoreCfg := br.NewConfiguration()
	backupRestoreCfg.HTTPClient = httpclient
	backupRestoreCfg.Host = serverlessURL.Host
	backupRestoreCfg.UserAgent = userAgent

	auditLogCfg := auditlog.NewConfiguration()
	auditLogCfg.HTTPClient = httpclient
	auditLogCfg.Host = serverlessURL.Host
	auditLogCfg.UserAgent = userAgent

	cdcCfg := cdc.NewConfiguration()
	cdcCfg.HTTPClient = httpclient
	cdcCfg.Host = serverlessURL.Host
	cdcCfg.UserAgent = userAgent

	privateLinkCfg := privatelink.NewConfiguration()
	privateLinkCfg.HTTPClient = httpclient
	privateLinkCfg.Host = serverlessURL.Host
	privateLinkCfg.UserAgent = userAgent

	return branch.NewAPIClient(branchCfg), cluster.NewAPIClient(clusterCfg),
		br.NewAPIClient(backupRestoreCfg),
		imp.NewAPIClient(importCfg), export.NewAPIClient(exportCfg),
		iam.NewAPIClient(iamCfg), auditlog.NewAPIClient(auditLogCfg), cdc.NewAPIClient(cdcCfg),
		privatelink.NewAPIClient(privateLinkCfg), nil

}

func NewDigestTransport(publicKey, privateKey string) http.RoundTripper {
	return &digest.Transport{
		Username:  publicKey,
		Password:  privateKey,
		Transport: NewDebugTransport(http.DefaultTransport),
	}
}

func NewBearTokenTransport(token string) http.RoundTripper {
	return NewTransportWithBearToken(NewDebugTransport(http.DefaultTransport), token)
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
	path := "<path>"
	if resp.Request != nil {
		path = fmt.Sprintf("[%s %s]", resp.Request.Method, resp.Request.URL.Path)
	}
	traceId := "<trace_id>"
	if resp.Header.Get("X-Debug-Trace-Id") != "" {
		traceId = resp.Header.Get("X-Debug-Trace-Id")
	}
	return fmt.Errorf("%s[%s][%s] %s", path, err.Error(), traceId, body)
}

func (d *ClientDelegate) CreateChangefeed(ctx context.Context, clusterId string, body *cdc.ChangefeedServiceCreateChangefeedBody) (*cdc.Changefeed, error) {
	r := d.cdc.ChangefeedServiceAPI.ChangefeedServiceCreateChangefeed(ctx, clusterId)
	if body != nil {
		r = r.Body(*body)
	}
	resp, h, err := r.Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) DeleteChangefeed(ctx context.Context, clusterId, changefeedId string) (*cdc.Changefeed, error) {
	resp, h, err := d.cdc.ChangefeedServiceAPI.ChangefeedServiceDeleteChangefeed(ctx, clusterId, changefeedId).Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) EditChangefeed(ctx context.Context, clusterId, changefeedId string, body *cdc.ChangefeedServiceEditChangefeedBody) (*cdc.Changefeed, error) {
	r := d.cdc.ChangefeedServiceAPI.ChangefeedServiceEditChangefeed(ctx, clusterId, changefeedId)
	if body != nil {
		r = r.Body(*body)
	}
	resp, h, err := r.Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) GetChangefeed(ctx context.Context, clusterId, changefeedId string) (*cdc.Changefeed, error) {
	resp, h, err := d.cdc.ChangefeedServiceAPI.ChangefeedServiceGetChangefeed(ctx, clusterId, changefeedId).Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) ListChangefeeds(ctx context.Context, clusterId string, pageSize *int32, pageToken *string, changefeedType *cdc.ChangefeedServiceListChangefeedsChangefeedTypeParameter) (*cdc.Changefeeds, error) {
	r := d.cdc.ChangefeedServiceAPI.ChangefeedServiceListChangefeeds(ctx, clusterId)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	if changefeedType != nil {
		r = r.ChangefeedType(*changefeedType)
	}
	resp, h, err := r.Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) StartChangefeed(ctx context.Context, clusterId, changefeedId string) (*cdc.Changefeed, error) {
	r := d.cdc.ChangefeedServiceAPI.ChangefeedServiceStartChangefeed(ctx, clusterId, changefeedId)
	r = r.Body(map[string]interface{}{})
	resp, h, err := r.Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) StopChangefeed(ctx context.Context, clusterId, changefeedId string) (*cdc.Changefeed, error) {
	r := d.cdc.ChangefeedServiceAPI.ChangefeedServiceStopChangefeed(ctx, clusterId, changefeedId)
	r = r.Body(map[string]interface{}{})
	resp, h, err := r.Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) TestChangefeed(ctx context.Context, clusterId string, body *cdc.ChangefeedServiceTestChangefeedBody) (map[string]interface{}, error) {
	r := d.cdc.ChangefeedServiceAPI.ChangefeedServiceTestChangefeed(ctx, clusterId)
	if body != nil {
		r = r.Body(*body)
	}
	resp, h, err := r.Execute()
	return resp, parseError(err, h)
}

func (d *ClientDelegate) DescribeSchemaTable(ctx context.Context, clusterId string, body *cdc.ChangefeedServiceDescribeSchemaTableBody) (*cdc.DescribeSchemaTableResp, error) {
	r := d.cdc.ChangefeedServiceAPI.ChangefeedServiceDescribeSchemaTable(ctx, clusterId)
	if body != nil {
		r = r.Body(*body)
	}
	resp, h, err := r.Execute()
	return resp, parseError(err, h)
}
