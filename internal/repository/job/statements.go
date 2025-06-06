package jobRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository/util"
)

func (js JobStatement) Query(destination interface{}) error {
	return js.Statement.QueryContext(js.ctx, js.db, destination)
}

func (js JobStatement) Exec() (sql.Result, error) {
	return js.Statement.ExecContext(js.ctx, js.db)
}

func (jb *jobRepository) createAllStatement(jobs []model.Job) JobStatement {
	statement := table.JobTable.INSERT(*table.Job, table.Job.JobType, table.Job.Status, table.Job.Data, table.Job.Parent, table.Job.Priority).
		MODELS(jobs).
		RETURNING(table.Job.AllColumns)

	util.DebugCheck(jb.env, statement)

	return JobStatement{db: jb.db, Statement: statement}
}

func (jb *jobRepository) getNextJobStatement() JobStatement {
	statement := table.Job.SELECT(table.Job.AllColumns).
		FROM(table.Job).
		WHERE(table.Job.Status.EQ(postgres.NewEnumValue(string(model.JobStatusEnum_NotStarted)))).
		ORDER_BY(table.Job.Priority.ASC(), table.Job.Created.ASC()).
		LIMIT(1)

	util.DebugCheck(jb.env, statement)

	return JobStatement{statement, jb.db, jb.ctx}
}

func (jb *jobRepository) updateJobStatusStatement(model *model.Job) JobStatement {
	statement := table.Job.UPDATE(table.Job.Modified, table.Job.Status, table.Job.Outcome).
		MODEL(model).
		WHERE(table.Job.ID.EQ(postgres.UUID(model.ID)))

	util.DebugCheck(jb.env, statement)

	return JobStatement{statement, jb.db, jb.ctx}
}
