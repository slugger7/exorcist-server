package mocks

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	libraryPathService "github.com/slugger7/exorcist/internal/service/library_path"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

var stackCount = 0

type MockService struct {
	userService        userService.IUserService
	libraryService     libraryService.ILibraryService
	libraryPathService libraryPathService.ILibraryPathService
}

type MockUserService struct {
	returningModel *model.User // deprecated
	returningError error       // deprecated
	mockModels     map[int][]model.User
	mockErrors     map[int]error
	mockModel      map[int]*model.User
}

type MockLibraryService struct {
	returningModel *model.Library // deprecated
	returningError error          // deprecated
	mockModels     map[int][]model.Library
	mockErrors     map[int]error
	mockModel      map[int]*model.Library
}

type MockLibaryPathService struct {
	mockModels map[int][]model.LibraryPath
	mockErrors map[int]error
}

type MockServices struct {
	libraryService     MockLibraryService
	userService        MockUserService
	libraryPathService MockLibaryPathService
}

func (ms MockService) UserService() userService.IUserService {
	return ms.userService
}

func (ms MockService) LibraryService() libraryService.ILibraryService {
	return ms.libraryService
}

func (ms MockService) LibraryPathService() libraryPathService.ILibraryPathService {
	return ms.libraryPathService
}

func (mus MockUserService) CreateUser(username, password string) (*model.User, error) {
	return mus.returningModel, mus.returningError
}

func (mus MockUserService) ValidateUser(username, password string) (*model.User, error) {
	return mus.returningModel, mus.returningError
}

func (ls MockLibraryService) CreateLibrary(actual model.Library) (*model.Library, error) {
	return ls.returningModel, ls.returningError
}

func (ls MockLibraryService) GetLibraries() ([]model.Library, error) {
	stackCount = stackCount + 1
	return ls.mockModels[stackCount-1], ls.mockErrors[stackCount-1]
}

func (ls MockLibraryService) GetLibraryById(uuid.UUID) (*model.Library, error) {
	panic("not implemented")
}

func (lps MockLibaryPathService) Create(*model.LibraryPath) (*model.LibraryPath, error) {
	panic("not implemented")
}

func setupMockUserService() MockUserService {
	mockModels := make(map[int][]model.User)
	mockErrors := make(map[int]error)
	mockModel := make(map[int]*model.User)
	return MockUserService{mockModels: mockModels, mockErrors: mockErrors, mockModel: mockModel}
}

func setupMockLibraryService() MockLibraryService {
	mockModels := make(map[int][]model.Library)
	mockErrors := make(map[int]error)
	mockModel := make(map[int]*model.Library)
	return MockLibraryService{mockModels: mockModels, mockErrors: mockErrors, mockModel: mockModel}
}
func setupMockLibraryPathService() MockLibaryPathService {
	mockModels := make(map[int][]model.LibraryPath)
	mockErrors := make(map[int]error)
	return MockLibaryPathService{mockModels, mockErrors}
}

func SetupMockService() (*MockService, *MockServices) {
	stackCount = 0

	mockServices := &mockServices{
		userService:        setupMockUserService(),
		libraryService:     setupMockLibraryService(),
		libraryPathService: setupMockLibraryPathService(),
	}
	ms := &MockService{
		userService:        mockServices.userService,
		libraryService:     mockServices.libraryService,
		libraryPathService: mockServices.libraryPathService,
	}
	return ms, mockServices
}
