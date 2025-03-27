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

package cloud

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/br"

	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/branch"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/export"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/imp"

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
	CreateTime *time.Time
}

type Export struct {
	DisplayName string
	ID          string
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
	return fmt.Sprintf("%s(%s)", e.DisplayName, e.ID)
}

type Import struct {
	ID     string
	Status *imp.ImportStateEnum
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
			*projectItems[0].Id,
			*projectItems[0].Name,
		}, nil
	}

	var items = make([]interface{}, 0, len(projectItems))
	for _, item := range projectItems {
		items = append(items, &Project{
			ID:   *item.Id,
			Name: *item.Name,
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
			ID:          *item.ClusterId,
			DisplayName: item.DisplayName,
			UserPrefix:  *item.UserPrefix,
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
			ID:          *item.BranchId,
			DisplayName: item.DisplayName,
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
			ID:          *item.ExportId,
			DisplayName: *item.DisplayName,
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

func GetSelectedRunningExport(ctx context.Context, clusterID string, pageSize int64, client TiDBCloudClient) (*Export, error) {
	_, exportItems, err := RetrieveExports(ctx, clusterID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(exportItems))
	for _, item := range exportItems {
		if *item.State == export.EXPORTSTATEENUM_RUNNING {
			items = append(items, &Export{
				ID:          *item.ExportId,
				DisplayName: *item.DisplayName,
			})
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no running exports found")
	}

	model, err := ui.InitialSelectModel(items, "Choose the running export:")
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
		if *item.Target.Type == export.EXPORTTARGETTYPEENUM_LOCAL && *item.State == export.EXPORTSTATEENUM_SUCCEEDED {
			items = append(items, &Export{
				ID:          *item.ExportId,
				DisplayName: *item.DisplayName,
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
			ID:         *item.BackupId,
			Name:       *item.Name,
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
func GetSelectedImport(ctx context.Context, cID string, pageSize int64, client TiDBCloudClient, statusFilter []imp.ImportStateEnum) (*Import, error) {
	_, importItems, err := RetrieveImports(ctx, cID, pageSize, client)
	if err != nil {
		return nil, err
	}

	var items = make([]interface{}, 0, len(importItems))
	for _, item := range importItems {
		if len(statusFilter) != 0 && !slices.Contains(statusFilter, *item.State) {
			continue
		}

		items = append(items, &Import{
			ID:     *item.ImportId,
			Status: item.State,
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
	imp, err := p.Run()
	if err != nil {
		return nil, err
	}
	if m, _ := imp.(ui.SelectModel); m.Interrupted {
		return nil, util.InterruptError
	}
	res := imp.(ui.SelectModel).GetSelectedItem()
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
			UserName: *item.UserName,
			Role:     util.GetDisplayRole(*item.BuiltinRole, item.CustomRoles),
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

func RetrieveProjects(ctx context.Context, pageSize int64, d TiDBCloudClient) (int64, []iam.ApiProject, error) {
	var items []iam.ApiProject
	pageSizeInt32 := int32(pageSize)
	var pageToken *string

	// loop to get all projects
	for {
		projects, err := d.ListProjects(ctx, &pageSizeInt32, pageToken)
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, projects.Projects...)

		pageToken = projects.NextPageToken
		if util.IsNilOrEmpty(pageToken) {
			break
		}
	}
	return int64(len(items)), items, nil
}

func RetrieveClusters(ctx context.Context, pID string, pageSize int64, d TiDBCloudClient) (int64, []cluster.TidbCloudOpenApiserverlessv1beta1Cluster, error) {
	var items []cluster.TidbCloudOpenApiserverlessv1beta1Cluster
	pageSizeInt32 := int32(pageSize)
	var pageToken *string
	var filter *string
	if pID != "" {
		projectFilter := fmt.Sprintf("projectId=%s", pID)
		filter = &projectFilter
	}

	clusters, err := d.ListClusters(ctx, filter, &pageSizeInt32, nil, nil, nil)
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, clusters.Clusters...)
	for {
		pageToken = clusters.NextPageToken
		if util.IsNilOrEmpty(pageToken) {
			break
		}
		clusters, err = d.ListClusters(ctx, filter, &pageSizeInt32, pageToken, nil, nil)
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, clusters.Clusters...)
	}
	return int64(len(items)), items, nil
}

func RetrieveBranches(ctx context.Context, cID string, pageSize int64, d TiDBCloudClient) (int64, []branch.Branch, error) {
	var items []branch.Branch
	pageSizeInt32 := int32(pageSize)
	var pageToken *string

	branches, err := d.ListBranches(ctx, cID, &pageSizeInt32, nil)
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, branches.Branches...)
	// loop to get all branches
	for {
		pageToken = branches.NextPageToken
		if util.IsNilOrEmpty(pageToken) {
			break
		}
		branches, err = d.ListBranches(ctx, cID, &pageSizeInt32, pageToken)
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, branches.Branches...)
	}
	return int64(len(items)), items, nil
}

func RetrieveExports(ctx context.Context, cID string, pageSize int64, d TiDBCloudClient) (int64, []export.Export, error) {
	var items []export.Export
	pageSizeInt32 := int32(pageSize)
	var pageToken *string

	orderBy := "create_time desc"
	exports, err := d.ListExports(ctx, cID, &pageSizeInt32, nil, &orderBy)
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, exports.Exports...)
	// loop to get all branches
	for {
		pageToken = exports.NextPageToken
		if util.IsNilOrEmpty(pageToken) {
			break
		}
		exports, err = d.ListExports(ctx, cID, &pageSizeInt32, pageToken, &orderBy)
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, exports.Exports...)
	}
	return int64(len(items)), items, nil
}

func RetrieveServerlessBackups(ctx context.Context, cID string, pageSize int32, d TiDBCloudClient) (int64, []br.V1beta1Backup, error) {
	var items []br.V1beta1Backup
	var pageToken *string

	backups, err := d.ListBackups(ctx, &cID, &pageSize, nil)
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, backups.Backups...)
	// loop to get all backups
	for {
		pageToken = backups.NextPageToken
		if util.IsNilOrEmpty(pageToken) {
			break
		}
		backups, err = d.ListBackups(ctx, &cID, &pageSize, pageToken)
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, backups.Backups...)
	}
	return int64(len(items)), items, nil
}

func RetrieveImports(context context.Context, cID string, pageSize int64, d TiDBCloudClient) (int64, []imp.Import, error) {
	orderBy := "create_time desc"
	ps := int32(pageSize)
	var items []imp.Import
	imports, err := d.ListImports(context, cID, &ps, nil, &orderBy)
	if err != nil {
		return 0, nil, errors.Trace(err)
	}
	items = append(items, imports.Imports...)
	var pageToken *string
	for {
		pageToken = imports.NextPageToken
		if util.IsNilOrEmpty(pageToken) {
			break
		}
		imports, err = d.ListImports(context, cID, &ps, pageToken, &orderBy)
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, imports.Imports...)
	}
	return int64(len(items)), items, nil
}

func RetrieveSQLUsers(ctx context.Context, cID string, pageSize int64, d TiDBCloudClient) (int64, []iam.ApiSqlUser, error) {
	var items []iam.ApiSqlUser

	pageSizeInt32 := int32(pageSize)
	var pageToken *string
	// loop to get all SQL users
	for {
		sqlUsers, err := d.ListSQLUsers(ctx, cID, &pageSizeInt32, pageToken)
		if err != nil {
			return 0, nil, errors.Trace(err)
		}
		items = append(items, sqlUsers.SqlUsers...)

		pageToken = sqlUsers.NextPageToken
		if util.IsNilOrEmpty(pageToken) {
			break
		}
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
			ID:          *item.BranchId,
			DisplayName: item.DisplayName,
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

func GetAllExportFiles(ctx context.Context, cID string, eID string, d TiDBCloudClient) ([]export.ExportFile, error) {
	var items []export.ExportFile
	var pageSize int32 = 1000
	var pageToken *string
	exportFilesResp, err := d.ListExportFiles(ctx, cID, eID, &pageSize, nil, false)
	if err != nil {
		return nil, errors.Trace(err)
	}
	items = append(items, exportFilesResp.Files...)
	for {
		pageToken = exportFilesResp.NextPageToken
		if util.IsNilOrEmpty(pageToken) {
			break
		}
		exportFilesResp, err = d.ListExportFiles(ctx, cID, eID, &pageSize, pageToken, false)
		if err != nil {
			return nil, errors.Trace(err)
		}
		items = append(items, exportFilesResp.Files...)
	}
	return items, nil
}
