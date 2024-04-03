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
	pingchatClient "tidbcloud-cli/pkg/tidbcloud/pingchat/client"
	pingchatOp "tidbcloud-cli/pkg/tidbcloud/pingchat/client/operations"
	branchClient "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client"
	branchOp "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client/branch_service"
	"tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"
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

	ListProjects(params *iam.ListProjectsParams) (res *project.ListProjectsOK, err error)

	CancelImport(params *serverlessImportOp.ImportServiceCancelImportParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCancelImportOK, error)

	CreateImport(params *serverlessImportOp.ImportServiceCreateImportParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceCreateImportOK, error)

	GetImport(params *serverlessImportOp.ImportServiceGetImportParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceGetImportOK, error)

	ListImports(params *serverlessImportOp.ImportServiceListImportsParams, opts ...serverlessImportOp.ClientOption) (*serverlessImportOp.ImportServiceListImportsOK, error)

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
	is  *iam.Service
	bc  *branchClient.TidbcloudServerless
	pc  *pingchatClient.TidbcloudPingchat
	sc  *serverlessClient.TidbcloudServerless
	brc *brClient.TidbcloudServerless
	sic *serverlessImportClient.TidbcloudServerless
}

func NewClientDelegateWithToken(token string, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewBearTokenTransport(token)
	client := resty.New()
	client.SetAuthToken(token)
	debug := os.Getenv(config.DebugEnv) != ""
	client.SetDebug(debug)
	bc, sc, pc, brc, sic, is, err := NewApiClient(transport, apiUrl, serverlessEndpoint, iamEndpoint, client)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		bc:  bc,
		sc:  sc,
		pc:  pc,
		brc: brc,
		is:  is,
		sic: sic,
	}, nil
}

func NewClientDelegateWithApiKey(publicKey string, privateKey string, apiUrl string, serverlessEndpoint string, iamEndpoint string) (*ClientDelegate, error) {
	transport := NewDigestTransport(publicKey, privateKey)
	client := resty.New()
	client.SetDigestAuth(publicKey, privateKey)
	bc, sc, pc, brc, sic, is, err := NewApiClient(transport, apiUrl, serverlessEndpoint, iamEndpoint, client)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		bc:  bc,
		sc:  sc,
		pc:  pc,
		brc: brc,
		is:  is,
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

func (d *ClientDelegate) ListProjects(params *iam.ListProjectsParams) (res *project.ListProjectsOK, err error) {
	return d.is.ListProjects(params)
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

func NewApiClient(rt http.RoundTripper, apiUrl string, serverlessEndpoint string, iamEndpoint string, client *resty.Client) (*branchClient.TidbcloudServerless, *serverlessClient.TidbcloudServerless, *pingchatClient.TidbcloudPingchat, *brClient.TidbcloudServerless, *serverlessImportClient.TidbcloudServerless, *iam.Service, error) {
	httpclient := &http.Client{
		Transport: rt,
	}

	// v1beta api
	u, err := prop.ValidateApiUrl(apiUrl)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	transport := httpTransport.NewWithClient(u.Host, u.Path, []string{u.Scheme}, httpclient)

	// v1beta1 api
	serverlessURL, err := prop.ValidateApiUrl(serverlessEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	serverlessTransport := httpTransport.NewWithClient(serverlessURL.Host, serverlessClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	branchTransport := httpTransport.NewWithClient(serverlessURL.Host, branchClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	backRestoreTransport := httpTransport.NewWithClient(serverlessURL.Host, branchClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)
	importTransport := httpTransport.NewWithClient(serverlessURL.Host, branchClient.DefaultBasePath, []string{serverlessURL.Scheme}, httpclient)

	iamUrl, err := prop.ValidateApiUrl(iamEndpoint)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return branchClient.New(branchTransport, strfmt.Default), serverlessClient.New(serverlessTransport, strfmt.Default),
		pingchatClient.New(transport, strfmt.Default), brClient.New(backRestoreTransport, strfmt.Default),
		serverlessImportClient.New(importTransport, strfmt.Default), iam.NewIamService(client, iamUrl.String()), nil
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
