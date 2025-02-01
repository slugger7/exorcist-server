package jobRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/enum"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type JobStatement struct {
	postgres.Statement
	db *sql.DB
}

type IJobRepository interface {
	FetchNextJob() JobStatement
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

func (s *JobRepository) FetchNextJob() JobStatement {
	statment := table.Job.SELECT(table.Job.AllColumns).
		FROM(table.Job).
		WHERE(table.Job.Status.EQ(enum.JobStatusEnum.NotStarted)).
		ORDER_BY(table.Job.Created.ASC()).
		LIMIT(1)

	util.DebugCheck(s.Env, statment)

	return JobStatement{statment, s.db}
}

func (js JobStatement) Query(destination interface{}) error {
	return js.Statement.Query(js.db, destination)
}
