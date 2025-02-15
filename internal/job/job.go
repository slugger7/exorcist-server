package job

import (
	"sync"

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
	for {
		select {
		case val, ok := <-jr.ch:
			if !ok {
				// Cleanup methods can be run from here
				jr.wg.Done()
				return
			}

			_ = val

			// run the next job
			// when starting a job it might be worth pushing to the channel from a goroutine
		}
	}
}
