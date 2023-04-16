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

// Delete mocks base method.
func (m *MockCarsController) Delete(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", w, r)
}

// Delete indicates an expected call of Delete.
func (mr *MockCarsControllerMockRecorder) Delete(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCarsController)(nil).Delete), w, r)
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

// List mocks base method.
func (m *MockCarsController) List(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "List", w, r)
}

// List indicates an expected call of List.
func (mr *MockCarsControllerMockRecorder) List(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCarsController)(nil).List), w, r)
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

// MockUsersController is a mock of UsersController interface.
type MockUsersController struct {
	ctrl     *gomock.Controller
	recorder *MockUsersControllerMockRecorder
}

// MockUsersControllerMockRecorder is the mock recorder for MockUsersController.
type MockUsersControllerMockRecorder struct {
	mock *MockUsersController
}

// NewMockUsersController creates a new mock instance.
func NewMockUsersController(ctrl *gomock.Controller) *MockUsersController {
	mock := &MockUsersController{ctrl: ctrl}
	mock.recorder = &MockUsersControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsersController) EXPECT() *MockUsersControllerMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUsersController) Delete(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", w, r)
}

// Delete indicates an expected call of Delete.
func (mr *MockUsersControllerMockRecorder) Delete(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUsersController)(nil).Delete), w, r)
}

// FullUpdate mocks base method.
func (m *MockUsersController) FullUpdate(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FullUpdate", w, r)
}

// FullUpdate indicates an expected call of FullUpdate.
func (mr *MockUsersControllerMockRecorder) FullUpdate(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullUpdate", reflect.TypeOf((*MockUsersController)(nil).FullUpdate), w, r)
}

// Get mocks base method.
func (m *MockUsersController) Get(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Get", w, r)
}

// Get indicates an expected call of Get.
func (mr *MockUsersControllerMockRecorder) Get(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUsersController)(nil).Get), w, r)
}

// SignUp mocks base method.
func (m *MockUsersController) SignUp(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SignUp", w, r)
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUsersControllerMockRecorder) SignUp(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUsersController)(nil).SignUp), w, r)
}

// MockCitiesController is a mock of CitiesController interface.
type MockCitiesController struct {
	ctrl     *gomock.Controller
	recorder *MockCitiesControllerMockRecorder
}

// MockCitiesControllerMockRecorder is the mock recorder for MockCitiesController.
type MockCitiesControllerMockRecorder struct {
	mock *MockCitiesController
}

// NewMockCitiesController creates a new mock instance.
func NewMockCitiesController(ctrl *gomock.Controller) *MockCitiesController {
	mock := &MockCitiesController{ctrl: ctrl}
	mock.recorder = &MockCitiesControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCitiesController) EXPECT() *MockCitiesControllerMockRecorder {
	return m.recorder
}

// ListNames mocks base method.
func (m *MockCitiesController) ListNames(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListNames", w, r)
}

// ListNames indicates an expected call of ListNames.
func (mr *MockCitiesControllerMockRecorder) ListNames(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNames", reflect.TypeOf((*MockCitiesController)(nil).ListNames), w, r)
}
