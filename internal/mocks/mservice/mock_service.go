package mservice

import (
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	libraryPathService "github.com/slugger7/exorcist/internal/service/library_path"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

var stackCount = 0

func incStack() int {
	stackCount++
	return stackCount - 1
}

type MockService struct {
	userService        userService.IUserService
	libraryService     libraryService.ILibraryService
	libraryPathService libraryPathService.ILibraryPathService
}

type MockServices struct {
	LibraryService     MockLibraryService
	UserService        MockUserService
	LibraryPathService *MockLibaryPathService
}

func (ms MockService) User() userService.IUserService {
	return ms.userService
}

func (ms MockService) Library() libraryService.ILibraryService {
	return ms.libraryService
}

func (ms MockService) LibraryPath() libraryPathService.ILibraryPathService {
	return ms.libraryPathService
}

func SetupMockService() (*MockService, *MockServices) {
	stackCount = 0

	mockServices := &MockServices{
		UserService:        SetupMockUserService(),
		LibraryService:     SetupMockLibraryService(),
		LibraryPathService: SetupMockLibraryPathService(),
	}
	ms := &MockService{
		userService:        mockServices.UserService,
		libraryService:     mockServices.LibraryService,
		libraryPathService: mockServices.LibraryPathService,
	}
	return ms, mockServices
}
