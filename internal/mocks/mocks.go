package mocks

type MockFixture[T any] struct {
	MockModels map[int][]T
	MockModel  map[int]*T
	MockError  map[int]error
}

func SetupMockFixture[C any]() *MockFixture[C] {
	models := make(map[int][]C)
	model := make(map[int]*C)
	errors := make(map[int]error)

	return &MockFixture[C]{models, model, errors}
}
