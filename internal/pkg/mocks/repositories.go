// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/ports/repositories.go

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	domain "github.com/Edigiraldo/car-rent/internal/core/domain"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockDatabase) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockDatabaseMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDatabase)(nil).Close))
}

// GetDBHandle mocks base method.
func (m *MockDatabase) GetDBHandle() *sql.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDBHandle")
	ret0, _ := ret[0].(*sql.DB)
	return ret0
}

// GetDBHandle indicates an expected call of GetDBHandle.
func (mr *MockDatabaseMockRecorder) GetDBHandle() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDBHandle", reflect.TypeOf((*MockDatabase)(nil).GetDBHandle))
}

// MockCarsRepo is a mock of CarsRepo interface.
type MockCarsRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCarsRepoMockRecorder
}

// MockCarsRepoMockRecorder is the mock recorder for MockCarsRepo.
type MockCarsRepoMockRecorder struct {
	mock *MockCarsRepo
}

// NewMockCarsRepo creates a new mock instance.
func NewMockCarsRepo(ctrl *gomock.Controller) *MockCarsRepo {
	mock := &MockCarsRepo{ctrl: ctrl}
	mock.recorder = &MockCarsRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCarsRepo) EXPECT() *MockCarsRepoMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockCarsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCarsRepoMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCarsRepo)(nil).Delete), ctx, id)
}

// FullUpdate mocks base method.
func (m *MockCarsRepo) FullUpdate(ctx context.Context, dc domain.Car) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullUpdate", ctx, dc)
	ret0, _ := ret[0].(error)
	return ret0
}

// FullUpdate indicates an expected call of FullUpdate.
func (mr *MockCarsRepoMockRecorder) FullUpdate(ctx, dc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullUpdate", reflect.TypeOf((*MockCarsRepo)(nil).FullUpdate), ctx, dc)
}

// Get mocks base method.
func (m *MockCarsRepo) Get(ctx context.Context, ID uuid.UUID) (domain.Car, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, ID)
	ret0, _ := ret[0].(domain.Car)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCarsRepoMockRecorder) Get(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCarsRepo)(nil).Get), ctx, ID)
}

// Insert mocks base method.
func (m *MockCarsRepo) Insert(ctx context.Context, dc domain.Car) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, dc)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockCarsRepoMockRecorder) Insert(ctx, dc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockCarsRepo)(nil).Insert), ctx, dc)
}

// List mocks base method.
func (m *MockCarsRepo) List(ctx context.Context, cityName, from_car_id string, limit uint16) ([]domain.Car, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, cityName, from_car_id, limit)
	ret0, _ := ret[0].([]domain.Car)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCarsRepoMockRecorder) List(ctx, cityName, from_car_id, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCarsRepo)(nil).List), ctx, cityName, from_car_id, limit)
}

// MockUsersRepo is a mock of UsersRepo interface.
type MockUsersRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUsersRepoMockRecorder
}

// MockUsersRepoMockRecorder is the mock recorder for MockUsersRepo.
type MockUsersRepoMockRecorder struct {
	mock *MockUsersRepo
}

// NewMockUsersRepo creates a new mock instance.
func NewMockUsersRepo(ctrl *gomock.Controller) *MockUsersRepo {
	mock := &MockUsersRepo{ctrl: ctrl}
	mock.recorder = &MockUsersRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsersRepo) EXPECT() *MockUsersRepoMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUsersRepo) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUsersRepoMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUsersRepo)(nil).Delete), ctx, id)
}

// FullUpdate mocks base method.
func (m *MockUsersRepo) FullUpdate(ctx context.Context, du domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullUpdate", ctx, du)
	ret0, _ := ret[0].(error)
	return ret0
}

// FullUpdate indicates an expected call of FullUpdate.
func (mr *MockUsersRepoMockRecorder) FullUpdate(ctx, du interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullUpdate", reflect.TypeOf((*MockUsersRepo)(nil).FullUpdate), ctx, du)
}

// Get mocks base method.
func (m *MockUsersRepo) Get(ctx context.Context, ID uuid.UUID) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, ID)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUsersRepoMockRecorder) Get(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUsersRepo)(nil).Get), ctx, ID)
}

// Insert mocks base method.
func (m *MockUsersRepo) Insert(ctx context.Context, du domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, du)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockUsersRepoMockRecorder) Insert(ctx, du interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUsersRepo)(nil).Insert), ctx, du)
}

// MockCitiesRepo is a mock of CitiesRepo interface.
type MockCitiesRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCitiesRepoMockRecorder
}

// MockCitiesRepoMockRecorder is the mock recorder for MockCitiesRepo.
type MockCitiesRepoMockRecorder struct {
	mock *MockCitiesRepo
}

// NewMockCitiesRepo creates a new mock instance.
func NewMockCitiesRepo(ctrl *gomock.Controller) *MockCitiesRepo {
	mock := &MockCitiesRepo{ctrl: ctrl}
	mock.recorder = &MockCitiesRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCitiesRepo) EXPECT() *MockCitiesRepoMockRecorder {
	return m.recorder
}

// GetIdByName mocks base method.
func (m *MockCitiesRepo) GetIdByName(ctx context.Context, name string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIdByName", ctx, name)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIdByName indicates an expected call of GetIdByName.
func (mr *MockCitiesRepoMockRecorder) GetIdByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIdByName", reflect.TypeOf((*MockCitiesRepo)(nil).GetIdByName), ctx, name)
}

// GetNameByID mocks base method.
func (m *MockCitiesRepo) GetNameByID(ctx context.Context, ID uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNameByID", ctx, ID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNameByID indicates an expected call of GetNameByID.
func (mr *MockCitiesRepoMockRecorder) GetNameByID(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNameByID", reflect.TypeOf((*MockCitiesRepo)(nil).GetNameByID), ctx, ID)
}

// ListNames mocks base method.
func (m *MockCitiesRepo) ListNames(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNames", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNames indicates an expected call of ListNames.
func (mr *MockCitiesRepoMockRecorder) ListNames(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNames", reflect.TypeOf((*MockCitiesRepo)(nil).ListNames), ctx)
}

// MockReservationsRepo is a mock of ReservationsRepo interface.
type MockReservationsRepo struct {
	ctrl     *gomock.Controller
	recorder *MockReservationsRepoMockRecorder
}

// MockReservationsRepoMockRecorder is the mock recorder for MockReservationsRepo.
type MockReservationsRepoMockRecorder struct {
	mock *MockReservationsRepo
}

// NewMockReservationsRepo creates a new mock instance.
func NewMockReservationsRepo(ctrl *gomock.Controller) *MockReservationsRepo {
	mock := &MockReservationsRepo{ctrl: ctrl}
	mock.recorder = &MockReservationsRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReservationsRepo) EXPECT() *MockReservationsRepoMockRecorder {
	return m.recorder
}

// Insert mocks base method.
func (m *MockReservationsRepo) Insert(ctx context.Context, dr domain.Reservation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, dr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockReservationsRepoMockRecorder) Insert(ctx, dr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockReservationsRepo)(nil).Insert), ctx, dr)
}
