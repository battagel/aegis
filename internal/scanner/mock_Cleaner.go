// Code generated by mockery v2.23.1. DO NOT EDIT.

package scanner

import (
	object "aegis/internal/object"

	mock "github.com/stretchr/testify/mock"
)

// MockCleaner is an autogenerated mock type for the Cleaner type
type MockCleaner struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockCleaner) Cleanup(_a0 *object.Object, _a1 bool, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(*object.Object, bool, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockCleaner interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockCleaner creates a new instance of MockCleaner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockCleaner(t mockConstructorTestingTNewMockCleaner) *MockCleaner {
	mock := &MockCleaner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
