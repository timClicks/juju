// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/modelgeneration (interfaces: GenerationModel)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	modelgeneration "github.com/juju/juju/apiserver/facades/client/modelgeneration"
	reflect "reflect"
)

// MockGenerationModel is a mock of GenerationModel interface
type MockGenerationModel struct {
	ctrl     *gomock.Controller
	recorder *MockGenerationModelMockRecorder
}

// MockGenerationModelMockRecorder is the mock recorder for MockGenerationModel
type MockGenerationModelMockRecorder struct {
	mock *MockGenerationModel
}

// NewMockGenerationModel creates a new mock instance
func NewMockGenerationModel(ctrl *gomock.Controller) *MockGenerationModel {
	mock := &MockGenerationModel{ctrl: ctrl}
	mock.recorder = &MockGenerationModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGenerationModel) EXPECT() *MockGenerationModelMockRecorder {
	return m.recorder
}

// AddGeneration mocks base method
func (m *MockGenerationModel) AddGeneration() error {
	ret := m.ctrl.Call(m, "AddGeneration")
	ret0, _ := ret[0].(error)
	return ret0
}

// AddGeneration indicates an expected call of AddGeneration
func (mr *MockGenerationModelMockRecorder) AddGeneration() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddGeneration", reflect.TypeOf((*MockGenerationModel)(nil).AddGeneration))
}

// HasNextGeneration mocks base method
func (m *MockGenerationModel) HasNextGeneration() (bool, error) {
	ret := m.ctrl.Call(m, "HasNextGeneration")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasNextGeneration indicates an expected call of HasNextGeneration
func (mr *MockGenerationModelMockRecorder) HasNextGeneration() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasNextGeneration", reflect.TypeOf((*MockGenerationModel)(nil).HasNextGeneration))
}

// NextGeneration mocks base method
func (m *MockGenerationModel) NextGeneration() (modelgeneration.Generation, error) {
	ret := m.ctrl.Call(m, "NextGeneration")
	ret0, _ := ret[0].(modelgeneration.Generation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NextGeneration indicates an expected call of NextGeneration
func (mr *MockGenerationModelMockRecorder) NextGeneration() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextGeneration", reflect.TypeOf((*MockGenerationModel)(nil).NextGeneration))
}
