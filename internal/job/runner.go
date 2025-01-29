package job

import (
	"database/sql"
	"log"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
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
		job.Status = model.JobStatusEnum_InProgress
		err := updateJobStatus(db, *job)
		if err != nil {
			log.Printf("Failed to update job status to %v: %v", job.Status, err)
			// try to add the error message to the job (there might be something else wrong and this might fail too.)
			break
		}

		switch job.JobType {
		case model.JobTypeEnum_GenerateChecksum:

		}
	}
}

func fetchNextJob(db *sql.DB) *model.Job {
	var jobs []struct {
		model.Job
	}
	err := jobRepository.FetchNextJob().
		Query(db, &jobs)
	if err != nil {
		log.Printf("Could not fetch jobs: %v", err)
		return nil
	}

	return &jobs[len(jobs)-1].Job
}

func updateJobStatus(db *sql.DB, job model.Job) error {
	_, err := table.Job.UPDATE(table.Job.Status).
		SET(table.Job.Status.SET(postgres.NewEnumValue(string(job.Status)))).
		WHERE(table.Job.ID.EQ(postgres.UUID(job.ID))).
		Exec(db)
	return err
}
