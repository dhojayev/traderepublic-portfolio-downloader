// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination repository_mock.go -package=transaction
//

// Package transaction is a generated GoMock package.
package transaction

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepositoryInterface) Create(model *Model) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", model)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryInterfaceMockRecorder) Create(model any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepositoryInterface)(nil).Create), model)
}
