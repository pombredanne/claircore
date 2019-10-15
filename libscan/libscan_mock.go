// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quay/claircore/libscan (interfaces: Libscan)

// Package libscan is a generated GoMock package.
package libscan

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	claircore "github.com/quay/claircore"
	reflect "reflect"
)

// MockLibscan is a mock of Libscan interface
type MockLibscan struct {
	ctrl     *gomock.Controller
	recorder *MockLibscanMockRecorder
}

// MockLibscanMockRecorder is the mock recorder for MockLibscan
type MockLibscanMockRecorder struct {
	mock *MockLibscan
}

// NewMockLibscan creates a new mock instance
func NewMockLibscan(ctrl *gomock.Controller) *MockLibscan {
	mock := &MockLibscan{ctrl: ctrl}
	mock.recorder = &MockLibscanMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLibscan) EXPECT() *MockLibscanMockRecorder {
	return m.recorder
}

// Scan mocks base method
func (m *MockLibscan) Scan(arg0 context.Context, arg1 *claircore.Manifest) (<-chan *claircore.ScanReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", arg0, arg1)
	ret0, _ := ret[0].(<-chan *claircore.ScanReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scan indicates an expected call of Scan
func (mr *MockLibscanMockRecorder) Scan(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockLibscan)(nil).Scan), arg0, arg1)
}

// ScanReport mocks base method
func (m *MockLibscan) ScanReport(arg0 context.Context, arg1 string) (*claircore.ScanReport, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScanReport", arg0, arg1)
	ret0, _ := ret[0].(*claircore.ScanReport)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ScanReport indicates an expected call of ScanReport
func (mr *MockLibscanMockRecorder) ScanReport(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScanReport", reflect.TypeOf((*MockLibscan)(nil).ScanReport), arg0, arg1)
}
