package mrepository

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
	*MockJobRepo
}

func SetupMockRespository() *MockRepository {
	stackCount = 0
	return &MockRepository{
		MockLibraryRepo:     SetupMockLibraryRepo(),
		MockLibraryPathRepo: SetupMockLibraryPathRepository(),
		MockUserRepo:        SetupMockUserRepository(),
		MockVideoRepo:       SetupMockVideoRepository(),
		MockJobRepo:         SetupMockJobRepo(),
	}
}

func (mr MockRepository) Health() map[string]string {
	panic("not implemented")
}
func (mr MockRepository) Close() error {
	panic("not implemented")
}
