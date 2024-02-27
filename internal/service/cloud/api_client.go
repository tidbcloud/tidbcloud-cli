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
	serverlessClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client"
	serverlessOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
	brClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client"
	brOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"
	serverlessImportClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client"
	serverlessImportOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"

	apiClient "github.com/c4pt0r/go-tidbcloud-sdk-v1/client"
	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"
	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
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

	StartUpload(params *serverlessImportOp.ImportServiceStartUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceStartUploadOK, error)

	CompleteMultipartUpload(params *serverlessImportOp.ImportServiceCompleteMultipartUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCompleteMultipartUploadOK, error)

	CancelMultipartUpload(params *serverlessImportOp.ImportServiceCancelMultipartUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCancelMultipartUploadOK, error)
}

type ClientDelegate struct {
	c   *apiClient.GoTidbcloud
	ic  *importClient.TidbcloudImport
	cc  *connectInfoClient.TidbcloudConnectInfo
	bc  *branchClient.TidbcloudServerless
	pc  *pingchatClient.TidbcloudPingchat
	sc  *serverlessClient.TidbcloudServerless
	brc *brClient.TidbcloudServerless
	sic *serverlessImportClient.TidbcloudServerless
}

func NewClientDelegate(publicKey string, privateKey string, apiUrl string, serverlessEndpoint string) (*ClientDelegate, error) {
	c, ic, cc, bc, sc, pc, brc, sic, err := NewApiClient(publicKey, privateKey, apiUrl, serverlessEndpoint)
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

func (d *ClientDelegate) StartUpload(params *serverlessImportOp.ImportServiceStartUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceStartUploadOK, error) {
	return d.sic.ImportService.ImportServiceStartUpload(params, opts...)
}

func (d *ClientDelegate) CompleteMultipartUpload(params *serverlessImportOp.ImportServiceCompleteMultipartUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCompleteMultipartUploadOK, error) {
	return d.sic.ImportService.ImportServiceCompleteMultipartUpload(params, opts...)
}

func (d *ClientDelegate) CancelMultipartUpload(params *serverlessImportOp.ImportServiceCancelMultipartUploadParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCancelMultipartUploadOK, error) {
	return d.sic.ImportService.ImportServiceCancelMultipartUpload(params, opts...)
}

func NewApiClient(publicKey string, privateKey string, apiUrl string, serverlessEndpoint string) (*apiClient.GoTidbcloud, *importClient.TidbcloudImport, *connectInfoClient.TidbcloudConnectInfo, *branchClient.TidbcloudServerless, *serverlessClient.TidbcloudServerless, *pingchatClient.TidbcloudPingchat, *brClient.TidbcloudServerless, *serverlessImportClient.TidbcloudServerless, error) {
	httpclient := &http.Client{
		Transport: NewTransportWithAgent(&digest.Transport{
			Username: publicKey,
			Password: privateKey,
		}, fmt.Sprintf("%s/%s", config.CliName, version.Version)),
	}

	// v1beta api
	u, err := prop.ValidateApiUrl(apiUrl)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	transport := httpTransport.NewWithClient(u.Host, u.Path, []string{u.Scheme}, httpclient)

	// v1beta1 api
	serverlessURL, err := prop.ValidateApiUrl(serverlessEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	serverlessTransport := httpTransport.NewWithClient(serverlessURL.Host, serverlessClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	branchTransport := httpTransport.NewWithClient(serverlessURL.Host, branchClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	backRestoreTransport := httpTransport.NewWithClient(serverlessURL.Host, branchClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	importTransport := httpTransport.NewWithClient(serverlessURL.Host, branchClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)

	return apiClient.New(transport, strfmt.Default), importClient.New(transport, strfmt.Default), connectInfoClient.New(transport, strfmt.Default),
		branchClient.New(branchTransport, strfmt.Default), serverlessClient.New(serverlessTransport, strfmt.Default),
		pingchatClient.New(transport, strfmt.Default), brClient.New(backRestoreTransport, strfmt.Default),
		serverlessImportClient.New(importTransport, strfmt.Default), nil
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
