// Code generated by MockGen. DO NOT EDIT.
// Source: psql/psql_vacancy.go

// Package psqlmock is a generated GoMock package.
package psqlmock

import (
	domain "HnH/internal/domain"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIVacancyRepository is a mock of IVacancyRepository interface.
type MockIVacancyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIVacancyRepositoryMockRecorder
}

// MockIVacancyRepositoryMockRecorder is the mock recorder for MockIVacancyRepository.
type MockIVacancyRepositoryMockRecorder struct {
	mock *MockIVacancyRepository
}

// NewMockIVacancyRepository creates a new mock instance.
func NewMockIVacancyRepository(ctrl *gomock.Controller) *MockIVacancyRepository {
	mock := &MockIVacancyRepository{ctrl: ctrl}
	mock.recorder = &MockIVacancyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIVacancyRepository) EXPECT() *MockIVacancyRepositoryMockRecorder {
	return m.recorder
}

// AddToFavourite mocks base method.
func (m *MockIVacancyRepository) AddToFavourite(ctx context.Context, userID, vacancyID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToFavourite", ctx, userID, vacancyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToFavourite indicates an expected call of AddToFavourite.
func (mr *MockIVacancyRepositoryMockRecorder) AddToFavourite(ctx, userID, vacancyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToFavourite", reflect.TypeOf((*MockIVacancyRepository)(nil).AddToFavourite), ctx, userID, vacancyID)
}

// AddVacancy mocks base method.
func (m *MockIVacancyRepository) AddVacancy(ctx context.Context, empID int, vacancy *domain.DbVacancy) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddVacancy", ctx, empID, vacancy)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddVacancy indicates an expected call of AddVacancy.
func (mr *MockIVacancyRepositoryMockRecorder) AddVacancy(ctx, empID, vacancy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVacancy", reflect.TypeOf((*MockIVacancyRepository)(nil).AddVacancy), ctx, empID, vacancy)
}

// DeleteEmpVacancy mocks base method.
func (m *MockIVacancyRepository) DeleteEmpVacancy(ctx context.Context, empID, vacancyID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEmpVacancy", ctx, empID, vacancyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEmpVacancy indicates an expected call of DeleteEmpVacancy.
func (mr *MockIVacancyRepositoryMockRecorder) DeleteEmpVacancy(ctx, empID, vacancyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEmpVacancy", reflect.TypeOf((*MockIVacancyRepository)(nil).DeleteEmpVacancy), ctx, empID, vacancyID)
}

// DeleteFromFavourite mocks base method.
func (m *MockIVacancyRepository) DeleteFromFavourite(ctx context.Context, userID, vacancyID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFromFavourite", ctx, userID, vacancyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFromFavourite indicates an expected call of DeleteFromFavourite.
func (mr *MockIVacancyRepositoryMockRecorder) DeleteFromFavourite(ctx, userID, vacancyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFromFavourite", reflect.TypeOf((*MockIVacancyRepository)(nil).DeleteFromFavourite), ctx, userID, vacancyID)
}

// GetAllVacancies mocks base method.
func (m *MockIVacancyRepository) GetAllVacancies(ctx context.Context) ([]domain.DbVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllVacancies", ctx)
	ret0, _ := ret[0].([]domain.DbVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllVacancies indicates an expected call of GetAllVacancies.
func (mr *MockIVacancyRepositoryMockRecorder) GetAllVacancies(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllVacancies", reflect.TypeOf((*MockIVacancyRepository)(nil).GetAllVacancies), ctx)
}

// GetCompanyName mocks base method.
func (m *MockIVacancyRepository) GetCompanyName(ctx context.Context, vacancyID int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyName", ctx, vacancyID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyName indicates an expected call of GetCompanyName.
func (mr *MockIVacancyRepositoryMockRecorder) GetCompanyName(ctx, vacancyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyName", reflect.TypeOf((*MockIVacancyRepository)(nil).GetCompanyName), ctx, vacancyID)
}

// GetEmpId mocks base method.
func (m *MockIVacancyRepository) GetEmpId(ctx context.Context, vacancyID int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmpId", ctx, vacancyID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmpId indicates an expected call of GetEmpId.
func (mr *MockIVacancyRepositoryMockRecorder) GetEmpId(ctx, vacancyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmpId", reflect.TypeOf((*MockIVacancyRepository)(nil).GetEmpId), ctx, vacancyID)
}

// GetEmpVacanciesByIds mocks base method.
func (m *MockIVacancyRepository) GetEmpVacanciesByIds(ctx context.Context, empID int, idList []int) ([]domain.DbVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmpVacanciesByIds", ctx, empID, idList)
	ret0, _ := ret[0].([]domain.DbVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmpVacanciesByIds indicates an expected call of GetEmpVacanciesByIds.
func (mr *MockIVacancyRepositoryMockRecorder) GetEmpVacanciesByIds(ctx, empID, idList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmpVacanciesByIds", reflect.TypeOf((*MockIVacancyRepository)(nil).GetEmpVacanciesByIds), ctx, empID, idList)
}

// GetEmployerInfo mocks base method.
func (m *MockIVacancyRepository) GetEmployerInfo(ctx context.Context, employerID int) (string, string, string, []domain.DbVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmployerInfo", ctx, employerID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].([]domain.DbVacancy)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// GetEmployerInfo indicates an expected call of GetEmployerInfo.
func (mr *MockIVacancyRepositoryMockRecorder) GetEmployerInfo(ctx, employerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmployerInfo", reflect.TypeOf((*MockIVacancyRepository)(nil).GetEmployerInfo), ctx, employerID)
}

// GetFavourite mocks base method.
func (m *MockIVacancyRepository) GetFavourite(ctx context.Context, userID int) ([]domain.DbVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavourite", ctx, userID)
	ret0, _ := ret[0].([]domain.DbVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavourite indicates an expected call of GetFavourite.
func (mr *MockIVacancyRepositoryMockRecorder) GetFavourite(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavourite", reflect.TypeOf((*MockIVacancyRepository)(nil).GetFavourite), ctx, userID)
}

// GetFavouriteFlags mocks base method.
func (m *MockIVacancyRepository) GetFavouriteFlags(ctx context.Context, userID int, vacID ...int) (map[int]bool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, userID}
	for _, a := range vacID {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFavouriteFlags", varargs...)
	ret0, _ := ret[0].(map[int]bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavouriteFlags indicates an expected call of GetFavouriteFlags.
func (mr *MockIVacancyRepositoryMockRecorder) GetFavouriteFlags(ctx, userID interface{}, vacID ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, userID}, vacID...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavouriteFlags", reflect.TypeOf((*MockIVacancyRepository)(nil).GetFavouriteFlags), varargs...)
}

// GetUserVacancies mocks base method.
func (m *MockIVacancyRepository) GetUserVacancies(ctx context.Context, userID int) ([]domain.DbVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserVacancies", ctx, userID)
	ret0, _ := ret[0].([]domain.DbVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserVacancies indicates an expected call of GetUserVacancies.
func (mr *MockIVacancyRepositoryMockRecorder) GetUserVacancies(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserVacancies", reflect.TypeOf((*MockIVacancyRepository)(nil).GetUserVacancies), ctx, userID)
}

// GetVacanciesByIds mocks base method.
func (m *MockIVacancyRepository) GetVacanciesByIds(ctx context.Context, idList []int) ([]domain.DbVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVacanciesByIds", ctx, idList)
	ret0, _ := ret[0].([]domain.DbVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVacanciesByIds indicates an expected call of GetVacanciesByIds.
func (mr *MockIVacancyRepositoryMockRecorder) GetVacanciesByIds(ctx, idList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVacanciesByIds", reflect.TypeOf((*MockIVacancyRepository)(nil).GetVacanciesByIds), ctx, idList)
}

// GetVacancy mocks base method.
func (m *MockIVacancyRepository) GetVacancy(ctx context.Context, vacancyID int) (*domain.DbVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVacancy", ctx, vacancyID)
	ret0, _ := ret[0].(*domain.DbVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVacancy indicates an expected call of GetVacancy.
func (mr *MockIVacancyRepositoryMockRecorder) GetVacancy(ctx, vacancyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVacancy", reflect.TypeOf((*MockIVacancyRepository)(nil).GetVacancy), ctx, vacancyID)
}

// UpdateEmpVacancy mocks base method.
func (m *MockIVacancyRepository) UpdateEmpVacancy(ctx context.Context, empID, vacancyID int, vacancy *domain.DbVacancy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmpVacancy", ctx, empID, vacancyID, vacancy)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmpVacancy indicates an expected call of UpdateEmpVacancy.
func (mr *MockIVacancyRepositoryMockRecorder) UpdateEmpVacancy(ctx, empID, vacancyID, vacancy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmpVacancy", reflect.TypeOf((*MockIVacancyRepository)(nil).UpdateEmpVacancy), ctx, empID, vacancyID, vacancy)
}
