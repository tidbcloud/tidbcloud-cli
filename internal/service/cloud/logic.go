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
	"context"
	"fmt"
	"math"

	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	branchApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client/branch_service"
	branchModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/models"
	iamApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/client/account"
	iamModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/models"
	serverlessApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
	serverlessModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/models"
	brApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"
	brModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/models"
	exportApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/client/export_service"
	exportModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/models"
	importApi "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"
	importModel "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-openapi/strfmt"
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

type Export struct {
	ID string
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

func (e Export) String() string {
	return e.ID
}

type Import struct {
	ID     string
	Status *importModel.V1beta1ImportState
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

func GetSelectedExport(clusterID string, pageSize int64, client TiDBCloudClient) (*Export, error) {
	_, exportItems, err := RetrieveExports(clusterID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(exportItems))
	for _, item := range exportItems {
		items = append(items, &Export{
			ID: item.ExportID,
		})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no available exports found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the export")
	if err != nil {
		return nil, errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	eModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if m, _ := eModel.(ui.SelectModel); m.Interrupted {
		return nil, util.InterruptError
	}
	export := eModel.(ui.SelectModel).GetSelectedItem().(*Export)
	return export, nil
}

func GetSelectedLocalExport(clusterID string, pageSize int64, client TiDBCloudClient) (*Export, error) {
	_, exportItems, err := RetrieveExports(clusterID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(exportItems))
	for _, item := range exportItems {
		if item.Target.Type == exportModel.TargetTargetTypeLOCAL && item.State == exportModel.V1beta1ExportStateSUCCEEDED {
			items = append(items, &Export{
				ID: item.ExportID,
			})
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no available exports found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the succeed local type export")
	if err != nil {
		return nil, errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	eModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if m, _ := eModel.(ui.SelectModel); m.Interrupted {
		return nil, util.InterruptError
	}
	export := eModel.(ui.SelectModel).GetSelectedItem().(*Export)
	return export, nil
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
func GetSelectedImport(ctx context.Context, cID string, pageSize int64, client TiDBCloudClient, statusFilter []importModel.V1beta1ImportState) (*Import, error) {
	_, importItems, err := RetrieveImports(ctx, cID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(importItems))
	for _, item := range importItems {
		if len(statusFilter) != 0 && !util.ElemInSlice(statusFilter, item.State) {
			continue
		}

		items = append(items, &Import{
			ID:     item.ID,
			Status: &item.State,
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

func RetrieveProjects(pageSize int64, d TiDBCloudClient) (int64, []*iamModel.APIProject, error) {
	var items []*iamModel.APIProject
	var pageToken string

	params := iamApi.NewGetV1beta1ProjectsParams().WithPageSize(&pageSize)
	projects, err := d.ListProjects(params)
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, projects.Payload.Projects...)
	// loop to get all projects
	for {
		pageToken = projects.Payload.NextPageToken
		if pageToken == "" {
			break
		}
		projects, err = d.ListProjects(params.WithPageSize(&pageSize).WithPageToken(&pageToken))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, projects.Payload.Projects...)
	}
	return int64(len(items)), items, nil
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

func RetrieveExports(cID string, pageSize int64, d TiDBCloudClient) (int64, []*exportModel.V1beta1Export, error) {
	var items []*exportModel.V1beta1Export
	pageSizeInt32 := int32(pageSize)
	var pageToken string

	params := exportApi.NewExportServiceListExportsParams().WithClusterID(cID).WithPageSize(&pageSizeInt32)
	exports, err := d.ListExports(params)
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, exports.Payload.Exports...)
	// loop to get all branches
	for {
		pageToken = exports.Payload.NextPageToken
		if pageToken == "" {
			break
		}
		exports, err = d.ListExports(params.WithPageSize(&pageSizeInt32).WithPageToken(&pageToken))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, exports.Payload.Exports...)
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

func RetrieveImports(context context.Context, cID string, pageSize int64, d TiDBCloudClient) (uint64, []*importModel.V1beta1Import, error) {
	params := importApi.NewImportServiceListImportsParams().WithClusterID(cID)
	ps := int32(pageSize)
	var total uint64 = math.MaxUint64
	var page int32 = 1
	var items []*importModel.V1beta1Import
	// loop to get all clusters
	for uint64((page-1)*ps) < total {
		imports, err := d.ListImports(params.WithPage(&page).WithPageSize(&ps).WithContext(context))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}

		total = uint64(imports.Payload.Total)
		if err != nil {
			return 0, nil, errors.Annotate(err, " failed parse total import number.")
		}
		page += 1
		items = append(items, imports.Payload.Imports...)
	}
	return total, items, nil
}
