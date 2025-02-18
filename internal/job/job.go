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

				switch job.JobType {
				case model.JobTypeEnum_ScanPath:
					if err := jr.ScanPath(job); err != nil {
						jr.logger.Errorf("Scan path finished with errors", err)
						job.Status = model.JobStatusEnum_Failed
						if erro := jr.repo.Job().UpdateJobStatus(job); erro != nil {
							jr.logger.Errorf("Could not update job status after error. Killing to prevent infinite loop: %v", erro)
							return
						}
					}
				case model.JobTypeEnum_GenerateChecksum:
					if err := jr.GenerateChecksum(job); err != nil {
						panic("not implemented")
					}
				default:
					jr.logger.Errorf("Job of type %v is not implemented", job.JobType)
					job.Status = model.JobStatusEnum_Failed
					errorMessage := `{"error":"can't run job due to no job runner implemented"}`
					job.Data = &errorMessage
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
