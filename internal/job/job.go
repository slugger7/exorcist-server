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
					jr.logger.Errorf("Could not get the next job: %v", err)
				}
				if job == nil {
					jr.logger.Info("No jobs to run. Waiting for next signal")
					break
				}

				switch job.JobType {
				case model.JobTypeEnum_ScanPath:
					jr.ScanPath(job)
				default:
					jr.logger.Errorf("Job of type %v is not implemented", job.JobType)
				}
			}
		}
	}
}
