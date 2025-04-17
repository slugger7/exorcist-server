package job

import (
	"context"
	"fmt"
	"sync"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
	"github.com/slugger7/exorcist/internal/service"
)

type JobRunner struct {
	env         *environment.EnvironmentVariables
	service     service.IService
	repo        repository.IRepository
	logger      logger.ILogger
	ch          chan bool
	shutdownCtx context.Context
	wg          *sync.WaitGroup
}

var jobRunnerInstance *JobRunner

func New(
	env *environment.EnvironmentVariables,
	serv service.IService,
	repo repository.IRepository,
	logger logger.ILogger,
	shutdownCtx context.Context,
	wg *sync.WaitGroup,
) chan bool {
	ch := make(chan bool)
	if jobRunnerInstance == nil {
		jobRunnerInstance = &JobRunner{
			env:         env,
			service:     serv,
			repo:        repo,
			logger:      logger,
			ch:          ch,
			wg:          wg,
			shutdownCtx: shutdownCtx,
		}

		wg.Add(1)
		go jobRunnerInstance.loop()
	}

	return ch
}

func (jr *JobRunner) loop() {
	defer jr.wg.Done()

	jr.logger.Infof("Running jobs")
	for {
		select {
		case <-jr.shutdownCtx.Done():
			jr.logger.Debug("Shutdown signal received. Shutting down")
			return
		case _, ok := <-jr.ch:
			if !ok {
				jr.logger.Debug("Channel closed. stopping loop")
				return
			}

			jr.logger.Info("Processing jobs")
			if err := jr.processJobs(); err != nil {
				jr.logger.Errorf("Error received while processing jobs. Stopping job runner", err.Error())
				return
			}
		}
	}
}

func (jr *JobRunner) processJobs() error {
	for {
		select {
		case <-jr.shutdownCtx.Done():
			return fmt.Errorf("shutdown signal received. Stopping job loop")
		default:
			job, err := jr.repo.Job().GetNextJob()
			if err != nil {
				return errs.BuildError(err, "Failed to fetch next job")
			}
			if job == nil {
				jr.logger.Info("No jobs to run. Waiting for next signal")
				return nil
			}

			job.Status = model.JobStatusEnum_InProgress
			if err := jr.repo.Job().UpdateJobStatus(job); err != nil {
				return errs.BuildError(err, "Failed to update job status")
			}

			switch job.JobType {
			case model.JobTypeEnum_ScanPath:
				if err := jr.ScanPath(job); err != nil {
					jr.logger.Errorf("Scan path finished with errors", err.Error())
					job.Status = model.JobStatusEnum_Failed
					if erro := jr.repo.Job().UpdateJobStatus(job); erro != nil {
						return errs.BuildError(erro, "Could not update job status after error. Killing to prevent infinite loop")
					}
				}
			case model.JobTypeEnum_GenerateChecksum:
				if err := jr.GenerateChecksum(job); err != nil {
					jr.logger.Errorf("Generate checksum finished with errors: %v", err.Error())
					job.Status = model.JobStatusEnum_Failed
					if erro := jr.repo.Job().UpdateJobStatus(job); erro != nil {
						return errs.BuildError(erro, "Could not update job status after error. Killing to prevent infinite loop")
					}
				}
			case model.JobTypeEnum_GenerateThumbnail:
				if err := jr.GenerateThumbnail(job); err != nil {
					jr.logger.Errorf("Generate thumbnail finished with errors: %v", err.Error())
					job.Status = model.JobStatusEnum_Failed
					if erro := jr.repo.Job().UpdateJobStatus(job); erro != nil {
						return errs.BuildError(erro, "Could not update job status after error. Killing to prevent infinite loop")
					}
				}
			default:
				jr.logger.Errorf("Job of type %v is not implemented", job.JobType)
				job.Status = model.JobStatusEnum_Cancelled
				errorMessage := `{"error":"can't run job due to no job runner implemented"}`
				job.Outcome = &errorMessage
				if err := jr.repo.Job().UpdateJobStatus(job); err != nil {
					return errs.BuildError(err, "Could not update not implemented job %v. Killing to prevent infinite loop", job.JobType)
				}
			}

			job.Status = model.JobStatusEnum_Completed
			if err := jr.repo.Job().UpdateJobStatus(job); err != nil {
				return errs.BuildError(err, "Could not update job status after success. Killing to prevent infinite loop")
			}
		}
	}
}
