// Copyright 2024 PingCAP, Inc.
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
	"slices"

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
	UserPrefix  string
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
		return fmt.Sprintf("%s(main)", b.DisplayName)
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

type SQLUser struct {
	UserName string
	Role     string
}

func (s SQLUser) String() string {
	return fmt.Sprintf("%s(%s)", s.UserName, s.Role)
}

func GetSelectedProject(ctx context.Context, pageSize int64, client TiDBCloudClient) (*Project, error) {
	_, projectItems, err := RetrieveProjects(ctx, pageSize, client)
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
	projectModel, err := p.Run()
	if err != nil {
		return nil, err
	}
	if m, _ := projectModel.(ui.SelectModel); m.Interrupted {
		return nil, util.InterruptError
	}
	res := projectModel.(ui.SelectModel).GetSelectedItem()
	if res == nil {
		return nil, errors.New("no project selected")
	}
	return res.(*Project), nil
}

func GetSelectedCluster(ctx context.Context, projectID string, pageSize int64, client TiDBCloudClient) (*Cluster, error) {
	_, clusterItems, err := RetrieveClusters(ctx, projectID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(clusterItems))
	for _, item := range clusterItems {
		items = append(items, &Cluster{
			ID:          item.ClusterID,
			DisplayName: *item.DisplayName,
			UserPrefix:  item.UserPrefix,
		})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no available clusters found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the cluster:")
	if err != nil {
		return nil, errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	clusterModel, err := p.Run()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if m, _ := clusterModel.(ui.SelectModel); m.Interrupted {
		return nil, util.InterruptError
	}
	cluster := clusterModel.(ui.SelectModel).GetSelectedItem()
	if cluster == nil {
		return nil, errors.New("no cluster selected")
	}
	return cluster.(*Cluster), nil
}

func GetSelectedField(mutableFields []string) (string, error) {
	var items = make([]interface{}, 0, len(mutableFields))
	for _, item := range mutableFields {
		items = append(items, item)
	}
	model, err := ui.InitialSelectModel(items, "Choose the field to update:")
	if err != nil {
		return "", errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	fieldModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := fieldModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	field := fieldModel.(ui.SelectModel).GetSelectedItem()
	if field == nil {
		return "", errors.New("no field selected")
	}
	return field.(string), nil
}

func GetSelectedBool(notice string) (bool, error) {
	items := []interface{}{
		"true",
		"false",
	}

	model, err := ui.InitialSelectModel(items, notice)
	if err != nil {
		return false, errors.Trace(err)
	}

	p := tea.NewProgram(model)
	bModel, err := p.Run()
	if err != nil {
		return false, errors.Trace(err)
	}
	if m, _ := bModel.(ui.SelectModel); m.Interrupted {
		return false, util.InterruptError
	}
	value := bModel.(ui.SelectModel).GetSelectedItem()
	if value == nil {
		return false, errors.New("no value selected")
	}
	return value.(string) == "true", nil
}

func GetSpendingLimitField(mutableFields []string) (string, error) {
	var items = make([]interface{}, 0, len(mutableFields))
	for _, item := range mutableFields {
		items = append(items, item)
	}
	model, err := ui.InitialSelectModel(items, "Choose the type of spending limit:")
	if err != nil {
		return "", errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	fieldModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := fieldModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	field := fieldModel.(ui.SelectModel).GetSelectedItem()
	if field == nil {
		return "", errors.New("no type selected")
	}
	return field.(string), nil
}

func GetSelectedBranch(ctx context.Context, clusterID string, pageSize int64, client TiDBCloudClient) (*Branch, error) {
	_, branchItems, err := RetrieveBranches(ctx, clusterID, pageSize, client)
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

	model, err := ui.InitialSelectModel(items, "Choose the branch:")
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
	branch := bModel.(ui.SelectModel).GetSelectedItem()
	if branch == nil {
		return nil, errors.New("no branch selected")
	}
	return branch.(*Branch), nil
}

func GetSelectedExport(ctx context.Context, clusterID string, pageSize int64, client TiDBCloudClient) (*Export, error) {
	_, exportItems, err := RetrieveExports(ctx, clusterID, pageSize, client)
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

	model, err := ui.InitialSelectModel(items, "Choose the export:")
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
	export := eModel.(ui.SelectModel).GetSelectedItem()
	if export == nil {
		return nil, errors.New("no export selected")
	}
	return export.(*Export), nil
}

func GetSelectedLocalExport(ctx context.Context, clusterID string, pageSize int64, client TiDBCloudClient) (*Export, error) {
	_, exportItems, err := RetrieveExports(ctx, clusterID, pageSize, client)
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

	model, err := ui.InitialSelectModel(items, "Choose the succeeded local type export:")
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
	export := eModel.(ui.SelectModel).GetSelectedItem()
	if export == nil {
		return nil, errors.New("no export selected")
	}
	return export.(*Export), nil
}

func GetSelectedServerlessBackup(ctx context.Context, clusterID string, pageSize int32, client TiDBCloudClient) (*ServerlessBackup, error) {
	_, backupItems, err := RetrieveServerlessBackups(ctx, clusterID, pageSize, client)
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

	model, err := ui.InitialSelectModel(items, "Choose the backup:")
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
	backup := bModel.(ui.SelectModel).GetSelectedItem()
	if backup == nil {
		return nil, errors.New("no backup selected")
	}
	return backup.(*ServerlessBackup), nil
}

func GetSelectedRestoreMode() (string, error) {
	items := []interface{}{
		RestoreModeSnapshot,
		RestoreModePointInTime,
	}

	model, err := ui.InitialSelectModel(items, "Choose the restore mode:")
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
	restoreMode := bModel.(ui.SelectModel).GetSelectedItem()
	if restoreMode == nil {
		return "", errors.New("no restore mode selected")
	}
	return restoreMode.(string), nil
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
		if len(statusFilter) != 0 && !slices.Contains(statusFilter, item.State) {
			continue
		}

		items = append(items, &Import{
			ID:     item.ID,
			Status: &item.State,
		})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no available imports found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the import task:")
	if err != nil {
		return nil, err
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	importModel, err := p.Run()
	if err != nil {
		return nil, err
	}
	if m, _ := importModel.(ui.SelectModel); m.Interrupted {
		return nil, util.InterruptError
	}
	res := importModel.(ui.SelectModel).GetSelectedItem()
	if res == nil {
		return nil, errors.New("no import task selected")
	}
	return res.(*Import), nil
}

func GetSelectedBuiltinRole() (string, error) {
	items := []interface{}{
		util.GetDisplayRole(util.ADMIN_ROLE, []string{}),
		util.GetDisplayRole(util.READWRITE_ROLE, []string{}),
		util.GetDisplayRole(util.READONLY_ROLE, []string{}),
	}

	model, err := ui.InitialSelectModel(items, "Choose the role:")
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
	role := bModel.(ui.SelectModel).GetSelectedItem()
	if role == nil {
		return "", errors.New("no role selected")
	}
	return role.(string), nil
}

func GetSelectedSQLUser(ctx context.Context, clusterID string, pageSize int64, client TiDBCloudClient) (string, error) {
	_, sqlUserItems, err := RetrieveSQLUsers(ctx, clusterID, pageSize, client)
	if err != nil {
		return "", err
	}

	var items = make([]interface{}, 0, len(sqlUserItems))
	for _, item := range sqlUserItems {
		items = append(items, &SQLUser{
			UserName: item.UserName,
			Role:     util.GetDisplayRole(item.BuiltinRole, item.CustomRoles),
		})
	}
	if len(items) == 0 {
		return "", fmt.Errorf("no available sql-users found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the SQL user:")
	if err != nil {
		return "", err
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	sqlUserModel, err := p.Run()
	if err != nil {
		return "", err
	}
	if m, _ := sqlUserModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	res := sqlUserModel.(ui.SelectModel).GetSelectedItem()
	if res == nil {
		return "", errors.New("no SQL user selected")
	}

	return res.(*SQLUser).UserName, nil

}

func RetrieveProjects(ctx context.Context, pageSize int64, d TiDBCloudClient) (int64, []*iamModel.APIProject, error) {
	var items []*iamModel.APIProject
	var pageToken string

	params := iamApi.NewGetV1beta1ProjectsParams().WithPageSize(&pageSize).WithContext(ctx)
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

func RetrieveClusters(ctx context.Context, pID string, pageSize int64, d TiDBCloudClient) (int64, []*serverlessModel.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	params := serverlessApi.NewServerlessServiceListClustersParams().WithContext(ctx)
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

func RetrieveBranches(ctx context.Context, cID string, pageSize int64, d TiDBCloudClient) (int64, []*branchModel.V1beta1Branch, error) {
	var items []*branchModel.V1beta1Branch
	pageSizeInt32 := int32(pageSize)
	var pageToken string

	params := branchApi.NewBranchServiceListBranchesParams().WithClusterID(cID).WithContext(ctx)
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

func RetrieveExports(ctx context.Context, cID string, pageSize int64, d TiDBCloudClient) (int64, []*exportModel.V1beta1Export, error) {
	var items []*exportModel.V1beta1Export
	pageSizeInt32 := int32(pageSize)
	var pageToken string

	orderBy := "create_time desc"
	params := exportApi.NewExportServiceListExportsParams().WithClusterID(cID).WithPageSize(&pageSizeInt32).
		WithOrderBy(&orderBy).WithContext(ctx)
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

func RetrieveServerlessBackups(ctx context.Context, cID string, pageSize int32, d TiDBCloudClient) (int64, []*brModel.V1beta1Backup, error) {
	var items []*brModel.V1beta1Backup
	var pageToken string

	params := brApi.NewBackupRestoreServiceListBackupsParams().WithClusterID(cID).WithContext(ctx)
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

func RetrieveImports(context context.Context, cID string, pageSize int64, d TiDBCloudClient) (int64, []*importModel.V1beta1Import, error) {
	params := importApi.NewImportServiceListImportsParams().WithClusterID(cID).WithContext(context)
	ps := int32(pageSize)
	var items []*importModel.V1beta1Import
	imports, err := d.ListImports(params.WithPageSize(&ps))
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, imports.Payload.Imports...)
	var pageToken string
	for {
		pageToken = imports.Payload.NextPageToken
		if pageToken == "" {
			break
		}
		imports, err = d.ListImports(params.WithPageToken(&pageToken).WithPageSize(&ps))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, imports.Payload.Imports...)
	}
	return int64(len(items)), items, nil
}

func RetrieveSQLUsers(ctx context.Context, cID string, pageSize int64, d TiDBCloudClient) (int64, []*iamModel.APISQLUser, error) {
	var items []*iamModel.APISQLUser
	var pageToken string

	params := iamApi.NewGetV1beta1ClustersClusterIDSQLUsersParams().
		WithClusterID(cID).
		WithPageSize(&pageSize).
		WithContext(ctx)
	users, err := d.ListSQLUsers(params)
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, users.Payload.SQLUsers...)
	// loop to get all SQL users
	for {
		pageToken = users.Payload.NextPageToken
		if pageToken == "" {
			break
		}
		users, err = d.ListSQLUsers(params.WithPageSize(&pageSize).WithPageToken(&pageToken))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, users.Payload.SQLUsers...)
	}
	return int64(len(items)), items, nil
}

func GetSelectedParentID(ctx context.Context, cluster *Cluster, pageSize int64, client TiDBCloudClient) (string, error) {
	clusterID := cluster.ID
	_, branchItems, err := RetrieveBranches(ctx, clusterID, pageSize, client)
	if err != nil {
		return "", err
	}
	// If there is no branch, return the clusterID directly.
	if len(branchItems) == 0 {
		return clusterID, nil
	}

	var items = make([]interface{}, 0, len(branchItems)+1)
	items = append(items, &Branch{
		ID:          clusterID,
		DisplayName: cluster.DisplayName,
		IsCluster:   true,
	})
	for _, item := range branchItems {
		items = append(items, &Branch{
			ID:          item.BranchID,
			DisplayName: *item.DisplayName,
			IsCluster:   false,
		})
	}

	model, err := ui.InitialSelectModel(items, "Choose the parent:")
	if err != nil {
		return "", errors.Trace(err)
	}
	itemsPerPage := 6
	model.EnablePagination(itemsPerPage)
	model.EnableFilter()

	p := tea.NewProgram(model)
	bModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := bModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	parent := bModel.(ui.SelectModel).GetSelectedItem()
	if parent == nil {
		return "", errors.New("no parent selected")
	}
	return parent.(*Branch).ID, nil
}
