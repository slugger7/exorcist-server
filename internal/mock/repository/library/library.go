// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/library/library.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/library/library.go
//

// Package mock_libraryRepository is a generated GoMock package.
package mock_libraryRepository

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	model "github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	dto "github.com/slugger7/exorcist/internal/dto"
	models "github.com/slugger7/exorcist/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockLibraryRepository is a mock of LibraryRepository interface.
type MockLibraryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLibraryRepositoryMockRecorder
	isgomock struct{}
}

// MockLibraryRepositoryMockRecorder is the mock recorder for MockLibraryRepository.
type MockLibraryRepositoryMockRecorder struct {
	mock *MockLibraryRepository
}

// NewMockLibraryRepository creates a new mock instance.
func NewMockLibraryRepository(ctrl *gomock.Controller) *MockLibraryRepository {
	mock := &MockLibraryRepository{ctrl: ctrl}
	mock.recorder = &MockLibraryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLibraryRepository) EXPECT() *MockLibraryRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockLibraryRepository) Create(name string) (*model.Library, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name)
	ret0, _ := ret[0].(*model.Library)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockLibraryRepositoryMockRecorder) Create(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLibraryRepository)(nil).Create), name)
}

// GetAll mocks base method.
func (m *MockLibraryRepository) GetAll() ([]model.Library, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]model.Library)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockLibraryRepositoryMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockLibraryRepository)(nil).GetAll))
}

// GetById mocks base method.
func (m *MockLibraryRepository) GetById(arg0 uuid.UUID) (*model.Library, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", arg0)
	ret0, _ := ret[0].(*model.Library)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockLibraryRepositoryMockRecorder) GetById(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockLibraryRepository)(nil).GetById), arg0)
}

// GetByName mocks base method.
func (m *MockLibraryRepository) GetByName(name string) (*model.Library, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", name)
	ret0, _ := ret[0].(*model.Library)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockLibraryRepositoryMockRecorder) GetByName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockLibraryRepository)(nil).GetByName), name)
}

// GetMedia mocks base method.
func (m *MockLibraryRepository) GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMedia", id, userId, search)
	ret0, _ := ret[0].(*dto.PageDTO[models.MediaOverviewModel])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMedia indicates an expected call of GetMedia.
func (mr *MockLibraryRepositoryMockRecorder) GetMedia(id, userId, search any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMedia", reflect.TypeOf((*MockLibraryRepository)(nil).GetMedia), id, userId, search)
}

// Update mocks base method.
func (m_2 *MockLibraryRepository) Update(m model.Library) (*model.Library, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", m)
	ret0, _ := ret[0].(*model.Library)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockLibraryRepositoryMockRecorder) Update(m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockLibraryRepository)(nil).Update), m)
}
