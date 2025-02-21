// Deprecated: use mocks in mock folder instead
package mservice

import (
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	libraryPathService "github.com/slugger7/exorcist/internal/service/library_path"
	userService "github.com/slugger7/exorcist/internal/service/user"
	videoService "github.com/slugger7/exorcist/internal/service/video"
)

var stackCount = 0

func incStack() int {
	stackCount++
	return stackCount - 1
}

// Deprecated: moved to mockgen in mock folder
type MockService struct {
	user        userService.IUserService
	library     libraryService.ILibraryService
	libraryPath libraryPathService.ILibraryPathService
	video       videoService.IVideoService
}

// Deprecated: moved to mockgen in mock folder
type MockServices struct {
	Library     *MockLibraryService
	User        *MockUserService
	LibraryPath *MockLibaryPathService
	Video       *MockVideoService
}

// Deprecated: moved to mockgen in mock folder
func SetupMockService() (*MockService, *MockServices) {
	stackCount = 0

	mockServices := &MockServices{
		User:        SetupMockUserService(),
		Library:     SetupMockLibraryService(),
		LibraryPath: SetupMockLibraryPathService(),
		Video:       SetupMockVideoService(),
	}
	ms := &MockService{
		user:        mockServices.User,
		library:     mockServices.Library,
		libraryPath: mockServices.LibraryPath,
		video:       mockServices.Video,
	}
	return ms, mockServices
}
