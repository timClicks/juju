// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/api/instancemutater (interfaces: MutaterMachine)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	instancemutater "github.com/juju/juju/api/instancemutater"
	watcher "github.com/juju/juju/core/watcher"
	names_v2 "gopkg.in/juju/names.v2"
	reflect "reflect"
)

// MockMutaterMachine is a mock of MutaterMachine interface
type MockMutaterMachine struct {
	ctrl     *gomock.Controller
	recorder *MockMutaterMachineMockRecorder
}

// MockMutaterMachineMockRecorder is the mock recorder for MockMutaterMachine
type MockMutaterMachineMockRecorder struct {
	mock *MockMutaterMachine
}

// NewMockMutaterMachine creates a new mock instance
func NewMockMutaterMachine(ctrl *gomock.Controller) *MockMutaterMachine {
	mock := &MockMutaterMachine{ctrl: ctrl}
	mock.recorder = &MockMutaterMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMutaterMachine) EXPECT() *MockMutaterMachineMockRecorder {
	return m.recorder
}

// CharmProfilingInfo mocks base method
func (m *MockMutaterMachine) CharmProfilingInfo(arg0 []string) (*instancemutater.ProfileInfo, error) {
	ret := m.ctrl.Call(m, "CharmProfilingInfo", arg0)
	ret0, _ := ret[0].(*instancemutater.ProfileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CharmProfilingInfo indicates an expected call of CharmProfilingInfo
func (mr *MockMutaterMachineMockRecorder) CharmProfilingInfo(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmProfilingInfo", reflect.TypeOf((*MockMutaterMachine)(nil).CharmProfilingInfo), arg0)
}

// RemoveUpgradeCharmProfileData mocks base method
func (m *MockMutaterMachine) RemoveUpgradeCharmProfileData(arg0 string) error {
	ret := m.ctrl.Call(m, "RemoveUpgradeCharmProfileData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUpgradeCharmProfileData indicates an expected call of RemoveUpgradeCharmProfileData
func (mr *MockMutaterMachineMockRecorder) RemoveUpgradeCharmProfileData(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUpgradeCharmProfileData", reflect.TypeOf((*MockMutaterMachine)(nil).RemoveUpgradeCharmProfileData), arg0)
}

// SetCharmProfiles mocks base method
func (m *MockMutaterMachine) SetCharmProfiles(arg0 []string) error {
	ret := m.ctrl.Call(m, "SetCharmProfiles", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCharmProfiles indicates an expected call of SetCharmProfiles
func (mr *MockMutaterMachineMockRecorder) SetCharmProfiles(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCharmProfiles", reflect.TypeOf((*MockMutaterMachine)(nil).SetCharmProfiles), arg0)
}

// SetUpgradeCharmProfileComplete mocks base method
func (m *MockMutaterMachine) SetUpgradeCharmProfileComplete(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "SetUpgradeCharmProfileComplete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUpgradeCharmProfileComplete indicates an expected call of SetUpgradeCharmProfileComplete
func (mr *MockMutaterMachineMockRecorder) SetUpgradeCharmProfileComplete(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUpgradeCharmProfileComplete", reflect.TypeOf((*MockMutaterMachine)(nil).SetUpgradeCharmProfileComplete), arg0, arg1)
}

// Tag mocks base method
func (m *MockMutaterMachine) Tag() names_v2.MachineTag {
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names_v2.MachineTag)
	return ret0
}

// Tag indicates an expected call of Tag
func (mr *MockMutaterMachineMockRecorder) Tag() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockMutaterMachine)(nil).Tag))
}

// WatchUnits mocks base method
func (m *MockMutaterMachine) WatchUnits() (watcher.StringsWatcher, error) {
	ret := m.ctrl.Call(m, "WatchUnits")
	ret0, _ := ret[0].(watcher.StringsWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchUnits indicates an expected call of WatchUnits
func (mr *MockMutaterMachineMockRecorder) WatchUnits() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchUnits", reflect.TypeOf((*MockMutaterMachine)(nil).WatchUnits))
}
