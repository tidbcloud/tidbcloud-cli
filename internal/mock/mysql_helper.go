// Code generated by mockery v2.27.1. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MySQLHelper is an autogenerated mock type for the MySQLHelper type
type MySQLHelper struct {
	mock.Mock
}

// CheckMySQLClient provides a mock function with given fields:
func (_m *MySQLHelper) CheckMySQLClient() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DownloadCaFile provides a mock function with given fields: caFile
func (_m *MySQLHelper) DownloadCaFile(caFile string) error {
	ret := _m.Called(caFile)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(caFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DumpFromMySQL provides a mock function with given fields: ctx, command, sqlCacheFile
func (_m *MySQLHelper) DumpFromMySQL(ctx context.Context, command []string, sqlCacheFile string) error {
	ret := _m.Called(ctx, command, sqlCacheFile)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, string) error); ok {
		r0 = rf(ctx, command, sqlCacheFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateSqlCachePath provides a mock function with given fields:
func (_m *MySQLHelper) GenerateSqlCachePath() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ImportToServerless provides a mock function with given fields: ctx, command, sqlCacheFile
func (_m *MySQLHelper) ImportToServerless(ctx context.Context, command []string, sqlCacheFile string) error {
	ret := _m.Called(ctx, command, sqlCacheFile)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, string) error); ok {
		r0 = rf(ctx, command, sqlCacheFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMySQLHelper interface {
	mock.TestingT
	Cleanup(func())
}

// NewMySQLHelper creates a new instance of MySQLHelper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMySQLHelper(t mockConstructorTestingTNewMySQLHelper) *MySQLHelper {
	mock := &MySQLHelper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
