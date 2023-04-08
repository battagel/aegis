// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ScanCollector is an autogenerated mock type for the ScanCollector type
type ScanCollector struct {
	mock.Mock
}

// CleanFile provides a mock function with given fields:
func (_m *ScanCollector) CleanFile() {
	_m.Called()
}

// FileScanned provides a mock function with given fields:
func (_m *ScanCollector) FileScanned() {
	_m.Called()
}

// InfectedFile provides a mock function with given fields:
func (_m *ScanCollector) InfectedFile() {
	_m.Called()
}

// ScanError provides a mock function with given fields:
func (_m *ScanCollector) ScanError() {
	_m.Called()
}

// ScanTime provides a mock function with given fields: _a0
func (_m *ScanCollector) ScanTime(_a0 float64) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewScanCollector interface {
	mock.TestingT
	Cleanup(func())
}

// NewScanCollector creates a new instance of ScanCollector. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewScanCollector(t mockConstructorTestingTNewScanCollector) *ScanCollector {
	mock := &ScanCollector{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
