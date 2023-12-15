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
	"github.com/go-openapi/strfmt"
	"math"
	"strconv"

	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	connectInfoApi "tidbcloud-cli/pkg/tidbcloud/connect_info/client/connect_info_service"
	connectInfoModel "tidbcloud-cli/pkg/tidbcloud/connect_info/models"
	importApi "tidbcloud-cli/pkg/tidbcloud/import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/import/models"
	branchApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client/branch_service"
	branchModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/models"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
	serverlessModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/models"
	brApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"
	brModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/models"

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
	ID          string
	Name        string
	DisplayName string
}

type Branch struct {
	ID          string
	DisplayName string
	IsCluster   bool
}

type Region struct {
	Name        string
	DisplayName string
	Provider    string
}

type ServerlessBackup struct {
	ID         string
	Name       string
	CreateTime strfmt.DateTime
}

func (c ServerlessBackup) String() string {
	return fmt.Sprintf("%s(%s)", c.CreateTime, c.ID)
}

const (
	RestoreModeSnapshot    = "snapshot"
	RestoreModePointInTime = "point-in-time"
)

func (r Region) String() string {
	return r.DisplayName
}

func (b Branch) String() string {
	if b.IsCluster {
		return "main(the cluster)"
	}
	return fmt.Sprintf("%s(%s)", b.DisplayName, b.ID)
}

func (c Cluster) String() string {
	return fmt.Sprintf("%s(%s)", c.DisplayName, c.ID)
}

type Import struct {
	ID     string
	Status *importModel.OpenapiGetImportRespStatus
}

func (i Import) String() string {
	return fmt.Sprintf("%s(%s)", i.ID, *i.Status)
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
		return nil, util.InterruptError
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
			ID:          item.ClusterID,
			DisplayName: *item.DisplayName,
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
		return nil, util.InterruptError
	}
	cluster := clusterModel.(ui.SelectModel).GetSelectedItem().(*Cluster)
	return cluster, nil
}

func GetSelectedField(mutableFields []string) (string, error) {
	var items = make([]interface{}, 0, len(mutableFields))
	for _, item := range mutableFields {
		items = append(items, item)
	}
	model, err := ui.InitialSelectModel(items, "Choose the field to update")
	if err != nil {
		return "", errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	fieldModel, err := p.StartReturningModel()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := fieldModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	field := fieldModel.(ui.SelectModel).GetSelectedItem().(string)
	return field, nil
}

func GetSpendingLimitField(mutableFields []string) (string, error) {
	var items = make([]interface{}, 0, len(mutableFields))
	for _, item := range mutableFields {
		items = append(items, item)
	}
	model, err := ui.InitialSelectModel(items, "Choose the type of spending limit")
	if err != nil {
		return "", errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	fieldModel, err := p.StartReturningModel()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := fieldModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	field := fieldModel.(ui.SelectModel).GetSelectedItem().(string)
	return field, nil
}

func GetSelectedBranch(clusterID string, pageSize int64, client TiDBCloudClient) (*Branch, error) {
	_, branchItems, err := RetrieveBranches(clusterID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(branchItems))
	for _, item := range branchItems {
		items = append(items, &Branch{
			ID:          item.BranchID,
			DisplayName: *item.DisplayName,
			IsCluster:   false,
		})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no available branches found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the branch")
	if err != nil {
		return nil, errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	bModel, err := p.StartReturningModel()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if m, _ := bModel.(ui.SelectModel); m.Interrupted {
		return nil, util.InterruptError
	}
	branch := bModel.(ui.SelectModel).GetSelectedItem().(*Branch)
	return branch, nil
}

func GetSelectedServerlessBackup(clusterID string, pageSize int32, client TiDBCloudClient) (*ServerlessBackup, error) {
	_, backupItems, err := RetrieveServerlessBackups(clusterID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(backupItems))
	for _, item := range backupItems {
		items = append(items, &ServerlessBackup{
			ID:         item.BackupID,
			Name:       item.Name,
			CreateTime: item.CreateTime,
		})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no available backups found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the backup")
	if err != nil {
		return nil, errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	bModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if m, _ := bModel.(ui.SelectModel); m.Interrupted {
		return nil, util.InterruptError
	}
	backup := bModel.(ui.SelectModel).GetSelectedItem().(*ServerlessBackup)
	return backup, nil
}

func GetSelectedRestoreMode() (string, error) {
	items := []interface{}{
		RestoreModeSnapshot,
		RestoreModePointInTime,
	}

	model, err := ui.InitialSelectModel(items, "Choose the restore mode")
	if err != nil {
		return "", errors.Trace(err)
	}

	model.EnableFilter()
	p := tea.NewProgram(model)
	bModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := bModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	restoreMode := bModel.(ui.SelectModel).GetSelectedItem().(string)
	return restoreMode, nil
}

// GetSelectedImport get the selected import task. statusFilter is used to filter the available options, only imports has status in statusFilter will be available.
// statusFilter with no filter will mark all the import tasks as available options just like statusFilter with all status.
func GetSelectedImport(pID string, cID string, pageSize int64, client TiDBCloudClient, statusFilter []importModel.OpenapiGetImportRespStatus) (*Import, error) {
	_, importItems, err := RetrieveImports(pID, cID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(importItems))
	for _, item := range importItems {
		if len(statusFilter) != 0 && !util.ElemInSlice(statusFilter, *item.Status) {
			continue
		}

		items = append(items, &Import{
			ID:     item.ID,
			Status: item.Status,
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
		return nil, util.InterruptError
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

func RetrieveClusters(pID string, pageSize int64, d TiDBCloudClient) (int64, []*serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	params := serverlessApi.NewServerlessServiceListClustersParams()
	if pID != "" {
		projectFilter := fmt.Sprintf("projectId=%s", pID)
		params.WithFilter(&projectFilter)
	}
	pageSizeInt32 := int32(pageSize)
	var pageToken string
	var items []*serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster
	clusters, err := d.ListClustersOfProject(params.WithPageSize(&pageSizeInt32))
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, clusters.Payload.Clusters...)
	for {
		pageToken = clusters.Payload.NextPageToken
		if pageToken == "" {
			break
		}
		clusters, err = d.ListClustersOfProject(params.WithPageToken(&pageToken).WithPageSize(&pageSizeInt32))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, clusters.Payload.Clusters...)
	}
	return int64(len(items)), items, nil
}

func RetrieveBranches(cID string, pageSize int64, d TiDBCloudClient) (int64, []*branchModel.V1beta1Branch, error) {
	var items []*branchModel.V1beta1Branch
	pageSizeInt32 := int32(pageSize)
	var pageToken string

	params := branchApi.NewBranchServiceListBranchesParams().WithClusterID(cID)
	branches, err := d.ListBranches(params.WithPageSize(&pageSizeInt32))
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, branches.Payload.Branches...)
	// loop to get all branches
	for {
		pageToken = branches.Payload.NextPageToken
		if pageToken == "" {
			break
		}
		branches, err = d.ListBranches(params.WithPageSize(&pageSizeInt32).WithPageToken(&pageToken))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, branches.Payload.Branches...)
	}
	return int64(len(items)), items, nil
}

func RetrieveServerlessBackups(cID string, pageSize int32, d TiDBCloudClient) (int64, []*brModel.V1beta1Backup, error) {
	var items []*brModel.V1beta1Backup
	var pageToken string

	params := brApi.NewBackupRestoreServiceListBackupsParams().WithClusterID(cID)
	backups, err := d.ListBackups(params.WithPageSize(&pageSize))
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, backups.Payload.Backups...)
	// loop to get all backups
	for {
		pageToken = backups.Payload.NextPageToken
		if pageToken == "" {
			break
		}
		backups, err = d.ListBackups(params.WithPageSize(&pageSize).WithPageToken(&pageToken))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, backups.Payload.Backups...)
	}
	return int64(len(items)), items, nil
}

func RetrieveImports(pID string, cID string, pageSize int64, d TiDBCloudClient) (uint64, []*importModel.OpenapiGetImportResp, error) {
	params := importApi.NewListImportsParams().WithProjectID(pID).WithClusterID(cID)
	ps := int32(pageSize)
	var total uint64 = math.MaxUint64
	var page int32 = 1
	var items []*importModel.OpenapiGetImportResp
	// loop to get all clusters
	for uint64((page-1)*ps) < total {
		imports, err := d.ListImports(params.WithPage(&page).WithPageSize(&ps))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}

		total, err = strconv.ParseUint(*imports.Payload.Total, 0, 64)
		if err != nil {
			return 0, nil, errors.Annotate(err, " failed parse total import number.")
		}
		page += 1
		items = append(items, imports.Payload.Imports...)
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
		return "", util.InterruptError
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
		return "", util.InterruptError
	}
	operatingSystem := connectClientModel.(ui.SelectModel).GetSelectedItem().(string)
	return operatingSystem, nil
}
