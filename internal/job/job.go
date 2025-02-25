package job

import (
	"sync"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
	"github.com/slugger7/exorcist/internal/service"
)

type JobRunner struct {
	env     *environment.EnvironmentVariables
	service service.IService
	repo    repository.IRepository
	logger  logger.ILogger
	ch      chan bool
	wg      *sync.WaitGroup
}

var jobRunnerInstance *JobRunner

func New(
	env *environment.EnvironmentVariables,
	serv service.IService,
	repo repository.IRepository,
	logger logger.ILogger,
	wg *sync.WaitGroup,
) chan bool {
	ch := make(chan bool)
	if jobRunnerInstance == nil {
		jobRunnerInstance = &JobRunner{
			env:     env,
			service: serv,
			repo:    repo,
			logger:  logger,
			ch:      ch,
			wg:      wg,
		}

		wg.Add(1)
		go jobRunnerInstance.loop()
	}

	return ch
}

func (jr *JobRunner) loop() {
	jr.logger.Infof("Running jobs")
	for {
		select {
		case _, ok := <-jr.ch:
			if !ok {
				// Cleanup methods can be run from here
				jr.wg.Done()
				return
			}

			jr.logger.Info("Checking for jobs")
			for {
				job, err := jr.repo.Job().GetNextJob()
				if err != nil {
					jr.logger.Errorf("Could not get the next job: %v", err.Error())
				}
				if job == nil {
					jr.logger.Info("No jobs to run. Waiting for next signal")
					break
				}

				job.Status = model.JobStatusEnum_InProgress
				if err := jr.repo.Job().UpdateJobStatus(job); err != nil {
					jr.logger.Errorf("could not update job (%v) status to in progress: %v", job.ID, err.Error())
					return
				}

				switch job.JobType {
				case model.JobTypeEnum_ScanPath:
					if err := jr.ScanPath(job); err != nil {
						jr.logger.Errorf("Scan path finished with errors", err.Error())
						job.Status = model.JobStatusEnum_Failed
						if erro := jr.repo.Job().UpdateJobStatus(job); erro != nil {
							jr.logger.Errorf("Could not update job status after error. Killing to prevent infinite loop: %v", erro.Error())
							return
						}
					}
				case model.JobTypeEnum_GenerateChecksum:
					if err := jr.GenerateChecksum(job); err != nil {
						jr.logger.Errorf("Generate checksum finished with errors: %v", err.Error())
						job.Status = model.JobStatusEnum_Failed
						if erro := jr.repo.Job().UpdateJobStatus(job); erro != nil {
							jr.logger.Errorf("Could not update job status after error. Killing to prevent infinite loop: %v", erro.Error())
							return
						}
					}
				case model.JobTypeEnum_GenerateThumbnail:
					if err := jr.GenerateThumbnail(job); err != nil {
						jr.logger.Errorf("Generate thumbnail finished with errors: %v", err.Error())
						job.Status = model.JobStatusEnum_Failed
						if erro := jr.repo.Job().UpdateJobStatus(job); erro != nil {
							jr.logger.Errorf("Could not update job status after error. Killing to prevent infinite loop: %v", erro.Error())
						}
					}
				default:
					jr.logger.Errorf("Job of type %v is not implemented", job.JobType)
					job.Status = model.JobStatusEnum_Cancelled
					errorMessage := `{"error":"can't run job due to no job runner implemented"}`
					job.Outcome = &errorMessage
					if err := jr.repo.Job().UpdateJobStatus(job); err != nil {
						jr.logger.Errorf("Could not update not implemented job %v. Killing to prevent infinite loop: %v", job.JobType, err.Error())
						return
					}
				}

				job.Status = model.JobStatusEnum_Completed
				if err := jr.repo.Job().UpdateJobStatus(job); err != nil {
					jr.logger.Errorf("Could not update job status after success. Killing to prevent infinite loop: %v", err)
				}
			}
		}
	}
}
