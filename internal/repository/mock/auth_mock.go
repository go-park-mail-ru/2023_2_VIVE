// Code generated by MockGen. DO NOT EDIT.
// Source: grpc/auth.go

// Package psqlmock is a generated GoMock package.
package psqlmock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIAuthRepository is a mock of IAuthRepository interface.
type MockIAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAuthRepositoryMockRecorder
}

// MockIAuthRepositoryMockRecorder is the mock recorder for MockIAuthRepository.
type MockIAuthRepositoryMockRecorder struct {
	mock *MockIAuthRepository
}

// NewMockIAuthRepository creates a new mock instance.
func NewMockIAuthRepository(ctrl *gomock.Controller) *MockIAuthRepository {
	mock := &MockIAuthRepository{ctrl: ctrl}
	mock.recorder = &MockIAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAuthRepository) EXPECT() *MockIAuthRepositoryMockRecorder {
	return m.recorder
}

// AddSession mocks base method.
func (m *MockIAuthRepository) AddSession(ctx context.Context, sessionID string, userID int, expiryUnixSeconds int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSession", ctx, sessionID, userID, expiryUnixSeconds)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSession indicates an expected call of AddSession.
func (mr *MockIAuthRepositoryMockRecorder) AddSession(ctx, sessionID, userID, expiryUnixSeconds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSession", reflect.TypeOf((*MockIAuthRepository)(nil).AddSession), ctx, sessionID, userID, expiryUnixSeconds)
}

// DeleteSession mocks base method.
func (m *MockIAuthRepository) DeleteSession(ctx context.Context, sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", ctx, sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockIAuthRepositoryMockRecorder) DeleteSession(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockIAuthRepository)(nil).DeleteSession), ctx, sessionID)
}

// GetUserIdBySession mocks base method.
func (m *MockIAuthRepository) GetUserIdBySession(ctx context.Context, sessionID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIdBySession", ctx, sessionID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIdBySession indicates an expected call of GetUserIdBySession.
func (mr *MockIAuthRepositoryMockRecorder) GetUserIdBySession(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIdBySession", reflect.TypeOf((*MockIAuthRepository)(nil).GetUserIdBySession), ctx, sessionID)
}

// ValidateSession mocks base method.
func (m *MockIAuthRepository) ValidateSession(ctx context.Context, sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateSession", ctx, sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateSession indicates an expected call of ValidateSession.
func (mr *MockIAuthRepositoryMockRecorder) ValidateSession(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSession", reflect.TypeOf((*MockIAuthRepository)(nil).ValidateSession), ctx, sessionID)
}
