package mrepository

import (
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
)

var stackCount = 0

func incStack() int {
	stackCount++
	return stackCount - 1
}

type MockRepository struct {
	*MockLibraryRepo
	*MockLibraryPathRepo
	*MockUserRepo
	*MockVideoRepo
}

func SetupMockRespository() *MockRepository {
	stackCount = 0
	return &MockRepository{
		MockLibraryRepo:     SetupMockLibraryRepo(),
		MockLibraryPathRepo: SetupMockLibraryPathRepository(),
		MockUserRepo:        SetupMockUserRepository(),
		MockVideoRepo:       SetupMockVideoRepository(),
	}
}

func (mr MockRepository) Health() map[string]string {
	panic("not implemented")
}
func (mr MockRepository) Close() error {
	panic("not implemented")
}
func (mr MockRepository) Job() jobRepository.IJobRepository {
	panic("not implemented")
}
