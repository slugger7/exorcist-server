package mrepository

import (
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

var stackCount = 0

func incStack() int {
	stackCount = stackCount + 1
	return stackCount - 1
}

type MockRepo struct {
	MockLibraryRepo
}

func SetupMockRespository() *MockRepo {
	stackCount = 0
	return &MockRepo{MockLibraryRepo: *SetupMockLibraryRepository()}
}

func (mr MockRepo) LibraryRepo() libraryRepository.ILibraryRepository {
	return mr.MockLibraryRepo
}

func (mr MockRepo) Health() map[string]string {
	panic("not implemented")
}
func (mr MockRepo) Close() error {
	panic("not implemented")
}
func (mr MockRepo) JobRepo() jobRepository.IJobRepository {
	panic("not implemented")
}
func (mr MockRepo) LibraryPathRepo() libraryPathRepository.ILibraryPathRepository {
	panic("not implemented")
}
func (mr MockRepo) VideoRepo() videoRepository.IVideoRepository {
	panic("not implemented")
}
func (mr MockRepo) UserRepo() userRepository.IUserRepository {
	panic("not implemented")
}
