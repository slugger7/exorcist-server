package mocks

// Deprecated: moved to mockgen in mock folder
type MockFixture[T any] struct {
	MockModels map[int][]T
	MockModel  map[int]*T
	MockError  map[int]error
}

// Deprecated: moved to mockgen in mock folder
func SetupMockFixture[C any]() *MockFixture[C] {
	models := make(map[int][]C)
	model := make(map[int]*C)
	errors := make(map[int]error)

	return &MockFixture[C]{models, model, errors}
}
