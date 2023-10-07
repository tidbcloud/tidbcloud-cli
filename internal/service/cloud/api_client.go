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
	branchClient "tidbcloud-cli/pkg/tidbcloud/branch/client"
	branchOp "tidbcloud-cli/pkg/tidbcloud/branch/client/branch_service"
	connectInfoClient "tidbcloud-cli/pkg/tidbcloud/connect_info/client"
	connectInfoOp "tidbcloud-cli/pkg/tidbcloud/connect_info/client/connect_info_service"
	importClient "tidbcloud-cli/pkg/tidbcloud/import/client"
	importOp "tidbcloud-cli/pkg/tidbcloud/import/client/import_service"
	serverlessClient "tidbcloud-cli/pkg/tidbcloud/serverless/client"
	serverlessOp "tidbcloud-cli/pkg/tidbcloud/serverless/client/serverless_service"

	apiClient "github.com/c4pt0r/go-tidbcloud-sdk-v1/client"
	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"
	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/icholy/digest"
)

const (
	DefaultApiUrl    = "https://api.tidbcloud.com"
	DefaultNewApiUrl = "https://serverless.tidbapi.com"
	userAgent        = "User-Agent"
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

	GetBranch(params *branchOp.GetBranchParams, opts ...branchOp.ClientOption) (*branchOp.GetBranchOK, error)

	ListBranches(params *branchOp.ListBranchesParams, opts ...branchOp.ClientOption) (*branchOp.ListBranchesOK, error)

	CreateBranch(params *branchOp.CreateBranchParams, opts ...branchOp.ClientOption) (*branchOp.CreateBranchOK, error)

	DeleteBranch(params *branchOp.DeleteBranchParams, opts ...branchOp.ClientOption) (*branchOp.DeleteBranchOK, error)
}

type ClientDelegate struct {
	c  *apiClient.GoTidbcloud
	ic *importClient.TidbcloudImport
	cc *connectInfoClient.TidbcloudConnectInfo
	bc *branchClient.TidbcloudBranch
	sc *serverlessClient.TidbcloudServerless
}

func NewClientDelegate(publicKey string, privateKey string, apiUrl string, newApiUrl string) (*ClientDelegate, error) {
	c, ic, cc, bc, sc, err := NewApiClient(publicKey, privateKey, apiUrl, newApiUrl)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		c:  c,
		ic: ic,
		cc: cc,
		bc: bc,
		sc: sc,
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

func (d *ClientDelegate) GetBranch(params *branchOp.GetBranchParams, opts ...branchOp.ClientOption) (*branchOp.GetBranchOK, error) {
	r, err := d.bc.BranchService.GetBranch(params, opts...)
	err = parseBranchError(err)
	return r, err
}

func (d *ClientDelegate) ListBranches(params *branchOp.ListBranchesParams, opts ...branchOp.ClientOption) (*branchOp.ListBranchesOK, error) {
	r, err := d.bc.BranchService.ListBranches(params, opts...)
	err = parseBranchError(err)
	return r, err
}

func (d *ClientDelegate) CreateBranch(params *branchOp.CreateBranchParams, opts ...branchOp.ClientOption) (*branchOp.CreateBranchOK, error) {
	r, err := d.bc.BranchService.CreateBranch(params, opts...)
	err = parseBranchError(err)
	return r, err
}

func (d *ClientDelegate) DeleteBranch(params *branchOp.DeleteBranchParams, opts ...branchOp.ClientOption) (*branchOp.DeleteBranchOK, error) {
	r, err := d.bc.BranchService.DeleteBranch(params, opts...)
	err = parseBranchError(err)
	return r, err
}

func NewApiClient(publicKey string, privateKey string, apiUrl string, newApiUrl string) (*apiClient.GoTidbcloud, *importClient.TidbcloudImport,
	*connectInfoClient.TidbcloudConnectInfo, *branchClient.TidbcloudBranch, *serverlessClient.TidbcloudServerless, error) {
	httpclient := &http.Client{
		Transport: NewTransportWithAgent(&digest.Transport{
			Username: publicKey,
			Password: privateKey,
		}, fmt.Sprintf("%s/%s", config.CliName, version.Version)),
	}

	// Parse the URL
	u, err := prop.ValidateApiUrl(apiUrl)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	transport := httpTransport.NewWithClient(u.Host, u.Path, []string{u.Scheme}, httpclient)

	// new api
	newU, err := prop.ValidateApiUrl(newApiUrl)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	newTransport := httpTransport.NewWithClient(newU.Host, serverlessClient.DefaultBasePath, []string{newU.Scheme}, httpclient)

	return apiClient.New(transport, strfmt.Default), importClient.New(transport, strfmt.Default), connectInfoClient.New(transport, strfmt.Default),
		branchClient.New(transport, strfmt.Default), serverlessClient.New(newTransport, strfmt.Default), nil
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

func parseBranchError(err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *branchOp.DeleteBranchDefault:
		msg := "[DELETE /api/v1beta/clusters/{cluster_id}/branches/{branch_id}] DeleteBranch"
		// return by api gateway
		if e.Payload == nil || e.Payload.Error == nil {
			return fmt.Errorf(msg+" [%d] unknown error", e.Code())
		}
		// return by serverless-svc
		return fmt.Errorf(msg+" [%d] %+v", e.Payload.Error.Code, e.Payload.Error.Message)
	case *branchOp.GetBranchDefault:
		msg := "[GET /api/v1beta/clusters/{cluster_id}/branches/{branch_id}] GetBranch"
		// return by api gateway
		if e.Payload == nil || e.Payload.Error == nil {
			return fmt.Errorf(msg+" [%d] unknown error", e.Code())
		}
		// return by serverless-svc
		return fmt.Errorf(msg+" [%d] %+v", e.Payload.Error.Code, e.Payload.Error.Message)
	case *branchOp.ListBranchesDefault:
		msg := "[GET /api/v1beta/clusters/{cluster_id}/branches] ListBranches"
		// return by api gateway
		if e.Payload == nil || e.Payload.Error == nil {
			return fmt.Errorf(msg+" [%d] unknown error", e.Code())
		}
		// return by serverless-svc
		return fmt.Errorf(msg+" [%d] %+v", e.Payload.Error.Code, e.Payload.Error.Message)
	case *branchOp.CreateBranchDefault:
		msg := "[POST /api/v1beta/clusters/{cluster_id}/branches] CreateBranch"
		// return by api gateway
		if e.Payload == nil || e.Payload.Error == nil {
			return fmt.Errorf(msg+" [%d] unknown error", e.Code())
		}
		// return by serverless-svc
		return fmt.Errorf(msg+" [%d] %+v", e.Payload.Error.Code, e.Payload.Error.Message)
	default:
		return err
	}
}
