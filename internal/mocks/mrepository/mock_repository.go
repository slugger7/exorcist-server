package mrepository

var stackCount = 0

func incStack() int {
	stackCount++
	return stackCount - 1
}

// Deprecated: moved to mockgen in mock folder
type MockRepository struct {
	*MockLibraryRepo
	*MockLibraryPathRepo
	*MockUserRepo
	*MockVideoRepo
	*MockJobRepo
}

// Deprecated: moved to mockgen in mock folder
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

// Deprecated: moved to mockgen in mock folder
func (mr MockRepository) Health() map[string]string {
	panic("not implemented")
}

// Deprecated: moved to mockgen in mock folder
func (mr MockRepository) Close() error {
	panic("not implemented")
}
