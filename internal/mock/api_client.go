// Code generated by mockery v2.42.1. DO NOT EDIT.

package mock

import (
	branch_service "tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/client/branch_service"
	backup_restore_service "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_br/client/backup_restore_service"

	iam "tidbcloud-cli/pkg/tidbcloud/v1beta1/iam"

	import_service "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/client/import_service"

	mock "github.com/stretchr/testify/mock"

	operations "tidbcloud-cli/pkg/tidbcloud/pingchat/client/operations"

	project "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"

	serverless_service "tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/client/serverless_service"
)

// TiDBCloudClient is an autogenerated mock type for the TiDBCloudClient type
type TiDBCloudClient struct {
	mock.Mock
}

// CancelImport provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) CancelImport(params *import_service.ImportServiceCancelImportParams, opts ...import_service.ClientOption) (*import_service.ImportServiceCancelImportOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CancelImport")
	}

	var r0 *import_service.ImportServiceCancelImportOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceCancelImportParams, ...import_service.ClientOption) (*import_service.ImportServiceCancelImportOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceCancelImportParams, ...import_service.ClientOption) *import_service.ImportServiceCancelImportOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*import_service.ImportServiceCancelImportOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*import_service.ImportServiceCancelImportParams, ...import_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CancelMultipartUpload provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) CancelMultipartUpload(params *import_service.ImportServiceCancelMultipartUploadParams, opts ...import_service.ClientOption) (*import_service.ImportServiceCancelMultipartUploadOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CancelMultipartUpload")
	}

	var r0 *import_service.ImportServiceCancelMultipartUploadOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceCancelMultipartUploadParams, ...import_service.ClientOption) (*import_service.ImportServiceCancelMultipartUploadOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceCancelMultipartUploadParams, ...import_service.ClientOption) *import_service.ImportServiceCancelMultipartUploadOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*import_service.ImportServiceCancelMultipartUploadOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*import_service.ImportServiceCancelMultipartUploadParams, ...import_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Chat provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) Chat(params *operations.ChatParams, opts ...operations.ClientOption) (*operations.ChatOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Chat")
	}

	var r0 *operations.ChatOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*operations.ChatParams, ...operations.ClientOption) (*operations.ChatOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*operations.ChatParams, ...operations.ClientOption) *operations.ChatOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*operations.ChatOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*operations.ChatParams, ...operations.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CompleteMultipartUpload provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) CompleteMultipartUpload(params *import_service.ImportServiceCompleteMultipartUploadParams, opts ...import_service.ClientOption) (*import_service.ImportServiceCompleteMultipartUploadOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CompleteMultipartUpload")
	}

	var r0 *import_service.ImportServiceCompleteMultipartUploadOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceCompleteMultipartUploadParams, ...import_service.ClientOption) (*import_service.ImportServiceCompleteMultipartUploadOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceCompleteMultipartUploadParams, ...import_service.ClientOption) *import_service.ImportServiceCompleteMultipartUploadOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*import_service.ImportServiceCompleteMultipartUploadOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*import_service.ImportServiceCompleteMultipartUploadParams, ...import_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateBranch provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) CreateBranch(params *branch_service.BranchServiceCreateBranchParams, opts ...branch_service.ClientOption) (*branch_service.BranchServiceCreateBranchOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateBranch")
	}

	var r0 *branch_service.BranchServiceCreateBranchOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*branch_service.BranchServiceCreateBranchParams, ...branch_service.ClientOption) (*branch_service.BranchServiceCreateBranchOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*branch_service.BranchServiceCreateBranchParams, ...branch_service.ClientOption) *branch_service.BranchServiceCreateBranchOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*branch_service.BranchServiceCreateBranchOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*branch_service.BranchServiceCreateBranchParams, ...branch_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCluster provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) CreateCluster(params *serverless_service.ServerlessServiceCreateClusterParams, opts ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceCreateClusterOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateCluster")
	}

	var r0 *serverless_service.ServerlessServiceCreateClusterOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceCreateClusterParams, ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceCreateClusterOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceCreateClusterParams, ...serverless_service.ClientOption) *serverless_service.ServerlessServiceCreateClusterOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serverless_service.ServerlessServiceCreateClusterOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*serverless_service.ServerlessServiceCreateClusterParams, ...serverless_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateImport provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) CreateImport(params *import_service.ImportServiceCreateImportParams, opts ...import_service.ClientOption) (*import_service.ImportServiceCreateImportOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateImport")
	}

	var r0 *import_service.ImportServiceCreateImportOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceCreateImportParams, ...import_service.ClientOption) (*import_service.ImportServiceCreateImportOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceCreateImportParams, ...import_service.ClientOption) *import_service.ImportServiceCreateImportOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*import_service.ImportServiceCreateImportOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*import_service.ImportServiceCreateImportParams, ...import_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteBackup provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) DeleteBackup(params *backup_restore_service.BackupRestoreServiceDeleteBackupParams, opts ...backup_restore_service.ClientOption) (*backup_restore_service.BackupRestoreServiceDeleteBackupOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteBackup")
	}

	var r0 *backup_restore_service.BackupRestoreServiceDeleteBackupOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*backup_restore_service.BackupRestoreServiceDeleteBackupParams, ...backup_restore_service.ClientOption) (*backup_restore_service.BackupRestoreServiceDeleteBackupOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*backup_restore_service.BackupRestoreServiceDeleteBackupParams, ...backup_restore_service.ClientOption) *backup_restore_service.BackupRestoreServiceDeleteBackupOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*backup_restore_service.BackupRestoreServiceDeleteBackupOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*backup_restore_service.BackupRestoreServiceDeleteBackupParams, ...backup_restore_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteBranch provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) DeleteBranch(params *branch_service.BranchServiceDeleteBranchParams, opts ...branch_service.ClientOption) (*branch_service.BranchServiceDeleteBranchOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteBranch")
	}

	var r0 *branch_service.BranchServiceDeleteBranchOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*branch_service.BranchServiceDeleteBranchParams, ...branch_service.ClientOption) (*branch_service.BranchServiceDeleteBranchOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*branch_service.BranchServiceDeleteBranchParams, ...branch_service.ClientOption) *branch_service.BranchServiceDeleteBranchOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*branch_service.BranchServiceDeleteBranchOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*branch_service.BranchServiceDeleteBranchParams, ...branch_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCluster provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) DeleteCluster(params *serverless_service.ServerlessServiceDeleteClusterParams, opts ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceDeleteClusterOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCluster")
	}

	var r0 *serverless_service.ServerlessServiceDeleteClusterOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceDeleteClusterParams, ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceDeleteClusterOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceDeleteClusterParams, ...serverless_service.ClientOption) *serverless_service.ServerlessServiceDeleteClusterOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serverless_service.ServerlessServiceDeleteClusterOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*serverless_service.ServerlessServiceDeleteClusterParams, ...serverless_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBackup provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) GetBackup(params *backup_restore_service.BackupRestoreServiceGetBackupParams, opts ...backup_restore_service.ClientOption) (*backup_restore_service.BackupRestoreServiceGetBackupOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetBackup")
	}

	var r0 *backup_restore_service.BackupRestoreServiceGetBackupOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*backup_restore_service.BackupRestoreServiceGetBackupParams, ...backup_restore_service.ClientOption) (*backup_restore_service.BackupRestoreServiceGetBackupOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*backup_restore_service.BackupRestoreServiceGetBackupParams, ...backup_restore_service.ClientOption) *backup_restore_service.BackupRestoreServiceGetBackupOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*backup_restore_service.BackupRestoreServiceGetBackupOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*backup_restore_service.BackupRestoreServiceGetBackupParams, ...backup_restore_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBranch provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) GetBranch(params *branch_service.BranchServiceGetBranchParams, opts ...branch_service.ClientOption) (*branch_service.BranchServiceGetBranchOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetBranch")
	}

	var r0 *branch_service.BranchServiceGetBranchOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*branch_service.BranchServiceGetBranchParams, ...branch_service.ClientOption) (*branch_service.BranchServiceGetBranchOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*branch_service.BranchServiceGetBranchParams, ...branch_service.ClientOption) *branch_service.BranchServiceGetBranchOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*branch_service.BranchServiceGetBranchOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*branch_service.BranchServiceGetBranchParams, ...branch_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCluster provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) GetCluster(params *serverless_service.ServerlessServiceGetClusterParams, opts ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceGetClusterOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetCluster")
	}

	var r0 *serverless_service.ServerlessServiceGetClusterOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceGetClusterParams, ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceGetClusterOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceGetClusterParams, ...serverless_service.ClientOption) *serverless_service.ServerlessServiceGetClusterOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serverless_service.ServerlessServiceGetClusterOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*serverless_service.ServerlessServiceGetClusterParams, ...serverless_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetImport provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) GetImport(params *import_service.ImportServiceGetImportParams, opts ...import_service.ClientOption) (*import_service.ImportServiceGetImportOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetImport")
	}

	var r0 *import_service.ImportServiceGetImportOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceGetImportParams, ...import_service.ClientOption) (*import_service.ImportServiceGetImportOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceGetImportParams, ...import_service.ClientOption) *import_service.ImportServiceGetImportOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*import_service.ImportServiceGetImportOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*import_service.ImportServiceGetImportParams, ...import_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListBackups provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) ListBackups(params *backup_restore_service.BackupRestoreServiceListBackupsParams, opts ...backup_restore_service.ClientOption) (*backup_restore_service.BackupRestoreServiceListBackupsOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListBackups")
	}

	var r0 *backup_restore_service.BackupRestoreServiceListBackupsOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*backup_restore_service.BackupRestoreServiceListBackupsParams, ...backup_restore_service.ClientOption) (*backup_restore_service.BackupRestoreServiceListBackupsOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*backup_restore_service.BackupRestoreServiceListBackupsParams, ...backup_restore_service.ClientOption) *backup_restore_service.BackupRestoreServiceListBackupsOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*backup_restore_service.BackupRestoreServiceListBackupsOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*backup_restore_service.BackupRestoreServiceListBackupsParams, ...backup_restore_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListBranches provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) ListBranches(params *branch_service.BranchServiceListBranchesParams, opts ...branch_service.ClientOption) (*branch_service.BranchServiceListBranchesOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListBranches")
	}

	var r0 *branch_service.BranchServiceListBranchesOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*branch_service.BranchServiceListBranchesParams, ...branch_service.ClientOption) (*branch_service.BranchServiceListBranchesOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*branch_service.BranchServiceListBranchesParams, ...branch_service.ClientOption) *branch_service.BranchServiceListBranchesOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*branch_service.BranchServiceListBranchesOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*branch_service.BranchServiceListBranchesParams, ...branch_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListClustersOfProject provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) ListClustersOfProject(params *serverless_service.ServerlessServiceListClustersParams, opts ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceListClustersOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListClustersOfProject")
	}

	var r0 *serverless_service.ServerlessServiceListClustersOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceListClustersParams, ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceListClustersOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceListClustersParams, ...serverless_service.ClientOption) *serverless_service.ServerlessServiceListClustersOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serverless_service.ServerlessServiceListClustersOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*serverless_service.ServerlessServiceListClustersParams, ...serverless_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListImports provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) ListImports(params *import_service.ImportServiceListImportsParams, opts ...import_service.ClientOption) (*import_service.ImportServiceListImportsOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListImports")
	}

	var r0 *import_service.ImportServiceListImportsOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceListImportsParams, ...import_service.ClientOption) (*import_service.ImportServiceListImportsOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceListImportsParams, ...import_service.ClientOption) *import_service.ImportServiceListImportsOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*import_service.ImportServiceListImportsOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*import_service.ImportServiceListImportsParams, ...import_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProjects provides a mock function with given fields: params
func (_m *TiDBCloudClient) ListProjects(params *iam.ListProjectsParams) (*project.ListProjectsOK, error) {
	ret := _m.Called(params)

	if len(ret) == 0 {
		panic("no return value specified for ListProjects")
	}

	var r0 *project.ListProjectsOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam.ListProjectsParams) (*project.ListProjectsOK, error)); ok {
		return rf(params)
	}
	if rf, ok := ret.Get(0).(func(*iam.ListProjectsParams) *project.ListProjectsOK); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.ListProjectsOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam.ListProjectsParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProviderRegions provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) ListProviderRegions(params *serverless_service.ServerlessServiceListRegionsParams, opts ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceListRegionsOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListProviderRegions")
	}

	var r0 *serverless_service.ServerlessServiceListRegionsOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceListRegionsParams, ...serverless_service.ClientOption) (*serverless_service.ServerlessServiceListRegionsOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServiceListRegionsParams, ...serverless_service.ClientOption) *serverless_service.ServerlessServiceListRegionsOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serverless_service.ServerlessServiceListRegionsOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*serverless_service.ServerlessServiceListRegionsParams, ...serverless_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PartialUpdateCluster provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) PartialUpdateCluster(params *serverless_service.ServerlessServicePartialUpdateClusterParams, opts ...serverless_service.ClientOption) (*serverless_service.ServerlessServicePartialUpdateClusterOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PartialUpdateCluster")
	}

	var r0 *serverless_service.ServerlessServicePartialUpdateClusterOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServicePartialUpdateClusterParams, ...serverless_service.ClientOption) (*serverless_service.ServerlessServicePartialUpdateClusterOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*serverless_service.ServerlessServicePartialUpdateClusterParams, ...serverless_service.ClientOption) *serverless_service.ServerlessServicePartialUpdateClusterOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serverless_service.ServerlessServicePartialUpdateClusterOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*serverless_service.ServerlessServicePartialUpdateClusterParams, ...serverless_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Restore provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) Restore(params *backup_restore_service.BackupRestoreServiceRestoreParams, opts ...backup_restore_service.ClientOption) (*backup_restore_service.BackupRestoreServiceRestoreOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Restore")
	}

	var r0 *backup_restore_service.BackupRestoreServiceRestoreOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*backup_restore_service.BackupRestoreServiceRestoreParams, ...backup_restore_service.ClientOption) (*backup_restore_service.BackupRestoreServiceRestoreOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*backup_restore_service.BackupRestoreServiceRestoreParams, ...backup_restore_service.ClientOption) *backup_restore_service.BackupRestoreServiceRestoreOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*backup_restore_service.BackupRestoreServiceRestoreOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*backup_restore_service.BackupRestoreServiceRestoreParams, ...backup_restore_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StartUpload provides a mock function with given fields: params, opts
func (_m *TiDBCloudClient) StartUpload(params *import_service.ImportServiceStartUploadParams, opts ...import_service.ClientOption) (*import_service.ImportServiceStartUploadOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for StartUpload")
	}

	var r0 *import_service.ImportServiceStartUploadOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceStartUploadParams, ...import_service.ClientOption) (*import_service.ImportServiceStartUploadOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*import_service.ImportServiceStartUploadParams, ...import_service.ClientOption) *import_service.ImportServiceStartUploadOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*import_service.ImportServiceStartUploadOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*import_service.ImportServiceStartUploadParams, ...import_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTiDBCloudClient creates a new instance of TiDBCloudClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTiDBCloudClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *TiDBCloudClient {
	mock := &TiDBCloudClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
