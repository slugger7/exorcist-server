package jobRepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository/util"
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
	GetAll(dto.JobSearchDTO) (*models.Page[model.Job], error)
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

func (r *JobRepository) GetAll(m dto.JobSearchDTO) (*models.Page[model.Job], error) {
	if m.Limit == 0 {
		m.Limit = 100
	}

	statement := table.Job.SELECT(table.Job.AllColumns).
		FROM(table.Job).
		ORDER_BY(m.OrderBy.ToColumn()).
		LIMIT(int64(m.Limit)).
		OFFSET(int64(m.Skip))

	countStatement := table.Job.SELECT(postgres.COUNT(table.Job.ID).AS("total")).FROM(table.Job)

	var whereExpression postgres.BoolExpression
	if m.Parent == nil {
		whereExpression = table.Job.Parent.IS_NULL()
	} else {
		id, _ := uuid.Parse(*m.Parent)
		whereExpression = table.Job.Parent.EQ(postgres.UUID(id))
	}

	statusExpressions := make([]postgres.Expression, len(m.Statuses))
	for i, s := range m.Statuses {
		statusExpressions[i] = postgres.NewEnumValue(string(s))
	}
	if len(statusExpressions) > 0 {
		whereExpression = whereExpression.AND(table.Job.Status.IN(statusExpressions...))
	}

	jobTypeExpression := make([]postgres.Expression, len(m.JobTypes))
	for i, t := range m.JobTypes {
		jobTypeExpression[i] = postgres.NewEnumValue(string(t))
	}
	if len(jobTypeExpression) > 0 {
		whereExpression = whereExpression.AND(table.Job.JobType.IN(jobTypeExpression...))
	}

	statement = statement.WHERE(whereExpression)
	countStatement = countStatement.WHERE(whereExpression)

	util.DebugCheck(r.Env, statement)
	util.DebugCheck(r.Env, countStatement)

	var totalStruct struct {
		Total int
	}
	if err := countStatement.QueryContext(r.ctx, r.db, &totalStruct); err != nil {
		return nil, errs.BuildError(err, "could not query jobs for total")
	}

	var jobsStruct []struct{ model.Job }
	if err := statement.QueryContext(r.ctx, r.db, &jobsStruct); err != nil {
		return nil, errs.BuildError(err, "could not get jobs with %v", m)
	}

	var jobs = make([]model.Job, len(jobsStruct))
	for i, j := range jobsStruct {
		jobs[i] = j.Job
	}

	return &models.Page[model.Job]{
		Total: totalStruct.Total,
		Limit: m.Limit,
		Skip:  m.Skip,
		Data:  jobs,
	}, nil
}
