package mrepository

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
)

type MockJobRepo mocks.MockFixture[model.Job]

func SetupMockJobRepo() *MockJobRepo {
	x := MockJobRepo(*mocks.SetupMockFixture[model.Job]())
	return &x
}

func (mr MockRepository) Job() jobRepository.IJobRepository {
	return mr.MockJobRepo
}

func (m *MockJobRepo) CreateAll(jobs []model.Job) ([]model.Job, error) {
	stack := incStack()
	return m.MockModels[stack], m.MockError[stack]
}

func (m *MockJobRepo) GetNextJob() (*model.Job, error) {
	stack := incStack()
	return m.MockModel[stack], m.MockError[stack]
}
