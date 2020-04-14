// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/provider/lxd (interfaces: LXCConfigReader)

// Package lxd is a generated GoMock package.
package lxd

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLXCConfigReader is a mock of LXCConfigReader interface
type MockLXCConfigReader struct {
	ctrl     *gomock.Controller
	recorder *MockLXCConfigReaderMockRecorder
}

// MockLXCConfigReaderMockRecorder is the mock recorder for MockLXCConfigReader
type MockLXCConfigReaderMockRecorder struct {
	mock *MockLXCConfigReader
}

// NewMockLXCConfigReader creates a new mock instance
func NewMockLXCConfigReader(ctrl *gomock.Controller) *MockLXCConfigReader {
	mock := &MockLXCConfigReader{ctrl: ctrl}
	mock.recorder = &MockLXCConfigReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLXCConfigReader) EXPECT() *MockLXCConfigReaderMockRecorder {
	return m.recorder
}

// ReadCert mocks base method
func (m *MockLXCConfigReader) ReadCert(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadCert", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadCert indicates an expected call of ReadCert
func (mr *MockLXCConfigReaderMockRecorder) ReadCert(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCert", reflect.TypeOf((*MockLXCConfigReader)(nil).ReadCert), arg0)
}

// ReadConfig mocks base method
func (m *MockLXCConfigReader) ReadConfig(arg0 string) (LXCConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadConfig", arg0)
	ret0, _ := ret[0].(LXCConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadConfig indicates an expected call of ReadConfig
func (mr *MockLXCConfigReaderMockRecorder) ReadConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadConfig", reflect.TypeOf((*MockLXCConfigReader)(nil).ReadConfig), arg0)
}
