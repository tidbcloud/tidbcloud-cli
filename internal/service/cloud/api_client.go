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
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/pingchat"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/auditlog"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/br"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/branch"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/dm"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"

	"github.com/icholy/digest"
)

const (
	DefaultApiUrl             = "https://api.tidbcloud.com"
	DefaultServerlessEndpoint = "https://serverless.tidbapi.com"
	DefaultIAMEndpoint        = "https://iam.tidbapi.com"
)

type TiDBCloudClient interface {
	CreateCluster(ctx context.Context, body *cluster.TidbCloudOpenApiserverlessv1beta1Cluster) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	DeleteCluster(ctx context.Context, clusterId string) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	GetCluster(ctx context.Context, clusterId string, view cluster.ServerlessServiceGetClusterViewParameter) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

	ListClusters(ctx context.Context, filter *string, pageSize *int32, pageToken *string, orderBy *string, skip *int32) (*cluster.TidbCloudOpenApiserverlessv1beta1ListClustersResponse, error)

	PartialUpdateCluster(ctx context.Context, clusterId string, body *cluster.V1beta1ServerlessServicePartialUpdateClusterBody) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error)

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

	Chat(ctx context.Context, chatInfo *pingchat.PingchatChatInfo) (*pingchat.PingchatChatResponse, error)

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

	DownloadAuditLogs(ctx context.Context, clusterID string, body *auditlog.AuditLogServiceDownloadAuditLogsBody) (*auditlog.DownloadAuditLogsResponse, error)

	ListAuditLogs(ctx context.Context, clusterID string, pageSize *int32, pageToken *string, date *string) (*auditlog.ListAuditLogsResponse, error)

	CreateAuditLogFilterRule(ctx context.Context, clusterID string, body *auditlog.AuditLogServiceCreateAuditLogFilterRuleBody) (*auditlog.AuditLogFilterRule, error)

	DeleteAuditLogFilterRule(ctx context.Context, clusterID, name string) (*auditlog.AuditLogFilterRule, error)

	GetAuditLogFilterRule(ctx context.Context, clusterID, name string) (*auditlog.AuditLogFilterRule, error)

	ListAuditLogFilterRules(ctx context.Context, clusterID string) (*auditlog.ListAuditLogFilterRulesResponse, error)

	UpdateAuditLogFilterRule(ctx context.Context, clusterID, name string, body *auditlog.AuditLogServiceUpdateAuditLogFilterRuleBody) (*auditlog.AuditLogFilterRule, error)

	// DM (Data Migration) operations
	Precheck(ctx context.Context, clusterID string, body *dm.DMServicePrecheckBody) (*dm.CreateDMPrecheckResp, error)
	CreateTask(ctx context.Context, clusterID string, body *dm.DMServiceCreateTaskBody) (*dm.DMTask, error)
	GetTask(ctx context.Context, clusterID string, taskID string) (*dm.DMTask, error)
	ListTasks(ctx context.Context, clusterID string, pageSize *int32, pageToken *string, orderBy *string) (*dm.ListDMTasksResp, error)
	DeleteTask(ctx context.Context, clusterID string, taskID string) (*dm.DMTask, error)
	// used to operate task status, i.e. pause/resume task
	OperateTask(ctx context.Context, clusterID string, taskID string, body *dm.DMServiceOperateTaskBody) error
	CancelPrecheck(ctx context.Context, clusterID string, precheckID string) error
	GetPrecheck(ctx context.Context, clusterID string, precheckID string) (*dm.DMPrecheck, error)
}
type ClientDelegate struct {
	ic  *iam.APIClient
	bc  *branch.APIClient
	pc  *pingchat.APIClient
	brc *br.APIClient
	sc  *cluster.APIClient
	sic *imp.APIClient
	ec  *export.APIClient
	alc *auditlog.APIClient
	dmc *dm.APIClient
}

func NewClientDelegateWithToken(token string, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewBearTokenTransport(token)
	bc, sc, pc, brc, sic, ec, ic, alc, dmc, err := NewApiClient(transport, apiUrl, serverlessEndpoint, iamEndpoint)
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
		alc: alc,
		dmc: dmc,
	}, nil
}

func NewClientDelegateWithApiKey(publicKey string, privateKey string, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewDigestTransport(publicKey, privateKey)
	bc, sc, pc, brc, sic, ec, ic, alc, dmc, err := NewApiClient(transport, apiUrl, serverlessEndpoint, iamEndpoint)
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
		alc: alc,
		dmc: dmc,
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

func (d *ClientDelegate) GetCluster(ctx context.Context, clusterId string, view cluster.ServerlessServiceGetClusterViewParameter) (*cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	r := d.sc.ServerlessServiceAPI.ServerlessServiceGetCluster(ctx, clusterId)
	r = r.View(view)
	c, h, err := r.Execute()
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

func (d *ClientDelegate) Chat(ctx context.Context, chatInfo *pingchat.PingchatChatInfo) (*pingchat.PingchatChatResponse, error) {
	r := d.pc.PingChatServiceAPI.Chat(ctx)
	if chatInfo != nil {
		r = r.ChatInfo(*chatInfo)
	}
	resp, h, err := r.Execute()
	return resp, parseError(err, h)
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

func (d *ClientDelegate) ListAuditLogs(ctx context.Context, clusterID string, pageSize *int32, pageToken *string, date *string) (*auditlog.ListAuditLogsResponse, error) {
	r := d.alc.AuditLogServiceAPI.AuditLogServiceListAuditLogs(ctx, clusterID)
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

func (d *ClientDelegate) DownloadAuditLogs(ctx context.Context, clusterID string, body *auditlog.AuditLogServiceDownloadAuditLogsBody) (*auditlog.DownloadAuditLogsResponse, error) {
	r := d.alc.AuditLogServiceAPI.AuditLogServiceDownloadAuditLogs(ctx, clusterID)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) CreateAuditLogFilterRule(ctx context.Context, clusterID string, body *auditlog.AuditLogServiceCreateAuditLogFilterRuleBody) (*auditlog.AuditLogFilterRule, error) {
	r := d.alc.AuditLogServiceAPI.AuditLogServiceCreateAuditLogFilterRule(ctx, clusterID)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) DeleteAuditLogFilterRule(ctx context.Context, clusterID, name string) (*auditlog.AuditLogFilterRule, error) {
	res, h, err := d.alc.AuditLogServiceAPI.AuditLogServiceDeleteAuditLogFilterRule(ctx, clusterID, name).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) GetAuditLogFilterRule(ctx context.Context, clusterID, name string) (*auditlog.AuditLogFilterRule, error) {
	res, h, err := d.alc.AuditLogServiceAPI.AuditLogServiceGetAuditLogFilterRule(ctx, clusterID, name).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) ListAuditLogFilterRules(ctx context.Context, clusterID string) (*auditlog.ListAuditLogFilterRulesResponse, error) {
	res, h, err := d.alc.AuditLogServiceAPI.AuditLogServiceListAuditLogFilterRules(ctx, clusterID).Execute()
	return res, parseError(err, h)
}

func (d *ClientDelegate) UpdateAuditLogFilterRule(ctx context.Context, clusterID, name string, body *auditlog.AuditLogServiceUpdateAuditLogFilterRuleBody) (*auditlog.AuditLogFilterRule, error) {
	r := d.alc.AuditLogServiceAPI.AuditLogServiceUpdateAuditLogFilterRule(ctx, clusterID, name)
	if body != nil {
		r = r.Body(*body)
	}
	res, h, err := r.Execute()
	return res, parseError(err, h)
}

func NewApiClient(rt http.RoundTripper, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*branch.APIClient, *cluster.APIClient, *pingchat.APIClient, *br.APIClient, *imp.APIClient, *export.APIClient, *iam.APIClient, *auditlog.APIClient, *dm.APIClient, error) {
	httpclient := &http.Client{
		Transport: rt,
	}

	// v1beta api
	u, err := prop.ValidateApiUrl(apiUrl)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
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

	pingchatCfg := pingchat.NewConfiguration()
	pingchatCfg.HTTPClient = httpclient
	pingchatCfg.Host = u.Host
	pingchatCfg.UserAgent = userAgent

	auditLogCfg := auditlog.NewConfiguration()
	auditLogCfg.HTTPClient = httpclient
	auditLogCfg.Host = serverlessURL.Host
	auditLogCfg.UserAgent = userAgent

	dmCfg := dm.NewConfiguration()
	dmCfg.HTTPClient = httpclient
	dmCfg.Host = serverlessURL.Host
	dmCfg.UserAgent = userAgent

	return branch.NewAPIClient(branchCfg), cluster.NewAPIClient(clusterCfg),
		pingchat.NewAPIClient(pingchatCfg), br.NewAPIClient(backupRestoreCfg),
		imp.NewAPIClient(importCfg), export.NewAPIClient(exportCfg),
		iam.NewAPIClient(iamCfg), auditlog.NewAPIClient(auditLogCfg),
		dm.NewAPIClient(dmCfg), nil
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

// DM (Data Migration) method implementations

func (d *ClientDelegate) Precheck(ctx context.Context, clusterID string, body *dm.DMServicePrecheckBody) (*dm.CreateDMPrecheckResp, error) {
	r := d.dmc.DMAPI.DMServicePrecheck(ctx, clusterID)
	if body != nil {
		r = r.Body(*body)
	}
	result, h, err := r.Execute()
	return result, parseError(err, h)
}

func (d *ClientDelegate) CreateTask(ctx context.Context, clusterID string, body *dm.DMServiceCreateTaskBody) (*dm.DMTask, error) {
	r := d.dmc.DMServiceAPI.DMServiceCreateTask(ctx, clusterID)
	if body != nil {
		r = r.Body(*body)
	}
	result, h, err := r.Execute()
	return result, parseError(err, h)
}

func (d *ClientDelegate) GetTask(ctx context.Context, clusterID string, taskID string) (*dm.DMTask, error) {
	result, h, err := d.dmc.DMServiceAPI.DMServiceGetTask(ctx, clusterID, taskID).Execute()
	return result, parseError(err, h)
}

func (d *ClientDelegate) ListTasks(ctx context.Context, clusterID string, pageSize *int32, pageToken *string, orderBy *string) (*dm.ListDMTasksResp, error) {
	r := d.dmc.DMServiceAPI.DMServiceListTasks(ctx, clusterID)
	if pageSize != nil {
		r = r.PageSize(*pageSize)
	}
	if pageToken != nil {
		r = r.PageToken(*pageToken)
	}
	if orderBy != nil {
		r = r.OrderBy(*orderBy)
	}
	result, h, err := r.Execute()
	return result, parseError(err, h)
}

func (d *ClientDelegate) DeleteTask(ctx context.Context, clusterID string, taskID string) (*dm.DMTask, error) {
	result, h, err := d.dmc.DMServiceAPI.DMServiceDeleteTask(ctx, clusterID, taskID).Execute()
	return result, parseError(err, h)
}

func (d *ClientDelegate) OperateTask(ctx context.Context, clusterID string, taskID string, body *dm.DMServiceOperateTaskBody) error {
	r := d.dmc.DMServiceAPI.DMServiceOperateTask(ctx, clusterID, taskID)
	if body != nil {
		r = r.Body(*body)
	}
	_, h, err := r.Execute()
	return parseError(err, h)
}

func (d *ClientDelegate) CancelPrecheck(ctx context.Context, clusterID string, precheckID string) error {
	_, h, err := d.dmc.DMServiceAPI.DMServiceCancelPrecheck(ctx, clusterID, precheckID).Execute()
	return parseError(err, h)
}

func (d *ClientDelegate) GetPrecheck(ctx context.Context, clusterID string, precheckID string) (*dm.DMPrecheck, error) {
	result, h, err := d.dmc.DMServiceAPI.DMServiceGetPrecheck(ctx, clusterID, precheckID).Execute()
	return result, parseError(err, h)
}
