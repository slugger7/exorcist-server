package jobRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type JobStatement struct {
	postgres.Statement
	db *sql.DB
}

type IJobRepository interface {
	CreateAll(jobs []model.Job) ([]model.Job, error)
}

type JobRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var jobRepoInstance *JobRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) IJobRepository {
	if jobRepoInstance != nil {
		return jobRepoInstance
	}
	jobRepoInstance = &JobRepository{
		db:  db,
		Env: env,
	}
	return jobRepoInstance
}

func (js JobStatement) Query(destination interface{}) error {
	return js.Statement.Query(js.db, destination)
}

func (j *JobRepository) CreateAll(jobs []model.Job) ([]model.Job, error) {
	if len(jobs) == 0 {
		return jobs, nil
	}
	var newJobs []struct{ model.Job }
	if err := j.createAllStatement(jobs).Query(&newJobs); err != nil {
		return nil, errs.BuildError(err, "error when creating jobs")
	}

	jobModels := []model.Job{}
	for _, j := range newJobs {
		jobModels = append(jobModels, j.Job)
	}

	return jobModels, nil
}
