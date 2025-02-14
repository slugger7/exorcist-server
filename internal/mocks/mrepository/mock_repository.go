package mrepository

import (
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

var stackCount = 0

func incStack() int {
	stackCount++
	return stackCount - 1
}

type MockRepo struct {
	MockLibraryRepo
	MockLibraryPathRepo
}

func SetupMockRespository() *MockRepo {
	stackCount = 0
	return &MockRepo{
		MockLibraryRepo:     *SetupMockLibraryRepository(),
		MockLibraryPathRepo: *SetupMockLibraryPathRepository(),
	}
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
func (mr MockRepo) VideoRepo() videoRepository.IVideoRepository {
	panic("not implemented")
}
func (mr MockRepo) UserRepo() userRepository.IUserRepository {
	panic("not implemented")
}
