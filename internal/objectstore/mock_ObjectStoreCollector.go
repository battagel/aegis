// Code generated by mockery v2.23.1. DO NOT EDIT.

package objectstore

import mock "github.com/stretchr/testify/mock"

// MockObjectStoreCollector is an autogenerated mock type for the ObjectStoreCollector type
type MockObjectStoreCollector struct {
	mock.Mock
}

// GetObject provides a mock function with given fields:
func (_m *MockObjectStoreCollector) GetObject() {
	_m.Called()
}

// GetObjectTagging provides a mock function with given fields:
func (_m *MockObjectStoreCollector) GetObjectTagging() {
	_m.Called()
}

// PutObject provides a mock function with given fields:
func (_m *MockObjectStoreCollector) PutObject() {
	_m.Called()
}

// PutObjectTagging provides a mock function with given fields:
func (_m *MockObjectStoreCollector) PutObjectTagging() {
	_m.Called()
}

// RemoveObject provides a mock function with given fields:
func (_m *MockObjectStoreCollector) RemoveObject() {
	_m.Called()
}

type mockConstructorTestingTNewMockObjectStoreCollector interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockObjectStoreCollector creates a new instance of MockObjectStoreCollector. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockObjectStoreCollector(t mockConstructorTestingTNewMockObjectStoreCollector) *MockObjectStoreCollector {
	mock := &MockObjectStoreCollector{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}