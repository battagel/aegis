// Code generated by mockery v2.23.1. DO NOT EDIT.

package scanner

import mock "github.com/stretchr/testify/mock"

// MockAuditLogger is an autogenerated mock type for the AuditLogger type
type MockAuditLogger struct {
	mock.Mock
}

// Log provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *MockAuditLogger) Log(_a0 string, _a1 string, _a2 string, _a3 string, _a4 string, _a5 string) {
	_m.Called(_a0, _a1, _a2, _a3, _a4, _a5)
}

type mockConstructorTestingTNewMockAuditLogger interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAuditLogger creates a new instance of MockAuditLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAuditLogger(t mockConstructorTestingTNewMockAuditLogger) *MockAuditLogger {
	mock := &MockAuditLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
