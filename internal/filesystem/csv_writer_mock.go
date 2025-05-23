// Code generated by MockGen. DO NOT EDIT.
// Source: csv_writer.go
//
// Generated by this command:
//
//	mockgen -source=csv_writer.go -destination csv_writer_mock.go -package=filesystem
//

// Package filesystem is a generated GoMock package.
package filesystem

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCSVWriterInterface is a mock of CSVWriterInterface interface.
type MockCSVWriterInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCSVWriterInterfaceMockRecorder
	isgomock struct{}
}

// MockCSVWriterInterfaceMockRecorder is the mock recorder for MockCSVWriterInterface.
type MockCSVWriterInterfaceMockRecorder struct {
	mock *MockCSVWriterInterface
}

// NewMockCSVWriterInterface creates a new mock instance.
func NewMockCSVWriterInterface(ctrl *gomock.Controller) *MockCSVWriterInterface {
	mock := &MockCSVWriterInterface{ctrl: ctrl}
	mock.recorder = &MockCSVWriterInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCSVWriterInterface) EXPECT() *MockCSVWriterInterfaceMockRecorder {
	return m.recorder
}

// Write mocks base method.
func (m *MockCSVWriterInterface) Write(filepath string, entry CSVEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", filepath, entry)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write.
func (mr *MockCSVWriterInterfaceMockRecorder) Write(filepath, entry any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockCSVWriterInterface)(nil).Write), filepath, entry)
}
