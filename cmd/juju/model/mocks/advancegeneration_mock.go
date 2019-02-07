// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/cmd/juju/model (interfaces: AdvanceGenerationCommandAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAdvanceGenerationCommandAPI is a mock of AdvanceGenerationCommandAPI interface
type MockAdvanceGenerationCommandAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAdvanceGenerationCommandAPIMockRecorder
}

// MockAdvanceGenerationCommandAPIMockRecorder is the mock recorder for MockAdvanceGenerationCommandAPI
type MockAdvanceGenerationCommandAPIMockRecorder struct {
	mock *MockAdvanceGenerationCommandAPI
}

// NewMockAdvanceGenerationCommandAPI creates a new mock instance
func NewMockAdvanceGenerationCommandAPI(ctrl *gomock.Controller) *MockAdvanceGenerationCommandAPI {
	mock := &MockAdvanceGenerationCommandAPI{ctrl: ctrl}
	mock.recorder = &MockAdvanceGenerationCommandAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAdvanceGenerationCommandAPI) EXPECT() *MockAdvanceGenerationCommandAPIMockRecorder {
	return m.recorder
}

// AdvanceGeneration mocks base method
func (m *MockAdvanceGenerationCommandAPI) AdvanceGeneration(arg0 string, arg1 []string) error {
	ret := m.ctrl.Call(m, "AdvanceGeneration", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AdvanceGeneration indicates an expected call of AdvanceGeneration
func (mr *MockAdvanceGenerationCommandAPIMockRecorder) AdvanceGeneration(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdvanceGeneration", reflect.TypeOf((*MockAdvanceGenerationCommandAPI)(nil).AdvanceGeneration), arg0, arg1)
}

// Close mocks base method
func (m *MockAdvanceGenerationCommandAPI) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockAdvanceGenerationCommandAPIMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAdvanceGenerationCommandAPI)(nil).Close))
}
