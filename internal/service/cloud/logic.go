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
	"math"
	"os"

	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	connectInfoApi "tidbcloud-cli/pkg/tidbcloud/connect_info/client/connect_info_service"
	connectInfoModel "tidbcloud-cli/pkg/tidbcloud/connect_info/models"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/import_operations"
	projectApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
)

type Project struct {
	ID   string
	Name string
}

func (p Project) String() string {
	return fmt.Sprintf("%s(%s)", p.Name, p.ID)
}

type Cluster struct {
	ID   string
	Name string
}

func (c Cluster) String() string {
	return fmt.Sprintf("%s(%s)", c.Name, c.ID)
}

type Import struct {
	ID    string
	Phase string
}

func (i Import) String() string {
	return fmt.Sprintf("%s(%s)", i.ID, i.Phase)
}

func GetSelectedProject(pageSize int64, client TiDBCloudClient) (*Project, error) {
	_, projectItems, err := RetrieveProjects(pageSize, client)
	if err != nil {
		return nil, err
	}

	// If there is only one project, return it directly.
	if len(projectItems) == 1 {
		return &Project{
			projectItems[0].ID,
			projectItems[0].Name,
		}, nil
	}

	var items = make([]interface{}, 0, len(projectItems))
	for _, item := range projectItems {
		items = append(items, &Project{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	model, err := ui.InitialSelectModel(items, "Choose the project:")
	if err != nil {
		return nil, err
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	projectModel, err := p.StartReturningModel()
	if err != nil {
		return nil, err
	}
	if m, _ := projectModel.(ui.SelectModel); m.Interrupted {
		os.Exit(130)
	}
	res := projectModel.(ui.SelectModel).GetSelectedItem().(*Project)
	return res, nil
}

func GetSelectedCluster(projectID string, pageSize int64, client TiDBCloudClient) (*Cluster, error) {
	_, clusterItems, err := RetrieveClusters(projectID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(clusterItems))
	for _, item := range clusterItems {
		items = append(items, &Cluster{
			ID:   *(item.ID),
			Name: item.Name,
		})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no available clusters found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the cluster")
	if err != nil {
		return nil, errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	clusterModel, err := p.StartReturningModel()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if m, _ := clusterModel.(ui.SelectModel); m.Interrupted {
		os.Exit(130)
	}
	cluster := clusterModel.(ui.SelectModel).GetSelectedItem().(*Cluster)
	return cluster, nil
}

// GetSelectedImport get the selected import task. statusFilter is used to filter the available options, only imports has status in statusFilter will be available.
// statusFilter with no filter will mark all the import tasks as available options just like statusFilter with all status.
func GetSelectedImport(pID string, cID string, pageSize int64, client TiDBCloudClient, statusFilter []string) (*Import, error) {
	_, importItems, err := RetrieveImports(pID, cID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(importItems))
	for _, item := range importItems {
		if len(statusFilter) != 0 && !util.ElemInSlice(statusFilter, *item.Status.Phase) {
			continue
		}

		items = append(items, &Import{
			ID:    *item.Metadata.ID,
			Phase: *item.Status.Phase,
		})
	}
	model, err := ui.InitialSelectModel(items, "Choose the import task:")
	if err != nil {
		return nil, err
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	importModel, err := p.StartReturningModel()
	if err != nil {
		return nil, err
	}
	if m, _ := importModel.(ui.SelectModel); m.Interrupted {
		os.Exit(130)
	}
	res := importModel.(ui.SelectModel).GetSelectedItem().(*Import)
	return res, nil
}

func RetrieveProjects(size int64, d TiDBCloudClient) (int64, []*projectApi.ListProjectsOKBodyItemsItems0, error) {
	params := projectApi.NewListProjectsParams()
	var total int64 = math.MaxInt64
	var page int64 = 1
	var pageSize = size
	var items []*projectApi.ListProjectsOKBodyItemsItems0
	for (page-1)*pageSize < total {
		projects, err := d.ListProjects(params.WithPage(&page).WithPageSize(&pageSize))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}

		total = *projects.Payload.Total
		page += 1
		items = append(items, projects.Payload.Items...)
	}
	return total, items, nil
}

func RetrieveClusters(pID string, pageSize int64, d TiDBCloudClient) (int64, []*clusterApi.ListClustersOfProjectOKBodyItemsItems0, error) {
	params := clusterApi.NewListClustersOfProjectParams().WithProjectID(pID)
	var total int64 = math.MaxInt64
	var page int64 = 1
	var items []*clusterApi.ListClustersOfProjectOKBodyItemsItems0
	// loop to get all clusters
	for (page-1)*pageSize < total {
		clusters, err := d.ListClustersOfProject(params.WithPage(&page).WithPageSize(&pageSize))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}

		total = *clusters.Payload.Total
		page += 1
		items = append(items, clusters.Payload.Items...)
	}
	return total, items, nil
}

func RetrieveImports(pID string, cID string, pageSize int64, d TiDBCloudClient) (int64, []*import_operations.ListImportTasksOKBodyItemsItems0, error) {
	params := import_operations.NewListImportTasksParams().WithProjectID(pID).WithClusterID(cID)
	var total int64 = math.MaxInt64
	var page int64 = 1
	var items []*import_operations.ListImportTasksOKBodyItemsItems0
	// loop to get all clusters
	for (page-1)*pageSize < total {
		imports, err := d.ListImportTasks(params.WithPage(&page).WithPageSize(&pageSize))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}

		total = *imports.Payload.Total
		page += 1
		items = append(items, imports.Payload.Items...)
	}
	return total, items, nil
}

func RetrieveConnectInfo(d TiDBCloudClient) (*connectInfoModel.ConnectInfo, error) {
	params := connectInfoApi.NewGetInfoParams()
	connectInfo, err := d.GetConnectInfo(params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return connectInfo.Payload, nil
}

func GetSelectedConnectClient(connectClientList []string) (string, error) {
	s := make([]interface{}, len(connectClientList))
	for i, v := range connectClientList {
		s[i] = v
	}
	model, err := ui.InitialSelectModel(s, "Choose the client")
	if err != nil {
		return "", errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()
	p := tea.NewProgram(model)
	connectClientModel, err := p.StartReturningModel()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := connectClientModel.(ui.SelectModel); m.Interrupted {
		os.Exit(130)
	}
	connectClient := connectClientModel.(ui.SelectModel).GetSelectedItem().(string)
	return connectClient, nil
}

func GetSelectedConnectOs(osList []string) (string, error) {
	s := make([]interface{}, len(osList))
	for i, v := range osList {
		s[i] = v
	}
	model, err := ui.InitialSelectModel(s, "Choose the operating system")
	if err != nil {
		return "", errors.Trace(err)
	}
	p := tea.NewProgram(model)
	connectClientModel, err := p.StartReturningModel()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := connectClientModel.(ui.SelectModel); m.Interrupted {
		os.Exit(130)
	}
	operatingSystem := connectClientModel.(ui.SelectModel).GetSelectedItem().(string)
	return operatingSystem, nil
}
