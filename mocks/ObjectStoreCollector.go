// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ObjectStoreCollector is an autogenerated mock type for the ObjectStoreCollector type
type ObjectStoreCollector struct {
	mock.Mock
}

// GetObject provides a mock function with given fields:
func (_m *ObjectStoreCollector) GetObject() {
	_m.Called()
}

// GetObjectTagging provides a mock function with given fields:
func (_m *ObjectStoreCollector) GetObjectTagging() {
	_m.Called()
}

// PutObject provides a mock function with given fields:
func (_m *ObjectStoreCollector) PutObject() {
	_m.Called()
}

// PutObjectTagging provides a mock function with given fields:
func (_m *ObjectStoreCollector) PutObjectTagging() {
	_m.Called()
}

// RemoveObject provides a mock function with given fields:
func (_m *ObjectStoreCollector) RemoveObject() {
	_m.Called()
}

type mockConstructorTestingTNewObjectStoreCollector interface {
	mock.TestingT
	Cleanup(func())
}

// NewObjectStoreCollector creates a new instance of ObjectStoreCollector. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewObjectStoreCollector(t mockConstructorTestingTNewObjectStoreCollector) *ObjectStoreCollector {
	mock := &ObjectStoreCollector{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
