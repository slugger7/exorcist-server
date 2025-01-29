package job

import (
	"database/sql"
	"log"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/enum"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
)

type Service struct {
	Running bool
}

var service *Service

func updateRunning(value bool) {
	service.Running = value
}

func RunJobs(db *sql.DB) {
	if service == nil {
		service = &Service{
			Running: false,
		}
	}
	if service.Running {
		return
	}

	service.Running = true
	defer updateRunning(false)

	for {
		job := fetchNextJob(db)
		if job == nil {
			break
		}
	}
}

func fetchNextJob(db *sql.DB) *model.Job {
	var jobs []struct {
		model.Job
	}
	err := table.Job.SELECT(table.Job.AllColumns).
		FROM(table.Job).
		WHERE(table.Job.Status.EQ(enum.JobStatusEnum.NotStarted)).
		ORDER_BY(table.Job.Created.ASC()).
		LIMIT(1).
		Query(db, &jobs)
	if err != nil {
		log.Printf("Could not fetch jobs: %v", err)
		return nil
	}

	return &jobs[len(jobs)-1].Job
}
