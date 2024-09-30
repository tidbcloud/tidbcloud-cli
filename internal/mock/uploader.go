// Code generated by mockery v2.43.0. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	s3 "github.com/tidbcloud/tidbcloud-cli/internal/service/aws/s3"
)

// Uploader is an autogenerated mock type for the Uploader type
type Uploader struct {
	mock.Mock
}

// SetConcurrency provides a mock function with given fields: concurrency
func (_m *Uploader) SetConcurrency(concurrency int) error {
	ret := _m.Called(concurrency)

	if len(ret) == 0 {
		panic("no return value specified for SetConcurrency")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(concurrency)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Upload provides a mock function with given fields: ctx, input
func (_m *Uploader) Upload(ctx context.Context, input *s3.PutObjectInput) (string, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *s3.PutObjectInput) (string, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *s3.PutObjectInput) string); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *s3.PutObjectInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUploader creates a new instance of Uploader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUploader(t interface {
	mock.TestingT
	Cleanup(func())
}) *Uploader {
	mock := &Uploader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
