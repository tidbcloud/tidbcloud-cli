// Code generated by mockery v2.43.1. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// EventsSender is an autogenerated mock type for the EventsSender type
type EventsSender struct {
	mock.Mock
}

// SendEvent provides a mock function with given fields: event
func (_m *EventsSender) SendEvent(event interface{}) error {
	ret := _m.Called(event)

	if len(ret) == 0 {
		panic("no return value specified for SendEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEventsSender creates a new instance of EventsSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventsSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventsSender {
	mock := &EventsSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
