package jobRepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type JobStatement struct {
	postgres.Statement
	db  *sql.DB
	ctx context.Context
}

type IJobRepository interface {
	CreateAll(jobs []model.Job) ([]model.Job, error)
	GetNextJob() (*model.Job, error)
	UpdateJobStatus(model *model.Job) error
}

type JobRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
	ctx context.Context
}

var jobRepoInstance *JobRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) IJobRepository {
	if jobRepoInstance != nil {
		return jobRepoInstance
	}
	jobRepoInstance = &JobRepository{
		db:  db,
		Env: env,
		ctx: context,
	}
	return jobRepoInstance
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

func (j *JobRepository) GetNextJob() (*model.Job, error) {
	var job []struct{ model.Job }
	if err := j.getNextJobStatement().Query(&job); err != nil {
		return nil, errs.BuildError(err, "could not get next job")
	}
	if len(job) == 1 {
		return &job[len(job)-1].Job, nil
	}

	return nil, nil
}

func (j *JobRepository) UpdateJobStatus(model *model.Job) error {
	model.Modified = time.Now()
	if _, err := j.updateJobStatusStatement(model).Exec(); err != nil {
		return errs.BuildError(err, "could not update job %v status to %v", model.ID, model.Status)
	}

	return nil
}
