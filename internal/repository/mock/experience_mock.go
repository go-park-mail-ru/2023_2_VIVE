// Code generated by MockGen. DO NOT EDIT.
// Source: psql/psql_experience.go

// Package psqlmock is a generated GoMock package.
package psqlmock

import (
	domain "HnH/internal/domain"
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIExperienceRepository is a mock of IExperienceRepository interface.
type MockIExperienceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIExperienceRepositoryMockRecorder
}

// MockIExperienceRepositoryMockRecorder is the mock recorder for MockIExperienceRepository.
type MockIExperienceRepositoryMockRecorder struct {
	mock *MockIExperienceRepository
}

// NewMockIExperienceRepository creates a new mock instance.
func NewMockIExperienceRepository(ctrl *gomock.Controller) *MockIExperienceRepository {
	mock := &MockIExperienceRepository{ctrl: ctrl}
	mock.recorder = &MockIExperienceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIExperienceRepository) EXPECT() *MockIExperienceRepositoryMockRecorder {
	return m.recorder
}

// AddExperience mocks base method.
func (m *MockIExperienceRepository) AddExperience(ctx context.Context, cvID int, experience domain.DbExperience) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddExperience", ctx, cvID, experience)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddExperience indicates an expected call of AddExperience.
func (mr *MockIExperienceRepositoryMockRecorder) AddExperience(ctx, cvID, experience interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddExperience", reflect.TypeOf((*MockIExperienceRepository)(nil).AddExperience), ctx, cvID, experience)
}

// AddTxExperiences mocks base method.
func (m *MockIExperienceRepository) AddTxExperiences(ctx context.Context, tx *sql.Tx, cvID int, experiences []domain.DbExperience) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTxExperiences", ctx, tx, cvID, experiences)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTxExperiences indicates an expected call of AddTxExperiences.
func (mr *MockIExperienceRepositoryMockRecorder) AddTxExperiences(ctx, tx, cvID, experiences interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTxExperiences", reflect.TypeOf((*MockIExperienceRepository)(nil).AddTxExperiences), ctx, tx, cvID, experiences)
}

// DeleteTxExperiences mocks base method.
func (m *MockIExperienceRepository) DeleteTxExperiences(ctx context.Context, tx *sql.Tx, cvID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTxExperiences", ctx, tx, cvID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTxExperiences indicates an expected call of DeleteTxExperiences.
func (mr *MockIExperienceRepositoryMockRecorder) DeleteTxExperiences(ctx, tx, cvID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTxExperiences", reflect.TypeOf((*MockIExperienceRepository)(nil).DeleteTxExperiences), ctx, tx, cvID)
}

// DeleteTxExperiencesByIDs mocks base method.
func (m *MockIExperienceRepository) DeleteTxExperiencesByIDs(ctx context.Context, tx *sql.Tx, expIds []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTxExperiencesByIDs", ctx, tx, expIds)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTxExperiencesByIDs indicates an expected call of DeleteTxExperiencesByIDs.
func (mr *MockIExperienceRepositoryMockRecorder) DeleteTxExperiencesByIDs(ctx, tx, expIds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTxExperiencesByIDs", reflect.TypeOf((*MockIExperienceRepository)(nil).DeleteTxExperiencesByIDs), ctx, tx, expIds)
}

// GetCVExperiencesIDs mocks base method.
func (m *MockIExperienceRepository) GetCVExperiencesIDs(ctx context.Context, cvID int) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCVExperiencesIDs", ctx, cvID)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCVExperiencesIDs indicates an expected call of GetCVExperiencesIDs.
func (mr *MockIExperienceRepositoryMockRecorder) GetCVExperiencesIDs(ctx, cvID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCVExperiencesIDs", reflect.TypeOf((*MockIExperienceRepository)(nil).GetCVExperiencesIDs), ctx, cvID)
}

// GetTxExperiences mocks base method.
func (m *MockIExperienceRepository) GetTxExperiences(ctx context.Context, tx *sql.Tx, cvID int) ([]domain.DbExperience, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxExperiences", ctx, tx, cvID)
	ret0, _ := ret[0].([]domain.DbExperience)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxExperiences indicates an expected call of GetTxExperiences.
func (mr *MockIExperienceRepositoryMockRecorder) GetTxExperiences(ctx, tx, cvID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxExperiences", reflect.TypeOf((*MockIExperienceRepository)(nil).GetTxExperiences), ctx, tx, cvID)
}

// GetTxExperiencesByIds mocks base method.
func (m *MockIExperienceRepository) GetTxExperiencesByIds(ctx context.Context, tx *sql.Tx, cvIDs []int) ([]domain.DbExperience, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxExperiencesByIds", ctx, tx, cvIDs)
	ret0, _ := ret[0].([]domain.DbExperience)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTxExperiencesByIds indicates an expected call of GetTxExperiencesByIds.
func (mr *MockIExperienceRepositoryMockRecorder) GetTxExperiencesByIds(ctx, tx, cvIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxExperiencesByIds", reflect.TypeOf((*MockIExperienceRepository)(nil).GetTxExperiencesByIds), ctx, tx, cvIDs)
}

// UpdateTxExperiences mocks base method.
func (m *MockIExperienceRepository) UpdateTxExperiences(ctx context.Context, tx *sql.Tx, cvID int, experiences []domain.DbExperience) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTxExperiences", ctx, tx, cvID, experiences)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTxExperiences indicates an expected call of UpdateTxExperiences.
func (mr *MockIExperienceRepositoryMockRecorder) UpdateTxExperiences(ctx, tx, cvID, experiences interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTxExperiences", reflect.TypeOf((*MockIExperienceRepository)(nil).UpdateTxExperiences), ctx, tx, cvID, experiences)
}
