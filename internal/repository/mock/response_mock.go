// Code generated by MockGen. DO NOT EDIT.
// Source: psql/psql_response.go

// Package psqlmock is a generated GoMock package.
package psqlmock

import (
	domain "HnH/internal/domain"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIResponseRepository is a mock of IResponseRepository interface.
type MockIResponseRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIResponseRepositoryMockRecorder
}

// MockIResponseRepositoryMockRecorder is the mock recorder for MockIResponseRepository.
type MockIResponseRepositoryMockRecorder struct {
	mock *MockIResponseRepository
}

// NewMockIResponseRepository creates a new mock instance.
func NewMockIResponseRepository(ctrl *gomock.Controller) *MockIResponseRepository {
	mock := &MockIResponseRepository{ctrl: ctrl}
	mock.recorder = &MockIResponseRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIResponseRepository) EXPECT() *MockIResponseRepositoryMockRecorder {
	return m.recorder
}

// GetAttachedCVs mocks base method.
func (m *MockIResponseRepository) GetAttachedCVs(ctx context.Context, vacancyID int) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttachedCVs", ctx, vacancyID)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttachedCVs indicates an expected call of GetAttachedCVs.
func (mr *MockIResponseRepositoryMockRecorder) GetAttachedCVs(ctx, vacancyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttachedCVs", reflect.TypeOf((*MockIResponseRepository)(nil).GetAttachedCVs), ctx, vacancyID)
}

// GetUserResponses mocks base method.
func (m *MockIResponseRepository) GetUserResponses(ctx context.Context, userID int) ([]domain.ApiResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserResponses", ctx, userID)
	ret0, _ := ret[0].([]domain.ApiResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserResponses indicates an expected call of GetUserResponses.
func (mr *MockIResponseRepositoryMockRecorder) GetUserResponses(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserResponses", reflect.TypeOf((*MockIResponseRepository)(nil).GetUserResponses), ctx, userID)
}

// GetVacanciesIdsByCVId mocks base method.
func (m *MockIResponseRepository) GetVacanciesIdsByCVId(ctx context.Context, cvID int) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVacanciesIdsByCVId", ctx, cvID)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVacanciesIdsByCVId indicates an expected call of GetVacanciesIdsByCVId.
func (mr *MockIResponseRepositoryMockRecorder) GetVacanciesIdsByCVId(ctx, cvID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVacanciesIdsByCVId", reflect.TypeOf((*MockIResponseRepository)(nil).GetVacanciesIdsByCVId), ctx, cvID)
}

// RespondToVacancy mocks base method.
func (m *MockIResponseRepository) RespondToVacancy(ctx context.Context, vacancyID, cvID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RespondToVacancy", ctx, vacancyID, cvID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RespondToVacancy indicates an expected call of RespondToVacancy.
func (mr *MockIResponseRepositoryMockRecorder) RespondToVacancy(ctx, vacancyID, cvID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RespondToVacancy", reflect.TypeOf((*MockIResponseRepository)(nil).RespondToVacancy), ctx, vacancyID, cvID)
}
