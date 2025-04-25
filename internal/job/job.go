package job

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
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

			jobFunc, err := jr.jobFuncResolver(job.JobType)
			if err != nil {
				job.Status = model.JobStatusEnum_Cancelled
				errorMessage := jr.marshallJobError(err.Error())
				job.Outcome = &errorMessage
				if err := jr.repo.Job().UpdateJobStatus(job); err != nil {
					return errs.BuildError(err, "Could not update not implemented job %v. Killing to prevent infinite loop", job.JobType)
				}
			}

			if err := jobFunc(job); err != nil {
				jr.logger.Errorf("Job finished with errors: %v", err.Error())
				job.Status = model.JobStatusEnum_Failed
				errText := jr.marshallJobError(err.Error())
				job.Outcome = &errText
				if erro := jr.repo.Job().UpdateJobStatus(job); erro != nil {
					return errs.BuildError(erro, "Could not update job status after error. Killing to prevent infinite loop")
				}
			}

			job.Status = model.JobStatusEnum_Completed
			if err := jr.repo.Job().UpdateJobStatus(job); err != nil {
				return errs.BuildError(err, "Could not update job status after success. Killing to prevent infinite loop")
			}
		}
	}
}

type JobFunc func(*model.Job) error

func (jr *JobRunner) jobFuncResolver(jobType model.JobTypeEnum) (JobFunc, error) {
	var f JobFunc
	switch jobType {
	case model.JobTypeEnum_ScanPath:
		f = func(j *model.Job) error {
			return jr.ScanPath(j)
		}
	case model.JobTypeEnum_GenerateChecksum:
		f = func(j *model.Job) error {
			return jr.GenerateChecksum(j)
		}
	case model.JobTypeEnum_GenerateThumbnail:
		f = func(j *model.Job) error {
			return jr.GenerateThumbnail(j)
		}
	default:
		return nil, fmt.Errorf("no implementation to run job type %v", jobType)
	}
	return f, nil
}

func (jr *JobRunner) marshallJobError(e string) string {
	data, err := json.Marshal(models.JobError{
		Error: e,
	})
	if err != nil {
		jr.logger.Errorf("Could not marshall erorr: %v", err.Error())
		return "could not marshall error. check logs"
	}
	return string(data)
}
