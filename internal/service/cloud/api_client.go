// Copyright 2022 PingCAP, Inc.
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
	"fmt"
	"net/http"
	"os"

	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/version"
	connectInfoClient "tidbcloud-cli/pkg/tidbcloud/connect_info/client"
	connectInfoOp "tidbcloud-cli/pkg/tidbcloud/connect_info/client/connect_info_service"
	importClient "tidbcloud-cli/pkg/tidbcloud/import/client"
	importOp "tidbcloud-cli/pkg/tidbcloud/import/client/import_service"
	pingchatClient "tidbcloud-cli/pkg/tidbcloud/pingchat/client"
	pingchatOp "tidbcloud-cli/pkg/tidbcloud/pingchat/client/operations"
	branchClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client"
	branchOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client/branch_service"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"
	serverlessClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client"
	serverlessOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
	brClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client"
	brOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"
	expClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/client"
	expOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/client/export_service"

	apiClient "github.com/c4pt0r/go-tidbcloud-sdk-v1/client"
	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"
	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-resty/resty/v2"
	"github.com/icholy/digest"
)

const (
	DefaultApiUrl             = "https://" + apiClient.DefaultHost
	DefaultServerlessEndpoint = "https://" + serverlessClient.DefaultHost
	userAgent                 = "User-Agent"
)

type TiDBCloudClient interface {
	CreateCluster(params *serverlessOp.ServerlessServiceCreateClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceCreateClusterOK, error)

	DeleteCluster(params *serverlessOp.ServerlessServiceDeleteClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceDeleteClusterOK, error)

	GetCluster(params *serverlessOp.ServerlessServiceGetClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceGetClusterOK, error)

	ListClustersOfProject(params *serverlessOp.ServerlessServiceListClustersParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceListClustersOK, error)

	PartialUpdateCluster(params *serverlessOp.ServerlessServicePartialUpdateClusterParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServicePartialUpdateClusterOK, error)

	ListProviderRegions(params *serverlessOp.ServerlessServiceListRegionsParams, opts ...serverlessOp.ClientOption) (*serverlessOp.ServerlessServiceListRegionsOK, error)

	ListProjects(params *project.ListProjectsParams, opts ...project.ClientOption) (*project.ListProjectsOK, error)

	CancelImport(params *importOp.CancelImportParams, opts ...importOp.ClientOption) (*importOp.CancelImportOK, error)

	CreateImport(params *importOp.CreateImportParams, opts ...importOp.ClientOption) (*importOp.CreateImportOK, error)

	GetImport(params *importOp.GetImportParams, opts ...importOp.ClientOption) (*importOp.GetImportOK, error)

	ListImports(params *importOp.ListImportsParams, opts ...importOp.ClientOption) (*importOp.ListImportsOK, error)

	GenerateUploadURL(params *importOp.GenerateUploadURLParams, opts ...importOp.ClientOption) (*importOp.GenerateUploadURLOK, error)

	PreSignedUrlUpload(url *string, uploadFile *os.File, size int64) error

	GetConnectInfo(params *connectInfoOp.GetInfoParams, opts ...connectInfoOp.ClientOption) (*connectInfoOp.GetInfoOK, error)

	GetBranch(params *branchOp.BranchServiceGetBranchParams, opts ...branchOp.ClientOption) (*branchOp.BranchServiceGetBranchOK, error)

	ListBranches(params *branchOp.BranchServiceListBranchesParams, opts ...branchOp.ClientOption) (*branchOp.BranchServiceListBranchesOK, error)

	CreateBranch(params *branchOp.BranchServiceCreateBranchParams, opts ...branchOp.ClientOption) (*branchOp.BranchServiceCreateBranchOK, error)

	DeleteBranch(params *branchOp.BranchServiceDeleteBranchParams, opts ...branchOp.ClientOption) (*branchOp.BranchServiceDeleteBranchOK, error)

	Chat(params *pingchatOp.ChatParams, opts ...pingchatOp.ClientOption) (*pingchatOp.ChatOK, error)

	DeleteBackup(params *brOp.BackupRestoreServiceDeleteBackupParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceDeleteBackupOK, error)

	GetBackup(params *brOp.BackupRestoreServiceGetBackupParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceGetBackupOK, error)

	ListBackups(params *brOp.BackupRestoreServiceListBackupsParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceListBackupsOK, error)

	Restore(params *brOp.BackupRestoreServiceRestoreParams, opts ...brOp.ClientOption) (*brOp.BackupRestoreServiceRestoreOK, error)

	GetExport(params *expOp.ExportServiceGetExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceGetExportOK, error)

	CancelExport(params *expOp.ExportServiceCancelExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceCancelExportOK, error)

	CreateExport(params *expOp.ExportServiceCreateExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceCreateExportOK, error)

	DeleteExport(params *expOp.ExportServiceDeleteExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceDeleteExportOK, error)

	ListExports(params *expOp.ExportServiceListExportsParams, opts ...expOp.ClientOption) (*expOp.ExportServiceListExportsOK, error)

	DownloadExport(params *expOp.ExportServiceDownloadExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceDownloadExportOK, error)
}

type ClientDelegate struct {
	c   *apiClient.GoTidbcloud
	is  *iam.Service
	ic  *importClient.TidbcloudImport
	cc  *connectInfoClient.TidbcloudConnectInfo
	bc  *branchClient.TidbcloudServerless
	pc  *pingchatClient.TidbcloudPingchat
	sc  *serverlessClient.TidbcloudServerless
	brc *brClient.TidbcloudServerless
	ec  *expClient.TidbcloudServerless
}

func NewClientDelegateWithToken(token string, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewBearTokenTransport(token)
	client := resty.New()
	client.SetAuthToken(token)
	debug := os.Getenv(config.DebugEnv) != ""
	client.SetDebug(debug)
	c, ic, cc, bc, sc, pc, brc, ec,is, err := NewApiClient(transport, apiUrl, serverlessEndpoint, iamEndpoint, client)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		c:   c,
		ic:  ic,
		cc:  cc,
		bc:  bc,
		sc:  sc,
		pc:  pc,
		brc: brc,
		ec:  ec,
		is:  is,
	}, nil
}

func NewClientDelegateWithApiKey(publicKey string, privateKey string, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewDigestTransport(publicKey, privateKey)
	client := resty.New()
	client.SetDigestAuth(publicKey, privateKey)
	c, ic, cc, bc, sc, pc, brc,ec, is, err := NewApiClient(transport, apiUrl, serverlessEndpoint, iamEndpoint, client)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		c:   c,
		ic:  ic,
		cc:  cc,
		bc:  bc,
		sc:  sc,
		pc:  pc,
		brc: brc,
		ec:  ec,
		is:  is,
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

func (d *ClientDelegate) ListProjects(params *project.ListProjectsParams, opts ...project.ClientOption) (*project.ListProjectsOK, error) {
	return d.c.Project.ListProjects(params, opts...)
}

func (d *ClientDelegate) CancelImport(params *importOp.CancelImportParams, opts ...importOp.ClientOption) (*importOp.CancelImportOK, error) {
	return d.ic.ImportService.CancelImport(params, opts...)
}

func (d *ClientDelegate) CreateImport(params *importOp.CreateImportParams, opts ...importOp.ClientOption) (*importOp.CreateImportOK, error) {
	return d.ic.ImportService.CreateImport(params, opts...)
}

func (d *ClientDelegate) GetImport(params *importOp.GetImportParams, opts ...importOp.ClientOption) (*importOp.GetImportOK, error) {
	return d.ic.ImportService.GetImport(params, opts...)
}

func (d *ClientDelegate) ListImports(params *importOp.ListImportsParams, opts ...importOp.ClientOption) (*importOp.ListImportsOK, error) {
	return d.ic.ImportService.ListImports(params, opts...)
}

func (d *ClientDelegate) GenerateUploadURL(params *importOp.GenerateUploadURLParams, opts ...importOp.ClientOption) (*importOp.GenerateUploadURLOK, error) {
	return d.ic.ImportService.GenerateUploadURL(params, opts...)
}

func (d *ClientDelegate) PreSignedUrlUpload(url *string, uploadFile *os.File, size int64) error {
	request, err := http.NewRequest("PUT", *url, uploadFile)
	if err != nil {
		return err
	}
	request.ContentLength = size

	putRes, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer putRes.Body.Close()

	if putRes.StatusCode != http.StatusOK {
		return fmt.Errorf("upload file failed : %s, %s", putRes.Status, putRes.Body)
	}

	return nil
}

func (d *ClientDelegate) GetConnectInfo(params *connectInfoOp.GetInfoParams, opts ...connectInfoOp.ClientOption) (*connectInfoOp.GetInfoOK, error) {
	return d.cc.ConnectInfoService.GetInfo(params, opts...)
}

func (d *ClientDelegate) GetBranch(params *branchOp.BranchServiceGetBranchParams, opts ...branchOp.ClientOption) (*branchOp.BranchServiceGetBranchOK, error) {
	r, err := d.bc.BranchService.BranchServiceGetBranch(params, opts...)
	return r, err
}

func (d *ClientDelegate) ListBranches(params *branchOp.BranchServiceListBranchesParams, opts ...branchOp.ClientOption) (*branchOp.BranchServiceListBranchesOK, error) {
	r, err := d.bc.BranchService.BranchServiceListBranches(params, opts...)
	return r, err
}

func (d *ClientDelegate) CreateBranch(params *branchOp.BranchServiceCreateBranchParams, opts ...branchOp.ClientOption) (*branchOp.BranchServiceCreateBranchOK, error) {
	r, err := d.bc.BranchService.BranchServiceCreateBranch(params, opts...)
	return r, err
}

func (d *ClientDelegate) DeleteBranch(params *branchOp.BranchServiceDeleteBranchParams, opts ...branchOp.ClientOption) (*branchOp.BranchServiceDeleteBranchOK, error) {
	r, err := d.bc.BranchService.BranchServiceDeleteBranch(params, opts...)
	return r, err
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

func (d *ClientDelegate) GetExport(params *expOp.ExportServiceGetExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceGetExportOK, error) {
	return d.ec.ExportService.ExportServiceGetExport(params, opts...)
}

func (d *ClientDelegate) CancelExport(params *expOp.ExportServiceCancelExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceCancelExportOK, error) {
	return d.ec.ExportService.ExportServiceCancelExport(params, opts...)
}

func (d *ClientDelegate) CreateExport(params *expOp.ExportServiceCreateExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceCreateExportOK, error) {
	return d.ec.ExportService.ExportServiceCreateExport(params, opts...)
}

func (d *ClientDelegate) DeleteExport(params *expOp.ExportServiceDeleteExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceDeleteExportOK, error) {
	return d.ec.ExportService.ExportServiceDeleteExport(params, opts...)
}

func (d *ClientDelegate) ListExports(params *expOp.ExportServiceListExportsParams, opts ...expOp.ClientOption) (*expOp.ExportServiceListExportsOK, error) {
	return d.ec.ExportService.ExportServiceListExports(params, opts...)
}

func (d *ClientDelegate) DownloadExport(params *expOp.ExportServiceDownloadExportParams, opts ...expOp.ClientOption) (*expOp.ExportServiceDownloadExportOK, error) {
	return d.ec.ExportService.ExportServiceDownloadExport(params, opts...)
}

func NewApiClient(rt http.RoundTripper, apiUrl string, serverlessEndpoint string, iamEndpoint string, client *resty.Client) (*apiClient.GoTidbcloud, *importClient.TidbcloudImport, *connectInfoClient.TidbcloudConnectInfo, *branchClient.TidbcloudServerless, *serverlessClient.TidbcloudServerless, *pingchatClient.TidbcloudPingchat, *brClient.TidbcloudServerless, *expClient.TidbcloudServerless, *iam.Service, error) {
	httpclient := &http.Client{
		Transport: rt,
	}

	// v1beta api
	u, err := prop.ValidateApiUrl(apiUrl)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	transport := httpTransport.NewWithClient(u.Host, u.Path, []string{u.Scheme}, httpclient)

	// v1beta1 api (serverless)
	serverlessURL, err := prop.ValidateApiUrl(serverlessEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	serverlessTransport := httpTransport.NewWithClient(serverlessURL.Host, serverlessClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	branchTransport := httpTransport.NewWithClient(serverlessURL.Host, branchClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	backRestoreTransport := httpTransport.NewWithClient(serverlessURL.Host, branchClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	exportTransport := httpTransport.NewWithClient(serverlessURL.Host, expClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)

	iamUrl, err := prop.ValidateApiUrl(iamEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return apiClient.New(transport, strfmt.Default), importClient.New(transport, strfmt.Default), connectInfoClient.New(transport, strfmt.Default),
		branchClient.New(branchTransport, strfmt.Default), serverlessClient.New(serverlessTransport, strfmt.Default),
		pingchatClient.New(transport, strfmt.Default), brClient.New(backRestoreTransport, strfmt.Default), expClient.New(exportTransport, strfmt.Default), iam.NewIamService(client, iamUrl.String()), nil
}

func NewDigestTransport(publicKey, privateKey string) http.RoundTripper {
	return NewTransportWithAgent(&digest.Transport{
		Username: publicKey,
		Password: privateKey,
	}, fmt.Sprintf("%s/%s", config.CliName, version.Version))
}

func NewBearTokenTransport(token string) http.RoundTripper {
	return NewTransportWithAgent(NewTransportWithBearToken(http.DefaultTransport, token),
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
