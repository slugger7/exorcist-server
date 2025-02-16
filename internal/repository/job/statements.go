package jobRepository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository/util"
)

func (jb *JobRepository) createAllStatement(jobs []model.Job) JobStatement {
	statement := table.JobTable.INSERT(*table.Job, table.Job.JobType, table.Job.Status, table.Job.Data).
		MODELS(jobs).
		RETURNING(table.Job.ID)

	util.DebugCheck(jb.Env, statement)

	return JobStatement{db: jb.db, Statement: statement}
}

func (jb *JobRepository) getNextJobStatement() JobStatement {
	statement := table.Job.SELECT(table.Job.AllColumns).
		FROM(table.Job).
		WHERE(table.Job.Status.EQ(postgres.NewEnumValue(string(model.JobStatusEnum_NotStarted)))).
		ORDER_BY(table.Job.Created.ASC()).
		LIMIT(1)

	util.DebugCheck(jb.Env, statement)

	return JobStatement{statement, jb.db}
}
