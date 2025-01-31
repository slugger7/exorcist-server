package repository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/enum"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
)

func (s *DatabaseService) FetchNextJob() postgres.SelectStatement {
	statment := table.Job.SELECT(table.Job.AllColumns).
		FROM(table.Job).
		WHERE(table.Job.Status.EQ(enum.JobStatusEnum.NotStarted)).
		ORDER_BY(table.Job.Created.ASC()).
		LIMIT(1)

	s.DebugCheck(statment)

	return statment
}
