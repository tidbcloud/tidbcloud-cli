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

	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/prop"
	"tidbcloud-cli/internal/version"
	connectInfoClient "tidbcloud-cli/pkg/tidbcloud/connect_info/client"
	connectInfoOp "tidbcloud-cli/pkg/tidbcloud/connect_info/client/connect_info_service"

	apiClient "github.com/c4pt0r/go-tidbcloud-sdk-v1/client"
	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/import_operations"
	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"
	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/icholy/digest"
)

const (
	DefaultApiUrl = "https://api.tidbcloud.com"
	userAgent     = "User-Agent"
)

type TiDBCloudClient interface {
	CreateCluster(params *cluster.CreateClusterParams, opts ...cluster.ClientOption) (*cluster.CreateClusterOK, error)

	DeleteCluster(params *cluster.DeleteClusterParams, opts ...cluster.ClientOption) (*cluster.DeleteClusterOK, error)

	GetCluster(params *cluster.GetClusterParams, opts ...cluster.ClientOption) (*cluster.GetClusterOK, error)

	ListClustersOfProject(params *cluster.ListClustersOfProjectParams, opts ...cluster.ClientOption) (*cluster.ListClustersOfProjectOK, error)

	ListProviderRegions(params *cluster.ListProviderRegionsParams, opts ...cluster.ClientOption) (*cluster.ListProviderRegionsOK, error)

	ListProjects(params *project.ListProjectsParams, opts ...project.ClientOption) (*project.ListProjectsOK, error)

	UpdateImportTask(params *import_operations.UpdateImportTaskParams, opts ...import_operations.ClientOption) (*import_operations.UpdateImportTaskOK, error)

	CreateImportTask(params *import_operations.CreateImportTaskParams, opts ...import_operations.ClientOption) (*import_operations.CreateImportTaskOK, error)

	GetImportTask(params *import_operations.GetImportTaskParams, opts ...import_operations.ClientOption) (*import_operations.GetImportTaskOK, error)

	ListImportTasks(params *import_operations.ListImportTasksParams, opts ...import_operations.ClientOption) (*import_operations.ListImportTasksOK, error)

	UploadLocalFile(params *import_operations.UploadLocalFileParams, opts ...import_operations.ClientOption) (*import_operations.UploadLocalFileOK, error)

	PreviewImportData(params *import_operations.PreviewImportDataParams, opts ...import_operations.ClientOption) (*import_operations.PreviewImportDataOK, error)

	GetConnectInfo(params *connectInfoOp.GetInfoParams, opts ...connectInfoOp.ClientOption) (*connectInfoOp.GetInfoOK, error)
}

type ClientDelegate struct {
	c  *apiClient.GoTidbcloud
	cc *connectInfoClient.TidbcloudConnectInfo
}

func NewClientDelegate(publicKey string, privateKey string, apiUrl string) (*ClientDelegate, error) {
	c, cc, err := NewApiClient(publicKey, privateKey, apiUrl)
	if err != nil {
		return nil, err
	}
	return &ClientDelegate{
		c:  c,
		cc: cc,
	}, nil
}

func (d *ClientDelegate) CreateCluster(params *cluster.CreateClusterParams, opts ...cluster.ClientOption) (*cluster.CreateClusterOK, error) {
	return d.c.Cluster.CreateCluster(params, opts...)
}

func (d *ClientDelegate) DeleteCluster(params *cluster.DeleteClusterParams, opts ...cluster.ClientOption) (*cluster.DeleteClusterOK, error) {
	return d.c.Cluster.DeleteCluster(params, opts...)
}

func (d *ClientDelegate) GetCluster(params *cluster.GetClusterParams, opts ...cluster.ClientOption) (*cluster.GetClusterOK, error) {
	return d.c.Cluster.GetCluster(params, opts...)
}

func (d *ClientDelegate) ListProviderRegions(params *cluster.ListProviderRegionsParams, opts ...cluster.ClientOption) (*cluster.ListProviderRegionsOK, error) {
	return d.c.Cluster.ListProviderRegions(params, opts...)
}

func (d *ClientDelegate) ListClustersOfProject(params *cluster.ListClustersOfProjectParams, opts ...cluster.ClientOption) (*cluster.ListClustersOfProjectOK, error) {
	return d.c.Cluster.ListClustersOfProject(params, opts...)
}

func (d *ClientDelegate) ListProjects(params *project.ListProjectsParams, opts ...project.ClientOption) (*project.ListProjectsOK, error) {
	return d.c.Project.ListProjects(params, opts...)
}

func (d *ClientDelegate) UpdateImportTask(params *import_operations.UpdateImportTaskParams, opts ...import_operations.ClientOption) (*import_operations.UpdateImportTaskOK, error) {
	return d.c.ImportOperations.UpdateImportTask(params, opts...)
}

func (d *ClientDelegate) CreateImportTask(params *import_operations.CreateImportTaskParams, opts ...import_operations.ClientOption) (*import_operations.CreateImportTaskOK, error) {
	return d.c.ImportOperations.CreateImportTask(params, opts...)
}

func (d *ClientDelegate) GetImportTask(params *import_operations.GetImportTaskParams, opts ...import_operations.ClientOption) (*import_operations.GetImportTaskOK, error) {
	return d.c.ImportOperations.GetImportTask(params, opts...)
}

func (d *ClientDelegate) ListImportTasks(params *import_operations.ListImportTasksParams, opts ...import_operations.ClientOption) (*import_operations.ListImportTasksOK, error) {
	return d.c.ImportOperations.ListImportTasks(params, opts...)
}

func (d *ClientDelegate) UploadLocalFile(params *import_operations.UploadLocalFileParams, opts ...import_operations.ClientOption) (*import_operations.UploadLocalFileOK, error) {
	return d.c.ImportOperations.UploadLocalFile(params, opts...)
}

func (d *ClientDelegate) PreviewImportData(params *import_operations.PreviewImportDataParams, opts ...import_operations.ClientOption) (*import_operations.PreviewImportDataOK, error) {
	return d.c.ImportOperations.PreviewImportData(params, opts...)
}

func (d *ClientDelegate) GetConnectInfo(params *connectInfoOp.GetInfoParams, opts ...connectInfoOp.ClientOption) (*connectInfoOp.GetInfoOK, error) {
	return d.cc.ConnectInfoService.GetInfo(params, opts...)
}

func NewApiClient(publicKey string, privateKey string, apiUrl string) (*apiClient.GoTidbcloud, *connectInfoClient.TidbcloudConnectInfo, error) {
	httpclient := &http.Client{
		Transport: NewTransportWithAgent(&digest.Transport{
			Username: publicKey,
			Password: privateKey,
		}, fmt.Sprintf("%s/%s", config.CliName, version.Version)),
	}

	// Parse the URL
	u, err := prop.ValidateApiUrl(apiUrl)
	if err != nil {
		return nil, nil, err
	}

	transport := httpTransport.NewWithClient(u.Host, u.Path, []string{u.Scheme}, httpclient)
	return apiClient.New(transport, strfmt.Default), connectInfoClient.New(transport, strfmt.Default), nil
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
