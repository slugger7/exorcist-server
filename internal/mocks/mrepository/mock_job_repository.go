package mrepository

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
)

// Deprecated: moved to mockgen in mock folder
type MockJobRepo mocks.MockFixture[model.Job]

// Deprecated: moved to mockgen in mock folder
func SetupMockJobRepo() *MockJobRepo {
	x := MockJobRepo(*mocks.SetupMockFixture[model.Job]())
	return &x
}

// Deprecated: moved to mockgen in mock folder
func (mr MockRepository) Job() jobRepository.IJobRepository {
	return mr.MockJobRepo
}

// Deprecated: moved to mockgen in mock folder
func (m *MockJobRepo) CreateAll(jobs []model.Job) ([]model.Job, error) {
	stack := incStack()
	return m.MockModels[stack], m.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (m *MockJobRepo) GetNextJob() (*model.Job, error) {
	stack := incStack()
	return m.MockModel[stack], m.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (m *MockJobRepo) UpdateJobStatus(model *model.Job) error {
	stack := incStack()
	return m.MockError[stack]
}
