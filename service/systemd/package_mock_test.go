// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/service/systemd (interfaces: DBusAPI,FileSystemOps)

// Package systemd_test is a generated GoMock package.
package systemd_test

import (
	dbus "github.com/coreos/go-systemd/v22/dbus"
	gomock "github.com/golang/mock/gomock"
	os "os"
	reflect "reflect"
)

// MockDBusAPI is a mock of DBusAPI interface
type MockDBusAPI struct {
	ctrl     *gomock.Controller
	recorder *MockDBusAPIMockRecorder
}

// MockDBusAPIMockRecorder is the mock recorder for MockDBusAPI
type MockDBusAPIMockRecorder struct {
	mock *MockDBusAPI
}

// NewMockDBusAPI creates a new mock instance
func NewMockDBusAPI(ctrl *gomock.Controller) *MockDBusAPI {
	mock := &MockDBusAPI{ctrl: ctrl}
	mock.recorder = &MockDBusAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDBusAPI) EXPECT() *MockDBusAPIMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockDBusAPI) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockDBusAPIMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDBusAPI)(nil).Close))
}

// DisableUnitFiles mocks base method
func (m *MockDBusAPI) DisableUnitFiles(arg0 []string, arg1 bool) ([]dbus.DisableUnitFileChange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableUnitFiles", arg0, arg1)
	ret0, _ := ret[0].([]dbus.DisableUnitFileChange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisableUnitFiles indicates an expected call of DisableUnitFiles
func (mr *MockDBusAPIMockRecorder) DisableUnitFiles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableUnitFiles", reflect.TypeOf((*MockDBusAPI)(nil).DisableUnitFiles), arg0, arg1)
}

// EnableUnitFiles mocks base method
func (m *MockDBusAPI) EnableUnitFiles(arg0 []string, arg1, arg2 bool) (bool, []dbus.EnableUnitFileChange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableUnitFiles", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].([]dbus.EnableUnitFileChange)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EnableUnitFiles indicates an expected call of EnableUnitFiles
func (mr *MockDBusAPIMockRecorder) EnableUnitFiles(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableUnitFiles", reflect.TypeOf((*MockDBusAPI)(nil).EnableUnitFiles), arg0, arg1, arg2)
}

// GetUnitProperties mocks base method
func (m *MockDBusAPI) GetUnitProperties(arg0 string) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitProperties", arg0)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitProperties indicates an expected call of GetUnitProperties
func (mr *MockDBusAPIMockRecorder) GetUnitProperties(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitProperties", reflect.TypeOf((*MockDBusAPI)(nil).GetUnitProperties), arg0)
}

// GetUnitTypeProperties mocks base method
func (m *MockDBusAPI) GetUnitTypeProperties(arg0, arg1 string) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitTypeProperties", arg0, arg1)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitTypeProperties indicates an expected call of GetUnitTypeProperties
func (mr *MockDBusAPIMockRecorder) GetUnitTypeProperties(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitTypeProperties", reflect.TypeOf((*MockDBusAPI)(nil).GetUnitTypeProperties), arg0, arg1)
}

// LinkUnitFiles mocks base method
func (m *MockDBusAPI) LinkUnitFiles(arg0 []string, arg1, arg2 bool) ([]dbus.LinkUnitFileChange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LinkUnitFiles", arg0, arg1, arg2)
	ret0, _ := ret[0].([]dbus.LinkUnitFileChange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LinkUnitFiles indicates an expected call of LinkUnitFiles
func (mr *MockDBusAPIMockRecorder) LinkUnitFiles(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LinkUnitFiles", reflect.TypeOf((*MockDBusAPI)(nil).LinkUnitFiles), arg0, arg1, arg2)
}

// ListUnits mocks base method
func (m *MockDBusAPI) ListUnits() ([]dbus.UnitStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUnits")
	ret0, _ := ret[0].([]dbus.UnitStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUnits indicates an expected call of ListUnits
func (mr *MockDBusAPIMockRecorder) ListUnits() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUnits", reflect.TypeOf((*MockDBusAPI)(nil).ListUnits))
}

// Reload mocks base method
func (m *MockDBusAPI) Reload() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reload")
	ret0, _ := ret[0].(error)
	return ret0
}

// Reload indicates an expected call of Reload
func (mr *MockDBusAPIMockRecorder) Reload() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reload", reflect.TypeOf((*MockDBusAPI)(nil).Reload))
}

// StartUnit mocks base method
func (m *MockDBusAPI) StartUnit(arg0, arg1 string, arg2 chan<- string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartUnit", arg0, arg1, arg2)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartUnit indicates an expected call of StartUnit
func (mr *MockDBusAPIMockRecorder) StartUnit(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartUnit", reflect.TypeOf((*MockDBusAPI)(nil).StartUnit), arg0, arg1, arg2)
}

// StopUnit mocks base method
func (m *MockDBusAPI) StopUnit(arg0, arg1 string, arg2 chan<- string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopUnit", arg0, arg1, arg2)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StopUnit indicates an expected call of StopUnit
func (mr *MockDBusAPIMockRecorder) StopUnit(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopUnit", reflect.TypeOf((*MockDBusAPI)(nil).StopUnit), arg0, arg1, arg2)
}

// MockFileSystemOps is a mock of FileSystemOps interface
type MockFileSystemOps struct {
	ctrl     *gomock.Controller
	recorder *MockFileSystemOpsMockRecorder
}

// MockFileSystemOpsMockRecorder is the mock recorder for MockFileSystemOps
type MockFileSystemOpsMockRecorder struct {
	mock *MockFileSystemOps
}

// NewMockFileSystemOps creates a new mock instance
func NewMockFileSystemOps(ctrl *gomock.Controller) *MockFileSystemOps {
	mock := &MockFileSystemOps{ctrl: ctrl}
	mock.recorder = &MockFileSystemOpsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileSystemOps) EXPECT() *MockFileSystemOpsMockRecorder {
	return m.recorder
}

// Remove mocks base method
func (m *MockFileSystemOps) Remove(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockFileSystemOpsMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockFileSystemOps)(nil).Remove), arg0)
}

// RemoveAll mocks base method
func (m *MockFileSystemOps) RemoveAll(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll
func (mr *MockFileSystemOpsMockRecorder) RemoveAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*MockFileSystemOps)(nil).RemoveAll), arg0)
}

// WriteFile mocks base method
func (m *MockFileSystemOps) WriteFile(arg0 string, arg1 []byte, arg2 os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFile indicates an expected call of WriteFile
func (mr *MockFileSystemOpsMockRecorder) WriteFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFile", reflect.TypeOf((*MockFileSystemOps)(nil).WriteFile), arg0, arg1, arg2)
}
