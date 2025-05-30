// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/user/user.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/user/user.go
//

// Package mock_userRepository is a generated GoMock package.
package mock_userRepository

import (
	reflect "reflect"

	postgres "github.com/go-jet/jet/v2/postgres"
	uuid "github.com/google/uuid"
	model "github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	gomock "go.uber.org/mock/gomock"
)

// MockIUserRepository is a mock of IUserRepository interface.
type MockIUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIUserRepositoryMockRecorder
	isgomock struct{}
}

// MockIUserRepositoryMockRecorder is the mock recorder for MockIUserRepository.
type MockIUserRepositoryMockRecorder struct {
	mock *MockIUserRepository
}

// NewMockIUserRepository creates a new mock instance.
func NewMockIUserRepository(ctrl *gomock.Controller) *MockIUserRepository {
	mock := &MockIUserRepository{ctrl: ctrl}
	mock.recorder = &MockIUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserRepository) EXPECT() *MockIUserRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockIUserRepository) CreateUser(user model.User) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockIUserRepositoryMockRecorder) CreateUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockIUserRepository)(nil).CreateUser), user)
}

// GetById mocks base method.
func (m *MockIUserRepository) GetById(id uuid.UUID) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIUserRepositoryMockRecorder) GetById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIUserRepository)(nil).GetById), id)
}

// GetUserByUsername mocks base method.
func (m *MockIUserRepository) GetUserByUsername(username string, columns ...postgres.Projection) (*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{username}
	for _, a := range columns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserByUsername", varargs...)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockIUserRepositoryMockRecorder) GetUserByUsername(username any, columns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{username}, columns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockIUserRepository)(nil).GetUserByUsername), varargs...)
}

// GetUserByUsernameAndPassword mocks base method.
func (m *MockIUserRepository) GetUserByUsernameAndPassword(username, password string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsernameAndPassword", username, password)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsernameAndPassword indicates an expected call of GetUserByUsernameAndPassword.
func (mr *MockIUserRepositoryMockRecorder) GetUserByUsernameAndPassword(username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsernameAndPassword", reflect.TypeOf((*MockIUserRepository)(nil).GetUserByUsernameAndPassword), username, password)
}

// UpdatePassword mocks base method.
func (m *MockIUserRepository) UpdatePassword(user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockIUserRepositoryMockRecorder) UpdatePassword(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockIUserRepository)(nil).UpdatePassword), user)
}
