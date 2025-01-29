package jobRepository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/enum"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository"
)

func FetchNextJob() postgres.SelectStatement {
	statment := table.Job.SELECT(table.Job.AllColumns).
		FROM(table.Job).
		WHERE(table.Job.Status.EQ(enum.JobStatusEnum.NotStarted)).
		ORDER_BY(table.Job.Created.ASC()).
		LIMIT(1)

	repository.DebugCheck(statment)

	return statment
}
