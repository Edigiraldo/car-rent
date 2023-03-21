// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/ports/controllers.go

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHeathController is a mock of HeathController interface.
type MockHeathController struct {
	ctrl     *gomock.Controller
	recorder *MockHeathControllerMockRecorder
}

// MockHeathControllerMockRecorder is the mock recorder for MockHeathController.
type MockHeathControllerMockRecorder struct {
	mock *MockHeathController
}

// NewMockHeathController creates a new mock instance.
func NewMockHeathController(ctrl *gomock.Controller) *MockHeathController {
	mock := &MockHeathController{ctrl: ctrl}
	mock.recorder = &MockHeathControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHeathController) EXPECT() *MockHeathControllerMockRecorder {
	return m.recorder
}

// Pong mocks base method.
func (m *MockHeathController) Pong(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Pong", w, r)
}

// Pong indicates an expected call of Pong.
func (mr *MockHeathControllerMockRecorder) Pong(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pong", reflect.TypeOf((*MockHeathController)(nil).Pong), w, r)
}

// MockCarsController is a mock of CarsController interface.
type MockCarsController struct {
	ctrl     *gomock.Controller
	recorder *MockCarsControllerMockRecorder
}

// MockCarsControllerMockRecorder is the mock recorder for MockCarsController.
type MockCarsControllerMockRecorder struct {
	mock *MockCarsController
}

// NewMockCarsController creates a new mock instance.
func NewMockCarsController(ctrl *gomock.Controller) *MockCarsController {
	mock := &MockCarsController{ctrl: ctrl}
	mock.recorder = &MockCarsControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCarsController) EXPECT() *MockCarsControllerMockRecorder {
	return m.recorder
}

// FullUpdate mocks base method.
func (m *MockCarsController) FullUpdate(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FullUpdate", w, r)
}

// FullUpdate indicates an expected call of FullUpdate.
func (mr *MockCarsControllerMockRecorder) FullUpdate(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullUpdate", reflect.TypeOf((*MockCarsController)(nil).FullUpdate), w, r)
}

// Get mocks base method.
func (m *MockCarsController) Get(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Get", w, r)
}

// Get indicates an expected call of Get.
func (mr *MockCarsControllerMockRecorder) Get(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCarsController)(nil).Get), w, r)
}

// Register mocks base method.
func (m *MockCarsController) Register(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", w, r)
}

// Register indicates an expected call of Register.
func (mr *MockCarsControllerMockRecorder) Register(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockCarsController)(nil).Register), w, r)
}
