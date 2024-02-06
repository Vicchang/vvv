// Code generated by MockGen. DO NOT EDIT.
// Source: pod.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPodService is a mock of PodService interface.
type MockPodService struct {
	ctrl     *gomock.Controller
	recorder *MockPodServiceMockRecorder
}

// MockPodServiceMockRecorder is the mock recorder for MockPodService.
type MockPodServiceMockRecorder struct {
	mock *MockPodService
}

// NewMockPodService creates a new mock instance.
func NewMockPodService(ctrl *gomock.Controller) *MockPodService {
	mock := &MockPodService{ctrl: ctrl}
	mock.recorder = &MockPodServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPodService) EXPECT() *MockPodServiceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockPodService) Add(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add.
func (mr *MockPodServiceMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockPodService)(nil).Add), arg0)
}
